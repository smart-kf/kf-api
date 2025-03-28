package kfbackend

type DomainOrderRequest struct {
	PayAddress string `json:"payAddress" doc:"支付地址" binding:"required"`
}

type DomainOrder struct {
	OrderNo     string `json:"orderNo" gorm:"column:order_no;unique;type:varchar(255)" doc:"订单号"`
	ToAddress   string `json:"toAddress" gorm:"column:to_address;type:varchar(255)" doc:"接收地址"`
	FromAddress string `json:"fromAddress" gorm:"column:from_address;type:varchar(255)" doc:"支付地址"`
	Price       int64  `json:"price" gorm:"column:price" doc:"支付价格"`                   // 1usdt = 1 * 10e6
	Status      int8   `json:"status" gorm:"column:status" doc:"状态: 1=创建,2=支付成功,3=失败"` // 支付状态
	ConfirmTime int64  `json:"confirmTime" gorm:"column:confirm_time" doc:"确认时间"`      // 支付确认时间
	ExpireTime  int64  `json:"expire_time" gorm:"column:expire_time" doc:"过期时间"`       // 过期时间
	Domain      string `json:"domain" gorm:"column:domain" doc:"域名"`                   // 域名地址
	TradeId     string `json:"tradeId" gorm:"column:trade_id" doc:"区块链交易id"`           // 交易支付的id
	PayUrl      string `json:"payUrl" gorm:"column:pay_url"`                           // 支付地址
}

type DomainOrderList struct {
	List []*DomainOrder `json:"list" doc:"订单列表"`
}
