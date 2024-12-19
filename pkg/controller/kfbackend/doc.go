package kfbackend

import "github.com/clearcodecn/swaggos"

func SwaggerDoc(group *swaggos.Group) {
	bg := group.Group("/kf-be").Tag("客服后台")
	bg.Post("/login").Body(LoginRequest{}).JSON(LoginResponse{}).Description("登陆接口")

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
		// TODO
		chat.Post("/").Body(LoginRequest{}).JSON(LoginResponse{}).Description("")
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

// TODO 临时放下 后面迁移到各个文件内
type LoginRequest struct {
	CardID   string `json:"cardID" binding:"required"` // 卡密id
	Password string `json:"password"`                  // 密码可选
}

type LoginResponse struct {
	Notice string `json:"notice"` // 公告通知
}

type QRCodeRequest struct{}
type QRCodeResponse struct {
	URL           string         `json:"qrcodeUrl"`     // 主站二维码图片地址
	HealthAt      int64          `json:"healthAt"`      // 主站通过健康检查的时间
	Enable        bool           `json:"enable"`        // 启用停用状态
	EnableNewUser bool           `json:"enableNewUser"` // 启用停用新粉状态
	Domains       []QRCodeDomain `json:"domains"`       // 域名列表
}

type QRCodeDomain struct {
	Domain   string `json:"domain"`   // 站点域名
	HealthAt int64  `json:"healthAt"` // 通过健康检查的时间
	CreateAt int64  `json:"createAt"` // 添加创建时间
	Remark   string `json:"remark"`   // 备注
	URL      string `json:"url"`      // 二维码图片地址
}

type QRCodeSwitchRequest struct{}
type QRCodeSwitchResponse struct {
	URL      string `json:"qrcodeUrl"` // 主站二维码图片地址
	HealthAt int64  `json:"healthAt"`  // 主站通过健康检查的时间
}

type QRCodeOnOffRequest struct {
	OnOff        *bool `json:"onoff"`        // 开关：所有二维码的所有用户都不能进入
	OnOffNewUser *bool `json:"onoffNewUser"` // 开关：老用户可进，新用户不能进
}
type QRCodeOnOffResponse struct{}
