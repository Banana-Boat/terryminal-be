package main

import (
	"os"

	"github.com/Banana-Boat/terryminal/terminal-service/internal/util"
	"github.com/Banana-Boat/terryminal/terminal-service/internal/ws"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// 美化zerolog
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	/* 加载配置 */
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Error().Err(err).Msg("cannot load config")
		return
	}

	/* 运行Websocket服务 */
	runWSServer(config)
}

func runWSServer(config util.Config) {
	wsServer := ws.NewWSServer(config)

	log.Info().Msgf(
		"websocket server started at %s:%s successfully",
		config.TerminalWSServerHost, config.TerminalWSServerPort,
	)
	if err := wsServer.Start(); err != nil {
		log.Error().Err(err).Msg("cannot start websocket server")
	}
}
