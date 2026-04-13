package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
)

// CreateContainer создает контейнер, запускает его и оставляет висеть
//
// Возвращает ID контейнера
func (c *Client) CreateContainer(ctx context.Context, imageName string, memoryLimitBytes int64) (string, error) {
	var hostConfig *container.HostConfig
	if memoryLimitBytes > 0 {
		hostConfig = &container.HostConfig{
			Resources: container.Resources{
				Memory:     memoryLimitBytes,
				MemorySwap: memoryLimitBytes,
			},
		}
	}

	resp, err := c.docker.ContainerCreate(ctx,
		&container.Config{
			Image:        imageName,
			WorkingDir:   "/app",
			Cmd:          []string{"tail", "-f", "/dev/null"},
			AttachStdin:  true,
			OpenStdin:    true,
			StdinOnce:    false,
			Tty:          false,
			AttachStdout: true,
			AttachStderr: true,
		},
		hostConfig,
		nil,
		nil,
		"",
	)
	if err != nil {
		return "", fmt.Errorf("docker.ContainerCreate: %w", err)
	}

	err = c.docker.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return "", fmt.Errorf("docker.ContainerStart: %w", err)
	}

	return resp.ID, nil
}

// RemoveContainer .
func (c *Client) RemoveContainer(ctx context.Context, containerID string) error {
	return c.docker.ContainerRemove(ctx, containerID, container.RemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
}
