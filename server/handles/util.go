package handles

import (
	"net/url"
	"path/filepath"

	"github.com/OpenListTeam/OpenList/v4/internal/driver"
	"github.com/OpenListTeam/OpenList/v4/internal/errs"
	"github.com/OpenListTeam/OpenList/v4/internal/fs"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	"github.com/OpenListTeam/OpenList/v4/server/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// 路径权限验证
func pathPermissionValidate(c *gin.Context, path *string, write bool) error {
	p, err := url.PathUnescape(*path)
	if err != nil {
		common.ErrorResp(c, err, 400)
		return err
	}
	user := c.MustGet("user").(*model.User)
	ap, err := user.JoinPath(p)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return err
	}
	meta, err := op.GetNearestMeta(filepath.Dir(ap))
	if err != nil && !errors.Is(errors.Cause(err), errs.MetaNotFound) {
		common.ErrorResp(c, err, 500)
		return err
	}
	password := c.GetHeader("Password")
	if !(common.CanAccess(user, meta, ap, password)) {
		return errs.PermissionDenied
	}

	if write && (user.CanWrite() || common.CanWrite(meta, filepath.Dir(ap))) {
		return errs.PermissionDenied
	}
	path = &ap
	return nil
}

func stroageNouploadCheck(c *gin.Context, path string) (driver.Driver, int, error) {
	storage, err := fs.GetStorage(path, &fs.GetStoragesArgs{})
	if err != nil {

		return nil, 400, err
	}
	if storage.Config().NoUpload {
		return nil, 405, errors.New("Current storage doesn't support upload")
	}
	return storage, 0, nil
}
