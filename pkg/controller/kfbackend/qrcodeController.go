package kfbackend

import (
	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"
	"github.com/smart-fm/kf-api/pkg/common"
	"github.com/smart-fm/kf-api/pkg/repository"
)

type QRCodeController struct {
	BaseController
}

type QRCodeRequest struct{}
type QRCodeResponse struct {
	URL           string         `json:"qrcodeUrl,omitempty" doc:"主站二维码图片地址"`
	HealthAt      int64          `json:"healthAt,omitempty" doc:"主站通过健康检查的时间 毫秒"`
	Enable        bool           `json:"enable,omitempty" doc:"启用停用状态"`
	EnableNewUser bool           `json:"enableNewUser,omitempty" doc:"启用停用新粉状态"`
	Domains       []QRCodeDomain `json:"domains,omitempty" doc:"域名列表"`
}

func (c *QRCodeController) List(ctx *gin.Context) {
	var req QRCodeRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	reqCtx := ctx.Request.Context()

	cardID := common.GetKFCardID(ctx)

	var kfsetting repository.KFSettingRepository
	setting, ok, err := kfsetting.GetByCardID(reqCtx, cardID)
	if err != nil {
		xlogger.Error(reqCtx, "查询客服设置失败", xlogger.Err(err), xlogger.Any("cardId", cardID))
		c.Error(ctx, err)
		return
	}

	enable := true
	enableNewUser := true

	if ok {
		enable = setting.QRCodeEnabled
		enableNewUser = setting.QRCodeEnabledNewUser
	}

	// TODO 二维码资源图
	c.Success(ctx, QRCodeResponse{
		URL:           "",
		HealthAt:      0,
		Enable:        enable,
		EnableNewUser: enableNewUser,

		// TODO 计费域名
		Domains: []QRCodeDomain{},
	})
}
