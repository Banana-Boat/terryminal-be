package grpc

import "github.com/Banana-Boat/terryminal/terminal-service/internal/util"

type GRPCServer struct {
	config util.Config
}

func NewGRPCServer(config util.Config) *GRPCServer {
	server := &GRPCServer{
		config: config,
	}

	return server
}

func (server *GRPCServer) Start() error {
	return nil
}
