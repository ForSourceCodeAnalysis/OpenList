package baidu_netdisk

import (
	"path"
	"strconv"
	"time"

	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
)

type TokenErrResp struct {
	ErrorDescription string `json:"error_description"`
	Error            string `json:"error"`
}

type File struct {
	//TkbindId     int    `json:"tkbind_id"`
	//OwnerType    int    `json:"owner_type"`
	Category int `json:"category"`
	//RealCategory string `json:"real_category"`
	FsId int64 `json:"fs_id"`
	//OperId      int   `json:"oper_id"`
	Thumbs struct {
		//Icon string `json:"icon"`
		Url3 string `json:"url3"`
		//Url2 string `json:"url2"`
		//Url1 string `json:"url1"`
	} `json:"thumbs"`
	//Wpfile         int    `json:"wpfile"`

	Size int64 `json:"size"`
	//ExtentTinyint7 int    `json:"extent_tinyint7"`
	Path string `json:"path"`
	//Share          int    `json:"share"`
	//Pl             int    `json:"pl"`
	ServerFilename string `json:"server_filename"`
	Md5            string `json:"md5"`
	//OwnerId        int    `json:"owner_id"`
	//Unlist int `json:"unlist"`
	Isdir int `json:"isdir"`

	// list resp
	ServerCtime int64 `json:"server_ctime"`
	ServerMtime int64 `json:"server_mtime"`
	LocalMtime  int64 `json:"local_mtime"`
	LocalCtime  int64 `json:"local_ctime"`
	//ServerAtime    int64    `json:"server_atime"` `

	// only create and precreate resp
	Ctime int64 `json:"ctime"`
	Mtime int64 `json:"mtime"`
}

func fileToObj(f File) *model.ObjThumb {
	if f.ServerFilename == "" {
		f.ServerFilename = path.Base(f.Path)
	}
	if f.ServerCtime == 0 {
		f.ServerCtime = f.Ctime
	}
	if f.ServerMtime == 0 {
		f.ServerMtime = f.Mtime
	}
	return &model.ObjThumb{
		Object: model.Object{
			ID:       strconv.FormatInt(f.FsId, 10),
			Path:     f.Path,
			Name:     f.ServerFilename,
			Size:     f.Size,
			Modified: time.Unix(f.ServerMtime, 0),
			Ctime:    time.Unix(f.ServerCtime, 0),
			IsFolder: f.Isdir == 1,

			// 直接获取的MD5是错误的
			HashInfo: utils.NewHashInfo(utils.MD5, DecryptMd5(f.Md5)),
		},
		Thumbnail: model.Thumbnail{Thumbnail: f.Thumbs.Url3},
	}
}

type ListResp struct {
	Errno     int    `json:"errno"`
	GuidInfo  string `json:"guid_info"`
	List      []File `json:"list"`
	RequestId int64  `json:"request_id"`
	Guid      int    `json:"guid"`
}

type DownloadResp struct {
	Errmsg string `json:"errmsg"`
	Errno  int    `json:"errno"`
	List   []struct {
		//Category    int    `json:"category"`
		//DateTaken   int    `json:"date_taken,omitempty"`
		Dlink string `json:"dlink"`
		//Filename    string `json:"filename"`
		//FsId        int64  `json:"fs_id"`
		//Height      int    `json:"height,omitempty"`
		//Isdir       int    `json:"isdir"`
		//Md5         string `json:"md5"`
		//OperId      int    `json:"oper_id"`
		//Path        string `json:"path"`
		//ServerCtime int    `json:"server_ctime"`
		//ServerMtime int    `json:"server_mtime"`
		//Size        int    `json:"size"`
		//Thumbs      struct {
		//	Icon string `json:"icon,omitempty"`
		//	Url1 string `json:"url1,omitempty"`
		//	Url2 string `json:"url2,omitempty"`
		//	Url3 string `json:"url3,omitempty"`
		//} `json:"thumbs"`
		//Width int `json:"width,omitempty"`
	} `json:"list"`
	//Names struct {
	//} `json:"names"`
	RequestId string `json:"request_id"`
}

type DownloadResp2 struct {
	Errno int `json:"errno"`
	Info  []struct {
		//ExtentTinyint4 int `json:"extent_tinyint4"`
		//ExtentTinyint1 int `json:"extent_tinyint1"`
		//Bitmap string `json:"bitmap"`
		//Category int `json:"category"`
		//Isdir int `json:"isdir"`
		//Videotag int `json:"videotag"`
		Dlink string `json:"dlink"`
		//OperID int64 `json:"oper_id"`
		//PathMd5 int `json:"path_md5"`
		//Wpfile int `json:"wpfile"`
		//LocalMtime int `json:"local_mtime"`
		/*Thumbs struct {
			Icon string `json:"icon"`
			URL3 string `json:"url3"`
			URL2 string `json:"url2"`
			URL1 string `json:"url1"`
		} `json:"thumbs"`*/
		//PlaySource int `json:"play_source"`
		//Share int `json:"share"`
		//FileKey string `json:"file_key"`
		//Errno int `json:"errno"`
		//LocalCtime int `json:"local_ctime"`
		//Rotate int `json:"rotate"`
		//Metadata time.Time `json:"metadata"`
		//Height int `json:"height"`
		//SampleRate int `json:"sample_rate"`
		//Width int `json:"width"`
		//OwnerType int `json:"owner_type"`
		//Privacy int `json:"privacy"`
		//ExtentInt3 int64 `json:"extent_int3"`
		//RealCategory string `json:"real_category"`
		//SrcLocation string `json:"src_location"`
		//MetaInfo string `json:"meta_info"`
		//ID string `json:"id"`
		//Duration int `json:"duration"`
		//FileSize string `json:"file_size"`
		//Channels int `json:"channels"`
		//UseSegment int `json:"use_segment"`
		//ServerCtime int `json:"server_ctime"`
		//Resolution string `json:"resolution"`
		//OwnerID int `json:"owner_id"`
		//ExtraInfo string `json:"extra_info"`
		//Size int `json:"size"`
		//FsID int64 `json:"fs_id"`
		//ExtentTinyint3 int `json:"extent_tinyint3"`
		//Md5 string `json:"md5"`
		//Path string `json:"path"`
		//FrameRate int `json:"frame_rate"`
		//ExtentTinyint2 int `json:"extent_tinyint2"`
		//ServerFilename string `json:"server_filename"`
		//ServerMtime int `json:"server_mtime"`
		//TkbindID int `json:"tkbind_id"`
	} `json:"info"`
	RequestID int64 `json:"request_id"`
}

type PrecreateResp struct {
	Errno      int   `json:"errno"`
	RequestID  int64 `json:"request_id"`
	ReturnType int   `json:"return_type"`

	// return_type=1
	Path      string `json:"path"`
	Uploadid  string `json:"uploadid"`
	BlockList []int  `json:"block_list"`

	// return_type=2
	File File `json:"info"`
}

type QuotaResp struct {
	Errno     int    `json:"errno"`
	RequestId int64  `json:"request_id"`
	Total     uint64 `json:"total"`
	Used      uint64 `json:"used"`
	//Free      uint64 `json:"free"`
	//Expire    bool   `json:"expire"`
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
