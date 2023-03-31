package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	containerTypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type bashContainer struct {
	id   string
	name string
}

func NewBashContainer(name string) (*bashContainer, error) {
	ctx := context.Background()

	/* 创建docker client */
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	/* 拉取docker image */
	imageName := "bash:alpine3.16"
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	defer out.Close()
	// io.Copy(os.Stdout, out) // 打印拉取进度（输出流）

	// /* 创建docker container */
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, name)
	if err != nil {
		return nil, err
	}

	// /* 启动容器 */
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	bashContainer := &bashContainer{
		id:   resp.ID,
		name: name,
	}
	return bashContainer, nil
}

func (container *bashContainer) Stop() error {
	/* 创建docker client */
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	/* 停止容器 */
	if err := cli.ContainerStop(context.Background(), container.id, containerTypes.StopOptions{}); err != nil {
		return err
	}

	return nil
}

func (container *bashContainer) Remove() error {
	/* 创建docker client */
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	/* 停止容器 */
	if err := cli.ContainerRemove(context.Background(), container.id, types.ContainerRemoveOptions{}); err != nil {
		return err
	}

	return nil
}
