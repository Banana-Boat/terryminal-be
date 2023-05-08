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

const maxMsgNum = 5 // 历史消息条数上限
// 系统用户prompt
const systemPrompt = `你是Terryminal平台的机器人，你叫Terry，你的职责是为用户解答有关Linux命令方面的疑问。
当用户首次向你问好时，请简短地介绍你自己。`

func (server *Server) chat(ctx *gin.Context) {
	msgBuf := "" // 消息缓冲区

	/* 解析请求参数 */
	var req chatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.SSEvent("error", err.Error()) // 向客户端发送错误信息
		return
	}
	messages := req.Messages
	// 限制历史消息条数，以减少token开销
	if len(messages) >= maxMsgNum {
		messages = messages[len(messages)-maxMsgNum:]
	}
	// 添加系统角色的prompt到消息列表开头
	_messages := []openai.ChatCompletionMessage{{
		Role:    openai.ChatMessageRoleSystem,
		Content: systemPrompt,
	}}
	messages = append(
		_messages,
		messages...,
	)
	log.Info().Msgf("Chat request: %+v", messages)

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
				/* 计算token开销 */
				log.Info().Msg("Stream read finished, message: " + msgBuf)
				_messages := append(
					messages,
					openai.ChatCompletionMessage{
						Content: msgBuf,
						Role:    openai.ChatMessageRoleAssistant,
					},
				)
				tokens, err := CalTokenCost(_messages, openai.GPT3Dot5Turbo)
				if err != nil {
					log.Err(err).Msg("CalTokenCost failed")
				} else {
					log.Info().Msgf("Token cost: %d", tokens)
				}

				/* 向客户端发送结束信息 */
				ctx.SSEvent("end", "")
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
		msgBuf += resp.Choices[0].Delta.Content
		ctx.SSEvent(
			"message",
			gin.H{
				"content": resp.Choices[0].Delta.Content,
			},
		)
		return true
	})

}
