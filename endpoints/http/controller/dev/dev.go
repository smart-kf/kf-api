package dev

import (
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/endpoints/http/middleware"
)

type DevController struct {
	middleware.BaseController
}

type PushMsgRequest struct {
	KfId    int64  `json:"kfId" doc:"客服id" binding:"required"`
	GuestId int64  `json:"guestId" doc:"客户id" binding:"required"`
	Content string `json:"content" doc:"消息内容" binding:"required"`
	MsgType string `json:"msgType" doc:"消息类型: image || text || video" binding:"required"`
	IsKF    int    `json:"isKf" doc:"是谁发的消息,1=客服发给客户，2=客户发给客服" binding:"required"`
}

func (d *DevController) PushMsg(ctx *gin.Context) {
	// var req PushMsgRequest
	// if !d.BindAndValidate(ctx, &req) {
	// 	return
	// }
	// var msg = &api.MessageTypeMsg{
	// 	MsgType:     req.MsgType,
	// 	GuestId:     req.GuestId,
	// 	GuestName:   "jerry",
	// 	GuestAvatar: "https://api.smartkf.top/static/avatar/guest.png",
	// 	KfName:      "客服小二",
	// 	KfAvatar:    "https://api.smartkf.top/static/avatar/kf.png",
	// 	MsgTime:     time.Now().Unix(),
	// 	KfId:        req.KfId,
	// 	Content:     req.Content,
	// 	IsKF:        req.IsKF,
	// }
	// svc := imMessage.NewMessageService()
	// err := svc.ReceiveMessage(ctx.Request.Context(), api.OpClientMsg, msg)
	// if err != nil {
	// 	d.Error(ctx, err)
	// 	return
	// }
	d.Success(ctx, nil)
}
