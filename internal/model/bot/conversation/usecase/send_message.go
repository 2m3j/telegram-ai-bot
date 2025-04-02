package usecase

import (
	"context"
	"fmt"
	"log/slog"

	botapi "bot/internal/pkg/bot"
	"bot/internal/pkg/decorator"
)

type SendMessage struct {
	UserID  uint64
	ChatID  uint64
	Message string
}

type SendMessageHandler decorator.CommandHandler[SendMessage]

type sendMessageHandler struct {
	bot botapi.Bot
}

func NewSendMessageHandler(bot botapi.Bot, logger *slog.Logger, metricsClient decorator.MetricsClient) SendMessageHandler {
	return decorator.ApplyCommandDecorators[SendMessage](
		&sendMessageHandler{bot: bot},
		logger,
		metricsClient,
	)
}

func (h *sendMessageHandler) Handle(ctx context.Context, c SendMessage) error {
	err := h.bot.SendMessage(ctx, c.ChatID, c.Message)
	if err != nil {
		return fmt.Errorf("failed to send message (chat_id=%d, message='%s'): %w", c.ChatID, c.Message, err)
	}
	return nil
}
