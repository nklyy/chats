package support_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"noname-realtime-support-chat/internal/support"
	"noname-realtime-support-chat/internal/support/jwt"
	"noname-realtime-support-chat/internal/support/jwt/mocks"
	"testing"
)

func TestNewMiddleware(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tests := []struct {
		name   string
		jwtSvc jwt.Service
		logger *zap.SugaredLogger
		expect func(*testing.T, support.Middleware, error)
	}{
		{
			name:   "should return middleware",
			jwtSvc: mock_jwt.NewMockService(controller),
			logger: &zap.SugaredLogger{},
			expect: func(t *testing.T, m support.Middleware, err error) {
				assert.NotNil(t, m)
				assert.Nil(t, err)
			},
		},
		{
			name:   "should return invalid jwt service",
			jwtSvc: nil,
			logger: &zap.SugaredLogger{},
			expect: func(t *testing.T, m support.Middleware, err error) {
				assert.Nil(t, m)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid jwt service")
			},
		},
		{
			name:   "should return invalid logger",
			jwtSvc: mock_jwt.NewMockService(controller),
			logger: nil,
			expect: func(t *testing.T, m support.Middleware, err error) {
				assert.Nil(t, m)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid logger")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := support.NewMiddleware(tc.jwtSvc, tc.logger)
			tc.expect(t, svc, err)
		})
	}
}
