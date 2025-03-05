package dao

import (
	"gorm.io/gorm"

	"github.com/smart-fm/kf-api/endpoints/common"
)

// KFMessage 消息
type KFMessage struct {
	gorm.Model
	MsgId   string             `gorm:"column:msg_id"`
	MsgType common.MessageType `gorm:"column:msg_type;type:varchar(128)" json:"type"` // 消息类型
	GuestId string             `gorm:"column:guest_id" json:"guestId"`                // 客服id
	CardId  string             `gorm:"column:card_id" json:"cardId"`                  // 卡密id
	Content string             `gorm:"column:content;type:longtext;" json:"content"`  // 内容.
	IsKf    int                `gorm:"column:is_kf;type:tinyint(4)"`                  // 是否是客服
}

func (KFMessage) TableName() string {
	return "kf_message"
}
