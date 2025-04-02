package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"bot/internal/model/ai/conversation/entity"
	"bot/internal/pkg/db/mysql"
	trmsql "github.com/avito-tech/go-transaction-manager/drivers/sql/v2"
)

const tableMessage = "ai_message"

type MessageMysqlRepository struct {
	*mysql.Client
	ctxGetter *trmsql.CtxGetter
}

func NewMessageMysqlRepository(mysql *mysql.Client, ctxGetter *trmsql.CtxGetter) *MessageMysqlRepository {
	return &MessageMysqlRepository{mysql, ctxGetter}
}

func (r *MessageMysqlRepository) Add(ctx context.Context, e *entity.Message) error {
	query, args, err := r.Builder.
		Insert(tableMessage).
		Columns(
			"id",
			"conversation_id",
			"status",
			"user_message",
			"assistant_platform",
			"assistant_model",
			"assistant_message",
			"updated_at",
			"created_at").
		Values(
			e.ID.String(),
			e.ConversationID.String(),
			e.Status,
			e.UserMessage,
			e.AssistantPlatform,
			e.AssistantModel,
			e.AssistantMessage,
			e.UpdatedAt,
			e.CreatedAt.Format(time.DateTime),
		).ToSql()
	if err != nil {
		return fmt.Errorf("error building SQL query for adding new message: %w", err)
	}
	_, err = r.ctxGetter.DefaultTrOrDB(ctx, r.Pool).ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error executing SQL query for adding new message: %w", err)
	}

	return nil
}

func (r *MessageMysqlRepository) Update(ctx context.Context, e *entity.Message) error {
	upData := map[string]interface{}{
		"status":            e.Status,
		"assistant_message": e.AssistantMessage,
		"updated_at":        e.UpdatedAt.Format(time.DateTime),
	}
	query, args, err := r.Builder.
		Update(tableMessage).Where("id = ?", e.ID.String()).SetMap(upData).Limit(1).ToSql()
	if err != nil {
		return fmt.Errorf("error building SQL query for updating message: %w", err)
	}
	_, err = r.ctxGetter.DefaultTrOrDB(ctx, r.Pool).ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error executing SQL query for updating message: %w", err)
	}
	return nil
}

func (r *MessageMysqlRepository) FindById(ctx context.Context, ID entity.MessageID) (*entity.Message, error) {
	query := "SELECT * FROM " + tableMessage + " WHERE id=? LIMIT 1"
	e, err := r.hydrate(r.ctxGetter.DefaultTrOrDB(ctx, r.Pool).QueryRowContext(ctx, query, ID.String()))
	if err != nil {
		return e, fmt.Errorf("error finding message by ID %s: %w", ID, err)
	}
	return e, nil
}

func (r *MessageMysqlRepository) FindByCriteria(ctx context.Context, criteria MessageCriteria, sort MessageSort, limit uint64, offset uint64) ([]*entity.Message, error) {
	list, err := r.findByCriteria(ctx, criteria, sort, limit, offset)
	if err != nil {
		return list, fmt.Errorf("error finding messages by criteria %+v: %w", criteria, err)
	}
	return list, nil
}

func (r *MessageMysqlRepository) findByCriteria(ctx context.Context, criteria MessageCriteria, sort MessageSort, limit uint64, offset uint64) ([]*entity.Message, error) {
	qb := r.Builder.Select("*").From(tableMessage)
	if criteria.ConversationID.Valid {
		qb = qb.Where("conversation_id = ?", criteria.ConversationID.String)
	}
	if criteria.CreatedAtFrom.Valid {
		qb = qb.Where("created_at >= ?", criteria.CreatedAtFrom.Time.Format(time.DateTime))
	}
	if criteria.CreatedAtTo.Valid {
		qb = qb.Where("created_at <= ?", criteria.CreatedAtTo.Time.Format(time.DateTime))
	}
	if criteria.Status.Valid {
		qb = qb.Where("status = ?", criteria.Status.String)
	}
	orderBy := make([]string, 0)
	if sort.ID.Valid {
		if sort.ID.Bool {
			orderBy = append(orderBy, "id ASC")
		} else {
			orderBy = append(orderBy, "id DESC")
		}
	}
	if len(orderBy) > 0 {
		qb = qb.OrderBy(orderBy...)
	}
	if limit != 0 {
		qb = qb.Limit(limit)
		qb = qb.Offset(offset)
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	rows, err := r.ctxGetter.DefaultTrOrDB(ctx, r.Pool).QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	entities := make([]*entity.Message, 0)
	for rows.Next() {
		var e *entity.Message
		e, err = r.hydrate(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to hydrate message: %w", err)
		}
		entities = append(entities, e)
	}

	return entities, rows.Err()
}

func (r *MessageMysqlRepository) hydrate(scan scanner) (*entity.Message, error) {
	var createdAtRaw, updatedAtRaw []uint8
	var id, convID string
	e := &entity.Message{}
	err := scan.Scan(
		&id,
		&convID,
		&e.Status,
		&e.UserMessage,
		&e.AssistantPlatform,
		&e.AssistantModel,
		&e.AssistantMessage,
		&updatedAtRaw,
		&createdAtRaw,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return e, fmt.Errorf("failed to scan message: %w", err)
		}
	}
	e.ID, err = entity.ParseMessageID(id)
	if err != nil {
		return e, fmt.Errorf("failed to parse message id: %w", err)
	}
	e.ConversationID, err = entity.ParseConversationID(id)
	if err != nil {
		return e, fmt.Errorf("failed to parse conversation id: %w", err)
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

func (r *MessageMysqlRepository) parseDateTime(datetime []uint8) (time.Time, error) {
	parsedTime, err := time.Parse(time.DateTime, string(datetime))
	if err != nil {
		return parsedTime, fmt.Errorf("failed to parse datetime: %w", err)
	}
	return parsedTime, nil
}
