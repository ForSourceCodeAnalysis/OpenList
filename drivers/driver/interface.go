package driver

import (
	"context"
	"time"
)

// 接口定义尽量遵循单一职责原则， 接口本身支持嵌套，扩展更灵活
// 输入参数不定义成接口，而是结构体，因为输入本身就需要具体的实现，定义成接口，反而不方便
// 输出则尽量定义成接口，方便后续处理
// 所以这里定义的接口有两种，
// 	第一种是为驱动定义，驱动需要实现对应的方法
//	第二种是为返回数据定义，返回数据需要实现对应方法，以Resp结尾进行区分

// IDriver 网盘驱动
// 字段，方法的定义保持在对应的驱动内
// 驱动内的逻辑保持仅驱动相关，具体业务相关的，应该放在外部层处理
type IDriver interface {
	IMeta
	IAPI
}

// IMeta 网盘元信息
type IMeta interface {
	Init() error
	// Name 返回驱动名称
	Name() string

	// SortSupported 是否支持排序
	SortSupported() bool

	// PaganitionSupported 是否支持分页
	// PaganitionSupported() bool

	// ConfigMeta 驱动配置描述
	// 这个功能主要是给前端使用的。添加存储时，每个驱动的配置项是不同的，通过这个方法返回
	// ConfigMeta() []FieldMeta
}

// IAPI 网盘API
type IAPI interface {
	// List files in the path
	// if identify files by path, need to set ID with path,like path.Join(dir.GetID(), obj.GetName())
	// if identify files by id, need to set ID with corresponding id
	List(ctx context.Context, req DListReq) (IListResp, error)
	// Get file info
	Get(ctx context.Context, req DGetReq) (IObjResp, error)
	// MkDir create a directory
	MkDir(ctx context.Context, req DMkdirReq) error
	// Rename rename file, support batch rename
	Rename(ctx context.Context, req DRenameReq) error
	// Move move files
	Move(ctx context.Context, req DMoveReq) error
	// Copy copy files
	Copy(ctx context.Context, req DCopyReq) error
	// Remove delete files
	Remove(ctx context.Context, req DRemoveReq) error
}

// IIDPathResp get ID and Path
type IIDPathResp interface {
	GetID() string
	GetPath() string
}

// IModTimeResp 返回修改时间
type IModTimeResp interface {
	ModTime() *time.Time
}

// ICreateTimeResp 创建时间
type ICreateTimeResp interface {
	CreateTime() *time.Time
}

// IListResp 列表对象
type IListResp interface {
	ListData() []IObjResp
	Offset() string
}

// IObjResp 网盘对象
type IObjResp interface {
	GetSize() int64
	GetName() string
	// ModTime() time.Time
	// CreateTime() time.Time
	IsDir() bool
	IIDPathResp
}

// IHashResp 网盘对象Hash
type IHashResp interface {
	GetHash() map[string]string
}

// IThumbnailResp 网盘对象缩略图
type IThumbnailResp interface {
	Thumbnail() string
}

// ILink 获取网盘对象链接
type ILink interface {
	Link(ctx context.Context, req DLinkReq) ILinkResp
}

// ILinkResp 网盘对象链接信息
type ILinkResp interface {
	GetURL() string
}

// 上传功能略微有些复杂，涉及到分片上传，单次上传，每个网盘的实现也不太一样
// 如果网盘原生支持分片上传，那么可以直接利用官方api实现
// 如果网盘不支持分片上传，可以通过openlist缓存的方式实现
// 上传逻辑如下：
// 1. 客户端先请求Preup接口，获取预上传信息，前后端都需要缓存必要的返回信息，以便出现问题可以根据相关信息恢复
//
// 2. 客户端根据预上传信息进行分割上传
// 3. 预上传完成后，客户端请求IUploadDone接口，完成上传

// IPreupload 预上传
type IPreupload interface {
	Preup(ctx context.Context, req DPreupReq) (IPreupResp, error)
}

// IPreupResp 预上传响应
type IPreupResp interface {
	GetPreupID() string
	GetSliceSize() int64
	IsRapidUpload() bool // 是否秒传
	GetBlockList() []int // 需要上传的块索引
}

// IUploadBlockListResp 上传块列表
// type IUploadBlockListResp interface {
// 	GetBlockList() []int
// }

// IUploadServerResp 获取上传服务器
// type IUploadServerResp interface {
// 	GetUploadServer() string
// }

// ISliceUplod 分片上传
type ISliceUplod interface {
	PutSlice(ctx context.Context, req DSliceUploadReq) error
}

// IUploadDone 分片上传完成
type IUploadDone interface {
	UploadDone(ctx context.Context, req DUploadDoneReq) error
}

// IPut 单次上传
// 限速和进度处理
type IPut interface {
	// Put a file (provided as a FileStreamer) into the driver
	// Besides the most basic upload functionality, the following features also need to be implemented:
	// 1. Canceling (when `<-ctx.Done()` returns), which can be supported by the following methods:
	//   (1) Use request methods that carry context, such as the following:
	//      a. http.NewRequestWithContext
	//      b. resty.Request.SetContext
	//      c. s3manager.Uploader.UploadWithContext
	//      d. utils.CopyWithCtx
	//   (2) Use a `driver.ReaderWithCtx` or `driver.NewLimitedUploadStream`
	//   (3) Use `utils.IsCanceled` to check if the upload has been canceled during the upload process,
	//       this is typically applicable to chunked uploads.
	// 2. Submit upload progress (via `up`) in real-time. There are three recommended ways as follows:
	//   (1) Use `utils.CopyWithCtx`
	//   (2) Use `driver.ReaderUpdatingProgress`
	//   (3) Use `driver.Progress` with `io.TeeReader`
	// 3. Slow down upload speed (via `stream.ServerUploadLimit`). It requires you to wrap the read stream
	//    in a `driver.RateLimitReader` or a `driver.RateLimitFile` after calculating the file's hash and
	//    before uploading the file or file chunks. Or you can directly call `driver.ServerUploadLimitWaitN`
	//    if your file chunks are sufficiently small (less than about 50KB).
	// NOTE that the network speed may be significantly slower than the stream's read speed. Therefore, if
	// you use a `errgroup.Group` to upload each chunk in parallel, you should use `Group.SetLimit` to
	// limit the maximum number of upload threads, preventing excessive memory usage caused by buffering
	// too many file chunks awaiting upload.
	Put(ctx context.Context, req DPutReq) error
}

