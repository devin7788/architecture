package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:64"`
	Password string `gorm:"size:128"` // 建议存加盐后的hash密码
}
