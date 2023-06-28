package config

import (
	"net/url"
)

type Config struct {
	HTTP   HTTPConfig
	Auth   AuthConfig
	Redis  RedisConfig
	Log    LogConfig
	GitHub GitHubConfig
	Bot    BotConfig
}

type HTTPConfig struct {
	ListenAddress string  `envconfig:"HTTP_LISTEN" default:":8080"`
	BaseURL       url.URL `envconfig:"HTTP_BASE_URL" required:"true"`
}

type RedisConfig struct {
	DSN string `envconfig:"REDIS_DSN" required:"true"`
}

type TelegramConfig struct {
	BotToken string `envconfig:"TELEGRAM_BOT_TOKEN" required:"true"`
}

type BotConfig struct {
	Env               Environment `envconfig:"APP_ENV" default:"dev"`
	WebHookSecret     string      `envconfig:"BOT_WEBHOOK_SECRET" required:"true"`
	ListenTimeout     int         `envconfig:"BOT_LISTEN_TIMEOUT" default:"60"`
	MessageBufferSize int         `envconfig:"BOT_MSG_BUFFER_SIZE" default:"100"`
	WorkerPoolSize    int         `envconfig:"BOT_WORKER_POOL_SIZE" default:"3"`
	Telegram          TelegramConfig
}
