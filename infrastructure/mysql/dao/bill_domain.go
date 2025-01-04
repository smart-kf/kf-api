package dao

import "gorm.io/gorm"

type BillDomain struct {
	gorm.Model
	TopName  string `gorm:"column:top_domain"` // 顶级域名
	IsPublic bool   `gorm:"column:is_public"`  // 是否是共享
	IsBind   bool   `gorm:"column:is_bind"`    // 是否有卡密绑定该域名.
	Status   int    `gorm:"column:status"`     // 1:正常，2禁用
}

func (BillDomain) TableName() string {
	return "bill_domain"
}
