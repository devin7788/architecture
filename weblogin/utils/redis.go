package utils

import (
	"context"
	"time"

	"weblogin/config"
)

var ctx = context.Background()

func SetTokenCache(token string, username string, expiration time.Duration) error {
	return config.RDB.Set(ctx, token, username, expiration).Err()
}

func GetUsernameByToken(token string) (string, error) {
	return config.RDB.Get(ctx, token).Result()
}

func DeleteToken(token string) error {
	return config.RDB.Del(ctx, token).Err()
}
