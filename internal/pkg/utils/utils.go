package utils

import (
	"strings"
	"time"
)

// 生成文章Slug
func GenerateSlug(title string) string {
	t := strings.ToLower(title)
	t = strings.ReplaceAll(t, " ", "-")
	// 可进一步加随机字符串保证唯一性
	return t + "-" + RandString(6)
}

// 生成随机 6 位字符串（简单示例）
func RandString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[RandInt(len(letters))]
	}
	return string(b)
}

// 简单随机数生成
func RandInt(max int) int {
	return int(time.Now().UnixNano() % int64(max))
}
