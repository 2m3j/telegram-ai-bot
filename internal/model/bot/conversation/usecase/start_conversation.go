package usecase

import (
	"context"
	"fmt"
	"log/slog"

	botapi "bot/internal/pkg/bot"
	"bot/internal/pkg/decorator"
)

type StartConversation struct {
	UserID uint64
	ChatID uint64
}

type StartConversationHandler decorator.CommandHandler[StartConversation]

type startConversationHandler struct {
	bot botapi.Bot
}

func NewStartConvHandler(bot botapi.Bot, logger *slog.Logger, metricsClient decorator.MetricsClient) StartConversationHandler {
	return decorator.ApplyCommandDecorators[StartConversation](
		&startConversationHandler{bot: bot},
		logger,
		metricsClient,
	)
}

func (h *startConversationHandler) Handle(ctx context.Context, c StartConversation) error {
	message := "Добро пожаловать в чат OpenAI, зайдайте интересующий вас вопрос. Для очистĸи ĸонтеĸста диалога используйте команду /clear."
	err := h.bot.SendMessage(ctx, c.ChatID, message)
	if err != nil {
		return fmt.Errorf("failed to send message (chat_id=%d, message='%s'): %w", c.ChatID, message, err)
	}
	return nil
}
