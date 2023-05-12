package main

import (
	"os"

	"github.com/Banana-Boat/terryminal/chatbot-service/internal/grpc"
	"github.com/Banana-Boat/terryminal/chatbot-service/internal/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// 美化zerolog
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	/* 加载配置 */
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
		return
	}

	/* 运行 gRPC 服务 */
	runGRPCServer(config)
}

func runGRPCServer(config util.Config) {
	server := grpc.NewServer(config)

	if err := server.Start(); err != nil {
		log.Error().Err(err).Msg("failed to start gRPC server")
	}
}
