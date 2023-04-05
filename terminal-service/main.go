package main

import (
	"github.com/Banana-Boat/terryminal/terminal-service/internal/util"
	"github.com/Banana-Boat/terryminal/terminal-service/internal/ws"
	"github.com/rs/zerolog/log"
)

func main() {
	/* 加载配置 */
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
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
		log.Fatal().Err(err).Msg("cannot start websocket server")
	}
}
