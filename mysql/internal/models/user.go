package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint64  `gorm:"primaryKey;autoIncrement"`
	Username     string  `gorm:"size:64;not null;uniqueIndex"`
	Email        *string `gorm:"size:128;uniqueIndex"`
	Phone        *string `gorm:"size:20;uniqueIndex"`
	PasswordHash string  `gorm:"size:256;not null"`
	Status       int8    `gorm:"default:1;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
