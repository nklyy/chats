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
	err := c.Connection.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Printf("failed to set read deadline %v", err)
	}
	c.Connection.SetPongHandler(func(string) error {
		err := c.Connection.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Printf("failed to set read deadline %v", err)
			return err
		}
		return nil
	})

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		msgHandleFunc(jsonMessage)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := c.Connection.Close()
		if err != nil {
			log.Printf("failed to close connection %v", err)
		}
	}()
	for {
		select {
		case message, ok := <-c.Send:
			err := c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Printf("failed to set write deadline %v", err)
			}
			if !ok {
				// The WsServer closed the channel.
				err := c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Printf("failed to write message %v", err)
				}
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
				_, err := w.Write([]byte{'\n'})
				if err != nil {
					log.Printf("failed to write message %v", err)
				}
				_, err = w.Write(<-c.Send)
				if err != nil {
					log.Printf("failed to write message %v", err)
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			err := c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Printf("failed to set write deadline %v", err)
			}
			if err := c.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
