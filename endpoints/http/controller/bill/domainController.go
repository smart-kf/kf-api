package bill

import (
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/http/vo/bill"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

type DomainController struct {
	BaseController
}

func (c *DomainController) AddDomain(ctx *gin.Context) {
	var req bill.AddDomainRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	reqCtx := ctx.Request.Context()

	// 1. 判重.
	var domainRepo repository.BillDomainRepository
	_, exist, err := domainRepo.FindByTopName(reqCtx, req.TopName)
	if err != nil {
		c.Error(ctx, err)
		return
	}

	if exist {
		c.Error(ctx, xerrors.NewCustomError("域名已存在，请勿重复添加"))
		return
	}

	// 2. 添加域名
	var domain = dao.BillDomain{
		TopName:  req.TopName,
		IsPublic: *req.IsPublic,
		IsBind:   false,
		Status:   req.Status,
	}

	tx := mysql.DB()
	if err := tx.Create(&domain).Error; err != nil {
		c.Error(ctx, err)
		return
	}

	c.Success(ctx, domain)
}

func (c *DomainController) ListDomain(ctx *gin.Context) {
	var req bill.ListDomainRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	reqCtx := ctx.Request.Context()
	var repo repository.BillDomainRepository

	list, count, err := repo.List(
		reqCtx, &repository.ListDomainOption{
			PageRequest: &req.PageRequest,
		},
	)
	if err != nil {
		c.Error(ctx, err)
		return
	}
	var rsp bill.ListDomainResponse
	for _, item := range list {
		rsp.List = append(
			rsp.List, &bill.BillDomainResponse{
				ID:         int64(item.ID),
				TopName:    item.TopName,
				CreateTime: item.CreatedAt.Unix(),
				IsPublic:   item.IsPublic,
				IsBind:     item.IsBind,
				Status:     item.Status,
			},
		)
	}
	rsp.Total = count
	c.Success(ctx, rsp)
}

func (c *DomainController) DeleteDomain(ctx *gin.Context) {
	var req bill.DeleteDomainRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	reqCtx := ctx.Request.Context()
	var repo repository.BillDomainRepository
	data, exist, err := repo.FindByID(reqCtx, req.ID)
	if err != nil {
		c.Error(ctx, err)
		return
	}
	if !exist {
		c.Error(ctx, xerrors.NewParamsErrors("数据不存在"))
		return
	}
	if data.IsBind {
		c.Error(ctx, xerrors.NewCustomError("域名已经被卡密绑定，请解绑在操作"))
		return
	}

	err = repo.DeleteByID(ctx, req.ID)
	if err != nil {
		c.Error(ctx, err)
		return
	}
	c.Success(ctx, nil)
}
