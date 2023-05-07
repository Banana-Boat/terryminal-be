package http

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Banana-Boat/terryminal/bot-service/internal/util"
	"github.com/sashabaranov/go-openai"
)

type Stream struct {
	reader   *bufio.Reader  // 用于读取流数据
	response *http.Response // 用于关闭连接
}

type Api2dClient struct {
	Url    string // API2D的API地址
	ApiKey string // API2D的Forward Key
}

func NewApi2dClient(config util.Config) *Api2dClient {
	return &Api2dClient{
		Url:    config.Api2dUrl,
		ApiKey: config.Api2dKey,
	}
}

func (api2dClient *Api2dClient) CreateStream(body openai.ChatCompletionRequest) (*Stream, error) {
	/* 根据传入请求体创建请求 */
	client := &http.Client{}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", api2dClient.Url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	/* 设置请求头 */
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api2dClient.ApiKey))

	/* 发起请求，并返回stream对象 */
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(resp.Body)

	return &Stream{
		reader:   reader,
		response: resp,
	}, nil
}

func (stream *Stream) Close() {
	stream.response.Body.Close()
}

// 用于读取流数据。isValid表示是否为有效数据，api会返回空数据，忽略即可
func (stream *Stream) Recv() (resp *openai.ChatCompletionStreamResponse, isValid bool, err error) {
	line, err := stream.reader.ReadBytes('\n')
	if err != nil {
		return nil, false, err
	}

	/* 处理data前缀 */
	var headerData = []byte("data: ")
	line = bytes.TrimSpace(line)
	if !bytes.HasPrefix(line, headerData) { // 无效数据，忽略即可
		return nil, true, nil
	}
	line = bytes.TrimPrefix(line, headerData)

	/* 处理结束符 */
	if string(line) == "[DONE]" {
		return nil, false, io.EOF
	}

	/* 解析数据，并返回 */
	var respData openai.ChatCompletionStreamResponse
	err = json.Unmarshal(line, &respData)
	if err != nil {
		return nil, false, err
	}

	return &respData, false, nil
}
