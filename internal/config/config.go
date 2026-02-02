// internal/config/config.go
package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds application configuration
// Use github.com/spf13/viper for env & file support
type Config struct {
	Server   ServerConfig `mapstructure:"server"`
	Database MySQLConfig  `mapstructure:"database"`
	JWT      JWTConfig    `mapstructure:"jwt"`
}

type ServerConfig struct {
	Addr         string        `mapstructure:"addr"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type MySQLConfig struct {
	DSN             string        `mapstructure:"dsn"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}
type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	ExpireTime time.Duration `mapstructure:"expire_time"`
}

// 这里是一个全局变量，只提供一个Getter
var cfg *Config

// C returns global config
func C() *Config {
	if cfg == nil {
		panic("config not initialized")
	}
	return cfg
}

// Load loads config from config.yaml and environment variables
func Load() error {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("./configs")

	v.SetEnvPrefix("APP")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return err
	}

	cfg = &c
	return nil
}
