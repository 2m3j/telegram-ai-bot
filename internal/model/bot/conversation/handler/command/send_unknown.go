package command

import (
	"context"
	"fmt"

	botconvuc "bot/internal/model/bot/conversation/usecase"
	botapi "bot/internal/pkg/bot"
)

type SendUnknownHandler struct {
	handler botconvuc.SendMessageHandler
}

func NewSendUnknownHandler(handler botconvuc.SendMessageHandler) *SendUnknownHandler {
	return &SendUnknownHandler{handler: handler}
}

func (c *SendUnknownHandler) Handle(ctx context.Context, im *botapi.IncomingMessage) error {
	err := c.handler.Handle(ctx, botconvuc.SendMessage{UserID: im.UserID, ChatID: im.ChatID, Message: "Неизвестная команда, попробуйте снова."})
	if err != nil {
		return fmt.Errorf("failed to send unknown command message: %w", err)
	}
	return nil
}
