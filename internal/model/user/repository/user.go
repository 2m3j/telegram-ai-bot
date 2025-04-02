package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"bot/internal/model/user/entity"
	"bot/internal/pkg/db/mysql"
	trmsql "github.com/avito-tech/go-transaction-manager/drivers/sql/v2"
)

const tableUser = "user"

type UserMysqlRepository struct {
	*mysql.Client
	ctxGetter *trmsql.CtxGetter
}

func NewUserMysqlRepository(mysql *mysql.Client, ctxGetter *trmsql.CtxGetter) *UserMysqlRepository {
	return &UserMysqlRepository{mysql, ctxGetter}
}

func (r *UserMysqlRepository) Add(ctx context.Context, e *entity.User) error {
	query, args, err := r.Builder.
		Insert(tableUser).
		Columns(
			"id",
			"ai_platform",
			"ai_model",
			"username",
			"first_name",
			"last_name",
			"language_code",
			"updated_at",
			"created_at",
		).
		Values(
			e.ID,
			e.AIPlatform,
			e.AIModel,
			e.Username,
			e.FirstName,
			e.LastName,
			e.LanguageCode,
			e.UpdatedAt.Format(time.DateTime),
			e.CreatedAt.Format(time.DateTime),
		).ToSql()
	if err != nil {
		return fmt.Errorf("error building SQL query for adding new user: %w", err)
	}
	_, err = r.ctxGetter.DefaultTrOrDB(ctx, r.Pool).ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error executing SQL query for adding new user: %w", err)
	}
	return nil
}

func (r *UserMysqlRepository) Update(ctx context.Context, e *entity.User) error {
	upData := map[string]interface{}{
		"ai_platform":   e.AIPlatform,
		"ai_model":      e.AIModel,
		"username":      e.Username,
		"first_name":    e.FirstName,
		"last_name":     e.LastName,
		"language_code": e.LanguageCode,
		"updated_at":    e.UpdatedAt.Format(time.DateTime),
		"created_at":    e.CreatedAt.Format(time.DateTime),
	}
	query, args, err := r.Builder.
		Update(tableUser).Where("id = ?", e.ID).SetMap(upData).Limit(1).ToSql()
	if err != nil {
		return fmt.Errorf("error building SQL query for updating user: %w", err)
	}
	_, err = r.ctxGetter.DefaultTrOrDB(ctx, r.Pool).ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error executing SQL query for updating user: %w", err)
	}
	return nil
}

func (r *UserMysqlRepository) FindByID(ctx context.Context, id uint64) (*entity.User, error) {
	e, err := r.findByID(ctx, id)
	if err != nil {
		return e, fmt.Errorf("error finding user by ID %d: %w", id, err)
	}
	return e, nil
}

func (r *UserMysqlRepository) GetByID(ctx context.Context, id uint64) (*entity.User, error) {
	e, err := r.findByID(ctx, id)
	if err != nil {
		return e, fmt.Errorf("error finding user by ID %d: %w", id, err)
	}
	if e == nil {
		return nil, fmt.Errorf("user not found for ID: %d", id)
	}
	return e, nil
}

func (r *UserMysqlRepository) findByID(ctx context.Context, id uint64) (*entity.User, error) {
	query := "SELECT * FROM " + tableUser + " WHERE id=? LIMIT 1"
	e, err := r.hydrate(r.ctxGetter.DefaultTrOrDB(ctx, r.Pool).QueryRowContext(ctx, query, id))
	if err != nil {
		return nil, fmt.Errorf("failed to hydrate user: %w", err)
	}
	return e, nil
}

func (r *UserMysqlRepository) hydrate(row *sql.Row) (*entity.User, error) {
	var updatedAtRaw, createdAtRaw []uint8
	var err error
	e := &entity.User{}
	err = row.Scan(
		&e.ID,
		&e.AIPlatform,
		&e.AIModel,
		&e.Username,
		&e.FirstName,
		&e.LastName,
		&e.LanguageCode,
		&updatedAtRaw,
		&createdAtRaw,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return e, fmt.Errorf("failed to scan user: %w", err)
		}
	}

	parsedUpdatedAt, err := r.parseDateTime(updatedAtRaw)
	if err != nil {
		return e, fmt.Errorf("failed to parse updated_at: %w", err)
	}
	e.UpdatedAt = parsedUpdatedAt

	parsedCreatedAt, err := r.parseDateTime(createdAtRaw)
	if err != nil {
		return e, fmt.Errorf("failed to parse created_at: %w", err)
	}
	e.CreatedAt = parsedCreatedAt
	return e, nil
}

func (r *UserMysqlRepository) parseDateTime(datetime []uint8) (time.Time, error) {
	parsedTime, err := time.Parse(time.DateTime, string(datetime))
	if err != nil {
		return parsedTime, fmt.Errorf("failed to parse datetime: %w", err)
	}
	return parsedTime, nil
}
