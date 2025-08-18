package open123

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/OpenListTeam/OpenList/v4/drivers/base"
	"github.com/go-resty/resty/v2"
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

func (d *Open123) getFiles(parentFileId int64, limit int, lastFileId int64) (*FileListInfo, error) {
	resp := &BaseResp{
		Code: -1,
		Data: &FileListInfo{},
	}

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

func (d *Open123) getDownloadInfo(fileID int64) (*DownloadInfo, error) {
	resp := BaseResp{
		Code: -1,
		Data: &DownloadInfo{},
	}

	_, err := d.Request(baseURL+downloadInfoAPI, http.MethodGet, func(req *resty.Request) {
		req.SetQueryParams(map[string]string{
			"fileId": strconv.FormatInt(fileID, 10),
		})
	}, &resp)

	return resp.Data.(*DownloadInfo), err
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

func (d *Open123) rename(fileID int64, fileName string) error {
	_, err := d.Request(baseURL+renameAPI, http.MethodPut, func(req *resty.Request) {
		req.SetBody(base.Json{
			"fileId":   fileID,
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

func (d *Open123) trash(fileID int64) error {
	_, err := d.Request(baseURL+trashAPI, http.MethodPost, func(req *resty.Request) {
		req.SetBody(base.Json{
			"fileIDs": []int64{fileID},
		})
	}, nil)

	return err
}

func (d *Open123) uploadSlice(req *UploadSliceReq) error {
	_, err := d.Request(req.Server+uploadSliceAPI, http.MethodPost, func(rt *resty.Request) {
		rt.SetHeader("Content-Type", "multipart/form-data")
		rt.SetMultipartFormData(map[string]string{
			"preuploadID": req.PreuploadID,
			"sliceMD5":    req.SliceMD5,
			"sliceNo":     strconv.FormatInt(int64(req.SliceNo), 10),
		})
		rt.SetMultipartField("slice", req.Name, "multipart/form-data", req.Slice)
	}, nil)
	return err
}

func (d *Open123) sliceUpComplete(uploadID string) error {
	r := &BaseResp{
		Code: -1,
		Data: &SingleUploadResp{},
	}
	_, err := d.Request(baseURL+uploadCompleteV2API, http.MethodPost, func(req *resty.Request) {
		req.SetBody(base.Json{
			"preuploadID": uploadID,
		})
	}, r)
	if err != nil {
		return err
	}
	rd := r.Data.(*SingleUploadResp)
	if rd.Completed {
		return nil
	}
	return errors.New("upload uncomplete")

}

func (d *Open123) getUploadServer() (string, error) {
	r := &BaseResp{
		Code: -1,
		Data: []string{},
	}
	_, err := d.Request(baseURL+uploadDomainAPI, "GET", nil, r)
	return r.Data.([]string)[0], err
}

func (d *Open123) singleUpload(req *SingleUploadReq) error {
	url, err := d.getUploadServer()
	if err != nil {
		return err
	}
	r := &BaseResp{
		Code: -1,
		Data: &SingleUploadResp{},
	}
	_, err = d.Request(url+uploadSingleCreateAPI, "POST", func(rt *resty.Request) {
		rt.SetHeader("Content-Type", "multipart/form-data")
		rt.SetMultipartFormData(map[string]string{
			"parentFileID": strconv.FormatInt(req.ParentFileID, 10),
			"filename":     req.FileName,
			"size":         strconv.FormatInt(req.Size, 10),
			"etag":         req.Etag,
			"duplicate":    strconv.Itoa(req.Duplicate),
		})
	}, r)
	if err != nil {
		return err
	}
	rd := r.Data.(*SingleUploadResp)
	if rd.Completed {
		return nil
	}
	return errors.New("upload uncomplete")
}
