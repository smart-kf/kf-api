package dao

import "gorm.io/gorm"

type BillAccount struct {
	gorm.Model

	Username string `gorm:"column:username;unique;type:varchar(255)"`
	Password string `gorm:"column:password;type:varchar(255)"`
}

func (BillAccount) TableName() string {
	return "bill_account"
}
