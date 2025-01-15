package kfbackend

import (
	"errors"
	"fmt"
	"time"

	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"

	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/pkg/utils"
)

type QRCodeController struct {
	BaseController
}

func (c *QRCodeController) List(ctx *gin.Context) {
	var req kfbackend.QRCodeRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	reqCtx := ctx.Request.Context()

	cardID := common.GetKFCardID(ctx)

	var kfsetting repository.KFSettingRepository
	setting, ok, err := kfsetting.GetByCardID(reqCtx, cardID)
	if err != nil {
		xlogger.Error(reqCtx, "查询客服设置失败", xlogger.Err(err), xlogger.Any("cardId", cardID))
		c.Error(ctx, err)
		return
	}

	enable := true
	enableNewUser := true

	if ok {
		enable = setting.QRCodeEnabled
		enableNewUser = setting.QRCodeEnabledNewUser
	}

	var qrCodeDomainRepo repository.QRCodeDomainRepository
	// 获取二维码列表.
	qrCode, exist, err := qrCodeDomainRepo.FindByCardID(reqCtx, cardID)
	if err != nil {
		xlogger.Error(reqCtx, "查询客服设置失败", xlogger.Err(err), xlogger.Any("cardId", cardID))
		c.Error(ctx, err)
		return
	}
	if !exist {
		c.Error(ctx, errors.New("查找二维码失败"))
		return
	}

	// 获取域名列表
	qrCodeDomain, err := qrCodeDomainRepo.FindDomain(reqCtx, int(qrCode.ID))
	if err != nil {
		c.Error(ctx, err)
		return
	}
	firstDomain := ""
	if len(qrCodeDomain) > 0 {
		firstDomain = qrCodeDomain[0].Domain
	}
	var domains []kfbackend.QRCodeDomain
	for _, item := range qrCodeDomain {
		domains = append(
			domains, kfbackend.QRCodeDomain{
				Domain:    item.Domain,
				CreateAt:  item.CreatedAt.Unix(),
				IsPrivate: item.IsPrivate,
			},
		)
	}
	c.Success(
		ctx, kfbackend.QRCodeResponse{
			QRCodeURL:     fmt.Sprintf("%s%s", firstDomain, qrCode.Path),
			HealthAt:      0,
			Enable:        enable,
			EnableNewUser: enableNewUser,
			Version:       qrCode.Version,
			Domains:       domains,
		},
	)
}

// Switch 更换二维码图片
func (c *QRCodeController) Switch(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	cardID := common.GetKFCardID(ctx)
	var qrCodeDomainRepo repository.QRCodeDomainRepository
	// 获取二维码列表.
	qrCode, exist, err := qrCodeDomainRepo.FindByCardID(reqCtx, cardID)
	if err != nil {
		xlogger.Error(reqCtx, "Switch-err", xlogger.Err(err), xlogger.Any("cardId", cardID))
		c.Error(ctx, err)
		return
	}
	if !exist {
		c.Error(ctx, errors.New("查找二维码失败"))
		return
	}
	qrCode.ID = 0
	qrCode.CreatedAt = time.Now()
	qrCode.Version++
	qrCode.Path = utils.RandomPath()
	err = qrCodeDomainRepo.CreateOne(reqCtx, qrCode)
	if err != nil {
		xlogger.Error(reqCtx, "Switch-error2", xlogger.Err(err), xlogger.Any("cardId", cardID))
		c.Error(ctx, err)
		return
	}
	qrCodeDomain, err := qrCodeDomainRepo.FindDomain(reqCtx, int(qrCode.ID))
	if err != nil {
		c.Error(ctx, err)
		return
	}
	firstDomain := ""
	if len(qrCodeDomain) > 0 {
		firstDomain = qrCodeDomain[0].Domain
	}
	c.Success(
		ctx, kfbackend.QRCodeSwitchResponse{
			QRCodeURL: fmt.Sprintf("%s%s", firstDomain, qrCode.Path),
		},
	)
}

// OnOff 二维码功能开关
func (c *QRCodeController) OnOff(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	cardID := common.GetKFCardID(ctx)

	var req kfbackend.QRCodeOnOffRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	var kfsetting repository.KFSettingRepository
	setting, ok, err := kfsetting.GetByCardID(reqCtx, cardID)
	if err != nil {
		xlogger.Error(reqCtx, "查询客服设置失败", xlogger.Err(err), xlogger.Any("cardId", cardID))
		c.Error(ctx, err)
		return
	}

	if !ok {
		setting = &dao.KFSettings{}
		setting.CardID = cardID
	}

	if req.OnOffNewUser != nil {
		setting.QRCodeEnabledNewUser = *req.OnOffNewUser
	}

	if req.OnOff != nil {
		setting.QRCodeEnabled = *req.OnOff
	}

	err = kfsetting.SaveOne(reqCtx, setting)
	if err != nil {
		xlogger.Error(reqCtx, "保存客服设置失败", xlogger.Err(err), xlogger.Any("cardId", cardID))
		c.Error(ctx, err)
		return
	}

	c.Success(ctx, kfbackend.QRCodeOnOffResponse{})
}
