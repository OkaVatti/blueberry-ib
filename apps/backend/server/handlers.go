package server

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

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

	// seed default boards if not exists
	seedDefaultBoards(db)

	go s.WSHub.Run()

	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("server", s)
			return next(c)
		}
	})

	// routes
	api := e.Group("/api/v1")
	api.POST("/auth/register", s.Register)
	api.POST("/auth/login", s.Login)

	api.GET("/boards", s.ListBoards)
	api.POST("/boards", s.AuthMiddleware(s.CreateBoard)) // protected

	api.GET("/boards/:slug/posts", s.ListPosts)
	api.POST("/boards/:slug/posts", s.AuthMiddleware(s.CreatePost))

	api.POST("/posts/:id/comments", s.AuthMiddleware(s.CreateComment))
	api.GET("/posts/:id/comments", s.ListComments)

	api.POST("/posts/:id/vote", s.AuthMiddleware(s.VotePost))

	e.GET("/ws", func(c echo.Context) error {
		s.WSHub.ServeWS(c.Response(), c.Request())
		return nil
	})

	return e.Start(":" + cfg.Port)
}

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

func (s *Server) CreatePost(c echo.Context) error {
	type req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid"})
	}
	slug := c.Param("slug")
	var board Board
	if err := s.DB.Where("slug = ?", slug).First(&board).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "board not found"})
	}
	user, err := userFromClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauth"})
	}
	post := Post{
		BoardID: board.ID,
		UserID:  user.ID,
		Title:   r.Title,
		Content: r.Content,
	}
	if err := s.DB.Create(&post).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create post"})
	}
	// broadcast via websocket
	s.WSHub.broadcast <- WSMessage{Type: "post.created", Data: post}
	return c.JSON(http.StatusCreated, post)
}

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
	comment := Comment{
		PostID:  pid,
		UserID:  user.ID,
		Content: r.Content,
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
