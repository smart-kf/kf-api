package notify

import (
	"errors"
	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"
	"github.com/smart-fm/kf-api/pkg/common"
)

type NotifyController struct {
	BaseController
}

type WebsocketAuthRequest struct {
	Header struct {
		UserAgent []string `json:"User-Agent"`
	} `json:"header"`
	Query struct {
		Platform []string `json:"platform"`
		Token    []string `json:"token"`
	} `json:"query"`
	RemoteAddress string `json:"remoteAddress"`
}

func (r WebsocketAuthRequest) GetPlatform() string {
	if len(r.Query.Platform) > 0 {
		return r.Query.Platform[0]
	}
	return ""
}

func (r WebsocketAuthRequest) GetToken() string {
	if len(r.Query.Token) > 0 {
		return r.Query.Token[0]
	}
	return ""
}

func (r WebsocketAuthRequest) GetUserAgent() string {
	if len(r.Header.UserAgent) > 0 {
		return r.Header.UserAgent[0]
	}
	return ""
}

func (r WebsocketAuthRequest) GetIP() string {
	return r.RemoteAddress
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
		_, err = common.VerifyKFBackendToken(req.GetToken())
	case "kf": // 前台
		_, err = common.VerifyKFToken(req.GetToken())
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
