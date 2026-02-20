package teldrive

import (
	"net/http"

	"github.com/OpenListTeam/OpenList/v4/drivers/base"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/go-resty/resty/v2"
	"golang.org/x/net/context"
)

// BatchRemove 批量删除
func (d *Teldrive) BatchRemove(ctx context.Context, srcDir model.Obj, objs []model.Obj) error {
	ids := []string{}
	for _, obj := range objs {
		ids = append(ids, obj.GetID())
	}
	body := base.Json{
		"ids": ids,
	}
	return d.request(http.MethodPost, "/api/files/delete", func(req *resty.Request) {
		req.SetBody(body)
	}, nil)
}
