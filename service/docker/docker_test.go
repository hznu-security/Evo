package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"testing"
)

func TestDockerHost(t *testing.T) {
	cli, err := client.NewClient(fmt.Sprintf("tcp://%s", "192.168.154.1:2375"), "1.41", nil, nil)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	ctx := context.Background()
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	t.Log(images)
}
