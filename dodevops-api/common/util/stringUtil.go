package util

import (
	"crypto/rand"
	"encoding/base64"
	"strconv"
)

// StringToUint 将字符串转换为uint
func StringToUint(s string) uint {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return uint(i)
}

// StringToInt 将字符串转换为int
func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// GenerateRandomString 生成安全的随机字符串
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err) // 随机数生成失败是严重错误
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}
