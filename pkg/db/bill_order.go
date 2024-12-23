package db

import "gorm.io/gorm"

type Orders struct {
	gorm.Model
	CardID         string `json:"cardId" gorm:"column:card_id;unique"`
	OrderNo        string `json:"orderNo" gorm:"column:order_no;unique"`
	PayUsdtAddress string `json:"payUsdtAddress" gorm:"column:pay_usdt_address"`
	Price          int64  `json:"price" gorm:"column:price"`              // 1usdt = 1 * 10e6
	Status         int8   `json:"status" gorm:"column:status"`            // 支付状态
	ConfirmTime    int64  `json:"confirmTime" gorm:"column:confirm_time"` // 创建时间
	ExpireTime     int64  `json:"expire_time" gorm:"column:expire_time"`  // 过期时间
	Ip             string `json:"ip" gorm:"column:ip"`                    // ip地址
	Area           string `json:"area" gorm:"column:area"`                // ip对应的地区
	Version        int    `json:"version" gorm:"column:version"`
}

func (Orders) TableName() string {
	return "orders"
}
