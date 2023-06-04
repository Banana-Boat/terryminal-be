package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	MainServerHost     string `mapstructure:"MAIN_SERVER_HOST"`
	MainServerPort     string `mapstructure:"MAIN_SERVER_PORT"`
	ChatbotServiceHost string `mapstructure:"CHATBOT_SERVICE_HOST"`
	ChatbotServicePort string `mapstructure:"CHATBOT_SERVICE_PORT"`

	MigrationFileUrl string `mapstructure:"MIGRATION_FILE_URL"`
	DBDriver         string `mapstructure:"DB_DRIVER"`
	DBUsername       string `mapstructure:"DB_USERNAME"`
	DBPassword       string `mapstructure:"DB_PASSWORD"`
	DBHost           string `mapstructure:"DB_HOST"`
	DBPort           string `mapstructure:"DB_PORT"`
	DBName           string `mapstructure:"DB_NAME"`

	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`

	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`

	BasePtyHost      string `mapstructure:"BASE_PTY_HOST"`
	BasePtyPort      string `mapstructure:"BASE_PTY_PORT"`
	BasePtyImageName string `mapstructure:"BASE_PTY_IMAGE_NAME"`
	BasePtyNetwork   string `mapstructure:"BASE_PTY_NETWORK"`
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
