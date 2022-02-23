package support

import (
	"context"
	"errors"
	"go.uber.org/zap"
)

type Service interface {
	GetSupportById(ctx context.Context, id string) (*DTO, error)
	CreateSupport(ctx context.Context, dto *CreateSupportDTO) (string, error)
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

	support.RemovePassword()

	return MapToDTO(support), nil
}

func (s *service) CreateSupport(ctx context.Context, dto *CreateSupportDTO) (string, error) {
	support, err := NewSupport(dto.Email, dto.Name, dto.Password, s.salt)
	if err != nil {
		s.logger.Errorf("failed to create new support %v", err)
		return "", err
	}

	id, err := s.repository.CreateSupport(ctx, support)
	if err != nil {
		s.logger.Errorf("failed to save support %v", err)
		return "", err
	}

	return id, nil
}
