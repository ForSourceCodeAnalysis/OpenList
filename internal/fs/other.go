package fs

import (
	"context"
	"path/filepath"

	"github.com/OpenListTeam/OpenList/v4/internal/driver"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	"github.com/OpenListTeam/OpenList/v4/internal/task"
	"github.com/pkg/errors"
)

func makeDir(ctx context.Context, path string, lazyCache ...bool) error {
	storage, actualPath, err := op.GetStorageAndActualPath(path)
	if err != nil {
		return errors.WithMessage(err, "failed get storage")
	}
	return op.MakeDir(ctx, storage, actualPath, lazyCache...)
}

func rename(ctx context.Context, srcPath, dstName string, lazyCache ...bool) error {
	storage, srcActualPath, err := op.GetStorageAndActualPath(srcPath)
	if err != nil {
		return errors.WithMessage(err, "failed get storage")
	}
	return op.Rename(ctx, storage, srcActualPath, dstName, lazyCache...)
}

func batchRename(ctx context.Context, srcPath string, renameObjs []model.RenameObj, lazyCache ...bool) error {
	storage, srcActualPath, err := op.GetStorageAndActualPath(srcPath)
	if err != nil {
		return errors.WithMessage(err, "failed get storage")
	}

	return op.BatchRename(ctx, storage, srcActualPath, renameObjs, lazyCache...)

}

func remove(ctx context.Context, path string) error {
	storage, actualPath, err := op.GetStorageAndActualPath(path)
	if err != nil {
		return errors.WithMessage(err, "failed get storage")
	}
	return op.Remove(ctx, storage, actualPath)
}
func batchRemove(ctx context.Context, srcpath string, objs []model.IDName) error {
	storage, actualPath, err := op.GetStorageAndActualPath(srcpath)
	if err != nil {
		return errors.WithMessage(err, "failed get storage")
	}
	switch storage.(type) {
	case driver.BatchRemove:
		return op.BatchRemove(ctx, storage, actualPath, objs)
	default:
		for _, obj := range objs {
			err := op.Remove(ctx, storage, filepath.Join(actualPath, obj.Name))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func other(ctx context.Context, args model.FsOtherArgs) (interface{}, error) {
	storage, actualPath, err := op.GetStorageAndActualPath(args.Path)
	if err != nil {
		return nil, errors.WithMessage(err, "failed get storage")
	}
	args.Path = actualPath
	return op.Other(ctx, storage, args)
}

type TaskData struct {
	task.TaskExtension
	Status        string        `json:"-"` //don't save status to save space
	SrcActualPath string        `json:"src_path"`
	DstActualPath string        `json:"dst_path"`
	SrcStorage    driver.Driver `json:"-"`
	DstStorage    driver.Driver `json:"-"`
	SrcStorageMp  string        `json:"src_storage_mp"`
	DstStorageMp  string        `json:"dst_storage_mp"`
}

func (t *TaskData) GetStatus() string {
	return t.Status
}
