package chat

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"
	"noname-realtime-support-chat/internal/chat/room"
	"noname-realtime-support-chat/internal/user"
	"noname-realtime-support-chat/pkg/jwt"
	"time"
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
	userCtxValue := ctx.Value(contextKey("user"))
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
				s.createRoomIfDoesntExist(ctx, client, u)
			} else {
				client.Room = r
				r.Clients[client] = true
			}
		} else {
			s.createRoomIfDoesntExist(ctx, client, u)
		}
	}

	if u.Support && u.RoomName != nil {
		r := s.findRoom(ctx, *u.RoomName)
		if r == nil {
			var emptyRoom string

			uEntity, _ := user.MapToEntity(u)
			uEntity.SetRoom(&emptyRoom)

			err := s.userSvc.UpdateUser(ctx, user.MapToDTO(uEntity))
			if err != nil {
				msg, err := s.encodeMessage(room.MessageResponse{
					Action:  "",
					Message: nil,
					From:    "",
					Error:   "failed update user",
				})
				if err != nil {
					s.logger.Errorf("failed to encode message %v", err)
				}

				client.Send <- msg
				return
			}
		} else {
			r.Clients[client] = true
		}
	}

	s.clients[client] = true
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

func (s *service) createRoomIfDoesntExist(ctx context.Context, client *room.Client, u *user.DTO) {
	newRoomId, err := uuid.NewUUID()
	if err != nil {
		s.logger.Errorf("failed create uuid %cv", err)
	}

	newRoom, err := s.roomSvc.CreateRoom(ctx, newRoomId.String(), u)
	if err != nil {
		s.logger.Errorf("failed create room %v", err)
		msg, err := s.encodeMessage(room.MessageResponse{
			Action:  "",
			Message: nil,
			From:    "",
			Error:   "failed create room",
		})
		if err != nil {
			s.logger.Errorf("failed to create room %v", err)
		}

		client.Send <- msg
		return
	}

	client.Room = newRoom
	newRoom.Clients[client] = true
	go newRoom.RunRoom(s.redisClient)
	s.rooms[newRoom] = true
}

func (s *service) encodeMessage(msg room.MessageResponse) ([]byte, error) {
	encMsg, err := json.Marshal(msg)
	if err != nil {
		s.logger.Errorf("failed to encode message %v", err)
		return nil, err
	}
	return encMsg, err
}

func (s *service) messageHandler(jsonMessage []byte) {
	var message room.Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		s.logger.Errorf("Error on unmarshal JSON message %s", err)
		return
	}

	switch message.Action {
	case "publish-room":
		uPayload, err := s.jwtSvc.ParseToken(message.Token, true)
		if err != nil {
			s.logger.Errorf("failed to parse token %v", err)
		}

		dbUser, err := s.userSvc.GetUserById(context.Background(), uPayload.Id, false)
		if err != nil {
			s.logger.Errorf("failed to get user %v", err)
		}

		dbRoom, err := s.roomSvc.GetRoomByName(context.Background(), *dbUser.RoomName)
		if err != nil {
			s.logger.Errorf("failed to get room %v", err)
		}

		for r := range s.rooms {
			if r.Name == *dbUser.RoomName {
				var msg []*room.RoomMessage

				if dbRoom.Messages == nil {
					msg = append(msg, &room.RoomMessage{
						Id:      dbUser.ID,
						Time:    time.Now(),
						Message: message.Message,
					})
				} else {
					msg = append(*dbRoom.Messages, &room.RoomMessage{
						Id:      dbUser.ID,
						Time:    time.Now(),
						Message: message.Message,
					})
				}
				dbRoom.Messages = &msg

				////////
				err := s.roomSvc.UpdateRoom(context.Background(), dbRoom)
				if err != nil {

					r.Broadcast <- &room.BroadcastMessage{
						Action: message.Action,
						Message: room.MessageResponse{
							Action:  "",
							Message: nil,
							From:    "",
							Error:   "failed update room",
						},
						RoomName: *dbUser.RoomName,
					}
				}

				r.Broadcast <- &room.BroadcastMessage{
					Action: message.Action,
					Message: room.MessageResponse{
						Action:  message.Action,
						Message: &message.Message,
						From:    dbUser.ID,
						Error:   nil,
					},
					RoomName: *dbUser.RoomName,
				}
			}
		}
	case "disconnect":
		uPayload, err := s.jwtSvc.ParseToken(message.Token, true)
		if err != nil {
			s.logger.Errorf("failed to parse token %v", err)
		}

		dbUser, err := s.userSvc.GetUserById(context.Background(), uPayload.Id, false)
		if err != nil {
			s.logger.Errorf("failed to get user %v", err)
		}

		for r := range s.rooms {
			if r.Name == *dbUser.RoomName {
				for client := range r.Clients {
					rUser, err := s.userSvc.GetUserById(context.Background(), client.Id, true)
					if err != nil {
						s.logger.Errorf("failed to get user %v", err)
					}

					userEntity, _ := user.MapToEntity(rUser)
					userEntity.SetRoom(nil)
					userEntity.SetFreeStatus(true)
					err = s.userSvc.UpdateUser(context.Background(), user.MapToDTO(userEntity))
					if err != nil {
						msg, _ := s.encodeMessage(room.MessageResponse{
							Action:  "",
							Message: nil,
							From:    "",
							Error:   "failed update user",
						})
						client.Send <- msg
						return
					}

					close(client.Send)
					err = client.Connection.Close()
					if err != nil {
						s.logger.Errorf("failed close connection %v", err)
					}

					for serverClient := range s.clients {
						if serverClient.Id == client.Id {
							delete(s.clients, client)
						}
					}
				}

				err := s.roomSvc.DeleteRoom(context.Background(), r.Name)
				if err != nil {
					s.logger.Errorf("failed ddelete room %v", err)
				}
				delete(s.rooms, r)
				break
			}
		}
	}
}
