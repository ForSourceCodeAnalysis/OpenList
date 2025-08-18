package db

import "github.com/OpenListTeam/OpenList/v4/internal/model/tables"

// CreateQueueItem 创建队列
func CreateQueueItem(q *tables.QueueItem) error {
	return db.Create(q).Error
}
