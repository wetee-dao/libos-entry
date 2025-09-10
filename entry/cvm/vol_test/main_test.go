package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"testing"
)

func TestKey(t *testing.T) {
	// 生成随机密钥
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("生成随机密钥失败: %v", err)
	}

	bt, err := os.ReadFile("/run/tee-vol/luks_key.bin")
	if err != nil {
		t.Fatalf("读取密钥文件失败: %v", err)
	}
	if len(bt) != 32 {
		t.Fatalf("密钥文件长度错误")
	}
	fmt.Println(bt)
}
