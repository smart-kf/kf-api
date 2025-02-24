package kfbackend

import (
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
)

type DomainOrderController struct {
	BaseController
}

func (c *DomainOrderController) GetDomainPrice(ctx *gin.Context) {
	c.Success(ctx, config.GetConfig().DomainPrice)
}

func (c *DomainOrderController) CreateOrder(ctx *gin.Context) {
	var req kfbackend.DomainOrderRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	// reqCtx := ctx.Request.Context()
	//
	// var domainRepo repository.BillDomainRepository
	// cnt, err := domainRepo.CountPrivateDomain(reqCtx)
	// if err != nil {
	// 	c.Error(ctx, err)
	// 	return
	// }
	// tx, newCtx := mysql.Begin(reqCtx)
}
