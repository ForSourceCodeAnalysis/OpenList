package base

import "github.com/go-resty/resty/v2"

// Json map[string]any 别名
type Json map[string]any

// TokenResp 获取token的返回
type TokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ReqCallback 请求回调
type ReqCallback func(req *resty.Request)
