package dao

import (
	"gorm.io/gorm"
)

// KFMessage 消息
type KFMessage struct {
	gorm.Model
	ID       uint64      `gorm:"primaryKey;type;bigint unsigned;autoIncrement"`    // 消息自增id
	Type     MessageType `gorm:"column:type;type:tinyint(4)" json:"type"`          // 消息类型
	From     string      `gorm:"column:from" json:"from"`                          // 发送方id
	FromType ChatObjType `gorm:"column:from_type;type:tinyint(4)" json:"fromType"` // 发送方类型
	To       string      `gorm:"column:to" json:"to"`                              // 接收方id
	ToType   ChatObjType `gorm:"column:to_type;type:tinyint(4)" json:"toType"`     // 接收方类型
	ReadAt   int64       `gorm:"column:read_at" json:"read_at"`                    // 接收方已读消息的时间
	Content  string      `gorm:"column:content;type:longtext;" json:"content"`     // 内容.
}

type MessageType int8

const (
	MessageTypeText  MessageType = iota // 文本
	MessageTypeVoice                    // 语音
	MessageTypeImage                    // 图片
	MessageTypeVideo                    // 视频
	MessageTypeUrl                      // 网址
	MessageTypeFile                     // 其他文件
)

func (m MessageType) ToMaterialType() MaterialType {
	return MaterialType(m)
}

// ChatObjType 聊天对象的类型
type ChatObjType int8

const (
	ChatObjTypeSys          ChatObjType = iota // 系统
	ChatObjTypeExternalUser                    // 访客 即用户/粉丝
	ChatObjTypeUser                            // 员工 即客服
	ChatObjTypeGroup                           // 群组
)

func (KFMessage) TableName() string {
	return "kf_message"
}
