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
	Type          *dao.MessageType
	From          *string
	FromType      *dao.ChatObjType
	To            *string
	ToType        *dao.ChatObjType
	ScrollRequest *common.ScrollRequest
}

func (r *KFMessageRepository) List(ctx context.Context, options *ListMsgOption) ([]*dao.KFMessage, error) {
	tx := mysql.GetDBFromContext(ctx).Debug()

	if len(options.CardID) == 0 {
		return nil, errors.New("cardID is required")
	}

	tx = tx.Where("card_id = ?", options.CardID)

	if options.From != nil {
		tx = tx.Where("from = ?", options.From)
	}

	if options.FromType != nil {
		tx = tx.Where("from_type = ?", options.FromType)
	}

	if options.To != nil {
		tx = tx.Where("to = ?", options.To)
	}

	if options.ToType != nil {
		tx = tx.Where("to_type = ?", options.ToType)
	}

	if options.Type != nil {
		tx = tx.Where("type = ?", options.Type)
	}

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

func (r *KFMessageRepository) BatchUpdateReadAt(ctx context.Context, ids []uint64, readAt int64) error {
	tx := mysql.GetDBFromContext(ctx)
	res := tx.Model(&dao.KFMessage{}).Where("id in ?", ids).Where("read_at <= 0").Updates(dao.KFMessage{ReadAt: readAt})
	if err := res.Error; err != nil {
		xlogger.Error(ctx, "BatchUpdateReadAt-failed", xlogger.Err(err))
		return err
	}
	return nil
}
