package driver

import (
	// open123 "github.com/OpenListTeam/OpenList/v4/drivers/123_open"
	"github.com/Xhofe/go-cache"
)

// DriverCache 缓存
var DriverCache = cache.NewMemCache(cache.WithShards[any](135))
