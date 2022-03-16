package support

import (
	"context"
	"errors"
	"go.uber.org/zap"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type Service interface {
	GetSupportById(ctx context.Context, id string) (*DTO, error)
	GetSupportByEmail(ctx context.Context, email string) (*DTO, error)
	CreateSupport(ctx context.Context, email, name, password string) (*DTO, error)
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

func (s *service) GetSupportById(ctx context.Context, id string) (*DTO, error) {
	support, err := s.repository.GetSupportById(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to get support: %v", err)
		return nil, err
	}

	return MapToDTO(support), nil
}

func (s *service) GetSupportByEmail(ctx context.Context, id string) (*DTO, error) {
	support, err := s.repository.GetSupportByEmail(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to get support: %v", err)
		return nil, err
	}

	return MapToDTO(support), nil
}

func (s *service) CreateSupport(ctx context.Context, email, name, password string) (*DTO, error) {
	support, err := NewSupport(email, name, password, &s.salt)
	if err != nil {
		s.logger.Errorf("failed to create new support %v", err)
		return nil, ErrFailedCreateSupport
	}

	_, err = s.repository.CreateSupport(ctx, support)
	if err != nil {
		s.logger.Errorf("failed to save support %v", err)
		return nil, err
	}

	return MapToDTO(support), nil
}
