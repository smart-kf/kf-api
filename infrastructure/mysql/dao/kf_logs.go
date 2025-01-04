package dao

import (
	"gorm.io/gorm"
)

// KFLog 客服后台审计日志.
type KFLog struct {
	gorm.Model
	CardID     string `json:"card_id" gorm:"column:card_id"`
	HandleFunc string `json:"handle_func" gorm:"column:handle_func"` // 操作类型
	Content    string `json:"content" gorm:"column:content;text"`    // 操作内容
}

func (KFLog) TableName() string {
	return "kf_log"
}