// type PutURL interface {
// 	// PutURL directly put a URL into the storage
// 	// Applicable to index-based drivers like URL-Tree or drivers that support uploading files as URLs
// 	// Called when using SimpleHttp for offline downloading, skipping creating a download task
// 	PutURL(ctx context.Context, dstDir model.Obj, name, url string) error
// }

// type PutResult interface {
// 	// Put a file (provided as a FileStreamer) into the driver and return the put obj
// 	// Besides the most basic upload functionality, the following features also need to be implemented:
// 	// 1. Canceling (when `<-ctx.Done()` returns), which can be supported by the following methods:
// 	//   (1) Use request methods that carry context, such as the following:
// 	//      a. http.NewRequestWithContext
// 	//      b. resty.Request.SetContext
// 	//      c. s3manager.Uploader.UploadWithContext
// 	//      d. utils.CopyWithCtx
// 	//   (2) Use a `driver.ReaderWithCtx` or `driver.NewLimitedUploadStream`
// 	//   (3) Use `utils.IsCanceled` to check if the upload has been canceled during the upload process,
// 	//       this is typically applicable to chunked uploads.
// 	// 2. Submit upload progress (via `up`) in real-time. There are three recommended ways as follows:
// 	//   (1) Use `utils.CopyWithCtx`
// 	//   (2) Use `driver.ReaderUpdatingProgress`
// 	//   (3) Use `driver.Progress` with `io.TeeReader`
// 	// 3. Slow down upload speed (via `stream.ServerUploadLimit`). It requires you to wrap the read stream
// 	//    in a `driver.RateLimitReader` or a `driver.RateLimitFile` after calculating the file's hash and
// 	//    before uploading the file or file chunks. Or you can directly call `driver.ServerUploadLimitWaitN`
// 	//    if your file chunks are sufficiently small (less than about 50KB).
// 	// NOTE that the network speed may be significantly slower than the stream's read speed. Therefore, if
// 	// you use a `errgroup.Group` to upload each chunk in parallel, you should use `Group.SetLimit` to
// 	// limit the maximum number of upload threads, preventing excessive memory usage caused by buffering
// 	// too many file chunks awaiting upload.
// 	Put(ctx context.Context, dstDir model.Obj, file model.FileStreamer, up UpdateProgress) (model.Obj, error)
// }

// type PutURLResult interface {
// 	// PutURL directly put a URL into the storage
// 	// Applicable to index-based drivers like URL-Tree or drivers that support uploading files as URLs
// 	// Called when using SimpleHttp for offline downloading, skipping creating a download task
// 	PutURL(ctx context.Context, dstDir model.Obj, name, url string) (model.Obj, error)
// }

// type ArchiveReader interface {
// 	// GetArchiveMeta get the meta-info of an archive
// 	// return errs.WrongArchivePassword if the meta-info is also encrypted but provided password is wrong or empty
// 	// return errs.NotImplement to use internal archive tools to get the meta-info, such as the following cases:
// 	// 1. the driver do not support the format of the archive but there may be an internal tool do
// 	// 2. handling archives is a VIP feature, but the driver does not have VIP access
// 	GetArchiveMeta(ctx context.Context, obj model.Obj, args model.ArchiveArgs) (model.ArchiveMeta, error)
// 	// ListArchive list the children of model.ArchiveArgs.InnerPath in the archive
// 	// return errs.NotImplement to use internal archive tools to list the children
// 	// return errs.NotSupport if the folder structure should be acquired from model.ArchiveMeta.GetTree
// 	ListArchive(ctx context.Context, obj model.Obj, args model.ArchiveInnerArgs) ([]model.Obj, error)
// 	// Extract get url/filepath/reader of a file in the archive
// 	// return errs.NotImplement to use internal archive tools to extract
// 	Extract(ctx context.Context, obj model.Obj, args model.ArchiveInnerArgs) (*model.Link, error)
// }

// type ArchiveGetter interface {
// 	// ArchiveGet get file by inner path
// 	// return errs.NotImplement to use internal archive tools to get the children
// 	// return errs.NotSupport if the folder structure should be acquired from model.ArchiveMeta.GetTree
// 	ArchiveGet(ctx context.Context, obj model.Obj, args model.ArchiveInnerArgs) (model.Obj, error)
// }

// type ArchiveDecompress interface {
// 	ArchiveDecompress(ctx context.Context, srcObj, dstDir model.Obj, args model.ArchiveDecompressArgs) error
// }

// type ArchiveDecompressResult interface {
// 	// ArchiveDecompress decompress an archive
// 	// when args.PutIntoNewDir, the new sub-folder should be named the same to the archive but without the extension
// 	// return each decompressed obj from the root path of the archive when args.PutIntoNewDir is false
// 	// return only the newly created folder when args.PutIntoNewDir is true
// 	// return errs.NotImplement to use internal archive tools to decompress
// 	ArchiveDecompress(ctx context.Context, srcObj, dstDir model.Obj, args model.ArchiveDecompressArgs) ([]model.Obj, error)
// }
