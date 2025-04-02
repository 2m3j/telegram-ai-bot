package handler

import (
	"context"
	"fmt"

	"bot/internal/model/user/entity"
	userrepo "bot/internal/model/user/repository"
	userusecase "bot/internal/model/user/usecase"
	aiapi "bot/internal/pkg/ai"
	botapi "bot/internal/pkg/bot"
)

type UserProcess struct {
	aiPlatform          aiapi.PlatformName
	aiModel             aiapi.PlatformName
	userRepository      userrepo.UserRepository
	registrationHandler userusecase.RegistrationHandler
}

func NewUserProcess(
	aiPlatform aiapi.PlatformName,
	aiModel aiapi.PlatformName,
	userRepository userrepo.UserRepository,
	registrationHandler userusecase.RegistrationHandler,
) *UserProcess {
	return &UserProcess{
		aiPlatform:          aiPlatform,
		aiModel:             aiModel,
		userRepository:      userRepository,
		registrationHandler: registrationHandler,
	}
}

func (h *UserProcess) CreateOrUpdate(ctx context.Context, im *botapi.IncomingMessage) (*entity.User, error) {
	user, err := h.userRepository.FindByID(ctx, im.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user with ID %d: %w", im.UserID, err)
	}
	if user == nil {
		err = h.registrationHandler.Handle(ctx, userusecase.Registration{UserID: im.UserID, AIPlatform: h.aiPlatform, AIModel: h.aiModel, Username: im.Username, FirstName: im.FirstName, LastName: im.LastName, LanguageCode: im.LanguageCode})
		if err != nil {
			return nil, fmt.Errorf("failed to register user with ID %d: %w", im.UserID, err)
		}
		user, err = h.userRepository.GetByID(ctx, im.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve newly registered user with ID %d: %w", im.UserID, err)
		}
	}
	if !user.IsEqualUserInfo(im.Username, im.FirstName, im.LastName, im.LanguageCode) {
		user.ChangeUserInfo(im.Username, im.FirstName, im.LastName, im.LanguageCode)
		err = h.userRepository.Update(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("failed to update user information for user with ID %d: %w", im.UserID, err)
		}
	}
	return user, err
}
