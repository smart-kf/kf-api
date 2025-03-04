package kfbackend

import (
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

type QRCodeRequest struct{}
type QRCodeResponse struct {
	Domains []QRCodeDomain `json:"domains,omitempty" doc:"域名列表"`
}

type QRCodeDomain struct {
	Id        int    `json:"id" doc:"id"`
	QRCodeURL string `json:"qrCodeUrl" doc:"二维码的内容，前端拿到字符串渲染二维码"`
	Domain    string `json:"domain,omitempty" doc:"站点域名"`
	CreateAt  int64  `json:"createAt,omitempty" doc:"添加创建时间秒"`
	IsPrivate bool   `json:"isPrivate" doc:"是否是私有域名"`
	Status    int    `json:"status" doc:"域名状态: 1=正常，2=微信封禁, 3=系统封禁"`
}

type QRCodeSwitchRequest struct{}
type QRCodeSwitchResponse struct {
	QRCodeURL string `json:"qrCodeUrl" doc:"二维码的内容，前端拿到字符串渲染二维码"`
}

type QRCodeOnOffRequest struct {
	Id         int64 `json:"id" doc:"域名id"`
	Status     int   `json:"status" doc:"状态"`
	DisableOld bool  `json:"disableOld" doc:"是否停用所有老码"`
}

func (r QRCodeOnOffRequest) Validate() error {
	if r.Id == 0 && !r.DisableOld {
		return xerrors.NewCustomError("参数错误")
	}

	if r.Status != 0 {
		switch r.Status {
		case constant.QRCodeNormal:
		case constant.QRCodeDisable:
		case constant.QRCodeStopGetNewFans:
		default:
			return xerrors.NewCustomError("参数错误")
		}
	}

	return nil
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
	GuestId     string `json:"guestId" binding:"required" doc:" 即一个会话中粉丝的id"`
	PageSize    uint   `json:"pageSize,omitempty" doc:"分页大小,默认20"`
	LastMsgTime int64  `json:"lastMsgTime" doc:"最老的消息时间"`
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
	MsgId   string             `doc:"消息id" json:"msgId"`
	MsgType common.MessageType `json:"msgType" doc:"消息类型:text||image||video"`               // 消息类型
	GuestId string             `gorm:"column:guest_id" json:"guestId" doc:"客户id"`             // 客服id
	CardId  string             `gorm:"column:card_id" json:"cardId" doc:"卡密id: 只有后台有值"` // 卡密id
	Content string             `json:"content" doc:"消息的内容"`                                // 内容.
	IsKf    int                `gorm:"column:is_kf;type:tinyint(4)" json:"isKf"`                // 是否是客服
	MsgTime int64              `json:"msgTime" doc:"消息创建时间 单位秒 可用作合并时间窗口"`
}

type User struct {
	UUID        string `json:"uuid" gorm:"column:uuid;unique;type:varchar(255)" doc:"用户的uuid"`             // 用户的uuid，不用主键做业务.
	Avatar      string `json:"avatar" gorm:"column:avatar;type:varchar(255)" doc:"头像地址，存储的是相对路径"` // 头像地址，存储的是相对路径
	NickName    string `json:"nickName" gorm:"column:nick_name;type:varchar(255)" doc:"昵称"`
	RemarkName  string `json:"remarkName" gorm:"column:remark_name" doc:"备注名称"`          // 备注名称.
	Mobile      string `json:"mobile" gorm:"column:mobile" doc:"手机号"`                     // 手机号
	Comments    string `json:"comments" gorm:"column:comments" doc:"备注信息"`               // 备注信息
	IP          string `json:"ip" gorm:"column:ip;type:varchar(255)" doc:"注册ip"`           // 注册ip
	Area        string `json:"area" gorm:"column:area;type:varchar(255)" doc:"ip对应的地区"` // ip对应的地区
	UserAgent   string `json:"userAgent" gorm:"column:user_agent;type:varchar(1000)" doc:"浏览器user-agent"`
	Browser     string `json:"browser"  doc:"浏览器 Chrome/Safari/firfox/..."`                          // 浏览器 Chrome/Safari/firfox/...
	Device      string `json:"device" doc:"设备类型： iphone、android、"`                                  // 设备类型： iphone、android、
	IsProxy     int    `json:"isProxy" gorm:"column:is_proxy" doc:"是否使用了代理ip访问: 1=是，2=不是."` // 是否使用了代理ip访问: 1=是，2=不是.
	IsEmulator  int    `json:"isEmulator" gorm:"column:is_emulator" doc:"是否是模拟器 1=是，2=不是"`     // 是否是模拟器 1=是，2=不是
	Source      string `json:"source" gorm:"column:source" doc:"来源"`                                  // 来源
	OfflineAt   int64  `json:"offlineAt" gorm:"column:offline_at" doc:"离线时间 秒"`                    // ws断开链接时记录
	NetworkType string `json:"networkType" gorm:"column:network_type" doc:"网络类型:wifi/4G/5G"`        // wifi/4G/5G
	ScanCount   int64  `json:"scanCount" gorm:"column:scan_count" doc:"扫码次数"`
	TopAt       int64  `json:"topAt" gorm:"column:top_at" doc:"置顶时间 >0则是置顶 秒"`
	BlockAt     int64  `json:"blockAt" gorm:"column:block_at" doc:"拉黑时间 >0则是拉黑 秒"`
	LastChatAt  int64  `json:"lastChatAt" gorm:"column:last_chat_at" doc:"最近聊天时间 毫秒"`
	LastMsgID   uint64 `json:"lastMsgID" doc:"最近一次由该用户发送的消息id"`
	IsOnline    bool   `json:"isOnline" doc:"是否在线"`
}

type BatchSendRequest struct {
	GuestId []string `json:"guestId" doc:"客户id" binding:"required"`
	Message Message  `json:"message" doc:"消息体" binding:"required"`
}
