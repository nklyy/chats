package room_test

import (
	"support-chat/internal/chat/room"
	mock_room "support-chat/internal/chat/room/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tests := []struct {
		name    string
		roomSvc room.Service
		expect  func(*testing.T, *room.Handler, error)
	}{
		{
			name:    "should return service",
			roomSvc: mock_room.NewMockService(controller),
			expect: func(t *testing.T, s *room.Handler, err error) {
				assert.NotNil(t, s)
				assert.Nil(t, err)
			},
		},
		{
			name:    "should return invalid room service",
			roomSvc: nil,
			expect: func(t *testing.T, s *room.Handler, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[chat_room_handler] invalid room service")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := room.NewHandler(tc.roomSvc)
			tc.expect(t, svc, err)
		})
	}
}
