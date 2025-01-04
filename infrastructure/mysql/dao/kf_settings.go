package dao

import "gorm.io/gorm"

type KFSettings struct {
	gorm.Model
	Nickname        string `json:"nickname"`
	AvatarURL       string `json:"avatar_url"`
	WSFilter        bool   `json:"wsFilter"`        // 开启之后，检测ws行为.
	WechatFilter    bool   `json:"wechatFilter"`    // 非微信浏览器不能访问 .
	AppleFilter     bool   `json:"appleFilter"`     // 苹果手机过滤器，开启后，只有苹果手机能访问,
	IPProxyFilter   bool   `json:"iPProxyFilter"`   // 代理ip过滤，开启后，代理ip不能访问
	DeviceFilter    bool   `json:"DeviceFilter"`    // 设备异常过滤.
	SimulatorFilter bool   `json:"simulatorFilter"` // 模拟器过滤，开启后，模拟器不能访问
	Notice          string `json:"notice"`          // 滚动公告. maxLength 255
	NewMessageVoice bool   `json:"newMessageVoice"` // 消息提示音.
}

func (KFSettings) TableName() string {
	return "kf_settings"
}
