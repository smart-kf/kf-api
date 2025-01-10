package event

import (
	"context"

	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/pkg/safe"
)

/* 调用下面的注册方法 监听事件, 单个事件、多个钩子，并发回调. */

// OnUserOnline 注册事件
func OnUserOnline(fn UserOnlineEvent) {
	userOnlineEvents = append(userOnlineEvents, fn)
}

// OnUserOffline 注册事件
func OnUserOffline(fn UserOfflineEvent) {
	userOfflineEvents = append(userOfflineEvents, fn)
}

// OnKfBackendOnline 注册事件
func OnKfBackendOnline(fn KfBackendOnlineEvent) {
	kfBackendOnlineEvent = append(kfBackendOnlineEvent, fn)
}

// OnKfBackendOffline 注册事件
func OnKfBackendOffline(fn KfBackendOfflineEvent) {
	kfBackendOfflineEvent = append(kfBackendOfflineEvent, fn)
}

// UserOnlineEvent 客服前台用户上线事件
type UserOnlineEvent func(userId string)

// UserOfflineEvent 客服前台用户离线事件
type UserOfflineEvent func(userId string)

type KfBackendOnlineEvent func(cardId string)

type KfBackendOfflineEvent func(cardId string)

var userOnlineEvents []UserOnlineEvent

var userOfflineEvents []UserOfflineEvent

var kfBackendOnlineEvent []KfBackendOnlineEvent

var kfBackendOfflineEvent []KfBackendOfflineEvent

func TriggerEvent(ctx context.Context, event string, platform string, userId string, cardId string) {
	switch event {
	case constant.EventOnline:
		switch platform {
		case constant.PlatformKfFe:
			for _, e := range userOnlineEvents {
				fn := e
				safe.Go(
					func() {
						fn(userId)
					},
				)
			}

		case constant.PlatformKfBe:
			for _, e := range userOnlineEvents {
				fn := e
				safe.Go(
					func() {
						fn(cardId)
					},
				)
			}
		}
	case constant.EventOffline:
		switch platform {
		case constant.PlatformKfFe:
			for _, e := range userOnlineEvents {
				fn := e
				safe.Go(
					func() {
						fn(userId)
					},
				)
			}
		case constant.PlatformKfBe:
			for _, e := range userOnlineEvents {
				fn := e
				safe.Go(
					func() {
						fn(cardId)
					},
				)
			}
		}
	}
}
