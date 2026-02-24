// simple-bank/worker/task_send_verify_email.go
package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (d *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(
		TaskSendVerifyEmail,
		jsonPayload,
		opts...,
	)

	info, err := d.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	// optional log
	fmt.Printf(
		"enqueued task: type=%s queue=%s max_retry=%d\n",
		task.Type(),
		info.Queue,
		info.MaxRetry,
	)

	return nil
}

func (p *RedisTaskProcessor) ProcessSendVerifyEmail(
	ctx context.Context,
	task *asynq.Task,
) error {

	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		// payload is invalid → retrying will not help
		return fmt.Errorf("unmarshal payload failed: %w", asynq.SkipRetry)
	}

	user, err := p.store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			// user does not exist → permanent error
			return fmt.Errorf("user not found: %w", asynq.SkipRetry)
		}
		// DB temporary error → retry
		return err
	}

	// TODO: generate verify token & save to DB
	// TODO: send verification email

	log.Info().
		Str("task", task.Type()).
		Str("username", user.Username).
		Str("email", user.Email).
		Msg("processed send verify email task")

	return nil
}
