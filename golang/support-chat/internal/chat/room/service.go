package room

import (
	"context"
	"errors"
	"support-chat/internal/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type Service interface {
	GetRoomByName(ctx context.Context, name string) (*DTO, error)
	GetRoomWithFormatMessages(ctx context.Context, name, userId string) ([]*FormatMessages, error)
	CreateRoom(ctx context.Context, name string, user *user.DTO) (*Room, error)
	UpdateRoom(ctx context.Context, dto *DTO) error
	DeleteRoom(ctx context.Context, name string) error
}

type service struct {
	repository Repository
	userSvc    user.Service
	logger     *zap.SugaredLogger
}

func NewService(repository Repository, userSvc user.Service, logger *zap.SugaredLogger) (Service, error) {
	if repository == nil {
		return nil, errors.New("[chat_room_service] invalid repository")
	}
	if userSvc == nil {
		return nil, errors.New("[chat_room_service] invalid user service")
	}
	if logger == nil {
		return nil, errors.New("[chat_room_service] invalid logger")
	}

	return &service{repository: repository, userSvc: userSvc, logger: logger}, nil
}

func (s *service) GetRoomByName(ctx context.Context, name string) (*DTO, error) {
	room, err := s.repository.GetRoom(ctx, bson.M{"name": name})
	if err != nil {
		s.logger.Errorf("failed to get room: %v", err)
		return nil, err
	}

	return MapToDTO(room), nil
}

func (s *service) GetRoomWithFormatMessages(ctx context.Context, name, userId string) ([]*FormatMessages, error) {
	room, err := s.repository.GetRoom(ctx, bson.M{"name": name})
	if err != nil {
		s.logger.Errorf("failed to get room: %v", err)
		return nil, err
	}

	var msg []*FormatMessages
	if room.Messages != nil {
		for _, message := range *room.Messages {
			if message.Id == userId {
				msg = append(msg, &FormatMessages{
					To:      message.Id,
					Message: message.Message,
					Time:    message.Time,
				})
			} else {
				msg = append(msg, &FormatMessages{
					From:    message.Id,
					Message: message.Message,
					Time:    message.Time,
				})
			}
		}
	}

	return msg, nil
}

func (s *service) CreateRoom(ctx context.Context, roomName string, u *user.DTO) (*Room, error) {
	room, err := NewRoom(roomName)
	if err != nil {
		s.logger.Errorf("failed to create new user %v", err)
		return nil, ErrFailedCreateRoom
	}

	m := &Model{
		ID:       room.ID,
		Name:     room.Name,
		Messages: nil,
	}

	_, err = s.repository.CreateRoom(ctx, m)
	if err != nil {
		s.logger.Errorf("failed to create room %v", err)
		return nil, err
	}

	userEntity, _ := user.MapToEntity(u)
	userEntity.SetRoom(&roomName)
	userEntity.SetFreeStatus(true)

	userDto := user.MapToDTO(userEntity)

	err = s.userSvc.UpdateUser(ctx, userDto)
	if err != nil {
		s.logger.Errorf("failed to update user %v", err)
		return nil, err
	}

	return room, nil
}

func (s *service) UpdateRoom(ctx context.Context, dto *DTO) error {
	// map dto to user entity
	model, err := MapToEntity(dto)
	if err != nil {
		return err
	}

	// update user in storage by email
	if err = s.repository.UpdateRoom(ctx, model); err != nil {
		s.logger.Errorf("failed to save room in db: %v", err)
		return err
	}
	return nil
}

func (s *service) DeleteRoom(ctx context.Context, name string) error {
	err := s.repository.DeleteRoom(ctx, name)
	if err != nil {
		s.logger.Errorf("failed to delete room in db: %v", err)
		return err
	}
	return nil
}
