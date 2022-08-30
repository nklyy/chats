package room_test

import (
	"support-chat/internal/chat/room"
	mock_room "support-chat/internal/chat/room/mocks"
	"support-chat/internal/user"
	mock_user "support-chat/internal/user/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewService(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tests := []struct {
		name       string
		repository room.Repository
		userSvc    user.Service
		logger     *zap.SugaredLogger
		expect     func(*testing.T, room.Service, error)
	}{
		{
			name:       "should return service",
			repository: mock_room.NewMockRepository(controller),
			userSvc:    mock_user.NewMockService(controller),
			logger:     &zap.SugaredLogger{},
			expect: func(t *testing.T, s room.Service, err error) {
				assert.NotNil(t, s)
				assert.Nil(t, err)
			},
		},
		{
			name:       "should return invalid repository",
			repository: nil,
			userSvc:    mock_user.NewMockService(controller),
			logger:     &zap.SugaredLogger{},
			expect: func(t *testing.T, s room.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[chat_room_service] invalid repository")
			},
		},
		{
			name:       "should return invalid user service",
			repository: mock_room.NewMockRepository(controller),
			userSvc:    nil,
			logger:     &zap.SugaredLogger{},
			expect: func(t *testing.T, s room.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[chat_room_service] invalid user service")
			},
		},
		{
			name:       "should return invalid logger",
			repository: mock_room.NewMockRepository(controller),
			userSvc:    mock_user.NewMockService(controller),
			logger:     nil,
			expect: func(t *testing.T, s room.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[chat_room_service] invalid logger")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := room.NewService(tc.repository, tc.userSvc, tc.logger)
			tc.expect(t, svc, err)
		})
	}
}
