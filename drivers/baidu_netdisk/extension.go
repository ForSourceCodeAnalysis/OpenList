package baidu_netdisk

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/OpenListTeam/OpenList/v4/drivers/base"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/model/reqres"
	"github.com/OpenListTeam/OpenList/v4/internal/model/tables"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

func (d *BaiduNetdisk) GetUploadInfo() *model.UploadInfo {
	return &model.UploadInfo{
		SliceHashNeed:    true,
		HashMd5Need:      true,
		HashMd5256KBNeed: true,
	}
}

func (d *BaiduNetdisk) BatchRename(ctx context.Context, obj model.Obj, renameObjs []model.RenameObj) error {
	data := []base.Json{}
	for _, ro := range renameObjs {
		data = append(data, base.Json{
			"path":    filepath.Join(obj.GetPath(), ro.SrcName),
			"newname": ro.NewName,
		})
	}

	_, err := d.manage("rename", data)
	return err
}

func (d *BaiduNetdisk) BatchRemove(ctx context.Context, srcObj model.Obj, objs []model.IDName) error {
	data := []string{}
	for _, obj := range objs {
		data = append(data, filepath.Join(srcObj.GetPath(), obj.Name))
	}

	_, err := d.manage("delete", data)
	return err
}

// SliceUpload 上传分片
func (d *BaiduNetdisk) SliceUpload(c context.Context, req *tables.SliceUpload, sliceno uint, fd io.Reader) error {
	fp := filepath.Join(req.DstPath, req.Name)
	if sliceno == 0 { //第一个分片需要先执行预上传
		rtype := 1
		if req.Overwrite {
			rtype = 3
		}
		precreateResp, err := d.precreate(&PrecreateReq{
			Path:       fp,
			Size:       req.Size,
			Isdir:      0,
			BlockList:  strings.Split(req.SliceHash, ","),
			Autoinit:   1,
			Rtype:      rtype,
			ContentMd5: req.HashMd5,
			SliceMd5:   req.HashMd5256KB,
		})
		if err != nil {
			return err
		}
		req.PreupID = precreateResp.Uploadid
	}
	err := d.uploadSlice(c, map[string]string{
		"method":       "upload",
		"access_token": d.AccessToken,
		"type":         "tmpfile",
		"path":         fp,
		"uploadid":     req.PreupID,
		"partseq":      strconv.Itoa(int(sliceno)),
	}, req.Name, fd)
	return err

}

// Preup 预上传(自定以接口，为了适配自定义的分片上传)
func (d *BaiduNetdisk) Preup(ctx context.Context, srcobj model.Obj, req *reqres.PreupReq) (*model.PreupInfo, error) {
	return &model.PreupInfo{
		SliceSize: d.getSliceSize(req.Size),
	}, nil
}

// UploadSliceComplete 分片上传完成
func (d *BaiduNetdisk) UploadSliceComplete(ctx context.Context, su *tables.SliceUpload) error {
	fp := filepath.Join(su.DstPath, su.Name)
	rsp := &SliceUpCompleteResp{}
	t := time.Now().Unix()
	sh, err := json.Marshal(strings.Split(su.SliceHash, ","))
	if err != nil {
		return err
	}
	b, err := d.create(fp, int64(su.Size), 0, su.PreupID, string(sh), rsp, t, t)
	if err != nil {
		log.Error(err, rsp, string(b))
	}
	return err
}

func (d *BaiduNetdisk) precreate(req *PrecreateReq) (*PrecreateResp, error) {
	bl, err := json.Marshal(req.BlockList)
	if err != nil {
		log.Errorf("json.Marshal error: %v", err)
		return nil, err
	}
	b := map[string]string{
		"path":        req.Path,
		"size":        strconv.Itoa(int(req.Size)),
		"isdir":       strconv.Itoa(req.Isdir),
		"autoinit":    strconv.Itoa(req.Autoinit),
		"rtype":       strconv.Itoa(req.Rtype),
		"block_list":  string(bl),
		"content-md5": req.ContentMd5,
		"slice-md5":   req.SliceMd5,
	}

	res := &PrecreateResp{}
	r, err := d.request("https://pan.baidu.com/rest/2.0/xpan/file", http.MethodPost, func(rt *resty.Request) {
		rt.SetQueryParam("method", "precreate").
			SetFormData(b)

	}, res)
	if err != nil {
		log.Errorf("baidu_netdisk precreate error: %s, %v", string(r), err)
		return nil, err
	}
	return res, nil

}

