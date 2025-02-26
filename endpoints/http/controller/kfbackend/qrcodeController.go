package kfbackend

import (
	"errors"
	"fmt"

	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/endpoints/cron/kflog"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
	"github.com/smart-fm/kf-api/pkg/utils"
	"github.com/smart-fm/kf-api/pkg/xerrors"
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
	var qrCodeDomainRepo repository.QRCodeDomainRepository
	// 获取域名列表
	qrCodeDomain, err := qrCodeDomainRepo.FindDomain(reqCtx, cardID)
	if err != nil {
		c.Error(ctx, err)
		return
	}
	var domains []kfbackend.QRCodeDomain
	for _, item := range qrCodeDomain {
		domains = append(
			domains, kfbackend.QRCodeDomain{
				Id:        int(item.ID),
				QRCodeURL: item.Path,
				Domain:    item.Domain,
				CreateAt:  item.CreatedAt.Unix(),
				IsPrivate: item.IsPrivate,
				Status:    item.Status,
			},
		)
	}
	c.Success(
		ctx, kfbackend.QRCodeResponse{
			Domains: domains,
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
	qrCode.Version++
	qrCode.ID = 0
	qrCode.Status = constant.QRCodeNormal
	qrCode.Path = "/s/" + utils.RandomPath()

	err = qrCodeDomainRepo.CreateOne(reqCtx, qrCode)
	if err != nil {
		xlogger.Error(reqCtx, "Switch-error2", xlogger.Err(err), xlogger.Any("cardId", cardID))
		c.Error(ctx, err)
		return
	}

	kflog.AddKFLog(cardID, "二维码", "更换了二维码图片", utils.ClientIP(ctx))

	c.Success(
		ctx, kfbackend.QRCodeSwitchResponse{
			QRCodeURL: fmt.Sprintf("%s%s", qrCode.Domain, qrCode.Path),
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
	var qrCodeDomainRepo repository.QRCodeDomainRepository

	if req.DisableOld {
		domain, ok, err := qrCodeDomainRepo.FindByCardID(ctx, cardID)
		if err != nil {
			c.Error(ctx, err)
			return
		}
		if !ok {
			c.Error(ctx, xerrors.NewCustomError("查找数据失败"))
			return
		}
		if err := qrCodeDomainRepo.DisableOld(reqCtx, cardID, domain.Domain, int64(domain.ID)); err != nil {
			c.Error(ctx, xerrors.NewCustomError("操作失败"))
			return
		}
		c.Success(ctx, nil)
		return
	}

	if req.Id <= 0 {
		c.Success(ctx, nil)
		return
	}

	_, ok, err := qrCodeDomainRepo.FindByIdAndCardID(reqCtx, req.Id, cardID)
	if err != nil {
		c.Error(ctx, err)
		return
	}
	if !ok {
		c.Error(ctx, xerrors.NewCustomError("查找数据失败"))
		return
	}

	err = qrCodeDomainRepo.UpdateStatus(reqCtx, cardID, req.Id, req.Status)
	if err != nil {
		c.Error(ctx, xerrors.NewCustomError("操作失败"))
		return
	}

	kflog.AddKFLog(cardID, "二维码", "更新了二维码开关", utils.ClientIP(ctx))
	c.Success(ctx, kfbackend.QRCodeOnOffResponse{})
}
