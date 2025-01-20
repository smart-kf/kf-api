package kfbackend

import (
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type QRCodeRequest struct{}
type QRCodeResponse struct {
	// URL           string         `json:"qrcodeUrl,omitempty" doc:"主站二维码图片地址"`
	QRCodeURL     string         `json:"qrCodeUrl" doc:"二维码的内容，前端拿到字符串渲染二维码"`
	HealthAt      int64          `json:"healthAt,omitempty" doc:"主站通过健康检查的时间 毫秒"`
	Enable        bool           `json:"enable,omitempty" doc:"启用停用状态"`
	EnableNewUser bool           `json:"enableNewUser,omitempty" doc:"启用停用新粉状态"`
	Version       int            `json:"version" doc:"版本号"`
	Domains       []QRCodeDomain `json:"domains,omitempty" doc:"域名列表"`
}

type QRCodeDomain struct {
	Domain    string `json:"domain,omitempty" doc:"站点域名"`
	HealthAt  int64  `json:"healthAt,omitempty" doc:"通过健康检查的时间 毫秒"`
	CreateAt  int64  `json:"createAt,omitempty" doc:"添加创建时间 毫秒"`
	IsPrivate bool   `json:"isPrivate" doc:"是否是私有域名"`
}

type QRCodeSwitchRequest struct{}
type QRCodeSwitchResponse struct {
	QRCodeURL string `json:"qrCodeUrl" doc:"二维码的内容，前端拿到字符串渲染二维码"`
}

type QRCodeOnOffRequest struct {
	OnOff        *bool `json:"onoff" doc:"开关：所有二维码的所有用户都不能进入"`
	OnOffNewUser *bool `json:"onoffNewUser" doc:"开关：老用户可进，新用户不能进"`
}
type QRCodeOnOffResponse struct{}

type ChatListRequest struct {
	SearchBy             string       `json:"searchBy" doc:"模糊搜索 用户id/昵称/手机号/备注"`
	ListType             ChatListType `json:"listType" doc:"列表类型 0:全部(默认) 1:消息未读 2:拉黑访客"`
	common.ScrollRequest `json:",inline"`
}

type ChatListType int8

const (
	ChatListTypeDefault ChatListType = 0
	ChatListTypeUnread  ChatListType = 1
	ChatListTypeBlock   ChatListType = 2
)

type ChatListResponse struct {
	Chats []*Chat `json:"chats,omitempty" doc:"会话列表"`
}

type ChatType int8

const (
	ChatTypeSingle ChatType = iota // 单聊
	ChatTypeGroup                  // 群聊
)

type Chat struct {
	User         User     `json:"user" doc:"访客信息"`
	LastChatAt   int64    `json:"lastChatAt" doc:"最近聊天时间 毫秒"`
	LastMessage  *Message `json:"lastMessage" doc:"最近一次聊天的消息内容"`
	UnreadMsgCnt int64    `json:"unreadMsgCnt" doc:"未读消息数"`
}

type MsgListRequest struct {
	GuestId              string `json:"guestId" binding:"required" doc:" 即一个会话中粉丝的id"`
	common.ScrollRequest `json:",inline"`
}

type MsgListResponse struct {
	Messages []*Message `json:"messages" doc:"消息列表"`
}

type ReadMsgRequest struct {
	MsgIDs []uint64 `json:"msgIDs" doc:"已读的消息ids"`
}

type ReadMsgResponse struct {
}

type BatchOpUserRequest struct {
	UserIDs []string `json:"userIDs" doc:"粉丝ids"`
	Op      UserOp   `json:"op" doc:"操作 1:置顶 2:取消置顶 3:拉黑 4:取消拉黑"`
}

type UserOp int8

const (
	_ UserOp = iota
	UserOpTop
	UserOpTopUndo
	UserOpBlock
	UserOpBlockUndo
)

type BatchOpUserResponse struct {
}

type UpdateUserRequest struct {
	ID         string `json:"id" doc:"粉丝id"`
	RemarkName string `json:"remarkName" doc:"备注名称"`
	Mobile     string `json:"mobile" doc:"手机号"`
	Comments   string `json:"comments" doc:"备注信息"`
}

type UpdateUserResponse struct {
}

type Message struct {
	MsgId    string             `doc:"消息id" json:"msgId"`
	MsgType  common.MessageType `json:"type" doc:"消息类型:text||image||video"`              // 消息类型
	GuestId  string             `gorm:"column:guest_id" json:"guestId" doc:"客户id"`       // 客服id
	CardId   string             `gorm:"column:card_id" json:"cardId" doc:"卡密id: 只有后台有值"` // 卡密id
	Content  string             `json:"content" doc:"消息的内容"`                             // 内容.
	IsKf     int                `gorm:"column:is_kf;type:tinyint(4)"`                    // 是否是客服
	CreateAt int64              `json:"createAt" doc:"消息创建时间 单位秒 可用作合并时间窗口"`
}

type Material struct {
	Type  dao.MaterialType `json:"type" doc:"资源类型 0:文本 1:语音 2:图片 3:视频 4:网址 5:其他文件"`
	Text  Text             `json:"text,omitempty" doc:"文本"`
	Voice Voice            `json:"voice,omitempty" doc:"语音" `
	Image Image            `json:"image,omitempty" doc:"图片"`
	Video Video            `json:"video,omitempty" doc:"视频"`
	URL   URL              `json:"url,omitempty" doc:"网址"`
	File  File             `json:"file,omitempty" doc:"文件"`
}

type Text struct {
	Content string `json:"content,omitempty" doc:"文本内容"`
}

type Voice struct {
	URL string `json:"url,omitempty" doc:"语音媒体文件的url地址"`
}

type Image struct {
	URL string `json:"url,omitempty" doc:"图片文件的url地址"`
}

type Video struct {
	URL         string `json:"url,omitempty" doc:"视频文件的url地址"`
	CoverURL    string `json:"coverUrl,omitempty" doc:"视频首帧图url地址"`
	DurationSec int64  `json:"duration,omitempty" doc:"视频时长 单位秒"`
}

type URL struct {
	URL      string `json:"url,omitempty" doc:"网址地址"`
	CoverURL string `json:"coverUrl,omitempty" doc:"网址封面图地址"`
	Desc     string `json:"desc,omitempty" doc:"网站描述"`
}

type File struct {
	URL  string `json:"url,omitempty" doc:"文件的url地址"`
	Type string `json:"type,omitempty" doc:"文件的类型格式 如zip/txt/..."`
	Size int64  `json:"size,omitempty" doc:"文件的大小 单位字节"`
}

type User struct {
	CardID      string `json:"card_id" gorm:"column:card_id;index" doc:"卡密id"`                    // 卡密id
	UUID        string `json:"uuid" gorm:"column:uuid;unique;type:varchar(255)" doc:"用户的uuid"`    // 用户的uuid，不用主键做业务.
	Avatar      string `json:"avatar" gorm:"column:avatar;type:varchar(255)" doc:"头像地址，存储的是相对路径"` // 头像地址，存储的是相对路径
	NickName    string `json:"nickName" gorm:"column:nick_name;type:varchar(255)" doc:"昵称"`
	RemarkName  string `json:"remarkName" gorm:"column:remark_name" doc:"备注名称"`         // 备注名称.
	Mobile      string `json:"mobile" gorm:"column:mobile" doc:"手机号"`                   // 手机号
	Comments    string `json:"comments" gorm:"column:comments" doc:"备注信息"`              // 备注信息
	IP          string `json:"ip" gorm:"column:ip;type:varchar(255)" doc:"注册ip"`        // 注册ip
	Area        string `json:"area" gorm:"column:area;type:varchar(255)" doc:"ip对应的地区"` // ip对应的地区
	UserAgent   string `json:"userAgent" gorm:"column:user_agent;type:varchar(1000)" doc:"浏览器user-agent"`
	Browser     string `json:"browser"  doc:"浏览器 Chrome/Safari/firfox/..."`                  // 浏览器 Chrome/Safari/firfox/...
	Device      string `json:"device" doc:"设备类型： iphone、android、"`                           // 设备类型： iphone、android、
	IsProxy     int    `json:"isProxy" gorm:"column:is_proxy" doc:"是否使用了代理ip访问: 1=是，2=不是."`  // 是否使用了代理ip访问: 1=是，2=不是.
	IsEmulator  int    `json:"isEmulator" gorm:"column:is_emulator" doc:"是否是模拟器 1=是，2=不是"`   // 是否是模拟器 1=是，2=不是
	Source      string `json:"source" gorm:"column:source" doc:"来源"`                         // 来源
	OfflineAt   int64  `json:"offlineAt" gorm:"column:offline_at" doc:"离线时间 秒"`              // ws断开链接时记录
	NetworkType string `json:"networkType" gorm:"column:network_type" doc:"网络类型:wifi/4G/5G"` // wifi/4G/5G
	ScanCount   int64  `json:"scanCount" gorm:"column:scan_count" doc:"扫码次数"`
	TopAt       int64  `json:"topAt" gorm:"column:top_at" doc:"置顶时间 >0则是置顶 秒"`
	BlockAt     int64  `json:"blockAt" gorm:"column:block_at" doc:"拉黑时间 >0则是拉黑 秒"`
	LastChatAt  int64  `json:"lastChatAt" gorm:"column:last_chat_at" doc:"最近聊天时间 毫秒"`
	LastMsgID   uint64 `json:"lastMsgID" doc:"最近一次由该用户发送的消息id"`
	IsOnline    bool   `json:"isOnline" doc:"是否在线"`
}
