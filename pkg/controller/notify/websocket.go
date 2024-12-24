package notify

import (
	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type NotifyController struct {
	BaseController
}

func (c *NotifyController) WebsocketAuth(ctx *gin.Context) {
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	xlogger.Info(ctx.Request.Context(), "websocketAuth", xlogger.Any("request", string(body)))
	c.Success(ctx, gin.H{})
}
