package repository

import (
	"context"
	"database/sql"

	"bot/internal/model/ai/conversation/entity"
)

type scanner interface {
	Scan(dest ...any) error
}

type ConversationRepository interface {
	Add(ctx context.Context, e *entity.Conversation) error
	Update(ctx context.Context, e *entity.Conversation) error
	FindById(ctx context.Context, ID entity.ConversationID) (*entity.Conversation, error)
	FindOneByCriteria(ctx context.Context, criteria ConversationCriteria, sort ConversationSort, offset uint64) (*entity.Conversation, error)
	FindByCriteria(ctx context.Context, criteria ConversationCriteria, sort ConversationSort, limit uint64, offset uint64) ([]*entity.Conversation, error)
}

type MessageRepository interface {
	Add(ctx context.Context, e *entity.Message) error
	Update(ctx context.Context, e *entity.Message) error
	FindById(ctx context.Context, ID entity.MessageID) (*entity.Message, error)
	FindByCriteria(ctx context.Context, criteria MessageCriteria, sort MessageSort, limit uint64, offset uint64) ([]*entity.Message, error)
}

type ConversationCriteria struct {
	UserID   sql.NullInt64
	Finished sql.NullBool
}

func NewConversationCriteria() ConversationCriteria {
	return ConversationCriteria{}
}
func (c ConversationCriteria) WithUserID(userID uint64) ConversationCriteria {
	c.UserID = sql.NullInt64{Int64: int64(userID), Valid: true}
	return c
}
func (c ConversationCriteria) WithFinished(finished bool) ConversationCriteria {
	c.Finished = sql.NullBool{Bool: finished, Valid: true}
	return c
}

type ConversationSort struct {
	ID sql.NullBool
}

func NewConversationSort() ConversationSort {
	return ConversationSort{}
}

func (s ConversationSort) WithID(asc bool) ConversationSort {
	s.ID = sql.NullBool{Bool: asc, Valid: true}
	return s
}

type MessageCriteria struct {
	ConversationID sql.NullString
	CreatedAtFrom  sql.NullTime
	CreatedAtTo    sql.NullTime
	Status         sql.NullString
}

func NewMessageCriteria() MessageCriteria {
	return MessageCriteria{}
}

func (m MessageCriteria) WithConversationID(ConversationID entity.ConversationID) MessageCriteria {
	m.ConversationID = sql.NullString{String: ConversationID.String(), Valid: true}
	return m
}
func (m MessageCriteria) WithStatus(status string) MessageCriteria {
	m.Status = sql.NullString{String: status, Valid: true}
	return m
}

type MessageSort struct {
	ID sql.NullBool
}

func NewMessageSort() MessageSort {
	return MessageSort{}
}

func (s MessageSort) WithID(asc bool) MessageSort {
	s.ID = sql.NullBool{Bool: asc, Valid: true}
	return s
}
