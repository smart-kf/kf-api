package kfbackend

// TODO 临时放下 后面迁移到各个文件内

type QRCodeDomain struct {
	Domain   string `json:"domain,omitempty" doc:"站点域名"`
	HealthAt int64  `json:"healthAt,omitempty" doc:"通过健康检查的时间 毫秒"`
	CreateAt int64  `json:"createAt,omitempty" doc:"添加创建时间 毫秒"`
	Remark   string `json:"remark,omitempty" doc:"备注"`
	URL      string `json:"url,omitempty" doc:"二维码图片地址"`
}

type QRCodeOnOffRequest struct {
	OnOff        *bool `json:"onoff" doc:"开关：所有二维码的所有用户都不能进入"`
	OnOffNewUser *bool `json:"onoffNewUser" doc:"开关：老用户可进，新用户不能进"`
}
type QRCodeOnOffResponse struct{}

type ChatListRequest struct {
	SearchBy string `json:"searchBy" doc:"模糊搜索 用户id/昵称/手机号/备注"`
	ListType int    `json:"listType" doc:"列表类型 0:全部(默认) 1:消息未读 2:拉黑访客"`
	ScrollID string `json:"scrollID" doc:"滚页id 即上页最后一条会话的最近聊天时间 请求首页时不传"`
}

type ChatListResponse struct {
	Chats []Chat `json:"chats,omitempty" doc:"会话列表"`
}

type ChatType int8

const (
	ChatTypeSingle ChatType = iota // 单聊
	ChatTypeGroup                  // 群聊
)

type Chat struct {
	Type         ChatType     `json:"type" doc:"会话类型 0:单聊(默认) 1:群聊(暂不做)"`
	ExternalUser ExternalUser `json:"externalUser" doc:"外部访客信息"`
	LastChatAt   int64        `json:"lastChatAt" doc:"最近聊天时间 毫秒"`
	LastMessage  Message      `json:"LastMessage" doc:"最近一次聊天的消息内容"`
	UnreadMsgCnt int64        `json:"unreadMsgCnt" doc:"未读消息数"`
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

// ChatObjType 聊天对象的类型
type ChatObjType int8

const (
	ChatObjTypeSys          ChatObjType = iota // 系统
	ChatObjTypeExternalUser                    // 访客 即用户/粉丝
	ChatObjTypeUser                            // 员工 即客服
)

type Message struct {
	Content  Material    `json:"content" doc:"消息的内容"`
	From     string      `json:"from" doc:"发送方id"`
	FromType ChatObjType `json:"fromType" doc:"发送方类型 0:系统 1:访客 2:客服"`
	To       string      `json:"to" doc:"接收方id"`
	ToType   ChatObjType `json:"toType" doc:"接收方类型 0:系统 1:访客 2:客服"`
}

type Material struct {
	Type  MaterialType `json:"type" doc:"资源类型 0:文本 1:语音 2:图片 3:视频 4:网址 5:其他文件"`
	Text  Text         `json:"text,omitempty" doc:"文本"`
	Voice Voice        `json:"voice,omitempty" doc:"语音" `
	Image Image        `json:"image,omitempty" doc:"图片"`
	Video Video        `json:"video,omitempty" doc:"视频"`
	URL   URL          `json:"url,omitempty" doc:"网址"`
	File  File         `json:"file,omitempty" doc:"文件"`
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

type ExternalUser struct {
	Avatar   string `json:"avatar" doc:"头像"`
	NickName string `json:"nickName,omitempty" doc:"昵称"`
	IsOnline bool   `json:"isOnline,omitempty" doc:"是否在线"`
}
