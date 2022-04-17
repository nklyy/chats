package jwt_test

import (
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"noname-support-chat/pkg/jwt"
	"testing"
)

func TestNewJwtService(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	secretAccessKey := "key"
	secretRefreahsKey := "key"
	expiryAccess := 1  // minutes
	expiryRefresh := 2 // minutes
	autoLogout := 3    // minutes

	tests := []struct {
		name             string
		secretKeyAccess  string
		expiryAccess     *int
		secretKeyRefresh string
		expiryRefresh    *int
		autoLogout       *int
		redisClient      *redis.Client
		expect           func(*testing.T, jwt.Service, error)
	}{
		{
			name:             "should return service",
			secretKeyAccess:  secretAccessKey,
			expiryAccess:     &expiryAccess,
			secretKeyRefresh: secretRefreahsKey,
			expiryRefresh:    &expiryRefresh,
			autoLogout:       &autoLogout,
			redisClient:      &redis.Client{},
			expect: func(t *testing.T, s jwt.Service, err error) {
				assert.NotNil(t, s)
				assert.Nil(t, err)
			},
		},
		{
			name:             "should return invalid jwt access secret key",
			secretKeyAccess:  "",
			expiryAccess:     &expiryAccess,
			secretKeyRefresh: secretRefreahsKey,
			expiryRefresh:    &expiryRefresh,
			autoLogout:       &autoLogout,
			redisClient:      &redis.Client{},
			expect: func(t *testing.T, s jwt.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid jwt access secret key")
			},
		},
		{
			name:             "should return invalid jwt expiry access",
			secretKeyAccess:  secretAccessKey,
			expiryAccess:     nil,
			secretKeyRefresh: secretRefreahsKey,
			expiryRefresh:    &expiryRefresh,
			autoLogout:       &autoLogout,
			redisClient:      &redis.Client{},
			expect: func(t *testing.T, s jwt.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid jwt expiry access")
			},
		},
		{
			name:             "should return invalid jwt refresh secret key",
			secretKeyAccess:  secretAccessKey,
			expiryAccess:     &expiryAccess,
			secretKeyRefresh: "",
			expiryRefresh:    &expiryRefresh,
			autoLogout:       &autoLogout,
			redisClient:      &redis.Client{},
			expect: func(t *testing.T, s jwt.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid jwt refresh secret key")
			},
		},
		{
			name:             "should return invalid jwt expiry refresh",
			secretKeyAccess:  secretAccessKey,
			expiryAccess:     &expiryAccess,
			secretKeyRefresh: secretRefreahsKey,
			expiryRefresh:    nil,
			autoLogout:       &autoLogout,
			redisClient:      &redis.Client{},
			expect: func(t *testing.T, s jwt.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid jwt expiry refresh")
			},
		},
		{
			name:             "should return invalid jwt auto logout",
			secretKeyAccess:  secretAccessKey,
			expiryAccess:     &expiryAccess,
			secretKeyRefresh: secretRefreahsKey,
			expiryRefresh:    &expiryRefresh,
			autoLogout:       nil,
			redisClient:      &redis.Client{},
			expect: func(t *testing.T, s jwt.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid jwt auto logout")
			},
		},
		{
			name:             "should return invalid redis client",
			secretKeyAccess:  secretAccessKey,
			expiryAccess:     &expiryAccess,
			secretKeyRefresh: secretRefreahsKey,
			expiryRefresh:    &expiryRefresh,
			autoLogout:       &autoLogout,
			redisClient:      nil,
			expect: func(t *testing.T, s jwt.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid redis client")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := jwt.NewJwtService(tc.secretKeyAccess,
				tc.expiryAccess,
				tc.secretKeyRefresh,
				tc.expiryRefresh,
				tc.autoLogout,
				tc.redisClient)
			tc.expect(t, svc, err)
		})
	}
}
