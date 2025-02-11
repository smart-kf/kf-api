package dao

import "gorm.io/gorm"

type Orders struct {
	gorm.Model
	CardID      string `json:"cardId" gorm:"column:card_id;type:varchar(255)" doc:"卡密id"`
	PackageId   string `json:"packageId" gorm:"column:package_id;type:varchar(255)" doc:"套餐id"` // 套餐id
	PackageDay  int    `json:"packageDay" gorm:"column:package_day" doc:"套餐天数"`                 // 套餐时长
	OrderNo     string `json:"orderNo" gorm:"column:order_no;unique;type:varchar(255)" doc:"订单号"`
	ToAddress   string `json:"toAddress" gorm:"column:to_address;type:varchar(255)" doc:"接收地址"`
	FromAddress string `json:"fromAddress" gorm:"column:from_address;type:varchar(255)" doc:"支付地址"`
	Price       int64  `json:"price" gorm:"column:price" doc:"支付价格"`                      // 1usdt = 1 * 10e6
	Status      int8   `json:"status" gorm:"column:status" doc:"状态: 1=创建,2=支付成功,3=失败"`    // 支付状态
	ConfirmTime int64  `json:"confirmTime" gorm:"column:confirm_time" doc:"确认时间"`         // 支付确认时间
	ExpireTime  int64  `json:"expire_time" gorm:"column:expire_time" doc:"过期时间"`          // 过期时间
	Ip          string `json:"ip" gorm:"column:ip;type:varchar(255)" doc:"下单ip地址"`        // ip地址
	Area        string `json:"area" gorm:"column:area;type:varchar(255)" doc:"下单ip对应的地区"` // ip对应的地区
	Version     int    `json:"version" gorm:"column:version" doc:"版本号"`
	Email       string `json:"email" doc:"邮箱地址"`                             // 邮箱
	TradeId     string `json:"tradeId" gorm:"column:trade_id" doc:"区块链交易id"` // 交易支付的id
}

func (Orders) TableName() string {
	return "orders"
}
