package kfbackend

import (
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type QRCodeRequest struct{}
type QRCodeResponse struct {
	URL           string         `json:"qrcodeUrl,omitempty" doc:"主站二维码图片地址"`
	HealthAt      int64          `json:"healthAt,omitempty" doc:"主站通过健康检查的时间 毫秒"`
	Enable        bool           `json:"enable,omitempty" doc:"启用停用状态"`
	EnableNewUser bool           `json:"enableNewUser,omitempty" doc:"启用停用新粉状态"`
	Domains       []QRCodeDomain `json:"domains,omitempty" doc:"域名列表"`
}

type QRCodeDomain struct {
	Domain   string `json:"domain,omitempty" doc:"站点域名"`
	HealthAt int64  `json:"healthAt,omitempty" doc:"通过健康检查的时间 毫秒"`
	CreateAt int64  `json:"createAt,omitempty" doc:"添加创建时间 毫秒"`
	Remark   string `json:"remark,omitempty" doc:"备注"`
	URL      string `json:"url,omitempty" doc:"二维码图片地址"`
}

type QRCodeSwitchRequest struct{}
type QRCodeSwitchResponse struct {
	URL      string `json:"qrcodeUrl,omitempty" doc:"主站二维码图片地址"`
	HealthAt int64  `json:"healthAt,omitempty" doc:"主站通过健康检查的时间 毫秒"`
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
	Chats []Chat `json:"chats,omitempty" doc:"会话列表"`
}

type ChatType int8

const (
	ChatTypeSingle ChatType = iota // 单聊
	ChatTypeGroup                  // 群聊
)

type Chat struct {
	Type         ChatType `json:"type" doc:"会话类型 0:单聊(默认) 1:群聊(暂不做)"`
	User         User     `json:"user" doc:"访客信息"`
	LastChatAt   int64    `json:"lastChatAt" doc:"最近聊天时间 毫秒"`
	LastMessage  Message  `json:"lastMessage" doc:"最近一次聊天的消息内容"`
	UnreadMsgCnt int64    `json:"unreadMsgCnt" doc:"未读消息数"`
}

type MsgListRequest struct {
	FromTos              []string `json:"fromTos" doc:"发送方id和接收方id数组 即一个会话中客服和粉丝的ids"`
	common.ScrollRequest `json:",inline"`
}

type MsgListResponse struct {
	Messages []Message `json:"messages" doc:"消息列表"`
}

type ReadMsgRequest struct {
	MsgIDs []uint64 `json:"msgIDs" doc:"已读的消息ids"`
}

type ReadMsgResponse struct {
}

type BatchOpUserRequest struct {
	UserIDs []uint `json:"userIDs" doc:"粉丝ids"`
	Op      UserOp `json:"op" doc:"操作 1:置顶 2:取消置顶 3:拉黑 4:取消拉黑"`
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
	ID         uint   `json:"id" doc:"粉丝id"`
	RemarkName string `json:"remarkName" doc:"备注名称"`
	Mobile     string `json:"mobile" doc:"手机号"`
	Comments   string `json:"comments" doc:"备注信息"`
}

type UpdateUserResponse struct {
}

type Message struct {
	ID       uint64          `json:"id" doc:"消息自增id 可用作排序"`
	Content  string          `json:"content" doc:"消息的内容"`
	From     string          `json:"from" doc:"发送方id"`
	FromType dao.ChatObjType `json:"fromType" doc:"发送方类型 0:系统 1:访客 2:客服"`
	To       string          `json:"to" doc:"接收方id"`
	ToType   dao.ChatObjType `json:"toType" doc:"接收方类型 0:系统 1:访客 2:客服"`
	ReadAt   int64           `json:"readAt" doc:"接收方是否已读 如果已读则存有已读时间 单位秒"`
	CreateAt int64           `json:"create_at" doc:"消息创建时间 单位秒 可用作合并时间窗口"`
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
	UUID        string `json:"uuid,omitempty" doc:"访客/粉丝/用户的id"`
	Avatar      string `json:"avatar,omitempty" doc:"头像 相对地址"`
	NickName    string `json:"nickName,omitempty" doc:"昵称"`
	IsOnline    bool   `json:"isOnline,omitempty" doc:"是否在线"`
	RemarkName  string `json:"remarkName,omitempty" doc:"备注名称"`
	Mobile      string `json:"mobile,omitempty" doc:"手机号"`
	Comments    string `json:"comments,omitempty" doc:"备注信息"`
	IP          string `json:"ip,omitempty" doc:"注册ip"`
	Area        string `json:"area,omitempty" doc:"ip对应的地区"`
	Browser     string `json:"browser,omitempty" doc:"浏览器 Chrome/Safari/firfox/..."`
	Device      string `json:"device,omitempty" doc:"设备类型： iphone/android/..."`
	IsProxy     int    `json:"isProxy,omitempty" doc:"是否使用了代理ip访问: 1=是，2=不是."`
	IsEmulator  int    `json:"isEmulator,omitempty" doc:"是否是模拟器 1=是，2=不是"`
	Source      string `json:"source,omitempty" doc:"来源"`
	OfflineAt   int64  `json:"offlineAt,omitempty" doc:"离线时间 秒"`
	NetworkType string `json:"networkType,omitempty" doc:"网络类型 wifi/4G/5G"`
	ScanCount   int64  `json:"scanCount,omitempty" doc:"扫码次数"`
	TopAt       int64  `json:"topAt,omitempty" doc:"置顶时间 >0则是置顶 秒"`
	BlockAt     int64  `json:"blockAt,omitempty" doc:"拉黑时间 >0则是拉黑 秒"`
	LastChatAt  int64  `json:"lastChatAt,omitempty" doc:"最近聊天时间 毫秒"`
}
