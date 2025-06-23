package main

import (
	"github.com/gin-gonic/gin"
	"weblogin/config"
	"weblogin/handler"
)

func main() {
	config.Init()

	r := gin.Default()

	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	auth := r.Group("/api")
	auth.Use(handler.AuthMiddleware())
	{
		auth.GET("/profile", func(c *gin.Context) {
			username, _ := c.Get("username")
			c.JSON(200, gin.H{"user": username})
		})
	}

	r.Run(":8080")
}
