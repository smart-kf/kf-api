package kffrontend

import "github.com/gin-gonic/gin"

type MsgController struct {
	BaseController
}

type MsgListRequest struct {
	LastMsgTime int64 `json:"lastMsgTime"`
}

func (q *BaseController) MsgList(ctx *gin.Context) {
	var req MsgListRequest
	if !q.BindAndValidate(ctx, &req) {
		return
	}
	
	return
}
