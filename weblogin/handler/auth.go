package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"weblogin/config"
	"weblogin/model"
	"weblogin/utils"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 密码hash
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	user := model.User{
		Username: req.Username,
		Password: string(hash),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	token, err := utils.GenerateToken(user.Username, utils.AccessToken, time.Hour*1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	// 缓存token到Redis，24小时过期
	if err := utils.SetTokenCache(token, user.Username, 24*time.Hour); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "缓存token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// 中间件验证token有效性
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "缺少token"})
			return
		}

		claims, err := utils.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token无效"})
			return
		}

		username, err := utils.GetUsernameByToken(token)
		if err != nil || username != claims.Username {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token验证失败"})
			return
		}

		// 设置上下文用户信息
		c.Set("username", claims.Username)
		c.Next()
	}
}
