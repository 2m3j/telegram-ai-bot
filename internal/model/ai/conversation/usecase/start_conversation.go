package usecase

import (
	"context"
	"fmt"
	"log/slog"

	conventity "bot/internal/model/ai/conversation/entity"
	convrepo "bot/internal/model/ai/conversation/repository"
	aiapi "bot/internal/pkg/ai"
	botapi "bot/internal/pkg/bot"
	"bot/internal/pkg/decorator"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

type StartConversation struct {
	ConversationID conventity.ConversationID
	MessageID      conventity.MessageID
	UserID         uint64
	ChatID         uint64
	UserMessage    string
	AIPlatform     string
	AIModel        string
}

func (sc StartConversation) GetChatID() uint64 {
	return sc.ChatID
}

type StartConversationHandler decorator.CommandHandler[StartConversation]

type startConversationHandler struct {
	convRepository convrepo.ConversationRepository
	msgRepository  convrepo.MessageRepository
	trm            trm.Manager
	bot            botapi.Bot
	ai             *aiapi.Client
}

func NewStartConversationHandler(
	convRepository convrepo.ConversationRepository,
	msgRepository convrepo.MessageRepository,
	trm trm.Manager,
	bot botapi.Bot,
	ai *aiapi.Client,
	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) StartConversationHandler {
	return decorator.ApplyCommandDecorators[StartConversation](
		decorator.ApplyCommandTypingDecorator[StartConversation](
			&startConversationHandler{
				convRepository,
				msgRepository,
				trm,
				bot,
				ai,
			},
			bot,
		),
		logger,
		metricsClient,
	)
}

func (h *startConversationHandler) Handle(ctx context.Context, c StartConversation) error {

	conv := conventity.NewConversation(c.ConversationID, c.UserID)
	m := conventity.NewMessage(c.MessageID, c.ConversationID, c.UserMessage, c.AIPlatform, c.AIModel)
	var err error
	err = h.trm.Do(ctx, func(ctx context.Context) error {
		err = h.convRepository.Add(ctx, conv)
		if err != nil {
			return fmt.Errorf("failed to add new conversation %+v: %w", conv, err)
		}
		err = h.msgRepository.Add(ctx, m)
		if err != nil {
			return fmt.Errorf("failed to add new message %+v: %w", m, err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to transaction execution: %w", err)
	}

	var msgAssistantResp, msgStatus, botMsg string
	var responseError error
	msgAssistantResp, responseError = h.ai.Request(ctx, c.AIPlatform, c.AIModel, c.UserMessage, make([]aiapi.RequestHistory, 0))
	if responseError != nil {
		msgAssistantResp = responseError.Error()
		botMsg = "К сожалению, я не могу ответить сейчас на вопрос, повторите позже."
		msgStatus = conventity.StatusError
	} else {
		botMsg = msgAssistantResp
		msgStatus = conventity.StatusSuccess
	}
	m.Change(msgAssistantResp, msgStatus)
	if err = h.msgRepository.Update(ctx, m); err != nil {
		return fmt.Errorf("failing to update message status to %s (message_id=%s): %w", msgStatus, m.ID.String(), err)
	}
	err = h.bot.SendMessage(
		ctx,
		c.ChatID,
		botMsg,
	)
	if err != nil {
		return fmt.Errorf("failed to send message (chat_id=%d, message='%s'): %w", c.ChatID, botMsg, err)
	}
	return nil
}
