// internal/config/config.go
package config

import (
	"github.com/spf13/viper"
	"strings"
	"time"
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

func setDefaults(v *viper.Viper) {
	// server
	v.SetDefault("server.addr", ":8080")
	v.SetDefault("server.read_timeout", "5s")
	v.SetDefault("server.write_timeout", "5s")

	// database
	v.SetDefault("database.max_open_conns", 50)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", "1h")

	// jwt
	v.SetDefault("jwt.expire_time", "24h")
}
