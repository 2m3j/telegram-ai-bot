package decorator

import (
	"context"
	"fmt"
	"log/slog"
)

type commandLoggingDecorator[C any] struct {
	base   CommandHandler[C]
	logger *slog.Logger
}

func (d commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	if d.logger != nil {
		handlerType := generateActionName(cmd)

		logger := d.logger.With(
			slog.String("command", handlerType),
			slog.String("command_body", fmt.Sprintf("%+v", cmd)),
		)
		logger.Debug("Executing command")
		defer func() {
			if err == nil {
				logger.Info("Command executed successfully")
			} else {
				logger.Error("Failed to execute command", slog.String("error", err.Error()))
			}
		}()
	}
	return d.base.Handle(ctx, cmd)
}
