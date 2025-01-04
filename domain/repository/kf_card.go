package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	xlogger "github.com/clearcodecn/log"
	"gorm.io/gorm"

	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/pkg/common"
)

const (
	kfCardRedisKeyPrefix = "kf.card"
)

type KFCardRepository struct{}

func (r *KFCardRepository) CreateBatch(ctx context.Context, cards []*dao.KFCard) error {
	tx := mysql.GetDBFromContext(ctx)
	res := tx.Model(&dao.KFCard{}).CreateInBatches(cards, len(cards))
	if err := res.Error; err != nil {
		xlogger.Error(ctx, "CreateBatch-failed", xlogger.Err(err))
		return err
	}
	return nil
}

func (r *KFCardRepository) GetByID(ctx context.Context, id uint) (*dao.KFCard, bool, error) {
	tx := mysql.GetDBFromContext(ctx)
	var card dao.KFCard
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
func (r *KFCardRepository) UpdateOne(ctx context.Context, card *dao.KFCard) error {
	tx := mysql.GetDBFromContext(ctx)
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

func (r *KFCardRepository) List(ctx context.Context, options *ListCardOption) ([]*dao.KFCard, int64, error) {
	tx := mysql.GetDBFromContext(ctx)
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
	return common.Paginate[*dao.KFCard](tx, options.PageRequest)
}

func (r *KFCardRepository) FindByCardID(ctx context.Context, cardID string) (*dao.KFCard, bool, error) {
	cacheKey := fmt.Sprintf("%s.%s", kfCardRedisKeyPrefix, cardID)
	if card := getCacheByKey[dao.KFCard](ctx, cacheKey); card != nil {
		return card, true, nil
	}
	tx := mysql.GetDBFromContext(ctx)
	var card dao.KFCard
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
