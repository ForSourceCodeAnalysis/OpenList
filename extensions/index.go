package extensions

import (
	"github.com/OpenListTeam/OpenList/v4/extensions/backup"
	"github.com/OpenListTeam/OpenList/v4/extensions/queue"
	"github.com/gin-gonic/gin"
)

// RegisterRoute register extension routes
func RegisterRoute(g map[string]*gin.RouterGroup) {
	backup.Route(g["backup"])
}

// Init extension
func Init() {
	// 在使用队列相关的操作前，要确保队列已经初始化了
	queue.Init()

	backup.Init()

	// 启动队列前，需要确保已经注册了任务处理函数，queue.RegisterHandler()
	queue.Start()
}
