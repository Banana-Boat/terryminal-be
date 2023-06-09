package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	ChatbotServerHost string `mapstructure:"CHATBOT_SERVER_HOST"`
	ChatbotServerPort string `mapstructure:"CHATBOT_SERVER_PORT"`

	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`

	Api2dUrl string `mapstructure:"API2D_URL"`
	Api2dKey string `mapstructure:"API2D_KEY"`
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
