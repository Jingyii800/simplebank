package worker

import (
	"context"
	"encoding/json"
	"fmt"

	db "github.com/Jingyii800/simplebank/db/sqlc"
	"github.com/Jingyii800/simplebank/util"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type PayloadSendVerifyEmail struct {
	// username is enough for the worker to retrieve all info from db
	Username string `json:"username"`
}

const TaskSendVerifyEmail = "task:send_verify_email"

// implement the func from interface
func (distributor *RedisTaskDistributor) DistributeTaskSendverifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	// serialize object version payload to json
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	// create a new task
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)

	// send to Redis queue
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enque task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")

	return nil
}

// implement from interface
func (processor *RedisTaskProcessor) ProcessTaskVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	// retrieve from database
	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		// No need to check the user is exist or not, because it will automatically
		// retry and run out of the max times and then return fail to get user
		// if err == sql.ErrNoRows {
		// 	return fmt.Errorf("user doesn't exist: %w", asynq.SkipRetry)
		// }
		return fmt.Errorf("failed to get user: %w", err)
	}

	// create an email
	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}
	// send email to user
	subject := "Welcome to SimpleBank"
	verifyURL := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=%d&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s, <br/>
	Thank you for registering with us! <br/>
	Please <a href = "%s">Click here</a> to verify your email. <br/>`, user.FullName, verifyURL)
	to := []string{user.Email}

	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("Failed to send verify email: %w", err)
	}

	// log info
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("enail", user.Email).Msg("processed task")

	return nil

}
