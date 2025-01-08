package http

import (
	"fmt"

	xlogger "github.com/clearcodecn/log"
	"github.com/clearcodecn/swaggos"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/captcha"

	"github.com/smart-fm/kf-api/endpoints/http/controller/bill"
	bill2 "github.com/smart-fm/kf-api/endpoints/http/controller/billfrontend"
	dev2 "github.com/smart-fm/kf-api/endpoints/http/controller/dev"
	"github.com/smart-fm/kf-api/endpoints/http/controller/kfbackend"
	"github.com/smart-fm/kf-api/endpoints/http/controller/kffrontend"
	notify2 "github.com/smart-fm/kf-api/endpoints/http/controller/notify"
	"github.com/smart-fm/kf-api/endpoints/http/controller/qrcode"
	"github.com/smart-fm/kf-api/pkg/version"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/endpoints/http/middleware"
	"github.com/smart-fm/kf-api/pkg/utils"
)

func Run() error {
	conf := config.GetConfig()
	g := gin.New()
	g.Use(gin.Recovery())
	g.Use(cors.Default())

	var logConfig xlogger.GinLogConfigure
	logConfig.LogIP(utils.ClientIP)
	logConfig.SkipPrefix("/static", "/favico.ico")
	if conf.Debug {
		logConfig.EnableRequestBody()
	}
	g.Use(xlogger.GinLog(logConfig))

	g.GET(
		"/healthy", func(ctx *gin.Context) {
			ctx.String(200, "success")
		},
	)
	g.GET(
		"/version", func(ctx *gin.Context) {
			ctx.String(200, fmt.Sprintf("git-version is: %s", version.Version))
		},
	)

	g.Static("/static", conf.Web.StaticDir)

	registerRouter(g)
	swaggerAPI(g)

	return g.Run(conf.Web.String())
}

func registerRouter(g *gin.Engine) {
	api := g.Group("/api")

	public := api.Group("/public")
	{
		var publicController kfbackend.PublicController
		public.GET("/captchaId", publicController.GetCaptchaId)           // 获取验证码id
		public.GET("/captcha/*action", gin.WrapH(captcha.Server(80, 40))) // 显示验证码图片
	}

	// 计费
	var bc bill.BaseController
	bgUnAuth := api.Group("/bill")
	{
		bgUnAuth.POST("/login", bc.Login)
	}

	bgAuth := api.Group("/bill", middleware.BillAuthMiddleware())
	{
		var cardController bill.CardController
		cardGroup := bgAuth.Group("/card")
		{
			cardGroup.POST("/batch-add", cardController.BatchAddCard)
			cardGroup.POST("/updateStatus", cardController.UpdateStatus)
			cardGroup.POST("/list", cardController.List)
		}

		var domainController bill.DomainController
		domainGroup := bgAuth.Group("/domain")
		{
			domainGroup.POST("/add", domainController.AddDomain)
			domainGroup.POST("/list", domainController.ListDomain)
			domainGroup.POST("/del", domainController.DeleteDomain)
		}
	}

	var authController kfbackend.AuthController
	api.POST("/kf-be/login", authController.Login)

	// 客服后台
	kf := api.Group("/kf-be", middleware.KFAuthMiddleware())
	{

		var qrcodeController kfbackend.QRCodeController
		qrCodeGroup := kf.Group("/qrcode")
		{
			qrCodeGroup.GET("/", qrcodeController.List)
			qrCodeGroup.POST("/switch", qrcodeController.Switch)
			qrCodeGroup.POST("/on-off", qrcodeController.OnOff)
		}

		var chatController kfbackend.ChatController
		chatGroup := kf.Group("/chat")
		{
			chatGroup.POST("/list", chatController.List)
			chatGroup.POST("/msgs", chatController.Msgs)
			chatGroup.POST("/msgs/read", chatController.MsgsRead)
			chatGroup.POST("/extUsers/op", chatController.ExtUserOp)
			chatGroup.POST("/extUser/update", chatController.ExtUserUpdate)
		}

		var sysConfController kfbackend.SysConfController
		settingGroup := kf.Group("/sysConf")
		{
			settingGroup.GET("/", sysConfController.Get)
			settingGroup.POST("/", sysConfController.Post)
		}
	}

	// 客服前台
	kffe := api.Group("/kf-fe")
	{
		kffe.GET("/qrcode/*action")
	}

	// 二维码入口.
	{
		var qrCodeControlelr qrcode.QRCodeController
		{
			g.Any("/s/*action", qrCodeControlelr.Scan)
		}
	}

	// 计费前台.
	billFe := api.Group("/bill-fe")
	{
		var orderController bill2.OrderController
		billFe.POST("/order/create", orderController.CreateOrder)
		billFe.POST("/order/notify", orderController.Notify) // TODO:: 加密处理
	}

	// 内部调用: websocket on auth 回调.
	internal := g.Group("/internal")
	{
		api := internal.Group("/api")
		{
			var nc notify2.NotifyController
			api.POST("websocket-auth", nc.WebsocketAuth)
		}
	}

	// dev push接口.
	dev := api.Group("/dev")
	{
		var dc dev2.DevController
		dev.POST("/push", dc.PushMsg)
	}
}

func swaggerAPI(g *gin.Engine) {
	swag := swaggos.Default()
	swag.Response(200, new(middleware.BaseResponse))

	swag.JWT("Authorization")
	apiGroup := swag.Group("/api")

	// 计费后台的swagger 文档.
	bill.SwaggerDoc(apiGroup)
	kfbackend.SwaggerDoc(apiGroup)
	kffrontend.SwaggerDoc(apiGroup)
	bill2.SwaggerDoc(apiGroup)
	dev2.SwaggerDoc(apiGroup)

	// swagger json 服务
	g.GET("/_doc", gin.WrapH(swag))

	// swagger ui 服务
	g.Any("/doc/*action", gin.WrapH(swaggos.UI("/doc", "/_doc")))

	fmt.Println("swagger ui: " + "http://" + config.GetConfig().Web.String() + "/doc")
}
