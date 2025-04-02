package ai

import (
	"context"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
)

type DeepSeekProvider struct {
	client *deepseek.Client
}

func NewDeepSeekProvider(token string) *DeepSeekProvider {
	client := deepseek.NewClient(token, "https://api.deepseek.com/")
	return &DeepSeekProvider{client: client}
}

func (c *DeepSeekProvider) Request(ctx context.Context, model string, message string, history []RequestHistory) (string, error) {
	var response *deepseek.ChatCompletionResponse
	var err error
	if len(history) == 0 {
		response, err = c.client.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
			Model: model,
			Messages: []deepseek.ChatCompletionMessage{
				{Role: constants.ChatMessageRoleUser, Content: message},
			},
		})
	} else {
		chatMessages := make([]deepseek.ChatCompletionMessage, 2*len(history)+1)
		var chatMessageIndex = 0
		for _, item := range history {
			chatMessages[chatMessageIndex] = deepseek.ChatCompletionMessage{
				Role:    constants.ChatMessageRoleUser,
				Content: item.UserMessage,
			}
			chatMessageIndex++
			chatMessages[chatMessageIndex] = deepseek.ChatCompletionMessage{
				Role:    constants.ChatMessageRoleAssistant,
				Content: item.AssistantMessage,
			}
			chatMessageIndex++
		}
		chatMessages[chatMessageIndex] = deepseek.ChatCompletionMessage{
			Role:    constants.ChatMessageRoleUser,
			Content: message,
		}
		response, err = c.client.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
			Model:    model,
			Messages: chatMessages,
		})
	}
	if err != nil {
		return "", err
	}
	return response.Choices[0].Message.Content, nil
}
