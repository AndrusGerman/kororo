package config

import (
	"os"

	"kororo/internal/core/domain/types"
)

type Config struct {
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Host() string {
	return ":" + c.Port()
}

func (c *Config) Port() string {
	return c.getValue("PORT", "4500")
}

func (c *Config) Database() types.Database {
	return "KororoDB"
}

func (c *Config) UriMongo() string {
	return "mongodb://user:passwod@localhost:27017/"
}

func (c *Config) getValue(envName string, defaultValue string) string {
	if os.Getenv(envName) != "" {
		return os.Getenv(envName)
	}
	return defaultValue
}
