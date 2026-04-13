package config

import (
	"time"
)

type Environment string

const (
	Local      Environment = "dev"
	Stage      Environment = "stg"
	Production Environment = "prod"
)

type Config struct {
	Env      Environment
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Host                    string        `mapstructure:"host"`
	GRPCPort                int           `mapstructure:"grpc_port"`
	GatewayPort             int           `mapstructure:"http_port"`
	GracefulShutdownTimeout time.Duration `mapstructure:"graceful_shutdown_timeout"`
}
