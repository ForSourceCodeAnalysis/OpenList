package open123

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"errors"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/OpenListTeam/OpenList/v4/internal/driver"
	"github.com/OpenListTeam/OpenList/v4/internal/errs"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/Xhofe/go-cache"
)

// List 获取目录下的文件列表
func (d *Open123) List(ctx context.Context, req *driver.DListReq) (driver.IListResp, error) {
	pid, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return nil, err
	}
	// 分页获取
	if req.Page > 0 && req.Offset != "" {
		if req.PageSize <= 0 || req.PageSize > 100 {
			req.PageSize = 100
		}
		lastFileID, err := strconv.ParseInt(req.Offset, 10, 64)
		if err != nil {
			return nil, err
		}

		r, err := d.getFiles(pid, req.PageSize, lastFileID)
		if err != nil {
			return nil, err
		}
		return r, nil

	}
	var lastFileID int64 = 0
	res := make([]*Item, 0)

	// 获取全部
	for lastFileID != -1 {

		lr, err := d.getFiles(pid, 100, lastFileID)
		if err != nil {
			return nil, err
		}
		// 目前123panAPI请求，trashed失效，只能通过遍历过滤
		for i := range lr.FileList {
			if lr.FileList[i].Trashed == 0 {
				res = append(res, lr.FileList[i])
			}
		}
		lastFileID = lr.LastFileID
	}
	return &ListObj{FileList: res, LastFileID: lastFileID}, nil
}

// Link 获取链接
func (d *Open123) Link(ctx context.Context, req *driver.DLinkReq) (driver.ILinkResp, error) {
	fileID, _ := strconv.ParseInt(req.ID, 10, 64)

	res, err := d.getDownloadInfo(fileID)
	if err != nil {
		return nil, err
	}

	return res, nil
	link := model.Link{URL: res.DownloadUrl}
	return &link, nil
}

// MkDir 创建文件夹
func (d *Open123) MkDir(ctx context.Context, req *driver.DMkdirReq) error {
	parentFileID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return err
	}

	return d.mkdir(parentFileID, req.Name)
}

// Move 移动文件
func (d *Open123) Move(ctx context.Context, req *driver.DMoveReq) error {
	toParentFileID, err := strconv.ParseInt(req.DstDir.ID, 10, 64)
	if err != nil {
		return err
	}
	ids := make([]int64, len(req.SrcObjs))
	for i, srcObj := range req.SrcObjs {
		id, err := strconv.ParseInt(srcObj.ID, 10, 64)
		if err != nil {
			return err
		}
		ids[i] = id

	}

	return d.move(ids, toParentFileID)
}

func (d *Open123) Rename(ctx context.Context, srcObj model.Obj, newName string) error {
	fileId, _ := strconv.ParseInt(srcObj.GetID(), 10, 64)

	return d.rename(fileId, newName)
}

func (d *Open123) BatchRename(ctx context.Context, obj model.Obj, renameObjs []model.RenameObj) error {
	rl := []string{}
	for _, ro := range renameObjs {
		fileID, err := strconv.ParseInt(ro.ID, 10, 64)
		if err != nil {
			return err
		}
		rl = append(rl, fmt.Sprintf("%d|%s", fileID, ro.NewName))

	}

	return d.batchRename(rl)
}

func (d *Open123) Copy(ctx context.Context, srcObj, dstDir model.Obj) error {
	// 尝试使用上传+MD5秒传功能实现复制
	// 1. 创建文件
	// parentFileID 父目录id，上传到根目录时填写 0
	parentFileId, err := strconv.ParseInt(dstDir.GetID(), 10, 64)
	if err != nil {
		return fmt.Errorf("parse parentFileID error: %v", err)
	}
	etag := srcObj.(File).Etag
	createResp, err := d.create(parentFileId, srcObj.GetName(), etag, srcObj.GetSize(), 2, false)
	if err != nil {
		return err
	}
	// 是否秒传
	if createResp.Data.Reuse {
		return nil
	}
	return errs.NotSupport
}

func (d *Open123) Remove(ctx context.Context, obj model.Obj) error {
	fileId, _ := strconv.ParseInt(obj.GetID(), 10, 64)

	return d.trash(fileId)
}

