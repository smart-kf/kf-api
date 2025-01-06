package producer

import (
	"os"
	"time"

	"github.com/nsqio/go-nsq"

	"github.com/smart-fm/kf-api/config"
)

var NSQProducer *nsq.Producer

func InitProducer() {
	hostname, _ := os.Hostname()
	cfg := nsq.NewConfig()
	cfg.DialTimeout = time.Duration(config.GetConfig().NSQ.Timeout) * time.Second
	cfg.ReadTimeout = time.Duration(config.GetConfig().NSQ.Timeout) * time.Second
	cfg.WriteTimeout = time.Duration(config.GetConfig().NSQ.Timeout) * time.Second
	cfg.ClientID = hostname
	cfg.Hostname = hostname + "-kf-api"
	cfg.UserAgent = "go-" + hostname + "-kf-api"
	p, err := nsq.NewProducer(config.GetConfig().NSQ.Addrs[0], cfg)
	if err != nil {
		panic(err)
	}
	NSQProducer = p
}
