package main

import (
	"fmt"

	"example.mysql/internal/db"
	"example.mysql/internal/service"
)

func main() {
	db.InitMySQL()

	// 创建用户
	user, err := service.CreateUser("testuser", "secure123")
	if err != nil {
		panic(err)
	}
	fmt.Printf("User created: %+v\n", user)

	// 登录验证
	user2, err := service.GetUserByUsername("testuser")
	if err != nil {
		panic(err)
	}

	if service.VerifyPassword(user2, "secure123") {
		fmt.Println("Login successful")
	} else {
		fmt.Println("Invalid credentials")
	}
}