// func (d *Open123) Put(ctx context.Context, dstDir model.Obj, file model.FileStreamer, up driver.UpdateProgress) error {
// 	parentFileId, err := strconv.ParseInt(dstDir.GetID(), 10, 64)
// 	etag := file.GetHash().GetHash(utils.MD5)

// 	if len(etag) < utils.MD5.Width {
// 		_, etag, err = stream.CacheFullInTempFileAndHash(file, utils.MD5)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	createResp, err := d.create(parentFileId, file.GetName(), etag, file.GetSize(), 2, false)
// 	if err != nil {
// 		return err
// 	}
// 	if createResp.Data.Reuse {
// 		return nil
// 	}
// 	up(10)

// 	return d.Upload(ctx, file, createResp, up)
// }

// Copy 复制文件
func (d *Open123) Copy(ctx context.Context, req *driver.DCopyReq) error {
	return errs.NotSupport
}

// Remove 删除文件
func (d *Open123) Remove(ctx context.Context, req *driver.DRemoveReq) error {
	ids := make([]int64, 0, len(req.Objs))
	for i, obj := range req.Objs {
		id, err := strconv.ParseInt(obj.ID, 10, 64)
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}
	//每次最多100
	for cids := range slices.Chunk(ids, 100) {
		if err := d.trash(cids); err != nil {
			return err
		}
	}
	return nil
}

func (d *Open123) GetDetails(ctx context.Context) (*model.StorageDetails, error) {
	userInfo, err := d.getUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	total := userInfo.Data.SpacePermanent + userInfo.Data.SpaceTemp
	free := total - userInfo.Data.SpaceUsed
	return &model.StorageDetails{
		DiskUsage: model.DiskUsage{
			TotalSpace: total,
			FreeSpace:  free,
		},
	}, nil
}

var (
	_ driver.Driver    = (*Open123)(nil)
	_ driver.PutResult = (*Open123)(nil)
)

// Preup 上传预处理
func (d *Open123) Preup(ctx context.Context, req *driver.DPreupReq) (driver.IPreupResp, error) {
	etag := req.SrcObj.Hash[driver.HashTypeMD5]
	if etag == "" {
		return nil, errors.New("md5 is empty")
	}

	ParentFileID, err := strconv.ParseInt(req.DstDir.ID, 10, 64)
	if err != nil {
		return nil, err
	}

	duplicate := 1
	if req.OverWrite {
		duplicate = 2
	}

	preupres, err := d.preup(&PreupReq{
		ParentFileID: ParentFileID,
		Filename:     req.SrcObj.Name,
		Etag:         etag,
		Size:         req.SrcObj.Size,
		Duplicate:    duplicate,
		ContainDir:   false,
	})
	if err != nil {
		return nil, err
	}
	// 预上传成功，缓存信息，官方文档没有明确说明预上传的有效期，暂且设置为48h
	driver.DriverCache.Set(preupres.PreuploadID, SliceUploadCache{
		Filename:          req.SrcObj.Name,
		Hash:              etag,
		PreupID:           preupres.PreuploadID,
		Size:              req.SrcObj.Size,
		SliceSize:         preupres.SliceSize,
		UploadServer:      preupres.Servers[0],
		UploadedBlockList: []int{},
		SliceHash:         req.SrcObj.SliceHash,
	}, cache.WithEx[any](time.Hour*48))
	bl := []int{}
	for range req.SrcObj.Size / preupres.SliceSize {
		bl = append(bl, 0)

	}
	preupinfo := &PreupInfo{
		BlockList: bl,
		PreupResp: preupres,
	}

	return preupinfo, nil

}

// getUpload
func getUploadSlices(size int64, sliceSize int64, uploaded []int) []int {
	totalSlices := size/sliceSize + 1

	exists := make(map[int]struct{})
	for _, chunk := range uploaded {
		exists[chunk] = struct{}{}
	}

	missing := []int{}
	for i := 0; i < int(totalSlices); i++ {
		if _, ok := exists[i]; !ok {
			missing = append(missing, i)
		}
	}

	return missing
}

// PutSlice 上传分片
func (d *Open123) PutSlice(req *driver.DSliceUploadReq) error {

	return nil
}
