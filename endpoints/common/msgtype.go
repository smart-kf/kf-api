package common

type MessageType string

const (
	MessageTypeText  MessageType = "text"  // 文本
	MessageTypeVoice MessageType = "voice" // 语音
	MessageTypeImage MessageType = "image" // 图片
	MessageTypeVideo MessageType = "video" // 视频
	MessageTypeUrl   MessageType = "link"  // 网址
	MessageTypeFile  MessageType = "file"  // 其他文件
)
