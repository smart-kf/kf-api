package bill

import (
	"github.com/clearcodecn/swaggos"
)

func SwaggerDoc(group *swaggos.Group) {
	bg := group.Group("/bill").Tag("计费后台-未授权接口")
	bg.Post("/login").Body(LoginRequest{}).JSON(LoginResponse{})

	billGroup := group.Group("/bill")

	cardC := billGroup.Group("/card").Tag("计费后台-卡密管理")
	cardC.Post("/batch-add").Body(BatchAddCardRequest{}).JSON(BatchAddResponse{})
}
