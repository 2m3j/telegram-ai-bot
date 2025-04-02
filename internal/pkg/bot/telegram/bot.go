package telegram

import (
	"bytes"
	"context"
	"strings"
	"time"

	"bot/internal/pkg/bot"
	"bot/internal/pkg/bot/telegram/dialog"
	"bot/internal/pkg/markdown"
	tgmd "github.com/Mad-Pixels/goldmark-tgmd"
	tbotapi "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	maxMessageLen         = 4096
	SendChatActionTimeout = 5
)

type Bot struct {
	token         string
	bot           *tbotapi.Bot
	errorsHandler tbotapi.ErrorsHandler
	menuCommands  map[string]string
	handler       bot.Handler
}

func NewBot(token string, errorsHandler tbotapi.ErrorsHandler, menuCommands map[string]string) *Bot {
	return &Bot{token: token, errorsHandler: errorsHandler, menuCommands: menuCommands}
}

func (b *Bot) Start(ctx context.Context, handler bot.Handler) error {
	b.handler = handler
	opts := []tbotapi.Option{
		tbotapi.WithDefaultHandler(
			func(ctx context.Context, tbot *tbotapi.Bot, update *models.Update) {
				if update.Message == nil {
					return
				}
				isCommand := false
				for _, entity := range update.Message.Entities {
					if entity.Type == models.MessageEntityTypeBotCommand {
						isCommand = true
						break
					}
				}
				handler.Handle(
					ctx,
					&bot.IncomingMessage{
						IsCommand:    isCommand,
						UserID:       uint64(update.Message.From.ID),
						Username:     update.Message.From.Username,
						FirstName:    update.Message.From.FirstName,
						LastName:     update.Message.From.LastName,
						LanguageCode: update.Message.From.LanguageCode,
						ChatID:       uint64(update.Message.Chat.ID),
						Message:      update.Message.Text,
						Date:         time.Unix(int64(update.Message.Date), 0),
					},
				)
			},
		),
		tbotapi.WithErrorsHandler(b.errorsHandler),
	}

	var err error
	b.bot, err = tbotapi.New(b.token, opts...)
	if err != nil {
		return err
	}
	b.setMenuCommands(ctx)
	b.bot.Start(ctx)
	return nil
}

func (b *Bot) SendMessage(ctx context.Context, chatID uint64, message string) error {
	md := tgmd.TGMD()
	var buf bytes.Buffer
	if err := md.Convert([]byte(message), &buf); err != nil {
		return err
	}
	var parts []string
	text := buf.String()
	if len(text) > maxMessageLen {
		parts = markdown.Split(text, maxMessageLen)
	} else {
		parts = append(parts, text)
	}
	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			continue
		}
		_, err := b.bot.SendMessage(ctx, &tbotapi.SendMessageParams{
			ChatID:    chatID,
			Text:      part,
			ParseMode: models.ParseModeMarkdown,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) setMenuCommands(ctx context.Context) {
	var err error
	var bCommands = make([]models.BotCommand, len(b.menuCommands))
	index := 0
	for com, description := range b.menuCommands {
		bCommands[index] = models.BotCommand{Command: com, Description: description}
		index++
	}
	_, err = b.bot.SetMyCommands(ctx, &tbotapi.SetMyCommandsParams{Commands: bCommands})
	if err != nil {
		b.errorsHandler(err)
	}
}

func (b *Bot) SendChatAction(ctx context.Context, chatID uint64, action bot.ChatAction) {
	ticker := time.NewTicker(SendChatActionTimeout * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := b.sendChatAction(ctx, chatID, action)
			if err != nil {
				b.errorsHandler(err)
			}
		}
	}
}

func (b *Bot) ShowDialog(ctx context.Context, chatID uint64, nodes []bot.DialogNode, startNode string) error {
	if b.handler == nil {
		panic("bot not started")
	}
	var err error
	d := dialog.New(b.bot, b.handler, nodes)
	_, err = d.Show(ctx, b.bot, chatID, startNode)
	return err
}

func (b *Bot) sendChatAction(ctx context.Context, chatID uint64, action bot.ChatAction) error {
	_, err := b.bot.SendChatAction(ctx, &tbotapi.SendChatActionParams{ChatID: chatID, Action: models.ChatAction(action)})
	return err
}
