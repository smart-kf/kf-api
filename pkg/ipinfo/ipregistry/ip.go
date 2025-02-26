package ipregistry

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/smart-fm/kf-api/pkg/ipinfo/common"
)

type RegistryInfo struct {
	Ip       string `json:"ip"`
	Type     string `json:"type"`
	Hostname string `json:"hostname"`
	Carrier  struct {
		Name interface{} `json:"name"`
		Mcc  interface{} `json:"mcc"`
		Mnc  interface{} `json:"mnc"`
	} `json:"carrier"`
	Company struct {
		Domain string `json:"domain"`
		Name   string `json:"name"`
		Type   string `json:"type"`
	} `json:"company"`
	Connection struct {
		Asn          int    `json:"asn"`
		Domain       string `json:"domain"`
		Organization string `json:"organization"`
		Route        string `json:"route"`
		Type         string `json:"type"`
	} `json:"connection"`
	Currency struct {
		Code         string `json:"code"`
		Name         string `json:"name"`
		NameNative   string `json:"name_native"`
		Plural       string `json:"plural"`
		PluralNative string `json:"plural_native"`
		Symbol       string `json:"symbol"`
		SymbolNative string `json:"symbol_native"`
		Format       struct {
			DecimalSeparator string `json:"decimal_separator"`
			GroupSeparator   string `json:"group_separator"`
			Negative         struct {
				Prefix string `json:"prefix"`
				Suffix string `json:"suffix"`
			} `json:"negative"`
			Positive struct {
				Prefix string `json:"prefix"`
				Suffix string `json:"suffix"`
			} `json:"positive"`
		} `json:"format"`
	} `json:"currency"`
	Location struct {
		Continent struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"continent"`
		Country struct {
			Area              int           `json:"area"`
			Borders           []interface{} `json:"borders"`
			CallingCode       string        `json:"calling_code"`
			Capital           string        `json:"capital"`
			Code              string        `json:"code"`
			Name              string        `json:"name"`
			Population        int           `json:"population"`
			PopulationDensity float64       `json:"population_density"`
			Flag              struct {
				Emoji        string `json:"emoji"`
				EmojiUnicode string `json:"emoji_unicode"`
				Emojitwo     string `json:"emojitwo"`
				Noto         string `json:"noto"`
				Twemoji      string `json:"twemoji"`
				Wikimedia    string `json:"wikimedia"`
			} `json:"flag"`
			Languages []struct {
				Code   string `json:"code"`
				Name   string `json:"name"`
				Native string `json:"native"`
			} `json:"languages"`
			Tld string `json:"tld"`
		} `json:"country"`
		Region struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"region"`
		City      string  `json:"city"`
		Postal    string  `json:"postal"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Language  struct {
			Code   string `json:"code"`
			Name   string `json:"name"`
			Native string `json:"native"`
		} `json:"language"`
		InEu bool `json:"in_eu"`
	} `json:"location"`
	Security struct {
		IsAbuser        bool `json:"is_abuser"`
		IsAttacker      bool `json:"is_attacker"`
		IsBogon         bool `json:"is_bogon"`
		IsCloudProvider bool `json:"is_cloud_provider"`
		IsProxy         bool `json:"is_proxy"`
		IsRelay         bool `json:"is_relay"`
		IsTor           bool `json:"is_tor"`
		IsTorExit       bool `json:"is_tor_exit"`
		IsVpn           bool `json:"is_vpn"`
		IsAnonymous     bool `json:"is_anonymous"`
		IsThreat        bool `json:"is_threat"`
	} `json:"security"`
	TimeZone struct {
		Id               string    `json:"id"`
		Abbreviation     string    `json:"abbreviation"`
		CurrentTime      time.Time `json:"current_time"`
		Name             string    `json:"name"`
		Offset           int       `json:"offset"`
		InDaylightSaving bool      `json:"in_daylight_saving"`
	} `json:"time_zone"`
	UserAgent struct {
		Header       string `json:"header"`
		Name         string `json:"name"`
		Type         string `json:"type"`
		Version      string `json:"version"`
		VersionMajor string `json:"version_major"`
		Device       struct {
			Brand string `json:"brand"`
			Name  string `json:"name"`
			Type  string `json:"type"`
		} `json:"device"`
		Engine struct {
			Name         string `json:"name"`
			Type         string `json:"type"`
			Version      string `json:"version"`
			VersionMajor string `json:"version_major"`
		} `json:"engine"`
		Os struct {
			Name    string `json:"name"`
			Type    string `json:"type"`
			Version string `json:"version"`
		} `json:"os"`
	} `json:"user_agent"`
}

func Crawl(ctx context.Context, setting common.Setting) (common.IpInfo, error) {
	req, err := http.NewRequest(
		http.MethodGet, "https://api.ipregistry.co/"+setting.Ip+"?hostname=true&key="+setting.
			ApiKey, nil,
	)
	if err != nil {
		return common.IpInfo{}, err
	}
	req.Header.Set("origin", "https://ipregistry.co")
	req.Header.Set("referer", "https://ipregistry.co/")
	req.Header.Set("User-Agent", setting.UserAgent)

	resp, err := setting.HttpClient.Do(req)
	if err != nil {
		return common.IpInfo{}, err
	}
	defer resp.Body.Close()

	var registryInfo RegistryInfo
	if err := common.UnmarshalResponse(resp, &registryInfo); err != nil {
		return common.IpInfo{}, err
	}

	return common.IpInfo{
		IsChina:         registryInfo.Location.Country.Name == "China",
		Lat:             registryInfo.Location.Latitude,
		Lon:             registryInfo.Location.Longitude,
		IsVpn:           registryInfo.Security.IsVpn,
		IsProxy:         registryInfo.Security.IsProxy,
		IsCloudProvider: registryInfo.Security.IsCloudProvider,
		UserAgentName:   registryInfo.UserAgent.Name,
		UserAgentOs:     fmt.Sprintf("%s/%s", registryInfo.UserAgent.Os.Name, registryInfo.UserAgent.Os.Version),
	}, nil
}
