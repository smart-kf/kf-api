package kfbackend

import (
	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/http/middleware"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type SettingController struct {
	BaseController
}

type GetSysConfRequest struct{}

type GetSysConfResponse struct {
	CardID          string `json:"cardId" doc:"卡密id"`
	Nickname        string `json:"nickname" doc:"昵称"`
	AvatarURL       string `json:"avatarUrl" doc:"头像地址"`
	WSFilter        bool   `json:"wsFilter" doc:"开启之后 检测ws行为"`
	WechatFilter    bool   `json:"wechatFilter" doc:"非微信浏览器不能访问"`
	AppleFilter     bool   `json:"appleFilter" doc:"苹果手机过滤器，开启后，只有苹果手机能访问"`
	IPProxyFilter   bool   `json:"iPProxyFilter" doc:"代理ip过滤，开启后，代理ip不能访问"`
	DeviceFilter    bool   `json:"DeviceFilter" doc:"设备异常过滤"`
	SimulatorFilter bool   `json:"simulatorFilter" doc:"模拟器过滤，开启后，模拟器不能访问"`
	Notice          string `json:"notice" doc:"滚动公告"`
	NewMessageVoice bool   `json:"newMessageVoice" doc:"消息提示音"`
}

func (c *SettingController) Get(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	cardID := middleware.GetKFCardID(ctx)

	var kfsetting repository.KFSettingRepository
	setting, ok, err := kfsetting.GetByCardID(reqCtx, cardID)
	if err != nil {
		xlogger.Error(reqCtx, "查询客服设置失败", xlogger.Err(err), xlogger.Any("cardId", cardID))
		c.Error(ctx, err)
		return
	}

	if !ok {
		setting = &dao.KFSettings{}
		setting.CardID = cardID
	}

	sysConf := GetSysConfResponse{
		CardID:          setting.CardID,
		Nickname:        setting.Nickname,
		AvatarURL:       setting.AvatarURL,
		WSFilter:        setting.WSFilter,
		WechatFilter:    setting.WechatFilter,
		AppleFilter:     setting.AppleFilter,
		IPProxyFilter:   setting.IPProxyFilter,
		DeviceFilter:    setting.DeviceFilter,
		SimulatorFilter: setting.SimulatorFilter,
		Notice:          setting.Notice,
		NewMessageVoice: setting.NewMessageVoice,
	}
	c.Success(ctx, sysConf)

}
