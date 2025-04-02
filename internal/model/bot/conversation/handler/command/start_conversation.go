package command

import (
	"context"
	"fmt"

	botconvuc "bot/internal/model/bot/conversation/usecase"
	botapi "bot/internal/pkg/bot"
)

type StartConversationHandler struct {
	handler botconvuc.StartConversationHandler
}

func NewStartConversationHandler(handler botconvuc.StartConversationHandler) *StartConversationHandler {
	return &StartConversationHandler{handler: handler}
}

func (c *StartConversationHandler) Handle(ctx context.Context, im *botapi.IncomingMessage) error {
	err := c.handler.Handle(ctx, botconvuc.StartConversation{UserID: im.UserID, ChatID: im.ChatID})
	if err != nil {
		return fmt.Errorf("failed to start conversation: %w", err)
	}
	return nil
}
