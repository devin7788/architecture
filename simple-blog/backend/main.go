package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	. "simple-blog/api"
)

func main() {
	db, _ := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	Migrate(db)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/login", LoginHandler)
	auth := r.Group("/", AuthMiddleware())
	{
		auth.POST("/articles", CreateArticleHandler(db))
		auth.DELETE("/articles/:id", DeleteArticleHandler(db))
	}
	r.GET("/articles", ListArticlesHandler(db))
	r.GET("/articles/:id", GetArticleHandler(db))

	r.Run(":8080")
}
