package dao

import "gorm.io/gorm"

type KFSettings struct {
	gorm.Model
	CardID           string `json:"cardId" gorm:"column:card_id;unique;type:varchar(255)" doc:"卡密id"`
	Nickname         string `json:"nickname" gorm:"column:nickname;type:varchar(255)"`
	AvatarURL        string `json:"avatarUrl" gorm:"column:avatar_url;type:varchar(255)"`
	WSFilter         bool   `json:"wsFilter" gorm:"column:ws_filter"`                               // 开启之后，检测ws行为.
	WechatFilter     bool   `json:"wechatFilter" gorm:"column:wechat_filter"`                       // 非微信浏览器不能访问 .
	AppleFilter      bool   `json:"appleFilter" gorm:"column:apple_filter"`                         // 苹果手机过滤器，开启后，只有苹果手机能访问,
	IPProxyFilter    bool   `json:"iPProxyFilter" gorm:"column:ip_proxy_filter"`                    // 代理ip过滤，开启后，代理ip不能访问
	DeviceFilter     bool   `json:"DeviceFilter" gorm:"column:device_filter"`                       // 设备异常过滤.
	SimulatorFilter  bool   `json:"simulatorFilter" gorm:"column:simulator_filter"`                 // 模拟器过滤，开启后，模拟器不能访问
	Notice           string `json:"notice" gorm:"column:notice;type:longtext"`                      // 滚动公告. maxLength 255
	NewMessageVoice  bool   `json:"newMessageVoice" gorm:"column:new_message_voice"`                // 消息提示音.
	QRCodeEnabled    bool   `json:"qrCodeEnabled" gorm:"column:qrcode_enabled"`                     // 二维码是否启用
	QRCodeScanFilter int    `json:"qrCodeScanFilter" gorm:"column:qrcode_scan_filter;default(1);" ` // doc:"扫码过滤: 1=关闭，2 =过滤机房，3=过滤非大陆，4 =过滤机房及非大陆"
}

func (KFSettings) TableName() string {
	return "kf_settings"
}

func NewDefaultKFSettings(cardID string) *KFSettings {
	return &KFSettings{
		CardID:          cardID,
		Nickname:        "客服",
		AvatarURL:       "/static/avatar/kf.png",
		WSFilter:        false,
		WechatFilter:    false,
		AppleFilter:     false,
		IPProxyFilter:   false,
		DeviceFilter:    false,
		SimulatorFilter: false,
		NewMessageVoice: true,
		QRCodeEnabled:   true,
	}
}
