package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/skip2/go-qrcode"
	"github.com/smart-fm/kf-api/config"
	"os"
	"path/filepath"
)

const fileType = "png"

// QRCodeSize 二维码的尺寸 长396 宽396
const QRCodeSize = 396

// QRCodeMidPicSize 二维码中间图的尺寸 长91 宽91
const QRCodeMidPicSize = 91

// DrawQRCodeNX 租户不存在该二维码content时则进行绘制保存 cardID:卡密id content:识别二维码出来的内容
func DrawQRCodeNX(cardID, content string) (url string, err error) {
	if len(content) == 0 {
		return "", errors.New("待绘制二维码的信息为空")
	}

	// 加上卡密id 按租户文件夹隔离
	dir := filepath.Join(config.GetConfig().Web.StaticDir, "qrcode", cardID)

	// 创建保存地址 判断一下是否存在目录  不存在则创建
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("绘制二维码失败 %+v", err.Error())
	}

	fileName := fmt.Sprintf("%s.%s", base64.StdEncoding.EncodeToString([]byte(content)), fileType)

	filePath := filepath.Join(dir, fileName)

	info, err := os.Stat(filePath)
	if os.IsNotExist(err) || info.IsDir() || info.Size() == 0 {
		// 不存在 or 非文件 or 大小为空 则创建
		err = qrcode.WriteFile(content, qrcode.Medium, QRCodeSize, filePath)
		if err != nil {
			return "", err
		}
	}

	return filePath, nil
}
