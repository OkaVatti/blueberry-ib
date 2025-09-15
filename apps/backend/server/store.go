package server

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGorm(dbURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// automigrate
	if err := db.AutoMigrate(&User{}, &Board{}, &Post{}, &Comment{}, &Vote{}); err != nil {
		return nil, err
	}
	return db, nil
}
