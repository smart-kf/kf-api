package kfbackend

import (
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/service/orders"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
)

type DomainOrderController struct {
	BaseController
}

func (c *DomainOrderController) GetDomainPrice(ctx *gin.Context) {
	price := caches.BillSettingCacheInstance.GetDomainPrice()
	c.Success(ctx, price)
}

func (c *DomainOrderController) CreateOrder(ctx *gin.Context) {
	var req kfbackend.DomainOrderRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	reqCtx := ctx.Request.Context()

	rsp, err := orders.CreateDomainOrder(
		reqCtx, orders.CreateDomainDTO{
			FromAddress: req.PayAddress,
		},
	)

	if err != nil {
		c.Error(ctx, err)
		return
	}

	c.Success(ctx, rsp)
}

func (c *DomainOrderController) OrderList(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	rsp, err := orders.ListOrder(
		reqCtx, common.GetKFCardID(ctx),
	)

	if err != nil {
		c.Error(ctx, err)
		return
	}

	c.Success(ctx, rsp)
}
