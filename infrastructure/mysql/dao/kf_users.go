package dao

import "gorm.io/gorm"

type KfUsers struct {
	gorm.Model
	CardID        string `json:"card_id" gorm:"column:card_id;index"`               // 卡密id
	UUID          string `json:"uuid" gorm:"column:uuid;unique;type:varchar(32)"`   // 用户的uuid，不用主键做业务.
	Avatar        string `json:"avatar" gorm:"column:avatar;type:varchar(255)"`     // 头像地址，存储的是相对路径
	Nickname      string `json:"nickname" gorm:"column:nickname;type:varchar(255)"` // 昵称.
	RemarkName    string `json:"remarkName" gorm:"column:remark_name"`              // 备注名称.
	Mobile        string `json:"mobile" gorm:"column:mobile"`                       // 手机号
	Comments      string `json:"comments" gorm:"column:comments"`                   // 备注信息
	IP            string `json:"ip" gorm:"column:ip;type:varchar(255)"`             // 注册ip
	Area          string `json:"area" gorm:"column:area;type:varchar(255)"`         // ip对应的地区
	OfflineTime   int64  `json:"offline_time" gorm:"column:offline_time"`           // 离线时间
	Device        string `json:"device" gorm:"column:device"`                       // 设备类型： iphone、android、
	Browser       string `json:"browser" gorm:"column:browser"`                     // 浏览器类型
	ScanQRCodeCnt int    `json:"scanQRCodeCnt" gorm:"column:scan_qrcode_cnt"`       // 扫码次数.
	IsTop         int    `json:"isTop" gorm:"column:is_top"`                        // 聊天置顶：1=置顶，2=不置顶
	IsBlack       int    `json:"isBlack" gorm:"column:is_black"`                    // 是否被拉黑: 1=拉黑, 2=不拉黑
	IsSimulator   int    `json:"isSimulator" gorm:"column:is_simulator"`            // 是否是模拟器: 1=是模拟器，2=不是模拟器
	IsProxy       int    `json:"isProxy" gorm:"column:is_proxy"`                    // 是否使用了代理ip访问: 1=是，2=不是.
	Source        string `json:"source" gorm:"column:source"`                       // 来源
}

// 还有2个字段，存储在redis:
// 1. onlineStatus = 是否在线
// 2. unreadMsgCnt = 未读消息条数

func (KfUsers) TableName() string {
	return "kf_users"
}
