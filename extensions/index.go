package extensions

import (
	"github.com/OpenListTeam/OpenList/v4/extensions/backup"
	// "github.com/OpenListTeam/OpenList/v4/extensions/queue"

	// "github.com/OpenListTeam/OpenList/v4/extensions/cron"
	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/gin-gonic/gin"
)

// RegisterRoute register extension routes
func RegisterRoute(g map[string]*gin.RouterGroup) {
	backup.Route(g["backup"])
	// cron.Route(g["cron"])
}

// Init extension
func Init() {

	// 如果扩展依赖队列（队列依赖redis），请放到下面的if范围内初始化
	if conf.Conf.Redis.Addr != "" {
		// 在使用队列相关的操作前，要确保队列已经初始化了
		// queue.Init()

		// 依赖队列的扩展，放在队列初始化之后
		backup.Init()

		// 启动队列前，需要确保已经注册了任务处理函数，queue.RegisterHandler()
		// queue.Start()
	}

	// 不依赖队列的扩展，可以在下面初始化
	// example.Init()

}
