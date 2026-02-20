package reqres

import "github.com/OpenListTeam/OpenList/v4/extensions/models"

// BatchRenameReq 批量重命名请求
type BatchRenameReq struct {
	RenameObjects []models.RenameObj `json:"rename_objects"`
}

// RemoveReq 删除请求
type RemoveReq struct {
	Names []string `json:"names"`
}
