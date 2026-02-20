package middlewares

import (
	"net/url"

	"github.com/OpenListTeam/OpenList/v4/extensions/constants"
	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/internal/errs"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	"github.com/OpenListTeam/OpenList/v4/server/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// 文件操作相关中间件
// 把文件路径校验及鉴权统一放在这里处理，并将处理结果放到上下文里，避免后续重复处理
// 使用这个中间件的前提是，请求头中必须携带 Src-Dir 字段，表示操作文件的目录
// 一般来讲，文件相关操作都是需要路径的，所以这个文件头是必须的
// 主要包括：
// 1. 用户实际路径的转换
// 2. 读写权限判定
// 3. 获取存储驱动及网盘实际路径

type permissionFunc func(user *model.User, meta *model.Meta, srcDir string, password string) bool

// BatchRename 批量重命名中间件
func BatchRename(c *gin.Context) {
	fs(c, func(user *model.User, meta *model.Meta, srcDir string, password string) bool {
		return user.CanRename()
	})
}

// BatchRemove 批量删除中间件
func BatchRemove(c *gin.Context) {
	fs(c, func(user *model.User, meta *model.Meta, srcDir string, password string) bool {
		return user.CanRemove()
	})
}

// 文件夹操作中间件，用户解析，源文件夹解析，权限校验
func fs(c *gin.Context, permission permissionFunc) {
	// 解析源文件夹
	srcDir := c.GetHeader("Src-Dir")
	password := c.GetHeader("Password")
	cleanSrcDir, err := url.PathUnescape(srcDir)
	if err != nil {
		common.ErrorResp(c, err, 400)
		c.Abort()
		return
	}
	//用户
	user := c.Request.Context().Value(conf.UserKey).(*model.User)
	userSrcDir, err := user.JoinPath(cleanSrcDir)
	if err != nil {
		common.ErrorResp(c, err, 403)
		c.Abort()
		return
	}
	meta, err := op.GetNearestMeta(userSrcDir)
	if err != nil {
		if !errors.Is(errors.Cause(err), errs.MetaNotFound) {
			common.ErrorResp(c, err, 500, true)
			c.Abort()
			return
		}
	}
	//权限校验
	if !permission(user, meta, userSrcDir, password) {
		common.ErrorResp(c, errs.PermissionDenied, 403)
		c.Abort()
		return
	}
	// 解析存储驱动及实际路径
	storage, actualPath, err := op.GetStorageAndActualPath(userSrcDir)
	if err != nil {
		common.ErrorResp(c, err, 400)
		c.Abort()
		return
	}
	// if storage.Config().NoUpload {
	// 	common.ErrorStrResp(c, "Current storage doesn't support upload", 403)
	// 	c.Abort()
	// 	return
	// }
	common.GinWithValue(c, constants.StorageKey, storage)
	common.GinWithValue(c, constants.SrcDirKey, actualPath) //这里的路径已经是网盘真实路径了

	c.Next()
}
