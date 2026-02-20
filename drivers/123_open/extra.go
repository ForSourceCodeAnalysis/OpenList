package _123_open

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/OpenListTeam/OpenList/v4/drivers/base"
	"github.com/OpenListTeam/OpenList/v4/extensions/models"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/go-resty/resty/v2"
)

var (
	batchRename = InitApiInfo(Api+"/api/v1/file/rename", 1)
)

// BatchRename 批量重命名文件
func (d *Open123) BatchRename(ctx context.Context, srcDir model.Obj, renameObjs []models.DriverRenameObj) error {
	rl := []string{}
	for _, ro := range renameObjs {
		fileID := ro.GetID()
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

func (d *Open123) batchRename(renamelist []string) error {
	_, err := d.Request(batchRename, http.MethodPost, func(req *resty.Request) {
		req.SetBody(base.Json{
			"renameList": renamelist,
		})
	}, nil)
	return err
}

// BatchRemove 批量删除
func (d *Open123) BatchRemove(ctx context.Context, srcDir model.Obj, objs []model.Obj) error {

	ids := []int64{}
	for _, obj := range objs {
		id, err := strconv.ParseInt(obj.GetID(), 10, 64)
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}
	//每次最多100
	for cids := range slices.Chunk(ids, 100) {
		if err := d.trashX(cids); err != nil {
			return err
		}
	}
	return nil
}

func (d *Open123) trashX(fileIds []int64) error {
	_, err := d.Request(Trash, http.MethodPost, func(req *resty.Request) {
		req.SetBody(base.Json{
			"fileIDs": fileIds,
		})
	}, nil)
	if err != nil {
		return err
	}

	return nil
}
