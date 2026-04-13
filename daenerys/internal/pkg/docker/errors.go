package docker

import "errors"

var (
	TimeLimitExceeded   = errors.New("time limit exceeded")
	MemoryLimitExceeded = errors.New("memory limit exceeded")
	RuntimeError        = errors.New("runtime error")
)
