package dao

import (
	"gorm.io/gorm"
)

// BillLog 计费后台审计日志.
type BillLog struct {
	gorm.Model
	Operator   string `json:"operator" gorm:"column:operator"`       // 操作人
	HandleFunc string `json:"handle_func" gorm:"column:handle_func"` // 操作类型
	Content    string `json:"content" gorm:"column:content;text"`    // 操作内容
}

func (BillLog) TableName() string {
	return "bill_log"
}
