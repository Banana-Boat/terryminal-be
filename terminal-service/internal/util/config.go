package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	TerminalServerHost string `mapstructure:"TERMINAL_SERVER_HOST"`
	TerminalServerPort string `mapstructure:"TERMINAL_SERVER_PORT"`

	BasePtyHost      string `mapstructure:"BASE_PTY_HOST"`
	BasePtyPort      string `mapstructure:"BASE_PTY_PORT"`
	BasePtyImageName string `mapstructure:"BASE_PTY_IMAGE_NAME"`
	BasePtyNetwork   string `mapstructure:"BASE_PTY_NETWORK"`

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
