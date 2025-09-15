package server

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

// Server struct kept as before
type Server struct {
	DB     *gorm.DB
	Config Config
	WSHub  *WSHub
}

func Run(cfg Config) error {
	db, err := NewGorm(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	s := &Server{
		DB:     db,
		Config: cfg,
		WSHub:  NewWSHub(),
	}

	// seed default boards
	seedDefaultBoards(db)

	go s.WSHub.Run()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:9876",
			"http://127.0.0.1:9876",
			"http://localhost:5173",
			"http://127.0.0.1:5173",
		},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	// attach server instance
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("server", s)
			return next(c)
		}
	})

	api := e.Group("/api/v1")
	api.GET("/health", func(c echo.Context) error { return c.JSON(http.StatusOK, echo.Map{"status": "ok"}) })

	// auth
	api.POST("/auth/register", s.Register)
	api.POST("/auth/login", s.Login)

	// boards
	api.GET("/boards", s.ListBoards)
	api.POST("/boards", s.AuthMiddleware(s.CreateBoard))

	// posts
	api.GET("/boards/:slug/posts", s.ListPosts)
	api.POST("/boards/:slug/posts", s.AuthMiddleware(s.CreatePost))

	// comments
	api.GET("/posts/:id/comments", s.ListComments)
	api.POST("/posts/:id/comments", s.AuthMiddleware(s.CreateComment))

	// threads
	api.GET("/boards/:slug/threads", s.ListThreadsForBoard)
	api.POST("/boards/:slug/threads", s.AuthMiddleware(s.CreateThread))
	api.GET("/threads/:id", s.GetThread)
	api.GET("/threads/:id/posts", s.ListPostsForThread)
	api.POST("/threads/:id/posts", s.AuthMiddleware(s.CreatePostInThread))

	// votes
	api.POST("/posts/:id/vote", s.AuthMiddleware(s.VotePost))
	api.POST("/comments/:id/vote", s.AuthMiddleware(s.VoteComment))
	api.POST("/threads/:id/vote", s.AuthMiddleware(s.VoteThread))

	// ws (keep at root)
	e.GET("/ws", func(c echo.Context) error {
		s.WSHub.ServeWS(c.Response(), c.Request())
		return nil
	})

	return e.Start(":" + cfg.Port)
}

// seedDefaultBoards creates basic boards if they don't exist
func seedDefaultBoards(db *gorm.DB) {
	defaults := []Board{
		{Name: "Random", Slug: "b", Description: "Random board", IsDefault: true},
		{Name: "Technology", Slug: "tech", Description: "Technology board", IsDefault: true},
		{Name: "Art", Slug: "art", Description: "Art and creative", IsDefault: true},
	}
	for _, b := range defaults {
		var exists Board
		if err := db.Where("slug = ?", b.Slug).First(&exists).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&b)
			}
		}
	}
}

// --- Helpers to obtain server & user from context ---
func serverFromContext(c echo.Context) *Server {
	return c.Get("server").(*Server)
}

func userFromClaims(c echo.Context) (*User, error) {
	token := c.Request().Header.Get("Authorization")
	s := serverFromContext(c)
	claims, err := ParseJWT(s.Config.JWTSecret, token)
	if err != nil {
		return nil, err
	}
	var user User
	if err := s.DB.First(&user, "id = ?", claims.UserID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// --- Handlers (concise implementations) ---

// Register expects: { "username","email","password" }
func (s *Server) Register(c echo.Context) error {
	type req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid"})
	}
	hash, err := HashPassword(r.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "hash failed"})
	}
	user := User{
		Username:     strings.TrimSpace(r.Username),
		Email:        strings.TrimSpace(r.Email),
		PasswordHash: hash,
		Role:         "user",
	}
	if err := s.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "user already exists or invalid"})
	}
	token, _ := GenerateJWT(s.Config.JWTSecret, &user)
	return c.JSON(http.StatusCreated, echo.Map{"token": token, "user": echo.Map{"id": user.ID, "username": user.Username}})
}

// Login expects: { "usernameOrEmail", "password" }
func (s *Server) Login(c echo.Context) error {
	type req struct {
		UsernameOrEmail string `json:"usernameOrEmail"`
		Password        string `json:"password"`
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid"})
	}
	var user User
	if err := s.DB.Where("username = ? OR email = ?", r.UsernameOrEmail, r.UsernameOrEmail).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid credentials"})
	}
	ok, err := VerifyPassword(r.Password, user.PasswordHash)
	if err != nil || !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid credentials"})
	}
	token, _ := GenerateJWT(s.Config.JWTSecret, &user)
	return c.JSON(http.StatusOK, echo.Map{"token": token, "user": echo.Map{"id": user.ID, "username": user.Username}})
}

func (s *Server) ListBoards(c echo.Context) error {
	var boards []Board
	s.DB.Order("created_at desc").Find(&boards)
	return c.JSON(http.StatusOK, boards)
}

