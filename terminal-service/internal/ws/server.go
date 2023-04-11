package ws

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/Banana-Boat/terryminal/terminal-service/internal/pb"
	"github.com/Banana-Boat/terryminal/terminal-service/internal/pty"
	"github.com/Banana-Boat/terryminal/terminal-service/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type WSServer struct {
	config util.Config
	server *gin.Engine
}

// 单个socket连接的上下文
type WSContext struct {
	conn             net.Conn
	config           util.Config
	basePtyContainer *pty.PtyContainer
	gRPCConnection   *grpc.ClientConn
	basePtyClient    pb.BasePtyClient
	ptyStream        pb.BasePty_RunCmdClient
}

func NewWSServer(config util.Config) *WSServer {
	server := gin.Default()

	server.GET("/", func(c *gin.Context) {
		conn, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
		if err != nil {
			log.Error().Err(err).Msg("cannot upgrade http to websocket")
			return
		}
		defer conn.Close()

		wsCtx := &WSContext{
			conn:   conn,
			config: config,
		}
		for {
			msg, _, err := wsutil.ReadClientData(conn)
			if err != nil {
				break
			}

			/* 解析 message */
			var wsMsg Message
			if err := json.Unmarshal(msg, &wsMsg); err != nil {
				log.Error().Err(err).Msg("cannot unmarshal message")
				break
			}

			/* 根据Event字段进行路由 */
			routeByEvent(wsCtx, wsMsg)
		}
	})

	wsServer := &WSServer{
		config: config,
		server: server,
	}

	return wsServer
}

func routeByEvent(wsCtx *WSContext, wsMsg Message) {
	switch wsMsg.Event {
	case "launch":
		/* 将Data字段解析为对应结构体 */
		var data LaunchClientData
		if err := mapstructure.Decode(wsMsg.Data, &data); err != nil {
			log.Error().Err(err).Msg("cannot decode data")
			return
		}

		launchHandle(wsCtx, data.ContainerName, wsCtx.config)

	case "close":
		closeHandle(wsCtx)

	case "run-cmd":
		/* 将Data字段解析为对应结构体 */
		var data RunCmdClientData
		if err := mapstructure.Decode(wsMsg.Data, &data); err != nil {
			log.Error().Err(err).Msg("cannot decode data")
			return
		}

		runCmdHandle(wsCtx, data.Cmd)
	}
}

func (wsServer *WSServer) Start() error {
	if err := wsServer.server.Run(
		fmt.Sprintf("%s:%s", wsServer.config.TerminalWSServerHost, wsServer.config.TerminalWSServerPort),
	); err != nil {
		return err
	}

	return nil
}
