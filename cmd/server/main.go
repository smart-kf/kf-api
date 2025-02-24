package main

import (
	"flag"
	"log"
	http2 "net/http"
	"time"

	xlogger "github.com/clearcodecn/log"
	"golang.org/x/sync/errgroup"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/domain/caches"
	"github.com/smart-fm/kf-api/domain/repository"
	"github.com/smart-fm/kf-api/endpoints/cron/billlog"
	"github.com/smart-fm/kf-api/endpoints/cron/kflog"
	"github.com/smart-fm/kf-api/endpoints/http"
	"github.com/smart-fm/kf-api/endpoints/nsq/producer"
	"github.com/smart-fm/kf-api/infrastructure/mysql"
	"github.com/smart-fm/kf-api/infrastructure/nsq"
	"github.com/smart-fm/kf-api/infrastructure/redis"
	"github.com/smart-fm/kf-api/pkg/datagen/kffe"
)

var configName string
var generateFakeMessage bool

func init() {
	flag.StringVar(&configName, "c", "config.yaml", "配置文件")
	flag.BoolVar(&generateFakeMessage, "gen", false, "生成假消息数据")
}

func main() {
	flag.Parse()
	conf := config.Load(configName)
	initLogger(conf)
	mysql.Load()
	redis.InitRedis()
	nsq.InitNSQ()
	producer.InitProducer()

	if generateFakeMessage {
		doGenerateFakeMessage()
		return
	}
	caches.InitCacheInstances()

	go func() {
		setRepo := repository.BillSettingRepository{}
		setRepo.InitDefault()
	}()

	var (
		eg       errgroup.Group
		stopChan = make(chan struct{})
	)
	eg.Go(
		func() error {
			task := billlog.InitBillLogBackgroundTask(1*time.Minute, 100) // 1分钟清空buffer
			task.Start(stopChan)
			return nil
		},
	)
	eg.Go(
		func() error {
			task := kflog.InitKFLogBackgroundTask(1*time.Minute, 10000)
			task.Start(stopChan)
			return nil
		},
	)
	eg.Go(
		func() error {
			nsq.StartConsume(stopChan)
			return nil
		},
	)

	if err := http.Run(); err != nil {
		close(stopChan)
		return
	}

	log.Fatal(eg.Wait())
}

func initLogger(conf *config.Config) {
	// xlogger.AddHook(func(ctx context.Context) xlogger.Field {
	//	reqid, ok := ctx.Value("reqid").(string)
	//	if !ok {
	//		return xlogger.Field{}
	//	}
	//	return xlogger.Any("reqid", reqid)
	// })
	logger, err := xlogger.NewLog(
		xlogger.Config{
			Level:  conf.Log.Level,
			Format: conf.Log.Format,
			File:   conf.Log.File,
		},
	)

	if err != nil {
		panic(err)
	}

	xlogger.SetGlobal(logger)
}

func doGenerateFakeMessage() {
	cli := http2.Client{
		Timeout: 1 * time.Second,
	}
	_, err := cli.Get("https://www.google.com")
	if err != nil {
		// 自动探测.
		kffe.Gen(
			kffe.GenRequest{
				Host:   "http://localhost:8081",
				QRCode: "/s/EkCLyM/BJfmus/ak8BXI.html",
				CardId: `TM-J9pWlL8GfI`,
			},
		)
	} else {
		// 自动探测.
		kffe.Gen(
			kffe.GenRequest{
				Host:   "http://localhost:8081",
				QRCode: "/s/2VoVue/rfQj7G/oLqFXg.html",
				CardId: "TM-Tsmab1509q",
			},
		)
	}
}
