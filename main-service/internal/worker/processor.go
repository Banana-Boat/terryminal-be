package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Banana-Boat/terryminal/main-service/internal/util"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type TaskProcessor struct {
	server *asynq.Server
	config util.Config
}

func NewTaskProcessor(config util.Config) *TaskProcessor {
	redisOpt := asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
	}
	server := asynq.NewServer(redisOpt, asynq.Config{})

	return &TaskProcessor{
		server: server,
		config: config,
	}
}

func (processor *TaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendMail, processor.ProcessTaskSendMail)

	return processor.server.Start(mux)
}

/* 执行SendMail任务 */
func (processor *TaskProcessor) ProcessTaskSendMail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendMail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	/* 待实现email发送逻辑 */

	log.Info().Msg("task processed")

	return nil
}
