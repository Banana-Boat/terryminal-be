package http

import (
	"fmt"

	"github.com/Banana-Boat/terryminal/bot-service/internal/util"
	"github.com/Banana-Boat/terryminal/bot-service/internal/worker"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config          util.Config
	taskDistributor *worker.TaskDistributor
	router          *gin.Engine
}

func NewServer(config util.Config, taskDistributor *worker.TaskDistributor) *Server {
	server := &Server{
		config:          config,
		taskDistributor: taskDistributor,
	}

	server.setupRouter()

	return server
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/chat", server.chat)

	server.router = router
}

func (server *Server) Start() error {
	if err := server.router.Run(
		fmt.Sprintf("%s:%s", server.config.BotHttpServerHost, server.config.BotHttpServerPort),
	); err != nil {
		return err
	}

	return nil
}

func wrapResponse(flag bool, msg string, data gin.H) gin.H {
	var _msg string
	if msg == "" {
		_msg = "OK"
	} else {
		_msg = msg
	}

	return gin.H{
		"flag": flag,
		"msg":  _msg,
		"data": data,
	}
}

// func (server *Server) SendMail(ctx context.Context, req *pb.SendMailRequest) (*pb.SendMailResponse, error) {
// 	payload := &worker.PayloadSendMail{
// 		DestAddr: req.GetDestAddr(),
// 		Content:  req.GetContent(),
// 	}

// 	err := server.taskDistributor.DistributeTaskSendMail(ctx, payload, asynq.ProcessIn(10*time.Second))
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "task distribution failed")
// 	}

// 	resp := &pb.SendMailResponse{
// 		CreatedAt: timestamppb.Now(),
// 	}
// 	return resp, nil
// }
