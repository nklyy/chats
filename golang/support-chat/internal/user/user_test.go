package user_test

import (
	"support-chat/internal/user"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
		expect   func(*testing.T, *user.User, error)
	}{
		{
			testName: "should return user",
			email:    "email",
			name:     "name",
			password: "password",
			salt:     &salt,
			expect: func(t *testing.T, support *user.User, err error) {
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
			expect: func(t *testing.T, s *user.User, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[user] invalid email")
			},
		},
		{
			testName: "should return name error",
			email:    "email",
			name:     "",
			password: "password",
			salt:     &salt,
			expect: func(t *testing.T, s *user.User, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[user] invalid name")
			},
		},
		{
			testName: "should return password error",
			email:    "email",
			name:     "name",
			password: "",
			salt:     &salt,
			expect: func(t *testing.T, s *user.User, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[user] invalid password")
			},
		},
		{
			testName: "should return salt error",
			email:    "email",
			name:     "name",
			password: "password",
			salt:     nil,
			expect: func(t *testing.T, s *user.User, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[user] invalid salt")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			svc, err := user.NewUser(tc.email, tc.name, tc.password, tc.salt)
			tc.expect(t, svc, err)
		})
	}
}
