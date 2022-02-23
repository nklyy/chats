package config

import (
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
	}

	tests := []struct {
		name      string
		args      args
		want      *Config
		wantError bool
	}{
		{
			name: "Test config file!",
			args: args{
				env: env{
					port:        ":5000",
					environment: "development",
					mongoDbName: "example",
					mongoDbUrl:  "http://127.0.0.1",
					salt:        "11",
				},
			},
			want: &Config{
				PORT:        ":5000",
				Environment: "development",
				MongoDbName: "example",
				MongoDbUrl:  "http://127.0.0.1",
				Salt:        11,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setEnv(test.args.env)

			got, err := Get("..")
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
