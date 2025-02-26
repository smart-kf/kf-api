package kffrontend

import (
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/endpoints/http/vo/kffrontend"
)

type CheckController struct {
	BaseController
}

func (c *CheckController) Check(ctx *gin.Context) {
	var req kffrontend.QRCodeScanRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	// reqCtx := ctx.Request.Context()
	//
	// ok, qrcode, card := c.getCard(ctx, req)
	// if !ok {
	// 	return
	// }
	// setting, err := caches.KfSettingCache.GetOne(ctx, card.CardID)
	// if err != nil {
	// 	c.Error(ctx, xerrors.NewCustomError("卡密不存在"))
	// 	return
	// }
	//
	// kfToken := common.GetKFToken(reqCtx)
	//
	// // 1.检测扫码引粉
	// if qrcode.Status != dao.QRCodeStatusOK {
	// 	switch qrcode.Status {
	// 	case dao.QRCodeStatusDisable:
	// 		// 二维码停用.
	// 		xlogger.Info(reqCtx, "禁止访问", xlogger.Any("cause", "二维码停用"))
	// 		c.Error(ctx, xerrors.CheckError)
	// 		return
	// 	case dao.QRCodeStatusStop:
	// 		if kfToken == "" {
	// 			xlogger.Info(reqCtx, "禁止访问", xlogger.Any("cause", "暂停引新粉"))
	// 			c.Error(ctx, xerrors.CheckError)
	// 			return
	// 		}
	// 	}
	// }
}
