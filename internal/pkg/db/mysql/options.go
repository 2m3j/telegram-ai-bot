package mysql

import "time"

type Option func(c *Client)

func WithConnMaxLifetime(connMaxLifetime time.Duration) Option {
	return func(c *Client) {
		c.connMaxLifetime = connMaxLifetime
	}
}
func WithMaxOpenConns(maxOpenConns int) Option {
	return func(c *Client) {
		c.maxOpenConns = maxOpenConns
	}
}

func WithMaxIdleConns(maxIdleConns int) Option {
	return func(c *Client) {
		c.maxIdleConns = maxIdleConns
	}
}
