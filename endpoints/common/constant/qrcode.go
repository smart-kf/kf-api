package constant

// 二维码状态
const (
	QRCodeNormal         = 1 // 正常
	QRCodeDisable        = 2 // 失效
	QRCodeStopGetNewFans = 3 // 暂停引新粉
)

// 域名状态
const (
	DomainStatusNormal = iota + 1
	DomainStatusDisable
)
