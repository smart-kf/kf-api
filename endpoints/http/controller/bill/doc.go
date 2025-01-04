package bill

import (
	"github.com/clearcodecn/swaggos"

	"github.com/smart-fm/kf-api/endpoints/http/vo/bill"
	"github.com/smart-fm/kf-api/pkg/common"
)

func SwaggerDoc(group *swaggos.Group) {
	bg := group.Group("/bill").Tag("计费后台-未授权接口")
	bg.Post("/login").Body(bill.LoginRequest{}).JSON(bill.LoginResponse{})

	billGroup := group.Group("/bill")

	cardC := billGroup.Group("/card").Tag("计费后台-卡密管理")
	cardC.Post("/batch-add").Body(bill.BatchAddCardRequest{}).JSON(bill.BatchAddResponse{}).Description("批量添加卡片")
	cardC.Post("/updateStatus").Body(bill.UpdateStatusRequest{}).JSON(common.EmptyResponse{}).Description("更新卡片出售状态")
	cardC.Post("/list").Body(bill.ListCardRequest{}).JSON(bill.ListCardResponse{}).Description("卡密列表")
}
