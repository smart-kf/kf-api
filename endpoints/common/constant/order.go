package constant

const (
	OrderStatusCreated = iota + 1 // 已创建
	OrderStatusPay                // 已支付
	OrderStatusCancel             // 已取消
)

const (
	OrderExpireZSetKey = "kf.order.expire"
)
