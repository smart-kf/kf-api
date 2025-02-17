package kfbackend

import (
	"github.com/clearcodecn/swaggos"

	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
)

func SwaggerDoc(group *swaggos.Group) {
	bg := group.Group("/kf-be").Tag("客服后台")

	public := group.Group("/public").Tag("公开接口")
	{
		public.Get("/captchaId").JSON(kfbackend.GetQRCodeIDResponse{}).Description("获取验证码id")
		public.Get("/captcha/:captchaId.png").Description("获取验证码图片")
		bg.Post("/login").Body(kfbackend.LoginRequest{}).JSON(kfbackend.LoginResponse{}).Description("登陆接口").Tag("公开接口")
	}

	bg.Post("/upload").Tag("公共接口").
		FormFile(
			"file",
			swaggos.Attribute{Required: true, Description: "文件, 不超过32M"},
		).
		Form("fileType", swaggos.Attribute{Required: true, Description: "文件类型: image || video"}).
		JSON(kfbackend.UploadResponse{}).Description("文件上传")

	bg.Post("/logout").JSON(common.EmptyResponse{}).Description("退出")

	qrcode := bg.Group("/qrcode").Tag("二维码管理")
	{
		qrcode.Get("/").
			QueryObject(kfbackend.QRCodeRequest{}).JSON(kfbackend.QRCodeResponse{}).
			Description("获取二维码和域名列表的接口")

		qrcode.Post("/switch").
			Body(kfbackend.QRCodeSwitchRequest{}).JSON(kfbackend.QRCodeSwitchResponse{}).
			Description("更换二维码")

		qrcode.Post("/on-off").
			Body(kfbackend.QRCodeOnOffRequest{}).JSON(kfbackend.QRCodeOnOffResponse{}).
			Description("二维码功能开关")
	}

	chat := bg.Group("/chat").Tag("聊天管理")
	{
		chat.Post("/list").Body(kfbackend.ChatListRequest{}).JSON(kfbackend.ChatListResponse{}).Description("会话列表")
		chat.Post("/msgs").Body(kfbackend.MsgListRequest{}).JSON(kfbackend.MsgListResponse{}).Description("消息列表 按消息id倒序滚页查询")
		chat.Post("/batchsend").Body(kfbackend.BatchSendRequest{}).JSON(common.EmptyResponse{}).Description(
			"群发消息",
		)
	}

	user := bg.Group("/user").Tag("客服后台客户信息")
	{
		user.Get("").FormObject(kfbackend.GetKfUserInfoRequest{}).JSON(kfbackend.User{}).Description("获取客户信息")
		user.Post("/update").Body(kfbackend.UpdateUserInfoRequest{}).JSON(common.EmptyResponse{}).Description(
			"更新用户信息",
		)
	}

	// msgLib := bg.Group("/msgLib").Tag("话术管理")
	// {
	// 	// TODO
	// 	msgLib.Post("/").Body(kfbackend.LoginRequest{}).JSON(kfbackend.LoginResponse{})
	// }
	//
	// sysLog := bg.Group("/sysLog").Tag("操作日志")
	// {
	// 	// TODO
	// 	sysLog.Post("/").Body(kfbackend.LoginRequest{}).JSON(kfbackend.LoginResponse{})
	// }

	sysConf := bg.Group("/sysConf").Tag("系统配置")
	{
		sysConf.Get("/").JSON(GetSysConfResponse{})
		sysConf.Post("/").Body(PostSysConfRequest{}).JSON(PostSysConfResponse{})
	}

	// 欢迎语.
	wel := bg.Group("/welcome").Tag("客服后台-欢迎语、智能回复、快捷回复")
	{
		wel.Post("/upsert").Body(kfbackend.UpsertWelcomeMsgRequest{}).JSON(common.EmptyResponse{}).Description(
			"创建||更新欢迎语",
		)
		wel.Get("/list").FormObject(kfbackend.ListAllRequest{}).JSON([]*kfbackend.KfWelcomeMessageResp{}).Description("欢迎语列表")
		wel.Post("/del").Body(kfbackend.DeleteWelcomeRequest{}).JSON(common.EmptyResponse{}).Description("删除欢迎语")
		wel.Post("/copy").Body(kfbackend.CopyCardMsgRequest{}).JSON(common.EmptyResponse{}).Description("复制话术")
	}

	// 日志
	log := bg.Group("/log")
	{
		log.Get("/list").QueryObject(kfbackend.LogRequest{}).JSON(kfbackend.ListLogResponse{}).Description("日志列表")
	}
}
