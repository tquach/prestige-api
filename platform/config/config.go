package config

import (
	"os"

	"github.com/tquach/prestige-api/platform/redis"
)

// Env is the current environment
var Env string

// EnvType is the current env
type EnvType string

// Environment enums
const (
	Local = "local"
	Dev   = "dev"
	Test  = "test"
	Prod  = "prod"
)

// EnvSettings holds environment settings.
type EnvSettings struct {
	Name          string               `yaml:"name"`
	Debug         bool                 `yaml:"debug"`
	SecretKey     string               `yaml:"secret_key"`
	AppSettings   AppSettings          `yaml:"app"`
	Database      DatabaseSettings     `yaml:"database"`
	RedisSettings redis.ServerSettings `yaml:"redis"`
}

// AppSettings contains application settings
type AppSettings struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
}

// DatabaseSettings contain database connection settings
type DatabaseSettings struct {
	Driver   string `yaml:"driver"`
	URL      string `yaml:"url"`
	DBName   string `yaml:"dbname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	SSLMode  bool   `yaml:"sslmode"`
}

func init() {
	Env = os.Getenv("APP_ENV")
}
