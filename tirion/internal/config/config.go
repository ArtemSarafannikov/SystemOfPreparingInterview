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
}

// ServerConfig .
type ServerConfig struct {
	Host                    string        `mapstructure:"host"`
	GRPCPort                int           `mapstructure:"grpc_port"`
	GatewayPort             int           `mapstructure:"http_port"`
	GracefulShutdownTimeout time.Duration `mapstructure:"graceful_shutdown_timeout"`
}
