package ai

import (
	"context"
	"fmt"
)

type PlatformName = string
type ModelName = string

const (
	PlatformOpenAI        PlatformName = "OpenAI"
	ModelOpenAIGPT4oMini  ModelName    = "gpt-4o-mini"
	ModelOpenAIGPT4o      ModelName    = "gpt-4o"
	ModelOpenAIO1Mini     ModelName    = "o1-mini"
	ModelOpenAIO3Mini     ModelName    = "o3-mini"
	PlatformDeepSeek      PlatformName = "DeepSeek"
	ModelDeepSeekChat     ModelName    = "deepseek-chat"
	ModelDeepSeekCoder    ModelName    = "deepseek-coder"
	ModelDeepSeekReasoner ModelName    = "deepseek-reasoner"
)

type PlatformProvider interface {
	Request(ctx context.Context, model string, message string, history []RequestHistory) (string, error)
}

type ErrorsHandler func(err error)

type RequestHistory struct {
	UserMessage      string
	AssistantMessage string
}

type Client struct {
	platforms    map[PlatformName]PlatformProvider
	errorHandler ErrorsHandler
}

func NewClient(options ...Option) *Client {
	c := &Client{
		platforms: make(map[PlatformName]PlatformProvider),
	}
	for _, o := range options {
		o(c)
	}
	return c
}

func (c *Client) Request(ctx context.Context, platform PlatformName, model string, message string, history []RequestHistory) (string, error) {
	var r string
	var err error
	if c.platforms[platform] != nil {
		r, err = c.platforms[platform].Request(ctx, model, message, history)
	} else {
		r = ""
		err = fmt.Errorf("platform %s not found (model=%s)", platform, model)
	}
	if err != nil {
		c.errorHandler(err)
	}
	return r, err
}
