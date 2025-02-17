package kfbackend

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/cron/kflog"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
	"github.com/smart-fm/kf-api/pkg/utils"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

type GuestController struct {
	BaseController
}

func (c *GuestController) GetKfUserInfo(ctx *gin.Context) {
	var req kfbackend.GetKfUserInfoRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	reqCtx := ctx.Request.Context()

	cardId := common.GetKFCardID(reqCtx)
	var repo repository.KFUserRepository
	kfUser, ok, err := repo.FindByToken(reqCtx, cardId, req.UUID)
	if err != nil {
		c.Error(ctx, err)
		return
	}
	if !ok {
		c.Error(ctx, xerrors.NewCustomError("获取客户信息失败"))
		return
	}
	user := user2VO(ctx, kfUser)
	isOnline, err := caches.UserOnLineCacheInstance.IsUserOnline(reqCtx, cardId, kfUser.UUID)
	if err == nil {
		user.IsOnline = isOnline
	}
	c.Success(ctx, user)
}

func (c GuestController) UpdateUserInfo(ctx *gin.Context) {
	var req kfbackend.UpdateUserInfoRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	reqCtx := ctx.Request.Context()
	cardId := common.GetKFCardID(reqCtx)
	var repo repository.KFUserRepository
	kfUser, ok, err := repo.FindByToken(reqCtx, cardId, req.UUID)
	if err != nil {
		c.Error(ctx, err)
		return
	}
	if !ok {
		c.Error(ctx, xerrors.NewCustomError("获取客户信息失败"))
		return
	}
	var (
		oldName = kfUser.RemarkName
	)
	var update bool
	if req.IsUserInfo() {
		kfUser.RemarkName = req.RemarkName
		kfUser.Comments = req.Comments
		kfUser.Mobile = req.Mobile
		update = true
	}
	if req.IsBlock() {
		if req.Block == 1 {
			kfUser.BlockAt = time.Now().Unix()
		} else {
			kfUser.BlockAt = 0
		}
		update = true
	}
	if req.IsTop() {
		if req.Top == 1 {
			kfUser.TopAt = time.Now().Unix()
		} else {
			kfUser.TopAt = 0
		}
		update = true
	}

	if !update {
		c.Success(ctx, nil)
		return
	}

	if err := repo.SaveOne(ctx, kfUser); err != nil {
		c.Error(ctx, err)
		return
	}
	if oldName != "" && oldName != kfUser.RemarkName {
		kflog.AddKFLog(
			cardId,
			"客户",
			fmt.Sprintf("更新了客户信息: %s => %s", kfUser.RemarkName, kfUser.RemarkName),
			utils.ClientIP(ctx),
		)
	} else {
		kflog.AddKFLog(
			cardId,
			"客户",
			fmt.Sprintf("更新了 %s 的客户信息", kfUser.RemarkName),
			utils.ClientIP(ctx),
		)
	}
	c.Success(ctx, nil)
}
