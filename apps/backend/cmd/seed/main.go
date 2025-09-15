package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/okavatti/blueberry/backend/server"
	"gorm.io/gorm"
)

func main() {
	gofakeit.Seed(time.Now().UnixNano())
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable required")
	}
	db, err := server.NewGorm(dsn)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// make sure migrations ran; auto-migrate AuditLog (so seed creates entries)
	if err := db.AutoMigrate(&server.AuditLog{}); err != nil {
		log.Fatalf("auto migrate audit log: %v", err)
	}

	// create owner
	ensureOwner(db)

	// create boards
	ensureBoards(db)

	// how many users/posts/comments
	numUsers := 100
	numPosts := 500
	numComments := 2000

	if v := os.Getenv("SEED_USERS"); v != "" {
		if val, e := strconv.Atoi(v); e == nil {
			numUsers = val
		}
	}
	if v := os.Getenv("SEED_POSTS"); v != "" {
		if val, e := strconv.Atoi(v); e == nil {
			numPosts = val
		}
	}
	if v := os.Getenv("SEED_COMMENTS"); v != "" {
		if val, e := strconv.Atoi(v); e == nil {
			numComments = val
		}
	}

	users := make([]server.User, 0, numUsers)
	// create several admins/mods/coowner
	admins := []struct{ username, role string }{
		{"coowner", "coowner"},
		{"admin1", "admin"},
		{"mod1", "moderator"},
	}

	for _, a := range admins {
		u := createUser(db, a.username, a.username+"@example.test", "password123", a.role)
		users = append(users, *u)
	}

	// create random users
	for i := 0; i < numUsers; i++ {
		username := gofakeit.Username()
		email := gofakeit.Email()
		u := createUser(db, username, email, "password123", "user")
		users = append(users, *u)
	}

	// create posts
	posts := make([]server.Post, 0, numPosts)
	boardIDs := getBoardIDs(db)
	for i := 0; i < numPosts; i++ {
		author := users[rand.Intn(len(users))]
		boardID := boardIDs[rand.Intn(len(boardIDs))]
		title := gofakeit.Sentence(3)
		content := gofakeit.Paragraph(1, 3, 10, " ")
		html, _ := server.RenderMarkdown(content)
		p := server.Post{
			BoardID:     boardID,
			UserID:      author.ID,
			Title:       title,
			Content:     content,
			ContentHTML: html,
		}
		if err := db.Create(&p).Error; err == nil {
			posts = append(posts, p)
		}
	}

	// create comments on random posts by random users
	for i := 0; i < numComments; i++ {
		author := users[rand.Intn(len(users))]
		if len(posts) == 0 {
			break
		}
		post := posts[rand.Intn(len(posts))]
		content := gofakeit.Sentence(8)
		html, _ := server.RenderMarkdown(content)
		c := server.Comment{
			PostID:      post.ID,
			UserID:      author.ID,
			Content:     content,
			ContentHTML: html,
		}
		_ = db.Create(&c)
	}

	fmt.Println("seeding complete")
}

func ensureOwner(db *gorm.DB) {
	// owner username: lilithinaparka
	var u server.User
	err := db.Where("username = ?", "lilithinaparka").First(&u).Error
	if err == nil {
		fmt.Println("owner already exists")
		return
	}
	// create owner
	createUser(db, "lilithinaparka", "lilith@example.test", "ownerpassword", "owner")
	fmt.Println("created owner user: lilithinaparka")
}

func ensureBoards(db *gorm.DB) {
	defaults := []server.Board{
		{Name: "Random", Slug: "b", Description: "Random board", IsDefault: true},
		{Name: "Technology", Slug: "tech", Description: "Technology", IsDefault: true},
		{Name: "Art", Slug: "art", Description: "Art", IsDefault: true},
		{Name: "Show", Slug: "show", Description: "Showcase", IsDefault: true},
	}
	for _, b := range defaults {
		var exists server.Board
		if err := db.Where("slug = ?", b.Slug).First(&exists).Error; err != nil {
			db.Create(&b)
		}
	}
}

func createUser(db *gorm.DB, username, email, rawPassword, role string) *server.User {
	pwHash, err := server.HashPassword(rawPassword)
	if err != nil {
		log.Fatalf("hash error: %v", err)
	}
	u := &server.User{
		Username:     username,
		Email:        email,
		PasswordHash: pwHash,
		Role:         role,
	}
	if err := db.Create(u).Error; err != nil {
		// if exists, return existing
		var existing server.User
		if db.Where("username = ?", username).First(&existing).Error == nil {
			return &existing
		}
		log.Fatalf("create user failed: %v", err)
	}
	return u
}

func getBoardIDs(db *gorm.DB) []uuid.UUID {
	var boards []server.Board
	db.Find(&boards)
	out := make([]uuid.UUID, 0, len(boards))
	for _, b := range boards {
		out = append(out, b.ID)
	}
	return out
}
