package bill

import (
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/http/vo/bill"
)

type OrderController struct {
	BaseController
}

func (c *OrderController) List(ctx *gin.Context) {
	var req bill.ListOrderRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	reqCtx := ctx.Request.Context()

	req.OrderBy = "created_at"
	req.Asc = false

	repo := repository.BillOrderRepository{}
	list, cnt, err := repo.ListOrder(
		reqCtx, repository.ListOrderOptions{
			PageRequest: req.PageRequest,
			OrderNo:     req.OrderNo,
			TradeId:     req.TradeId,
			FromAddress: req.FromAddress,
			ToAddress:   req.ToAddress,
			Status:      req.Status,
		},
	)

	if err != nil {
		c.Error(ctx, err)
		return
	}

	c.Success(
		ctx, bill.ListOrderResponse{
			List:  list,
			Total: cnt,
		},
	)
}
