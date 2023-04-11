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

func launchHandle(wsCtx *WSContext, containerName string, config util.Config) {
	/* 创建容器并启动 */

	// 端口映射待去除！！！！
	basePtyContainer, err := pty.NewPtyContainer(
		config.BasePtyImageName, containerName, config.BasePtyNetwork,
		&pty.PtyPortMap{
			HostPort:      config.BasePtyPort,
			ContainerPort: config.BasePtyPort,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create pty container")
		sendMessage(wsCtx.conn, "launch", LaunchServerData{Result: false})
		return
	}
	if err = basePtyContainer.Start(); err != nil {
		log.Error().Err(err).Msg("failed to start pty container")
		sendMessage(wsCtx.conn, "launch", LaunchServerData{Result: false})
		return
	}

	/* 创建gRPC Client */
	gRPCConnection, err := grpc.Dial(
		fmt.Sprintf("%s:%s", config.BasePtyHost, config.BasePtyPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create gRPC client")
		sendMessage(wsCtx.conn, "launch", LaunchServerData{Result: false})
		return
	}
	basePtyClient := pb.NewBasePtyClient(gRPCConnection)

	/* 调用RunCmd方法获取数据流对象，创建go routine接受数据流的数据，转发到client */
	ptyStream, err := basePtyClient.RunCmd(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("failed to create gRPC client")
		sendMessage(wsCtx.conn, "launch", LaunchServerData{Result: false})
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

			sendMessage(wsCtx.conn, "run-cmd", RunCmdServerData{Result: resp.Result})
			log.Info().Msgf("run-cmd receive: %s", resp.Result)
		}
	}()

	/* 将对象存入context */
	wsCtx.basePtyClient = basePtyClient
	wsCtx.basePtyContainer = basePtyContainer
	wsCtx.gRPCConnection = gRPCConnection
	wsCtx.ptyStream = ptyStream

	/* 向客户端发送成功的消息 */
	log.Info().Msg("successed to start pty container and create gRPC client")
	sendMessage(wsCtx.conn, "launch", LaunchServerData{Result: true})
}

func runCmdHandle(wsCtx *WSContext, cmd string) {
	/* 如果传入为exit，则关闭gRPC连接 & 关闭并删除容器 */
	if cmd == "exit" { // 后续需要补充退出的命令 Ctr+D / Ctrl+C
		if err := wsCtx.gRPCConnection.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close gRPC Connection")
		}
		if err := wsCtx.basePtyContainer.Stop(); err != nil {
			log.Error().Err(err).Msg("failed to stop basePty container")
		}
		if err := wsCtx.basePtyContainer.Remove(); err != nil {
			log.Error().Err(err).Msg("failed to remove basePty container")
		}
		log.Info().Msg("successed to remove pty container and close gRPC client")
		return
	}

	log.Info().Msgf("run-cmd send: %s", cmd)
	wsCtx.ptyStream.Send(&pb.RunCmdRequest{
		Cmd: cmd,
	})
}

func closeHandle(wsCtx *WSContext) {
	sendMessage(wsCtx.conn, "close", CloseServerData{Result: true})

	wsCtx.conn.Close()
}
