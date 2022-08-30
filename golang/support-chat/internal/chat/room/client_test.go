package room_test

import (
	"support-chat/internal/chat/room"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	id := "id"

	tests := []struct {
		name   string
		id     string
		conn   *websocket.Conn
		expect func(*testing.T, *room.Client, error)
	}{
		{
			name: "should return service",
			id:   id,
			conn: &websocket.Conn{},
			expect: func(t *testing.T, s *room.Client, err error) {
				assert.NotNil(t, s)
				assert.Nil(t, err)
			},
		},
		{
			name: "should return invalid id",
			id:   "",
			conn: &websocket.Conn{},
			expect: func(t *testing.T, s *room.Client, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[chat_room_client] invalid id")
			},
		},
		{
			name: "should return invalid connection",
			id:   id,
			conn: nil,
			expect: func(t *testing.T, s *room.Client, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[chat_room_client] invalid websocket connection")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := room.NewClient(tc.id, tc.conn)
			tc.expect(t, svc, err)
		})
	}
}
