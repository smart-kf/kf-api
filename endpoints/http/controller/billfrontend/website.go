package billfrontend

import (
	"crypto/md5"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type WebsiteController struct{}

func (c *WebsiteController) Index(ctx *gin.Context) {
	ctx.HTML(200, "index.html", nil)
}

func (c *WebsiteController) Package(ctx *gin.Context) {
	pkgs := config.GetConfig().CardPackages
	var (
		dailyPrice   int64
		weeklyPrice  int64
		monthlyPrice int64
	)
	for _, p := range pkgs {
		if p.Name == "日卡" {
			dailyPrice = p.Price
		}
		if p.Name == "周卡" {
			weeklyPrice = p.Price
		}
		if p.Name == "月卡" {
			monthlyPrice = p.Price
		}
	}
	ctx.HTML(
		200, "package.html", gin.H{
			"dailyPrice":   dailyPrice,
			"weeklyPrice":  weeklyPrice,
			"monthlyPrice": monthlyPrice,
		},
	)
}

func (c *WebsiteController) Order(ctx *gin.Context) {
	packageId := ctx.Query("packageId")
	if packageId == "" {
		packageId = "daily"
	}
	if !lo.Contains([]string{"daily", "weekly", "monthly"}, packageId) {
		packageId = "daily"
	}

	ctx.HTML(
		200, "order.html", gin.H{
			"packageId": packageId,
		},
	)
}

func (c *WebsiteController) PaySuccess(ctx *gin.Context) {
	tradeId := ctx.Query("tradeId")
	orderId := ctx.Query("orderId")
	status := ctx.Query("status")
	sign := ctx.Query("sign")
	if tradeId == "" || orderId == "" || status == "" || sign == "" {
		ctx.Redirect(302, "/")
		return
	}
	if !validateSign(tradeId, orderId, status, sign) {
		ctx.Redirect(302, "/")
		return
	}

	db := mysql.GetDBFromContext(ctx.Request.Context())

	var order dao.Orders
	if err := db.Where("order_no = ?", orderId).First(&order).Error; err != nil {
		ctx.Error(err)
		return
	}

	ctx.HTML(
		200, "order-mail.html", gin.H{
			"order":            order,
			"KfManagerAddress": config.GetConfig().Web.KfManagerAddress,
		},
	)
}

func validateSign(tradeId, orderId string, status string, sign string) bool {
	token := config.GetConfig().Payment.Token

	var keys = []string{
		"orderId=" + orderId,
		"tradeId=" + tradeId,
		"status=" + status,
	}
	x := md5.New()
	for _, k := range keys {
		x.Write([]byte(k))
	}
	x.Write([]byte(token))
	a := x.Sum(nil)
	mySign := fmt.Sprintf("%x", a)
	return mySign == sign
}
