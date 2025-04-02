package handler

import (
	"context"
	"fmt"
	"time"

	botapi "bot/internal/pkg/bot"
)

type CommandHandler interface {
	Handle(ctx context.Context, incomingMessage *botapi.IncomingMessage) error
}

type ErrorsHandler func(err error)

type Handler struct {
	skipMessageTimeout    time.Duration
	messageHandler        *MessageHandler
	messageCommandHandler *MessageCommandHandler
	userProcess           *UserProcess
	errorsHandler         ErrorsHandler
}

func NewHandler(
	skipMessageTimeout time.Duration,
	messageHandler *MessageHandler,
	messageCommandHandler *MessageCommandHandler,
	userProcess *UserProcess,
	errorsHandler ErrorsHandler,
) *Handler {
	h := &Handler{
		skipMessageTimeout,
		messageHandler,
		messageCommandHandler,
		userProcess,
		errorsHandler,
	}
	return h
}

func (h *Handler) Handle(ctx context.Context, im *botapi.IncomingMessage) {
	if im.Date.Unix() <= time.Now().Add(-h.skipMessageTimeout).Unix() {
		return
	}
	user, err := h.userProcess.CreateOrUpdate(ctx, im)
	if err != nil {
		h.handleError(fmt.Errorf("failed to create or update user for message %+v: %w", im, err))
		return
	}
	if im.IsCommand {
		err = h.messageCommandHandler.Handle(ctx, im)
		if err != nil {
			h.handleError(fmt.Errorf("failed to handle command for message %+v: %w", im, err))
		}
	} else {
		err = h.messageHandler.Handle(ctx, im, user)
		if err != nil {
			h.handleError(fmt.Errorf("failed to handle message %+v: %w", im, err))
		}
	}
}

func (h *Handler) handleError(err error) {
	h.errorsHandler(fmt.Errorf("conversation handler error: %w", err))
}
