package entity

import (
	"time"

	uuidgen "bot/internal/pkg/uuid"
	"github.com/google/uuid"
)

const (
	StatusNew     = "new"
	StatusSuccess = "success"
	StatusError   = "error"
)

type MessageID uuid.UUID

func (c MessageID) String() string {
	return uuid.UUID(c).String()
}
func NextMessageID() MessageID {
	return MessageID(uuidgen.Next())
}
func ParseMessageID(s string) (MessageID, error) {
	parsedUUID, err := uuid.Parse(s)
	if err != nil {
		return MessageID{}, err
	}
	return MessageID(parsedUUID), nil
}

type Message struct {
	ID                MessageID
	ConversationID    ConversationID
	Status            string
	UserMessage       string
	AssistantPlatform string
	AssistantModel    string
	AssistantMessage  string
	UpdatedAt         time.Time
	CreatedAt         time.Time
}

func NewMessage(id MessageID, conversationID ConversationID, userMessage string, assistantPlatform string, assistantModel string) *Message {
	now := time.Now()
	return &Message{
		ID:                id,
		ConversationID:    conversationID,
		Status:            StatusNew,
		UserMessage:       userMessage,
		AssistantPlatform: assistantPlatform,
		AssistantModel:    assistantModel,
		AssistantMessage:  "",
		UpdatedAt:         now,
		CreatedAt:         now,
	}
}

func (m *Message) Change(assistantMessage string, status string) {
	m.AssistantMessage = assistantMessage
	m.Status = status
	m.UpdatedAt = time.Now()
}
