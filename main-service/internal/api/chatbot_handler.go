package api

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Banana-Boat/terryminal/main-service/internal/db"
	"github.com/Banana-Boat/terryminal/main-service/internal/pb"
)

type chatRequest struct {
	Messages []*pb.ChatRequest_ChatMessage `json:"messages" binding:"required"`
}

func (server *Server) handleChat(ctx *gin.Context) {
	tokenPayload := ctx.MustGet("token_payload").(*TokenPayload)

	/* 设置响应头 */
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")

	/* 判断用户token是否充足 */
	user, err := server.store.GetUserById(ctx, tokenPayload.ID)
	if err != nil {
		log.Error().Err(err).Msg("user not found")
		ctx.SSEvent("error", "对话失败")
		return
	}
	if user.ChatbotToken < 100 {
		log.Info().Msg("user's chatbot token is not enough")
		ctx.SSEvent("error", "Token不足100，请在个人信息页面充值后使用")
		return
	}

	/* 解析请求参数 */
	var req chatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("wrong request params")
		ctx.SSEvent("error", "参数不合法")
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
		ctx.SSEvent("error", "对话失败")
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
		ctx.SSEvent("error", "对话失败")
		return
	}

	/* 循环获取流数据，并转发 */
	ctx.Stream(func(w io.Writer) bool {
		resp, err := chatStream.Recv()

		if err != nil {
			if err == io.EOF {
				return false
			}
			log.Error().Err(err).Msg("failed to read stream")
			ctx.SSEvent("error", "对话失败")
			return false
		}

		/* 更新用户token数 */
		if resp.Event == "token" {
			user, err := server.store.GetUserById(ctx, tokenPayload.ID)
			if err != nil {
				log.Error().Err(err).Msg("user not found")
				return true
			}

			num, err := strconv.ParseInt(resp.Data, 10, 32)
			if err != nil {
				log.Error().Err(err).Msg("parse token num failed")
				return true
			}
			tokens := int32(num)
			if user.ChatbotToken < tokens {
				tokens = 0
			} else {
				tokens = user.ChatbotToken - tokens
			}

			arg := db.UpdateUserInfoParams{
				ID:           user.ID,
				Nickname:     user.Nickname,
				Password:     user.Password,
				ChatbotToken: tokens,
				UpdatedAt:    time.Now(),
			}
			err = server.store.UpdateUserInfo(ctx, arg)
			if err != nil {
				log.Error().Err(err).Msg("update user's chatbot token failed")
			}

			return true
		}

		ctx.SSEvent(resp.Event, resp.Data)

		return true
	})
}
