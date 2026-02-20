package fs

import (
	"context"
	"path"
	"slices"

	"github.com/OpenListTeam/OpenList/v4/extensions/interfaces"
	"github.com/OpenListTeam/OpenList/v4/extensions/models"
	"github.com/OpenListTeam/OpenList/v4/internal/driver"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// BatchRename 批量重命名
func BatchRename(ctx context.Context, storage driver.Driver, srcDir string, renameObjs []models.RenameObj) error {
	srcRawObj, err := op.Get(ctx, storage, srcDir)
	if err != nil {
		return errors.WithMessage(err, "failed to get src object")
	}
	srcObj := model.UnwrapObjName(srcRawObj)

	list, err := storage.List(ctx, srcObj, model.ListArgs{})
	if err != nil {
		logrus.Errorf("failed to list src object before rename,srcdir: %s, err: %+v", srcDir, err)
		return errors.WithMessage(err, "failed to list src object before rename")
	}
	var romap = make(map[string]string)
	for _, obj := range renameObjs {
		romap[obj.SrcName] = obj.NewName
	}
	var driverRenameObjs []models.DriverRenameObj
	for _, obj := range list {
		if newname, ok := romap[obj.GetName()]; ok {
			driverRenameObjs = append(driverRenameObjs, models.DriverRenameObj{NewName: newname, Obj: obj})
		}
	}

	switch s := storage.(type) {
	case interfaces.IBatchRename:
		err := s.BatchRename(ctx, srcObj, driverRenameObjs)
		if err == nil {
			op.Cache.DeleteDirectory(storage, srcDir)
			return nil
		}
		return err
	default:
		for _, renameObject := range renameObjs {
			err := op.Rename(ctx, storage, path.Join(srcDir, renameObject.SrcName), renameObject.NewName)
			if err != nil {
				logrus.Errorf("failed rename %s to %s: %+v", renameObject.SrcName, renameObject.NewName, err)
				return err
			}
		}
	}
	return nil
}

// BatchRemove 批量删除
func BatchRemove(ctx context.Context, storage driver.Driver, srcDir string, delObjs []string) error {
	srcobj, err := op.Get(ctx, storage, srcDir)
	if err != nil {
		return errors.WithMessage(err, "failed to get src object")
	}
	srcobj = model.UnwrapObjName(srcobj)

	list, err := storage.List(ctx, srcobj, model.ListArgs{})
	if err != nil {
		logrus.Errorf("failed to list src object before remove,srcDir: %s, err: %+v", srcDir, err)
		return errors.WithMessage(err, "failed to list src object before remove")
	}
	var driverObjs []model.Obj
	for _, obj := range list {
		if slices.Contains(delObjs, obj.GetName()) {
			driverObjs = append(driverObjs, obj)
		}
	}

	switch s := storage.(type) {
	case interfaces.IBatchRemove:
		err = s.BatchRemove(ctx, srcobj, driverObjs)
		if err == nil {
			op.Cache.DeleteDirectory(storage, srcDir)
			return nil
		}
		return err
	case driver.Remove:
		for _, obj := range driverObjs {
			err = s.Remove(ctx, model.UnwrapObjName(obj))
			if err != nil {
				logrus.Errorf("failed remove %s: %+v", obj.GetName(), err)
				return err
			}
		}
		op.Cache.DeleteDirectory(storage, srcDir)
	default:
		return errors.WithMessage(err, "driver not support remove")

	}
	return nil
}
