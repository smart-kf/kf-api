package repository

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/samber/lo"

	"github.com/smart-fm/kf-api/endpoints/common/constant"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
	"github.com/smart-fm/kf-api/pkg/xerrors"
)

type KfWelcomeMessageRepository struct{}

type CopyParams struct {
	FromCardId           string
	ToCardId             string
	ReplaceTargetContent string // 被替换
	ReplaceContent       string // 替换
	MsgType              string
}

func (w *KfWelcomeMessageRepository) CopyFromCard(ctx context.Context, param CopyParams) error {
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

func (w *KfWelcomeMessageRepository) UpsertOne(ctx context.Context, msg *dao.KfWelcomeMessage) error {
	db := mysql.GetDBFromContext(ctx)
	if msg.ID > 0 {
		var exist dao.KfWelcomeMessage
		err := db.Where("id = ? and msg_type = ? and deleted_at is null", msg.ID, msg.MsgType).First(&exist).Error
		if err != nil {
			return err
		}
		exist.Content = msg.Content
		exist.Type = msg.Type
		exist.Enable = msg.Enable
		exist.Sort = msg.Sort
		exist.Title = msg.Title
		err = db.Where("id = ? and msg_type = ?", msg.ID, msg.MsgType).Save(exist).Error
		if err != nil {
			return err
		}
	}
	return db.Create(msg).Error
}

func (w *KfWelcomeMessageRepository) List(
	ctx context.Context,
	cardId string,
	msgType string,
	page int64,
	pageSize int64,
) ([]*dao.KfWelcomeMessage, int64, error) {
	var (
		data []*dao.KfWelcomeMessage
		cnt  int64
	)
	tx := mysql.GetDBFromContext(ctx)
	tx = tx.Where("card_id = ? and msg_type = ?", cardId, msgType).Order("sort asc")
	tx = tx.Model(&dao.KfWelcomeMessage{}).Count(&cnt)
	if page != 0 && pageSize != 0 {
		tx = tx.Limit(int(pageSize)).Offset(int((page - 1) * pageSize))
	}
	err := tx.Find(&data).Error
	if err != nil {
		return nil, 0, err
	}

	return data, cnt, nil
}

func (w *KfWelcomeMessageRepository) Delete(ctx context.Context, cardId string, id int) (*dao.KfWelcomeMessage, error) {
	db := mysql.GetDBFromContext(ctx)

	var msg dao.KfWelcomeMessage
	err := db.Where("id = ? and card_id = ? and deleted_at is null").First(&msg).Error
	if err != nil {
		return nil, err
	}
	n := db.Where("id = ? and card_id = ?", id, cardId).Delete(&dao.KfWelcomeMessage{}).RowsAffected
	if n == 0 {
		return nil, err
	}

	return &msg, nil
}

// 同步处理 全文索引.
func (w *KfWelcomeMessageRepository) upsertIndex(ctx context.Context, msg *dao.KfWelcomeMessage) {
	if msg.MsgType != constant.SmartReply {
		return
	}
}
