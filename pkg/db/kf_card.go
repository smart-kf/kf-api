package db

import (
	"gorm.io/gorm"
	"std-api/pkg/constant"
	"time"
)

type KFCard struct {
	gorm.Model
	CardID        string               `json:"cardId" gorm:"column:card_id;unique" doc:"卡密id"`
	Password      string               `json:"password" gorm:"column:password" doc:"密码"`
	SaleStatus    constant.SaleStatus  `json:"saleStatus" gorm:"column:sale_status" doc:"销售状态"`
	LoginStatus   constant.LoginStatus `json:"loginStatus" gorm:"column:login_status" doc:"登录状态"`
	CardType      constant.CardType    `json:"cardType" gorm:"column:card_type" doc:"卡片类型"`
	Day           int                  `json:"day" gorm:"column:day" doc:"卡密的天数"`
	ExpireTime    int64                `json:"expireTime" gorm:"column:expire_time" doc:"过期时间"`
	LastLoginTime int64                `json:"lastLoginTime" gorm:"column:last_login_time" doc:"上次登录时间"`
	Version       int                  `json:"version" gorm:"version" doc:"乐观锁"`
}

func (c KFCard) HasExpire() bool {
	return time.Now().Unix() > c.ExpireTime
}
func (KFCard) TableName() string {
	return "kf_card"
}
