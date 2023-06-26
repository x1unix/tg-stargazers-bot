package config

type Config struct {
	HTTP     HTTPConfig
	Redis    RedisConfig
	Log      LogConfig
	Telegram TelegramConfig
	GitHub   GitHubConfig
}

type HTTPConfig struct {
	ListenAddress string `envconfig:"HTTP_LISTEN" default:":8080"`
}

type RedisConfig struct {
	DSN string `envconfig:"REDIS_DSN" required:"true"`
}

type TelegramConfig struct {
	BotToken string `envconfig:"TELEGRAM_BOT_TOKEN" required:"true"`
}
