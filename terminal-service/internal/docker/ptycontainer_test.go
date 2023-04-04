package docker

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/Banana-Boat/terryminal/terminal-service/internal/pb"
	"github.com/Banana-Boat/terryminal/terminal-service/internal/util"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestCreateBash(t *testing.T) {
	// github action 不执行该测试
	if testing.Short() {
		t.Skip("Skipping test in github action.")
	}

	/* 加载配置 */
	config, err := util.LoadConfig("../..")
	if err != nil {
		t.Error(err)
	}

	/* 容器创建并启动 */
	containerName := fmt.Sprint(time.Now().Unix())
	bashContainer, err := NewPtyContainer(config.BasePtyImageName, containerName, config.BasePtyNetwork, true)
	if err != nil {
		t.Error(err)
	}
	if err = bashContainer.Start(); err != nil {
		t.Error(err)
	}

	/* 创建gRPC Client，调用BasePtyService的RunCmd */
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", config.BasePtyHost, config.BasePtyPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()
	c := pb.NewBasePtyClient(conn)
	stream, err := c.RunCmd(context.Background())
	if err != nil {
		t.Error(err)
	}

	done := make(chan bool)
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				t.Error(err)
			}
			log.Info().Msg(resp.Result)
		}
	}()

	time.Sleep(time.Millisecond * 1000)
	stream.Send(&pb.RunCmdRequest{
		Cmd: "ls",
	})
	time.Sleep(time.Millisecond * 1000)
	stream.Send(&pb.RunCmdRequest{
		Cmd: "exit",
	})
	<-done

	/* 关闭并删除容器 */
	if err = bashContainer.Stop(); err != nil {
		t.Error(err)
	}
	if err = bashContainer.Remove(); err != nil {
		t.Error(err)
	}
}
