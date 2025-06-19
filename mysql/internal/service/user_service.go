package service

import (
	"example.mysql/internal/db"
	"example.mysql/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(username, password string) (*models.User, error) {
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	result := db.DB.Create(user)
	return user, result.Error
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := db.DB.Where("username = ?", username).First(&user)
	return &user, result.Error
}

func VerifyPassword(user *models.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}
