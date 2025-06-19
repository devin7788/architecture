package db

import (
	"fmt"
	"log"

	"example.mysql/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMySQL() {
	dsn := "root:admin123@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	fmt.Println("Connected to MySQL and migrated schema.")
}
