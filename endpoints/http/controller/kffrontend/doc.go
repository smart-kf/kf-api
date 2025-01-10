package kffrontend

import (
	"github.com/clearcodecn/swaggos"

	"github.com/smart-fm/kf-api/endpoints/http/vo/kffrontend"
)

func SwaggerDoc(group *swaggos.Group) {
	bg := group.Group("/kf-fe").Tag("客服前台")
	_ = bg

	bg.Post("/qrcode/scan").Body(kffrontend.QRCodeScanRequest{}).JSON(kffrontend.QRCodeScanResponse{}).Description(
		"扫码接口",
	)
}
