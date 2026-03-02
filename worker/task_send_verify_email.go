// simple-bank/worker/task_send_verify_email.go
package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	db "github.com/thinhcompany/simple-bank/db/sqlc"
	"github.com/thinhcompany/simple-bank/util"
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

	// Create verify email record
	verifyEmail, err := p.store.CreateVerifyEmail(
		ctx,
		db.CreateVerifyEmailParams{
			Username:   user.Username,
			Email:      user.Email,
			SecretCode: util.RandomString(32),
			ExpiredAt:  time.Now().Add(15 * time.Minute),
		},
	)
	if err != nil {
		return err
	}

	// Build verify URL
	verifyURL := fmt.Sprintf(
		"https://example.com/verify-email?id=%d&secret_code=%s",
		verifyEmail.ID,
		verifyEmail.SecretCode,
	)

	// Email content
	subject := "Verify your email address"
	content := fmt.Sprintf(`
		<h2>Hello %s 👋</h2>
		<p>Please click the link below to verify your email:</p>
		<a href="%s">Verify Email</a>
		<p>This link will expire in 15 minutes.</p>
	`, user.FullName, verifyURL)

	to := []string{user.Email}

	// Send email
	err = p.mailer.SendEmail(
		subject,
		content,
		to,
		nil,
		nil,
		nil,
	)
	if err != nil {
		return err
	}

	log.Info().
		Str("task", task.Type()).
		Str("username", user.Username).
		Str("to", user.Email).
		Str("subject", subject).
		Msg("processed send verify email task")

	return nil
}
