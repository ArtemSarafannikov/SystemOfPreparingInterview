package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/moby/moby/pkg/stdcopy"
)

// RunTest runs argv in the container with stdin wired to input (plus newline when non-empty).
func (c *Client) RunTest(ctx context.Context, containerID string, argv []string, input string, maxExecTime time.Duration) (Output, error) {
	output := Output{}

	execID, err := c.docker.ContainerExecCreate(ctx, containerID, container.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          argv,
	})
	if err != nil {
		return output, fmt.Errorf("ContainerExecCreate: %w", err)
	}

	resp, err := c.docker.ContainerExecAttach(ctx, execID.ID, container.ExecStartOptions{})
	if err != nil {
		return output, fmt.Errorf("ContainerExecAttach: %w", err)
	}
	defer resp.Close()

	ctxTimeout, cancel := context.WithTimeout(ctx, maxExecTime)
	defer cancel()

	var stdout, stderr bytes.Buffer
	doneCh := make(chan error, 1)

	go func() {
		if input != "" {
			_, err = resp.Conn.Write([]byte(input + "\n"))
			if err != nil {
				doneCh <- fmt.Errorf("write stdin: %w", err)
				return
			}
		}

		if closer, ok := resp.Conn.(interface{ CloseWrite() error }); ok {
			_ = closer.CloseWrite()
		}

		_, err = stdcopy.StdCopy(&stdout, &stderr, resp.Reader)
		doneCh <- err
	}()

	select {
	case <-ctxTimeout.Done():
		return output, TimeLimitExceeded
	case err = <-doneCh:
		if err != nil {
			return output, fmt.Errorf("io error: %w", err)
		}
	}

	inspect, err := c.docker.ContainerExecInspect(ctx, execID.ID)
	if err != nil {
		return output, fmt.Errorf("docker.ContainerExecInspect: %w", err)
	}

	output.Stdout = stdout.String()
	output.Stderr = stderr.String()

	if inspect.ExitCode == 137 {
		return output, MemoryLimitExceeded
	}

	if inspect.ExitCode != 0 || len(output.Stderr) > 0 {
		return output, RuntimeError
	}

	return output, nil
}

// CopyCodeToContainer .
func (c *Client) CopyCodeToContainer(ctx context.Context, containerID, dir, filename, code string) error {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	hdr := &tar.Header{
		Name: filename,
		Mode: 0644,
		Size: int64(len(code)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return fmt.Errorf("tw.WriteHeader: %w", err)
	}

	if _, err := tw.Write([]byte(code)); err != nil {
		return fmt.Errorf("tw.Write: %w", err)
	}
	if err := tw.Close(); err != nil {
		return fmt.Errorf("tw.Close: %w", err)
	}

	return c.docker.CopyToContainer(ctx, containerID, dir, buf, container.CopyToContainerOptions{})
}
