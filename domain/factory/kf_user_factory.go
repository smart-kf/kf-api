package factory

import (
	"regexp"
	"strings"

	"github.com/google/uuid"

	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/pkg/utils"
)

func FactoryNewKfUser(cardId string, ip string) dao.KfUsers {
	token := uuid.New().String()
	nickName := strings.ToUpper(utils.RandomWord(1)) + utils.RandomNumber(10)
	user := dao.KfUsers{
		CardID:        cardId,
		UUID:          token,
		Avatar:        "/static/avatar/guest.png", // 先写死.
		Nickname:      nickName,                   // 生成随机名称.
		RemarkName:    "",
		Mobile:        "",
		Comments:      "",
		IP:            "127.0.0.1", // 先写死
		Area:          "四川成都电信",    // 先写死
		OfflineTime:   0,
		Device:        "",
		Browser:       "",
		ScanQRCodeCnt: 0,
		IsTop:         0,
		IsBlack:       0,
		IsSimulator:   0,
		IsProxy:       0,
		Source:        "",
	}
}

func parseUserAgent(ua string) (string, string, string, bool, bool) {
	// 定义正则表达式
	iosRegex := regexp.MustCompile(`iPhone; CPU iPhone OS ([\d_]+)`)
	androidRegex := regexp.MustCompile(`Android.*; ([^;]+)`)
	emulatorRegex := regexp.MustCompile(`(sdk|emulator|Genymotion)`)

	var deviceType, brand, iosVersion string
	isEmulator := false
	isWeChat := false

	if iosMatches := iosRegex.FindStringSubmatch(ua); len(iosMatches) > 1 {
		deviceType = "iPhone"
		iosVersion = strings.ReplaceAll(iosMatches[1], "_", ".") // 将版本中的下划线替换为点
	} else if androidMatches := androidRegex.FindStringSubmatch(ua); len(androidMatches) > 1 {
		deviceType = "Android"
		brand = strings.TrimSpace(androidMatches[1])
	}

	if emulatorRegex.MatchString(ua) {
		isEmulator = true
	}

	// 检查是否为微信浏览器
	if strings.Contains(ua, "MicroMessenger") {
		isWeChat = true
	}

	return deviceType, brand, iosVersion, isEmulator, isWeChat
}
