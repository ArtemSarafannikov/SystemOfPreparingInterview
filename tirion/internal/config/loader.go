package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Load .
func Load(configPath string) (Config, error) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	viper.BindEnv("database.db_dsn", "DB_DSN")

	config := Config{}
	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("viper.ReadInConfig: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("viper.Unmarshal: %w", err)
	}
	return config, nil
}

// MustLoad .
func MustLoad(configPath string) Config {
	cfg, err := Load(configPath)
	if err != nil {
		panic(err)
	}
	return cfg
}

// MustLoadFromFlag .
func MustLoadFromFlag() Config {
	configPath := pflag.String("config", "", "Filepath to config")
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}

	return MustLoad(*configPath)
}
