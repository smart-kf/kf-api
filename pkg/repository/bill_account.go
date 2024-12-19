package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"std-api/pkg/db"
)

type BillAccountRepository struct {
	BaseRepository
}

func (r *BillAccountRepository) FindOneByUsername(ctx context.Context, username string) (*db.BillAccount, bool, error) {
	tx := r.getDB(ctx)

	var res db.BillAccount
	if err := tx.Where("username = ?", username).First(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &res, true, nil
}
