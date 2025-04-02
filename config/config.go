package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	App      App      `envPrefix:"APP_"`
	Telegram Telegram `envPrefix:"TELEGRAM_"`
	Mysql    Mysql    `envPrefix:"MYSQL_"`
	AI       AI       `envPrefix:"AI_"`
	Bot      Bot      `envPrefix:"BOT_"`
}
type App struct {
	Env string `env:"ENV" envDefault:"dev"`
}
type Telegram struct {
	Token string `env:"TOKEN,required"`
}
type AI struct {
	OpenAI   Openai   `envPrefix:"OPENAI_"`
	DeepSeek DeepSeek `envPrefix:"DEEP_SEEK_"`
}
type Openai struct {
	Token string `env:"TOKEN,required"`
}
type DeepSeek struct {
	Token string `env:"TOKEN,required"`
}

// todo add perform WarningTimeOut, WarningMessagesLen options
type Bot struct {
	SkipMessageTimeout uint   `env:"SKIP_MESSAGE_TIMEOUT" envDefault:"300"`
	WarningTimeOut     uint   `env:"WARNING_TIMEOUT" envDefault:"86400"`
	WarningMessagesLen uint   `env:"WARNING_MESSAGES_LEN" envDefault:"15"`
	DefaultAIPlatform  string `env:"DEFAULT_AI_PLATFORM" envDefault:"DeepSeek"`
	DefaultAIModel     string `env:"DEFAULT_AI_MODEL" envDefault:"deepseek-chat"`
}

type Mysql struct {
	Dsn                string `env:"DSN,required"`
	MaxOpenConnections uint   `env:"MAX_OPEN_CONNECTIONS" envDefault:"100"`
	MaxIdleConnections uint   `env:"MAX_IDLE_CONNECTIONS" envDefault:"100"`
	ConnTimeout        uint   `env:"CONN_TIMEOUT" envDefault:"300"`
}

func New() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	var cfg Config
	err = env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	return &cfg, nil
}
