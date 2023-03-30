package main

import (
	"fmt"
	"net"

	"github.com/Banana-Boat/terryminal/terryminal-terminal/internal/api"
	"github.com/Banana-Boat/terryminal/terryminal-terminal/internal/pb"
	"github.com/Banana-Boat/terryminal/terryminal-terminal/internal/util"
	"github.com/Banana-Boat/terryminal/terryminal-terminal/internal/worker"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	/* 加载配置 */
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config: ")
		return
	}

	/* 创建 Redis 的 Distributor & Processor */
	redisOPt := asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
	}
	taskDistributor := worker.NewTaskDistributor(redisOPt)
	go runTaskProcessor(redisOPt) // 创建 go routine

	/* 运行 gRPC 服务 */
	runGRPCServer(config, taskDistributor)
}

func runGRPCServer(config util.Config, taskDistributor *worker.TaskDistributor) error {
	server, err := api.NewServer(config, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server: ")
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMailServiceServer(grpcServer, server)
	reflection.Register(grpcServer) // 使得grpc客户端能够了解哪些rpc调用被服务端支持，以及如何调用

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.TerminalServerHost, config.TerminalServerPort))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
		return err
	}

	log.Info().Msgf("gRPC server started at %s:%s successfully", config.TerminalServerHost, config.TerminalServerPort)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start gRPC server")
		return err
	}

	return nil
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt) {
	taskProcessor := worker.NewTaskProcessor(redisOpt)
	log.Info().Msg("task processor started successfully")

	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}
