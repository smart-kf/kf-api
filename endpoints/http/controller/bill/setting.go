package bill

import (
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/http/middleware"
	"github.com/smart-fm/kf-api/endpoints/http/vo/bill"
	sdk "github.com/smart-fm/kf-api/pkg/usdtpayment"
	"github.com/smart-fm/kf-api/pkg/xerrors"
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

	setting.DomainPrice = req.DomainPrice

	err = repo.UpsertSettings(reqCtx, setting, true)
	if err != nil {
		c.Error(ctx, err)
		return
	}

	// 更新缓存.
	caches.BillSettingCacheInstance.DeleteCache(ctx.Request.Context())
	c.Success(ctx, nil)
}

func (c SettingController) AddressList(ctx *gin.Context) {
	conf := config.GetConfig().Payment
	client := sdk.NewUsdtPaymentClient(conf.Host, conf.Token, 30*time.Second)

	rsp, err := client.Management(
		ctx.Request.Context(), &sdk.ListAddressRequest{
			AppId: conf.AppId,
		},
	)

	if err != nil {
		c.Error(ctx, err)
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.String(200, string(rsp))
}

type UpsertAddressRequest struct {
	Id      int    `json:"id"`
	Enable  bool   `json:"enable"`
	Address string `json:"address" binding:"required"`
	Remark  string `json:"remark"`
}

func (c SettingController) UpsertAddress(ctx *gin.Context) {
	var req UpsertAddressRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	conf := config.GetConfig().Payment
	client := sdk.NewUsdtPaymentClient(conf.Host, conf.Token, 30*time.Second)

	rsp, err := client.Management(
		ctx.Request.Context(), &sdk.UpsertAddressRequest{
			Id:      req.Id,
			AppId:   conf.AppId,
			Enable:  req.Enable,
			Address: req.Address,
			Remark:  req.Remark,
		},
	)

	if err != nil {
		c.Error(ctx, err)
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.String(200, string(rsp))
}

type DeleteAddressRequest struct {
	Id    int    `json:"id"`
	AppId string `json:"appId"`
}

func (c SettingController) DelAddress(ctx *gin.Context) {
	var req DeleteAddressRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	conf := config.GetConfig().Payment
	client := sdk.NewUsdtPaymentClient(conf.Host, conf.Token, 30*time.Second)

	rsp, err := client.Management(
		ctx.Request.Context(), &sdk.DeleteAddressRequest{
			Id:    req.Id,
			AppId: conf.AppId,
		},
	)

	if err != nil {
		c.Error(ctx, err)
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.String(200, string(rsp))
}

func (c SettingController) GetTron(ctx *gin.Context) {
	conf := config.GetConfig().Payment
	client := sdk.NewUsdtPaymentClient(conf.Host, conf.Token, 30*time.Second)

	rsp, err := client.Management(
		ctx.Request.Context(), &sdk.GetTronRequest{},
	)

	if err != nil {
		c.Error(ctx, err)
		return
	}

	ctx.Header("Content-Type", "application/json")
	ctx.String(200, string(rsp))
}

type UpdateTronRequest struct {
	ApiKey      string `json:"apiKey" binding:"required"` // 可能有多个，以 , 分割.
	Proxy       string `json:"proxy"`
	Timeout     int    `json:"timeout"`
	TronNetwork string `json:"tron_network" binding:"required"`
	CronSecond  int    `json:"cron_second" binding:"required"` // 定时任务执行秒数间隔
}

func (c SettingController) UpsertTron(ctx *gin.Context) {
	var req UpdateTronRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	conf := config.GetConfig().Payment
	client := sdk.NewUsdtPaymentClient(conf.Host, conf.Token, 30*time.Second)
	rsp, err := client.Management(
		ctx.Request.Context(), &sdk.UpdateTronRequest{
			ApiKey:      req.ApiKey,
			Proxy:       req.Proxy,
			Timeout:     req.Timeout,
			TronNetwork: req.TronNetwork,
			CronSecond:  req.CronSecond,
		},
	)
	if err != nil {
		c.Error(ctx, err)
		return
	}
	ctx.Header("Content-Type", "application/json")
	ctx.String(200, string(rsp))
}

type UpdatePasswordRequest struct {
	OldPassword    string `json:"old_password" binding:"required"`
	NewPassword    string `json:"new_password"  binding:"required"`
	RepeatPassword string `json:"repeat_password"  binding:"required"`
}

func (c SettingController) UpdatePassword(ctx *gin.Context) {
	var req UpdatePasswordRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}
	reqCtx := ctx.Request.Context()

	if req.NewPassword != req.RepeatPassword {
		c.Error(ctx, xerrors.NewCustomError("两次密码不一致"))
		return
	}

	account := middleware.GetBillAccount(ctx)

	var br repository.BillAccountRepository
	acc, ok, err := br.FindOneByUsername(reqCtx, account.Username)
	if err != nil {
		c.Error(ctx, err)
		return
	}
	if !ok {
		c.Error(ctx, xerrors.NewCustomError("账号不存在"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(req.OldPassword)); err != nil {
		c.Error(ctx, xerrors.NewCustomError("密码错误"))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)

	if err != nil {
		return
	}

	acc.Password = string(hash)

	err = br.Save(reqCtx, acc)

	c.Success(ctx, nil)
}
