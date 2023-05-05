package main

import (
	"fmt"

	"github.com/Banana-Boat/terryminal/bot-service/internal/http"
	"github.com/Banana-Boat/terryminal/bot-service/internal/util"
	"github.com/Banana-Boat/terryminal/bot-service/internal/worker"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func main() {
	/* 加载配置 */
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
		return
	}

	/* 创建 Redis 的 Distributor & Processor */
	redisOPt := asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
	}
	taskDistributor := worker.NewTaskDistributor(redisOPt)
	go runTaskProcessor(redisOPt)

	/* 运行 http 服务 */
	runHttpServer(config, taskDistributor)
}

func runHttpServer(config util.Config, taskDistributor *worker.TaskDistributor) {
	server := http.NewServer(config, taskDistributor)

	if err := server.Start(); err != nil {
		log.Error().Err(err).Msg("cannot start server")
	}
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt) {
	taskProcessor := worker.NewTaskProcessor(redisOpt)
	log.Info().Msg("task processor started successfully")

	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}
