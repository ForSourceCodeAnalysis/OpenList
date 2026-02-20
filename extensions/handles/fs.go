package handles

import (
	"github.com/OpenListTeam/OpenList/v4/extensions/constants"
	"github.com/OpenListTeam/OpenList/v4/extensions/models/reqres"
	"github.com/OpenListTeam/OpenList/v4/extensions/services/fs"
	"github.com/OpenListTeam/OpenList/v4/internal/driver"
	"github.com/OpenListTeam/OpenList/v4/server/common"
	"github.com/gin-gonic/gin"
)

// BatchRename 批量重命名
func BatchRename(c *gin.Context) {
	var req reqres.BatchRenameReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	storage := c.Request.Context().Value(constants.StorageKey).(driver.Driver)
	srcDir := c.Request.Context().Value(constants.SrcDirKey).(string)

	err := fs.BatchRename(c.Request.Context(), storage, srcDir, req.RenameObjects)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}

	common.SuccessResp(c)
}

// BatchRemove 批量删除
func BatchRemove(c *gin.Context) {
	var req reqres.RemoveReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if len(req.Names) == 0 {
		common.ErrorStrResp(c, "Empty file names", 400)
		return
	}

	storage := c.Request.Context().Value(constants.StorageKey).(driver.Driver)
	srcDir := c.Request.Context().Value(constants.SrcDirKey).(string)

	err := fs.BatchRemove(c.Request.Context(), storage, srcDir, req.Names)

	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}

	common.SuccessResp(c)
}
