package kfbackend

import (
	"context"
	xlogger "github.com/clearcodecn/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/make-money-fast/captcha"
	"golang.org/x/crypto/bcrypt"
	"std-api/config"
	"std-api/pkg/caches"
	"std-api/pkg/constant"
	"std-api/pkg/db"
	"std-api/pkg/repository"
	"std-api/pkg/xerrors"
	"time"
)

type LoginRequest struct {
	CardID      string `json:"cardID" binding:"required" doc:"卡密id" validate:"required"`
	Password    string `json:"password" doc:"密码可选"`
	CaptchaID   string `json:"captchaId" doc:"验证码id,通过验证码接口获取"`
	CaptchaCode string `json:"captchaCode" doc:"验证码"`
}

func (r *LoginRequest) Validate(ctx context.Context) error {
	if caches.BillSettingCacheInstance.IsVerifyCodeEnable() {
		if r.CaptchaID == "" && r.CaptchaCode == "" {
			return xerrors.NewParamsErrors("请输入验证码")
		}

		if !captcha.VerifyString(ctx, r.CaptchaID, r.CaptchaCode) {
			return xerrors.NewParamsErrors("验证码错误")
		}
	}
	return nil
}

type LoginResponse struct {
	Token  string `json:"token" doc:"用户的token"`
	Notice string `json:"notice,omitempty" doc:"公告通知"`
}

type AuthController struct {
	BaseController
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
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
			card.ExpireTime = time.Now().Unix() + int64(constant.CardTimeExpire.Seconds())
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
	db.AddKFLog(card.CardID, "login", "登录成功")

	tk := jwt.New(jwt.SigningMethodHS256)
	tk.Claims = jwt.MapClaims{
		"cardId": card.CardID,
	}
	token, err := tk.SignedString([]byte(config.GetConfig().JwtKey))
	if err != nil {
		xlogger.Error(reqCtx, "签名token失败", xlogger.Err(err))
		c.Error(ctx, xerrors.NewCustomError("登录失败，请重试"))
		return
	}
	// 查找公告.
	c.Success(ctx, LoginResponse{
		Token:  token,
		Notice: caches.BillSettingCacheInstance.GetNotice(),
	})
}
