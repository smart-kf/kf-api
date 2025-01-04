package billlog

import (
	"context"
	"sync"
	"time"

	xlogger "github.com/clearcodecn/log"
	"gorm.io/gorm"

	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/mysql/dao"
)

var (
	_billBackgroundTask *BillLogBackGroundTask
	billTaskOnce        sync.Once
)

type BillLogBackGroundTask struct {
	buffer              chan *dao.BillLog
	batchCreateDuration time.Duration // seconds
}

func InitBillLogBackgroundTask(duration time.Duration, bufferSize int) *BillLogBackGroundTask {
	billTaskOnce.Do(
		func() {
			_billBackgroundTask = &BillLogBackGroundTask{
				batchCreateDuration: duration,
				buffer:              make(chan *dao.BillLog, bufferSize),
			}
		},
	)
	return _billBackgroundTask
}

func AddBillLog(operator string, handleFunc string, content string) {
	if _billBackgroundTask == nil {
		return
	}
	log := &dao.BillLog{
		Model:      gorm.Model{},
		Operator:   operator,
		HandleFunc: handleFunc,
		Content:    content,
	}
	select {
	case _billBackgroundTask.buffer <- log:
	default:
		xlogger.Warn(context.Background(), "billLogTaskFull", xlogger.Any("log", log))
	}
}

func (b *BillLogBackGroundTask) Start(stopChan chan struct{}) {
	if b.batchCreateDuration == 0 {
		b.batchCreateDuration = 1 * time.Minute
	}
	tk := time.NewTicker(b.batchCreateDuration)
	defer tk.Stop()

	defer func() {
		close(b.buffer)
	}()

	var (
		buffers []*dao.BillLog
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

func (b *BillLogBackGroundTask) create(logs []*dao.BillLog) {
	db := mysql.DB()
	if err := db.CreateInBatches(logs, len(logs)).Error; err != nil {
		xlogger.Error(context.Background(), "bill_log_task_create_failed", xlogger.Err(err))
	} else {
		xlogger.Info(context.Background(), "bill_log_task_create", xlogger.Any("logsLength", len(logs)))
	}
}
