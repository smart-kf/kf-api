package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/smart-fm/kf-api/config"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var (
	_logicClient *logicClient
	o            sync.Once
)

type logicClient struct {
	httpClient *resty.Client
}

func GetLogicClient() *logicClient {
	o.Do(func() {
		conf := config.GetConfig().HttpClient
		var hc = &http.Client{}

		if conf.Timeout != 0 {
			hc.Timeout = time.Duration(conf.Timeout) * time.Second
		}
		if conf.Proxy != "" {
			u, _ := url.Parse(conf.Proxy)
			hc.Transport = &http.Transport{
				Proxy: http.ProxyURL(u),
			}
		}
		r := resty.NewWithClient(hc)
		_logicClient = &logicClient{
			httpClient: r,
		}
	})
	return _logicClient
}

type resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (c *logicClient) PushMids(ctx context.Context, op int32, mids []int64, msg interface{}) error {
	req := c.httpClient.R()
	req.SetContext(ctx)

	q := make(url.Values)
	for _, mid := range mids {
		q.Add("mids", strconv.FormatInt(mid, 10))
	}
	q.Add("operation", strconv.FormatInt(int64(op), 10))

	logicAddr := config.GetConfig().HttpClient.LogicAddress
	rsp, err := req.Post(fmt.Sprintf("%s/goim/push/mids?%s", logicAddr, q.Encode()))
	if err != nil {
		return err
	}
	var res resp
	err = json.Unmarshal(rsp.Body(), &res)
	if err != nil {
		return err
	}
	if res.Code == 0 {
		return nil
	}
	return errors.New("push failed: " + res.Message)
}
