package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"sync"
)

type Config struct {
	PORT        string `required:"true" default:"5000" envconfig:"PORT"`
	Environment string `required:"true" envconfig:"APP_ENV"`
	Salt        string `required:"true" envconfig:"SALT"`
	MongoDb
	Redis
}

type MongoDb struct {
	MongoDbName string `required:"true" envconfig:"MONGO_DB_NAME"`
	MongoDbUrl  string `required:"true" envconfig:"MONGO_DB_URL"`
}

type Redis struct {
	RedisHost string `required:"true" envconfig:"REDIS_HOST"`
	RedisPort string `required:"true" envconfig:"REDIS_PORT"`
}

var (
	once   sync.Once
	config *Config
)

func Get() (*Config, error) {
	var err error
	once.Do(func() {
		var cfg Config
		// If you run it locally and through terminal please set up this in Load function (../.env)
		_ = godotenv.Load(".env")

		if err = envconfig.Process("", &cfg); err != nil {
			return
		}

		config = &cfg
	})

	return config, err
}
