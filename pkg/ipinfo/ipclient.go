package ipinfo

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/smart-fm/kf-api/config"
	"github.com/smart-fm/kf-api/pkg/ipinfo/common"
	"github.com/smart-fm/kf-api/pkg/ipinfo/ipregistry"
)

var (
	ipClient *IpClient
)

type IpClient struct {
	httpClient *http.Client
}

func (c *IpClient) init() {
	if c.httpClient != nil {
		return
	}
	c.httpClient = &http.Client{}
	proxy := config.GetConfig().Ip2Region.Proxy
	if proxy != "" {
		u, _ := url.Parse(proxy)
		if u != nil {
			c.httpClient.Transport = &http.Transport{
				Proxy: http.ProxyURL(u),
			}
		}
	}
	ti := 5
	if config.GetConfig().Ip2Region.Timeout != 0 {
		ti = config.GetConfig().Ip2Region.Timeout
	}
	c.httpClient.Timeout = time.Second * time.Duration(ti)
}

func (c *IpClient) Crawl(ctx context.Context, ua string, ip string) (common.IpInfo, error) {

	c.init()

	setting := common.Setting{
		ApiKey:     config.GetConfig().Ip2Region.RegistryApiKey,
		Ip:         ip,
		UserAgent:  ua,
		HttpClient: c.httpClient,
	}

	ipinfo, err := ipregistry.Crawl(ctx, setting)
	if err != nil {
		return common.IpInfo{}, err
	}

	loc := GetLocation(config.GetConfig().Ip2Region.XDBPath, ip)

	ipinfo.Country = loc.Country
	ipinfo.City = loc.City
	ipinfo.Net = loc.Net
	ipinfo.Province = loc.Province
	return ipinfo, nil
}

func Crawl(ctx context.Context, ua string, ip string) (common.IpInfo, error) {
	if ipClient == nil {
		ipClient = &IpClient{}
	}

	return ipClient.Crawl(ctx, ua, ip)
}
