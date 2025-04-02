package usecase_test

import (
	"context"
	"errors"
	"testing"

	trmmocks "bot/internal/mocks/github.com/avito-tech/go-transaction-manager/trm/v2"
	cnvent "bot/internal/model/ai/conversation/entity"
	aiconvrepomocks "bot/internal/model/ai/conversation/repository/mocks"
	cnvuc "bot/internal/model/ai/conversation/usecase"
	botmocks "bot/internal/pkg/bot/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEndConversationsHandler_Handle_Success(t *testing.T) {
	var userID, chatID uint64 = 1, 2
	ctx := context.Background()
	repo := new(aiconvrepomocks.MockConversationRepository)
	bot := new(botmocks.MockBot)
	trmng := new(trmmocks.MockManager)
	foundConv := cnvent.NewConversation(cnvent.NextConversationID(), userID)
	repo.On("FindByCriteria", ctx, mock.Anything, mock.Anything, uint64(0), uint64(0)).Return([]*cnvent.Conversation{foundConv}, nil)
	trmng.On(
		"Do",
		mock.Anything,
		mock.Anything,
	).Return(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)
	repo.On("Update", ctx, foundConv).Return(nil)
	bot.On("SendMessage", mock.Anything, chatID, mock.AnythingOfType("string")).Return(nil)
	handler := cnvuc.NewEndConversationsHandler(trmng, bot, repo, nil, nil)

	err := handler.Handle(ctx, cnvuc.EndConversations{UserID: userID, ChatID: chatID})

	assert.NoError(t, err)
	assert.False(t, foundConv.EndsAt.IsZero())
	repo.AssertExpectations(t)
	trmng.AssertExpectations(t)
	bot.AssertExpectations(t)
}

func TestEndConversationsHandler_Handle_NoActiveConversations(t *testing.T) {
	var userID, chatID uint64 = 1, 2
	ctx := context.Background()
	repo := new(aiconvrepomocks.MockConversationRepository)
	bot := new(botmocks.MockBot)
	trmng := new(trmmocks.MockManager)
	handler := cnvuc.NewEndConversationsHandler(trmng, bot, repo, nil, nil)
	repo.On("FindByCriteria", ctx, mock.Anything, mock.Anything, uint64(0), uint64(0)).Return([]*cnvent.Conversation{}, nil)
	bot.On("SendMessage", ctx, chatID, mock.AnythingOfType("string")).Return(nil)

	err := handler.Handle(ctx, cnvuc.EndConversations{UserID: userID, ChatID: chatID})

	assert.NoError(t, err)
	trmng.AssertNotCalled(t, "Do")
	repo.AssertNotCalled(t, "Update")
	repo.AssertExpectations(t)
	bot.AssertExpectations(t)
}

func TestEndConversationsHandler_Handle_FindError(t *testing.T) {
	ctx := context.Background()
	repo := new(aiconvrepomocks.MockConversationRepository)
	bot := new(botmocks.MockBot)
	trmng := new(trmmocks.MockManager)
	handler := cnvuc.NewEndConversationsHandler(trmng, bot, repo, nil, nil)
	repo.On("FindByCriteria", ctx, mock.Anything, mock.Anything, uint64(0), uint64(0)).Return([]*cnvent.Conversation{}, errors.New("find error"))

	err := handler.Handle(ctx, cnvuc.EndConversations{UserID: 1, ChatID: 2})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find active conversations")
}

func TestEndConversationsHandler_Handle_UpdateError(t *testing.T) {
	var userID, chatID uint64 = 1, 2
	ctx := context.Background()
	foundConv := cnvent.NewConversation(cnvent.NextConversationID(), userID)
	repo := new(aiconvrepomocks.MockConversationRepository)
	repo.On("FindByCriteria", ctx, mock.Anything, mock.Anything, uint64(0), uint64(0)).Return([]*cnvent.Conversation{foundConv}, nil)
	repo.On("Update", ctx, foundConv).Return(errors.New("update error"))
	bot := new(botmocks.MockBot)
	trmng := new(trmmocks.MockManager)
	trmng.On(
		"Do",
		mock.Anything,
		mock.Anything,
	).Return(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)
	handler := cnvuc.NewEndConversationsHandler(trmng, bot, repo, nil, nil)

	err := handler.Handle(ctx, cnvuc.EndConversations{UserID: userID, ChatID: chatID})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update conversation")
}

func TestEndConversationsHandler_Handle_SendMessageError(t *testing.T) {
	var userID, chatID uint64 = 1, 2
	ctx := context.Background()
	foundConv := cnvent.NewConversation(cnvent.NextConversationID(), userID)
	repo := new(aiconvrepomocks.MockConversationRepository)
	repo.On("FindByCriteria", ctx, mock.Anything, mock.Anything, uint64(0), uint64(0)).Return([]*cnvent.Conversation{foundConv}, nil)
	repo.On("Update", ctx, foundConv).Return(nil)
	trmng := new(trmmocks.MockManager)
	trmng.On(
		"Do",
		mock.Anything,
		mock.Anything,
	).Return(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	)
	bot := new(botmocks.MockBot)
	bot.On("SendMessage", ctx, chatID, mock.AnythingOfType("string")).Return(errors.New("send error"))
	handler := cnvuc.NewEndConversationsHandler(trmng, bot, repo, nil, nil)

	err := handler.Handle(ctx, cnvuc.EndConversations{UserID: userID, ChatID: chatID})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send message")
}
