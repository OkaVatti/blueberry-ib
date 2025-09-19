// apps/backend/server/models.go
package server

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username     string    `gorm:"uniqueIndex;not null"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	Role         string    `gorm:"default:user"` // user, mod, admin
	Ink          int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IsBanned     bool      `gorm:"default:false"`
	BannedUntil  time.Time `gorm:""`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

type Board struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"uniqueIndex;not null"`
	Slug        string    `gorm:"uniqueIndex;not null"`
	Description string
	IsDefault   bool
	CreatedByID *uuid.UUID
	CreatedBy   *User `gorm:"foreignKey:CreatedByID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (b *Board) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}

type Thread struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	BoardID   uuid.UUID `gorm:"type:uuid;index;not null"`
	Board     Board
	UserID    uuid.UUID `gorm:"type:uuid;index"`
	User      User
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Thread) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return
}

type Post struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	BoardID     uuid.UUID `gorm:"type:uuid;index;not null"`
	Board       Board
	ThreadID    *uuid.UUID `gorm:"type:uuid;index"` // nullable; set for threads
	Thread      *Thread
	UserID      uuid.UUID `gorm:"type:uuid;index"`
	User        User
	Title       string
	Content     string `gorm:"type:text"`
	ContentHTML string `gorm:"type:text"`
	Likes       int64
	Dislikes    int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}

type Comment struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	PostID      uuid.UUID `gorm:"type:uuid;index;not null"`
	Post        Post
	UserID      uuid.UUID `gorm:"type:uuid;index"`
	User        User
	Content     string `gorm:"type:text"`
	ContentHTML string `gorm:"type:text"`
	Likes       int64
	Dislikes    int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return
}

// Vote represents a like/dislike for a target (post, comment, thread)
type Vote struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;index;not null"`
	TargetID   uuid.UUID `gorm:"type:uuid;index;not null"`  // post, comment, or thread id
	TargetType string    `gorm:"type:varchar(10);not null"` // "post" | "comment" | "thread"
	Value      int8      // 1 for like, -1 for dislike
	CreatedAt  time.Time
}

func (v *Vote) BeforeCreate(tx *gorm.DB) (err error) {
	v.ID = uuid.New()
	return
}
