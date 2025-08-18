package bootstrap

import (
	queuehandler "github.com/OpenListTeam/OpenList/v4/internal/queue_handler"
	"github.com/OpenListTeam/OpenList/v4/pkg/queue"
)

// InitQueue 初始化队列
func InitQueue() {
	queue.Init()
	queue.RegisterHandler(queue.TaskTypeProxyUp, queuehandler.QueueHandlerProxyUpload)
	// queue.RegisterHandler(queue.TaskTypeMove, queuehandler.QHMoveCopyTask)
	// queue.RegisterHandler(queue.TaskTypeCopy, queuehandler.QHMoveCopyTask)
}
