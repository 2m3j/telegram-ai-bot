package handler

import (
	"context"
	"fmt"

	convent "bot/internal/model/ai/conversation/entity"
	convrepo "bot/internal/model/ai/conversation/repository"
	aiusecase "bot/internal/model/ai/conversation/usecase"
	userent "bot/internal/model/user/entity"
	botapi "bot/internal/pkg/bot"
)

type MessageHandler struct {
	aiCreateConvHandler    aiusecase.StartConversationHandler
	aiSendConvHandler      aiusecase.SendConversationMessageHandler
	conversationRepository convrepo.ConversationRepository
}

func NewMessageHandler(
	aiCreateConvHandler aiusecase.StartConversationHandler,
	aiSendConvHandler aiusecase.SendConversationMessageHandler,
	conversationRepository convrepo.ConversationRepository,
) *MessageHandler {
	return &MessageHandler{aiCreateConvHandler: aiCreateConvHandler, aiSendConvHandler: aiSendConvHandler, conversationRepository: conversationRepository}
}

func (h *MessageHandler) Handle(ctx context.Context, im *botapi.IncomingMessage, user *userent.User) error {
	var err error
	var foundConv *convent.Conversation
	foundConv, err = h.conversationRepository.FindOneByCriteria(
		ctx,
		convrepo.NewConversationCriteria().WithUserID(im.UserID).WithFinished(false),
		convrepo.NewConversationSort().WithID(false),
		0,
	)
	if err != nil {
		return fmt.Errorf("failed to find conversation for user with ID %d: %w", im.UserID, err)
	}
	var newMessageId = convent.NextMessageID()
	if foundConv == nil {
		var newConvId = convent.NextConversationID()
		err = h.aiCreateConvHandler.Handle(
			ctx,
			aiusecase.StartConversation{
				ConversationID: newConvId,
				MessageID:      newMessageId,
				UserID:         im.UserID,
				ChatID:         im.ChatID,
				UserMessage:    im.Message,
				AIPlatform:     user.AIPlatform,
				AIModel:        user.AIModel,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to create new conversation with ID '%s' for user with ID %d: %w", newConvId.String(), im.UserID, err)
		}
	} else {
		err = h.aiSendConvHandler.Handle(
			ctx,
			aiusecase.SendConversationMessage{
				ConversationID: foundConv.ID,
				MessageID:      newMessageId,
				UserID:         im.UserID,
				ChatID:         im.ChatID,
				UserMessage:    im.Message,
				AIPlatform:     user.AIPlatform,
				AIModel:        user.AIModel,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to send message to conversation with ID '%s' for user with ID '%d: %w", foundConv.ID.String(), im.UserID, err)
		}
	}

	return nil
}
