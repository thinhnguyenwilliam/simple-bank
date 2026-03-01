// simple-bank/worker/processor.go
package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
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
			Logger:      NewLogger(),
			Queues: map[string]int{
				"critical": 6,
				"default":  4,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(
				func(ctx context.Context, task *asynq.Task, err error) {
					log.Error().
						Err(err).
						Str("type", task.Type()).
						RawJSON("payload", task.Payload()).
						Msg("process task failed")
				},
			),
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
