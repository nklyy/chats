package room

import (
	"errors"
	"go.uber.org/zap"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type Service interface {
	CreateRoom(name string) (*Room, error)
}

type service struct {
	logger *zap.SugaredLogger
}

func NewService(logger *zap.SugaredLogger) (Service, error) {
	if logger == nil {
		return nil, errors.New("[chat_room_service] invalid logger")
	}

	return &service{logger: logger}, nil
}

func (s *service) CreateRoom(roomName string) (*Room, error) {
	room, err := NewRoom(roomName)
	if err != nil {
		s.logger.Errorf("failed to create new user %v", err)
		return nil, ErrFailedCreateRoom
	}

	return room, nil
}
