package worker

import (
	"context"

	db "github.com/Jingyii800/simplebank/db/sqlc"
	"github.com/Jingyii800/simplebank/mail"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	// register the task
	Start() error
	// handler has 2 inputs and an error output
	ProcessTaskVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	// need server
	server *asynq.Server
	// access to db
	store db.Store
	// send email
	mailer mail.EmailSender
}

// func to create a new processor
func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) TaskProcessor {
	// force redis to follow this log format
	logger := NewLogger()
	redis.SetLogger(logger)

	// can leave the config empty now to follow the pre-option
	server := asynq.NewServer(redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			// print more detailed error
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().Err(err).Str("type", task.Type()).
					Bytes("payload", task.Payload()).Msg("Process task failed")
			}),
			// print logs
			Logger: logger,
		})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
		mailer: mailer,
	}
}

// register the task to know how to process it
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	//             this is the pattern from distributor
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskVerifyEmail)

	return processor.server.Start(mux)
}
