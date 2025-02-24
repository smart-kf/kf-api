package kfbackend

import (
	"fmt"
	"time"

	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/endpoints/cron/kflog"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
	"github.com/smart-fm/kf-api/infrastructure/redis"
	"github.com/smart-fm/kf-api/pkg/utils"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

type AuthController struct {
	BaseController
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req kfbackend.LoginRequest
	if !c.BindAndValidate(ctx, &req) {
		return
	}

	reqCtx := ctx.Request.Context()

	var cardRepo repository.KFCardRepository
	card, ok, err := cardRepo.FindByCardID(reqCtx, req.CardID)
	if err != nil {
		xlogger.Error(reqCtx, "查询卡密失败", xlogger.Err(err), xlogger.Any("cardId", req.CardID))
		c.Error(ctx, err)
		return
	}

	if !ok {
		c.Error(ctx, xerrors.NewCustomError("卡密不存在"))
		return
	}

	// 判断卡密状态.
	if card.SaleStatus != constant.SaleStatusSold {
		c.Error(ctx, xerrors.NewCustomError("卡密不存在,请重试"))
		return
	}
	// 更新登录时间.
	// TODO:: 将这段逻辑，放到出售卡密的地方，这里暂时先这样处理了.
	// 不放也可以，不放的问题，有可能延长卡密时间
	card.LastLoginTime = time.Now().Unix()
	if card.LoginStatus == constant.LoginStatusUnLogin {
		card.LoginStatus = constant.LoginStatusLoginned
		if card.CardType == constant.CardTypeNormal {
			card.ExpireTime = time.Now().Unix() + int64(card.Day*86400) // 首次登录，设置卡密过期时间
		} else {
			// 测试卡.
			min := caches.BillSettingCacheInstance.GetTestingCardMinute()
			card.ExpireTime = time.Now().Unix() + int64(min*60)
		}
	} else {
		// 判断过期时间
		if card.HasExpire() {
			c.Error(ctx, xerrors.NewCustomError("卡密已过期，请续费!"))
			return
		}
	}

	// 判断密码.
	if card.Password != "" && req.Password == "" {
		c.Error(ctx, xerrors.NewCustomError("密码错误"))
		return
	}
	if card.Password != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(card.Password), []byte(req.Password)); err != nil {
			c.Error(ctx, xerrors.NewCustomError("密码错误"))
			return
		}
	}

	if err := cardRepo.UpdateOne(reqCtx, card); err != nil {
		// 这里可能是乐观锁问题.
		xlogger.Error(reqCtx, "更新卡密失败", xlogger.Err(err))
		c.Error(ctx, xerrors.NewCustomError("登录失败，请重试"))
		return
	}
	kflog.AddKFLog(card.CardID, "login", "登录成功", utils.ClientIP(ctx))

	token := uuid.New().String()

	err = caches.KfAuthCacheInstance.SetBackendToken(reqCtx, token, card.CardID)
	if err != nil {
		xlogger.Error(reqCtx, "SetBackendToken-failed", xlogger.Err(err))
		c.Error(ctx, xerrors.NewCustomError("登录失败，请重试"))
		return
	}
	// 将session 存储到 redis.
	redisClient := redis.GetRedisClient()
	redisClient.Set(reqCtx, fmt.Sprintf("kfbe.%s", token), card.CardID, 7*24*time.Hour) // 7 天token

	rsp := kfbackend.LoginResponse{
		Token:     token,
		CdnDomain: config.GetConfig().Web.CdnHost,
	}

	notice := caches.BillSettingCacheInstance.GetNotice()
	if notice.Enable {
		rsp.Notice = notice.Content
	}
	// 查找公告.
	c.Success(
		ctx, rsp,
	)
}

func (c *AuthController) Logout(ctx *gin.Context) {
	cardId := common.GetKFCardID(ctx.Request.Context())
	kflog.AddKFLog(cardId, "客户", "退出了系统", utils.ClientIP(ctx))

	c.Success(
		ctx, nil,
	)
}
