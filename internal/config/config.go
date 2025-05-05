package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server          ServerConfig
	Database        DatabaseConfig
	APILogin        string
	TokenCacheKey   string
	TokenTimeout    time.Duration
	APITimeout      time.Duration
	WaiterAPIURL    string
	WaiterAPIKey    string
	DefaultWaiterID string
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (c *DatabaseConfig) PostgresURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

func NewDefaultConfig() *Config {
	return &Config{
		APILogin:        "default_api_login",
		TokenCacheKey:   "iiko_token",
		TokenTimeout:    time.Hour,
		APITimeout:      time.Second * 15,
		WaiterAPIURL:    "https://api.waiter.iiko.ru",
		WaiterAPIKey:    "default_waiter_api_key",
		DefaultWaiterID: "default_waiter_id",
	}
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("ошибка чтения конфигурационного файла: %w", err)
	}

	var config Config

	config.Server.Port = viper.GetString("server.port")
	config.Server.ReadTimeout = viper.GetDuration("server.read_timeout")
	config.Server.WriteTimeout = viper.GetDuration("server.write_timeout")

	config.Database.Host = viper.GetString("database.host")
	config.Database.Port = viper.GetString("database.port")
	config.Database.User = viper.GetString("database.user")
	config.Database.Password = viper.GetString("database.password")
	config.Database.DBName = viper.GetString("database.dbname")
	config.Database.SSLMode = viper.GetString("database.sslmode")

	config.APILogin = viper.GetString("iiko.api_login")
	config.TokenCacheKey = viper.GetString("iiko.token_cache_key")
	config.TokenTimeout = viper.GetDuration("iiko.token_timeout")
	config.APITimeout = viper.GetDuration("iiko.api_timeout")
	config.WaiterAPIURL = viper.GetString("iiko.waiter_api_url")
	config.WaiterAPIKey = viper.GetString("iiko.waiter_api_key")
	config.DefaultWaiterID = viper.GetString("iiko.default_waiter_id")

	return &config, nil
}
