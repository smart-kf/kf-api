package bill

import (
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type ListOrderRequest struct {
	common.PageRequest
	OrderNo     string `json:"orderNo" doc:"订单id"`
	TradeId     string `json:"tradeId" doc:"区块链上交易id"`
	FromAddress string `json:"fromAddress" doc:"客户地址"`
	ToAddress   string `json:"toAddress" doc:"接收地址"`
	Status      int    `json:"status" doc:"1=等待支付,2=支付成功,3=已取消"`
}

type ListOrderResponse struct {
	List  []*dao.Orders `json:"list"`
	Total int64         `json:"total"`
}
