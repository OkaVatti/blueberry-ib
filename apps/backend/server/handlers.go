// apps/backend/server/handlers.go
package server

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// Seed default boards and admin user
	seedDefaults(db)

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
			"https://localhost:9876",
			"https://127.0.0.1:9876",
		},
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
	}))

	// Attach server to context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("server", s)
			return next(c)
		}
	})

	// API routes
	api := e.Group("/api/v1")
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"status": "ok", "time": time.Now()})
	})

	// Auth routes
	api.POST("/auth/register", s.Register)
	api.POST("/auth/login", s.Login)
	api.POST("/auth/logout", s.AuthMiddleware(s.Logout))
	api.GET("/auth/me", s.AuthMiddleware(s.GetMe))

	// Board routes
	api.GET("/boards", s.ListBoards)
	api.GET("/boards/:slug", s.GetBoard)
	api.POST("/boards", s.AuthMiddleware(s.CreateBoard))
	api.PUT("/boards/:slug", s.RequireRoles("owner", "coowner", "admin")(s.UpdateBoard))
	api.DELETE("/boards/:slug", s.RequireRoles("owner", "coowner")(s.DeleteBoard))

	// Board membership
	api.POST("/boards/:slug/join", s.AuthMiddleware(s.JoinBoard))
	api.POST("/boards/:slug/leave", s.AuthMiddleware(s.LeaveBoard))
	api.GET("/boards/:slug/members", s.ListBoardMembers)

	// Posts
	api.GET("/boards/:slug/posts", s.ListPosts)
	api.POST("/boards/:slug/posts", s.AuthMiddleware(s.CreatePost))
	api.GET("/posts/:id", s.GetPost)
	api.PUT("/posts/:id", s.AuthMiddleware(s.UpdatePost))
	api.DELETE("/posts/:id", s.AuthMiddleware(s.DeletePost))

	// Threads
	api.GET("/boards/:slug/threads", s.ListThreadsForBoard)
	api.POST("/boards/:slug/threads", s.AuthMiddleware(s.CreateThread))
	api.GET("/threads/:id", s.GetThread)
	api.GET("/threads/:id/posts", s.ListPostsForThread)
	api.POST("/threads/:id/posts", s.AuthMiddleware(s.CreatePostInThread))
	api.POST("/threads/:id/lock", s.RequireRoles("moderator", "admin", "coowner", "owner")(s.LockThread))
	api.POST("/threads/:id/pin", s.RequireRoles("moderator", "admin", "coowner", "owner")(s.PinThread))

	// Comments
	api.GET("/posts/:id/comments", s.ListComments)
	api.POST("/posts/:id/comments", s.AuthMiddleware(s.CreateComment))
	api.GET("/comments/:id", s.GetComment)
	api.DELETE("/comments/:id", s.AuthMiddleware(s.DeleteComment))

	// Voting/Interactions
	api.POST("/posts/:id/vote", s.AuthMiddleware(s.VotePost))
	api.POST("/comments/:id/vote", s.AuthMiddleware(s.VoteComment))
	api.POST("/threads/:id/vote", s.AuthMiddleware(s.VoteThread))
	api.POST("/posts/:id/repost", s.AuthMiddleware(s.RepostPost))
	api.POST("/posts/:id/quote-repost", s.AuthMiddleware(s.QuoteRepostPost))

	// User routes
	api.GET("/users/:username", s.GetUserProfile)
	api.GET("/users/:username/posts", s.GetUserPosts)
	api.GET("/users/:username/threads", s.GetUserThreads)
	api.GET("/users/:username/comments", s.GetUserComments)
	api.PUT("/users/:username", s.AuthMiddleware(s.UpdateUserProfile))

	// Search
	api.GET("/search", s.Search)

	// Admin routes
	api.GET("/admin/users", s.RequireRoles("owner", "coowner", "admin")(s.AdminListUsers))
	api.POST("/admin/users/:id/ban", s.RequireRoles("owner", "coowner", "admin", "moderator")(s.AdminBanUser))
	api.DELETE("/admin/posts/:id", s.RequireRoles("owner", "coowner", "admin", "moderator")(s.AdminDeletePost))
	api.DELETE("/admin/comments/:id", s.RequireRoles("owner", "coowner", "admin", "moderator")(s.AdminDeleteComment))
	api.GET("/admin/logs", s.RequireRoles("owner", "coowner", "admin")(s.AdminListLogs))

	// Markdown rendering
	api.POST("/render", s.RenderMarkdownEndpoint)

	// WebSocket
	e.GET("/ws", func(c echo.Context) error {
		s.WSHub.ServeWS(c.Response(), c.Request())
		return nil
	})

	return e.Start(":" + cfg.Port)
}

