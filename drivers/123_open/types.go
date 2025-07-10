package open123

import (
	"strconv"
	"time"

	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
)

// BaseResp 通用返回字段
type BaseResp struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	XTraceID string `json:"x-traceID"`
	Data     any    `json:"data"`
}

// AccessTokenInfo access tokens
// 个人开发者获取到的token格式
type AccessTokenInfo struct {
	AccessToken string    `json:"accessToken"`
	ExpiredAt   time.Time `json:"expiredAt"` // 2025-03-23T15:48:37+08:00
}

// RefreshTokenInfo refresh token
// 第三方应用获取到的token格式
type RefreshTokenInfo struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"` //单位s，比如7200，表示在2h后过期
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

// -----------------------------------------------------------------------------------
var _ driver.IObjResp = (*Item)(nil)
var _ driver.IHashResp = (*Item)(nil)

// Item 文件信息
type Item struct {
	FileID       int64      `json:"fileId"`
	FileName     string     `json:"filename"`
	Type         int        `json:"type"`
	Size         int64      `json:"size"`
	Etag         string     `json:"etag"`
	Status       int        `json:"status"`
	ParentFileID int        `json:"parentFileId"`
	Category     int        `json:"category"`
	Trashed      int        `json:"trashed"`
	PunishFlag   int        `json:"punishFlag"`
	S3KeyFlag    string     `json:"s3KeyFlag"`
	StorageNode  string     `json:"storageNode"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdateAt     *time.Time `json:"updateAt"`
}

// GetID 获取文件ID
func (i *Item) GetID() string {
	return strconv.FormatInt(i.FileID, 10)
}

// GetPath 获取文件路径
func (i *Item) GetPath() string {
	return ""
}

// GetName 获取文件名
func (i *Item) GetName() string {
	return i.FileName
}

// GetSize 获取文件大小
func (i *Item) GetSize() int64 {
	return i.Size
}

// ModTime 获取文件修改时间
func (i *Item) ModTime() *time.Time {
	return i.UpdateAt
}

// CreateTime 获取文件创建时间
func (i *Item) CreateTime() *time.Time {
	return i.CreatedAt
}

// IsDir 判断是否是文件夹
func (i *Item) IsDir() bool {
	return i.Type == 1
}

// GetHash 获取文件hash
func (i *Item) GetHash() map[string]string {
	return map[string]string{
		driver.HashTypeMD5: i.Etag,
	}
}

// ----------------------------------------------------------------------------------
var _ driver.IListResp = (*ListObj)(nil)

// ListObj 列表对象
type ListObj struct {
	FileList   []*Item `json:"fileList"`
	LastFileID int64   `json:"lastFileId"`
}

// ListData 获取列表数据
func (l *ListObj) ListData() []driver.IObjResp {
	result := make([]driver.IObjResp, len(l.FileList))
	for i, item := range l.FileList {
		result[i] = item
	}
	return result
}

// Offset 获取偏移
func (l *ListObj) Offset() string {
	return strconv.FormatInt(l.LastFileID, 10)
}

// -------------------------------------------------------------------------

// APIInfo api信息
type APIInfo struct {
	url   string
	qps   int
	token chan struct{}
}

// Require 获取token
func (a *APIInfo) Require() {
	if a.qps > 0 {
		a.token <- struct{}{}
	}
}

// Release 释放token
func (a *APIInfo) Release() {
	if a.qps > 0 {
		time.AfterFunc(time.Second, func() {
			<-a.token
		})
	}
}

// SetQPS 设置qps
func (a *APIInfo) SetQPS(qps int) {
	a.qps = qps
	a.token = make(chan struct{}, qps)
}

// NowLen 获取当前token数量
func (a *APIInfo) NowLen() int {
	return len(a.token)
}

// initAPIInfo 初始化api信息
func initAPIInfo(url string, qps int) *APIInfo {
	return &APIInfo{
		url:   url,
		qps:   qps,
		token: make(chan struct{}, qps),
	}
}

// ---------------------------------------------------------------------

var _ driver.ILinkResp = (*DownloadInfo)(nil)

// DownloadInfo 下载信息
type DownloadInfo struct {
	DownloadURL string `json:"download_url"`
}

// GetURL 获取下载地址
func (di *DownloadInfo) GetURL() string {
	return di.DownloadURL
}

// -------------------------------------------------------------

var _ driver.IPreupResp = (*PreupInfo)(nil)

// var _ driver.IUploadServerResp = (*PreupInfo)(nil)

// PreupReq 预上传请求
type PreupReq struct {
	ParentFileID int64  `json:"parentFileId"`
	Filename     string `json:"filename"`
	Etag         string `json:"etag"`
	Size         int64  `json:"size"`
	Duplicate    int    `json:"duplicate"`
	ContainDir   bool   `json:"containDir"`
}

// PreupResp 预上传响应
type PreupResp struct {
	FileID      int64    `json:"fileID"`
	PreuploadID string   `json:"preuploadID"`
	Reuse       bool     `json:"reuse"`
	SliceSize   int64    `json:"sliceSize"`
	Servers     []string `json:"servers"`
}

// PreupInfo 预上传信息
type PreupInfo struct {
	*PreupResp
	BlockList []int
}

// GetPreupID 获取预上传ID
func (p *PreupInfo) GetPreupID() string {
	return p.PreuploadID
}

// GetSliceSize 获取分片大小
func (p *PreupInfo) GetSliceSize() int64 {
	return p.SliceSize
}

func (f File) CreateTime() time.Time {
	// 返回的时间没有时区信息，默认 UTC+8
	loc := time.FixedZone("UTC+8", 8*60*60)
	parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", f.CreateAt, loc)
	if err != nil {
		return time.Now()
	}
	return parsedTime
}

func (f File) ModTime() time.Time {
	// 返回的时间没有时区信息，默认 UTC+8
	loc := time.FixedZone("UTC+8", 8*60*60)
	parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", f.UpdateAt, loc)
	if err != nil {
		return time.Now()
	}
	return parsedTime
// IsRapidUpload 是否是秒传
func (p *PreupInfo) IsRapidUpload() bool {
	return p.Reuse
}

// GetBlockList 获取需要上传的分片列表
func (p *PreupInfo) GetBlockList() []int {
	return p.BlockList
}

// // GetUploadServer 获取上传服务器
// func (p *PreupInfo) GetUploadServer() string {
// 	return p.Servers[0]
// }

// SliceUploadCache 分片上传缓存信息
type SliceUploadCache struct {
	Filename          string   `json:"filename"`
	Size              int64    `json:"size"`
	Hash              string   `json:"hash"`
	UploadServer      string   `json:"upload_server"`
	PreupID           string   `json:"preup_id"`
	SliceSize         int64    `json:"slice_size"`
	UploadedBlockList []int    `json:"uploaded_block_list"` //已上传分片列表
	SliceHash         []string `json:"slice_hash"`          //分片hash
}

// type UploadCreateResp struct {
// 	BaseResp
// 	Data struct {
// 		FileID      int64  `json:"fileID"`
// 		PreuploadID string `json:"preuploadID"`
// 		Reuse       bool   `json:"reuse"`
// 		SliceSize   int64  `json:"sliceSize"`
// 	} `json:"data"`
// }

// type UploadUrlResp struct {
// 	BaseResp
// 	Data struct {
// 		PresignedURL string `json:"presignedURL"`
// 	}
// }

// type UploadCompleteResp struct {
// 	BaseResp
// 	Data struct {
// 		Async     bool  `json:"async"`
// 		Completed bool  `json:"completed"`
// 		FileID    int64 `json:"fileID"`
// 	} `json:"data"`
// }

// type UploadAsyncResp struct {
// 	BaseResp
// 	Data struct {
// 		Completed bool  `json:"completed"`
// 		FileID    int64 `json:"fileID"`
// 	} `json:"data"`
// }

// type UploadResp struct {
// 	BaseResp
// 	Data struct {
// 		AccessKeyId     string `json:"AccessKeyId"`
// 		Bucket          string `json:"Bucket"`
// 		Key             string `json:"Key"`
// 		SecretAccessKey string `json:"SecretAccessKey"`
// 		SessionToken    string `json:"SessionToken"`
// 		FileId          int64  `json:"FileId"`
// 		Reuse           bool   `json:"Reuse"`
// 		EndPoint        string `json:"EndPoint"`
// 		StorageNode     string `json:"StorageNode"`
// 		UploadId        string `json:"UploadId"`
// 	} `json:"data"`
// }

type UserInfoResp struct {
	BaseResp
	Data struct {
		UID uint64 `json:"uid"`
		// Username       string `json:"username"`
		// DisplayName    string `json:"displayName"`
		// HeadImage      string `json:"headImage"`
		// Passport       string `json:"passport"`
		// Mail           string `json:"mail"`
		SpaceUsed      uint64 `json:"spaceUsed"`
		SpacePermanent uint64 `json:"spacePermanent"`
		SpaceTemp      uint64 `json:"spaceTemp"`
		// SpaceTempExpr  int64  `json:"spaceTempExpr"`
		// Vip            bool   `json:"vip"`
		// DirectTraffic  int64  `json:"directTraffic"`
		// IsHideUID      bool   `json:"isHideUID"`
	} `json:"data"`
}

type FileListResp struct {
	BaseResp
	Data struct {
		LastFileId int64  `json:"lastFileId"`
		FileList   []File `json:"fileList"`
	} `json:"data"`
}

type DownloadInfoResp struct {
	BaseResp
	Data struct {
		DownloadUrl string `json:"downloadUrl"`
	} `json:"data"`
}

type DirectLinkResp struct {
	BaseResp
	Data struct {
		URL string `json:"url"`
	} `json:"data"`
}

// 创建文件V2返回
type UploadCreateResp struct {
	BaseResp
	Data struct {
		FileID      int64    `json:"fileID"`
		PreuploadID string   `json:"preuploadID"`
		Reuse       bool     `json:"reuse"`
		SliceSize   int64    `json:"sliceSize"`
		Servers     []string `json:"servers"`
	} `json:"data"`
}

// 上传完毕V2返回
type UploadCompleteResp struct {
	BaseResp
	Data struct {
		Completed bool  `json:"completed"`
		FileID    int64 `json:"fileID"`
	} `json:"data"`
}
// type UserInfoResp struct {
// 	BaseResp
// 	Data struct {
// 		UID            int64  `json:"uid"`
// 		Username       string `json:"username"`
// 		DisplayName    string `json:"displayName"`
// 		HeadImage      string `json:"headImage"`
// 		Passport       string `json:"passport"`
// 		Mail           string `json:"mail"`
// 		SpaceUsed      int64  `json:"spaceUsed"`
// 		SpacePermanent int64  `json:"spacePermanent"`
// 		SpaceTemp      int64  `json:"spaceTemp"`
// 		SpaceTempExpr  string `json:"spaceTempExpr"`
// 		Vip            bool   `json:"vip"`
// 		DirectTraffic  int64  `json:"directTraffic"`
// 		IsHideUID      bool   `json:"isHideUID"`
// 	} `json:"data"`
// }
