// apps/backend/server/models.go
package server

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username      string    `gorm:"uniqueIndex;not null"`
	DisplayName   string
	Email         string `gorm:"uniqueIndex;not null"`
	PasswordHash  string `gorm:"not null"`
	Role          string `gorm:"default:user"` // user, moderator, admin, coowner, owner
	Ink           int64  `gorm:"default:0"`
	Bio           string `gorm:"type:text"`
	ProfileImage  string
	ProfileBanner string
	Birthday      *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	IsBanned      bool           `gorm:"default:false"`
	BannedUntil   *time.Time

	// Relations
	Posts    []Post    `gorm:"foreignKey:UserID"`
	Threads  []Thread  `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:UserID"`
	Votes    []Vote    `gorm:"foreignKey:UserID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	if u.DisplayName == "" {
		u.DisplayName = u.Username
	}
	return
}

type Board struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"uniqueIndex;not null"`
	Slug        string    `gorm:"uniqueIndex;not null"`
	Description string    `gorm:"type:text"`
	IsDefault   bool      `gorm:"default:false"`
	Settings    string    `gorm:"type:jsonb"` // JSON settings
	CreatedByID *uuid.UUID
	CreatedBy   *User `gorm:"foreignKey:CreatedByID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Relations
	Posts   []Post            `gorm:"foreignKey:BoardID"`
	Threads []Thread          `gorm:"foreignKey:BoardID"`
	Members []BoardMembership `gorm:"foreignKey:BoardID"`
}

func (b *Board) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}

type BoardMembership struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	BoardID  uuid.UUID `gorm:"type:uuid;index;not null"`
	Board    Board
	UserID   uuid.UUID `gorm:"type:uuid;index;not null"`
	User     User
	Role     string `gorm:"default:member"` // member, moderator, admin
	JoinedAt time.Time
}

func (bm *BoardMembership) BeforeCreate(tx *gorm.DB) (err error) {
	bm.ID = uuid.New()
	bm.JoinedAt = time.Now()
	return
}

type Thread struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	BoardID   uuid.UUID `gorm:"type:uuid;index;not null"`
	Board     Board
	UserID    uuid.UUID `gorm:"type:uuid;index"`
	User      User
	Title     string `gorm:"type:text"`
	IsLocked  bool   `gorm:"default:false"`
	IsPinned  bool   `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	Posts []Post `gorm:"foreignKey:ThreadID"`
}

func (t *Thread) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return
}

type Post struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	BoardID          uuid.UUID `gorm:"type:uuid;index;not null"`
	Board            Board
	ThreadID         *uuid.UUID `gorm:"type:uuid;index"` // nullable
	Thread           *Thread
	UserID           uuid.UUID `gorm:"type:uuid;index"`
	User             User
	Title            string `gorm:"type:text"`
	Content          string `gorm:"type:text"`
	ContentHTML      string `gorm:"type:text"`
	MediaAttachments string `gorm:"type:jsonb"` // JSON array of media
	Likes            int64  `gorm:"default:0"`
	Dislikes         int64  `gorm:"default:0"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`

	// Relations
	Comments []Comment `gorm:"foreignKey:PostID"`
	Votes    []Vote    `gorm:"foreignKey:TargetID"`
	Reposts  []Repost  `gorm:"foreignKey:OriginalID"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}

type Comment struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey"`
	PostID      *uuid.UUID `gorm:"type:uuid;index"` // nullable for thread comments
	Post        *Post
	ThreadID    *uuid.UUID `gorm:"type:uuid;index"` // nullable for post comments
	Thread      *Thread
	UserID      uuid.UUID `gorm:"type:uuid;index"`
	User        User
	ParentID    *uuid.UUID `gorm:"type:uuid;index"` // for nested comments
	Content     string     `gorm:"type:text"`
	ContentHTML string     `gorm:"type:text"`
	ThreadLevel int        `gorm:"default:0"` // nesting depth
	Likes       int64      `gorm:"default:0"`
	Dislikes    int64      `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Relations
	Replies []Comment `gorm:"foreignKey:ParentID"`
	Votes   []Vote    `gorm:"foreignKey:TargetID"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return
}

type Vote struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;index;not null"`
	User       User
	TargetID   uuid.UUID `gorm:"type:uuid;index;not null"`
	TargetType string    `gorm:"type:varchar(10);not null"` // "post", "comment", "thread"
	Value      int8      `gorm:"not null"`                  // 1 for like, -1 for dislike
	IsActive   bool      `gorm:"default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (v *Vote) BeforeCreate(tx *gorm.DB) (err error) {
	v.ID = uuid.New()
	return
}

type Repost struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID `gorm:"type:uuid;index;not null"`
	User         User
	OriginalType string    `gorm:"type:varchar(10);not null"` // "post", "thread", "comment"
	OriginalID   uuid.UUID `gorm:"type:uuid;index;not null"`
	BoardID      uuid.UUID `gorm:"type:uuid;index;not null"`
	Board        Board
	CreatedAt    time.Time
}

func (r *Repost) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return
}

type QuoteRepost struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID           uuid.UUID `gorm:"type:uuid;index;not null"`
	User             User
	OriginalType     string    `gorm:"type:varchar(10);not null"`
	OriginalID       uuid.UUID `gorm:"type:uuid;index;not null"`
	BoardID          uuid.UUID `gorm:"type:uuid;index;not null"`
	Board            Board
	Content          string `gorm:"type:text"`
	ContentHTML      string `gorm:"type:text"`
	MediaAttachments string `gorm:"type:jsonb"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (q *QuoteRepost) BeforeCreate(tx *gorm.DB) (err error) {
	q.ID = uuid.New()
	return
}
