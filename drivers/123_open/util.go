package open123

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/OpenListTeam/OpenList/v4/drivers/base"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// Request 发送api请求
func (d *Open123) Request(apiInfo *APIInfo, method string, callback base.ReqCallback, resp any) ([]byte, error) {
	// 检查token是否过期
	if d.IsAccessTokenExpired() {
		if err := d.flushAccessToken(); err != nil {
			return nil, err
		}
	}

	req := base.RestyClient.R().
		SetHeader("Platform", "open_platform")

	if callback != nil {
		callback(req)
	}
	if resp != nil {
		req.SetResult(resp)
	}

	log.Debugf("API: %s, QPS: %d, NowLen: %d", apiInfo.url, apiInfo.qps, apiInfo.NowLen())

	apiInfo.Require()
	defer apiInfo.Release()

	// 最多重试2次，共3次
	for range 3 {

		req.SetHeader("authorization", "Bearer "+d.AccessToken)
		res, err := req.Execute(method, apiInfo.url)
		if err != nil { // 内部错误，直接返回（RestyClient已经设置了重试次数）
			return nil, err
		}
		// 这个是必须的，避免出现虽然没发生内部错误，但是服务器出现问题
		if res.StatusCode() != http.StatusOK {
			return nil, fmt.Errorf("status code: %d, body: %s", res.StatusCode(), res.String())
		}

		body := res.Body()

		// 解析为通用响应
		var baseResp BaseResp
		if err = json.Unmarshal(body, &baseResp); err != nil {
			return nil, err
		}
		r := resp.(*BaseResp)

		switch baseResp.Code {
		case 0:
			return body, nil
		case 401: //token失效
			if err = d.flushAccessToken(); err != nil {
				return nil, err
			}
		case 429: // api请求过快
			time.Sleep(time.Second)
			log.Warningf("API: %s, QPS: %d, 请求太频繁，对应API提示过多请减小QPS", apiInfo.url, apiInfo.qps)
		default: //其他错误
			return nil, errors.New(baseResp.Message)
		}
	}
	return nil, fmt.Errorf("max retry count exceeded,api : %s", apiInfo.url)

}

func (d *Open123) flushAccessToken() error {
	// 第三方授权应用刷新token
	if d.RefreshToken != "" {
		rta := d.qpsInstance[refreshTokenAPI]
		rta.Require()
		defer rta.Release()

		r := &BaseResp{
			Code: -1,
			Data: &RefreshTokenInfo{},
		}
		res, err := base.RestyClient.R().
			SetHeaders(map[string]string{
				"Platform":     "open_platform",
				"Content-Type": "application/json",
			}).
			SetResult(r).
			SetQueryParams(map[string]string{
				"grant_type":    "refresh_token",
				"client_id":     d.ClientID,
				"client_secret": d.ClientSecret,
				"refresh_token": d.RefreshToken,
			}).
			Post(rta.url)
		if err != nil {
			return err
		}
		if res.StatusCode() != http.StatusOK {
			return fmt.Errorf("refresh token failed: %s", res.String())
		}

		if r.Code != 0 {
			return fmt.Errorf("refresh token failed: %s", r.Message)
		}
		rt := r.Data.(*RefreshTokenInfo)
		d.RefreshToken = rt.RefreshToken
		d.AccessToken = rt.AccessToken
		d.ExpiredAt = time.Now().Add(time.Duration(rt.ExpiresIn) * time.Second)
		return nil
	}
	// 个人开发者获取access token
	at := d.qpsInstance[accessTokenAPI]

	at.Require()
	defer at.Release()
	r := &BaseResp{
		Code: -1,
		Data: &AccessTokenInfo{},
	}
	res, err := base.RestyClient.R().
		SetHeaders(map[string]string{
			"Platform":     "open_platform",
			"Content-Type": "application/json",
		}).
		SetResult(r).
		SetQueryParams(map[string]string{
			"client_id":     d.ClientID,
			"client_secret": d.ClientSecret,
		}).
		Post(at.url)
	if err != nil {
		return err
	}
	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("refresh token failed: %s", res.String())
	}

	if r.Code != 0 {
		return fmt.Errorf("refresh token failed: %s", r.Message)
	}
	ac := r.Data.(*AccessTokenInfo)
	d.AccessToken = ac.AccessToken
	d.ExpiredAt = ac.ExpiredAt

	return nil
}

func (d *Open123) getFiles(parentFileID int64, limit int, lastFileID int64) (*ListObj, error) {
	resp := &BaseResp{
		Code: -1,
		Data: &ListObj{},
	}

	_, err := d.Request(d.qpsInstance[fileListAPI], http.MethodGet, func(req *resty.Request) {
		req.SetQueryParams(
			map[string]string{
				"parentFileId": strconv.FormatInt(parentFileID, 10),
				"limit":        strconv.Itoa(limit),
				"lastFileId":   strconv.FormatInt(lastFileID, 10),
				"trashed":      "false",
				"searchMode":   "",
				"searchData":   "",
			}).
			SetHeader("Content-Type", "application/json")
	}, resp)
	r := resp.Data.(*ListObj)

	return r, err
}

func (d *Open123) getDownloadInfo(fileID int64) (*DownloadInfo, error) {
	resp := &BaseResp{
		Code: -1,
		Data: &DownloadInfo{
			DownloadURL: "",
		},
	}

	_, err := d.Request(d.qpsInstance[downloadInfoAPI], http.MethodGet, func(req *resty.Request) {
		req.SetQueryParams(map[string]string{
			"fileId": strconv.FormatInt(fileID, 10),
		}).SetHeader("Content-Type", "application/json")
	}, &resp)
	r := resp.Data.(*DownloadInfo)
	return r, err
}

func (d *Open123) mkdir(parentID int64, name string) error {
	_, err := d.Request(d.qpsInstance[mkdirAPI], http.MethodPost, func(req *resty.Request) {
		req.SetBody(map[string]any{
			"parentID": parentID,
			"name":     name,
		}).SetHeader("Content-Type", "application/json")
	}, nil)

	return err
}

func (d *Open123) move(fileIDs []int64, toParentFileID int64) error {
	_, err := d.Request(d.qpsInstance[moveAPI], http.MethodPost, func(req *resty.Request) {
		req.SetBody(map[string]any{
			"fileIDs":        fileIDs,
			"toParentFileID": toParentFileID,
		}).SetHeader("Content-Type", "application/json")
	}, nil)

	return err
}

func (d *Open123) rename(renameList []string) error {
	_, err := d.Request(d.qpsInstance[renameAPI], http.MethodPost, func(req *resty.Request) {
		req.SetBody(map[string]any{
			"renameList": renameList,
		}).SetHeader("Content-Type", "application/json")
	}, nil)

	return err
}

func (d *Open123) trash(fileIDs []int64) error {
	_, err := d.Request(d.qpsInstance[trashAPI], http.MethodPost, func(req *resty.Request) {
		req.SetBody(map[string]any{
			"fileIDs": fileIDs,
		}).SetHeader("Content-Type", "application/json")
	}, nil)

	return err
}
