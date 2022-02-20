package support

import (
	"context"
	"errors"
	"go.uber.org/zap"
)

type Service interface {
	GetSupportById(ctx context.Context, id string) (*DTO, error)
}

type service struct {
	repository Repository
	logger     *zap.SugaredLogger
}

func NewService(repository Repository, logger *zap.SugaredLogger) (Service, error) {
	if repository == nil {
		return nil, errors.New("invalid repository")
	}
	if logger == nil {
		return nil, errors.New("invalid logger")
	}

	return &service{repository: repository, logger: logger}, nil
}

func (s *service) GetSupportById(ctx context.Context, id string) (*DTO, error) {
	support, err := s.repository.GetSupportById(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to get support: %v", err)
		return nil, err
	}

	return MapToDTO(support), nil
}
