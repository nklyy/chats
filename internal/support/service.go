package support

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"noname-realtime-support-chat/pkg/jwt"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type Service interface {
	Registration(ctx context.Context, dto *RegistrationDTO) (*string, error)
	Login(ctx context.Context, dto *LoginDTO) (*string, error)
	GetSupportById(ctx context.Context, id string) (*DTO, error)
}

type service struct {
	repository Repository
	logger     *zap.SugaredLogger
	salt       int
	jwtSvc     jwt.Service
}

func NewService(repository Repository, logger *zap.SugaredLogger, salt *int, jwtSvc jwt.Service) (Service, error) {
	if repository == nil {
		return nil, errors.New("invalid repository")
	}
	if logger == nil {
		return nil, errors.New("invalid logger")
	}
	if salt == nil {
		return nil, errors.New("invalid salt")
	}
	if jwtSvc == nil {
		return nil, errors.New("invalid jwt service")
	}

	return &service{repository: repository, logger: logger, salt: *salt, jwtSvc: jwtSvc}, nil
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

func (s *service) Registration(ctx context.Context, dto *RegistrationDTO) (*string, error) {
	support, err := NewSupport(dto.Email, dto.Name, dto.Password, &s.salt)
	if err != nil {
		s.logger.Errorf("failed to create new support %v", err)
		return nil, err
	}

	id, err := s.repository.CreateSupport(ctx, support)
	if err != nil {
		s.logger.Errorf("failed to save support %v", err)
		return nil, err
	}

	return &id, nil
}

func (s *service) Login(ctx context.Context, dto *LoginDTO) (*string, error) {
	support, err := s.repository.GetSupportByEmail(ctx, dto.Email)
	if err != nil {
		s.logger.Errorf("failed to find support %v", err)
		return nil, err
	}

	cp, err := support.CheckPassword(dto.Password)
	if !cp {
		s.logger.Errorf("failed to check password %v", err)
		return nil, err
	}

	token, err := s.jwtSvc.CreateJWT(support.Name, "support")
	if err != nil {
		s.logger.Errorf("failed to create jwt token %v", err)
		return nil, err
	}

	support.SetOnline()

	return token, nil
}
