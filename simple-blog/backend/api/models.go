package api

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       uint `gorm:"primaryKey"`
	Username string
	Password string
}

type Article struct {
	ID          uint `gorm:"primaryKey"`
	Title       string
	Content     string // Markdown 原文
	HTMLContent string // HTML渲染内容
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AuthorID    uint
	Author      User
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Article{})
}