func (d *BaiduNetdisk) locateUpload(req *LocateUploadReq) (string, error) {
	res := &LocateUploadResp{}
	_, err := d.get("/pcs/file", map[string]string{
		"method":         "locateupload",
		"appid":          "250528",
		"path":           req.Path,
		"uploadid":       req.UploadID,
		"upload_version": "2.0",
	}, res)
	if err != nil {
		return "", err
	}
	return res.Servers[0].Server, nil

}

// 请求参数
// 参数名称	类型	是否必填	示例	参数位置	描述
// method	string	是	precreate	URL参数	本接口固定为precreate
// access_token	string	是	12.a6b7dbd428f731035f771b8d15063f61.86400.1292922000-2346678-124328	URL参数	接口鉴权认证参数，标识用户
// path	string	是	/apps/appName/filename.jpg	RequestBody参数	上传后使用的文件绝对路径，需要urlencode
// size	int	是	4096	RequestBody参数	文件和目录两种情况：上传文件时，表示文件的大小，单位B；上传目录时，表示目录的大小，目录的话大小默认为0
// isdir	int	是	0	RequestBody参数	是否为目录，0 文件，1 目录
// block_list	string	是	["98d02a0f54781a93e354b1fc85caf488", "ca5273571daefb8ea01a42bfa5d02220"]	RequestBody参数	文件各分片MD5数组的json串。block_list的含义如下，如果上传的文件小于4MB，其md5值（32位小写）即为block_list字符串数组的唯一元素；如果上传的文件大于4MB，需要将上传的文件按照4MB大小在本地切分成分片，不足4MB的分片自动成为最后一个分片，所有分片的md5值（32位小写）组成的字符串数组即为block_list。
// autoinit	int	是	1	RequestBody参数	固定值1
// rtype	int	否	1	RequestBody参数	文件命名策略。
// 1 表示当path冲突时，进行重命名
// 2 表示当path冲突且block_list不同时，进行重命名
// 3 当云端存在同名文件时，对该文件进行覆盖
// uploadid	string	否	P1-MTAuMjI4LjQzLjMxOjE1OTU4NTg==	RequestBody参数	上传ID
// content-md5	string	否	b20f8ac80063505f264e5f6fc187e69a	RequestBody参数	文件MD5，32位小写
// slice-md5	string	否	9aa0aa691s5c0257c5ab04dd7eddaa47	RequestBody参数	文件校验段的MD5，32位小写，校验段对应文件前256KB
// local_ctime	string	否	1595919297	RequestBody参数	客户端创建时间， 默认为当前时间戳
// local_mtime	string	否	1595919297	RequestBody参数	客户端修改时间，默认为当前时间戳
type PrecreateReq struct {
	Path       string   `json:"path"`                  // 上传后使用的文件绝对路径（需urlencode）
	Size       int64    `json:"size"`                  // 文件或目录大小，单位B
	Isdir      int      `json:"isdir"`                 // 是否为目录，0 文件，1 目录
	BlockList  []string `json:"block_list"`            // 文件各分片MD5数组的json串
	Autoinit   int      `json:"autoinit"`              // 固定值1
	Rtype      int      `json:"rtype,omitempty"`       // 文件命名策略，非必填
	Uploadid   string   `json:"uploadid,omitempty"`    // 上传ID，非必填
	ContentMd5 string   `json:"content-md5,omitempty"` // 文件MD5，非必填
	SliceMd5   string   `json:"slice-md5,omitempty"`   // 文件校验段的MD5，非必填
	LocalCtime string   `json:"local_ctime,omitempty"` // 客户端创建时间，非必填
	LocalMtime string   `json:"local_mtime,omitempty"` // 客户端修改时间，非必填
}

