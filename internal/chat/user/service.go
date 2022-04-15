package user

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type Service interface {
	GetUserById(ctx context.Context, id string) (*DTO, error)
	GetUserByFingerprint(ctx context.Context, fingerprint string) (*DTO, error)
	GetFreeUser(ctx context.Context, userId string) (*DTO, error)
	CreateUser(ctx context.Context, fingerprint string) (*DTO, error)
	UpdateUser(ctx context.Context, userDTO *DTO) error
}

type service struct {
	repository Repository
	logger     *zap.SugaredLogger
	salt       string
}

func NewService(repository Repository, logger *zap.SugaredLogger, salt string) (Service, error) {
	if repository == nil {
		return nil, errors.New("invalid repository")
	}
	if logger == nil {
		return nil, errors.New("invalid logger")
	}
	if salt == "" {
		return nil, errors.New("invalid salt")
	}

	return &service{repository: repository, logger: logger, salt: salt}, nil
}

func (s *service) GetUserById(ctx context.Context, id string) (*DTO, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	user, err := s.repository.GetUser(ctx, bson.M{"_id": objId})
	if err != nil {
		s.logger.Errorf("failed to get user: %v", err)
		return nil, err
	}

	return MapToDTO(user), nil
}

func (s *service) GetUserByFingerprint(ctx context.Context, fingerprint string) (*DTO, error) {
	user, err := s.repository.GetUser(ctx, bson.M{"fingerprint": fingerprint})
	if err != nil {
		s.logger.Errorf("failed to get user: %v", err)
		return nil, err
	}

	return MapToDTO(user), nil
}

func (s *service) GetFreeUser(ctx context.Context, userId string) (*DTO, error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		s.logger.Errorf("failed conver user id: %v", err)
		return nil, err
	}

	users, err := s.repository.GetUsers(ctx, bson.M{"_id": bson.M{"$ne": id}, "room_name": bson.M{"$eq": nil}})
	if err != nil {
		s.logger.Errorf("failed to get user: %v", err)
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrNoUsersYet
	}

	rand.Seed(time.Now().Unix())

	return MapToDTO(users[rand.Intn(len(users))]), nil

	//user, err := s.repository.GetFreeUser(ctx)
	//if err != nil {
	//	s.logger.Errorf("failed to get user: %v", err)
	//	return nil, err
	//}
	//
	//userCtxValue := ctx.Value(contextKey("user"))
	//if userCtxValue == nil {
	//	s.logger.Error("Not authenticated")
	//	return nil, errors.New("not authenticated")
	//}
	//
	//ctxUserDto := userCtxValue.(DTO)
	//
	//ctxUserEntity, err := MapToEntity(&ctxUserDto)
	//if err != nil {
	//	s.logger.Error(err)
	//	return nil, err
	//}
	//
	//ctxUserEntity.SetRoom(user.RoomName)
	//
	//// update support
	//err = s.UpdateUser(ctx, MapToDTO(ctxUserEntity))
	//if err != nil {
	//	s.logger.Error(err)
	//	return nil, ErrFailedUpdateUser
	//}
	//
	//// update user
	//user.SetFreeStatus(false)
	//err = s.UpdateUser(ctx, MapToDTO(user))
	//if err != nil {
	//	s.logger.Error(err)
	//	return nil, ErrFailedUpdateUser
	//}
	//
	//user.RemovePassword()
	//
	//return MapToDTO(user), nil
}

func (s *service) CreateUser(ctx context.Context, fingerprint string) (*DTO, error) {
	user, err := NewUser(fingerprint, s.salt)
	if err != nil {
		s.logger.Errorf("failed to create new user %v", err)
		return nil, ErrFailedCreateUser
	}

	_, err = s.repository.CreateUser(ctx, user)
	if err != nil {
		s.logger.Errorf("failed to save user %v", err)
		return nil, err
	}

	return MapToDTO(user), nil
}

func (s *service) UpdateUser(ctx context.Context, userDTO *DTO) error {
	// map dto to user entity
	updateUser, err := MapToEntity(userDTO)
	if err != nil {
		return err
	}

	// update user in storage by email
	if err = s.repository.UpdateUser(ctx, updateUser); err != nil {
		s.logger.Errorf("failed to save user in db: %v", err)
		return err
	}
	return nil
}
