package constant

import "time"

type SaleStatus int8

const (
	SaleStatusOnSale  = iota + 1 // 销售中
	SaleStatusOffSale            // 下架
	SaleStatusSold               // 已售出
)

type LoginStatus int8

const (
	LoginStatusUnLogin  = iota + 1 // 未登录过，未使用
	LoginStatusLoginned            // 已登录过，已使用
)

type CardType int8

const (
	CardTypeNormal  = iota + 1 // 正式卡
	CardTypeTesting            // 测试卡
)

const CardTimeExpire = 15 * time.Minute // 测试卡片有效期15分钟.

const (
	WelcomeMsg = "welcome_msg"
	QuickReply = "quick_reply"
	SmartReply = "smart_reply"
)
