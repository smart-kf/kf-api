package kffrontend

type QRCodeScanRequest struct {
	Code string `json:"code" doc:"二维码的code" validate:"required" binding:"required"` // 二维码path
}

type KFUserInfo struct {
	UUID       string `json:"uuid" gorm:"column:uuid;unique;type:varchar(32)"` // 用户的uuid，不用主键做业务.
	Avatar     string `json:"avatar" gorm:"column:avatar;type:varchar(255)"`   // 头像地址，存储的是相对路径
	NickName   string `json:"nickName" gorm:"column:nick_name;type:varchar(255)" doc:"昵称"`
	WsHost     string `json:"wsHost" doc:"ws地址"`
	WsFullHost string `json:"wsFullHost" doc:"ws带域名协议地址"`
}

type QRCodeScanResponse struct {
	UserInfo  KFUserInfo `json:"userInfo" doc:"用户信息"`
	IsNewUser bool       `json:"isNewUser" doc:"是否是新用户"`
	CdnHost   string     `json:"cdnHost" doc:"cdn域名"`
}

type CheckResponse struct{}

type MsgListRequest struct {
	LastMsgTime int64 `json:"lastMsgTime"`
}

type SmartMsg struct {
	Id      int64  `json:"id"`
	Keyword string `json:"keyword"`
}
