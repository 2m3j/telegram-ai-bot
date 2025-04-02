package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"bot/config"
	botapi "bot/internal/app/bot"
	"bot/internal/pkg/db/mysql"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	cfg, err := config.New()
	failOnError(cancel, err)

	mysqlClient, err := initMysql(cfg.Mysql)
	defer func() {
		_ = mysqlClient.Close()
	}()
	failOnError(cancel, err)

	botClient := botapi.New(cfg, mysqlClient)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return botClient.Run(ctx)
	})

	g.Go(func() error {
		<-ctx.Done()
		//todo bot close or shutdown
		return nil
	})
	err = g.Wait()
	failOnError(cancel, err)
}

func initMysql(mysqlConfig config.Mysql) (*mysql.Client, error) {
	mysqlOpts := []mysql.Option{
		mysql.WithConnMaxLifetime(time.Duration(mysqlConfig.ConnTimeout) * time.Second),
		mysql.WithMaxOpenConns(int(mysqlConfig.MaxOpenConnections)),
		mysql.WithMaxIdleConns(int(mysqlConfig.MaxIdleConnections)),
	}
	mysqlClient, err := mysql.New(
		mysqlConfig.Dsn,
		mysqlOpts...,
	)
	if err != nil {
		return nil, err
	}
	return mysqlClient, nil
}

func failOnError(cancelFunc context.CancelFunc, err error) {
	if err == nil {
		return
	}
	cancelFunc()
	panic(err)
}
