package entity

import (
	"time"

	uuidgen "bot/internal/pkg/uuid"
	"github.com/google/uuid"
)

type ConversationID uuid.UUID

func (c ConversationID) String() string {
	return uuid.UUID(c).String()
}

func NextConversationID() ConversationID {
	return ConversationID(uuidgen.Next())
}
func ParseConversationID(s string) (ConversationID, error) {
	parsedUUID, err := uuid.Parse(s)
	if err != nil {
		return ConversationID{}, err
	}
	return ConversationID(parsedUUID), nil
}

type Conversation struct {
	ID        ConversationID
	UserID    uint64
	UpdatedAt time.Time
	StartedAt time.Time
	EndsAt    time.Time
}

func NewConversation(
	id ConversationID,
	userID uint64,
) *Conversation {
	now := time.Now()
	return &Conversation{
		ID:        id,
		UserID:    userID,
		UpdatedAt: now,
		StartedAt: now,
		EndsAt:    time.Time{},
	}
}

func (c *Conversation) End() {
	c.UpdatedAt = time.Now()
	c.EndsAt = c.UpdatedAt
}
