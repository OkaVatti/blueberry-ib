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

type Post struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	BoardID   uuid.UUID `gorm:"type:uuid;index;not null"`
	Board     Board
	UserID    uuid.UUID `gorm:"type:uuid;index"`
	User      User
	Title     string
	Content   string `gorm:"type:text"`
	Likes     int64
	Dislikes  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}

type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	PostID    uuid.UUID `gorm:"type:uuid;index;not null"`
	Post      Post
	UserID    uuid.UUID `gorm:"type:uuid;index"`
	User      User
	Content   string `gorm:"type:text"`
	Likes     int64
	Dislikes  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return
}

// Simple Like table so user can toggle like/dislike
type Vote struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;index;not null"`
	TargetID   uuid.UUID `gorm:"type:uuid;index;not null"`  // post or comment id
	TargetType string    `gorm:"type:varchar(10);not null"` // "post" or "comment"
	Value      int8      // 1 for like, -1 for dislike
	CreatedAt  time.Time
}

func (v *Vote) BeforeCreate(tx *gorm.DB) (err error) {
	v.ID = uuid.New()
	return
}
