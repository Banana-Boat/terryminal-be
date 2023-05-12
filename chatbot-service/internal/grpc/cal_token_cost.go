package grpc

import (
	"fmt"

	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
)

func CalTokenCost(messages []openai.ChatCompletionMessage, model string) (int, error) {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		return 0, fmt.Errorf("model %s not supported", model)
	}

	num_tokens := 0
	var tokens_per_message int
	var tokens_per_name int

	if model == openai.GPT3Dot5Turbo0301 || model == openai.GPT3Dot5Turbo {
		tokens_per_message = 4
		tokens_per_name = -1
	} else if model == openai.GPT40314 || model == openai.GPT4 {
		tokens_per_message = 3
		tokens_per_name = 1
	} else {
		fmt.Println("Warning: model not found. Using cl100k_base encoding.")
		tokens_per_message = 3
		tokens_per_name = 1
	}

	for _, message := range messages {
		num_tokens += tokens_per_message

		num_tokens += len(tkm.Encode(message.Content, nil, nil))
		num_tokens += len(tkm.Encode(message.Role, nil, nil))
		num_tokens += len(tkm.Encode(message.Name, nil, nil))
		if message.Name != "" {
			num_tokens += tokens_per_name
		}
	}
	num_tokens += 3

	return num_tokens, nil
}
