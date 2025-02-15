package kffrontend

import (
	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/factory"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kffrontend"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

type QRCodeController struct {
	BaseController
}

func (c *QRCodeController) Scan(ctx *gin.Context) {
	var req kffrontend.QRCodeScanRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	reqCtx := ctx.Request.Context()

	var qrcodeDomainRepo repository.QRCodeDomainRepository
	qrcodeDomain, exist, err := qrcodeDomainRepo.FindByPath(reqCtx, req.Code)
	if err != nil {
		xlogger.Error(reqCtx, "FindByPath failed", xlogger.Err(err))
		c.Error(ctx, err)
		return
	}

	if !exist {
		c.Error(ctx, xerrors.NewCustomError("二维码已失效"))
		return
	}

	cardID := qrcodeDomain.CardID
	card, err := caches.KfCardCacheInstance.GetCardByID(reqCtx, cardID)
	if err != nil {
		c.Error(ctx, err)
		return
	}

	// TODO:: 判断card的状态.
	// TODO:: 判断域名.
	// 先返回success, 使前端能联调.
	var (
		isNewUser = false
		user      *dao.KfUser
		userRepo  repository.KFUserRepository
	)
	// 1. 获取token，如果没有拿到token，则生成新token，生成新用户返回用户信息.
	kfToken := common.GetKFToken(reqCtx)
	if kfToken == "" {
		// 生成用户信息.
		// token := uuid.New().String()
		user = factory.FactoryNewKfUser(int64(card.ID), cardID, ctx.ClientIP())
		if err := userRepo.SaveOne(reqCtx, user); err != nil {
			xlogger.Error(reqCtx, "FindByPath failed", xlogger.Err(err))
			c.Error(ctx, err)
			return
		}
		isNewUser = true
		err := caches.KfAuthCacheInstance.SetFrontToken(reqCtx, kfToken, cardID)
		if err != nil {
			xlogger.Error(reqCtx, "SetFrontToken failed", xlogger.Err(err))
			c.Error(ctx, err)
			return
		}
	} else {
		user, err = caches.KfUserCacheInstance.GetDBUser(ctx, cardID, kfToken)
		if err != nil {
			xlogger.Error(reqCtx, "GetDBUser-failed", xlogger.Err(err))
			c.Error(ctx, err)
			return
		}
	}
	resp := kffrontend.QRCodeScanResponse{
		UserInfo: kffrontend.KFUserInfo{
			UUID:     user.UUID,
			Avatar:   user.Avatar,
			NickName: user.NickName,
		},
		IsNewUser: isNewUser,
	}

	// 获取个人信息，如果没有个人信息，渲染个人首页.
	// 用户鉴权、派发用户token、写Cookie 等操作.
	c.Success(ctx, resp)
}
