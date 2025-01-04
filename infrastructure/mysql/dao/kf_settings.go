package dao

import "gorm.io/gorm"

type KFSettings struct {
	gorm.Model
	CardID               string `json:"cardId" gorm:"column:card_id;unique" doc:"卡密id"`
	Nickname             string `json:"nickname" gorm:"column:nickname"`
	AvatarURL            string `json:"avatarUrl" gorm:"column:avatar_url"`
	WSFilter             bool   `json:"wsFilter" gorm:"column:ws_filter"`                           // 开启之后，检测ws行为.
	WechatFilter         bool   `json:"wechatFilter" gorm:"column:wechat_filter"`                   // 非微信浏览器不能访问 .
	AppleFilter          bool   `json:"appleFilter" gorm:"column:apple_filter"`                     // 苹果手机过滤器，开启后，只有苹果手机能访问,
	IPProxyFilter        bool   `json:"iPProxyFilter" gorm:"column:ip_proxy_filter"`                // 代理ip过滤，开启后，代理ip不能访问
	DeviceFilter         bool   `json:"DeviceFilter" gorm:"column:device_filter"`                   // 设备异常过滤.
	SimulatorFilter      bool   `json:"simulatorFilter" gorm:"column:simulator_filter"`             // 模拟器过滤，开启后，模拟器不能访问
	Notice               string `json:"notice" gorm:"column:notice"`                                // 滚动公告. maxLength 255
	NewMessageVoice      bool   `json:"newMessageVoice" gorm:"column:new_message_voice"`            // 消息提示音.
	QRCodeEnabled        bool   `json:"qrCodeEnabled" gorm:"column:qrcode_enabled"`                 // 二维码是否启用
	QRCodeVersion        int    `json:"qrCodeVersion" gorm:"column:qrcode_version"`                 // 二维码版本 每次切换后更新
	QRCodeEnabledNewUser bool   `json:"qrCodeEnabledNewUser" gorm:"column:qrcode_enabled_new_user"` // 启用停用新粉状态
}

func (KFSettings) TableName() string {
	return "kf_settings"
}
