// simple-bank/worker/processor.go
package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/thinhcompany/simple-bank/db/sqlc"
)

// TaskProcessor handle tasks background
type TaskProcessor interface {
	Start() error
	ProcessSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

// RedisTaskProcessor implements TaskProcessor
type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

// NewRedisTaskProcessor create worker server
func NewRedisTaskProcessor(
	redisOpt asynq.RedisClientOpt,
	store db.Store,
) TaskProcessor {

	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  4,
			},
		},
	)

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

// Start run worker
func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(
		TaskSendVerifyEmail,
		p.ProcessSendVerifyEmail,
	)

	return p.server.Run(mux)
}