// Seed defaults
func seedDefaults(db *gorm.DB) {
	// Create default boards
	defaultBoards := []Board{
		{Name: "Random", Slug: "b", Description: "Random posts and discussions", IsDefault: true},
		{Name: "Technology", Slug: "tech", Description: "Technology and programming", IsDefault: true},
		{Name: "Art", Slug: "art", Description: "Art, design, and creativity", IsDefault: true},
		{Name: "Gaming", Slug: "games", Description: "Video games and gaming culture", IsDefault: true},
		{Name: "Music", Slug: "music", Description: "Music discussion and sharing", IsDefault: true},
	}

	for _, b := range defaultBoards {
		var exists Board
		if err := db.Where("slug = ?", b.Slug).First(&exists).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&b)
			}
		}
	}

	// Create admin user if not exists
	var adminUser User
	if err := db.Where("username = ?", "admin").First(&adminUser).Error; err == gorm.ErrRecordNotFound {
		hash, _ := HashPassword("admin123")
		admin := User{
			Username:     "admin",
			DisplayName:  "Administrator",
			Email:        "admin@blueberry.local",
			PasswordHash: hash,
			Role:         "owner",
			Bio:          "Site administrator",
		}
		db.Create(&admin)
	}
}

// Helper functions
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
	// Check if user is banned
	if user.IsBanned {
		if user.BannedUntil != nil && time.Now().After(*user.BannedUntil) {
			// Ban expired, unban user
			user.IsBanned = false
			user.BannedUntil = nil
			s.DB.Save(&user)
		} else {
			return nil, echo.NewHTTPError(http.StatusForbidden, "user is banned")
		}
	}
	return &user, nil
}

// Auth endpoints
func (s *Server) Register(c echo.Context) error {
	type req struct {
		Username    string `json:"username" validate:"required,min=3,max=20"`
		Email       string `json:"email" validate:"required,email"`
		Password    string `json:"password" validate:"required,min=6"`
		DisplayName string `json:"displayName"`
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	// Check if user exists
	var existing User
	if err := s.DB.Where("username = ? OR email = ?", r.Username, r.Email).First(&existing).Error; err == nil {
		return c.JSON(http.StatusConflict, echo.Map{"error": "user already exists"})
	}

	hash, err := HashPassword(r.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "password hashing failed"})
	}

	user := User{
		Username:     strings.ToLower(strings.TrimSpace(r.Username)),
		DisplayName:  r.DisplayName,
		Email:        strings.ToLower(strings.TrimSpace(r.Email)),
		PasswordHash: hash,
		Role:         "user",
	}

	if user.DisplayName == "" {
		user.DisplayName = user.Username
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "registration failed"})
	}

	token, _ := GenerateJWT(s.Config.JWTSecret, &user)

	// Broadcast new user event
	s.WSHub.broadcast <- WSMessage{Type: "user.registered", Data: echo.Map{"username": user.Username}}

	return c.JSON(http.StatusCreated, echo.Map{
		"token": token,
		"user": echo.Map{
			"id":          user.ID,
			"username":    user.Username,
			"displayName": user.DisplayName,
			"role":        user.Role,
		},
	})
}

