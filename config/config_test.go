package config_test

import (
	"noname-realtime-support-chat/config"
	"os"
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	type env struct {
		port             string
		environment      string
		mongoDbName      string
		mongoDbUrl       string
		salt             string
		jwtSecretAccess  string
		jwtExpiryAccess  string
		jwtSecretRefresh string
		jwtExpiryRefresh string
		autoLogout       string
		redisHost        string
		redisPortAuth    string
		redisPortChat    string
		rabbitMqUrl      string
	}

	type args struct {
		env env
	}

	setEnv := func(env env) {
		os.Setenv("PORT", env.port)
		os.Setenv("ENVIRONMENT", env.environment)
		os.Setenv("MONGO_DB_NAME", env.mongoDbName)
		os.Setenv("MONGO_DB_URL", env.mongoDbUrl)
		os.Setenv("SALT", env.salt)
		os.Setenv("JWT_SECRET_ACCESS", env.jwtSecretAccess)
		os.Setenv("JWT_EXPIRY_ACCESS", env.jwtExpiryAccess)
		os.Setenv("JWT_SECRET_REFRESH", env.jwtSecretRefresh)
		os.Setenv("JWT_EXPIRY_REFRESH", env.jwtExpiryRefresh)
		os.Setenv("AUTO_LOGOUT", env.autoLogout)
		os.Setenv("REDIS_HOST", env.redisHost)
		os.Setenv("REDIS_PORT_AUTH", env.redisPortAuth)
		os.Setenv("REDIS_PORT_CHAT", env.redisPortChat)
		os.Setenv("AMQP_SERVER_URL", env.rabbitMqUrl)
	}

	tests := []struct {
		name      string
		args      args
		want      *config.Config
		wantError bool
	}{
		{
			name: "Test config file!",
			args: args{
				env: env{
					port:             ":5000",
					environment:      "development",
					mongoDbName:      "example",
					mongoDbUrl:       "http://127.0.0.1",
					salt:             "11",
					jwtSecretAccess:  "jwt",
					jwtExpiryAccess:  "100",
					jwtSecretRefresh: "asd",
					jwtExpiryRefresh: "300",
					autoLogout:       "3",
					redisHost:        "localhost",
					redisPortAuth:    "1234",
					redisPortChat:    "4321",
					rabbitMqUrl:      "amqp://guest:guest@message-broker:5672",
				},
			},
			want: &config.Config{
				PORT:             ":5000",
				Environment:      "development",
				MongoDbName:      "example",
				MongoDbUrl:       "http://127.0.0.1",
				Salt:             11,
				JwtSecretAccess:  "jwt",
				JwtExpiryAccess:  100,
				JwtSecretRefresh: "asd",
				JwtExpiryRefresh: 300,
				AutoLogout:       3,
				RedisHost:        "localhost",
				RedisPortAuth:    "1234",
				RedisPortChat:    "4321",
				RabbitMqUrl:      "amqp://guest:guest@message-broker:5672",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setEnv(test.args.env)

			got, err := config.Get("..")
			if (err != nil) != test.wantError {
				t.Errorf("Init() error = %v, wantErr %v", err, test.wantError)

				return
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Init() got = %v, want %v", got, test.want)
			}
		})
	}
}
