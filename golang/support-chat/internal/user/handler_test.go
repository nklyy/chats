package user_test

import (
	"support-chat/internal/user"
	mock_user "support-chat/internal/user/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tests := []struct {
		name    string
		userSvc user.Service
		expect  func(*testing.T, *user.Handler, error)
	}{
		{
			name:    "should return service",
			userSvc: mock_user.NewMockService(controller),
			expect: func(t *testing.T, s *user.Handler, err error) {
				assert.NotNil(t, s)
				assert.Nil(t, err)
			},
		},
		{
			name:    "should return invalid user service",
			userSvc: nil,
			expect: func(t *testing.T, s *user.Handler, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[user_handler] invalid user service")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := user.NewHandler(tc.userSvc)
			tc.expect(t, svc, err)
		})
	}
}
