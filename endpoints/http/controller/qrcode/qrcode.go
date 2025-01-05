package qrcode

import (
	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/domain/repository"
)

type QRCodeController struct{}

func (QRCodeController) Scan(ctx *gin.Context) {
	path := ctx.Param("action")
	reqCtx := ctx.Request.Context()

	var qrcodeDomainRepo repository.QRCodeDomainRepository
	qrcodeDomain, exist, err := qrcodeDomainRepo.FindByPath(reqCtx, path)
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

	_ = card
	// 获取个人信息，如果没有个人信息，渲染个人首页.
	// 用户鉴权、派发用户token、写Cookie 等操作.
	ctx.HTML(200, "front.html", gin.H{})
}
