package dao

import "gorm.io/gorm"

type KfUser struct {
	gorm.Model
	CardID      string `json:"card_id" gorm:"column:card_id;index"`              // 卡密id
	UUID        string `json:"uuid" gorm:"column:uuid;unique;type:varchar(255)"` // 用户的uuid，不用主键做业务.
	Avatar      string `json:"avatar" gorm:"column:avatar;type:varchar(255)"`    // 头像地址，存储的是相对路径
	NickName    string `json:"nickName" gorm:"column:nick_name;type:varchar(255)" doc:"昵称"`
	RemarkName  string `json:"remarkName" gorm:"column:remark_name"`      // 备注名称.
	Mobile      string `json:"mobile" gorm:"column:mobile"`               // 手机号
	Comments    string `json:"comments" gorm:"column:comments"`           // 备注信息
	IP          string `json:"ip" gorm:"column:ip;type:varchar(255)"`     // 注册ip
	Area        string `json:"area" gorm:"column:area;type:varchar(255)"` // ip对应的地区
	UserAgent   string `json:"userAgent" gorm:"column:user_agent;type:varchar(1000)" doc:"浏览器user-agent"`
	Browser     string `json:"browser" gorm:"column:browser;type:varchar(255)"`   // 浏览器 Chrome/Safari/firfox/...
	Device      string `json:"device" gorm:"column:device;type:varchar(50)"`      // 设备类型： iphone、android、
	IsProxy     int    `json:"isProxy" gorm:"column:is_proxy"`                    // 是否使用了代理ip访问: 1=是，2=不是.
	IsEmulator  int    `json:"isEmulator" gorm:"column:is_emulator"`              // 是否是模拟器 1=是，2=不是
	Source      string `json:"source" gorm:"column:source"`                       // 来源
	OfflineAt   int64  `json:"offlineAt" gorm:"column:offline_at" doc:"离线时间 秒"`   // ws断开链接时记录
	NetworkType string `json:"networkType" gorm:"column:network_type" doc:"网络类型"` // wifi/4G/5G
	ScanCount   int64  `json:"scanCount" gorm:"column:scan_count" doc:"扫码次数"`
	TopAt       int64  `json:"topAt" gorm:"column:top_at" doc:"置顶时间 >0则是置顶 秒"`
	BlockAt     int64  `json:"blockAt" gorm:"column:block_at" doc:"拉黑时间 >0则是拉黑 秒"`
	LastChatAt  int64  `json:"lastChatAt" gorm:"column:last_chat_at" doc:"最近聊天时间 毫秒"`
	LastMsgID   uint64 `json:"lastMsgID" gorm:"column:last_msg_id;type:bigint unsigned" doc:"最近一次由该用户发送的消息id"`
}

func (KfUser) TableName() string {
	return "kf_users"
}

// UserExtra 用户持久化存储属性
type UserExtra struct {
	LastChatTime    int64  `json:"last_chat_time"`    // 最近聊天时间
	LastOfflineTime int64  `json:"last_offline_time"` // 最近离线时间.
	LastMessageId   string `json:"last_message_id"`   // 最近消息id.
}
