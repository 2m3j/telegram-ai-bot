package dialog

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"bot/internal/pkg/bot"
	tbotapi "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type OnErrorHandler func(err error)

type Dialog struct {
	data              []string
	prefix            string
	onError           OnErrorHandler
	nodes             []bot.DialogNode
	commandHandler    bot.Handler
	callbackHandlerID string
}

func New(b *tbotapi.Bot, commandHandler bot.Handler, nodes []bot.DialogNode, opts ...Option) *Dialog {
	p := &Dialog{
		prefix:         tbotapi.RandomString(16),
		onError:        defaultOnError,
		nodes:          nodes,
		commandHandler: commandHandler,
	}

	for _, opt := range opts {
		opt(p)
	}

	p.callbackHandlerID = b.RegisterHandler(tbotapi.HandlerTypeCallbackQueryData, p.prefix, tbotapi.MatchTypePrefix, p.callback)

	return p
}

// Prefix returns the prefix of the widget
func (d *Dialog) Prefix() string {
	return d.prefix
}

func defaultOnError(err error) {
	log.Printf("[TG-UI-DIALOG] [ERROR] %s", err)
}

func (d *Dialog) showNode(ctx context.Context, b *tbotapi.Bot, chatID any, node bot.DialogNode) (*models.Message, error) {
	params := &tbotapi.SendMessageParams{
		ChatID:      chatID,
		Text:        node.Text,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: d.buildKB(node, d.prefix),
	}

	return b.SendMessage(ctx, params)
}

func (d *Dialog) Show(ctx context.Context, b *tbotapi.Bot, chatID any, nodeID string) (*models.Message, error) {
	node, ok := d.findNode(nodeID)
	if !ok {
		return nil, fmt.Errorf("failed to find node with id %s", nodeID)
	}

	return d.showNode(ctx, b, chatID, node)
}

func (d *Dialog) callback(ctx context.Context, b *tbotapi.Bot, update *models.Update) {
	ok, err := b.AnswerCallbackQuery(ctx, &tbotapi.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})
	if err != nil {
		d.onError(err)
	}
	if !ok {
		d.onError(fmt.Errorf("failed to answer callback query"))
	}

	nodeID := strings.TrimPrefix(update.CallbackQuery.Data, d.prefix)
	node, ok := d.findNode(nodeID)
	if !ok {
		d.onError(fmt.Errorf("failed to find node with id %s", nodeID))
		return
	}

	if "" != node.Command {
		d.commandHandler.Handle(
			ctx,
			&bot.IncomingMessage{
				IsCommand:    true,
				IsDialog:     true,
				UserID:       uint64(update.CallbackQuery.From.ID),
				Username:     update.CallbackQuery.From.Username,
				FirstName:    update.CallbackQuery.From.FirstName,
				LastName:     update.CallbackQuery.From.LastName,
				LanguageCode: update.CallbackQuery.From.LanguageCode,
				ChatID:       uint64(update.CallbackQuery.Message.Message.Chat.ID),
				Message:      node.Command,
				Date:         time.Unix(int64(update.CallbackQuery.Message.Message.Date), 0),
			})
		_, errDelete := b.DeleteMessage(ctx, &tbotapi.DeleteMessageParams{
			ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.Message.ID,
		})
		if errDelete != nil {
			d.onError(errDelete)
		}
		return
	}

	_, errEdit := b.EditMessageText(ctx, &tbotapi.EditMessageTextParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		Text:        node.Text,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: d.buildKB(node, d.prefix),
	})
	if errEdit != nil {
		d.onError(errEdit)
	}
	return
}

func (d *Dialog) findNode(id string) (bot.DialogNode, bool) {
	for _, node := range d.nodes {
		if node.ID == id {
			return node, true
		}
	}

	return bot.DialogNode{}, false
}

func (d *Dialog) buildKB(n bot.DialogNode, prefix string) models.ReplyMarkup {
	if len(n.Keyboard) == 0 {
		return nil
	}

	var kb [][]models.InlineKeyboardButton

	for _, row := range n.Keyboard {
		var kbRow []models.InlineKeyboardButton
		for _, btn := range row {
			b := models.InlineKeyboardButton{
				Text: btn.Text,
			}
			if btn.URL != "" {
				b.URL = btn.URL
			} else {
				b.CallbackData = prefix + btn.NodeID
			}
			kbRow = append(kbRow, b)
		}
		kb = append(kb, kbRow)
	}

	return models.InlineKeyboardMarkup{
		InlineKeyboard: kb,
	}
}
