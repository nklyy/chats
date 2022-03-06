package support_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"noname-realtime-support-chat/internal/support"
	mock_support "noname-realtime-support-chat/internal/support/mocks"
	"noname-realtime-support-chat/pkg/jwt"
	mock_jwt "noname-realtime-support-chat/pkg/jwt/mocks"
	"testing"
)

func TestNewService(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	salt := 10

	tests := []struct {
		name       string
		repository support.Repository
		logger     *zap.SugaredLogger
		salt       *int
		jwtSvc     jwt.Service
		expect     func(*testing.T, support.Service, error)
	}{
		{
			name:       "should return service",
			repository: mock_support.NewMockRepository(controller),
			logger:     &zap.SugaredLogger{},
			salt:       &salt,
			jwtSvc:     mock_jwt.NewMockService(controller),
			expect: func(t *testing.T, s support.Service, err error) {
				assert.NotNil(t, s)
				assert.Nil(t, err)
			},
		},
		{
			name:       "should return invalid repository",
			repository: nil,
			logger:     &zap.SugaredLogger{},
			salt:       &salt,
			jwtSvc:     mock_jwt.NewMockService(controller),
			expect: func(t *testing.T, s support.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid repository")
			},
		},
		{
			name:       "should return invalid logger",
			repository: mock_support.NewMockRepository(controller),
			logger:     nil,
			salt:       &salt,
			jwtSvc:     mock_jwt.NewMockService(controller),
			expect: func(t *testing.T, s support.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid logger")
			},
		},
		{
			name:       "should return invalid salt",
			repository: mock_support.NewMockRepository(controller),
			logger:     &zap.SugaredLogger{},
			salt:       nil,
			jwtSvc:     mock_jwt.NewMockService(controller),
			expect: func(t *testing.T, s support.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid salt")
			},
		},
		{
			name:       "should return invalid jwt service",
			repository: mock_support.NewMockRepository(controller),
			logger:     &zap.SugaredLogger{},
			salt:       &salt,
			jwtSvc:     nil,
			expect: func(t *testing.T, s support.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid jwt service")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := support.NewService(tc.repository, tc.logger, tc.salt, tc.jwtSvc)
			tc.expect(t, svc, err)
		})
	}
}
