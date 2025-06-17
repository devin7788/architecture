package main

import (
	"fmt"
	"time"
)

// 用户信息结构体
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

// 定义用户服务接口（面向抽象）
type UserService interface {
	CreateUser(username, email string) (*User, error)
	GetUserByID(id int64) (*User, error)
	DeactivateUser(id int64) error
}

// 一个简单的内存实现（实际项目中应对接数据库）
type InMemoryUserService struct {
	users  map[int64]*User
	nextID int64
}

// 构造函数
func NewInMemoryUserService() *InMemoryUserService {
	return &InMemoryUserService{
		users:  make(map[int64]*User),
		nextID: 1,
	}
}

// 实现 CreateUser 方法
func (s *InMemoryUserService) CreateUser(username, email string) (*User, error) {
	user := &User{
		ID:        s.nextID,
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
		IsActive:  true,
	}
	s.users[s.nextID] = user
	s.nextID++
	return user, nil
}

// 实现 GetUserByID 方法
func (s *InMemoryUserService) GetUserByID(id int64) (*User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("user with id %d not found", id)
	}
	return user, nil
}

// 实现 DeactivateUser 方法
func (s *InMemoryUserService) DeactivateUser(id int64) error {
	user, exists := s.users[id]
	if !exists {
		return fmt.Errorf("user with id %d not found", id)
	}
	user.IsActive = false
	return nil
}

// 测试入口
func main() {
	service := NewInMemoryUserService()

	user, _ := service.CreateUser("alice", "alice@example.com")
	fmt.Printf("User created: %+v\n", user)

	loaded, _ := service.GetUserByID(user.ID)
	fmt.Printf("User loaded: %+v\n", loaded)

	service.DeactivateUser(user.ID)
	fmt.Printf("User after deactivation: %+v\n", loaded)
}
