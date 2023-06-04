package worker

import (
	"context"
	"encoding/json"

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

/* 发布SendMail任务 */
func (distributor *TaskDistributor) DistributeTaskSendMail(
	ctx context.Context,
	payload *PayloadSendMail,
	opts ...asynq.Option,
) error {
	_payload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TaskSendMail, _payload, opts...)
	_, err = distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return err
	}

	log.Info().Msg("task enqueued")

	return nil
}
