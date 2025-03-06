package caches

import (
	"context"

	"github.com/samber/lo"
)

type UnReadHelper struct{}

func (r *UnReadHelper) GetUnReadUserIDs(ctx context.Context, cardId string) (map[string]int64, []string, error) {
	uid2Cnt, err := UserUnReadCacheInstance.GetUnReadUsers(ctx, cardId)
	if err != nil {
		return nil, nil, err
	}
	uids := lo.MapToSlice(
		uid2Cnt, func(k string, v int64) string {
			return k
		},
	)

	return uid2Cnt, uids, nil
}

func (r *UnReadHelper) GetUserExtras(ctx context.Context, cardId string) (map[string]int64, []string, error) {
	uid2Cnt, err := UserUnReadCacheInstance.GetUnReadUsers(ctx, cardId)
	if err != nil {
		return nil, nil, err
	}
	uids := lo.MapToSlice(
		uid2Cnt, func(k string, v int64) string {
			return k
		},
	)

	return uid2Cnt, uids, nil
}
