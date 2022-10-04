/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/13 19:43
 * 描述     ：关于容器的操作
 */

package docker

import (
	"Evo/auth"
	"context"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func StartContainer(image, name, net, ip string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
	}, nil, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			net: {
				IPAMConfig: &network.EndpointIPAMConfig{
					IPv4Address: ip,
				},
				IPAddress: ip,
			},
		},
	}, nil, name)
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

	res, _ := cli.ContainerExecAttach(ctx, id.ID, types.ExecStartCheck{})
	err = cli.ContainerExecStart(ctx, id.ID, types.ExecStartCheck{
		Tty: true,
	})
	if err != nil {
		return err
	}
	if _, err = res.Conn.Write([]byte(pwd + "\n")); err != nil {
		return err
	}

	if _, err = res.Conn.Write([]byte(pwd + "\n")); err != nil {
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
func ContainerExec(container string, command string) ([]byte, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}
	cmd := strings.Split(command, " ")
	id, _ := cli.ContainerExecCreate(ctx, container, types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		AttachStdin:  true,
		Cmd:          cmd,
	})
	attach, err := cli.ContainerExecAttach(ctx, id.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, err
	}
	err = cli.ContainerExecStart(ctx, id.ID, types.ExecStartCheck{
		Tty: true,
	})
	if err != nil {
		return nil, err
	}
	res := make([]byte, 0)
	_, err = attach.Conn.Read(res)
	if err != nil {
		return nil,err
	}
	return	res, nil
}
