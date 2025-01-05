package repository

import (
	"context"
	xlogger "github.com/clearcodecn/log"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type KFExternalUserRepository struct{}

func (r *KFExternalUserRepository) SaveOne(ctx context.Context, chat *dao.KFExternalUser) error {
	tx := mysql.GetDBFromContext(ctx)
	res := tx.Model(&dao.KFExternalUser{}).Save(chat)
	if err := res.Error; err != nil {
		xlogger.Error(ctx, "SaveOne-failed", xlogger.Err(err))
		return err
	}
	return nil
}

type ListExtUserOption struct {
	CardID        string
	SearchBy      string
	ListType      int
	ScrollRequest *common.ScrollRequest
}

func (r *KFExternalUserRepository) List(ctx context.Context, options *ListExtUserOption) ([]*dao.KFExternalUser, int64, error) {
	tx := mysql.GetDBFromContext(ctx)
	if options.CardID != "" {
		tx = tx.Where("card_id = ?", options.CardID)
	}
	return common.Scroll[*dao.KFExternalUser](tx, options.ScrollRequest)
}
