package config

import (
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PORT        string `required:"true" default:"5000" envconfig:"PORT"`
	Environment string `required:"true" envconfig:"APP_ENV"`
	Salt        int    `required:"true" envconfig:"SALT"`
	MongoDb
	Jwt
	Redis
}

type MongoDb struct {
	MongoDbName string `required:"true" envconfig:"MONGO_DB_NAME"`
	MongoDbUrl  string `required:"true" envconfig:"MONGO_DB_URL"`
}

type Jwt struct {
	JwtSecretAccess  string `required:"true" envconfig:"JWT_SECRET_ACCESS"`
	JwtExpiryAccess  int    `required:"true" envconfig:"JWT_EXPIRY_ACCESS"`
	JwtSecretRefresh string `required:"true" envconfig:"JWT_SECRET_REFRESH"`
	JwtExpiryRefresh int    `required:"true" envconfig:"JWT_EXPIRY_REFRESH"`
	AutoLogout       int    `required:"true" envconfig:"AUTO_LOGOUT"`
}

type Redis struct {
	RedisHostAuth string `required:"true" envconfig:"REDIS_HOST_AUTH"`
	RedisPortAuth string `required:"true" envconfig:"REDIS_PORT_AUTH"`
	RedisHostChat string `required:"true" envconfig:"REDIS_HOST_CHAT"`
	RedisPortChat string `required:"true" envconfig:"REDIS_PORT_CHAT"`
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
