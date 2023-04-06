package ws

import (
	"fmt"
	"net/http"

	"github.com/Banana-Boat/terryminal/terminal-service/internal/pb"
	"github.com/Banana-Boat/terryminal/terminal-service/internal/pty"
	"github.com/Banana-Boat/terryminal/terminal-service/internal/util"
	socketio "github.com/googollee/go-socket.io"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type WSServer struct {
	config util.Config
	server *socketio.Server
}

type wsContext struct {
	basePtyContainer *pty.PtyContainer
	gRPCConnection   *grpc.ClientConn
	basePtyClient    pb.BasePtyClient
	ptyStream        pb.BasePty_RunCmdClient
}

func NewWSServer(config util.Config) *WSServer {
	wsServer := &WSServer{
		config: config,
		server: socketio.NewServer(nil),
	}

	wsServer.setupRouter()

	return wsServer
}

func (wsServer *WSServer) setupRouter() {

	wsServer.server.OnConnect("/", func(s socketio.Conn) error {
		log.Info().Msgf("connected: %s", s.ID())
		return nil
	})

	wsServer.server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Info().Msgf("disconnect: %s, reason: %s", s.ID(), reason)
	})

	wsServer.server.OnError("/", func(s socketio.Conn, err error) {
		log.Error().Err(err).Msg("websocket error")
	})

	wsServer.server.OnEvent("/", "launch", func(s socketio.Conn, containerName string) {
		launchHandle(s, containerName, wsServer.config)
	})

	wsServer.server.OnEvent("/", "close", closeHandle)

	wsServer.server.OnEvent("/", "run-cmd", runCmdHandle)
}

func (wsServer *WSServer) Start() error {
	go wsServer.server.Serve()
	defer wsServer.server.Close()

	// 框架默认地址前缀，无需改动
	http.Handle("/socket.io/", wsServer.server)

	if err := http.ListenAndServe(
		fmt.Sprintf("%s:%s", wsServer.config.TerminalWSServerHost, wsServer.config.TerminalWSServerPort),
		nil,
	); err != nil {
		return err
	}

	return nil
}
