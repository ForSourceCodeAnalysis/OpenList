package driver

import (
	"io"
	"time"
)

// FieldMeta 配置字段元数据
type FieldMeta struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Default  string `json:"default"`
	Options  string `json:"options"`
	Required bool   `json:"required"`
	Help     string `json:"help"`
}

// 请求参数定义，D前缀是为了与外部接口区分
// 如果有新增的请求参数，可以直接在结构体中添加
// 如果添加的参数，不具有通用性，务必添加详尽的注释
// 添加的参数，根据需要传递，根据需要取用

// IDPath 网盘文件ID，路径
// 根据需要传递，根据需要取用
type IDPath struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

// DListReq 列表请求参数
type DListReq struct {
	IDPath
	Offset   string `json:"offset"` // 不同的网盘有不同的含义，对应：123pan lastFileId，baidunetdisk start
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

// DGetReq 获取文件信息请求参数
type DGetReq struct {
	IDPath
}

// DMkdirReq 创建文件夹请求参数
type DMkdirReq struct {
	IDPath
	Name string `json:"name"`
}

// RenameObj 重命名对象
type RenameObj struct {
	IDPath
	NewName string `json:"new_name"`
}

// DRenameReq 重命名文件请求参数
type DRenameReq struct {
	Objs []RenameObj `json:"objs"`
}

// DMoveReq 移动文件请求参数
type DMoveReq struct {
	SrcObjs []IDPath `json:"src_objs"`
	DstDir  IDPath   `json:"dst_dir"`
}

// DCopyReq 复制文件请求参数
type DCopyReq struct {
	SrcObjs []IDPath `json:"src_objs"`
	DstDir  IDPath   `json:"dst_dir"`
}

// DRemoveReq 删除文件请求参数
type DRemoveReq struct {
	Objs []IDPath `json:"objs"`
}

// DLinkReq 获取文件链接请求参数
type DLinkReq struct {
	IDPath
}

// DPreupReq 预上传请求参数
// 预上传有点问题，百度网盘预先规定了分片大小，所以上传请求需要传递分片hash，而123pan则是需要先请求接口返回分片大小
// 这就要求前端在发起请求时需要先知道网盘类型，才能正确传递参数
// 前端可以根据
type DPreupReq struct {
	DstDir    IDPath    `json:"dst_dir"`
	SrcObj    DPreupObj `json:"src_obj"`
	OverWrite bool      `json:"overwrite"`
}

// DPreupObj 预上传对象
type DPreupObj struct {
	IDPath
	Name      string              `json:"name"`
	Size      int64               `json:"size"`
	IsDir     bool                `json:"is_dir"`
	Hash      map[string]string   // 文件hash
	SliceHash []map[string]string //分片hash 按分片顺序 百度网盘
}

// DSliceUploadReq 分片上传请求参数
type DSliceUploadReq struct {
	// DPreupObj           // 需要比对之前缓存的预上传信息
	PreupID string    `json:"preup_id"`
	SliceNo int       `json:"slice_no"`
	DstPath IDPath    `json:"dst_path"` // 绝对路径，包含文件名，百度网盘需要
	File    io.Reader `json:"-"`
}

// DUploadDoneReq 分片上传完成请求参数
type DUploadDoneReq struct {
	PreupID    string    `json:"preup_id"`
	DstPath    IDPath    `json:"dst_path"` // 绝对路径，包含文件夹名
	Size       int64     `json:"size"`
	IsDir      bool      `json:"is_dir"`
	BlockList  []string  `json:"block_list"`  // 文件各分片hash，参考 https://pan.baidu.com/union/doc/rksg0sa17
	Rtype      string    `json:"rtype"`       // 文件命名策略，定义成string是为了方便后续兼容性 参考 https://pan.baidu.com/union/doc/rksg0sa17
	LocalCtime time.Time `json:"local_ctime"` // 本地创建时间 参考 https://pan.baidu.com/union/doc/rksg0sa17
	LocalMtime time.Time `json:"local_mtime"` // 本地修改时间 参考 https://pan.baidu.com/union/doc/rksg0sa17

}

// DPutReq 上传文件参数
type DPutReq struct {
	DstDir IDPath    `json:"dst_dir"`
	File   io.Reader `json:"-"`
}
