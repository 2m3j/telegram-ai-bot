package bot

import (
	"context"
	"time"
)

type ChatAction string
type Name string

const (
	ChatActionTyping ChatAction = "typing"
)

type Bot interface {
	Start(ctx context.Context, handler Handler) error
	SendMessage(ctx context.Context, chatID uint64, message string) error
	SendChatAction(ctx context.Context, chatID uint64, action ChatAction)
	ShowDialog(ctx context.Context, chatID uint64, nodes []DialogNode, startNode string) error
}

type Handler interface {
	Handle(ctx context.Context, message *IncomingMessage)
}

type ErrorsHandler func(err error)

type IncomingMessage struct {
	IsCommand    bool
	IsDialog     bool
	UserID       uint64
	Username     string
	FirstName    string
	LastName     string
	LanguageCode string
	ChatID       uint64
	Message      string
	Date         time.Time
}

type DialogButton struct {
	Text   string
	NodeID string
	URL    string
}

type DialogNode struct {
	ID       string
	Text     string
	Command  string
	Keyboard [][]DialogButton
}
