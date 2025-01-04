package kfbackend

import (
	"fmt"

	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"

	repository2 "github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/http/middleware"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
	"github.com/smart-fm/kf-api/pkg/utils"
)

type QRCodeController struct {
	BaseController
}

type QRCodeRequest struct{}
type QRCodeResponse struct {
	URL           string                   `json:"qrcodeUrl,omitempty" doc:"主站二维码图片地址"`
	HealthAt      int64                    `json:"healthAt,omitempty" doc:"主站通过健康检查的时间 毫秒"`
	Enable        bool                     `json:"enable,omitempty" doc:"启用停用状态"`
	EnableNewUser bool                     `json:"enableNewUser,omitempty" doc:"启用停用新粉状态"`
	Domains       []kfbackend.QRCodeDomain `json:"domains,omitempty" doc:"域名列表"`
}

func (c *QRCodeController) List(ctx *gin.Context) {
	var req QRCodeRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	reqCtx := ctx.Request.Context()

	cardID := middleware.GetKFCardID(ctx)

	var kfsetting repository2.KFSettingRepository
	setting, ok, err := kfsetting.GetByCardID(reqCtx, cardID)
	if err != nil {
		xlogger.Error(reqCtx, "查询客服设置失败", xlogger.Err(err), xlogger.Any("cardId", cardID))
		c.Error(ctx, err)
		return
	}

	enable := true
	enableNewUser := true
	ver := 0

	if ok {
		enable = setting.QRCodeEnabled
		enableNewUser = setting.QRCodeEnabledNewUser
		ver = setting.QRCodeVersion
	}

	baseDomain := "base.domain"                          // TODO 配置
	chatH5 := fmt.Sprintf("https://%s/todo", baseDomain) // TODO 配置 前端客服聊天C端入口地址
	if ver > 0 {
		chatH5 = fmt.Sprintf("%s?ver=%d", chatH5, ver)
	}

	static, err := utils.DrawQRCodeNX(cardID, chatH5)
	if err != nil {
		c.Error(ctx, err)
		return
	}

	c.Success(
		ctx, QRCodeResponse{
			URL:           fmt.Sprintf("https://%s/%s", baseDomain, static),
			HealthAt:      0,
			Enable:        enable,
			EnableNewUser: enableNewUser,

			// TODO 计费域名
			Domains: []kfbackend.QRCodeDomain{},
		},
	)
}

type QRCodeSwitchRequest struct{}
type QRCodeSwitchResponse struct {
	URL      string `json:"qrcodeUrl,omitempty" doc:"主站二维码图片地址"`
	HealthAt int64  `json:"healthAt,omitempty" doc:"主站通过健康检查的时间 毫秒"`
}

// Switch 更换二维码图片
func (c *QRCodeController) Switch(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	cardID := middleware.GetKFCardID(ctx)

	var kfsetting repository2.KFSettingRepository
	setting, ok, err := kfsetting.GetByCardID(reqCtx, cardID)
	if err != nil {
		xlogger.Error(reqCtx, "查询客服设置失败", xlogger.Err(err), xlogger.Any("cardId", cardID))
		c.Error(ctx, err)
		return
	}

	if !ok {
		setting.CardID = cardID
	}

	// 更新版本号
	setting.QRCodeVersion++

	err = kfsetting.SaveOne(reqCtx, setting)
	if err != nil {
		xlogger.Error(reqCtx, "保存客服设置失败", xlogger.Err(err), xlogger.Any("cardId", cardID))
		c.Error(ctx, err)
		return
	}

	baseDomain := "base.domain" // TODO 配置
	chatH5 := fmt.Sprintf(
		"https://%s/todo?ver=%d",
		baseDomain,
		setting.QRCodeVersion, // 使用新的版本号生成新的二维码
	) // TODO 配置 前端客服聊天C端入口地址

	resource, err := utils.DrawQRCodeNX(cardID, chatH5)
	if err != nil {
		c.Error(ctx, err)
		return
	}

	c.Success(
		ctx, QRCodeSwitchResponse{
			URL:      fmt.Sprintf("https://%s/%s", baseDomain, resource),
			HealthAt: 0,
		},
	)
}
