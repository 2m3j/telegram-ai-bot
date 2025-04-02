package bot

import (
	"context"

	"bot/config"
	"bot/internal/pkg/db/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	deps *Dependencies
}

func New(cfg *config.Config, mysqlClient *mysql.Client) *App {
	app := &App{NewDependencies(cfg, mysqlClient)}
	return app
}

func (app *App) Run(ctx context.Context) error {
	return app.deps.bot.Start(ctx, app.deps.botCnvHandler)
}
