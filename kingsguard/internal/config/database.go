package config

import "fmt"

type DatabaseConfig struct {
	Host               string `mapstructure:"host"`
	Port               int    `mapstructure:"port"`
	User               string `mapstructure:"user"`
	Password           string `mapstructure:"password"`
	Name               string `mapstructure:"name"`
	SslMode            string `mapstructure:"ssl_mode"`
	MaxOpenConnections int    `mapstructure:"max_open_connections"`
	Dsn                string `mapstructure:"db_dsn"`
}

func (cfg *DatabaseConfig) ConnString() string {
	if cfg.Dsn != "" {
		return cfg.Dsn
	}

	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s pool_max_conns=%d",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SslMode,
		cfg.MaxOpenConnections,
	)
}
