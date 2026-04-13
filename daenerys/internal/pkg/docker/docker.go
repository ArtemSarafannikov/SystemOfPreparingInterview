package docker

import (
	"github.com/moby/moby/client"
)

// Output .
type Output struct {
	Stdout string
	Stderr string
}

// Client .
type Client struct {
	docker *client.Client
}

// New .
func New(
	dockerClient *client.Client,
) *Client {
	docker := &Client{
		docker: dockerClient,
	}

	return docker
}
