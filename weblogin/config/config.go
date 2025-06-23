package config

import (
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"weblogin/model"
)

var (
	DB   *gorm.DB
	RDB  *redis.Client
	once sync.Once
)

func Init() {
	once.Do(func() {
		// 初始化MySQL
		dsn := "root:admin123@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
		var err error
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("数据库连接失败:", err)
		}

		// 自动迁移User表
		err = DB.AutoMigrate(&model.User{})
		if err != nil {
			log.Fatal("自动迁移失败:", err)
		}

		// 初始化Redis
		RDB = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "admin123", // no password set
			DB:       0,          // use default DB
		})
		_, err = RDB.Ping(RDB.Context()).Result()
		if err != nil {
			log.Fatal("Redis连接失败:", err)
		}
	})
}
