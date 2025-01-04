package billfront

type OrderNotifyRequest struct {
	OrderNo string `json:"orderNo" doc:"订单号" binding:"required" validate:"required"`
}
