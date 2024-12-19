package server

import (
	"fmt"
	xlogger "github.com/clearcodecn/log"
	"github.com/clearcodecn/swaggos"
	"github.com/gin-gonic/gin"
	"std-api/config"
	"std-api/pkg/controller"
	"std-api/pkg/controller/bill"
	"std-api/pkg/controller/kfbackend"
	"std-api/pkg/controller/kffrontend"
	"std-api/pkg/utils"
)

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

	return g.Run(conf.Web.String())
}

func registerRouter(g *gin.Engine) {
	api := g.Group("/api")

	// 计费
	var bc bill.BillController
	bg := api.Group("/bill")
	{
		bg.POST("/login", bc.Login)
	}

	// 客服后台
	kf := api.Group("/kf-be")
	{
		_ = kf
	}

	// 客服前台
	kffe := api.Group("/kf-fe")
	{
		_ = kffe
	}
}

func swaggerAPI(g *gin.Engine) {
	swag := swaggos.Default()
	swag.Response(200, new(controller.BaseResponse))

	swag.JWT("access_token")
	apiGroup := swag.Group("/api")

	// 计费后台的swagger 文档.
	bill.SwaggerDoc(apiGroup)
	kfbackend.SwaggerDoc(apiGroup)
	kffrontend.SwaggerDoc(apiGroup)

	// swagger json 服务
	g.GET("/_doc", gin.WrapH(swag))

	// swagger ui 服务
	g.Any("/doc/*action", gin.WrapH(swaggos.UI("/doc", "http://"+config.GetConfig().Web.String()+"/_doc")))

	fmt.Println("swagger ui: " + "http://" + config.GetConfig().Web.String() + "/doc")
}
