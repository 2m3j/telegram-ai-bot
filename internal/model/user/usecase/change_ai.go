package usecase

import (
	"context"
	"fmt"
	"log/slog"

	userrepo "bot/internal/model/user/repository"
	"bot/internal/pkg/decorator"
)

type ChangeAI struct {
	UserID     uint64
	AIPlatform string
	AIModel    string
}

type ChangeAIHandler decorator.CommandHandler[ChangeAI]

type changeAIHandler struct {
	userRepository userrepo.UserRepository
}

func NewChangeAIHandler(
	userRepository userrepo.UserRepository,
	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) ChangeAIHandler {
	return decorator.ApplyCommandDecorators[ChangeAI](
		&changeAIHandler{userRepository},
		logger,
		metricsClient,
	)
}

func (h *changeAIHandler) Handle(ctx context.Context, c ChangeAI) error {
	user, err := h.userRepository.GetByID(ctx, c.UserID)
	if err != nil {
		return fmt.Errorf("failed to retrieve user with ID %d: %w", c.UserID, err)
	}
	user.ChangeAISettings(c.AIPlatform, c.AIModel)
	err = h.userRepository.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user settings %+v: %w", c, err)
	}
	return nil
}
