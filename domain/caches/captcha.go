package caches

import (
	"context"
	"fmt"
	"time"

	xlogger "github.com/clearcodecn/log"
	"github.com/redis/go-redis/v9"

	"github.com/smart-fm/kf-api/pkg/utils"
)

const (
	captchaDuration = 5 * time.Minute
)

type CaptchaCache struct {
	db         *redis.Client
	keyPrefix  string
	expireTime time.Duration
}

func NewCaptchaCache(client *redis.Client) *CaptchaCache {
	return &CaptchaCache{
		db:         client,
		keyPrefix:  "kf.captcha.",
		expireTime: captchaDuration,
	}
}
func (c *CaptchaCache) key(id string) string {
	return fmt.Sprintf("%s.%s", c.keyPrefix, id)
}

func (c *CaptchaCache) Set(ctx context.Context, id string, digits []byte) {
	e := newEntity(id, digits)

	err := c.db.Set(ctx, c.key(id), utils.MustMarshalObject(e), c.expireTime).Err()
	if err != nil {
		xlogger.Error(ctx, "captchaCache set failed", xlogger.Err(err))
	}
}

func (c *CaptchaCache) Del(ctx context.Context, id string) {
	err := c.db.Del(ctx, c.key(id)).Err()
	if err != nil {
		xlogger.Error(ctx, "captchaCache del failed", xlogger.Err(err))
	}
}

func (c *CaptchaCache) Get(ctx context.Context, id string) []byte {
	res := c.db.Get(ctx, c.key(id))
	if res.Err() != nil {
		return nil
	}
	var e entity
	utils.MustUnMarshalObject([]byte(res.String()), &e)
	return []byte(e.CaptchaCode)
}

type entity struct {
	CaptchaId   string `json:"captchaId"`
	CaptchaCode string `json:"captchaCode"`
}

func newEntity(id string, code []byte) *entity {
	return &entity{
		CaptchaId:   id,
		CaptchaCode: string(code),
	}
}
