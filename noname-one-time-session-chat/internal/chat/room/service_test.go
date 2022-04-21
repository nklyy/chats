package room_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"noname-one-time-session-chat/internal/chat/room"
	"testing"
)

func TestNewService(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tests := []struct {
		name   string
		logger *zap.SugaredLogger
		expect func(*testing.T, room.Service, error)
	}{
		{
			name:   "should return service",
			logger: &zap.SugaredLogger{},
			expect: func(t *testing.T, s room.Service, err error) {
				assert.NotNil(t, s)
				assert.Nil(t, err)
			},
		},
		{
			name:   "should return invalid logger",
			logger: nil,
			expect: func(t *testing.T, s room.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[chat_room_service] invalid logger")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := room.NewService(tc.logger)
			tc.expect(t, svc, err)
		})
	}
}
