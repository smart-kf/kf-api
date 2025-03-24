package constant

import "time"

// 二维码状态
const (
	QRCodeNormal         = 1 // 正常
	QRCodeDisable        = 2 // 失效
	QRCodeStopGetNewFans = 3 // 暂停引新粉
)

// 域名状态
const (
	DomainStatusNormal  = iota + 1
	DomainStatusDisable = 2
	DomainStatusLocked  = 3
)

const (
	DomainExpireTime = 10 * time.Second
)

const (
	QRCodeFilterClose              = 1 // 关闭
	QRCodeFilterRoom               = 2 // 过滤机房
	QRCodeFilterNonMainland        = 3 // 过滤非大陆
	QRCodeFilterRoomAndNonMainland = 4 // 过滤机房及非大陆
)
