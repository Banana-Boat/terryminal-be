package grpc

import (
	"io"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"

	"github.com/Banana-Boat/terryminal/chatbot-service/internal/pb"
)

const maxMsgNum = 5 // 历史消息条数上限
// 系统用户prompt
const systemPrompt = `你是Terryminal平台的机器人，你叫Terry，你的职责是为用户解答有关Linux命令方面的疑问。
当用户首次向你问好时，请简短地介绍你自己。与linux无关的问题委婉地拒绝回答。`

func (server *Server) Chat(req *pb.ChatRequest, stream pb.Chatbot_ChatServer) error {
	msgBuf := "" // 消息缓冲区
	messages := preprocessMessages(req.GetMessages())
	log.Info().Msgf("Chat request: %+v", messages)

	/* 创建api2dClient，发送请求 */
	api2dClient := NewApi2dClient(server.config)
	api2dStream, err := api2dClient.CreateStream(messages)
	if err != nil {
		log.Error().Err(err).Msg("CreateStream failed")
	}
	defer api2dStream.Close()

	/* 循环获取流数据，并转发 */
	for {
		resp, isValid, err := api2dStream.Recv() // 从流中读取数据

		if err != nil {
			/* 流数据读取完毕 */
			if err == io.EOF {
				// 计算token开销
				log.Info().Msg("Stream read finished, message: " + msgBuf)
				messages := append(
					messages,
					openai.ChatCompletionMessage{
						Content: msgBuf,
						Role:    openai.ChatMessageRoleAssistant,
					},
				)
				tokens, err := CalTokenCost(messages, openai.GPT3Dot5Turbo)
				if err != nil {
					log.Error().Err(err).Msg("CalTokenCost failed")
				} else {
					log.Info().Msgf("Token cost: %d", tokens)
				}

				// 向客户端发送结束信息，退出循环
				stream.Send(&pb.ChatResponse{Event: "end", Data: ""})
				break
			}

			/* 流数据读取错误 */
			log.Error().Err(err).Msg("Stream read failed")
			stream.Send(&pb.ChatResponse{Event: "error", Data: err.Error()})
			break
		}

		/* 无效数据忽略即可，继续读取 */
		if isValid {
			continue
		}

		/* 向客户端发送流数据 */
		msgBuf += resp.Choices[0].Delta.Content
		stream.Send(&pb.ChatResponse{
			Event: "message",
			Data:  resp.Choices[0].Delta.Content,
		})
		continue
	}

	return nil
}

func preprocessMessages(messages []*pb.ChatRequest_ChatMessage) []openai.ChatCompletionMessage {
	// 限制历史消息条数，以减少token开销
	if len(messages) >= maxMsgNum {
		messages = messages[len(messages)-maxMsgNum:]
	}
	// 转换数组类型
	_messages := []openai.ChatCompletionMessage{}
	for _, msg := range messages {
		_messages = append(_messages,
			openai.ChatCompletionMessage{
				Role:    msg.Role,
				Content: msg.Content,
			},
		)
	}
	// 添加系统角色的prompt到消息列表开头
	head_messages := []openai.ChatCompletionMessage{{
		Role:    openai.ChatMessageRoleSystem,
		Content: systemPrompt,
	}}
	_messages = append(
		head_messages,
		_messages...,
	)
	return _messages
}
