package kffrontend

import (
	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
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

func (c *BaseController) GetSmartReplyKeywords(ctx *gin.Context) {
	var (
		reqCtx  = ctx.Request.Context()
		repo    repository.KfWelcomeMessageRepository
		kfToken = common.GetKFToken(reqCtx)
	)

	cardId, err := caches.KfAuthCacheInstance.GetFrontToken(ctx.Request.Context(), kfToken)
	msgs, _, err := repo.List(reqCtx, cardId, constant.SmartMsg, 1, 10)
	if err != nil {
		xlogger.Error(reqCtx, "FindAll failed", xlogger.Err(err))
		c.Error(ctx, err)
		return
	}

	var voList []kffrontend.SmartMsg
	lo.ForEach(
		msgs, func(item *dao.KfWelcomeMessage, index int) {
			voList = append(
				voList, kffrontend.SmartMsg{
					Id:      int64(item.ID),
					Keyword: item.Keyword,
				},
			)
		},
	)

	c.Success(ctx, voList)
}

func (c *BaseController) GetNotice(ctx *gin.Context) {
	var (
		reqCtx  = ctx.Request.Context()
		kfToken = common.GetKFToken(reqCtx)
	)

	cardId, err := caches.KfAuthCacheInstance.GetFrontToken(ctx.Request.Context(), kfToken)
	if err != nil {
		c.Error(ctx, err)
		return
	}

	setting, err := caches.KfSettingCache.GetOne(reqCtx, cardId)
	if err != nil {
		c.Success(
			ctx, gin.H{
				"notice": "",
			},
		)
		return
	}

	c.Success(
		ctx, gin.H{
			"notice": setting.Notice,
		},
	)
}
