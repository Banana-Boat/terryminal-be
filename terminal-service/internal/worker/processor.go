package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type TaskProcessor struct {
	server *asynq.Server
}

const TaskSendMail = "task:send_mail"

func NewTaskProcessor(redisOpt asynq.RedisClientOpt) *TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{})

	return &TaskProcessor{
		server: server,
	}
}

func (processor *TaskProcessor) ProcessTaskSendMail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendMail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("processed task")

	return nil
}

func (processor *TaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendMail, processor.ProcessTaskSendMail)

	return processor.server.Start(mux)
}
