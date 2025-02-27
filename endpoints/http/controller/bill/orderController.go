package bill

import (
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/http/vo/bill"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
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
	tx := mysql.GetDBFromContext(reqCtx)

	req.OrderBy = "created_at"
	req.Asc = false

	if req.FromAddress != "" {
		tx = tx.Where("from_address = ?", req.FromAddress)
	}
	if req.ToAddress != "" {
		tx = tx.Where("to_address = ?", req.ToAddress)
	}
	if req.TradeId != "" {
		tx = tx.Where("trade_id = ?", req.TradeId)
	}
	if req.OrderNo != "" {
		tx = tx.Where("order_no = ?", req.OrderNo)
	}
	if req.Status != 0 {
		tx = tx.Where("status = ?", req.Status)
	}

	list, cnt, err := common.Paginate[*dao.Orders](tx, &req.PageRequest)
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
