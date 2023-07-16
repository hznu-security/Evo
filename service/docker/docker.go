package docker

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/spf13/viper"
)

var (
	host    string
	version string
)

func InitDocker() {
	host = viper.GetString("docker.host")
	version = viper.GetString("docker.version")
}

func getClient() (*client.Client, error) {
	if host == "" {
		return client.NewClientWithOpts()
	}
	return client.NewClient(fmt.Sprintf("tcp://%s", host), version, nil, nil)
}
