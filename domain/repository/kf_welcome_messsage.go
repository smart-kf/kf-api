package repository

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/samber/lo"

	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

type KfWelcomeMessage struct{}

type CopyParams struct {
	FromCardId           string
	ToCardId             string
	ReplaceTargetContent string // 被替换
	ReplaceContent       string // 替换
	MsgType              string
}

func (w *KfWelcomeMessage) CopyFromCard(ctx context.Context, param CopyParams) error {
	db := mysql.GetDBFromContext(ctx)
	var content []*dao.KfWelcomeMessage
	err := db.Where(
		"card_id = ? and msg_type = ? and deleted_at is null", param.FromCardId,
		param.MsgType,
	).Find(&content).Error
	if err != nil {
		return err
	}
	var (
		newContent []*dao.KfWelcomeMessage
	)
	lo.ForEach(
		content, func(item *dao.KfWelcomeMessage, index int) {
			if param.ReplaceContent != "" && param.ReplaceTargetContent != "" {
				item.Content = strings.ReplaceAll(item.Content, param.ReplaceTargetContent, param.ReplaceContent)
			}
			if item.Keyword != "" && param.ReplaceContent != "" && param.ReplaceTargetContent != "" {
				item.Keyword = strings.ReplaceAll(item.Keyword, param.ReplaceTargetContent, param.ReplaceContent)
			}
			if item.Title != "" && param.ReplaceContent != "" && param.ReplaceTargetContent != "" {
				item.Title = strings.ReplaceAll(item.Title, param.ReplaceTargetContent, param.ReplaceContent)
			}
			if utf8.RuneCountInString(item.Content) > 255 || utf8.RuneCountInString(item.Keyword) > 255 || utf8.
				RuneCountInString(item.Keyword) > 255 {
				err = xerrors.NewCustomError("替换后内容长度需要<255个字符，请重试")
			}
			if err != nil {
				return
			}
			item.CardId = param.ToCardId
			item.ID = 0
			item.CreatedAt = time.Time{}
			item.UpdatedAt = time.Time{}
			newContent = append(newContent, item)
		},
	)
	if err != nil {
		return err
	}
	if len(newContent) > 0 {
		return db.CreateInBatches(newContent, len(newContent)).Error
	}
	return nil
}
