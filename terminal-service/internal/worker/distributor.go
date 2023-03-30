package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type TaskDistributor struct {
	client *asynq.Client
}

func NewTaskDistributor(redisOpt asynq.RedisClientOpt) *TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &TaskDistributor{
		client: client,
	}
}

type PayloadSendMail struct {
	DestAddr string `json:"destAddr"`
	Content  string `json:"content"`
}

func (distributor *TaskDistributor) DistributeTaskSendMail(
	ctx context.Context,
	payload *PayloadSendMail,
	opts ...asynq.Option,
) error {
	_payload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendMail, _payload, opts...)
	_, err = distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("enqueued task")

	return nil
}
