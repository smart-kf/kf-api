package billfront

type OrderNotifyRequest struct {
	TradeId   string `json:"trade_id"`
	OrderId   string `json:"order_id"`
	Status    int    `json:"status"`    // 1=wait，2=success，3=fail
	Timestamp int64  `json:"timestamp"` // 确认时间
	Address   string `json:"address"`
}