func (s *Server) CreateBoard(c echo.Context) error {
	type req struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid"})
	}
	user, err := userFromClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauth"})
	}
	board := Board{
		Name:        r.Name,
		Slug:        r.Slug,
		Description: r.Description,
		CreatedByID: &user.ID,
	}
	if err := s.DB.Create(&board).Error; err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "failed to create board"})
	}
	return c.JSON(http.StatusCreated, board)
}

func (s *Server) ListPosts(c echo.Context) error {
	slug := c.Param("slug")
	var board Board
	if err := s.DB.Where("slug = ?", slug).First(&board).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "board not found"})
	}
	var posts []Post
	s.DB.Where("board_id = ?", board.ID).Order("created_at desc").Find(&posts)
	return c.JSON(http.StatusOK, posts)
}

func (s *Server) VotePost(c echo.Context) error {
	type req struct {
		Value int8 `json:"value"` // 1 or -1
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid"})
	}
	if r.Value != 1 && r.Value != -1 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "value must be 1 or -1"})
	}
	postIDStr := c.Param("id")
	pid, err := uuid.Parse(postIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	user, err := userFromClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauth"})
	}
	// upsert vote
	var v Vote
	if err := s.DB.Where("user_id = ? AND target_id = ? AND target_type = ?", user.ID, pid, "post").First(&v).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			nv := Vote{
				UserID:     user.ID,
				TargetID:   pid,
				TargetType: "post",
				Value:      r.Value,
			}
			s.DB.Create(&nv)
		} else {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "db error"})
		}
	} else {
		// update
		v.Value = r.Value
		s.DB.Save(&v)
	}
	// recalc tally (simple)
	var likes int64
	var dislikes int64
	s.DB.Model(&Vote{}).Where("target_id = ? AND target_type = ? AND value = 1", pid, "post").Count(&likes)
	s.DB.Model(&Vote{}).Where("target_id = ? AND target_type = ? AND value = -1", pid, "post").Count(&dislikes)
	s.DB.Model(&Post{}).Where("id = ?", pid).Updates(map[string]interface{}{"likes": likes, "dislikes": dislikes})
	return c.JSON(http.StatusOK, echo.Map{"likes": likes, "dislikes": dislikes})
}

// CreateThread - create a thread and an initial post in that thread
func (s *Server) CreateThread(c echo.Context) error {
	type req struct {
		Title     string `json:"title"`
		PostTitle string `json:"post_title"`
		Content   string `json:"content"`
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid"})
	}
	user, err := userFromClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauth"})
	}
	slug := c.Param("slug")
	var board Board
	if err := s.DB.Where("slug = ?", slug).First(&board).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "board not found"})
	}
	thread := Thread{
		BoardID: board.ID,
		UserID:  user.ID,
		Title:   r.Title,
	}
	if err := s.DB.Create(&thread).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create thread"})
	}

	// create initial post in thread
	html, _ := RenderMarkdown(r.Content)
	post := Post{
		BoardID:     board.ID,
		ThreadID:    &thread.ID,
		UserID:      user.ID,
		Title:       r.PostTitle,
		Content:     r.Content,
		ContentHTML: html,
	}
	if err := s.DB.Create(&post).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create post"})
	}

	// broadcast events
	s.WSHub.broadcast <- WSMessage{Type: "thread.created", Data: thread}
	s.WSHub.broadcast <- WSMessage{Type: "post.created", Data: post}

	return c.JSON(http.StatusCreated, echo.Map{"thread": thread, "post": post})
}

func (s *Server) ListThreadsForBoard(c echo.Context) error {
	slug := c.Param("slug")
	var board Board
	if err := s.DB.Where("slug = ?", slug).First(&board).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "board not found"})
	}
	var threads []Thread
	s.DB.Where("board_id = ?", board.ID).Order("created_at desc").Find(&threads)
	return c.JSON(http.StatusOK, threads)
}

func (s *Server) GetThread(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var t Thread
	if err := s.DB.First(&t, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "thread not found"})
	}
	return c.JSON(http.StatusOK, t)
}

func (s *Server) ListPostsForThread(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var posts []Post
	s.DB.Where("thread_id = ?", id).Order("created_at asc").Find(&posts)
	return c.JSON(http.StatusOK, posts)
}

// CreatePost (enhanced to support optional thread_id sent by frontend)
func (s *Server) CreatePost(c echo.Context) error {
	type req struct {
		Title    string  `json:"title"`
		Content  string  `json:"content"`
		ThreadID *string `json:"thread_id,omitempty"`
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid"})
	}
	user, err := userFromClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauth"})
	}
	slug := c.Param("slug")
	var board Board
	if err := s.DB.Where("slug = ?", slug).First(&board).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "board not found"})
	}

	var threadID *uuid.UUID
	if r.ThreadID != nil && *r.ThreadID != "" {
		tid, err := uuid.Parse(*r.ThreadID)
		if err == nil {
			// verify thread exists and is in same board
			var th Thread
			if err := s.DB.First(&th, "id = ? AND board_id = ?", tid, board.ID).Error; err == nil {
				threadID = &tid
			}
		}
	}

	html, _ := RenderMarkdown(r.Content)
	post := Post{
		BoardID:     board.ID,
		ThreadID:    threadID,
		UserID:      user.ID,
		Title:       r.Title,
		Content:     r.Content,
		ContentHTML: html,
	}
	if err := s.DB.Create(&post).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create post"})
	}
	s.WSHub.broadcast <- WSMessage{Type: "post.created", Data: post}
	return c.JSON(http.StatusCreated, post)
}

