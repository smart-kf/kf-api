package dao

import (
	"encoding/json"

	"gorm.io/gorm"
)

type KfDomainOrders struct {
	gorm.Model
	CardID      string `json:"cardId" gorm:"column:card_id;type:varchar(255)" doc:"卡密id"`
	OrderNo     string `json:"orderNo" gorm:"column:order_no;unique;type:varchar(255)" doc:"订单号"`
	ToAddress   string `json:"toAddress" gorm:"column:to_address;type:varchar(255)" doc:"接收地址"`
	FromAddress string `json:"fromAddress" gorm:"column:from_address;type:varchar(255)" doc:"支付地址"`
	Price       int64  `json:"price" gorm:"column:price" doc:"支付价格"`                   // 1usdt = 1 * 10e6
	Status      int8   `json:"status" gorm:"column:status" doc:"状态: 1=创建,2=支付成功,3=失败"` // 支付状态
	ConfirmTime int64  `json:"confirmTime" gorm:"column:confirm_time" doc:"确认时间"`      // 支付确认时间
	ExpireTime  int64  `json:"expire_time" gorm:"column:expire_time" doc:"过期时间"`       // 过期时间
	Version     int    `json:"version" gorm:"column:version" doc:"版本号"`
	TradeId     string `json:"tradeId" gorm:"column:trade_id" doc:"区块链交易id"` // 交易支付的id
	DomainList  string `json:"domain_list" gorm:"column:domain_list"`
}

func (KfDomainOrders) TableName() string {
	return "kf_domain_orders"
}

func (o *KfDomainOrders) GetDomainList() []string {
	var res []string
	json.Unmarshal([]byte(o.DomainList), &res)
	return res
}

func (o *KfDomainOrders) SetDomainList(v []string) {
	data, _ := json.Marshal(v)
	o.DomainList = string(data)
}