// path	string	是	/apps/appName/filename.jpg	RequestBody参数	上传后使用的文件绝对路径，需要urlencode，需要与预上传precreate接口中的path保持一致
// size	string	是	4096	RequestBody参数	文件或目录的大小，必须要和文件真实大小保持一致，需要与预上传precreate接口中的size保持一致
// isdir	string	是	0	RequestBody参数	是否目录，0 文件、1 目录，需要与预上传precreate接口中的isdir保持一致
// block_list	json array	是	["7d57c40c9fdb4e4a32d533bee1a4e409"]	RequestBody参数	文件各分片md5数组的json串
// 需要与预上传precreate接口中的block_list保持一致，同时对应分片上传superfile2接口返回的md5，且要按照序号顺序排列，组成md5数组的json串。
// uploadid	string	是	N1-MjIwLjE4MS4zfgsdgewrSEEd=	RequestBody参数	预上传precreate接口下发的uploadid
// rtype	int	否	1	RequestBody参数	文件命名策略，默认0
// 0 为不重命名，返回冲突
// 1 为只要path冲突即重命名
// 2 为path冲突且block_list不同才重命名
// 3 为覆盖，需要与预上传precreate接口中的rtype保持一致
// local_ctime	int	否	1596009229	RequestBody参数	客户端创建时间(精确到秒)，默认为当前时间戳
// local_mtime	int	否	1596009229	RequestBody参数	客户端修改时间(精确到秒)，默认为当前时间戳
// zip_quality	int	否	70	RequestBody参数	图片压缩程度，有效值50、70、100（带此参数时，zip_sign 参数需要一并带上）
// zip_sign	int	否	7d57c40c9fdb4e4a32d533bee1a4e409	RequestBody参数	未压缩原始图片文件真实md5（带此参数时，zip_quality 参数需要一并带上）
// is_revision	int	否	0	RequestBody参数	是否需要多版本支持
// 1为支持，0为不支持， 默认为0 (带此参数会忽略重命名策略)
// mode	int	否	1	RequestBody参数	上传方式
// 1 手动、2 批量上传、3 文件自动备份
// 4 相册自动备份、5 视频自动备份
// exif_info	string	否	{"height":3024,"date_time_original":"2018:09:06 15:58:58","model":"iPhone 6s","width":4032,"date_time_digitized":"2018:09:06 15:58:58","date_time":"2018:09:06 15:58:58","orientation":6,"recovery":0}	RequestBody参数	json字符串，orientation、width、height、recovery为必传字段，其他字段如果没有可以不传
type SliceUpCompleteReq struct {
	Path       string   `json:"path"`                  // 上传后使用的文件绝对路径（需urlencode），与预上传precreate接口中的path保持一致
	Size       int64    `json:"size"`                  // 文件或目录的大小，必须与实际大小一致
	Isdir      int      `json:"isdir"`                 // 是否目录，0 文件、1 目录，与预上传precreate接口中的isdir保持一致
	BlockList  []string `json:"block_list"`            // 文件各分片md5数组的json串，与预上传precreate接口中的block_list保持一致
	Uploadid   string   `json:"uploadid"`              // 预上传precreate接口下发的uploadid
	Rtype      int      `json:"rtype,omitempty"`       // 文件命名策略，默认0
	LocalCtime int64    `json:"local_ctime,omitempty"` // 客户端创建时间(精确到秒)，默认为当前时间戳
	LocalMtime int64    `json:"local_mtime,omitempty"` // 客户端修改时间(精确到秒)，默认为当前时间戳
	ZipQuality int      `json:"zip_quality,omitempty"` // 图片压缩程度，有效值50、70、100（带此参数时，zip_sign 参数需要一并带上）
	ZipSign    string   `json:"zip_sign,omitempty"`    // 未压缩原始图片文件真实md5（带此参数时，zip_quality 参数需要一并带上）
	IsRevision int      `json:"is_revision,omitempty"` // 是否需要多版本支持，1为支持，0为不支持，默认为0
	Mode       int      `json:"mode,omitempty"`        // 上传方式，1手动、2批量上传、3文件自动备份、4相册自动备份、5视频自动备份
	ExifInfo   string   `json:"exif_info,omitempty"`   // exif信息，json字符串，orientation、width、height、recovery为必传字段
}

