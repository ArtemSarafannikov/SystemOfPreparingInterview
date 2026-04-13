package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/CodefriendOrg/daenerys/internal/pkg/logger"
	"github.com/docker/docker/api/types/image"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// PullImage .
func (c *Client) PullImage(ctx context.Context, imageName string) error {
	images, err := c.docker.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return fmt.Errorf("docker.ImageList: %w", err)
	}

	_, exist := lo.Find(images, func(item image.Summary) bool {
		_, exist := lo.Find(item.RepoTags, func(tag string) bool {
			return tag == imageName
		})
		return exist
	})
	if exist {
		return nil
	}

	reader, err := c.docker.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("docker.ImagePull: %w", err)
	}
	defer reader.Close()

	_, err = io.Copy(io.Discard, reader)
	if err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	logger.Infof(ctx, "image pull successful", zap.String("image", imageName))

	return nil
}
