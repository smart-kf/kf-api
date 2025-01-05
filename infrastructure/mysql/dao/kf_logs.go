package dao

import (
	"gorm.io/gorm"
)

// KFLog 客服后台审计日志.
type KFLog struct {
	gorm.Model
	CardID     string `json:"card_id" gorm:"column:card_id;type:varchar(255)"`
	HandleFunc string `json:"handle_func" gorm:"column:handle_func;type:varchar(255)"` // 操作类型
	Content    string `json:"content" gorm:"column:content;longtext;"`                 // 操作内容
}

func (KFLog) TableName() string {
	return "kf_log"
}
