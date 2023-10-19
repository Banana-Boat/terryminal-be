package pty

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	containerTypes "github.com/docker/docker/api/types/container"
	networkTypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type PtyPortMap struct {
	HostPort      string
	ContainerPort string
}

/* 创建容器 */
func NewPty(imageName string, containerName string, network string, ptyPortMap *PtyPortMap) (string, error) {
	ctx := context.Background()

	/* 创建docker client */
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}
	defer cli.Close()

	/* 如果不存在docker image则拉取 */
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return "", err
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
			return "", err
		}
		defer out.Close()
		io.Copy(os.Stdout, out) // 必须输出否则报错，待查明！！
	}

	/* 创建docker container */
	var hostConfig *containerTypes.HostConfig = nil
	var networkConfig = &networkTypes.NetworkingConfig{
		EndpointsConfig: map[string]*networkTypes.EndpointSettings{
			network: {},
		},
	}
	// 本地测试需要做端口映射，无需网络配置
	if ptyPortMap != nil {
		hostConfig = &containerTypes.HostConfig{
			PortBindings: nat.PortMap{
				nat.Port(fmt.Sprintf("%s/tcp", ptyPortMap.ContainerPort)): []nat.PortBinding{{
					HostPort: ptyPortMap.HostPort,
				}},
			},
		}
		networkConfig = nil
	}
	resp, err := cli.ContainerCreate(ctx,
		&containerTypes.Config{
			Image: imageName,
		},
		hostConfig,
		networkConfig, nil, containerName)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

/* 根据容器Id获取容信息 */
func GetPtyInfo(id string) (types.ContainerJSON, error) {
	/* 创建docker client */
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return types.ContainerJSON{}, err
	}
	defer cli.Close()

	/* 获取容器信息 */
	container, err := cli.ContainerInspect(context.Background(), id)
	if err != nil {
		return types.ContainerJSON{}, err
	}

	// 容器名第一个字符为/，需要去掉
	return container, nil
}

/* 启动容器 */
func StartPty(id string) error {
	/* 创建docker client */
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	// /* 启动容器 */
	if err := cli.ContainerStart(
		context.Background(),
		id,
		types.ContainerStartOptions{},
	); err != nil {
		return err
	}

	return nil
}

/* 停止容器 */
func StopPty(id string) error {
	/* 创建docker client */
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	/* 停止容器 */
	if err := cli.ContainerStop(
		context.Background(),
		id,
		containerTypes.StopOptions{},
	); err != nil {
		return err
	}

	return nil
}

/* 销毁容器 */
func RemovePty(id string) error {
	/* 创建docker client */
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	/* 删除容器 */
	if err := cli.ContainerRemove(
		context.Background(),
		id,
		types.ContainerRemoveOptions{},
	); err != nil {
		return err
	}

	return nil
}
