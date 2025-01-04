package dao

import (
	"gorm.io/gorm"
)

// BillLog 计费后台审计日志.
type BillLog struct {
	gorm.Model
	Operator   string `json:"operator" gorm:"column:operator;type:varchar(255)"`       // 操作人
	HandleFunc string `json:"handle_func" gorm:"column:handle_func;type:varchar(255)"` // 操作类型
	Content    string `json:"content" gorm:"column:content;longtext"`                  // 操作内容
}

func (BillLog) TableName() string {
	return "bill_log"
}
