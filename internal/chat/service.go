package chat

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"golang.org/x/crypto/scrypt"
	"net"
	"noname-realtime-support-chat/internal/chat/room"
	"noname-realtime-support-chat/internal/chat/user"
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
	userSvc     user.Service
	salt        string
	logger      *zap.SugaredLogger
}

func NewService(redisClient *redis.Client, roomSvc room.Service, userSvc user.Service, salt string, logger *zap.SugaredLogger) (Service, error) {
	if redisClient == nil {
		return nil, errors.New("invalid redis chat client")
	}
	if roomSvc == nil {
		return nil, errors.New("invalid room service")
	}
	if roomSvc == nil {
		return nil, errors.New("invalid user service")
	}
	if salt == "" {
		return nil, errors.New("invalid salt")
	}
	if logger == nil {
		return nil, errors.New("invalid logger")
	}
	return &service{
		redisClient: redisClient,
		clients:     make(map[*room.Client]bool),
		rooms:       make(map[*room.Room]bool),
		roomSvc:     roomSvc,
		userSvc:     userSvc,
		salt:        salt,
		logger:      logger,
	}, nil
}

func (s *service) Chat(ctx context.Context, ws *websocket.Conn) error {
	// IPv6
	host, port, err := net.SplitHostPort(ws.RemoteAddr().String())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(host, port)

	var usr *user.DTO

	hashedAddr, _ := s.createHash(host)
	usr, _ = s.userSvc.GetUserByIp(ctx, hashedAddr)
	if usr == nil {
		usr, _ = s.userSvc.CreateUser(ctx, host)
	}
	//fmt.Println(usr)

	//fmt.Println("SADASDASDASDASD", hashedAddr)
	newClient, _ := room.NewClient(usr.ID, ws)
	go newClient.WritePump()
	go newClient.ReadPump(s.messageHandler, hashedAddr)
	s.findCompanion(ctx, newClient, usr)

	return nil
}

func (s *service) createHash(data string) (string, error) {
	hash, err := scrypt.Key([]byte(data), []byte(s.salt), 16384, 8, 1, 32)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(hash), nil
}

func (s *service) findCompanion(ctx context.Context, client *room.Client, userDto *user.DTO) {
	if userDto.RoomName != nil {
		foundRoom := s.findRoom(ctx, *userDto.RoomName)

		if foundRoom == nil {
			newRoomId, err := uuid.NewUUID()
			if err != nil {
				s.logger.Errorf("failed create uuid %cv", err)
			}
			newRoom, _ := s.roomSvc.CreateRoom(ctx, newRoomId.String())

			for {
				freeUser, _ := s.userSvc.GetFreeUser(ctx, userDto.ID)
				if freeUser != nil {
					foundClients := s.findServerClients(freeUser.ID)
					if foundClients != nil {
						for _, foundClient := range foundClients {
							foundClient.Room = newRoom
							newRoom.Clients[foundClient] = true
						}

						// update free user
						freeUserEntity, _ := user.MapToEntity(freeUser)
						freeUserEntity.SetRoom(&newRoom.Name)
						s.userSvc.UpdateUser(ctx, user.MapToDTO(freeUserEntity))
						break
					}
				}
			}

			client.Room = newRoom
			newRoom.Clients[client] = true

			msg, _ := s.encodeMessage(room.MessageResponse{
				Action:  "connected",
				Message: nil,
				From:    "",
				Error:   nil,
			})
			client.Send <- msg

			// update user
			userEntity, _ := user.MapToEntity(userDto)
			userEntity.SetRoom(&newRoom.Name)
			s.userSvc.UpdateUser(ctx, user.MapToDTO(userEntity))

			go newRoom.RunRoom(s.redisClient)
			s.rooms[newRoom] = true
		} else {
			client.Room = foundRoom
			foundRoom.Clients[client] = true
		}
	}

	s.clients[client] = true
}

func (s *service) findServerClients(freeUSerId string) []*room.Client {
	var foundClients []*room.Client

	for serverClient := range s.clients {
		if serverClient.Id == freeUSerId {
			foundClients = append(foundClients, serverClient)
		}
	}

	return foundClients
}

func (s *service) findRoom(ctx context.Context, roomName string) *room.Room {
	var foundRoom *room.Room
	for r := range s.rooms {
		if r.Name == roomName {
			foundRoom = r
			break
		}
	}

	if foundRoom == nil {
		foundRoom = s.runRoomFromRepository(ctx, roomName)
	}

	return foundRoom
}

func (s *service) runRoomFromRepository(ctx context.Context, roomName string) *room.Room {
	var r *room.Room
	dbRoom, _ := s.roomSvc.GetRoomByName(ctx, roomName)
	if dbRoom != nil {
		r, _ = room.NewRoom(dbRoom.Name)
		go r.RunRoom(s.redisClient)
		s.rooms[r] = true
	}

	return r
}

//func (s *service) createRoomIfDoesntExist(ctx context.Context, client *room.Client, u *user.DTO) {
//	newRoomId, err := uuid.NewUUID()
//	if err != nil {
//		s.logger.Errorf("failed create uuid %cv", err)
//	}
//
//	newRoom, err := s.roomSvc.CreateRoom(ctx, newRoomId.String(), u)
//	if err != nil {
//		s.logger.Errorf("failed create room %v", err)
//		msg, err := s.encodeMessage(room.MessageResponse{
//			Action:  "",
//			Message: nil,
//			From:    "",
//			Error:   "failed create room",
//		})
//		if err != nil {
//			s.logger.Errorf("failed to create room %v", err)
//		}
//
//		client.Send <- msg
//		return
//	}
//
//	client.Room = newRoom
//	newRoom.Clients[client] = true
//	go newRoom.RunRoom(s.redisClient)
//	s.rooms[newRoom] = true
//}

func (s *service) encodeMessage(msg room.MessageResponse) ([]byte, error) {
	encMsg, err := json.Marshal(msg)
	if err != nil {
		s.logger.Errorf("failed to encode message %v", err)
		return nil, err
	}
	return encMsg, err
}
