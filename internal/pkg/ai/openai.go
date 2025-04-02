package ai

import (
	"context"

	openaiapi "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIProvider struct {
	client *openaiapi.Client
}

func NewOpenAIProvider(token string) *OpenAIProvider {
	var options = []option.RequestOption{option.WithAPIKey(token)}
	client := openaiapi.NewClient(options...)
	return &OpenAIProvider{client: client}
}

func (c *OpenAIProvider) Request(ctx context.Context, model string, message string, history []RequestHistory) (string, error) {
	var chatCompletion *openaiapi.ChatCompletion
	var err error
	if len(history) == 0 {
		chatCompletion, err = c.client.Chat.Completions.New(ctx, openaiapi.ChatCompletionNewParams{
			Messages: openaiapi.F([]openaiapi.ChatCompletionMessageParamUnion{
				openaiapi.UserMessage(message),
			}),
			Model: openaiapi.F(model),
		})
	} else {
		chatMessages := make([]openaiapi.ChatCompletionMessageParamUnion, 2*len(history)+1)
		var chatMessageIndex = 0
		for _, request := range history {
			chatMessages[chatMessageIndex] = openaiapi.UserMessage(request.UserMessage)
			chatMessageIndex++
			chatMessages[chatMessageIndex] = openaiapi.AssistantMessage(request.AssistantMessage)
			chatMessageIndex++
		}
		chatMessages[chatMessageIndex] = openaiapi.UserMessage(message)
		chatCompletion, err = c.client.Chat.Completions.New(ctx, openaiapi.ChatCompletionNewParams{
			Messages: openaiapi.F(chatMessages),
			Model:    openaiapi.F(model),
		})
	}
	if err != nil {
		return "", err
	}
	return chatCompletion.Choices[0].Message.Content, nil
}
