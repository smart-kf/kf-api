package dao

import "gorm.io/gorm"

type Orders struct {
	gorm.Model
	CardID         string `json:"cardId" gorm:"column:card_id;type:varchar(255)"`
	PackageId      string `json:"packageId" gorm:"column:package_id;type:varchar(255)"` // 套餐id
	PackageDay     int    `json:"packageDay" gorm:"column:package_day"`                 // 套餐时长
	OrderNo        string `json:"orderNo" gorm:"column:order_no;unique;type:varchar(255)"`
	PayUsdtAddress string `json:"payUsdtAddress" gorm:"column:pay_usdt_address;type:varchar(255)"`
	Price          int64  `json:"price" gorm:"column:price"`                 // 1usdt = 1 * 10e6
	Status         int8   `json:"status" gorm:"column:status"`               // 支付状态
	ConfirmTime    int64  `json:"confirmTime" gorm:"column:confirm_time"`    // 支付确认时间
	ExpireTime     int64  `json:"expire_time" gorm:"column:expire_time"`     // 过期时间
	Ip             string `json:"ip" gorm:"column:ip;type:varchar(255)"`     // ip地址
	Area           string `json:"area" gorm:"column:area;type:varchar(255)"` // ip对应的地区
	Version        int    `json:"version" gorm:"column:version"`
}

func (Orders) TableName() string {
	return "orders"
}
