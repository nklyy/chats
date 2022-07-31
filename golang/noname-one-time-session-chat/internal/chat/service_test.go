package chat_test

import (
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"noname-one-time-session-chat/internal/chat"
	"noname-one-time-session-chat/internal/chat/room"
	mock_room "noname-one-time-session-chat/internal/chat/room/mocks"
	"testing"
)

func TestNewService(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tests := []struct {
		name        string
		redisClient *redis.Client
		roomSvc     room.Service
		logger      *zap.SugaredLogger
		expect      func(*testing.T, chat.Service, error)
	}{
		{
			name:        "should return service",
			redisClient: &redis.Client{},
			roomSvc:     mock_room.NewMockService(controller),
			logger:      &zap.SugaredLogger{},
			expect: func(t *testing.T, s chat.Service, err error) {
				assert.NotNil(t, s)
				assert.Nil(t, err)
			},
		},
		{
			name:        "should return invalid redis chat client",
			redisClient: nil,
			roomSvc:     mock_room.NewMockService(controller),
			logger:      &zap.SugaredLogger{},
			expect: func(t *testing.T, s chat.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[chat_service] invalid redis chat client")
			},
		},
		{
			name:        "should return invalid room service",
			redisClient: &redis.Client{},
			roomSvc:     nil,
			logger:      &zap.SugaredLogger{},
			expect: func(t *testing.T, s chat.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[chat_service] invalid room service")
			},
		},
		{
			name:        "should return invalid logger",
			redisClient: &redis.Client{},
			roomSvc:     mock_room.NewMockService(controller),
			logger:      nil,
			expect: func(t *testing.T, s chat.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[chat_service] invalid logger")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := chat.NewService(tc.redisClient, tc.roomSvc, tc.logger)
			tc.expect(t, svc, err)
		})
	}
}
