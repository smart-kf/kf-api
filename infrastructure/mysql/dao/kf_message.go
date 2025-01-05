package dao

import "gorm.io/gorm"

type KFMessage struct {
	gorm.Model
	MsgType string `gorm:"column:msg_type;type:varchar(255)" json:"msgType;"` // 消息类型：text=文本，video=视频,image=图片
	KfId    int64  `gorm:"column:kf_id;" json:"kfId"`                         // 客服id
	GuestId int64  `gorm:"column:guest_id;" json:"guestId"`                   // 客户id
	Content string `gorm:"column:content;type:longtext;" json:"content"`      // 内容.
	City    string `gorm:"column:city;type:varchar(255)" json:"city"`         // 城市.
	Ip      string `gorm:"column:ip;type:varchar(255)" json:"ip"`             // ip 地址.
	IsKF    int    `gorm:"column:is_kf;" json:"isKf"`                         // 是否是客服：1=是客服消息，2=是客户消息.
}

func (KFMessage) TableName() string {
	return "kf_message"
}
