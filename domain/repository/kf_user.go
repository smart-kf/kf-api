package repository

import (
	"context"
	"time"

	xlogger "github.com/clearcodecn/log"

	"gorm.io/gorm"

	"github.com/smart-fm/kf-api/endpoints/http/vo/kfbackend"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type KFUserRepository struct{}

func (r *KFUserRepository) SaveOne(ctx context.Context, kfUser *dao.KfUser) error {
	tx := mysql.GetDBFromContext(ctx)
	res := tx.Model(&dao.KfUser{})
	if kfUser.ID > 0 {
		res = res.Where("id = ?", kfUser.ID)
	}
	res = res.Save(kfUser)
	if err := res.Error; err != nil {
		xlogger.Error(ctx, "SaveOne-failed", xlogger.Err(err))
		return err
	}
	return nil
}

func (r *KFUserRepository) FindByToken(ctx context.Context, cardId string, token string) (
	*dao.KfUser, bool, error,
) {
	var (
		kfUser dao.KfUser
		err    error
	)
	tx := mysql.GetDBFromContext(ctx)
	err = tx.Model(&dao.KfUser{}).Where("card_id = ? and uuid = ?", cardId, token).First(&kfUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		xlogger.Error(ctx, "FindByToken-failed", xlogger.Err(err))
		return nil, false, nil
	}
	return &kfUser, true, nil
}

type ListUserOption struct {
	CardID   string
	SearchBy string
	UUids    []string
	Blocked  bool
	ListType kfbackend.ChatListType
	Page     int
	PageSize int
}

func (r *KFUserRepository) List(ctx context.Context, options *ListUserOption) ([]*dao.KfUser, error) {
	tx := mysql.GetDBFromContext(ctx)
	tx = tx.Where("card_id = ?", options.CardID)

	// 用户id/昵称/手机号/备注
	if options.SearchBy != "" {
		searchBy := "%" + options.SearchBy + "%"
		tx = tx.Where("remark_name LIKE ? or mobile like ?", searchBy, searchBy)
	}

	if len(options.UUids) != 0 {
		tx = tx.Where("uuid in ?", options.UUids) // 有未读消息的访客
	}
	if options.Blocked {
		tx = tx.Where("block_at > 0") // 用拉黑时间来判断
	}
	var (
		res []*dao.KfUser
	)

	tx = tx.Order("top_at desc,block_at asc, last_chat_at desc")
	err := tx.Limit(options.PageSize).Offset((options.Page - 1) * options.PageSize).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
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

type UpdateColOption struct {
	UUID          string
	LastChatAt    int64
	LastMessageId string
	TopAt         int64
	BlockAt       int64
	ScanCount     int64
}

func (r *KFUserRepository) UpdateCol(ctx context.Context, cardID string, opt UpdateColOption) error {
	tx := mysql.GetDBFromContext(ctx)
	tx = tx.Model(&dao.KfUser{}).Where("card_id = ? and uuid = ?", cardID, opt.UUID)
	var updateCols = make(map[string]interface{})
	if opt.LastChatAt > 0 {
		updateCols["last_chat_at"] = opt.LastChatAt
	}
	if opt.LastMessageId != "" {
		updateCols["last_message_id"] = opt.LastMessageId
	}
	if opt.TopAt > 0 {
		updateCols["top_at"] = opt.TopAt
	}
	if opt.BlockAt > 0 {
		updateCols["block_at"] = opt.BlockAt
	}
	if opt.ScanCount > 0 {
		updateCols["scan_count"] = opt.ScanCount
	}

	if len(updateCols) > 0 {
		return tx.Updates(updateCols).Error
	}

	return nil
}
