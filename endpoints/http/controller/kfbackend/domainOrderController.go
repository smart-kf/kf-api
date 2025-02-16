package kfbackend

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
	"github.com/smart-fm/kf-api/pkg/xerrors"
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

	reqCtx := ctx.Request.Context()

	var domainRepo repository.BillDomainRepository
	cnt, err := domainRepo.CountPrivateDomain(reqCtx)
	if err != nil {
		c.Error(ctx, err)
		return
	}

	if cnt < req.Num {
		c.Error(ctx, xerrors.NewParamsErrors(fmt.Sprintf("域名库存不足，本次最多可购买: %d 个", req.Num)))
		return
	}
}
