package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomCard 生成卡密, 规则如下:
// TM- 开头 [a-zA-Z0-9]{10}
func RandomCard(n int) string {
	// 定义卡密的字符集
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// 生成10位随机字符
	cardKey := "TM-"
	for i := 0; i < n; i++ {
		cardKey += string(charset[rand.Intn(len(charset))])
	}
	// 返回生成的卡密
	return cardKey
}

func Random(n int) string {
	// 定义卡密的字符集
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// 生成10位随机字符
	key := ""
	for i := 0; i < n; i++ {
		key += string(charset[rand.Intn(len(charset))])
	}
	// 返回生成的卡密
	return key
}

func RandomOrderNo() string {
	s := time.Now().Format(`060102150405`)
	s += fmt.Sprintf("%d", rand.Intn(10000)+100)
	return s
}

// RandomPath 生成唯一的path
func RandomPath() string {
	part1 := Random(6)
	part2 := Random(6)
	part3 := Random(6)

	return fmt.Sprintf(
		"%s/%s/%s.html", part1, part2, part3,
	)
}
