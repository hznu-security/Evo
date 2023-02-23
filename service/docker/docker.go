package docker

import "github.com/docker/docker/client"

func getClient() (*client.Client, error) {
	return client.NewClient("tcp://192.168.154.128:2375", "1.41", nil, nil)
}
