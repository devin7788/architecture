package tokencache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenCache struct {
	redisClient *redis.Client
	ctx         context.Context
	ttl         time.Duration
}

func NewTokenCache(addr string, password string, db int, ttl time.Duration) *TokenCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &TokenCache{
		redisClient: rdb,
		ctx:         context.Background(),
		ttl:         ttl,
	}
}

// 保存token，key为用户ID，value为token字符串，带过期时间
func (c *TokenCache) SetToken(userID string, token string) error {
	return c.redisClient.Set(c.ctx, userID, token, c.ttl).Err()
}

// 获取token
func (c *TokenCache) GetToken(userID string) (string, error) {
	val, err := c.redisClient.Get(c.ctx, userID).Result()
	if errors.Is(err, redis.Nil) {
		return "", errors.New("token not found")
	}
	return val, err
}

// 删除token
func (c *TokenCache) DeleteToken(userID string) error {
	return c.redisClient.Del(c.ctx, userID).Err()
}
