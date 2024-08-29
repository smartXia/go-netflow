package rpc

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

type Autheticator interface {
	GetAuthHeaders() map[string]string
}

func NewAutheticator(appKey, appSecret string) Autheticator {
	if appKey == "" || appSecret == "" {
		panic("app key and secret MUST be provided")
	}
	return &autheticatorImpl{appKey: appKey, appSecret: appSecret}
}

type autheticatorImpl struct {
	appKey    string
	appSecret string
}

func (a *autheticatorImpl) GetAuthHeaders() map[string]string {
	now := time.Now().Unix()

	return map[string]string{
		"ak":        "gfjqXeKSKDIkkaIO9A7Pz5CV",
		"timestamp": strconv.FormatInt(now, 10),
		"sign":      getSign("A9QbSBGJXksjdopSDFaQ0Dnlr4ulx6", now),
	}
}

func getSign(appSecret string, now int64) string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s#%d", appSecret, now)))
	return hex.EncodeToString(h.Sum(nil))
}
