package usecase

import (
	"context"
	"fmt"
	"log/slog"

	cnvrepo "bot/internal/model/ai/conversation/repository"
	botapi "bot/internal/pkg/bot"
	"bot/internal/pkg/decorator"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

type EndConversations struct {
	UserID uint64
	ChatID uint64
}

type EndConversationsHandler decorator.CommandHandler[EndConversations]

type endConversationsHandler struct {
	trm            trm.Manager
	bot            botapi.Bot
	convRepository cnvrepo.ConversationRepository
}

func NewEndConversationsHandler(
	trm trm.Manager,
	bot botapi.Bot,
	convRepository cnvrepo.ConversationRepository,
	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) EndConversationsHandler {
	return decorator.ApplyCommandDecorators[EndConversations](
		&endConversationsHandler{trm: trm, bot: bot, convRepository: convRepository},
		logger,
		metricsClient,
	)
}

func (h *endConversationsHandler) Handle(ctx context.Context, c EndConversations) error {
	list, err := h.convRepository.FindByCriteria(ctx, cnvrepo.NewConversationCriteria().WithUserID(c.UserID).WithFinished(false), cnvrepo.NewConversationSort(), 0, 0)
	if err != nil {
		return fmt.Errorf("failed to find active conversations with user ID %d: %w", c.UserID, err)
	}
	if len(list) != 0 {
		err = h.trm.Do(ctx, func(ctx context.Context) error {
			for _, conv := range list {
				conv.End()
				if upErr := h.convRepository.Update(ctx, conv); upErr != nil {
					return fmt.Errorf("failed to update conversation %+v: %w", conv, upErr)
				}
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to transaction execution: %w", err)
		}
	}
	var msg string
	if len(list) > 0 {
		msg = "Контекст диалога очищен, о чём поговорим?"
	} else {
		msg = "Очистка контекста не требуется, о чём поговорим?"
	}
	if err = h.bot.SendMessage(ctx, c.ChatID, msg); err != nil {
		return fmt.Errorf("failed to send message (chat_id=%d, message='%s'): %w", c.ChatID, msg, err)
	}
	return nil
}
