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

const tableConversation = "ai_conversation"

type ConversationMysqlRepository struct {
	*mysql.Client
	ctxGetter *trmsql.CtxGetter
}

func NewConversationMysqlRepository(mysql *mysql.Client, ctxGetter *trmsql.CtxGetter) *ConversationMysqlRepository {
	return &ConversationMysqlRepository{mysql, ctxGetter}
}

func (r *ConversationMysqlRepository) Add(ctx context.Context, e *entity.Conversation) error {
	var endsAt sql.NullTime
	if !e.EndsAt.IsZero() {
		endsAt.Time = e.EndsAt
		endsAt.Valid = true
	}
	query, args, err := r.Builder.
		Insert(tableConversation).
		Columns(
			"id",
			"user_id",
			"started_at",
			"ended_at").
		Values(
			e.ID.String(),
			e.UserID,
			e.StartedAt.Format(time.DateTime),
			endsAt,
		).ToSql()
	if err != nil {
		return fmt.Errorf("error building SQL query for adding new conversation: %w", err)
	}
	_, err = r.ctxGetter.DefaultTrOrDB(ctx, r.Pool).ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error executing SQL query for adding new conversation: %w", err)
	}
	return nil
}

func (r *ConversationMysqlRepository) Update(ctx context.Context, e *entity.Conversation) error {
	var endsAt sql.NullTime
	if !e.EndsAt.IsZero() {
		endsAt.Time = e.EndsAt
		endsAt.Valid = true
	}

	upData := map[string]interface{}{
		"ended_at": endsAt,
	}
	query, args, err := r.Builder.
		Update(tableConversation).Where("id = ?", e.ID.String()).SetMap(upData).Limit(1).ToSql()
	if err != nil {
		return fmt.Errorf("error building SQL query for updating conversation: %w", err)
	}
	_, err = r.ctxGetter.DefaultTrOrDB(ctx, r.Pool).ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error executing SQL query for updating conversation: %w", err)
	}
	return nil
}

func (r *ConversationMysqlRepository) FindById(ctx context.Context, ID entity.ConversationID) (*entity.Conversation, error) {
	query := "SELECT * FROM " + tableConversation + " WHERE id=? LIMIT 1"
	e, err := r.hydrate(r.ctxGetter.DefaultTrOrDB(ctx, r.Pool).QueryRowContext(ctx, query, ID.String()))
	if err != nil {
		return e, fmt.Errorf("error finding conversation by ID %s: %w", ID, err)
	}
	return e, nil
}

func (r *ConversationMysqlRepository) FindOneByCriteria(ctx context.Context, criteria ConversationCriteria, sort ConversationSort, offset uint64) (*entity.Conversation, error) {
	list, err := r.findByCriteria(ctx, criteria, sort, 1, offset)
	if err != nil {
		return nil, fmt.Errorf("error finding conversations by criteria %+v: %w", criteria, err)
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

func (r *ConversationMysqlRepository) FindByCriteria(ctx context.Context, criteria ConversationCriteria, sort ConversationSort, limit uint64, offset uint64) ([]*entity.Conversation, error) {
	list, err := r.findByCriteria(ctx, criteria, sort, limit, offset)
	if err != nil {
		return list, fmt.Errorf("error finding conversations by criteria %+v: %w", criteria, err)
	}
	return list, nil
}

func (r *ConversationMysqlRepository) findByCriteria(ctx context.Context, criteria ConversationCriteria, sort ConversationSort, limit uint64, offset uint64) ([]*entity.Conversation, error) {
	qb := r.Builder.Select("*").From(tableConversation)
	if criteria.UserID.Valid {
		qb = qb.Where("user_id = ?", criteria.UserID.Int64)
	}
	if criteria.Finished.Valid {
		if criteria.Finished.Bool {
			qb = qb.Where("ended_at IS NOT NULL")
		} else {
			qb = qb.Where("ended_at IS NULL")
		}
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
	entities := make([]*entity.Conversation, 0)
	for rows.Next() {
		var e *entity.Conversation
		e, err = r.hydrate(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to hydrate conversation: %w", err)
		}
		entities = append(entities, e)
	}
	return entities, rows.Err()
}

func (r *ConversationMysqlRepository) hydrate(scan scanner) (*entity.Conversation, error) {
	var startedAtRaw []uint8
	var endedAtRaw []uint8
	var id string
	var err error
	e := &entity.Conversation{}
	err = scan.Scan(
		&id,
		&e.UserID,
		&startedAtRaw,
		&endedAtRaw,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return e, fmt.Errorf("failed to scan conversation: %w", err)
		}
	}
	e.ID, err = entity.ParseConversationID(id)
	if err != nil {
		return e, fmt.Errorf("failed to parse conversation id: %w", err)
	}

	parsedStartedAt, err := r.parseDateTime(startedAtRaw)
	if err != nil {
		return e, fmt.Errorf("failed to parse started_at: %w", err)
	}
	e.StartedAt = parsedStartedAt

	var parsedEndedAt time.Time
	if endedAtRaw != nil {
		parsedEndedAt, err = r.parseDateTime(endedAtRaw)
		if err != nil {
			return e, fmt.Errorf("failed to parse ended_at: %w", err)
		}
		e.EndsAt = parsedEndedAt
	}
	return e, nil
}

func (r *ConversationMysqlRepository) parseDateTime(datetime []uint8) (time.Time, error) {
	parsedTime, err := time.Parse(time.DateTime, string(datetime))
	if err != nil {
		return parsedTime, fmt.Errorf("failed to parse datetime: %w", err)
	}
	return parsedTime, nil
}
