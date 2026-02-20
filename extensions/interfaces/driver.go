package interfaces

import (
	"context"

	"github.com/OpenListTeam/OpenList/v4/extensions/models"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
)

// IBatchRename 批量重命名接口
type IBatchRename interface {
	BatchRename(ctx context.Context, srcDir model.Obj, renameObjs []models.DriverRenameObj) error
}

// IBatchRemove 批量删除接口
type IBatchRemove interface {
	BatchRemove(ctx context.Context, srcDir model.Obj, objs []model.Obj) error
}
