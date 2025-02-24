package bill

import (
	"github.com/clearcodecn/swaggos"

	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/http/vo/bill"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

func SwaggerDoc(group *swaggos.Group) {
	bg := group.Group("/bill").Tag("计费后台-未授权接口")
	bg.Post("/login").Body(bill.LoginRequest{}).JSON(bill.LoginResponse{})

	billGroup := group.Group("/bill")

	cardC := billGroup.Group("/card").Tag("计费后台-卡密管理")
	cardC.Post("/batch-add").Body(bill.BatchAddCardRequest{}).JSON(bill.BatchAddResponse{}).Description("批量添加卡片")
	cardC.Post("/updateStatus").Body(bill.UpdateStatusRequest{}).JSON(common.EmptyResponse{}).Description("更新卡片出售状态")
	cardC.Post("/list").Body(bill.ListCardRequest{}).JSON(bill.ListCardResponse{}).Description("卡密列表")

	domainC := billGroup.Group("/domain").Tag("计费后台-域名管理")
	domainC.Post("/add").Body(bill.AddDomainRequest{}).JSON(dao.BillDomain{}).Description("添加域名")
	domainC.Post("/list").Body(bill.ListDomainRequest{}).JSON(bill.ListDomainResponse{}).Description("域名列表")
	domainC.Post("/del").Body(bill.DeleteDomainRequest{}).JSON(common.EmptyResponse{}).Description("删除域名")

	orderC := billGroup.Group("/order").Tag("计费后台-订单管理")
	orderC.Post("/list").Body(bill.ListOrderRequest{}).JSON(bill.ListOrderResponse{}).Description("订单列表")

	settingC := billGroup.Group("/setting").Tag("计费后台-系统配置")
	settingC.Get("/get").JSON(bill.SettingRequest{}).Description("获取配置")
	settingC.Post("/update").Body(bill.SettingRequest{}).JSON(bill.SettingRequest{}).Description("更新配置")
}
