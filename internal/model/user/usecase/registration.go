package usecase

import (
	"context"
	"fmt"
	"log/slog"

	cstent "bot/internal/model/user/entity"
	userrepo "bot/internal/model/user/repository"
	"bot/internal/pkg/decorator"
)

type Registration struct {
	UserID       uint64
	AIPlatform   string
	AIModel      string
	Username     string
	FirstName    string
	LastName     string
	LanguageCode string
}

type RegistrationHandler decorator.CommandHandler[Registration]

type registrationHandler struct {
	userRepository userrepo.UserRepository
}

func NewRegistrationHandler(
	userRepository userrepo.UserRepository,
	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) RegistrationHandler {
	return decorator.ApplyCommandDecorators[Registration](
		&registrationHandler{
			userRepository,
		},
		logger,
		metricsClient,
	)
}

func (h *registrationHandler) Handle(ctx context.Context, c Registration) error {
	user := cstent.NewUser(c.UserID, c.AIPlatform, c.AIModel, c.Username, c.FirstName, c.LastName, c.LanguageCode)
	err := h.userRepository.Add(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to add new user %+v: %w", c, err)
	}
	return nil
}
