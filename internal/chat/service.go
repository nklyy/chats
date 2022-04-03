package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"
	"noname-realtime-support-chat/internal/chat/room"
	"noname-realtime-support-chat/internal/user"
	"noname-realtime-support-chat/pkg/jwt"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type Service interface {
	Chat(ctx context.Context, ws *websocket.Conn) error
}

type service struct {
	redisClient *redis.Client
	clients     map[*room.Client]bool
	rooms       map[*room.Room]bool
	roomSvc     room.Service
	jwtSvc      jwt.Service
	userSvc     user.Service
	logger      *zap.SugaredLogger
}

func NewService(redisClient *redis.Client, roomSvc room.Service, jwtSvc jwt.Service, userSvc user.Service, logger *zap.SugaredLogger) (Service, error) {
	if redisClient == nil {
		return nil, errors.New("invalid redis chat client")
	}
	if jwtSvc == nil {
		return nil, errors.New("invalid jwt service")
	}
	if roomSvc == nil {
		return nil, errors.New("invalid room service")
	}
	if roomSvc == nil {
		return nil, errors.New("invalid user service")
	}
	if logger == nil {
		return nil, errors.New("invalid logger")
	}
	return &service{
		logger:      logger,
		clients:     make(map[*room.Client]bool),
		rooms:       make(map[*room.Room]bool),
		roomSvc:     roomSvc,
		jwtSvc:      jwtSvc,
		userSvc:     userSvc,
		redisClient: redisClient,
	}, nil
}

func (s *service) Chat(ctx context.Context, ws *websocket.Conn) error {
	userCtxValue := ctx.Value("user")
	if userCtxValue == nil {
		log.Println("Not authenticated")
		return nil
	}

	u := userCtxValue.(user.DTO)
	c, err := room.NewClient(u.ID, ws)
	if err != nil {
		return err
	}

	go c.WritePump()
	go c.ReadPump(s.messageHandler)
	s.registerClientAndCreateRoom(ctx, c, &u)

	return nil
}

func (s *service) registerClientAndCreateRoom(ctx context.Context, client *room.Client, u *user.DTO) {
	if !u.Support {
		if u.RoomName != nil {
			r := s.findRoom(ctx, *u.RoomName)

			if r == nil {
				newRoomId, _ := uuid.NewUUID()
				r, _ := s.roomSvc.CreateRoom(ctx, newRoomId.String(), u)
				client.Room = r
				r.Clients[client] = true
				go r.RunRoom(s.redisClient)
			} else {
				client.Room = r
				r.Clients[client] = true
			}
		} else {
			newRoomId, _ := uuid.NewUUID()
			r, _ := s.roomSvc.CreateRoom(ctx, newRoomId.String(), u)
			client.Room = r
			r.Clients[client] = true
			go r.RunRoom(s.redisClient)
		}
	}

	if u.Support && u.RoomName != nil {
		r := s.findRoom(ctx, *u.RoomName)
		if r == nil {
			client.Send <- []byte("Room doesn't exist!")
			uEntity, _ := user.MapToEntity(u)
			var emptyRoom string
			uEntity.SetRoom(&emptyRoom)
			s.userSvc.UpdateUser(ctx, user.MapToDTO(uEntity))
		} else {
			r.Clients[client] = true
		}
	}

	s.clients[client] = true
}

func (s *service) findRoom(ctx context.Context, roomName string) *room.Room {
	var foundRoom *room.Room
	for r, _ := range s.rooms {
		if r.Name == roomName {
			foundRoom = r
			break
		}
	}

	if foundRoom == nil {
		foundRoom = s.runRoomFromRepo(ctx, roomName)
	}

	return foundRoom
}

func (s *service) runRoomFromRepo(ctx context.Context, roomName string) *room.Room {
	var r *room.Room
	dbRoom, _ := s.roomSvc.GetRoomByName(ctx, roomName)
	if dbRoom != nil {
		r, _ = room.NewRoom(dbRoom.Name)
		go r.RunRoom(s.redisClient)
		s.rooms[r] = true
	}

	return r
}

func (s *service) messageHandler(jsonMessage []byte) {
	var message room.Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}

	fmt.Println(message)
	switch message.Action {
	case "join-room":
		fmt.Println("ASDASDSADASDASDASDSADSADASDSA")
		var u *room.Client
		for cl, _ := range s.clients {
			fmt.Println(cl)
			if cl.Id == message.User {
				u = cl
				break
			}
		}

		r := s.findRoom(context.Background(), message.TargetRoom)
		if r == nil {
			u.Send <- []byte("Room doesn't exist")
			return
		}
		r.Clients[u] = true
		fmt.Println("ROOM JOIN", r.Name)

		//fmt.Println("USER JOIN", u)
		//
		//for cl, _ := range s.clients {
		//	if cl.Room != nil && cl.Room.Name == message.TargetRoom {
		//		cl.Room.Clients[u] = true
		//		fmt.Println("ROOM JOIN", cl.Room)
		//	}
		//}
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
				r.Broadcast <- &message
			}
		}
	case "disconnect":

		for cl, _ := range s.clients {
			if cl.Id == message.User {

				//for c, _ := range v.Room.clients {
				//	c.Connection.Close()
				//}

				cl.Connection.Close()
				delete(s.clients, cl)
				fmt.Println(s.clients, len(s.clients))
			}
		}
	}
}
