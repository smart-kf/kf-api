package kffrontend

import (
	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/http/middleware"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kffrontend"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

type BaseController struct {
	middleware.BaseController
}

func NewBaseController() *BaseController {
	return &BaseController{}
}

func (c *BaseController) getCard(ctx *gin.Context, req *kffrontend.QRCodeScanRequest) (
	bool,
	*dao.KFQRCodeDomain, *dao.KFCard,
) {
	reqCtx := ctx.Request.Context()

	var qrcodeDomainRepo repository.QRCodeDomainRepository
	qrcodeDomain, exist, err := qrcodeDomainRepo.FindByPath(reqCtx, req.Code)
	if err != nil {
		xlogger.Error(reqCtx, "FindByPath failed", xlogger.Err(err))
		c.Error(ctx, err)
		return false, nil, nil
	}

	if !exist {
		c.Error(ctx, xerrors.NewCustomError("二维码已失效"))
		return false, nil, nil
	}

	cardID := qrcodeDomain.CardID
	card, err := caches.KfCardCacheInstance.GetCardByID(reqCtx, cardID)
	if err != nil {
		c.Error(ctx, err)
		return false, nil, nil
	}

	return true, qrcodeDomain, card
}
