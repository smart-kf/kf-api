package repository

import (
	"context"
	"errors"
	"time"

	xlogger "github.com/clearcodecn/log"
	"gorm.io/gorm"

	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type KFUserRepository struct{}

func (r *KFUserRepository) SaveOne(ctx context.Context, kfUser *dao.KfUser) error {
	tx := mysql.GetDBFromContext(ctx)
	res := tx.Model(&dao.KfUser{}).Save(kfUser)
	if err := res.Error; err != nil {
		xlogger.Error(ctx, "SaveOne-failed", xlogger.Err(err))
		return err
	}
	return nil
}

func (r *KFUserRepository) FindByToken(ctx context.Context, token string) (kfUser *dao.KfUser, ok bool, err error) {
	tx := mysql.GetDBFromContext(ctx)
	res := tx.Model(&dao.KfUser{}).Where("uuid = ?", token).First(&kfUser)
	if err = res.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		xlogger.Error(ctx, "FindByToken-failed", xlogger.Err(err))
		return
	}
	return
}

type ListUserOption struct {
	CardID      string
	SearchBy    string
	UnreadUUIDs []string
	ListType    kfbackend.ChatListType

	ScrollRequest *common.ScrollRequest
}

func (r *KFUserRepository) List(ctx context.Context, options *ListUserOption) ([]*dao.KfUser, error) {
	tx := mysql.GetDBFromContext(ctx).Debug()

	if len(options.CardID) == 0 {
		return nil, errors.New("cardID is required")
	}

	tx = tx.Where("card_id = ?", options.CardID)

	// 用户id/昵称/手机号/备注
	if options.SearchBy != "" {
		tx.Where(
			tx.Where("nick_name LIKE ?", "%"+options.SearchBy+"%").
				Or("id LIKE ?", "%"+options.SearchBy+"%").
				Or("mobile LIKE ?", "%"+options.SearchBy+"%").
				Or("remark_name LIKE ?", "%"+options.SearchBy+"%").
				Or("comments LIKE ?", "%"+options.SearchBy+"%"),
		)
	}

	switch options.ListType {
	case kfbackend.ChatListTypeUnread:
		tx = tx.Where("uuid in ?", options.UnreadUUIDs) // 有未读消息的访客
	case kfbackend.ChatListTypeBlock:
		tx = tx.Where("block_at > 0") // 用拉黑时间来判断
	}

	return common.Scroll[*dao.KfUser](tx, options.ScrollRequest)
}

func (r *KFUserRepository) BatchUpdate(ctx context.Context, ids []string, u dao.KfUser) error {
	tx := mysql.GetDBFromContext(ctx)
	res := tx.Model(&dao.KfUser{}).Where("uuid in ?", ids).Updates(u)
	if err := res.Error; err != nil {
		xlogger.Error(ctx, "BatchUpdate-failed", xlogger.Err(err))
		return err
	}
	return nil
}

func (r *KFUserRepository) Offline(ctx context.Context, id string) error {
	tx := mysql.GetDBFromContext(ctx)
	res := tx.Model(&dao.KfUser{}).Where("uuid = ?", id).Update("offline_at", time.Now().Unix())
	if err := res.Error; err != nil {
		xlogger.Error(ctx, "Offline-failed", xlogger.Err(err))
		return err
	}
	return nil
}
