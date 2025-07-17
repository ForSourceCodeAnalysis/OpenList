package open123

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/OpenListTeam/OpenList/v4/drivers/base"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (d *Open123) Request(url, method string, callback base.ReqCallback, resp interface{}) ([]byte, error) {
	if d.ExpiredAt.Before(time.Now()) {
		if err := d.flushAccessToken(); err != nil {
			return nil, err
		}
	}

	req := base.RestyClient.R()
	req.SetHeaders(map[string]string{

		"platform":     "open_platform",
		"Content-Type": "application/json",
	})

	if callback != nil {
		callback(req)
	}
	if resp == nil {
		resp = &BaseResp{
			Code: -1,
		}
	}

	req.SetResult(resp)

	log.Debugf("API: %s", url)

	retryToken := true

	for i := range 3 {
		req.SetHeader("authorization", "Bearer "+d.AccessToken)

		res, err := req.Execute(method, url)
		if err != nil {
			return nil, err
		}
		if res.StatusCode() != http.StatusOK {
			return nil, fmt.Errorf("request failed, url: %s,statuscode: %d,message:%s", url, res.StatusCode(), res.String())
		}
		r := resp.(*BaseResp)

		switch r.Code {
		case 0:
			return res.Body(), nil
		case 401:
			if retryToken {
				retryToken = false
				if err := d.flushAccessToken(); err != nil {
					return nil, err
				}
			}
		case 429:
			time.Sleep(time.Second * time.Duration(i+1))
			log.Warningf("API: %s, 请求太频繁，对应API提示过多请减小QPS", url)
		}
	}
	return nil, fmt.Errorf("request failed, url: %s, max retry count excceed ", url)

}

func (d *Open123) flushAccessToken() error {
	// 第三方授权应用刷新token
	if d.RefreshToken != "" {
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
			Post(baseURL + refreshTokenAPI)

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
		Post(baseURL + accessTokenAPI)

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

func (d *Open123) SignURL(originURL, privateKey string, uid uint64, validDuration time.Duration) (newURL string, err error) {
	// 生成Unix时间戳
	ts := time.Now().Add(validDuration).Unix()

	// 生成随机数（建议使用UUID，不能包含中划线（-））
	rand := strings.ReplaceAll(uuid.New().String(), "-", "")

	// 解析URL
	objURL, err := url.Parse(originURL)
	if err != nil {
		return "", err
	}

	// 待签名字符串，格式：path-timestamp-rand-uid-privateKey
	unsignedStr := fmt.Sprintf("%s-%d-%s-%d-%s", objURL.Path, ts, rand, uid, privateKey)
	md5Hash := md5.Sum([]byte(unsignedStr))
	// 生成鉴权参数，格式：timestamp-rand-uid-md5hash
	authKey := fmt.Sprintf("%d-%s-%d-%x", ts, rand, uid, md5Hash)

	// 添加鉴权参数到URL查询参数
	v := objURL.Query()
	v.Add("auth_key", authKey)
	objURL.RawQuery = v.Encode()

	return objURL.String(), nil
}

func (d *Open123) getUserInfo(ctx context.Context) (*UserInfoResp, error) {
	var resp UserInfoResp

	if _, err := d.Request(UserInfo, http.MethodGet, func(req *resty.Request) {
		req.SetContext(ctx)
	}, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (d *Open123) getUID(ctx context.Context) (uint64, error) {
	if d.UID != 0 {
		return d.UID, nil
	}
	resp, err := d.getUserInfo(ctx)
	if err != nil {
		return 0, err
	}
	d.UID = resp.Data.UID
	return resp.Data.UID, nil
}

func (d *Open123) getFiles(parentFileId int64, limit int, lastFileId int64) (*FileListResp, error) {
	var resp FileListResp

	_, err := d.Request(baseURL+fileListAPI, http.MethodGet, func(req *resty.Request) {
		req.SetQueryParams(
			map[string]string{
				"parentFileId": strconv.FormatInt(parentFileId, 10),
				"limit":        strconv.Itoa(limit),
				"lastFileId":   strconv.FormatInt(lastFileId, 10),
				"trashed":      "false",
				"searchMode":   "",
				"searchData":   "",
			})
	}, resp)

	return resp.Data.(*FileListInfo), err
}

func (d *Open123) getDownloadInfo(fileId int64) (*DownloadInfo, error) {
	resp := BaseResp{
		Code: -1,
		Data: &DownloadInfo{},
	}

	_, err := d.Request(baseURL+downloadInfoAPI, http.MethodGet, func(req *resty.Request) {
		req.SetQueryParams(map[string]string{
			"fileId": strconv.FormatInt(fileId, 10),
		})
	}, &resp)

	return resp.Data.(*DownloadInfo), err
}

func (d *Open123) getDirectLink(fileId int64) (*DirectLinkResp, error) {
	var resp DirectLinkResp

	_, err := d.Request(DirectLink, http.MethodGet, func(req *resty.Request) {
		req.SetQueryParams(map[string]string{
			"fileID": strconv.FormatInt(fileId, 10),
		})
	}, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (d *Open123) mkdir(parentID int64, name string) error {
	_, err := d.Request(baseURL+mkdirAPI, http.MethodPost, func(req *resty.Request) {
		req.SetBody(base.Json{
			"parentID": strconv.FormatInt(parentID, 10),
			"name":     name,
		})
	}, nil)

	return err
}

func (d *Open123) move(fileID, toParentFileID int64) error {
	_, err := d.Request(baseURL+moveAPI, http.MethodPost, func(req *resty.Request) {
		req.SetBody(base.Json{
			"fileIDs":        []int64{fileID},
			"toParentFileID": toParentFileID,
		})
	}, nil)

	return err
}

func (d *Open123) rename(fileId int64, fileName string) error {
	_, err := d.Request(baseURL+renameAPI, http.MethodPut, func(req *resty.Request) {
		req.SetBody(base.Json{
			"fileId":   fileId,
			"fileName": fileName,
		})
	}, nil)

	return err
}

func (d *Open123) batchRename(renamelist []string) error {
	_, err := d.Request(baseURL+batchRenameAPI, http.MethodPost, func(req *resty.Request) {
		req.SetBody(base.Json{
			"renameList": renamelist,
		})
	}, nil)
	return err
}

func (d *Open123) trash(fileId int64) error {
	_, err := d.Request(baseURL+trashAPI, http.MethodPost, func(req *resty.Request) {
		req.SetBody(base.Json{
			"fileIDs": []int64{fileId},
		})
	}, nil)
	if err != nil {
		return err
	}

	return nil
}
