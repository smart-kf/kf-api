package kfbackend

import (
	"github.com/clearcodecn/swaggos"

	. "github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
)

func SwaggerDoc(group *swaggos.Group) {
	bg := group.Group("/kf-be").Tag("客服后台")

	public := group.Group("/public").Tag("公开接口")
	{
		public.Get("/captchaId").JSON(GetQRCodeIDResponse{}).Description("获取验证码id")
		public.Get("/captcha/:captchaId.png").Description("获取验证码图片")
		bg.Post("/login").Body(LoginRequest{}).JSON(LoginResponse{}).Description("登陆接口").Tag("公开接口")
	}

	qrcode := bg.Group("/qrcode").Tag("二维码管理").
		Header("authorization", "授权session", true)
	{
		qrcode.Post("/").
			Body(QRCodeRequest{}).JSON(QRCodeResponse{}).
			Description("获取二维码和域名列表的接口")

		qrcode.Post("/switch").
			Body(QRCodeSwitchRequest{}).JSON(QRCodeSwitchResponse{}).
			Description("更换二维码")

		qrcode.Post("/on-off").
			Body(QRCodeOnOffRequest{}).JSON(QRCodeOnOffResponse{}).
			Description("二维码功能开关")
	}

	chat := bg.Group("/chat").Tag("聊天管理").
		Header("authorization", "授权session", true)
	{
		chat.Post("/list").Body(ChatListRequest{}).JSON(ChatListResponse{}).Description("会话列表")
	}

	msgLib := bg.Group("/msgLib").Tag("话术管理").
		Header("authorization", "授权session", true)
	{
		// TODO
		msgLib.Post("/").Body(LoginRequest{}).JSON(LoginResponse{})
	}

	sysLog := bg.Group("/sysLog").Tag("操作日志").
		Header("authorization", "授权session", true)
	{
		// TODO
		sysLog.Post("/").Body(LoginRequest{}).JSON(LoginResponse{})
	}

	sysConf := bg.Group("/sysConf").Tag("系统配置").
		Header("authorization", "授权session", true)
	{
		// TODO
		sysConf.Post("/").Body(LoginRequest{}).JSON(LoginResponse{})
	}
}
