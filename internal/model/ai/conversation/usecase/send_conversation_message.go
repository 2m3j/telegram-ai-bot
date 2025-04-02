package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"bot/internal/model/ai/conversation/entity"
	"bot/internal/model/ai/conversation/repository"
	aiapi "bot/internal/pkg/ai"
	botapi "bot/internal/pkg/bot"
	"bot/internal/pkg/decorator"
)

type SendConversationMessage struct {
	ConversationID entity.ConversationID
	MessageID      entity.MessageID
	UserID         uint64
	ChatID         uint64
	UserMessage    string
	AIPlatform     string
	AIModel        string
}

func (sc SendConversationMessage) GetChatID() uint64 {
	return sc.ChatID
}

type SendConversationMessageHandler decorator.CommandHandler[SendConversationMessage]

type sendConversationMessageHandler struct {
	convRepository repository.ConversationRepository
	msgRepository  repository.MessageRepository
	bot            botapi.Bot
	ai             *aiapi.Client
}

func NewSendConversationMessageHandler(
	convRepository repository.ConversationRepository,
	msgRepository repository.MessageRepository,
	bot botapi.Bot,
	ai *aiapi.Client,
	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) SendConversationMessageHandler {
	h := &sendConversationMessageHandler{convRepository, msgRepository, bot, ai}
	return decorator.ApplyCommandDecorators[SendConversationMessage](
		decorator.ApplyCommandTypingDecorator[SendConversationMessage](h, bot),
		logger,
		metricsClient,
	)
}

func (h *sendConversationMessageHandler) Handle(ctx context.Context, c SendConversationMessage) error {
	foundConv, err := h.convRepository.FindById(ctx, c.ConversationID)
	if err != nil {
		return fmt.Errorf("failed to find conversation with ID %s: %w", c.ConversationID.String(), err)
	}
	if foundConv == nil {
		return fmt.Errorf("conversation with ID %s not found", c.ConversationID.String())
	}
	var history []*entity.Message
	history, err = h.msgRepository.FindByCriteria(
		ctx,
		repository.NewMessageCriteria().WithConversationID(foundConv.ID).WithStatus(entity.StatusSuccess),
		repository.NewMessageSort().WithID(true),
		0,
		0,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to fetch conversation message history with ID %s and status %s: %w",
			c.ConversationID.String(),
			entity.StatusSuccess,
			err,
		)
	}

	requestHistory := make([]aiapi.RequestHistory, len(history))
	for hIndex, hItem := range history {
		requestHistory[hIndex] = aiapi.RequestHistory{UserMessage: hItem.UserMessage, AssistantMessage: hItem.AssistantMessage}
	}
	m := entity.NewMessage(c.MessageID, c.ConversationID, c.UserMessage, c.AIPlatform, c.AIModel)
	err = h.msgRepository.Add(ctx, m)
	if err != nil {
		return fmt.Errorf("failed to add new message %+v: %w", m, err)
	}

	var msgAssistantResp, msgStatus, botMsg string
	var responseError error
	msgAssistantResp, responseError = h.ai.Request(ctx, c.AIPlatform, c.AIModel, c.UserMessage, requestHistory)
	if responseError != nil {
		msgAssistantResp = responseError.Error()
		botMsg = "К сожалению, я не могу ответить сейчас на вопрос, повторите позже."
		msgStatus = entity.StatusError
	} else {
		botMsg = msgAssistantResp
		msgStatus = entity.StatusSuccess
	}
	m.Change(msgAssistantResp, msgStatus)
	if err = h.msgRepository.Update(ctx, m); err != nil {
		return fmt.Errorf("failing to update message status to '%s' (message_id=%s): %w", msgStatus, m.ID.String(), err)
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
