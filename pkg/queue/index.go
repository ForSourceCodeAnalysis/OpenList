package queue

import (
	"context"

	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/internal/db"
	"github.com/OpenListTeam/OpenList/v4/internal/model/tables"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var client *asynq.Client
var server *asynq.Server
var mux *asynq.ServeMux

func Init() {

	redisOpt := asynq.RedisClientOpt{Addr: conf.Conf.Redis.Addr,
		Password: conf.Conf.Redis.Password,
		DB:       conf.Conf.Redis.DB}

	client = asynq.NewClient(redisOpt)

	server = asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)
	mux = asynq.NewServeMux()
}

// Start starts the queue server
func Start() {
	go func() {
		if err := server.Run(mux); err != nil {
			logrus.Fatal(errors.WithStack(err))
		}
	}()
}

// GetClient returns the client
func GetClient() *asynq.Client {
	return client
}

// RegisterHandler registers a handler for a task
func RegisterHandler(typename string, Handler func(ctx context.Context, task *asynq.Task) error) {
	mux.HandleFunc(typename, Handler)
}

// AddQueue adds a task to the queue
func AddQueue(taskType string, payload []byte, options ...asynq.Option) error {
	t := asynq.NewTask(taskType, payload, options...)

	ti, err := client.Enqueue(t)
	if err != nil {
		if errors.Is(err, asynq.ErrTaskIDConflict) {
			logrus.Warnf("task id conflict,%v", ti.ID)
			return nil
		}
		logrus.Error("add queue err :", err)
		return err
	}
	db.CreateQueueItem(&tables.QueueItem{
		TaskID:   ti.ID,
		TaskType: taskType,
		Payload:  string(payload),
	})
	return nil
}
