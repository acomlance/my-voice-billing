package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type DBConnector struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SQLDebug        bool          `mapstructure:"sql_debug"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
	Dir   string `mapstructure:"dir"`
}

type GrpcConfig struct {
	Port int `mapstructure:"port"`
}

type Config struct {
	Environment string                 `mapstructure:"environment"`
	Grpc        GrpcConfig             `mapstructure:"grpc"`
	Database    map[string]DBConnector `mapstructure:"database"`
	Log         LogConfig              `mapstructure:"log"`
}

func (c *Config) Connector(alias string) (DBConnector, bool) {
	if c.Database == nil {
		return DBConnector{}, false
	}
	conn, ok := c.Database[alias]
	return conn, ok
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	if cfg.Grpc.Port <= 0 {
		cfg.Grpc.Port = 50051
	}
	return &cfg, nil
}

func Validate(c *Config) error {
	if c.Grpc.Port <= 0 {
		return fmt.Errorf("grpc.port must be positive")
	}
	if _, ok := c.Connector("write"); !ok {
		return fmt.Errorf("database connector 'write' not found")
	}
	return nil
}
