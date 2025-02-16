package dao

import "gorm.io/gorm"

type KfWelcomeMessage struct {
	gorm.Model
	CardId  string `json:"card_id" gorm:"column:card_id;index"` // 卡密id
	Content string `json:"content" gorm:"type:text"`
	Type    string `json:"type" gorm:"type:varchar(255)"`
	Sort    int    `json:"sort"`   // 排序
	Enable  bool   `json:"enable"` // 是否启用.
}

func (KfWelcomeMessage) TableName() string {
	return "kf_welcome_message"
}
