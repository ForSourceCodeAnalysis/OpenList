package base

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/internal/net"
	"github.com/go-resty/resty/v2"
)

var (
	// NoRedirectClient is a resty client with no redirect
	NoRedirectClient *resty.Client
	// RestyClient is a resty client
	RestyClient *resty.Client
	// HTTPClient is a http client
	HttpClient *http.Client
)
var UserAgent = "Mozilla/5.0 (Macintosh; Apple macOS 15_5) AppleWebKit/537.36 (KHTML, like Gecko) Safari/537.36 Chrome/138.0.0.0"
var DefaultTimeout = time.Second * 30

// InitClient initializes the client
func InitClient() {
	NoRedirectClient = resty.New().
		SetHeader("user-agent", UserAgent).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: conf.Conf.TLSInsecureSkipVerify}).
		SetRedirectPolicy(
			resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}),
		)

	RestyClient = NewRestyClient()

	HttpClient = net.NewHttpClient()
}

// NewRestyClient returns a new resty client
func NewRestyClient() *resty.Client {
	client := resty.New().
		SetHeader("user-agent", UserAgent).
		SetRetryCount(3).
		SetRetryResetReaders(true).
		SetTimeout(DefaultTimeout).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: conf.Conf.TLSInsecureSkipVerify})
	return client
}

/// example
/// StatusCode 在200-299之间，且body格式为json/xml，会转换为result
/// StatusCode > 399，且body格式为json/xml，会转换为error
/// 其他形式，不会发生转换，可以通过res.String()获取
/// err和e的主要区别就是，err是请求没有成功（没有到达服务器的底层错误返回），e是请求成功，服务端返回的错误
/// 当出现下面这种情况： 服务请求成功，但是
///  StatusCode > 399，且body格式不为json/xml，此时result没有结果，error也没有结果
/// 如果仅根据 error和result判断结果就是不准确的，必须要判断 StatusCode
// res,err := base.RestyClient.R().

// 		SetResult(&resp).
// 		//
// 		SetError(&e).
// 		SetQueryParams(map[string]string{
// 			"grant_type":    "refresh_token",
// 			"refresh_token": d.RefreshToken,
// 			"client_id":     d.ClientID,
// 			"client_secret": d.ClientSecret,
// 		}).
// 		Get(u)
