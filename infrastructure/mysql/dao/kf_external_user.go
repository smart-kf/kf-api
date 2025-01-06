package dao

import (
	"gorm.io/gorm"
)

// KFExternalUser 访客
type KFExternalUser struct {
	gorm.Model
	CardID        string `json:"cardID" gorm:"column:card_id" doc:"卡密id"`
	Avatar        string `json:"avatar" gorm:"column:avatar" doc:"头像"`
	NickName      string `json:"nickName" gorm:"column:nick_name" doc:"昵称"`
	PhoneNumber   string `json:"phoneNumber" gorm:"column:phone_number" doc:"手机号"`
	Remark        string `json:"remark" gorm:"column:remark" doc:"备注"`
	UserAgent     string `json:"userAgent" gorm:"column:user_agent" doc:"浏览器属性"`
	City          string `json:"city" gorm:"column:city" doc:"城市"`
	IP            string `json:"ip" gorm:"column:ip" doc:"ip"`
	CreateAt      int64  `json:"createAt" gorm:"column:create_at" doc:"注册时间 秒"`
	OfflineAt     int64  `json:"offlineAt" gorm:"column:offline_at" doc:"离线时间 秒"`
	DeviceVersion string `json:"deviceVersion" gorm:"column:device_version" doc:"设备系统版本"`
	NetworkType   string `json:"networkType" gorm:"column:network_type" doc:"网络类型"`
	ScanCount     int64  `json:"scanCount" gorm:"column:scan_count" doc:"扫码次数"`
	TopAt         int64  `json:"topAt" gorm:"column:top_at" doc:"置顶时间 秒"`
	BlockAt       int64  `json:"blockAt" gorm:"column:block_at" doc:"拉黑时间 秒"`
	LastChatAt    int64  `json:"lastChatAt" gorm:"column:last_chat_at" doc:"最近聊天时间 毫秒"`
	LastMsgID     uint64 `json:"lastMsgID" gorm:"column:last_msg_id;type:bigint unsigned" doc:"最近一次由该用户发送的消息id"`
	UnreadMsgCnt  int64  `json:"unreadMsgCnt" gorm:"column:unread_msg_cnt" doc:"未读消息数"`
}

func (KFExternalUser) TableName() string {
	return "kf_external_user"
}
