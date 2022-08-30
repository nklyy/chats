package auth_test

import (
	"support-chat/internal/user/auth"
	mock_auth "support-chat/internal/user/auth/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tests := []struct {
		name    string
		authSvc auth.Service
		expect  func(*testing.T, *auth.Handler, error)
	}{
		{
			name:    "should return service",
			authSvc: mock_auth.NewMockService(controller),
			expect: func(t *testing.T, s *auth.Handler, err error) {
				assert.NotNil(t, s)
				assert.Nil(t, err)
			},
		},
		{
			name:    "should return invalid auth service",
			authSvc: nil,
			expect: func(t *testing.T, s *auth.Handler, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[chat_auth_handler] invalid auth service")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := auth.NewHandler(tc.authSvc)
			tc.expect(t, svc, err)
		})
	}
}
