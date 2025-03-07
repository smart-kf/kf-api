package kffrontend

import (
	"fmt"
	"strings"

	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/factory"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kffrontend"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/pkg/ipinfo"
	"github.com/smart-fm/kf-api/pkg/utils"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

func (c *QRCodeController) Check(ctx *gin.Context) {
	var req kffrontend.QRCodeScanRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	userAgent := ctx.Request.UserAgent()
	lowerUA := strings.ToLower(userAgent)
	ip := utils.ClientIP(ctx)
	reqCtx := ctx.Request.Context()

	ok, qrcode, card := c.getCard(ctx, &req)
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
		c.Error(ctx, xerrors.CheckError.Append("qrcode is disabled"))
		return
	}
	// 2.检测扫码引粉
	if qrcode.Status != constant.QRCodeNormal {
		switch qrcode.Status {
		case constant.QRCodeDisable:
			// 二维码停用.
			xlogger.Info(reqCtx, "禁止访问", xlogger.Any("cause", "二维码停用"))
			c.Error(ctx, xerrors.CheckError.Append("qrcode is disabled"))
			return
		case constant.QRCodeStopGetNewFans:
			if kfToken == "" {
				xlogger.Info(reqCtx, "禁止访问", xlogger.Any("cause", "暂停引新粉"))
				c.Error(ctx, xerrors.CheckError.Append("qrcode stop new fans"))
				return
			}
		}
	}
	// 3. 扫码过滤.
	if (setting.QRCodeScanFilter != constant.QRCodeFilterClose) || setting.IPProxyFilter {
		info, err := ipinfo.Crawl(reqCtx, userAgent, ip)
		// 未出错才拦截.
		if err == nil {
			if setting.IPProxyFilter {
				if info.IsProxy || info.IsVpn {
					c.Error(ctx, xerrors.CheckError.Append("proxy or vpn"))
					return
				}
			}
			switch setting.QRCodeScanFilter {
			case constant.QRCodeFilterRoom:
				if info.IsCloudProvider {
					c.Error(ctx, xerrors.CheckError.Append("cloud provider"))
					return
				}
			case constant.QRCodeFilterNonMainland:
				if !info.IsChina {
					c.Error(ctx, xerrors.CheckError.Append("non mainland"))
					return
				}
			case constant.QRCodeFilterRoomAndNonMainland:
				if info.IsCloudProvider || !info.IsChina {
					c.Error(ctx, xerrors.CheckError.Append("cloud provider or non mainland"))
					return
				}
			}
		}
	}
	// 4. 模拟器过滤
	if setting.SimulatorFilter {
		if strings.Contains(lowerUA, "android") || strings.Contains(
			lowerUA,
			"simulator",
		) {
			c.Error(ctx, xerrors.CheckError.Append("simulator"))
			return
		}
	}
	if setting.AppleFilter {
		if !strings.Contains(lowerUA, "iphone") {
			c.Error(ctx, xerrors.CheckError.Append("not iphone"))
			return
		}
	}
	if setting.WechatFilter {
		if !strings.Contains(lowerUA, "micromessenger") {
			c.Error(ctx, xerrors.CheckError.Append("not wechat"))
			return
		}
	}

	c.scan(ctx, &req)
}

func (c *QRCodeController) scan(ctx *gin.Context, req *kffrontend.QRCodeScanRequest) {
	ok, _, card := c.getCard(ctx, req)
	if !ok {
		return
	}

	reqCtx := ctx.Request.Context()

	cardID := card.CardID

	// 先返回success, 使前端能联调.
	var (
		isNewUser = false
		user      *dao.KfUser
		userRepo  repository.KFUserRepository
	)

	var err error
	// 1. 获取token，如果没有拿到token，则生成新token，生成新用户返回用户信息.
	kfToken := common.GetKFToken(reqCtx)
	if kfToken == "" {
		// 生成用户信息.
		// token := uuid.New().String()
		user = factory.FactoryNewKfUser(int64(card.ID), cardID, ctx.ClientIP())
		info, _ := ipinfo.Crawl(reqCtx, ctx.Request.UserAgent(), utils.ClientIP(ctx))
		user.IP = utils.ClientIP(ctx)
		user.Area = fmt.Sprintf("%s-%s-%s-%s", info.Country, info.Province, info.City, info.Net)
		user.Comments = ""
		user.Mobile = ""
		user.RemarkName = ""
		if info.IsProxy || info.IsVpn {
			user.IsProxy = 1
		}

		if err := userRepo.SaveOne(reqCtx, user); err != nil {
			xlogger.Error(reqCtx, "FindByPath failed", xlogger.Err(err))
			c.Error(ctx, err)
			return
		}
		kfToken = user.UUID
		isNewUser = true
		err = caches.KfAuthCacheInstance.SetFrontToken(reqCtx, kfToken, cardID)
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
			UUID:       user.UUID,
			Avatar:     user.Avatar,
			NickName:   user.NickName,
			WsHost:     config.GetConfig().SocketIO.Host,
			WsFullHost: config.GetConfig().SocketIO.FullHost,
		},
		IsNewUser: isNewUser,
	}

	c.Success(ctx, resp)
}
