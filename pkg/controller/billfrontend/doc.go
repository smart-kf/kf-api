package billfrontend

import (
	"github.com/clearcodecn/swaggos"
)

func SwaggerDoc(group *swaggos.Group) {
	bg := group.Group("/bill-fe").Tag("计费前台")
	bg.Post("/order/create").Body(CreateOrderRequest{}).JSON(CreateOrderResponse{}).Description("创建订单")
}
