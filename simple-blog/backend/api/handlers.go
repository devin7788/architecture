package api

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yuin/goldmark"
	"gorm.io/gorm"
)

var jwtKey = []byte("secret_key_123")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func generateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		c.Set("username", claims.Username)
		c.Next()
	}
}

// 简单登录接口，用户名密码写死 admin/admin
func LoginHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid params"})
		return
	}
	if req.Username == "admin" && req.Password == "admin" {
		token, err := generateToken(req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
}

func CreateArticleHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid params"})
			return
		}
		// Markdown转HTML
		var buf bytes.Buffer
		if err := goldmark.Convert([]byte(req.Content), &buf); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "markdown parse error"})
			return
		}

		article := Article{
			Title:       req.Title,
			Content:     req.Content,
			HTMLContent: buf.String(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			AuthorID:    1, // 简化写死
		}

		if err := db.Create(&article).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db insert failed"})
			return
		}
		c.JSON(http.StatusOK, article)
	}
}

func ListArticlesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var articles []Article
		if err := db.Order("created_at desc").Find(&articles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
			return
		}
		c.JSON(http.StatusOK, articles)
	}
}

func GetArticleHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var article Article
		if err := db.First(&article, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
			return
		}
		c.JSON(http.StatusOK, article)
	}
}

func DeleteArticleHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&Article{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "deleted"})
	}
}
