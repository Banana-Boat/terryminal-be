package main

import (
	"os"

	"github.com/Banana-Boat/terryminal/bot-service/internal/http"
	"github.com/Banana-Boat/terryminal/bot-service/internal/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// 美化zerolog
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	/* 加载配置 */
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
		return
	}

	/* 运行 http 服务 */
	runHttpServer(config)
}

func runHttpServer(config util.Config) {
	server := http.NewServer(config)

	if err := server.Start(); err != nil {
		log.Error().Err(err).Msg("failed to start server")
	}
}
