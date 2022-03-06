package support

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"noname-realtime-support-chat/pkg/errors"
	"testing"
)

func TestNewSupport(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	salt := 10

	tests := []struct {
		testName string
		email    string
		name     string
		password string
		salt     *int
		expect   func(*testing.T, *Support, error)
	}{
		{
			testName: "should return support",
			email:    "email",
			name:     "name",
			password: "password",
			salt:     &salt,
			expect: func(t *testing.T, support *Support, err error) {
				assert.NotNil(t, support)
				assert.Nil(t, err)
			},
		},
		{
			testName: "should return email error",
			email:    "",
			name:     "name",
			password: "password",
			salt:     &salt,
			expect: func(t *testing.T, support *Support, err error) {
				assert.Nil(t, support)
				assert.NotNil(t, err)
				assert.EqualError(t, err, errors.WithMessage(ErrInvalidEmail, "should be not empty").Error())
			},
		},
		{
			testName: "should return name error",
			email:    "email",
			name:     "",
			password: "password",
			salt:     &salt,
			expect: func(t *testing.T, support *Support, err error) {
				assert.Nil(t, support)
				assert.NotNil(t, err)
				assert.EqualError(t, err, errors.WithMessage(ErrInvalidName, "should be not empty").Error())
			},
		},
		{
			testName: "should return password error",
			email:    "email",
			name:     "name",
			password: "",
			salt:     &salt,
			expect: func(t *testing.T, support *Support, err error) {
				assert.Nil(t, support)
				assert.NotNil(t, err)
				assert.EqualError(t, err, errors.WithMessage(ErrInvalidPassword, "should be not empty").Error())
			},
		},
		{
			testName: "should return password error",
			email:    "email",
			name:     "name",
			password: "password",
			salt:     nil,
			expect: func(t *testing.T, support *Support, err error) {
				assert.Nil(t, support)
				assert.NotNil(t, err)
				assert.EqualError(t, err, errors.WithMessage(ErrInvalidSalt, "should be not empty").Error())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			svc, err := NewSupport(tc.email, tc.name, tc.password, tc.salt)
			tc.expect(t, svc, err)
		})
	}
}
