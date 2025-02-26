package kffrontend

import (
	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kffrontend"
	"github.com/smart-fm/kf-api/pkg/ipinfo"
	"github.com/smart-fm/kf-api/pkg/utils"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

type CheckController struct {
	BaseController
}

func (c *CheckController) Check(ctx *gin.Context) {
	var req kffrontend.QRCodeScanRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	ip := utils.ClientIP(ctx)
	reqCtx := ctx.Request.Context()

	ok, qrcode, card := c.getCard(ctx, req)
	if !ok {
		return
	}
	setting, err := caches.KfSettingCache.GetOne(ctx, card.CardID)
	if err != nil {
		c.Error(ctx, xerrors.NewCustomError("卡密不存在"))
		return
	}

	kfToken := common.GetKFToken(reqCtx)

	// 1. 全局开关检测
	if !setting.QRCodeEnabled {
		c.Error(ctx, xerrors.CheckError)
		return
	}
	// 2.检测扫码引粉
	if qrcode.Status != constant.QRCodeNormal {
		switch qrcode.Status {
		case constant.QRCodeDisable:
			// 二维码停用.
			xlogger.Info(reqCtx, "禁止访问", xlogger.Any("cause", "二维码停用"))
			c.Error(ctx, xerrors.CheckError)
			return
		case constant.QRCodeStopGetNewFans:
			if kfToken == "" {
				xlogger.Info(reqCtx, "禁止访问", xlogger.Any("cause", "暂停引新粉"))
				c.Error(ctx, xerrors.CheckError)
				return
			}
		}
	}
	// 3. 扫码过滤.
	if setting.QRCodeScanFilter != constant.QRCodeFilterClose {
		info, err := ipinfo.Crawl(reqCtx, ctx.Request.UserAgent(), ip)
		// 未出错才拦截.
		if err == nil {
			switch setting.QRCodeScanFilter {
			case constant.QRCodeFilterRoom:
				if info.IsCloudProvider {
					c.Error(ctx, xerrors.CheckError)
					return
				}
			case constant.QRCodeFilterNonMainland:
				if !info.IsChina {
					c.Error(ctx, xerrors.CheckError)
					return
				}
			case constant.QRCodeFilterRoomAndNonMainland:
				if info.IsCloudProvider || !info.IsChina {
					c.Error(ctx, xerrors.CheckError)
					return
				}
			}
		}
	}
}
