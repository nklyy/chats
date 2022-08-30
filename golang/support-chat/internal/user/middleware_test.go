package user_test

import (
	"support-chat/internal/user"
	mock_user "support-chat/internal/user/mocks"
	"support-chat/pkg/jwt"
	mock_jwt "support-chat/pkg/jwt/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"testing"
)

func TestNewMiddleware(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tests := []struct {
		name    string
		jwtSvc  jwt.Service
		userSvc user.Service
		logger  *zap.SugaredLogger
		expect  func(*testing.T, user.Middleware, error)
	}{
		{
			name:    "should return middleware",
			jwtSvc:  mock_jwt.NewMockService(controller),
			userSvc: mock_user.NewMockService(controller),
			logger:  &zap.SugaredLogger{},
			expect: func(t *testing.T, m user.Middleware, err error) {
				assert.NotNil(t, m)
				assert.Nil(t, err)
			},
		},
		{
			name:    "should return invalid jwt service",
			jwtSvc:  nil,
			userSvc: mock_user.NewMockService(controller),
			logger:  &zap.SugaredLogger{},
			expect: func(t *testing.T, m user.Middleware, err error) {
				assert.Nil(t, m)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[user_middleware] invalid jwt service")
			},
		},
		{
			name:    "should return invalid user service",
			jwtSvc:  mock_jwt.NewMockService(controller),
			userSvc: nil,
			logger:  &zap.SugaredLogger{},
			expect: func(t *testing.T, m user.Middleware, err error) {
				assert.Nil(t, m)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[user_middleware] invalid user service")
			},
		},
		{
			name:    "should return invalid logger",
			jwtSvc:  mock_jwt.NewMockService(controller),
			userSvc: mock_user.NewMockService(controller),
			logger:  nil,
			expect: func(t *testing.T, m user.Middleware, err error) {
				assert.Nil(t, m)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[user_middleware] invalid logger")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := user.NewMiddleware(tc.jwtSvc, tc.userSvc, tc.logger)
			tc.expect(t, svc, err)
		})
	}
}
