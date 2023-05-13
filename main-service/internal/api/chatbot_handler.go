package api

import (
	"context"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Banana-Boat/terryminal/main-service/internal/pb"
)

type chatRequest struct {
	Messages []*pb.ChatRequest_ChatMessage `json:"messages" binding:"required"`
}

func (server *Server) chatHandle(ctx *gin.Context) {
	/* 设置响应头 */
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")

	/* 解析请求参数 */
	var req chatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.SSEvent("error", err.Error()) // 向客户端发送错误信息
		return
	}

	/* 通过gRPC调用chatbot-service */
	gRPCConn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", server.config.ChatbotServiceHost, server.config.ChatbotServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect grpc server")
		ctx.SSEvent("error", err.Error()) // 向客户端发送错误信息
		return
	}
	defer gRPCConn.Close()

	gRPCClient := pb.NewChatbotClient(gRPCConn)
	chatStream, err := gRPCClient.Chat(
		context.Background(),
		&pb.ChatRequest{
			Messages: req.Messages,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to call grpc server")
		ctx.SSEvent("error", err.Error()) // 向客户端发送错误信息
		return
	}

	/* 循环获取流数据，并转发 */
	ctx.Stream(func(w io.Writer) bool {
		resp, err := chatStream.Recv()

		if err != nil {
			if err == io.EOF {
				return false
			}
			ctx.SSEvent("error", err.Error())
			return false
		}

		ctx.SSEvent(resp.Event, resp.Data)
		return true
	})
}
