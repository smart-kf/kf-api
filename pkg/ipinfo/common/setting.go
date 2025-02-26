package common

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Setting struct {
	ApiKey     string
	Ip         string
	UserAgent  string
	HttpClient *http.Client
}

type IpInfo struct {
	IsChina         bool
	Lat             float64
	Lon             float64
	IsVpn           bool
	IsProxy         bool
	IsCloudProvider bool // 是否是机房.
	Country         string
	City            string
	Province        string
	UserAgentName   string
	UserAgentOs     string
	Net             string
}

func UnmarshalResponse(resp *http.Response, v interface{}) error {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
