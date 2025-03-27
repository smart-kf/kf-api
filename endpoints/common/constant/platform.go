package constant

const (
	PlatformKfFe = "kf"
	PlatformKfBe = "kf-backend"
)

const (
	EventSessionId  = "sessionId"  // 新建连接、初始化sessionId事件
	EventMessage    = "message"    // 发送消息事件
	EventMessageAck = "messageAck" // 发送消息事件
	EventOnline     = "online"     // 上线事件
	EventOffline    = "offline"    // 下线事件
)

const (
	MsgTypeRead = "read"    // 读取消息事件
	MsgKeyword  = "keyword" // 读取消息事件
)

const (
	IsKf    = 1
	IsNotKf = 2
)
