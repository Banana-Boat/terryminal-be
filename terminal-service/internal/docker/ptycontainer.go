package docker

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	containerTypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type PtyContainer struct {
	ID            string
	ContainerName string
	ImageName     string
}

func NewPtyContainer(imageName string, containerName string) (*PtyContainer, error) {
	ctx := context.Background()

	/* 创建docker client */
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	/* 如果不存在docker image则拉取 */
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return nil, err
	}
	isImageExist := false
	for _, image := range images {
		if image.RepoTags[0] == imageName {
			isImageExist = true
			break
		}
	}
	if !isImageExist {
		out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			return nil, err
		}
		defer out.Close()
		io.Copy(os.Stdout, out) // 必须输出否则报错，待查明！！
	}

	// /* 创建docker container */
	resp, err := cli.ContainerCreate(ctx, &containerTypes.Config{
		Image: imageName,
	}, nil, nil, nil, containerName)
	if err != nil {
		return nil, err
	}

	return &PtyContainer{
		ID:            resp.ID,
		ContainerName: containerName,
		ImageName:     imageName,
	}, nil
}

func (ptyContainer *PtyContainer) Start() error {
	/* 创建docker client */
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	// /* 启动容器 */
	if err := cli.ContainerStart(context.Background(), ptyContainer.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

func (ptyContainer *PtyContainer) Stop() error {
	/* 创建docker client */
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	/* 停止容器 */
	if err := cli.ContainerStop(context.Background(), ptyContainer.ID, containerTypes.StopOptions{}); err != nil {
		return err
	}

	return nil
}

func (ptyContainer *PtyContainer) Remove() error {
	/* 创建docker client */
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	/* 停止容器 */
	if err := cli.ContainerRemove(context.Background(), ptyContainer.ID, types.ContainerRemoveOptions{}); err != nil {
		return err
	}

	return nil
}
