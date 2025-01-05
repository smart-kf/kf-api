package kfbackend

import (
	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/http/middleware"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
)

type ChatController struct {
	BaseController
}

func (c *ChatController) List(ctx *gin.Context) {
	var req kfbackend.ChatListRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	reqCtx := ctx.Request.Context()
	cardID := middleware.GetKFCardID(ctx)

	var repo repository.KFExternalUserRepository
	extUsers, cnt, err := repo.List(reqCtx, &repository.ListExtUserOption{})
	if err != nil {
		xlogger.Error(reqCtx, "查询客服设置失败", xlogger.Err(err), xlogger.Any("cardId", cardID))
		c.Error(ctx, err)
		return
	}

	_ = extUsers

	_ = cnt
}
