package models

import "github.com/OpenListTeam/OpenList/v4/internal/model"

// RenameObj 重命名对象
type RenameObj struct {
	SrcName string `json:"src_name"`
	NewName string `json:"new_name"`
}

// DriverRenameObj 驱动重命名对象
type DriverRenameObj struct {
	model.Obj
	NewName string `json:"new_name"`
}

// Hash 文件的哈希值
type Hash struct {
	Md5      string `json:"md5"`
	Md5256KB string `json:"md5_256kb"`
	Sha1     string `json:"sha1"`
}
