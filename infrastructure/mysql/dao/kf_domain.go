package dao

import "gorm.io/gorm"

type KfDomain struct {
	gorm.Model
}

func (KfDomain) TableName() string {
	return "kf_domain"
}
