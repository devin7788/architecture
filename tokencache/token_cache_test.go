package tokencache

import (
	"testing"
	"time"
)

func TestTokenCache(t *testing.T) {
	// 请保证本地/测试环境Redis启动且地址正确
	tc := NewTokenCache("localhost:6379", "admin123", 0, 2*time.Second)

	userID := "user123"
	token := "token_abc123"

	// 测试写入
	if err := tc.SetToken(userID, token); err != nil {
		t.Fatalf("SetToken error: %v", err)
	}

	// 测试读取
	got, err := tc.GetToken(userID)
	if err != nil {
		t.Fatalf("GetToken error: %v", err)
	}
	if got != token {
		t.Fatalf("Expected token %s, got %s", token, got)
	}

	// 测试过期（等待超过TTL）
	time.Sleep(3 * time.Second)
	_, err = tc.GetToken(userID)
	if err == nil {
		t.Fatalf("Expected token expired error, got nil")
	}

	// 再次写入，测试删除
	if err := tc.SetToken(userID, token); err != nil {
		t.Fatalf("SetToken error: %v", err)
	}
	if err := tc.DeleteToken(userID); err != nil {
		t.Fatalf("DeleteToken error: %v", err)
	}
	_, err = tc.GetToken(userID)
	if err == nil {
		t.Fatalf("Expected token not found after delete, got nil")
	}
}
