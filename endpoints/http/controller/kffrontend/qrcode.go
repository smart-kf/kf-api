package kffrontend

import (
	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/http/middleware"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kffrontend"
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
		ctx.String(200, "500 Internal Error")
		return
	}

	if !exist {
		ctx.String(200, "404 Not Found")
		return
	}

	card := qrcodeDomain.CardID
	// TODO:: 判断card的状态.
	// TODO:: 判断域名.
	// 先返回success, 使前端能联调.

	// 1. 获取token，如果没有拿到token，则生成新token，生成新用户返回用户信息.
	kfToken := middleware.GetKFToken(ctx)
	if kfToken == "" {
		// 生成用户信息.
		// token := uuid.New().String()

	}

	_ = card
	// 获取个人信息，如果没有个人信息，渲染个人首页.
	// 用户鉴权、派发用户token、写Cookie 等操作.
	ctx.HTML(200, "front.html", gin.H{})
}