// CreatePostInThread receives posts created directly in a thread URL
func (s *Server) CreatePostInThread(c echo.Context) error {
	type req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid"})
	}
	user, err := userFromClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauth"})
	}
	thIDStr := c.Param("id")
	tid, err := uuid.Parse(thIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var th Thread
	if err := s.DB.First(&th, "id = ?", tid).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "thread not found"})
	}

	html, _ := RenderMarkdown(r.Content)
	post := Post{
		BoardID:     th.BoardID,
		ThreadID:    &th.ID,
		UserID:      user.ID,
		Title:       r.Title,
		Content:     r.Content,
		ContentHTML: html,
	}
	if err := s.DB.Create(&post).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create post"})
	}
	s.WSHub.broadcast <- WSMessage{Type: "post.created", Data: post}
	return c.JSON(http.StatusCreated, post)
}

// CreateComment (renders markdown & stores sanitized HTML)
func (s *Server) CreateComment(c echo.Context) error {
	type req struct {
		Content string `json:"content"`
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid"})
	}
	postIDStr := c.Param("id")
	pid, err := uuid.Parse(postIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid post id"})
	}
	user, err := userFromClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauth"})
	}
	html, _ := RenderMarkdown(r.Content)
	comment := Comment{
		PostID:      pid,
		UserID:      user.ID,
		Content:     r.Content,
		ContentHTML: html,
	}
	if err := s.DB.Create(&comment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create comment"})
	}
	s.WSHub.broadcast <- WSMessage{Type: "comment.created", Data: comment}
	return c.JSON(http.StatusCreated, comment)
}

func (s *Server) ListComments(c echo.Context) error {
	postIDStr := c.Param("id")
	pid, err := uuid.Parse(postIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var comments []Comment
	s.DB.Where("post_id = ?", pid).Order("created_at asc").Find(&comments)
	return c.JSON(http.StatusOK, comments)
}

// Generic helper to upsert votes for target_type
func (s *Server) upsertVote(userID uuid.UUID, targetID uuid.UUID, targetType string, val int8) (likes, dislikes int64, err error) {
	// upsert vote record
	var v Vote
	if err := s.DB.Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).First(&v).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			nv := Vote{
				UserID:     userID,
				TargetID:   targetID,
				TargetType: targetType,
				Value:      val,
			}
			if err := s.DB.Create(&nv).Error; err != nil {
				return 0, 0, err
			}
		} else {
			return 0, 0, err
		}
	} else {
		// toggle logic: if same value, remove vote (toggle off)
		if v.Value == val {
			if err := s.DB.Delete(&v).Error; err != nil {
				return 0, 0, err
			}
		} else {
			v.Value = val
			if err := s.DB.Save(&v).Error; err != nil {
				return 0, 0, err
			}
		}
	}

	// recalc tallies
	var likesCount int64
	var dislikesCount int64
	s.DB.Model(&Vote{}).Where("target_id = ? AND target_type = ? AND value = 1", targetID, targetType).Count(&likesCount)
	s.DB.Model(&Vote{}).Where("target_id = ? AND target_type = ? AND value = -1", targetID, targetType).Count(&dislikesCount)

	// apply tally to the actual row
	switch targetType {
	case "post":
		s.DB.Model(&Post{}).Where("id = ?", targetID).Updates(map[string]interface{}{"likes": likesCount, "dislikes": dislikesCount})
	case "comment":
		s.DB.Model(&Comment{}).Where("id = ?", targetID).Updates(map[string]interface{}{"likes": likesCount, "dislikes": dislikesCount})
	case "thread":
		// store likes/dislikes on thread? optional fields not present; if desired add likes/dislikes columns to threads
	default:
		// nothing
	}

	return likesCount, dislikesCount, nil
}

func (s *Server) VoteComment(c echo.Context) error {
	type req struct {
		Value int8 `json:"value"`
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid"})
	}
	if r.Value != 1 && r.Value != -1 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "value must be 1 or -1"})
	}
	commentIDStr := c.Param("id")
	cid, err := uuid.Parse(commentIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	user, err := userFromClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauth"})
	}
	likes, dislikes, err := s.upsertVote(user.ID, cid, "comment", r.Value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "vote failed"})
	}
	return c.JSON(http.StatusOK, echo.Map{"likes": likes, "dislikes": dislikes})
}

func (s *Server) VoteThread(c echo.Context) error {
	type req struct {
		Value int8 `json:"value"`
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid"})
	}
	if r.Value != 1 && r.Value != -1 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "value must be 1 or -1"})
	}
	threadIDStr := c.Param("id")
	tid, err := uuid.Parse(threadIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	user, err := userFromClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauth"})
	}
	_, _, err = s.upsertVote(user.ID, tid, "thread", r.Value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "vote failed"})
	}
	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
