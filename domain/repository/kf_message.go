package repository

import (
	"context"
	"errors"
	xlogger "github.com/clearcodecn/log"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type KFMessageRepository struct{}

func (r *KFMessageRepository) SaveOne(ctx context.Context, chat *dao.KFMessage) error {
	tx := mysql.GetDBFromContext(ctx)
	res := tx.Model(&dao.KFMessage{}).Save(chat)
	if err := res.Error; err != nil {
		xlogger.Error(ctx, "SaveOne-failed", xlogger.Err(err))
		return err
	}
	return nil
}

type ListMsgOption struct {
	CardID        string
	ScrollRequest *common.ScrollRequest
}

func (r *KFMessageRepository) List(ctx context.Context, options *ListMsgOption) ([]*dao.KFMessage, error) {
	tx := mysql.GetDBFromContext(ctx).Debug()

	if len(options.CardID) == 0 {
		return nil, errors.New("cardID is required")
	}

	tx = tx.Where("card_id = ?", options.CardID)

	return common.Scroll[*dao.KFMessage](tx, options.ScrollRequest)
}

func (r *KFMessageRepository) ByIDs(ctx context.Context, cardID string, ids ...uint64) ([]*dao.KFMessage, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	tx := mysql.GetDBFromContext(ctx)

	if len(cardID) == 0 {
		return nil, errors.New("cardID is required")
	}

	tx = tx.Where("card_id = ?", cardID)
	tx = tx.Where("id in ?", ids)

	res := make([]*dao.KFMessage, 0)
	result := tx.
		Find(&res)
	if result.Error != nil {
		return nil, result.Error
	}

	return res, nil
}
