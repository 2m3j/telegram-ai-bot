package decorator

import (
	"context"

	botapi "bot/internal/pkg/bot"
)

type ChatIDProvider interface {
	GetChatID() uint64
}

type commandTypingDecorator[C ChatIDProvider] struct {
	base CommandHandler[C]
	bot  botapi.Bot
}

func ApplyCommandTypingDecorator[C ChatIDProvider](handler CommandHandler[C], bot botapi.Bot) CommandHandler[C] {
	return &commandTypingDecorator[C]{
		base: handler,
		bot:  bot,
	}
}

func (d commandTypingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	sendChatActionCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go d.bot.SendChatAction(sendChatActionCtx, cmd.GetChatID(), botapi.ChatActionTyping)
	return d.base.Handle(ctx, cmd)
}
