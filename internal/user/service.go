package user

import (
	"context"
	"errors"
	"go.uber.org/zap"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type Service interface {
	GetUserById(ctx context.Context, id string, withPassword bool) (*DTO, error)
	GetUserByEmail(ctx context.Context, email string, withPassword bool) (*DTO, error)
	CreateUser(ctx context.Context, email, name, password string) (*DTO, error)
}

type service struct {
	repository Repository
	logger     *zap.SugaredLogger
	salt       int
}

func NewService(repository Repository, logger *zap.SugaredLogger, salt *int) (Service, error) {
	if repository == nil {
		return nil, errors.New("invalid repository")
	}
	if logger == nil {
		return nil, errors.New("invalid logger")
	}
	if salt == nil {
		return nil, errors.New("invalid salt")
	}

	return &service{repository: repository, logger: logger, salt: *salt}, nil
}

func (s *service) GetUserById(ctx context.Context, id string, withPassword bool) (*DTO, error) {
	user, err := s.repository.GetUserById(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to get user: %v", err)
		return nil, err
	}

	if !withPassword {
		user.RemovePassword()
	}

	return MapToDTO(user), nil
}

func (s *service) GetUserByEmail(ctx context.Context, id string, withPassword bool) (*DTO, error) {
	user, err := s.repository.GetUserByEmail(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to get user: %v", err)
		return nil, err
	}

	if !withPassword {
		user.RemovePassword()
	}

	return MapToDTO(user), nil
}

func (s *service) CreateUser(ctx context.Context, email, name, password string) (*DTO, error) {
	user, err := NewUser(email, name, password, &s.salt)
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
