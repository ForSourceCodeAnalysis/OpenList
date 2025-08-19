package open123

import (
	"context"
	"fmt"
	"mime/multipart"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/OpenListTeam/OpenList/v4/internal/driver"
	"github.com/OpenListTeam/OpenList/v4/internal/errs"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/model/reqres"
	"github.com/OpenListTeam/OpenList/v4/internal/model/tables"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	"github.com/OpenListTeam/OpenList/v4/internal/stream"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Open123 struct {
	model.Storage
	Addition
}

func (d *Open123) Config() driver.Config {
	return config
}

func (d *Open123) GetAddition() driver.Additional {
	return &d.Addition
}

func (d *Open123) GetUploadInfo() *model.UploadInfo {
	return &model.UploadInfo{
		SliceHashNeed: true,
		HashMd5Need:   true,
	}
}

func (d *Open123) Init(ctx context.Context) error {
	if d.UploadThread < 1 || d.UploadThread > 32 {
		d.UploadThread = 3
	}

	return nil
}

func (d *Open123) Drop(ctx context.Context) error {
	op.MustSaveDriverStorage(d)
	return nil
}

func (d *Open123) List(ctx context.Context, dir model.Obj, args model.ListArgs) ([]model.Obj, error) {
	fileLastId := int64(0)
	parentFileId, err := strconv.ParseInt(dir.GetID(), 10, 64)
	if err != nil {
		return nil, err
	}
	res := make([]File, 0)

	for fileLastId != -1 {
		files, err := d.getFiles(parentFileId, 100, fileLastId)
		if err != nil {
			return nil, err
		}
		// 目前123panAPI请求，trashed失效，只能通过遍历过滤
		for i := range files.FileList {
			if files.FileList[i].Trashed == 0 {
				res = append(res, files.FileList[i])
			}
		}
		fileLastId = files.LastFileId
	}
	return utils.SliceConvert(res, func(src File) (model.Obj, error) {
		return src, nil
	})
}

func (d *Open123) Link(ctx context.Context, file model.Obj, args model.LinkArgs) (*model.Link, error) {
	fileId, _ := strconv.ParseInt(file.GetID(), 10, 64)

	if d.DirectLink {
		res, err := d.getDirectLink(fileId)
		if err != nil {
			return nil, err
		}

		if d.DirectLinkPrivateKey == "" {
			duration := 365 * 24 * time.Hour // 缓存1年
			return &model.Link{
				URL:        res.Data.URL,
				Expiration: &duration,
			}, nil
		}

		u, err := d.getUserInfo()
		if err != nil {
			return nil, err
		}

		duration := time.Duration(d.DirectLinkValidDuration) * time.Minute

		newURL, err := d.SignURL(res.Data.URL, d.DirectLinkPrivateKey,
			u.Data.UID, duration)
		if err != nil {
			return nil, err
		}

		return &model.Link{
			URL:        newURL,
			Expiration: &duration,
		}, nil
	}

	res, err := d.getDownloadInfo(fileId)
	if err != nil {
		return nil, err
	}

	link := model.Link{URL: res.DownloadUrl}
	return &link, nil
}

func (d *Open123) MakeDir(ctx context.Context, parentDir model.Obj, dirName string) error {
	parentFileId, _ := strconv.ParseInt(parentDir.GetID(), 10, 64)

	return d.mkdir(parentFileId, dirName)
}

func (d *Open123) Move(ctx context.Context, srcObj, dstDir model.Obj) error {
	toParentFileID, _ := strconv.ParseInt(dstDir.GetID(), 10, 64)

	return d.move(srcObj.(File).FileId, toParentFileID)
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
	// 每次最多30
	for names := range slices.Chunk(rl, 30) {
		if err := d.batchRename(names); err != nil {
			return err
		}
	}

	return nil
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
	createResp, err := d.uploadCreate(&UploadCreateReq{
		ParentFileID: uint64(parentFileId),
		FileName:     srcObj.GetName(),
		Etag:         etag,
		Size:         srcObj.GetSize(),
		Duplicate:    2,
		ContainDir:   false,
	})

	if err != nil {
		return err
	}
	// 是否秒传
	if createResp.Reuse {
		return nil
	}
	return errs.NotSupport
}

// Remove 删除文件
func (d *Open123) Remove(ctx context.Context, obj model.Obj) error {
	fileId, _ := strconv.ParseInt(obj.GetID(), 10, 64)

	return d.trash([]int64{fileId})
}

// BatchRemove 批量删除
func (d *Open123) BatchRemove(ctx context.Context, srcObj model.Obj, objs []model.IDName) error {

	ids := []int64{}
	for _, obj := range objs {
		id, err := strconv.ParseInt(obj.ID, 10, 64)
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}
	//每次最多100
	for cids := range slices.Chunk(ids, 100) {
		if err := d.trash(cids); err != nil {
			return err
		}
	}
	return nil
}

// Put 单次上传
func (d *Open123) Put(ctx context.Context, dstDir model.Obj, file model.FileStreamer, up driver.UpdateProgress) error {
	parentFileID, err := strconv.ParseInt(dstDir.GetID(), 10, 64)
	etag := file.GetHash().GetHash(utils.MD5)
	logrus.Infof("file %+v", file)

	if len(etag) < utils.MD5.Width {
		_, etag, err = stream.CacheFullAndHash(file, &up, utils.MD5)
		if err != nil {
			logrus.Errorf("cache full and hash error: %v", err)
			return err
		}
	}
	up(10)

	return d.singleUpload(&SingleUploadReq{
		ParentFileID: parentFileID,
		FileName:     file.GetName(),
		Etag:         etag,
		Size:         file.GetSize(),
		File:         file,
		Duplicate:    1,
	})
}

// Preup 预上传
func (d *Open123) Preup(c context.Context, srcobj model.Obj, req *reqres.PreupReq) (*model.PreupInfo, error) {
	pid, err := strconv.ParseUint(srcobj.GetID(), 10, 64)
	if err != nil {
		return nil, err
	}
	duplicate := 1
	if req.Overwrite {
		duplicate = 2
	}

	ucr := &UploadCreateReq{
		ParentFileID: pid,
		Etag:         req.Hash.Md5,
		FileName:     req.Name,
		Size:         int64(req.Size),
		Duplicate:    duplicate,
	}

	resp, err := d.uploadCreate(ucr)
	if err != nil {
		return nil, err
	}
	return &model.PreupInfo{
		PreupID:   resp.PreuploadID,
		Server:    resp.Servers[0],
		SliceSize: resp.SliceSize,
		Reuse:     resp.Reuse,
	}, nil
}

// UploadSlice 上传分片
func (d *Open123) SliceUpload(c context.Context, req *tables.SliceUpload, sliceno uint, fd multipart.File) error {
	sh := strings.Split(req.SliceHash, ",")
	r := &UploadSliceReq{
		Name:        req.Name,
		PreuploadID: req.PreupID,
		Server:      req.Server,
		Slice:       fd,
		SliceMD5:    sh[sliceno],
		SliceNo:     int(sliceno) + 1,
	}
	return d.uploadSlice(r)
}

// UploadSliceComplete 分片上传完成
func (d *Open123) UploadSliceComplete(c context.Context, su *tables.SliceUpload) error {

	return d.sliceUpComplete(su.PreupID)
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
// Preup 预上传
func (d *Open123) Preup(c context.Context, srcobj model.Obj, req *reqres.PreupReq) (*model.PreupInfo, error) {
	pid, err := strconv.ParseUint(srcobj.GetID(), 10, 64)
	if err != nil {
		return nil, err
	}
	duplicate := 1
	if req.Overwrite {
		duplicate = 2
	}

	ucr := &UploadCreateReq{
		ParentFileID: pid,
		Etag:         req.Hash.Md5,
		FileName:     req.Name,
		Size:         int64(req.Size),
		Duplicate:    duplicate,
	}

	resp, err := d.uploadCreate(ucr)
	if err != nil {
		return nil, err
	}
	return &model.PreupInfo{
		PreupID:   resp.PreuploadID,
		Server:    resp.Servers[0],
		SliceSize: resp.SliceSize,
		Reuse:     resp.Reuse,
	}, nil
}

// UploadSlice 上传分片
func (d *Open123) SliceUpload(c context.Context, req *tables.SliceUpload, sliceno uint, fd io.Reader) error {
	sh := strings.Split(req.SliceHash, ",")
	r := &UploadSliceReq{
		Name:        req.Name,
		PreuploadID: req.PreupID,
		Server:      req.Server,
		Slice:       fd,
		SliceMD5:    sh[sliceno],
		SliceNo:     int(sliceno) + 1,
	}
	return d.uploadSlice(r)
}

// UploadSliceComplete 分片上传完成
func (d *Open123) UploadSliceComplete(c context.Context, su *tables.SliceUpload) error {

	return d.sliceUpComplete(su.PreupID)
}

var _ driver.Driver = (*Open123)(nil)
