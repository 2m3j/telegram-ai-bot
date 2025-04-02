package handler

import (
	"context"
	"fmt"

	"bot/internal/model/bot/conversation/handler/command"
	botapi "bot/internal/pkg/bot"
)

type MessageCommandHandler struct {
	commandHandlers       map[string]CommandHandler
	unknownCommandHandler *command.SendUnknownHandler
}

func NewMessageCommandHandler(commandHandlers map[string]CommandHandler, unknownCommandHandler *command.SendUnknownHandler) *MessageCommandHandler {
	return &MessageCommandHandler{
		commandHandlers:       commandHandlers,
		unknownCommandHandler: unknownCommandHandler,
	}
}

func (h *MessageCommandHandler) Handle(ctx context.Context, incomingMessage *botapi.IncomingMessage) error {
	if h.commandHandlers[incomingMessage.Message] != nil {
		err := h.commandHandlers[incomingMessage.Message].Handle(ctx, incomingMessage)
		if err != nil {
			return fmt.Errorf("failed to handle command '%s': %w", incomingMessage.Message, err)
		}
	} else {
		err := h.unknownCommandHandler.Handle(ctx, incomingMessage)
		if err != nil {
			return fmt.Errorf("failed to handle unknown command '%s': %w", incomingMessage.Message, err)
		}
	}
	return nil
}
