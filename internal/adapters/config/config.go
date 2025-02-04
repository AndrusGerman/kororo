package config

import (
	"os"

	"kororo/internal/core/domain/types"

	"github.com/joho/godotenv"
)

type Config struct {
}

func NewConfig() *Config {

	godotenv.Load()
	return &Config{}
}

func (c *Config) Host() string {
	return ":" + c.Port()
}

func (c *Config) Debug() bool {
	return true
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

func (c *Config) GEMINI_API_KEY() string {
	return c.getValue("GEMINI_API_KEY", "")
}

func (c *Config) HUGGINGFACE_API_KEY() string {
	return c.getValue("HUGGINGFACE_API_KEY", "")
}

func (c *Config) getValue(envName string, defaultValue string) string {
	if os.Getenv(envName) != "" {

		return os.Getenv(envName)
	}
	return defaultValue
}
