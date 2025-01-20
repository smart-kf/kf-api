package dao

import "gorm.io/gorm"

type KfFile struct {
	gorm.Model
	Filename   string `json:"filename" gorm:"column:filename"`      // 文件原名字
	Ext        string `json:"ext" gorm:"column:ext"`                // 文件后缀.
	FileType   string `json:"fileType" gorm:"column:file_type"`     // 文件类型.
	Md5        string `json:"md5" gorm:"column:md5"`                // md5字符串.
	CardId     string `json:"cardId" gorm:"column:card_id"`         // 卡密id
	PublicPath string `json:"publicPath" gorm:"column:public_path"` // 访问路径.
}

func (KfFile) TableName() string {
	return "kf_file"
}
