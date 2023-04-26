package ws

import (
	"context"
	"fmt"
	"io"

	"github.com/Banana-Boat/terryminal/terminal-service/internal/pb"
	"github.com/Banana-Boat/terryminal/terminal-service/internal/pty"
	"github.com/Banana-Boat/terryminal/terminal-service/internal/util"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func startHandle(wsCtx *WSContext, ptyID string, config util.Config) {
	/* 创建容器并启动 */
	basePtyContainer, err := pty.NewPtyContainer(
		config.BasePtyImageName, ptyID, config.BasePtyNetwork, nil,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create pty container")
		sendMessage(wsCtx.conn, "start", StartServerData{PtyID: ptyID, Result: false})
		return
	}
	if err = basePtyContainer.Start(); err != nil {
		log.Error().Err(err).Msg("failed to start pty container")
		sendMessage(wsCtx.conn, "start", StartServerData{PtyID: ptyID, Result: false})
		return
	}

	/* 创建gRPC Client */
	gRPCConnection, err := grpc.Dial(
		fmt.Sprintf("%s:%s", ptyID, config.BasePtyPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create gRPC client")
		sendMessage(wsCtx.conn, "start", StartServerData{PtyID: ptyID, Result: false})
		return
	}
	basePtyClient := pb.NewBasePtyClient(gRPCConnection)

	/* 调用RunCmd方法获取数据流对象，创建go routine接受数据流的数据，转发到client */
	ptyStream, err := basePtyClient.RunCmd(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("failed to create gRPC client")
		sendMessage(wsCtx.conn, "start", StartServerData{PtyID: ptyID, Result: false})
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

			sendMessage(wsCtx.conn, "run-cmd", RunCmdServerData{PtyID: ptyID, IsError: false, Result: resp.Result})
			log.Info().Msgf("run-cmd receive: %s", resp.Result)
		}
	}()

	/* 将对象存入context */
	ptyHandler := &PtyHandler{
		container:  basePtyContainer,
		gRPCConn:   gRPCConnection,
		gRPCClient: basePtyClient,
		gRPCStream: ptyStream,
	}
	wsCtx.PtyHandlerMap[ptyID] = ptyHandler

	/* 向客户端发送成功的消息 */
	log.Info().Msg("successed to start pty container and create gRPC client")
	sendMessage(wsCtx.conn, "start", StartServerData{PtyID: ptyID, Result: true})
}

func runCmdHandle(wsCtx *WSContext, ptyID string, cmd string) {
	if cmd == "exit" { // 后续需要补充退出的命令 Ctr+D / Ctrl+C
		log.Warn().Msgf("receive invalid command: %s", cmd)
		sendMessage(wsCtx.conn, "run-cmd", RunCmdServerData{PtyID: ptyID, IsError: true, Result: "命令不合法"})
		return
	}

	log.Info().Msgf("run-cmd send: %s", cmd)
	wsCtx.PtyHandlerMap[ptyID].gRPCStream.Send(&pb.RunCmdRequest{
		Cmd: cmd,
	})
}

func endHandle(wsCtx *WSContext, ptyID string) {
	sendMessage(wsCtx.conn, "end", EndServerData{PtyID: ptyID, Result: true})
	destroy(wsCtx)
}

// 遍历所有的ptyHandler，关闭gRPC连接，停止容器，删除容器
func destroy(wsCtx *WSContext) {
	for ptyID, ptyHandler := range wsCtx.PtyHandlerMap {
		if err := ptyHandler.gRPCConn.Close(); err != nil {
			log.Error().Err(err).Msgf("PtyID: %s, failed to close gRPC Connection", ptyID)
		}
		if err := ptyHandler.container.Stop(); err != nil {
			log.Error().Err(err).Msgf("PtyID: %s, failed to stop basePty container", ptyID)
		}
		if err := ptyHandler.container.Remove(); err != nil {
			log.Error().Err(err).Msgf("PtyID: %s, failed to remove basePty container", ptyID)
		}
		log.Info().Msgf("PtyID: %s, successed to remove pty container and close gRPC client", ptyID)

		delete(wsCtx.PtyHandlerMap, ptyID)
	}
}
