package worker

import (
	"fmt"

	"github.com/Banana-Boat/terryminal/main-service/internal/util"
	"github.com/hibiken/asynq"
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
	mux.HandleFunc(TaskSendMail, processor.processTaskSendMail)

	return processor.server.Start(mux)
}