func (s *Server) Login(c echo.Context) error {
	type req struct {
		UsernameOrEmail string `json:"usernameOrEmail"`
		Password        string `json:"password"`
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	var user User
	if err := s.DB.Where("username = ? OR email = ?", r.UsernameOrEmail, r.UsernameOrEmail).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid credentials"})
	}

	ok, err := VerifyPassword(r.Password, user.PasswordHash)
	if err != nil || !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid credentials"})
	}

	// Check ban status
	if user.IsBanned {
		if user.BannedUntil != nil && time.Now().After(*user.BannedUntil) {
			// Ban expired
			user.IsBanned = false
			user.BannedUntil = nil
			s.DB.Save(&user)
		} else {
			return c.JSON(http.StatusForbidden, echo.Map{"error": "account is banned"})
		}
	}

	token, _ := GenerateJWT(s.Config.JWTSecret, &user)
	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
		"user": echo.Map{
			"id":          user.ID,
			"username":    user.Username,
			"displayName": user.DisplayName,
			"role":        user.Role,
			"ink":         user.Ink,
		},
	})
}

func (s *Server) Logout(c echo.Context) error {
	// In a real app, you might want to blacklist the token
	return c.JSON(http.StatusOK, echo.Map{"message": "logged out"})
}

func (s *Server) GetMe(c echo.Context) error {
	user, err := userFromClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id":           user.ID,
		"username":     user.Username,
		"displayName":  user.DisplayName,
		"email":        user.Email,
		"role":         user.Role,
		"ink":          user.Ink,
		"bio":          user.Bio,
		"profileImage": user.ProfileImage,
		"createdAt":    user.CreatedAt,
	})
}

// Board endpoints
func (s *Server) ListBoards(c echo.Context) error {
	var boards []Board
	query := s.DB.Order("is_default DESC, created_at DESC")

	// Optional filters
	if search := c.QueryParam("search"); search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Find(&boards)

	// Get member counts
	type BoardWithStats struct {
		Board
		MemberCount int64 `json:"memberCount"`
		PostCount   int64 `json:"postCount"`
	}

	results := make([]BoardWithStats, len(boards))
	for i, board := range boards {
		var memberCount, postCount int64
		s.DB.Model(&BoardMembership{}).Where("board_id = ?", board.ID).Count(&memberCount)
		s.DB.Model(&Post{}).Where("board_id = ?", board.ID).Count(&postCount)

		results[i] = BoardWithStats{
			Board:       board,
			MemberCount: memberCount,
			PostCount:   postCount,
		}
	}

	return c.JSON(http.StatusOK, results)
}

func (s *Server) GetBoard(c echo.Context) error {
	slug := c.Param("slug")
	var board Board
	if err := s.DB.Where("slug = ?", slug).First(&board).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "board not found"})
	}

	// Get stats
	var memberCount, postCount, threadCount int64
	s.DB.Model(&BoardMembership{}).Where("board_id = ?", board.ID).Count(&memberCount)
	s.DB.Model(&Post{}).Where("board_id = ?", board.ID).Count(&postCount)
	s.DB.Model(&Thread{}).Where("board_id = ?", board.ID).Count(&threadCount)

	return c.JSON(http.StatusOK, echo.Map{
		"board": board,
		"stats": echo.Map{
			"members": memberCount,
			"posts":   postCount,
			"threads": threadCount,
		},
	})
}

func (s *Server) CreateBoard(c echo.Context) error {
	type req struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	user, err := userFromClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	board := Board{
		Name:        r.Name,
		Slug:        strings.ToLower(strings.TrimSpace(r.Slug)),
		Description: r.Description,
		CreatedByID: &user.ID,
	}

	if err := s.DB.Create(&board).Error; err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "failed to create board"})
	}

	// Auto-join creator as admin
	membership := BoardMembership{
		BoardID: board.ID,
		UserID:  user.ID,
		Role:    "admin",
	}
	s.DB.Create(&membership)

	s.WSHub.broadcast <- WSMessage{Type: "board.created", Data: board}
	return c.JSON(http.StatusCreated, board)
}

func (s *Server) UpdateBoard(c echo.Context) error {
	slug := c.Param("slug")
	var board Board
	if err := s.DB.Where("slug = ?", slug).First(&board).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "board not found"})
	}

	type req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Settings    string `json:"settings"`
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	if r.Name != "" {
		board.Name = r.Name
	}
	if r.Description != "" {
		board.Description = r.Description
	}
	if r.Settings != "" {
		board.Settings = r.Settings
	}

	if err := s.DB.Save(&board).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "update failed"})
	}

	return c.JSON(http.StatusOK, board)
}

