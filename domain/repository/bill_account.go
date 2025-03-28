package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

type BillAccountRepository struct {
	BaseRepository
}

func (r *BillAccountRepository) FindOneByUsername(ctx context.Context, username string) (
	*dao.BillAccount,
	bool,
	error,
) {
	tx := r.getDB(ctx)

	var res dao.BillAccount
	if err := tx.Where("username = ?", username).First(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &res, true, nil
}

func (r *BillAccountRepository) Save(ctx context.Context, account *dao.BillAccount) error {
	tx := r.getDB(ctx)
	if err := tx.Where("id = ?", account.ID).Save(account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return nil
	}
	return nil
}
