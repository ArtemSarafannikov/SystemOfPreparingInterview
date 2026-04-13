package config

import (
	"time"
)

// Environment .
type Environment string

// .
const (
	Local      Environment = "dev"
	Stage      Environment = "stg"
	Production Environment = "prod"
)

// Config .
type Config struct {
	Env      Environment
	Server   ServerConfig
	Database DatabaseConfig
	Judge    JudgeConfig `mapstructure:"judge"`
}

// JudgeConfig limits sandbox containers and the River worker pool (see config example YAML).
type JudgeConfig struct {
	// SandboxMemoryOverheadMB is added to the problem memory_limit_kb when setting Docker cgroup Memory
	// so the interpreter and base image fit; user solution is still bounded by the problem limit in aggregate RSS.
	SandboxMemoryOverheadMB int64 `mapstructure:"sandbox_memory_overhead_mb"`
	MaxRiverWorkers         int   `mapstructure:"max_river_workers"`
	// MaxConcurrentSandboxes caps concurrent Docker sandboxes (0 = same as MaxRiverWorkers).
	MaxConcurrentSandboxes int `mapstructure:"max_concurrent_sandboxes"`
}

// ServerConfig .
type ServerConfig struct {
	Host                    string        `mapstructure:"host"`
	GRPCPort                int           `mapstructure:"grpc_port"`
	GatewayPort             int           `mapstructure:"http_port"`
	GracefulShutdownTimeout time.Duration `mapstructure:"graceful_shutdown_timeout"`
}
