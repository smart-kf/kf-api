package dao

import "gorm.io/gorm"

type BillAccount struct {
	gorm.Model

	Username string `gorm:"column:username;unique"`
	Password string `gorm:"column:password"`
}

func (BillAccount) TableName() string {
	return "bill_account"
}
