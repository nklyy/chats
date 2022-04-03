package old_chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"
	"noname-realtime-support-chat/pkg/jwt"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type Service interface {
	Chat(ws *websocket.Conn, token string, uID, rID string) error
	GetUser() (*Client, error)
}

type service struct {
	logger      *zap.SugaredLogger
	redisClient *redis.Client
	clients     map[*Client]bool
	rooms       map[*Room]bool
	jwtSvc      jwt.Service
}

func NewService(logger *zap.SugaredLogger, redisClient *redis.Client, clients map[*Client]bool, rooms map[*Room]bool, jwtSvc jwt.Service) (Service, error) {
	if logger == nil {
		return nil, errors.New("invalid logger")
	}
	if redisClient == nil {
		return nil, errors.New("invalid redis chat client")
	}
	if clients == nil {
		return nil, errors.New("invalid clients map")
	}
	if rooms == nil {
		return nil, errors.New("invalid rooms map")
	}
	if jwtSvc == nil {
		return nil, errors.New("invalid jwt service")
	}
	return &service{logger: logger, redisClient: redisClient, clients: clients, rooms: rooms, jwtSvc: jwtSvc}, nil
}

func (s *service) Chat(ws *websocket.Conn, token string, uID, rID string) error {
	if token == "" {
		fmt.Println(ws.RemoteAddr())
		var client *Client

		fmt.Println(uID, rID, "ASDASDASD")
		if rID != "" && uID != "" {
			client = &Client{
				Id:         uID,
				Room:       nil,
				Connection: ws,
				Free:       true,
				Support:    false,
				Send:       make(chan []byte, 256),
			}

			for room, _ := range s.rooms {
				if room.Name == rID {
					client.Room = room
					for c, _ := range room.clients {
						if !c.Support {
							fmt.Println(c.Id)
							delete(room.clients, c)
						}
					}
					room.clients[client] = true
					fmt.Println(len(room.clients))
					break
				}
			}

			s.clients[client] = true
			go s.readPump(ws)
			go s.writePump(ws, client)
			j, err := json.Marshal(map[string]string{"user_id": client.Id, "room_id": client.Room.Name})
			if err != nil {
				return err
			}

			client.Connection.WriteMessage(1, j)
		} else {
			userId, err := uuid.NewUUID()
			if err != nil {
				return err
			}

			roomId, _ := uuid.NewUUID()

			client = &Client{
				Id:         userId.String(),
				Room:       nil,
				Connection: ws,
				Free:       true,
				Support:    false,
				Send:       make(chan []byte, 256),
			}

			s.clients[client] = true
			go s.readPump(ws)
			go s.writePump(ws, client)
			//user := s.clients[client]

			j, err := json.Marshal(map[string]string{"user_id": client.Id, "room_id": roomId.String()})
			if err != nil {
				return err
			}

			client.Connection.WriteMessage(1, j)
			//go s.readPump(ws)
			//go s.writePump(ws, client)
			s.createRoom(client, roomId.String())
		}
	} else {
		//fmt.Println(token)
		_, err := s.jwtSvc.ParseToken(token, true)
		if err != nil {
			fmt.Println(err)
			ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			ws.Close()
			return nil
		}

		fmt.Println(ws.RemoteAddr())
		userId, err := uuid.NewUUID()
		if err != nil {
			return err
		}

		client := &Client{
			Id:         userId.String(),
			Room:       nil,
			Connection: ws,
			Free:       true,
			Support:    true,
			Send:       make(chan []byte, 256),
		}

		s.clients[client] = true
		//user := s.clients[client]

		j, err := json.Marshal(map[string]string{"user_id": client.Id})
		if err != nil {
			return err
		}
		client.Connection.WriteMessage(1, j)
		go s.readPump(ws)
		go s.writePump(ws, client)
	}

	return nil
}

func (s *service) createRoom(client *Client, roomId string) {
	room := NewRoom(roomId)
	room.clients[client] = true
	client.Room = room

	go client.Room.RunRoom(s.redisClient)
	s.rooms[room] = true
	client.Room.subscribeToRoomMessages(s.redisClient, roomId)
}

func (s *service) GetUser() (*Client, error) {
	for c, _ := range s.clients {
		if c.Room != nil && c.Free && !c.Support {
			c.Free = false
			user := c
			return user, nil
		}
	}
	return nil, errors.New("no users")
}

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

func (s *service) readPump(ws *websocket.Conn) {
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		fmt.Println(string(jsonMessage))
		s.messageHandler(jsonMessage)
		//client.handleNewMessage(jsonMessage)
	}
}

func (s *service) writePump(ws *websocket.Conn, client *Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		ws.Close()
	}()
	for {
		select {
		case message, ok := <-client.Send:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := ws.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = w.Write(message)
			if err != nil {
				return
			}

			//Attach queued chat messages to the current websocket message.
			n := len(client.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-client.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (s *service) messageHandler(jsonMessage []byte) {
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}

	fmt.Println(message)
	switch message.Action {
	//case "create-room":
	//	roomId, _ := uuid.NewUUID()
	//
	//	room := NewRoom(roomId.String())
	//
	//	for _, v := range s.clients {
	//		if v.Id == message.User {
	//			room.clients[v] = true
	//			v.Room = room
	//			j, err := json.Marshal(map[string]string{"room_id": v.Room.Name})
	//			if err != nil {
	//				log.Fatal(err)
	//			}
	//			v.Send <- j
	//			go v.Room.RunRoom(s.redisClient)
	//			v.Room.subscribeToRoomMessages(s.redisClient, roomId.String())
	//		}
	//	}
	case "join-room":
		var user *Client
		for c, _ := range s.clients {
			if c.Id == message.User {
				user = c
			}
		}
		fmt.Println("USER JOIN", user)

		for c, _ := range s.clients {
			if c.Room != nil && c.Room.Name == message.TargetRoom {
				c.Room.clients[user] = true
				fmt.Println("ROOM JOIN", c.Room)
			}
		}
	case "publish-room":
		//for c, _ := range s.clients {
		//	if c.Room != nil && c.Room.Name == message.TargetRoom {
		//		//j, err := json.Marshal(message)
		//		//if err != nil {
		//		//	log.Println(err)
		//		//}
		//		//v.Room.publishRoomMessage(s.redisClient, j)
		//		c.Room.broadcast <- &message
		//	}
		//}

		for r, _ := range s.rooms {
			if r.Name == message.TargetRoom {
				//j, err := json.Marshal(message)
				//if err != nil {
				//	log.Println(err)
				//}
				//v.Room.publishRoomMessage(s.redisClient, j)
				r.broadcast <- &message
			}
		}
	case "disconnect":
		for c, _ := range s.clients {
			if c.Id == message.User {

				//for c, _ := range v.Room.clients {
				//	c.Connection.Close()
				//}

				c.Connection.Close()
				delete(s.clients, c)
				fmt.Println(s.clients, len(s.clients))
			}
		}
	}
}
