package http

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
)

type chatRequest struct {
	Messages []openai.ChatCompletionMessage `json:"messages" binding:"required"`
}

func (server *Server) chat(ctx *gin.Context) {
	/* 解析请求参数 */
	var req chatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.SSEvent("error", err.Error()) // 向客户端发送错误信息
		return
	}
	messages := req.Messages

	/* 设置响应头 */
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")

	/* 创建api2dClient，发送请求 */
	api2dClient := NewApi2dClient(server.config)
	stream, err := api2dClient.CreateStream(messages)
	if err != nil {
		log.Err(err).Msg("CreateStream failed")
		ctx.SSEvent("error", err.Error()) // 向客户端发送错误信息
	}
	defer stream.Close()

	/* 循环获取流数据，并转发 */
	ctx.Stream(func(w io.Writer) bool {
		resp, isValid, err := stream.Recv() // 从流中读取数据
		if err != nil {
			/* 流数据读取完毕 */
			if err == io.EOF {
				log.Info().Msg("Stream read finished")
				ctx.SSEvent("end", "") // 向客户端发送结束信息
				return false
			}
			/* 流数据读取错误 */
			log.Err(err).Msg("Stream read failed")
			ctx.SSEvent("error", err.Error()) // 向客户端发送错误信息
			return false
		}
		/* 无效数据忽略即可，继续读取 */
		if isValid {
			return true
		}

		/* 向客户端发送流数据 */
		ctx.SSEvent(
			"message",
			gin.H{
				"content": resp.Choices[0].Delta.Content,
			},
		)
		return true
	})

}
