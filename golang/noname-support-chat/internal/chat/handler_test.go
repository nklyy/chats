package chat_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"noname-support-chat/internal/chat"
	mock_chat "noname-support-chat/internal/chat/mocks"
	"testing"
)

func TestNewHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tests := []struct {
		name    string
		chatSvc chat.Service
		expect  func(*testing.T, *chat.Handler, error)
	}{
		{
			name:    "should return service",
			chatSvc: mock_chat.NewMockService(controller),
			expect: func(t *testing.T, s *chat.Handler, err error) {
				assert.NotNil(t, s)
				assert.Nil(t, err)
			},
		},
		{
			name:    "should return invalid chat service",
			chatSvc: nil,
			expect: func(t *testing.T, s *chat.Handler, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[chat_handler] invalid chat service")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := chat.NewHandler(tc.chatSvc)
			tc.expect(t, svc, err)
		})
	}
}
