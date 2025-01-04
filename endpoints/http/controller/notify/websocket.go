package notify

import (
	"errors"

	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/endpoints/http/middleware"
)

type NotifyController struct {
	BaseController
}

type WebsocketAuthRequest struct {
	Token     string `json:"token"`
	Ip        string `json:"ip"`
	Platform  string `json:"platform"`
	UserAgent string `json:"userAgent"`
}

func (r WebsocketAuthRequest) GetPlatform() string {
	return r.Platform
}

func (r WebsocketAuthRequest) GetToken() string {
	return r.Token
}

func (r WebsocketAuthRequest) GetUserAgent() string {
	return r.UserAgent
}

func (r WebsocketAuthRequest) GetIP() string {
	return r.Ip
}

func (c *NotifyController) WebsocketAuth(ctx *gin.Context) {
	var req WebsocketAuthRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	var (
		err error
	)
	switch req.GetPlatform() {
	case "kf-backend": // 后台
		_, err = middleware.VerifyKFBackendToken(req.GetToken())
	case "kf": // 前台
		_, err = middleware.VerifyKFToken(req.GetToken())
	default:
		err = errors.New("invalid platform")
	}
	if err != nil {
		xlogger.Info(ctx.Request.Context(), "websocket auth failed: "+err.Error())
		c.Error(ctx, err)
		return
	}
	xlogger.Info(ctx.Request.Context(), "websocket auth success:"+req.GetIP())
	c.Success(ctx, gin.H{})
}
