package bill

import (
	"github.com/clearcodecn/swaggos"
	"std-api/pkg/common"
)

func SwaggerDoc(group *swaggos.Group) {
	bg := group.Group("/bill").Tag("计费后台-未授权接口")
	bg.Post("/login").Body(LoginRequest{}).JSON(LoginResponse{})

	billGroup := group.Group("/bill")

	cardC := billGroup.Group("/card").Tag("计费后台-卡密管理")
	cardC.Post("/batch-add").Body(BatchAddCardRequest{}).JSON(BatchAddResponse{}).Description("批量添加卡片")
	cardC.Post("/updateStatus").Body(UpdateStatusRequest{}).JSON(common.EmptyResponse{}).Description("更新卡片出售状态")
	cardC.Post("/list").Body(ListCardRequest{}).JSON(ListCardResponse{}).Description("卡密列表")
}
