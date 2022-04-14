package user_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"noname-realtime-support-chat/internal/chat/user"
	"noname-realtime-support-chat/pkg/errors"
	"testing"
)

func TestNewSupport(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	salt := "salt"

	tests := []struct {
		testName string
		ipAddr   string
		salt     string
		expect   func(*testing.T, *user.User, error)
	}{
		{
			testName: "should return user",
			ipAddr:   "127.0.0.1",
			salt:     salt,
			expect: func(t *testing.T, support *user.User, err error) {
				assert.NotNil(t, support)
				assert.Nil(t, err)
			},
		},
		{
			testName: "should return email error",
			ipAddr:   "",
			salt:     salt,
			expect: func(t *testing.T, s *user.User, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, errors.WithMessage(user.ErrInvalidIpAddress, "should be not empty").Error())
			},
		},
		{
			testName: "should return password error",
			ipAddr:   "127.0.0.1",
			salt:     "",
			expect: func(t *testing.T, s *user.User, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, errors.WithMessage(user.ErrInvalidSalt, "should be not empty").Error())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			svc, err := user.NewUser(tc.ipAddr, tc.salt)
			tc.expect(t, svc, err)
		})
	}
}
