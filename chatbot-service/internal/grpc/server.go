package grpc

import (
	"fmt"
	"net"

	"github.com/Banana-Boat/terryminal/chatbot-service/internal/pb"
	"github.com/Banana-Boat/terryminal/chatbot-service/internal/util"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedChatbotServer
	config util.Config
}

func NewServer(config util.Config) *Server {
	server := &Server{
		config: config,
	}

	return server
}

func (server *Server) Start() error {
	grpcServer := grpc.NewServer()
	pb.RegisterChatbotServer(grpcServer, server)
	reflection.Register(grpcServer) // 使得grpc客户端能够了解哪些rpc调用被服务端支持，以及如何调用

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.config.ChatbotServerHost, server.config.ChatbotServerPort))
	if err != nil {
		return err
	}

	log.Info().Msgf("gRPC server started at %s:%s successfully", server.config.ChatbotServerHost, server.config.ChatbotServerPort)
	if err = grpcServer.Serve(listener); err != nil {
		return err
	}

	return nil
}
