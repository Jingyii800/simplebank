package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	// handler
	DistributeTaskSendverifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	// a new client
	client := asynq.NewClient(redisOpt)
	// a pointer to the object
	// to force the RedisTaskDistributor to implement the TaskDistributor interface
	return &RedisTaskDistributor{
		client: client,
	}
}
