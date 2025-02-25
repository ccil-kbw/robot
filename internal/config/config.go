package config

var (
	Cfg = config{}
)

type config struct {
	Environment        string `env:"ENVIRONMENT" envDefault:"dev"`
	DiscordServerID    string `env:"DISCORD_SERVER_ID"`
	DiscordBotToken    string `env:"DISCORD_BOT_TOKEN"`
	DiscordBotAsPublic bool   `env:"DISCORD_BOT_AS_PUBLIC" envDefault:"true"`
}
