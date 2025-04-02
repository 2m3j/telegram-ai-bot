package decorator

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

func ApplyCommandDecorators[C any](handler CommandHandler[C], logger *slog.Logger, metricsClient MetricsClient) CommandHandler[C] {
	return &commandLoggingDecorator[C]{
		base: &commandMetricsDecorator[C]{
			base:   handler,
			client: metricsClient,
		},
		logger: logger,
	}
}

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
