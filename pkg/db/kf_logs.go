package db

import (
	"context"
	xlogger "github.com/clearcodecn/log"
	"gorm.io/gorm"
	"sync"
	"time"
)

// KFLog 客服后台审计日志.
type KFLog struct {
	gorm.Model
	CardID     string `json:"card_id" gorm:"column:card_id"`
	HandleFunc string `json:"handle_func" gorm:"column:handle_func"` // 操作类型
	Content    string `json:"content" gorm:"column:content;text"`    // 操作内容
}

func (KFLog) TableName() string {
	return "kf_log"
}

var (
	_kfBackgroundTask *KFLogBackGroundTask
	kfTaskOnce        sync.Once
)

type KFLogBackGroundTask struct {
	buffer              chan *KFLog
	batchCreateDuration time.Duration // seconds
}

func InitKFLogBackgroundTask(duration time.Duration, bufferSize int) *KFLogBackGroundTask {
	kfTaskOnce.Do(func() {
		_kfBackgroundTask = &KFLogBackGroundTask{
			batchCreateDuration: duration,
			buffer:              make(chan *KFLog, bufferSize),
		}
	})
	return _kfBackgroundTask
}

func AddKFLog(cardID string, handleFunc string, content string) {
	if _kfBackgroundTask == nil {
		return
	}
	log := &KFLog{
		Model:      gorm.Model{},
		CardID:     cardID,
		HandleFunc: handleFunc,
		Content:    content,
	}
	select {
	case _kfBackgroundTask.buffer <- log:
	default:
		xlogger.Warn(context.Background(), "KFLogTaskFull", xlogger.Any("log", log))
	}
}

func (b *KFLogBackGroundTask) Start(stopChan chan struct{}) {
	if b.batchCreateDuration == 0 {
		b.batchCreateDuration = 1 * time.Minute
	}
	tk := time.NewTicker(b.batchCreateDuration)
	defer tk.Stop()

	defer func() {
		close(b.buffer)
	}()

	var (
		buffers []*KFLog
	)
	for {
		select {
		case buf, ok := <-b.buffer:
			if !ok {
				return
			}
			buffers = append(buffers, buf)
		case <-tk.C:
			if len(buffers) == 0 {
				continue
			}
			b.create(buffers)
			buffers = buffers[:0]
		case <-stopChan:
			return
		}
	}
}

func (b *KFLogBackGroundTask) create(logs []*KFLog) {
	db := DB()
	if err := db.CreateInBatches(logs, len(logs)).Error; err != nil {
		xlogger.Error(context.Background(), "bill_log_task_create_failed", xlogger.Err(err))
	} else {
		xlogger.Info(context.Background(), "bill_log_task_create", xlogger.Any("logsLength", len(logs)))
	}
}
