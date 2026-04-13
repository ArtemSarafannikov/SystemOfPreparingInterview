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
	Env    Environment
	Server ServerConfig
}

// ServerConfig .
type ServerConfig struct {
	Host                    string        `mapstructure:"host"`
	GraphQLPort             int           `mapstructure:"gql_port"`
	GracefulShutdownTimeout time.Duration `mapstructure:"graceful_shutdown_timeout"`
	CORSAllowedOrigins      []string      `mapstructure:"cors_allowed_origins"`
	// WebsocketAllowedHosts is matched against http.Request.Host for GraphQL websocket upgrades.
	WebsocketAllowedHosts []string `mapstructure:"websocket_allowed_hosts"`
}