func (s *Server) DeleteBoard(c echo.Context) error {
	slug := c.Param("slug")
	var board Board
	if err := s.DB.Where("slug = ?", slug).First(&board).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "board not found"})
	}

	if board.IsDefault {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "cannot delete default board"})
	}

	if err := s.DB.Delete(&board).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "delete failed"})
	}

	s.WSHub.broadcast <- WSMessage{Type: "board.deleted", Data: echo.Map{"slug": slug}}
	return c.JSON(http.StatusOK, echo.Map{"message": "board deleted"})
}

// Board membership
func (s *Server) JoinBoard(c echo.Context) error {
	slug := c.Param("slug")
	user, err := userFromClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	var board Board
	if err := s.DB.Where("slug = ?", slug).First(&board).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "board not found"})
	}

	// Check if already member
	var existing BoardMembership
	if err := s.DB.Where("board_id = ? AND user_id = ?", board.ID, user.ID).First(&existing).Error; err == nil {
		return c.JSON(http.StatusConflict, echo.Map{"error": "already a member"})
	}

	membership := BoardMembership{
		BoardID: board.ID,
		UserID:  user.ID,
		Role:    "member",
	}

	if err := s.DB.Create(&membership).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "join failed"})
	}

	s.WSHub.broadcast <- WSMessage{
		Type: "board.joined",
		Data: echo.Map{"board": slug, "user": user.Username},
	}

	return c.JSON(http.StatusOK, membership)
}

func (s *Server) LeaveBoard(c echo.Context) error {
	slug := c.Param("slug")
	user, err := userFromClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	var board Board
	if err := s.DB.Where("slug = ?", slug).First(&board).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "board not found"})
	}

	if err := s.DB.Where("board_id = ? AND user_id = ?", board.ID, user.ID).Delete(&BoardMembership{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "leave failed"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "left board"})
}

func (s *Server) ListBoardMembers(c echo.Context) error {
	slug := c.Param("slug")
	var board Board
	if err := s.DB.Where("slug = ?", slug).First(&board).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "board not found"})
	}

	var members []struct {
		BoardMembership
		User User
	}

	s.DB.Table("board_memberships").
		Select("board_memberships.*, users.*").
		Joins("JOIN users ON users.id = board_memberships.user_id").
		Where("board_memberships.board_id = ?", board.ID).
		Scan(&members)

	return c.JSON(http.StatusOK, members)
}

// Post endpoints
func (s *Server) ListPosts(c echo.Context) error {
	slug := c.Param("slug")
	var board Board
	if err := s.DB.Where("slug = ?", slug).First(&board).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "board not found"})
	}

	limit := 20
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := 0
	if o := c.QueryParam("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	var posts []Post
	s.DB.Where("board_id = ? AND thread_id IS NULL", board.ID).
		Preload("User").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&posts)

	return c.JSON(http.StatusOK, posts)
}

func (s *Server) GetPost(c echo.Context) error {
	id := c.Param("id")
	pid, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	var post Post
	if err := s.DB.Preload("User").Preload("Board").First(&post, "id = ?", pid).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "post not found"})
	}

	return c.JSON(http.StatusOK, post)
}

func (s *Server) UpdatePost(c echo.Context) error {
	id := c.Param("id")
	pid, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	user, err := userFromClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	var post Post
	if err := s.DB.First(&post, "id = ?", pid).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "post not found"})
	}

	// Check ownership
	if post.UserID != user.ID && user.Role != "admin" && user.Role != "owner" && user.Role != "coowner" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "not authorized"})
	}

	type req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	if r.Title != "" {
		post.Title = r.Title
	}
	if r.Content != "" {
		post.Content = r.Content
		html, _ := RenderMarkdown(r.Content)
		post.ContentHTML = html
	}

	if err := s.DB.Save(&post).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "update failed"})
	}

	s.WSHub.broadcast <- WSMessage{Type: "post.updated", Data: post}
	return c.JSON(http.StatusOK, post)
}

