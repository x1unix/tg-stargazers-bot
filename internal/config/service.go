package config

type Config struct {
	HTTP   HTTPConfig
	Redis  RedisConfig
	Log    LogConfig
	GitHub GitHubConfig
	Bot    BotConfig
}

type HTTPConfig struct {
	ListenAddress string `envconfig:"HTTP_LISTEN" default:":8080"`
	BaseURL       string `envconfig:"HTTP_BASE_URL" required:"true"`
}

type RedisConfig struct {
	DSN string `envconfig:"REDIS_DSN" required:"true"`
}

type TelegramConfig struct {
	BotToken string `envconfig:"TELEGRAM_BOT_TOKEN" required:"true"`
}

type BotConfig struct {
	Env                 Environment `envconfig:"APP_ENV" default:"dev"`
	ListenTimeout       int         `envconfig:"BOT_LISTEN_TIMEOUT" default:"60"`
	MessageBufferSize   int         `envconfig:"BOT_MSG_BUFFER_SIZE" default:"100"`
	WebHookURLPath      string      `envconfig:"BOT_WEBHOOK_URL_PATH" default:"/webhooks/tg"`
	WorkerPoolSize      int         `envconfig:"BOT_WORKER_POOL_SIZE" default:"3"`
	UpdateWebhookOnBoot bool        `envconfig:"BOT_UPDATE_WEBHOOK_ON_BOOT" default:"false"`
	Telegram            TelegramConfig
}
