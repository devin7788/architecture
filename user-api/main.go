package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义注册请求体结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// POST /user/register 处理函数
func registerHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数不合法"})
		return
	}

	// 这里可以加入用户名是否存在判断、数据库保存等逻辑
	// 示例中简化为直接返回成功
	c.JSON(http.StatusOK, gin.H{
		"message": "register success",
		"user":    req.Username,
	})
}

func main() {
	r := gin.Default()
	r.POST("/user/register", registerHandler)
	r.Run(":8080")
}
