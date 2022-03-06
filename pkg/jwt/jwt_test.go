package jwt_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"noname-realtime-support-chat/pkg/jwt"
	"testing"
)

func TestNewJwtService(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	secretKey := "key"
	expiry := 1 // minutes

	tests := []struct {
		name      string
		secretKey string
		expiry    *int
		expect    func(*testing.T, jwt.Service, error)
	}{
		{
			name:      "should return service",
			secretKey: secretKey,
			expiry:    &expiry,
			expect: func(t *testing.T, s jwt.Service, err error) {
				assert.NotNil(t, s)
				assert.Nil(t, err)
			},
		},
		{
			name:      "should return invalid jwt secret key",
			secretKey: "",
			expiry:    &expiry,
			expect: func(t *testing.T, s jwt.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid jwt secret key")
			},
		},
		{
			name:      "should return invalid jwt expiry",
			secretKey: secretKey,
			expiry:    nil,
			expect: func(t *testing.T, s jwt.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid jwt expiry")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := jwt.NewJwtService(tc.secretKey, tc.expiry)
			tc.expect(t, svc, err)
		})
	}
}
