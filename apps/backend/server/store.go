// apps/backend/server/store.go
package server

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGorm(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// AutoMigrate all models
	if err := db.AutoMigrate(
		&User{},
		&Board{},
		&Thread{},
		&Post{},
		&Comment{},
		&Vote{},
		&AuditLog{},
		&Repost{},
		&QuoteRepost{},
		&BoardMembership{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
