package decorator

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type MetricsClient interface {
	Inc(key string, value int)
}

type commandMetricsDecorator[C any] struct {
	base   CommandHandler[C]
	client MetricsClient
}

func (d commandMetricsDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	if d.client != nil {
		start := time.Now()

		actionName := strings.ToLower(generateActionName(cmd))

		defer func() {
			end := time.Since(start)

			d.client.Inc(fmt.Sprintf("commands.%s.duration", actionName), int(end.Seconds()))

			if err == nil {
				d.client.Inc(fmt.Sprintf("commands.%s.success", actionName), 1)
			} else {
				d.client.Inc(fmt.Sprintf("commands.%s.failure", actionName), 1)
			}
		}()
	}
	return d.base.Handle(ctx, cmd)
}
