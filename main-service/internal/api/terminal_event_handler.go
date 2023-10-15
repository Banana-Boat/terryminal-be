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

func startEventHandle(wsCtx *WSContext, ptyID string, config util.Config) {
	/* 启动容器 */
	if err = pty.StartPty(ptyID); err != nil {
		log.Error().Err(err).Msg("failed to start pty container")
		sendMessage(wsCtx.conn, ptyID, "start", StartServerData{Result: false})
		return
	}

	/* 创建gRPC Client */
	ptyName, err := pty.GetPtyName(ptyID) // 获取容器名
	if err != nil {
		log.Error().Err(err).Msg("failed to get pty name")
		sendMessage(wsCtx.conn, ptyID, "start", StartServerData{Result: false})
		return
	}
	gRPCConnection, err := grpc.Dial(
		fmt.Sprintf("%s:%s", ptyName, config.BasePtyPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create gRPC client")
		sendMessage(wsCtx.conn, ptyID, "start", StartServerData{Result: false})
		return
	}
	client := pb.NewBasePtyClient(gRPCConnection)

	/* 调用RunCmd方法获取数据流对象，创建go routine接受数据流的数据，转发到client */
	ptyStream, err := client.RunCmd(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("failed to create gRPC client")
		sendMessage(wsCtx.conn, ptyID, "start", StartServerData{Result: false})
		return
	}
	go func() {
		for {
			resp, err := ptyStream.Recv()
			fmt.Println(resp)
			if err == io.EOF {
				return
			}
			if err != nil {
				return
			}

			sendMessage(wsCtx.conn, ptyID, "run-cmd", RunCmdServerData{IsError: false, Result: resp.Result})
		}
	}()

	/* 将对象存入context */
	ptyHandler := &PtyHandler{
		gRPCConn:   gRPCConnection,
		gRPCClient: client,
		gRPCStream: ptyStream,
	}
	wsCtx.PtyHandlerMap[ptyID] = ptyHandler

	/* 向客户端发送成功的消息 */
	log.Info().Msgf("PtyID: %s, start pty container and create gRPC client successfully", ptyID)
	sendMessage(wsCtx.conn, ptyID, "start", StartServerData{Result: true})
}

func runCmdEventHandle(wsCtx *WSContext, ptyID string, cmd string) {
	if cmd == "exit" { // 后续需要补充退出的命令 Ctr+D / Ctrl+C
		log.Warn().Msgf("receive invalid command: %s", cmd)
		sendMessage(wsCtx.conn, ptyID, "run-cmd", RunCmdServerData{IsError: true, Result: "命令不合法"})
		return
	}

	log.Info().Msgf("run-cmd send: %s", cmd)
	wsCtx.PtyHandlerMap[ptyID].gRPCStream.Send(&pb.RunCmdRequest{
		Cmd: cmd,
	})
}

func endEventHandle(wsCtx *WSContext, ptyID string) {
	if err := end(wsCtx, ptyID); err != nil {
		log.Error().Err(err).Msgf("PtyID: %s, failed to stop pty and close gRPC client", ptyID)
		sendMessage(wsCtx.conn, ptyID, "end", EndServerData{Result: false})
		return
	}

	sendMessage(wsCtx.conn, ptyID, "end", EndServerData{Result: true})
}

func end(wsCtx *WSContext, ptyID string) error {
	if wsCtx.PtyHandlerMap[ptyID] == nil {
		return fmt.Errorf("ptyID: %s not found", ptyID)
	}

	ptyHandler := wsCtx.PtyHandlerMap[ptyID]
	if err := ptyHandler.gRPCConn.Close(); err != nil {
		return err
	}
	if err := pty.StopPty(ptyID); err != nil {
		return err
	}

	delete(wsCtx.PtyHandlerMap, ptyID)

	log.Info().Msgf("PtyID: %s, stop pty and close gRPC client successfully", ptyID)
	return nil
}

/* 遍历所有的ptyHandler，关闭gRPC连接，停止容器 */
func endAll(wsCtx *WSContext) {
	for ptyID := range wsCtx.PtyHandlerMap {
		if err := end(wsCtx, ptyID); err != nil {
			log.Error().Err(err).Msgf("PtyID: %s, failed to stop pty and close gRPC client", ptyID)
			continue
		}
	}
}
