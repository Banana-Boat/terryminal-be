package http

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

const URL string = "https://openai.api2d.net/v1/chat/completions"
const API_KEY string = "fk200509-dDh7uhiMFxHi0gNhbrNZmjiOeKxa7hko"

func TestOpenAIAPI(t *testing.T) {
	client := &http.Client{}

	data := make(map[string]interface{})
	data["model"] = "gpt-3.5-turbo"
	data["max_tokens"] = 100
	data["stream"] = true
	data["stop"] = []string{"\n"}
	data["messages"] = []map[string]string{
		{
			"role":    "user",
			"content": "你能介绍一下linux命令吗？",
		},
	}
	_data, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", URL, bytes.NewBuffer(_data))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", API_KEY))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(line))
	}
}
