package repository

import (
	"context"
	"errors"
	xlogger "github.com/clearcodecn/log"
	"github.com/smart-fm/kf-api/endpoints/common"
	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type KFUserRepository struct{}

func (r *KFUserRepository) SaveOne(ctx context.Context, chat *dao.KfUser) error {
	tx := mysql.GetDBFromContext(ctx)
	res := tx.Model(&dao.KfUser{}).Save(chat)
	if err := res.Error; err != nil {
		xlogger.Error(ctx, "SaveOne-failed", xlogger.Err(err))
		return err
	}
	return nil
}

type ListExtUserOption struct {
	CardID        string
	SearchBy      string
	ListType      kfbackend.ChatListType
	ScrollRequest *common.ScrollRequest
}

func (r *KFUserRepository) List(ctx context.Context, options *ListExtUserOption) ([]*dao.KfUser, error) {
	tx := mysql.GetDBFromContext(ctx).Debug()

	if len(options.CardID) == 0 {
		return nil, errors.New("cardID is required")
	}

	tx = tx.Where("card_id = ?", options.CardID)

	// 用户id/昵称/手机号/备注
	if options.SearchBy != "" {
		tx.Where(tx.Where("nick_name LIKE ?", "%"+options.SearchBy+"%").
			Or("id LIKE ?", "%"+options.SearchBy+"%").
			Or("mobile LIKE ?", "%"+options.SearchBy+"%").
			Or("remark_name LIKE ?", "%"+options.SearchBy+"%").
			Or("comments LIKE ?", "%"+options.SearchBy+"%"))
	}

	switch options.ListType {
	case kfbackend.ChatListTypeUnread:
		tx = tx.Where("unread_msg_cnt > 0") // 未读消息
	case kfbackend.ChatListTypeBlock:
		tx = tx.Where("block_at > 0") // 用拉黑时间来判断
	}

	return common.Scroll[*dao.KfUser](tx, options.ScrollRequest)
}

func (r *KFUserRepository) BatchUpdate(ctx context.Context, ids []uint, u dao.KfUser) error {
	tx := mysql.GetDBFromContext(ctx)
	res := tx.Model(&dao.KfUser{}).Where("id in ?", ids).Updates(u)
	if err := res.Error; err != nil {
		xlogger.Error(ctx, "BatchUpdate-failed", xlogger.Err(err))
		return err
	}
	return nil
}