// errno	int	错误码
// fs_id	uint64	文件在云端的唯一标识ID
// md5	string	文件的MD5，只有提交文件时才返回，提交目录时没有该值
// server_filename	string	文件名
// category	int	分类类型, 1 视频 2 音频 3 图片 4 文档 5 应用 6 其他 7 种子
// path	string	上传后使用的文件绝对路径
// size	uint64	文件大小，单位B
// ctime	uint64	文件创建时间
// mtime	uint64	文件修改时间
// isdir	int	是否目录，0 文件、1 目录
type SliceUpCompleteResp struct {
	Errno          int    `json:"errno"`           // 错误码
	FsID           uint64 `json:"fs_id"`           // 文件在云端的唯一标识ID
	Md5            string `json:"md5,omitempty"`   // 文件的MD5，只有提交文件时才返回，提交目录时没有该值
	ServerFilename string `json:"server_filename"` // 文件名
	Category       int    `json:"category"`        // 分类类型, 1 视频 2 音频 3 图片 4 文档 5 应用 6 其他 7 种子
	Path           string `json:"path"`            // 上传后使用的文件绝对路径
	Size           uint64 `json:"size"`            // 文件大小，单位B
	Ctime          uint64 `json:"ctime"`           // 文件创建时间
	Mtime          uint64 `json:"mtime"`           // 文件修改时间
	Isdir          int    `json:"isdir"`           // 是否目录，0 文件、1 目录
}

// LocateUploadReq 获取上传域名请求
// method	string	是	locateupload	URL参数	本接口固定为locateupload
// appid	Integer	是	250528	URL参数	应用ID，本接口固定为250528
// access_token	string	是	12.a6b7dbd428f731035f771b8d15063f61.86400.1292922000-2346678-124328	URL参数	接口鉴权认证参数，标识用户
// path	string	是	/apps/appName/filename.jpg	URL参数	上传后使用的文件绝对路径，需要urlencode
// uploadid	string	是	P1-MTAuMjI4LjQzLjMxOjE1OTU4NTg==	URL参数	上传ID
// upload_version	string	是	2.0	URL参数	版本号，本接口固定为2.0
type LocateUploadReq struct {
	AppID         int    `json:"appid"`          // 应用ID，固定为250528
	Path          string `json:"path"`           // 上传后使用的文件绝对路径，需要urlencode
	UploadID      string `json:"uploadid"`       // 上传ID
	UploadVersion string `json:"upload_version"` // 版本号，固定为2.0
}

//	{
//		"bak_server": [],
//		"bak_servers": [
//			{
//				"server": "https://c.pcs.baidu.com"
//			}
//		],
//		"client_ip": "39.103.0.0",
//		"error_code": 0,
//		"error_msg": "",
//		"expire": 60,
//		"host": "c.pcs.baidu.com",
//		"newno": "",
//		"quic_server": [],
//		"quic_servers": [
//			{
//				"server": "https://panup.pcs.baidu.com"
//			}
//		],
//		"request_id": 2962549359693023232,
//		"server": [],
//		"server_time": 1715074756,
//		"servers": [
//			{
//				"server": "https://c3.pcs.baidu.com"
//			},
//			{
//				"server": "http://c3.pcs.baidu.com"
//			},
//			{
//				"server": "http://c3.pcs.baidu.com"
//			},
//			{
//				"server": "http://c2.pcs.baidu.com"
//			}
//		],
//		"sl": 0
//	}
type ServerInfo struct {
	Server string `json:"server"`
}

type LocateUploadResp struct {
	BakServer   []any        `json:"bak_server"`   // 兼容空数组
	BakServers  []ServerInfo `json:"bak_servers"`  // 备用服务器列表
	ClientIP    string       `json:"client_ip"`    // 客户端IP
	ErrorCode   int          `json:"error_code"`   // 错误码
	ErrorMsg    string       `json:"error_msg"`    // 错误信息
	Expire      int          `json:"expire"`       // 过期时间（秒）
	Host        string       `json:"host"`         // 主机
	Newno       string       `json:"newno"`        // 新编号，可能为空
	QuicServer  []any        `json:"quic_server"`  // 兼容空数组
	QuicServers []ServerInfo `json:"quic_servers"` // quic服务器列表
	RequestID   int64        `json:"request_id"`   // 请求ID
	Server      []any        `json:"server"`       // 兼容空数组
	ServerTime  int64        `json:"server_time"`  // 服务器时间
	Servers     []ServerInfo `json:"servers"`      // 服务器列表
	Sl          int          `json:"sl"`           // 未知字段
}
