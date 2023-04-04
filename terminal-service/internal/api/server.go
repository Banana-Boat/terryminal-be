package api

import (
	"context"
	"time"

	"github.com/Banana-Boat/terryminal/terminal-service/internal/pb"
	"github.com/Banana-Boat/terryminal/terminal-service/internal/util"
	"github.com/Banana-Boat/terryminal/terminal-service/internal/worker"
	"github.com/hibiken/asynq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedMailServiceServer // 使得还未被实现的rpc能够被接受
	config                            util.Config
	taskDistributor                   *worker.TaskDistributor
}

func NewServer(config util.Config, taskDistributor *worker.TaskDistributor) (*Server, error) {

	server := &Server{
		config:          config,
		taskDistributor: taskDistributor,
	}

	return server, nil
}

func (server *Server) SendMail(ctx context.Context, req *pb.SendMailRequest) (*pb.SendMailResponse, error) {
	payload := &worker.PayloadSendMail{
		DestAddr: req.GetDestAddr(),
		Content:  req.GetContent(),
	}

	err := server.taskDistributor.DistributeTaskSendMail(ctx, payload, asynq.ProcessIn(10*time.Second))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "task distribution failed")
	}

	resp := &pb.SendMailResponse{
		CreatedAt: timestamppb.Now(),
	}
	return resp, nil
}
