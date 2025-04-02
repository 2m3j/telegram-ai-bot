package repository

import (
	"context"

	"bot/internal/model/user/entity"
)

type UserRepository interface {
	Add(ctx context.Context, e *entity.User) error
	Update(ctx context.Context, e *entity.User) error
	FindByID(ctx context.Context, id uint64) (*entity.User, error)
	GetByID(ctx context.Context, id uint64) (*entity.User, error)
}
