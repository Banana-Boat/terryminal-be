package api

import (
	"context"
	"fmt"
	"io"

	"github.com/Banana-Boat/terryminal/main-service/internal/pb"
	"github.com/Banana-Boat/terryminal/main-service/internal/pty"
	"github.com/Banana-Boat/terryminal/main-service/internal/util"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func startEventHandle(wsCtx *WSContext, ptyId string, config util.Config) {
	/* 启动容器 */
	if err = pty.StartPty(ptyId); err != nil {
		log.Error().Err(err).Msg("failed to start pty container")
		sendMessage(wsCtx.conn, ptyId, "start", StartServerData{Result: false})
		return
	}

	/* 创建gRPC Client */
	ptyName, err := pty.GetPtyName(ptyId) // 获取容器名
	gRPCConnection, err := grpc.Dial(
		fmt.Sprintf("%s:%s", ptyName, config.BasePtyPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create gRPC client")
		sendMessage(wsCtx.conn, ptyId, "start", StartServerData{Result: false})
		return
	}
	basePtyClient := pb.NewBasePtyClient(gRPCConnection)

	/* 调用RunCmd方法获取数据流对象，创建go routine接受数据流的数据，转发到client */
	ptyStream, err := basePtyClient.RunCmd(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("failed to create gRPC client")
		sendMessage(wsCtx.conn, ptyId, "start", StartServerData{Result: false})
		return
	}
	go func() {
		for {
			resp, err := ptyStream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				return
			}

			sendMessage(wsCtx.conn, ptyId, "run-cmd", RunCmdServerData{IsError: false, Result: resp.Result})
		}
	}()

	/* 将对象存入context */
	ptyHandler := &PtyHandler{
		gRPCConn:   gRPCConnection,
		gRPCClient: basePtyClient,
		gRPCStream: ptyStream,
	}
	wsCtx.PtyHandlerMap[ptyId] = ptyHandler

	/* 向客户端发送成功的消息 */
	log.Info().Msgf("PtyId: %s, start pty container and create gRPC client successfully", ptyId)
	sendMessage(wsCtx.conn, ptyId, "start", StartServerData{Result: true})
}

func runCmdEventHandle(wsCtx *WSContext, ptyId string, cmd string) {
	if cmd == "exit" { // 后续需要补充退出的命令 Ctr+D / Ctrl+C
		log.Warn().Msgf("receive invalid command: %s", cmd)
		sendMessage(wsCtx.conn, ptyId, "run-cmd", RunCmdServerData{IsError: true, Result: "命令不合法"})
		return
	}

	log.Info().Msgf("run-cmd send: %s", cmd)
	wsCtx.PtyHandlerMap[ptyId].gRPCStream.Send(&pb.RunCmdRequest{
		Cmd: cmd,
	})
}

func endEventHandle(wsCtx *WSContext, ptyId string) {
	if err := end(wsCtx, ptyId); err != nil {
		log.Error().Err(err).Msgf("PtyID: %s, failed to stop pty and close gRPC client", ptyId)
		sendMessage(wsCtx.conn, ptyId, "end", EndServerData{Result: false})
		return
	}

	sendMessage(wsCtx.conn, ptyId, "end", EndServerData{Result: true})
}

func end(wsCtx *WSContext, ptyId string) error {
	if wsCtx.PtyHandlerMap[ptyId] == nil {
		return fmt.Errorf("ptyId: %s not found", ptyId)
	}

	ptyHandler := wsCtx.PtyHandlerMap[ptyId]
	if err := ptyHandler.gRPCConn.Close(); err != nil {
		return err
	}
	if err := pty.StopPty(ptyId); err != nil {
		return err
	}

	delete(wsCtx.PtyHandlerMap, ptyId)
	return nil
}

/* 遍历所有的ptyHandler，关闭gRPC连接，停止容器 */
func endAll(wsCtx *WSContext) {
	for ptyId := range wsCtx.PtyHandlerMap {
		if err := end(wsCtx, ptyId); err != nil {
			log.Error().Err(err).Msgf("PtyId: %s, failed to stop pty and close gRPC client", ptyId)
			continue
		}

		log.Info().Msgf("PtyId: %s, stop pty and close gRPC client successfully", ptyId)
	}
}
