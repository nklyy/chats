package room

import (
	"github.com/gorilla/websocket"
	"log"

	"noname-realtime-support-chat/pkg/errors"
	"time"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less than pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

type Client struct {
	Id         string          `json:"id"`
	Room       *Room           `json:"room"`
	Connection *websocket.Conn `json:"connection"`
	Send       chan []byte     `json:"send"`
}

func NewClient(id string, conn *websocket.Conn) (*Client, error) {
	if id == "" {
		return nil, errors.WithMessage(ErrInvalidId, "should be not empty")
	}
	if conn == nil {
		return nil, errors.WithMessage(ErrInvalidConnection, "should be not empty")
	}

	return &Client{
		Id:         id,
		Room:       nil,
		Connection: conn,
		Send:       make(chan []byte, 256),
	}, nil
}

type HandlerFunc func([]byte)

func (c *Client) ReadPump(msgHandleFunc HandlerFunc) {
	c.Connection.SetReadLimit(maxMessageSize)
	c.Connection.SetReadDeadline(time.Now().Add(pongWait))
	c.Connection.SetPongHandler(func(string) error { c.Connection.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		//fmt.Println(string(jsonMessage))

		msgHandleFunc(jsonMessage)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Connection.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = w.Write(message)
			if err != nil {
				return
			}

			//Attach queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
