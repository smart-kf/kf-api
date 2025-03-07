package httpclient

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

var (
	httpClient *resty.Client
)

func init() {
	httpClient = resty.New()
}

func SetClient(client *resty.Client) {
	httpClient = client
}

func doRequest(
	ctx context.Context,
	method string,
	url string,
	headers map[string]string,
	params interface{},
	out interface{},
) error {
	if headers == nil {
		headers = map[string]string{
			"Content-Type": "application/json",
		}
	}
	req := httpClient.R().SetContext(ctx).SetHeaders(headers)

	if strings.ToUpper(method) == http.MethodGet {
		queryParams := objectToMap(params)
		req.SetQueryParams(queryParams)
	} else {
		req.SetBody(params)
	}

	res, err := req.Execute(method, url)
	if err != nil {
		return err
	}
	return json.Unmarshal(res.Body(), out)
}

func objectToMap(obj interface{}) map[string]string {
	switch v := obj.(type) {
	case map[string]string:
		return v
	default:
		data, _ := json.Marshal(v)
		var result map[string]string
		json.Unmarshal(data, &result)
		return result
	}
}

func Get(ctx context.Context, url string, params map[string]string, out interface{}) error {
	return doRequest(ctx, "GET", url, nil, params, out)
}

func Post(ctx context.Context, url string, params interface{}, out interface{}) error {
	return doRequest(ctx, "POST", url, nil, params, out)
}

func PostHeader(ctx context.Context, url string, params interface{}, header map[string]string, out interface{}) error {
	return doRequest(ctx, "POST", url, header, params, out)
}
