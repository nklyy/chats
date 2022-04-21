package room_test

import (
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"noname-one-time-session-chat/internal/chat/room"
	"testing"
)

func TestNewClient(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	fingerprint := "fingerprint"

	tests := []struct {
		name        string
		fingerprint string
		conn        *websocket.Conn
		expect      func(*testing.T, *room.Client, error)
	}{
		{
			name:        "should return service",
			fingerprint: fingerprint,
			conn:        &websocket.Conn{},
			expect: func(t *testing.T, s *room.Client, err error) {
				assert.NotNil(t, s)
				assert.Nil(t, err)
			},
		},
		{
			name:        "should return invalid fingerprint",
			fingerprint: "",
			conn:        &websocket.Conn{},
			expect: func(t *testing.T, s *room.Client, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[chat_room_client] invalid fingerprint")
			},
		},
		{
			name:        "should return invalid connection",
			fingerprint: fingerprint,
			conn:        nil,
			expect: func(t *testing.T, s *room.Client, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[chat_room_client] invalid connection")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := room.NewClient(tc.fingerprint, tc.conn)
			tc.expect(t, svc, err)
		})
	}
}
