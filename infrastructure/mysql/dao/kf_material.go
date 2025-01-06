package dao

import (
	"gorm.io/gorm"
)

// KFMaterial 素材库（包含话术）
type KFMaterial struct {
	gorm.Model
	Type    MaterialType `gorm:"column:type;type:tinyint(4)" json:"type"`         // 素材类型
	BizType int8         `gorm:"column:biz_type;type:tinyint(8)" json:"biz_type"` // 业务类型
	Content string       `gorm:"column:content;type:longtext;" json:"content"`    // 内容
	Order   int8         `gorm:"column:order;type:tinyint(4)" json:"order"`       // 展示顺序
	Enable  bool         `gorm:"column:enable;type:bool" json:"enable"`           // 是否启用
}

type MaterialType int8

const (
	MaterialTypeText  MaterialType = iota // 文本
	MaterialTypeVoice                     // 语音
	MaterialTypeImage                     // 图片
	MaterialTypeVideo                     // 视频
	MaterialTypeUrl                       // 网址
	MaterialTypeFile                      // 其他文件
)

func (m MaterialType) ToMessageType() MessageType {
	return MessageType(m)
}

func (KFMaterial) TableName() string {
	return "kf_material"
}
