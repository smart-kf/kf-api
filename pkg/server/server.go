package server

import (
	"fmt"
	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"std-api/config"
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

	return g.Run(fmt.Sprintf("%s:%d", conf.Web.Addr, conf.Web.Port))
}

func registerRouter(g *gin.Engine) {

}
