package command

import (
	"context"
	"fmt"

	aiconvuc "bot/internal/model/ai/conversation/usecase"
	botapi "bot/internal/pkg/bot"
)

type ClearConversationContextHandler struct {
	handler aiconvuc.EndConversationsHandler
}

func NewClearConversationContextHandler(handler aiconvuc.EndConversationsHandler) *ClearConversationContextHandler {
	return &ClearConversationContextHandler{handler: handler}
}

func (c *ClearConversationContextHandler) Handle(ctx context.Context, im *botapi.IncomingMessage) error {
	err := c.handler.Handle(ctx, aiconvuc.EndConversations{UserID: im.UserID, ChatID: im.ChatID})
	if err != nil {
		return fmt.Errorf("failed to clear conversation context: %w", err)
	}
	return nil
}
