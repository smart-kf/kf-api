package kffrontend

import (
	"github.com/clearcodecn/swaggos"

	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kffrontend"
)

func SwaggerDoc(group *swaggos.Group) {
	bg := group.Group("/kf-fe").Tag("客服前台")
	_ = bg

	bg.Post("/qrcode/scan").Body(kffrontend.QRCodeScanRequest{}).JSON(kffrontend.QRCodeScanResponse{}).Description(
		"扫码接口",
	)

	bg.Post("/qrcode/check").Body(kffrontend.QRCodeScanRequest{}).JSON(common.EmptyResponse{}).Description(
		"检测接口",
	).Summary("返回错误，代表不能继续访问")

	bg.Post("/msg/list").Body(kffrontend.MsgListRequest{}).JSON(kfbackend.MsgListResponse{}).Description(
		"前台消息列表",
	)

	bg.Post("/upload").Tag("公共接口").
		FormFile(
			"file",
			swaggos.Attribute{Required: true, Description: "文件, 不超过32M"},
		).
		Form("fileType", swaggos.Attribute{Required: true, Description: "文件类型: image || video"}).
		JSON(kfbackend.UploadResponse{}).Description("文件上传")
}
