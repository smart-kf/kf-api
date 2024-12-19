package db

import "gorm.io/gorm"

type SomeTable struct {
	gorm.Model
}

func (SomeTable) TableName() string {
	return "some_table"
}
