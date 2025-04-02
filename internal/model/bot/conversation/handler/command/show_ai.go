package command

import (
	"context"
	"fmt"

	"bot/internal/model/bot/conversation/storage"
	botapi "bot/internal/pkg/bot"
)

type ShowAIHandler struct {
	bot botapi.Bot
}

func NewShowAIHandler(bot botapi.Bot) *ShowAIHandler {
	return &ShowAIHandler{bot: bot}
}

func (c *ShowAIHandler) Handle(ctx context.Context, im *botapi.IncomingMessage) error {
	err := c.bot.ShowDialog(ctx, im.ChatID, storage.AIMenu, "start")
	if err != nil {
		return fmt.Errorf("failed to show AI menu dialog: %w", err)
	}
	return nil
}
