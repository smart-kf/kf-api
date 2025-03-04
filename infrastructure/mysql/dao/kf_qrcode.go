package dao

import "gorm.io/gorm"

type KFQRCodeDomain struct {
	gorm.Model
	Version          int    `json:"version" gorm:"column:version"`                               // 版本号.
	ChangeQRCodeTime int64  `json:"changeQRCodeTime" gorm:"column:change_qrcode_time"`           // 换码时间
	Status           int    `json:"status" gorm:"column:status" doc:"二维码状态: 1=正常, 2=失效,3=暂停引新粉"` // 二维码状态
	Path             string `json:"path" gorm:"column:path;type:varchar(255)"`                   // 扫描路径
	CardID           string `json:"cardId" gorm:"card_id"`
	Domain           string `json:"domain" gorm:"column:domain;type:varchar(255)"` // 前台域名
	DomainId         int64  `json:"domainId" gorm:"column:domain_id"`              // bill_domain 的id
	IsPrivate        bool   `json:"is_private" gorm:"column:is_private"`           // 是否是私有
}

func (KFQRCodeDomain) TableName() string {
	return "kf_qrcode_domain"
}
