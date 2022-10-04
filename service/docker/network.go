/* COPYRIGHT NOTICE
 * 作者     ：ymk
 * 创建时间 ：2022/07/21 19:43
 * 描述     ：关于容器的操作
 */

package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func CreateNetwork(name string, subnet string) (Warning string, err error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return "", err
	}

	resp, err := cli.NetworkCreate(ctx, name, types.NetworkCreate{
		IPAM: &network.IPAM{
			Config: []network.IPAMConfig{
				{
					Subnet: subnet,
				},
			},
		},
	})
	if err != nil {
		return "", err
	}
	return resp.Warning, nil
}

func DelNetwork(name string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}
	err = cli.NetworkRemove(ctx, name)
	return err
}
