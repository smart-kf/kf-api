package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomCard 生成卡密, 规则如下:
// TM- 开头 [a-zA-Z0-9]{10}
func RandomCard() string {
	// 定义卡密的字符集
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-#"
	// 生成10位随机字符
	cardKey := "TM-"
	for i := 0; i < 10; i++ {
		cardKey += string(charset[rand.Intn(len(charset))])
	}
	// 返回生成的卡密
	return cardKey
}