// Repost functionality
func (s *Server) RepostPost(c echo.Context) error {
	id := c.Param("id")
	pid, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	user, err := userFromClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	type req struct {
		BoardID string `json:"boardId"`
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	boardID, err := uuid.Parse(r.BoardID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid board id"})
	}

	// Check if post exists
	var post Post
	if err := s.DB.First(&post, "id = ?", pid).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "post not found"})
	}

	// Check if already reposted by this user to this board
	var existing Repost
	if err := s.DB.Where("user_id = ? AND original_id = ? AND board_id = ?", user.ID, pid, boardID).First(&existing).Error; err == nil {
		return c.JSON(http.StatusConflict, echo.Map{"error": "already reposted"})
	}

	repost := Repost{
		UserID:       user.ID,
		OriginalType: "post",
		OriginalID:   pid,
		BoardID:      boardID,
	}

	if err := s.DB.Create(&repost).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "repost failed"})
	}

	s.WSHub.broadcast <- WSMessage{Type: "post.reposted", Data: repost}
	return c.JSON(http.StatusCreated, repost)
}

func (s *Server) QuoteRepostPost(c echo.Context) error {
	id := c.Param("id")
	pid, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	user, err := userFromClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	type req struct {
		BoardID string `json:"boardId"`
		Content string `json:"content"`
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	boardID, err := uuid.Parse(r.BoardID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid board id"})
	}

	// Check if post exists
	var post Post
	if err := s.DB.First(&post, "id = ?", pid).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "post not found"})
	}

	html, _ := RenderMarkdown(r.Content)

	quoteRepost := QuoteRepost{
		UserID:       user.ID,
		OriginalType: "post",
		OriginalID:   pid,
		BoardID:      boardID,
		Content:      r.Content,
		ContentHTML:  html,
	}

	if err := s.DB.Create(&quoteRepost).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "quote repost failed"})
	}

	s.WSHub.broadcast <- WSMessage{Type: "post.quote_reposted", Data: quoteRepost}
	return c.JSON(http.StatusCreated, quoteRepost)
}

// Thread endpoints
func (s *Server) LockThread(c echo.Context) error {
	id := c.Param("id")
	tid, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	var thread Thread
	if err := s.DB.First(&thread, "id = ?", tid).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "thread not found"})
	}

	thread.IsLocked = !thread.IsLocked
	if err := s.DB.Save(&thread).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "lock failed"})
	}

	action := "locked"
	if !thread.IsLocked {
		action = "unlocked"
	}

	s.WSHub.broadcast <- WSMessage{Type: "thread." + action, Data: echo.Map{"id": tid}}
	return c.JSON(http.StatusOK, echo.Map{"locked": thread.IsLocked})
}

func (s *Server) PinThread(c echo.Context) error {
	id := c.Param("id")
	tid, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	var thread Thread
	if err := s.DB.First(&thread, "id = ?", tid).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "thread not found"})
	}

	thread.IsPinned = !thread.IsPinned
	if err := s.DB.Save(&thread).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "pin failed"})
	}

	action := "pinned"
	if !thread.IsPinned {
		action = "unpinned"
	}

	s.WSHub.broadcast <- WSMessage{Type: "thread." + action, Data: echo.Map{"id": tid}}
	return c.JSON(http.StatusOK, echo.Map{"pinned": thread.IsPinned})
}

// User profile endpoints
func (s *Server) GetUserProfile(c echo.Context) error {
	username := c.Param("username")
	var user User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
	}

	// Don't expose sensitive data
	return c.JSON(http.StatusOK, echo.Map{
		"id":           user.ID,
		"username":     user.Username,
		"displayName":  user.DisplayName,
		"bio":          user.Bio,
		"profileImage": user.ProfileImage,
		"ink":          user.Ink,
		"createdAt":    user.CreatedAt,
	})
}

func (s *Server) GetUserPosts(c echo.Context) error {
	username := c.Param("username")
	var user User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
	}

	var posts []Post
	s.DB.Where("user_id = ?", user.ID).
		Preload("Board").
		Order("created_at desc").
		Limit(50).
		Find(&posts)

	return c.JSON(http.StatusOK, posts)
}

