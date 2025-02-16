package dao

import "gorm.io/gorm"

const (
	WelcomeMsg = "welcome_msg"
	QuickReply = "quick_reply"
)

type KfWelcomeMessage struct {
	gorm.Model
	CardId  string `json:"card_id" gorm:"column:card_id;index"` // 卡密id
	Content string `json:"content" gorm:"type:text"`
	Type    string `json:"type" gorm:"type:varchar(255)"`
	Sort    int    `json:"sort"`    // 排序
	Enable  bool   `json:"enable"`  // 是否启用.
	Keyword string `json:"keyword"` // 关键词.
	MsgType string `json:"msgType"` // 类型：welcome_msg=欢迎语, quick_reply=快捷回复.
}

func (KfWelcomeMessage) TableName() string {
	return "kf_welcome_message"
}
