package config_test

import (
	"noname-one-time-session-chat/config"
	"os"
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	type env struct {
		port        string
		environment string
		mongoDbName string
		mongoDbUrl  string
		salt        string
		redisHost   string
		redisPort   string
	}

	type args struct {
		env env
	}

	setEnv := func(env env) {
		os.Setenv("PORT", env.port)
		os.Setenv("APP_ENV", env.environment)
		os.Setenv("MONGO_DB_NAME", env.mongoDbName)
		os.Setenv("MONGO_DB_URL", env.mongoDbUrl)
		os.Setenv("SALT", env.salt)
		os.Setenv("REDIS_HOST", env.redisHost)
		os.Setenv("REDIS_PORT", env.redisPort)
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
					port:        "5000",
					environment: "development",
					mongoDbName: "example",
					mongoDbUrl:  "http://127.0.0.1",
					salt:        "salt",
					redisHost:   "localhost",
					redisPort:   "1234",
				},
			},
			want: &config.Config{
				PORT:        "5000",
				Environment: "development",
				Salt:        "salt",
				MongoDb: config.MongoDb{
					MongoDbName: "example",
					MongoDbUrl:  "http://127.0.0.1",
				},
				Redis: config.Redis{
					RedisHost: "localhost",
					RedisPort: "1234",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setEnv(test.args.env)

			got, err := config.Get()
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
