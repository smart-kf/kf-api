package factory

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/pkg/utils"
)

func FactoryNewKfUser(cardMainId int64, cardId string, ip string) *dao.KfUser {
	uid := strings.ReplaceAll(uuid.New().String(), "-", "")
	token := fmt.Sprintf("%d|%s", cardMainId, uid)
	nickName := strings.ToUpper(utils.RandomWord(1)) + utils.RandomNumber(10)
	user := dao.KfUser{
		CardID:     cardId,
		UUID:       token,
		Avatar:     "/static/avatar/guest.png", // 先写死.
		NickName:   nickName,                   // 生成随机名称.
		RemarkName: "备注名称123",
		Mobile:     "13300991111",
		Comments:   "这里是写死的长备注",
		IP:         "127.0.0.1", // 先写死
		Area:       "四川成都电信",    // 先写死
		OfflineAt:  0,
		Device:     "Iphone",
		IsProxy:    0,
		Source:     "",
	}
	return &user
}

func FactoryParseUserToken(token string) (cardMainId int64, err error) {
	arr := strings.Split(token, "|")
	if len(arr) != 2 {
		return 0, errors.New("wrong token parts1")
	}
	cardMainId, err = strconv.ParseInt(arr[0], 10, 64)
	if err != nil {
		return 0, err
	}
	return
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
