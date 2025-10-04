package driver

import (
	"context"
	"io"

	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/model/reqres"
	"github.com/OpenListTeam/OpenList/v4/internal/model/tables"
)

type BatchRename interface {
	BatchRename(ctx context.Context, obj model.Obj, renameObjs []model.RenameObj) error
}

type BatchRemove interface {
	BatchRemove(ctx context.Context, srcobj model.Obj, objs []model.IDName) error
}

// IUploadInfo 上传信息接口
type IUploadInfo interface {
	GetUploadInfo() *model.UploadInfo
}

// IPreup 预上传接口
type IPreup interface {
	Preup(ctx context.Context, srcobj model.Obj, req *reqres.PreupReq) (*model.PreupInfo, error)
}

// ISliceUpload 分片上传接口
type ISliceUpload interface {
	// SliceUpload 分片上传
	SliceUpload(ctx context.Context, req *tables.SliceUpload, sliceno uint, file io.Reader) error
}

// IUploadSliceComplete 分片上传完成接口
type IUploadSliceComplete interface {
	UploadSliceComplete(ctx context.Context, req *tables.SliceUpload) error
}
