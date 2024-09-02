package rpc

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type UrlProvider struct {
	ServerEndpoint string
	CommonHeaders  CommonHeadersProvider
}

func (p UrlProvider) Get(path string) (string, map[string]string) {
	return p.GetWithParams(path, nil)
}

func (p UrlProvider) GetWithParams(path string, params map[string]string) (string, map[string]string) {
	httpUrl := fmt.Sprintf("%s%s", p.ServerEndpoint, path)
	if !strings.HasPrefix(httpUrl, "http") {
		httpUrl = "https://" + httpUrl
	}

	headers := map[string]string{}
	for k, v := range p.CommonHeaders.GetCommonHeaders() {
		headers[k] = v
	}
	for k, v := range params {
		headers[k] = v
	}
	return httpUrl, headers
}

type CommonHeadersProvider struct {
	ImageVersion string
	DeviceId     string
	BizType      string
	Auth         Autheticator
}

func (p CommonHeadersProvider) GetCommonHeaders() map[string]string {
	//headers := p.Auth.GetAuthHeaders()
	now := time.Now().Unix()
	headers := map[string]string{
		"ak":        "gfjqXeKSKDIkkaIO9A7Pz5CV",
		"timestamp": strconv.FormatInt(now, 10),
		"sign":      getSign("A9QbSBGJXksjdopSDFaQ0Dnlr4ulx6", now),
	}

	headers["version"] = p.ImageVersion
	headers["bizType"] = p.BizType
	if headers["version"] == "" {
		headers["version"] = "unknown"
	}
	headers["agent"] = "v0.0.1"
	p.DeviceId = "c05803a7250ab9ccddb957122de312d0"
	if p.DeviceId != "" {
		headers["deviceID"] = p.DeviceId
	}
	return headers
}
