package kfbackend

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type LogController struct {
	BaseController
}

func (c *LogController) List(ctx *gin.Context) {
	var req kfbackend.LogRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	reqCtx := ctx.Request.Context()
	db := mysql.GetDBFromContext(reqCtx)

	var list []*dao.KFLog
	if req.BeginTime != 0 {
		db = db.Where("created_at >= ?", time.Unix(req.BeginTime/1000, 0))
	}
	if req.EndTime != 0 {
		db = db.Where("created_at <= ?", time.Unix(req.EndTime/1000, 0))
	}
	if req.Function != "" {
		db = db.Where("handle_func = ?", req.Function)
	}
	req.OrderBy = "id"
	req.Asc = false

	list, total, err := common.Paginate[*dao.KFLog](db, &req.PageRequest)
	if err != nil {
		c.Error(ctx, err)
		return
	}

	var rsp kfbackend.ListLogResponse
	lo.ForEach(
		list, func(item *dao.KFLog, index int) {
			rsp.List = append(
				rsp.List, &kfbackend.LogResponse{
					Id:         int64(item.ID),
					HandleFunc: item.HandleFunc,
					Content:    item.Content,
					CreateTime: item.CreatedAt.UnixMilli(),
					Ip:         item.Ip,
				},
			)
		},
	)
	rsp.Total = total
	c.Success(ctx, rsp)
}
