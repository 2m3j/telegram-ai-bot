package usecase_test

import (
	"context"
	"testing"
	"time"

	cnvent "bot/internal/model/ai/conversation/entity"
	aiconvrepomocks "bot/internal/model/ai/conversation/repository/mocks"
	cnvuc "bot/internal/model/ai/conversation/usecase"
	aiapi "bot/internal/pkg/ai"
	aiapimocks "bot/internal/pkg/ai/mocks"
	botapi "bot/internal/pkg/bot"
	botmocks "bot/internal/pkg/bot/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSendConversationHandler_Handle_Success(t *testing.T) {
	var userID, chatID uint64 = 1, 2
	var conversationID = cnvent.NextConversationID()
	var messageID = cnvent.NextMessageID()
	var aiPlatform = aiapi.PlatformDeepSeek
	var aiModel = aiapi.ModelDeepSeekChat
	foundConv := cnvent.NewConversation(conversationID, userID)
	var historyUserMessageText, historyAssistantMessageText, newUserMessageText, newAssistantMessageText = "Привет", "Привет, чем могу помочь?", "Как дела?", "Хорошо"
	foundMessages := []*cnvent.Message{
		{
			ID:                messageID,
			ConversationID:    foundConv.ID,
			UserMessage:       historyUserMessageText,
			AssistantPlatform: aiPlatform,
			AssistantModel:    aiModel,
			AssistantMessage:  historyAssistantMessageText,
			UpdatedAt:         time.Now(),
			CreatedAt:         time.Now(),
		},
	}
	ctx := context.Background()
	convRepo := new(aiconvrepomocks.MockConversationRepository)
	msgRepo := new(aiconvrepomocks.MockMessageRepository)
	bot := new(botmocks.MockBot)
	aiProvider := new(aiapimocks.MockPlatformProvider)
	aiClientOpts := []aiapi.Option{
		aiapi.WithErrorHandler(func(err error) {}),
		aiapi.WithPlatform(aiPlatform, aiProvider),
	}
	aiClient := aiapi.NewClient(aiClientOpts...)
	convRepo.On("FindById", ctx, conversationID).Return(foundConv, nil)
	msgRepo.On("FindByCriteria", ctx, mock.Anything, mock.Anything, uint64(0), uint64(0)).Return(foundMessages, nil)
	msgRepo.On("Add", ctx, mock.Anything).Return(nil)
	msgRepo.On("Update", ctx, mock.Anything).Return(nil)
	bot.On("SendChatAction", mock.Anything, chatID, botapi.ChatActionTyping)
	bot.On("SendMessage", mock.Anything, chatID, newAssistantMessageText).Return(nil)
	aiProvider.On("Request", ctx, aiModel, newUserMessageText, []aiapi.RequestHistory{{UserMessage: historyUserMessageText, AssistantMessage: historyAssistantMessageText}}).Return(newAssistantMessageText, nil)
	handler := cnvuc.NewSendConversationMessageHandler(convRepo, msgRepo, bot, aiClient, nil, nil)

	err := handler.Handle(
		ctx,
		cnvuc.SendConversationMessage{
			ConversationID: conversationID,
			MessageID:      messageID,
			UserID:         userID,
			ChatID:         chatID,
			UserMessage:    newUserMessageText,
			AIPlatform:     aiPlatform,
			AIModel:        aiModel,
		},
	)

	assert.NoError(t, err)
	convRepo.AssertExpectations(t)
	msgRepo.AssertExpectations(t)
	bot.AssertExpectations(t)
}
