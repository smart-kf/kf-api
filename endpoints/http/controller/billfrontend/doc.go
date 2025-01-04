package billfrontend

import (
	"github.com/clearcodecn/swaggos"

	"github.com/smart-fm/kf-api/endpoints/http/vo/billfrontend"
)

func SwaggerDoc(group *swaggos.Group) {
	bg := group.Group("/bill-fe").Tag("计费前台")
	bg.Post("/order/create").Body(billfrontend.CreateOrderRequest{}).JSON(billfrontend.CreateOrderResponse{}).Description("创建订单")
}
