package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	PORT             string `mapstructure:"PORT"`
	Environment      string `mapstructure:"APP_ENV"`
	MongoDbName      string `mapstructure:"MONGO_DB_NAME"`
	MongoDbUrl       string `mapstructure:"MONGO_DB_URL"`
	Salt             int    `mapstructure:"SALT"`
	JwtSecretAccess  string `mapstructure:"JWT_SECRET_ACCESS"`
	JwtExpiryAccess  int    `mapstructure:"JWT_EXPIRY_ACCESS"`
	JwtSecretRefresh string `mapstructure:"JWT_SECRET_REFRESH"`
	JwtExpiryRefresh int    `mapstructure:"JWT_EXPIRY_REFRESH"`
	AutoLogout       int    `mapstructure:"AUTO_LOGOUT"`
	RedisHost        string `mapstructure:"REDIS_HOST"`
	RedisPortAuth    string `mapstructure:"REDIS_PORT_AUTH"`
	RedisPortChat    string `mapstructure:"REDIS_PORT_CHAT"`
	RabbitMqUrl      string `mapstructure:"AMQP_SERVER_URL"`
}

func Get(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	//viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var configuration Config
	err = viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	return &configuration, nil
}
