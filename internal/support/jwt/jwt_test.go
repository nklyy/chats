package jwt_test

import (
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"noname-realtime-support-chat/internal/support/jwt"
	"testing"
)

func TestNewJwtService(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	secretKey := "key"
	expiry := 1 // minutes

	tests := []struct {
		name        string
		secretKey   string
		expiry      *int
		redisClient *redis.Client
		expect      func(*testing.T, jwt.Service, error)
	}{
		{
			name:        "should return service",
			secretKey:   secretKey,
			expiry:      &expiry,
			redisClient: &redis.Client{},
			expect: func(t *testing.T, s jwt.Service, err error) {
				assert.NotNil(t, s)
				assert.Nil(t, err)
			},
		},
		{
			name:        "should return invalid jwt secret key",
			secretKey:   "",
			expiry:      &expiry,
			redisClient: &redis.Client{},
			expect: func(t *testing.T, s jwt.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid jwt secret key")
			},
		},
		{
			name:        "should return invalid jwt expiry",
			secretKey:   secretKey,
			expiry:      nil,
			redisClient: &redis.Client{},
			expect: func(t *testing.T, s jwt.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid jwt expiry")
			},
		},
		{
			name:        "should return invalid jwt expiry",
			secretKey:   secretKey,
			expiry:      &expiry,
			redisClient: nil,
			expect: func(t *testing.T, s jwt.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid redis client")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := jwt.NewJwtService(tc.secretKey, tc.expiry, tc.redisClient)
			tc.expect(t, svc, err)
		})
	}
}