func (s *Server) GetUserThreads(c echo.Context) error {
	username := c.Param("username")
	var user User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
	}

	var threads []Thread
	s.DB.Where("user_id = ?", user.ID).
		Preload("Board").
		Order("created_at desc").
		Limit(50).
		Find(&threads)

	return c.JSON(http.StatusOK, threads)
}

func (s *Server) GetUserComments(c echo.Context) error {
	username := c.Param("username")
	var user User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
	}

	var comments []Comment
	s.DB.Where("user_id = ?", user.ID).
		Order("created_at desc").
		Limit(50).
		Find(&comments)

	return c.JSON(http.StatusOK, comments)
}

func (s *Server) UpdateUserProfile(c echo.Context) error {
	username := c.Param("username")
	user, err := userFromClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	// Check if updating own profile
	if user.Username != username && user.Role != "admin" && user.Role != "owner" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "not authorized"})
	}

	type req struct {
		DisplayName   string `json:"displayName"`
		Bio           string `json:"bio"`
		ProfileImage  string `json:"profileImage"`
		ProfileBanner string `json:"profileBanner"`
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	if r.DisplayName != "" {
		user.DisplayName = r.DisplayName
	}
	if r.Bio != "" {
		user.Bio = r.Bio
	}
	if r.ProfileImage != "" {
		user.ProfileImage = r.ProfileImage
	}
	if r.ProfileBanner != "" {
		user.ProfileBanner = r.ProfileBanner
	}

	if err := s.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "update failed"})
	}

	return c.JSON(http.StatusOK, user)
}

// Search endpoint
func (s *Server) Search(c echo.Context) error {
	q := c.QueryParam("q")
	if q == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "query required"})
	}

	searchType := c.QueryParam("type") // users, posts, threads, boards
	limit := 20

	results := echo.Map{}

	if searchType == "" || searchType == "boards" {
		var boards []Board
		s.DB.Where("name ILIKE ? OR description ILIKE ?", "%"+q+"%", "%"+q+"%").
			Limit(limit).
			Find(&boards)
		results["boards"] = boards
	}

	if searchType == "" || searchType == "posts" {
		var posts []Post
		s.DB.Where("title ILIKE ? OR content ILIKE ?", "%"+q+"%", "%"+q+"%").
			Preload("User").
			Preload("Board").
			Limit(limit).
			Find(&posts)
		results["posts"] = posts
	}

	if searchType == "" || searchType == "users" {
		var users []User
		s.DB.Where("username ILIKE ? OR display_name ILIKE ?", "%"+q+"%", "%"+q+"%").
			Limit(limit).
			Find(&users)

		// Filter sensitive data
		safeUsers := make([]echo.Map, len(users))
		for i, u := range users {
			safeUsers[i] = echo.Map{
				"id":          u.ID,
				"username":    u.Username,
				"displayName": u.DisplayName,
				"ink":         u.Ink,
			}
		}
		results["users"] = safeUsers
	}

	if searchType == "" || searchType == "threads" {
		var threads []Thread
		s.DB.Where("title ILIKE ?", "%"+q+"%").
			Preload("User").
			Preload("Board").
			Limit(limit).
			Find(&threads)
		results["threads"] = threads
	}

	return c.JSON(http.StatusOK, results)
}

// Comment endpoints
func (s *Server) GetComment(c echo.Context) error {
	id := c.Param("id")
	cid, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	var comment Comment
	if err := s.DB.Preload("User").First(&comment, "id = ?", cid).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "comment not found"})
	}

	return c.JSON(http.StatusOK, comment)
}

func (s *Server) DeleteComment(c echo.Context) error {
	id := c.Param("id")
	cid, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	user, err := userFromClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	var comment Comment
	if err := s.DB.First(&comment, "id = ?", cid).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "comment not found"})
	}

	// Check ownership or admin
	if comment.UserID != user.ID && user.Role != "admin" && user.Role != "moderator" && user.Role != "owner" && user.Role != "coowner" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "not authorized"})
	}

	if err := s.DB.Delete(&comment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "delete failed"})
	}

	s.WSHub.broadcast <- WSMessage{Type: "comment.deleted", Data: echo.Map{"id": cid}}
	return c.JSON(http.StatusOK, echo.Map{"message": "comment deleted"})
}
