/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/13 19:43
 * 描述     ：关于容器的操作
 */

package docker

import (
	"Evo/auth"
	"context"
	"github.com/docker/go-connections/nat"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func StartContainer(image, name string, portMap *nat.PortMap) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}
	exports := make(nat.PortSet)
	for k, _ := range *portMap {
		exports[k] = struct{}{}
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        image,
		ExposedPorts: exports,
	}, &container.HostConfig{
		PortBindings: *portMap,
	}, nil, nil, name)
	if err != nil {
		return err
	}

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return nil
}

// SetContainerSSH 给容器设置ssh账号密码
func SetContainerSSH(container, user, pwd string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}
	id, _ := cli.ContainerExecCreate(ctx, container, types.ExecConfig{
		AttachStderr: true,
		AttachStdin:  true,
		AttachStdout: true,
		Cmd:          []string{"passwd"},
	})
	attach, _ := cli.ContainerExecAttach(ctx, id.ID, types.ExecStartCheck{})
	err = cli.ContainerExecStart(ctx, id.ID, types.ExecStartCheck{
		Tty: true,
	})
	defer attach.Close()
	if err != nil {
		return err
	}
	if _, err = attach.Conn.Write([]byte(pwd + "\n")); err != nil {
		return err
	}

	if _, err = attach.Conn.Write([]byte(pwd + "\n")); err != nil {
		return err
	}
	return nil
}

// RemoveContainer kill并删除容器
func RemoveContainer(name string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}
	//  kill并删除容器
	err = cli.ContainerRemove(ctx, name, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		return err
	}
	return nil
}

// ResetContainer 重置容器,会将容器重启并设置ssh密码
func ResetContainer(name string) (pwd string, err error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return "", err
	}
	var timeout time.Duration
	timeout = time.Second * 3
	err = cli.ContainerRestart(ctx, name, &timeout)
	if err != nil {
		return "", err
	}

	pwd = auth.NewPwd()
	err = SetContainerSSH(name, "root", pwd)
	if err != nil {
		return "", err
	}
	return pwd, nil
}

// ContainerExec 通过exec在容器中执行命令
func ContainerExec(container string, command string) (types.ContainerExecInspect, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return types.ContainerExecInspect{}, err
	}
	id, _ := cli.ContainerExecCreate(ctx, container, types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		AttachStdin:  true,
		Cmd:          []string{command},
	})
	attach, err := cli.ContainerExecAttach(ctx, id.ID, types.ExecStartCheck{})
	if err != nil {
		return types.ContainerExecInspect{}, err
	}
	defer attach.Close()
	err = cli.ContainerExecStart(ctx, id.ID, types.ExecStartCheck{
		Tty: true,
	})
	inspect, err := cli.ContainerExecInspect(ctx, id.ID)
	return inspect, nil
}

func ParsePort(port string) nat.PortMap {
	portMap := make(map[nat.Port][]nat.PortBinding)
	portSet := make(nat.PortSet, 0)
	portBind := strings.Split(port, ",")
	for _, v := range portBind {
		ports := strings.Split(v, ":")
		portSet[nat.Port(ports[1])] = struct{}{}
		portMap[nat.Port(ports[1])] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: ports[0],
			},
		}
	}
	return portMap
}

func GetIp(name string) (string, error) {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts()
	res, err := cli.ContainerInspect(ctx, name)
	if err != nil {
		return "", err
	}
	return res.NetworkSettings.IPAddress, nil
}
