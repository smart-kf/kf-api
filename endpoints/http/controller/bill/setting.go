package bill

import (
	"github.com/gin-gonic/gin"

	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/http/vo/bill"
)

type SettingController struct {
	BaseController
}

func (c SettingController) Get(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	var repo repository.BillSettingRepository
	setting, err := repo.GetSetting(reqCtx)
	if err != nil {
		c.Error(ctx, err)
		return
	}
	c.Success(ctx, setting)
}

func (c SettingController) Set(ctx *gin.Context) {
	var req bill.SettingRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	// 更新缓存.
	caches.BillSettingCacheInstance.DeleteCache(ctx.Request.Context())

	reqCtx := ctx.Request.Context()
	var repo repository.BillSettingRepository
	setting, err := repo.GetSetting(reqCtx)
	if err != nil {
		c.Error(ctx, err)
		return
	}
	setting.DailyPackage.Price = req.DailyPackage.Price
	setting.DailyPackage.Days = req.DailyPackage.Days

	setting.WeeklyPackage.Days = req.WeeklyPackage.Days
	setting.WeeklyPackage.Days = req.WeeklyPackage.Days

	setting.MonthlyPackage.Days = req.MonthlyPackage.Days
	setting.MonthlyPackage.Days = req.MonthlyPackage.Days

	setting.Payment.Token = req.Payment.Token
	setting.Payment.PayUrl = req.Payment.PayUrl
	setting.Payment.AppId = req.Payment.AppId
	setting.Payment.Email = req.Payment.Email

	setting.Notice.Enable = req.Notice.Enable
	setting.Notice.Content = req.Notice.Content

	err = repo.UpsertSettings(reqCtx, setting, true)
	if err != nil {
		c.Error(ctx, err)
		return
	}

	// 更新缓存.
	caches.BillSettingCacheInstance.DeleteCache(ctx.Request.Context())
	c.Success(ctx, nil)
}
