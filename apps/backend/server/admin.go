package server

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// AuditLog model (mirror of DB table)
type AuditLog struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	ActorID    *uuid.UUID `gorm:"type:uuid;index"`
	Action     string     `gorm:"type:text;not null"`
	TargetType string     `gorm:"type:text"`
	TargetID   *uuid.UUID `gorm:"type:uuid"`
	Details    string     `gorm:"type:jsonb"`
	CreatedAt  time.Time
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New()
	return
}

// Helper to create audit log
func (s *Server) createAudit(actor *User, action string, targetType string, targetID *uuid.UUID, details string) {
	var actorID *uuid.UUID
	if actor != nil {
		actorID = &actor.ID
	}
	al := AuditLog{
		ActorID:    actorID,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Details:    details,
	}
	_ = s.DB.Create(&al).Error
}

// RequireRoles middleware factory
func (s *Server) RequireRoles(roles ...string) echo.MiddlewareFunc {
	allowed := map[string]bool{}
	for _, r := range roles {
		allowed[r] = true
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, err := userFromClaims(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauth"})
			}
			if allowed[user.Role] {
				return next(c)
			}
			return c.JSON(http.StatusForbidden, echo.Map{"error": "insufficient role"})
		}
	}
}

// Admin: list users (basic)
func (s *Server) AdminListUsers(c echo.Context) error {
	var users []User
	s.DB.Order("created_at desc").Find(&users)
	// mask password hash
	out := make([]map[string]interface{}, 0, len(users))
	for _, u := range users {
		out = append(out, map[string]interface{}{
			"id":           u.ID,
			"username":     u.Username,
			"email":        u.Email,
			"role":         u.Role,
			"is_banned":    u.IsBanned,
			"banned_until": u.BannedUntil,
			"ink":          u.Ink,
			"created_at":   u.CreatedAt,
		})
	}
	return c.JSON(http.StatusOK, out)
}

// Admin: ban user
func (s *Server) AdminBanUser(c echo.Context) error {
	type req struct {
		DurationMinutes int    `json:"duration_minutes"` // 0 for permanent
		Reason          string `json:"reason"`
		Unban           bool   `json:"unban"` // if true, remove ban
	}
	var r req
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid"})
	}
	userIDstr := c.Param("id")
	uid, err := uuid.Parse(userIDstr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	actor, _ := userFromClaims(c)

	var target User
	if err := s.DB.First(&target, "id = ?", uid).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
	}

	if r.Unban {
		target.IsBanned = false
		target.BannedUntil = time.Time{}
		if err := s.DB.Save(&target).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to unban"})
		}
		s.createAudit(actor, "unban_user", "user", &target.ID, `{"reason":""}`)
		return c.JSON(http.StatusOK, echo.Map{"status": "unbanned"})
	}

	target.IsBanned = true
	if r.DurationMinutes <= 0 {
		// permanent ban
		target.BannedUntil = time.Time{} // zero means permanent when is_banned true
	} else {
		target.BannedUntil = time.Now().Add(time.Duration(r.DurationMinutes) * time.Minute)
	}
	if err := s.DB.Save(&target).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to ban"})
	}
	s.createAudit(actor, "ban_user", "user", &target.ID, `{"duration_minutes":`+string(rune(r.DurationMinutes))+`,"reason":"`+r.Reason+`"}`)
	return c.JSON(http.StatusOK, echo.Map{"status": "banned"})
}

// Admin: delete post (soft delete by removing content & marking)
func (s *Server) AdminDeletePost(c echo.Context) error {
	postID := c.Param("id")
	pid, err := uuid.Parse(postID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var post Post
	if err := s.DB.First(&post, "id = ?", pid).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "post not found"})
	}
	actor, _ := userFromClaims(c)
	// For simplicity: delete row
	if err := s.DB.Delete(&post).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "delete failed"})
	}
	s.createAudit(actor, "delete_post", "post", &pid, `{"reason":"admin_delete"}`)
	return c.JSON(http.StatusOK, echo.Map{"status": "deleted"})
}

// Admin: delete comment
func (s *Server) AdminDeleteComment(c echo.Context) error {
	cid := c.Param("id")
	id, err := uuid.Parse(cid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var com Comment
	if err := s.DB.First(&com, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "comment not found"})
	}
	actor, _ := userFromClaims(c)
	if err := s.DB.Delete(&com).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "delete failed"})
	}
	s.createAudit(actor, "delete_comment", "comment", &id, `{"reason":"admin_delete"}`)
	return c.JSON(http.StatusOK, echo.Map{"status": "deleted"})
}

// Admin: list audit logs (paginated)
func (s *Server) AdminListLogs(c echo.Context) error {
	var logs []AuditLog
	// simple: latest 500
	s.DB.Order("created_at desc").Limit(500).Find(&logs)
	return c.JSON(http.StatusOK, logs)
}
