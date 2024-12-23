package repository

import (
	"context"
	"errors"
	"fmt"
	xlogger "github.com/clearcodecn/log"
	"gorm.io/gorm"
	"std-api/pkg/common"
	"std-api/pkg/constant"
	"std-api/pkg/db"
	"time"
)

const (
	kfCardRedisKeyPrefix = "kf.card"
)

type KFCardRepository struct{}

func (r *KFCardRepository) CreateBatch(ctx context.Context, cards []*db.KFCard) error {
	tx := db.GetDBFromContext(ctx)
	res := tx.Model(&db.KFCard{}).CreateInBatches(cards, len(cards))
	if err := res.Error; err != nil {
		xlogger.Error(ctx, "CreateBatch-failed", xlogger.Err(err))
		return err
	}
	return nil
}

func (r *KFCardRepository) GetByID(ctx context.Context, id uint) (*db.KFCard, bool, error) {
	tx := db.GetDBFromContext(ctx)
	var card db.KFCard
	res := tx.Where("id = ?", id).First(&card)
	if err := res.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &card, true, nil
}

// UpdateOne 更新一条数据
func (r *KFCardRepository) UpdateOne(ctx context.Context, card *db.KFCard) error {
	tx := db.GetDBFromContext(ctx)
	version := card.Version
	card.Version++
	res := tx.Where("id = ? and version = ?", card.ID, version).Updates(card)
	if res.Error != nil {
		xlogger.Error(ctx, "UpdateOne-failed", xlogger.Err(res.Error))
		return res.Error
	}
	deleteCacheByKey(ctx, fmt.Sprintf(kfCardRedisKeyPrefix, card.CardID))
	return nil
}

type ListCardOption struct {
	CardID             string
	CardType           constant.CardType
	SaleStatus         constant.SaleStatus
	LoginStatus        constant.LoginStatus
	ExpireStart        int64
	ExpireEnd          int64
	PageRequest        *common.PageRequest
	LastLoginTimeStart int64
	LastLoginTimeEnd   int64
}

func (r *KFCardRepository) List(ctx context.Context, options *ListCardOption) ([]*db.KFCard, int64, error) {
	tx := db.GetDBFromContext(ctx)
	if options.CardID != "" {
		tx = tx.Where("card_id = ?", options.CardID)
	}
	if options.CardType != 0 {
		tx = tx.Where("card_type = ?", options.CardType)
	}
	if options.SaleStatus != 0 {
		tx = tx.Where("sale_status = ?", options.SaleStatus)
	}
	if options.LoginStatus != 0 {
		tx = tx.Where("login_status = ?", options.LoginStatus)
	}
	if options.ExpireStart != 0 {
		tx = tx.Where("expire_time >= ?", options.ExpireStart)
	}
	if options.ExpireEnd != 0 {
		tx = tx.Where("expire_time <= ?", options.ExpireEnd)
	}
	if options.LastLoginTimeStart != 0 {
		tx = tx.Where("last_login_time >= ?", options.LastLoginTimeStart)
	}
	if options.LastLoginTimeEnd != 0 {
		tx = tx.Where("last_login_time <= ?", options.LastLoginTimeEnd)
	}
	return common.Paginate[*db.KFCard](tx, options.PageRequest)
}

func (r *KFCardRepository) FindByCardID(ctx context.Context, cardID string) (*db.KFCard, bool, error) {
	cacheKey := fmt.Sprintf("%s.%s", kfCardRedisKeyPrefix, cardID)
	if card := getCacheByKey[db.KFCard](ctx, cacheKey); card != nil {
		return card, true, nil
	}
	tx := db.GetDBFromContext(ctx)
	var card db.KFCard
	res := tx.Where("card_id = ?", cardID).First(&card)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	setCacheByKey(ctx, cacheKey, &card, 10*time.Minute)
	return &card, true, nil
}
