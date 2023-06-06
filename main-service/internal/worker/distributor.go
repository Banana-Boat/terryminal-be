package worker

import (
	"context"

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

/* 发布任务 */
func (distributor *TaskDistributor) DistributeTask(
	ctx context.Context,
	typeName string,
	payload []byte,
	opts ...asynq.Option,
) error {
	task := asynq.NewTask(typeName, payload, opts...)
	taskInfo, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return err
	}

	log.Info().Msgf("task enqueued: id=%s type=%s", taskInfo.ID, task.Type())
	return nil
}
