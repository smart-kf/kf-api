package kffrontend

import (
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kffrontend"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type MsgController struct {
	BaseController
}

func (c *BaseController) MsgList(ctx *gin.Context) {
	var req kffrontend.MsgListRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	var (
		reqCtx  = ctx.Request.Context()
		repo    repository.KFMessageRepository
		kfToken = common.GetKFToken(reqCtx)
	)

	cardId, err := caches.KfAuthCacheInstance.GetFrontToken(ctx.Request.Context(), kfToken)

	if err != nil {
		c.Error(ctx, err)
		return
	}

	option := &repository.ListMsgOption{
		CardID:   cardId,
		GuestId:  kfToken,
		PageSize: 30,
	}

	if req.LastMsgTime != 0 {
		option.LastMsgTime = time.Unix(req.LastMsgTime, 0)
	}

	resp, err := repo.List(
		ctx,
		option,
	)

	sort.Slice(
		resp, func(i, j int) bool {
			return !resp[i].CreatedAt.After(resp[j].CreatedAt)
		},
	)

	var res []*kfbackend.Message
	lo.ForEach(
		resp, func(item *dao.KFMessage, index int) {
			res = append(res, msg2VO(item))
		},
	)
	c.Success(
		ctx, kfbackend.MsgListResponse{
			Messages: res,
		},
	)
	return
}

func msg2VO(m *dao.KFMessage) *kfbackend.Message {
	vo := &kfbackend.Message{
		MsgId:   m.MsgId,
		MsgType: m.MsgType,
		GuestId: m.GuestId,
		Content: m.Content,
		IsKf:    m.IsKf,
		MsgTime: m.CreatedAt.Unix(),
	}
	return vo
}
