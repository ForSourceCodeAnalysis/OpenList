package middlewares

import (
	"net/url"
	stdpath "path"

	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/internal/errs"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	"github.com/OpenListTeam/OpenList/v4/server/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func FsUp(c *gin.Context) {
	path := c.GetHeader("File-Path")
	password := c.GetHeader("Password")
	path, err := url.PathUnescape(path)
	if err != nil {
		common.ErrorResp(c, err, 400)
		c.Abort()
		return
	}
	user := c.Request.Context().Value(conf.UserKey).(*model.User)
	path, err = user.JoinPath(path)
	if err != nil {
		common.ErrorResp(c, err, 403)
		return
	}
	meta, err := op.GetNearestMeta(stdpath.Dir(path))
	if err != nil {
		if !errors.Is(errors.Cause(err), errs.MetaNotFound) {
			common.ErrorResp(c, err, 500, true)
			c.Abort()
			return
		}
	}
	if !(common.CanAccess(user, meta, path, password) && (user.CanWrite() || common.CanWrite(meta, stdpath.Dir(path)))) {
		common.ErrorResp(c, errs.PermissionDenied, 403)
		c.Abort()
		return
	}
	storage, actualPath, err := op.GetStorageAndActualPath(path)
	if err != nil {
		common.ErrorResp(c, err, 400)
		c.Abort()
		return
	}
	if storage.Config().NoUpload {
		common.ErrorStrResp(c, "Current storage doesn't support upload", 403)
		c.Abort()
		return
	}

	common.GinWithValue(c, conf.StorageKey, storage)
	common.GinWithValue(c, conf.PathKey, actualPath)

	c.Next()
}
