package server

import (
	"fmt"
	xlogger "github.com/clearcodecn/log"
	"github.com/clearcodecn/swaggos"
	"github.com/gin-gonic/gin"
	"net/http"
	"std-api/config"
	"std-api/pkg/controller"
	"std-api/pkg/utils"
)

var httpServer http.Server

func Run() error {
	conf := config.GetConfig()

	g := gin.New()
	g.Use(gin.Recovery())

	var logConfig xlogger.GinLogConfigure
	logConfig.LogIP(utils.ClientIP)
	logConfig.SkipPrefix("/static", "/favico.ico")
	if conf.Debug {
		logConfig.EnableRequestBody()
	}
	g.Use(xlogger.GinLog(logConfig))

	g.GET("/healthy", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})

	registerRouter(g)
	swaggerAPI(g)

	return g.Run(fmt.Sprintf("%s:%d", conf.Web.Addr, conf.Web.Port))
}

func registerRouter(g *gin.Engine) {

}

type DemoRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DemoResponse struct {
	Token string `json:"token"`
}

func swaggerAPI(g *gin.Engine) {
	swag := swaggos.Default()
	swag.Response(200, new(controller.BaseResponse))

	swag.JWT("access_token")
	apiGroup := swag.Group("/api")
	apiGroup.Get("/demo").Body(&DemoRequest{}).JSON(&DemoResponse{})

	g.GET("/_doc", gin.WrapH(swag))
	g.GET("/doc", gin.WrapH(swaggos.UI("/_doc", "")))
}
