package mysql

import (
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
)

const (
	_defaultConnMaxLifetime = time.Minute
	_defaultMaxOpenConns    = 10
	_defaultMaxIdleConns    = 10
	_defaultPingAttempts    = 5
)

type Client struct {
	connMaxLifetime time.Duration
	maxOpenConns    int
	maxIdleConns    int
	Builder         squirrel.StatementBuilderType
	Pool            *sql.DB
}

func New(dsn string, options ...Option) (*Client, error) {
	c := &Client{
		connMaxLifetime: _defaultConnMaxLifetime,
		maxOpenConns:    _defaultMaxOpenConns,
		maxIdleConns:    _defaultMaxIdleConns,
	}
	for _, o := range options {
		o(c)
	}
	c.Builder = squirrel.StatementBuilder
	var err error
	c.Pool, err = connect(dsn)
	if err != nil {
		return nil, err
	}
	c.Pool.SetConnMaxLifetime(c.connMaxLifetime)
	c.Pool.SetMaxOpenConns(c.maxOpenConns)
	c.Pool.SetMaxIdleConns(c.maxIdleConns)
	return c, nil
}

func connect(dsn string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for i := 0; i < _defaultPingAttempts; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (m *Client) Close() error {
	if m.Pool != nil {
		err := m.Pool.Close()
		return err
	}
	return nil
}
