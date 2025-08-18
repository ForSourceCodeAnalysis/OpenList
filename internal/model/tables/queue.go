package tables

import "time"

const (
	//QueueItemStatusWaiting 队列等待中
	QueueItemStatusWaiting = iota
	//QueueItemStatusRunning 队列运行中
	QueueItemStatusRunning
	//QueueItemStatusSuccess 队列成功
	QueueItemStatusSuccess
	//QueueItemStatusFailed 队列失败
	QueueItemStatusFailed
	//QueueItemStatusCancelled 队列取消
	QueueItemStatusCancelled
)

// QueueItem 队列条目信息
type QueueItem struct {
	Base
	TaskID   string    `json:"task_id" gorm:"index:idx_task_id"`
	TaskType string    `json:"task_type"`
	Payload  string    `json:"payload"`
	Status   uint      `json:"status"`
	Message  string    `json:"message"`
	StartAt  time.Time `json:"start_at"`
	EndAt    time.Time `json:"end_at"`
}
