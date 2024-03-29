package chat

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"math/rand"
	"one-time-session-chat/internal/chat/room"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type Service interface {
	Chat(fingerprint string, ws *websocket.Conn) error
}

type service struct {
	redisClient *redis.Client
	clients     map[*room.Client]bool
	rooms       map[*room.Room]bool
	roomSvc     room.Service
	logger      *zap.SugaredLogger
}

func NewService(redisClient *redis.Client, roomSvc room.Service, logger *zap.SugaredLogger) (Service, error) {
	if redisClient == nil {
		return nil, errors.New("[chat_service] invalid redis chat client")
	}
	if roomSvc == nil {
		return nil, errors.New("[chat_service] invalid room service")
	}
	if logger == nil {
		return nil, errors.New("[chat_service] invalid logger")
	}
	return &service{
		redisClient: redisClient,
		clients:     make(map[*room.Client]bool),
		rooms:       make(map[*room.Room]bool),
		roomSvc:     roomSvc,
		logger:      logger,
	}, nil
}

func (s *service) Chat(fingerprint string, ws *websocket.Conn) error {
	if fingerprint == "" || fingerprint == "null" || fingerprint == "undefined" {
		msg, _ := s.encodeMessage(room.MessageResponse{
			Action:  "",
			Message: nil,
			From:    "",
			Error:   "invalid fingerprint",
		})

		err := ws.WriteMessage(1, msg)
		if err != nil {
			s.logger.Errorf("failed to write message %v", err)
		}

		return nil
	}

	newClient, _ := room.NewClient(fingerprint, ws)
	go newClient.WritePump()
	go newClient.ReadPump(s.messageHandler)
	s.cleanupOldConnections(fingerprint)
	s.findCompanion(newClient)

	return nil
}

func (s *service) findServerFreeUser(currentUserFingerprint string) []*room.Client {
	var foundFreeClients []*room.Client
	for client := range s.clients {
		if client.Fingerprint != currentUserFingerprint && client.Room == nil {
			foundFreeClients = append(foundFreeClients, client)
		}
	}

	return foundFreeClients
}

func (s *service) findCompanion(client *room.Client) {
	foundFreeClients := s.findServerFreeUser(client.Fingerprint)
	if foundFreeClients != nil {
		rand.Seed(time.Now().Unix())
		freeClient := foundFreeClients[rand.Intn(len(foundFreeClients))]

		// Create room
		newRoomId, err := uuid.NewUUID()
		if err != nil {
			s.logger.Errorf("failed create uuid %cv", err)
		}
		newRoom, _ := s.roomSvc.CreateRoom(newRoomId.String())

		msg, _ := s.encodeMessage(room.MessageResponse{
			Action:  "connected",
			Message: nil,
			From:    "",
			Error:   nil,
		})

		client.Room = newRoom
		newRoom.Clients[client] = true
		client.Send <- msg

		freeClient.Room = newRoom
		newRoom.Clients[freeClient] = true
		freeClient.Send <- msg

		go newRoom.RunRoom(s.redisClient, client, freeClient)
		s.rooms[newRoom] = true

	}
	s.clients[client] = true
}
