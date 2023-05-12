package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	GatewayServerHost   string `mapstructure:"GATEWAY_SERVER_HOST"`
	GatewayServerPort   string `mapstructure:"GATEWAY_SERVER_PORT"`
	TerminalServiceHost string `mapstructure:"TERMINAL_SERVICE_HOST"`
	TerminalServicePort string `mapstructure:"TERMINAL_SERVICE_PORT"`
	ChatbotServiceHost  string `mapstructure:"CHATBOT_SERVICE_HOST"`
	ChatbotServicePort  string `mapstructure:"CHATBOT_SERVICE_PORT"`

	MigrationFileUrl string `mapstructure:"MIGRATION_FILE_URL"`
	DBDriver         string `mapstructure:"DB_DRIVER"`
	DBUsername       string `mapstructure:"DB_USERNAME"`
	DBPassword       string `mapstructure:"DB_PASSWORD"`
	DBHost           string `mapstructure:"DB_HOST"`
	DBPort           string `mapstructure:"DB_PORT"`
	DBName           string `mapstructure:"DB_NAME"`

	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
