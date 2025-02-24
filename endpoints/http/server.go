package http

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	xlogger "github.com/clearcodecn/log"
	"github.com/clearcodecn/swaggos"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/captcha"

	"github.com/smart-fm/kf-api/data/website"
	"github.com/smart-fm/kf-api/endpoints/http/controller/bill"
	bill2 "github.com/smart-fm/kf-api/endpoints/http/controller/billfrontend"
	dev2 "github.com/smart-fm/kf-api/endpoints/http/controller/dev"
	"github.com/smart-fm/kf-api/endpoints/http/controller/kfbackend"
	"github.com/smart-fm/kf-api/endpoints/http/controller/kffrontend"
	notify2 "github.com/smart-fm/kf-api/endpoints/http/controller/notify"
	"github.com/smart-fm/kf-api/pkg/version"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/endpoints/http/middleware"
	"github.com/smart-fm/kf-api/pkg/utils"
)

func Run() error {
	conf := config.GetConfig()
	g := gin.New()
	g.Use(gin.Recovery())
	// g.Use(
	// 	func(ctx *gin.Context) {
	// 		origin := ctx.Request.Header.Get("Origin")
	// 		allowHeader := ctx.Request.Header.Get("Access-Control-Request-Headers")
	// 		methods := ctx.Request.Header.Get("Access-Control-Request-Method")
	//
	// 		ctx.Header("Access-Control-Allow-Origin", origin)
	// 		ctx.Header("Access-Control-Allow-Methods", methods)
	// 		ctx.Header("Access-Control-Allow-Headers", allowHeader)
	// 		ctx.Header("Access-Control-Allow-Credentials", "true")
	// 		if ctx.Request.Method == http.MethodOptions {
	// 			ctx.AbortWithStatus(204)
	// 			return
	// 		}
	// 		ctx.Next()
	// 	},
	// )
	g.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     []string{"*"},
				AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
				AllowHeaders:     []string{"authorization,x-requested-with,content-type"},
				AllowCredentials: false,
				MaxAge:           24 * time.Hour,
			},
		),
	)

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
	if conf.Debug {
		g.Static("/website/static", "data/website/static")
		g.LoadHTMLGlob("data/website/views/*.html")
	} else {
		g.StaticFS("/website/", http.FS(website.StaticFS))
		tpl, err := template.ParseFS(website.FS, "views/*.html")
		if err != nil {
			panic(err)
		}
		g.SetHTMLTemplate(tpl)
	}

	g.RedirectTrailingSlash = true
	g.RedirectFixedPath = true

	registerRouter(g)
	swaggerAPI(g)

	registerWebsiteRouter(g)

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

		var oc bill.OrderController
		orderGroup := bgAuth.Group("/order")
		{
			orderGroup.POST("/list", oc.List)
		}

		var settingController bill.SettingController
		settingGroup := bgAuth.Group("/setting")
		{
			settingGroup.GET("/get", settingController.Get)
			settingGroup.POST("/update", settingController.Set)
		}
	}

	var authController kfbackend.AuthController
	api.POST("/kf-be/login", authController.Login)

	// 客服后台
	kf := api.Group("/kf-be", middleware.KFBeAuthMiddleware())
	{
		var publicController kfbackend.PublicController
		kf.POST("/upload", publicController.Upload)
		kf.POST("/logout", authController.Logout)

		var qrcodeController kfbackend.QRCodeController
		qrCodeGroup := kf.Group("/qrcode")
		{
			qrCodeGroup.GET("", qrcodeController.List)
			qrCodeGroup.POST("/switch", qrcodeController.Switch)
			qrCodeGroup.POST("/on-off", qrcodeController.OnOff)
		}

		var chatController kfbackend.ChatController
		chatGroup := kf.Group("/chat")
		{
			chatGroup.POST("/list", chatController.List)
			chatGroup.POST("/msgs", chatController.Msgs)
			chatGroup.POST("/batchsend", chatController.BatchSend)
		}

		var sysConfController kfbackend.SysConfController
		settingGroup := kf.Group("/sysConf")
		{
			settingGroup.GET("", sysConfController.Get)
			settingGroup.POST("", sysConfController.Post)
		}

		var kfUserInfo kfbackend.GuestController
		userInfoGroup := kf.Group("/user")
		{
			userInfoGroup.GET("", kfUserInfo.GetKfUserInfo)
			userInfoGroup.POST("/update", kfUserInfo.UpdateUserInfo)
		}

		var wc kfbackend.WelcomeMsgController
		welGroup := kf.Group("/welcome")
		{
			welGroup.POST("/upsert", wc.Upsert)
			welGroup.GET("/list", wc.ListAll)
			welGroup.POST("/del", wc.Delete)
			welGroup.POST("/copy", wc.CopyCardMsg)
		}

		var logc kfbackend.LogController
		logGroup := kf.Group("/log")
		{
			logGroup.GET("/list", logc.List)
		}
	}

	// 客服前台
	kffe := api.Group("/kf-fe")
	kffe.Use(middleware.KFFeAuthMiddleware())
	{
		var qr kffrontend.QRCodeController
		kffe.POST("/qrcode/scan", qr.Scan)
	}

	// 计费前台.
	billFe := api.Group("/bill-fe")
	{
		// /api/bill-fe/order/notify
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

func registerWebsiteRouter(g *gin.Engine) {
	var wc bill2.WebsiteController
	g.GET("/", wc.Index)
	g.GET("/package.html", wc.Package)
	g.GET("/order", wc.Order)
	g.GET("/order/pay-success", wc.PaySuccess)
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
