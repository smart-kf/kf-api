package kflog

import (
	"context"
	"sync"
	"time"

	xlogger "github.com/clearcodecn/log"
	"gorm.io/gorm"

	"github.com/smart-fm/kf-api/infrastructure/mysql"
	. "github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

var (
	_kfBackgroundTask *KFLogBackGroundTask
	kfTaskOnce        sync.Once
)

type KFLogBackGroundTask struct {
	buffer              chan *KFLog
	batchCreateDuration time.Duration // seconds
}

func InitKFLogBackgroundTask(duration time.Duration, bufferSize int) *KFLogBackGroundTask {
	kfTaskOnce.Do(
		func() {
			_kfBackgroundTask = &KFLogBackGroundTask{
				batchCreateDuration: duration,
				buffer:              make(chan *KFLog, bufferSize),
			}
		},
	)
	return _kfBackgroundTask
}

func AddKFLog(cardID string, handleFunc string, content string, ip string) {
	if _kfBackgroundTask == nil {
		return
	}
	fn, ok := Functions[handleFunc]
	if !ok {
		fn = handleFunc
	}

	log := &KFLog{
		Model:      gorm.Model{},
		CardID:     cardID,
		HandleFunc: fn,
		Content:    content,
		Ip:         ip,
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
	db := mysql.DB()
	if err := db.CreateInBatches(logs, len(logs)).Error; err != nil {
		xlogger.Error(context.Background(), "bill_log_task_create_failed", xlogger.Err(err))
	} else {
		xlogger.Info(context.Background(), "bill_log_task_create", xlogger.Any("logsLength", len(logs)))
	}
}
