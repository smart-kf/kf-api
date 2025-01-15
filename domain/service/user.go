package service

import (
	"context"
	xlogger "github.com/clearcodecn/log"
	"github.com/smart-fm/kf-api/domain/event"
	"github.com/smart-fm/kf-api/domain/repository"
)

// UserService 访客service
type UserService struct{}

func init() {
	svc := UserService{}
	event.OnUserOffline(svc.OnUserOffline)
}

// OnUserOffline 访客下线记录下线时间
func (u *UserService) OnUserOffline(userId string) {
	repo := repository.KFUserRepository{}
	err := repo.Offline(context.TODO(), userId)
	if err != nil {
		xlogger.Error(context.TODO(), "UserService OnUserOffline fail", xlogger.Err(err))
	}
}
