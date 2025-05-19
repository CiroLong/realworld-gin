package config

// Config holds application configuration
// Use github.com/spf13/viper for env & file support
type Config struct {
	Server   ServerConfig
	Database DBConfig
	JWT      JWTConfig
}

type ServerConfig struct {
}
type DBConfig struct {
}
type JWTConfig struct{}

func Load() *Config {
	// initialize Viper, read env and config.yaml
	// return populated Config
	return &Config{}
}
