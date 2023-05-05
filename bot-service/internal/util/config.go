package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	BotHttpServerHost string `mapstructure:"BOT_HTTP_SERVER_HOST"`
	BotHttpServerPort string `mapstructure:"BOT_HTTP_SERVER_PORT"`

	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`
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
