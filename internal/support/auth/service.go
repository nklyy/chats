package auth

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"noname-realtime-support-chat/internal/support"
	"noname-realtime-support-chat/internal/support/jwt"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type Service interface {
	Registration(ctx context.Context, dto *RegistrationDTO) (*string, error)
	Login(ctx context.Context, dto *LoginDTO) (*string, *string, error)
	Refresh(ctx context.Context, dto *RefreshDTO) (*string, *string, error)
	Logout(ctx context.Context, dto *LogoutDTO) error
}

type service struct {
	supportSvc support.Service
	logger     *zap.SugaredLogger
	jwtSvc     jwt.Service
}

func NewService(supportSvc support.Service, logger *zap.SugaredLogger, jwtSvc jwt.Service) (Service, error) {
	if supportSvc == nil {
		return nil, errors.New("invalid support service")
	}
	if logger == nil {
		return nil, errors.New("invalid logger")
	}
	if jwtSvc == nil {
		return nil, errors.New("invalid jwt service")
	}

	return &service{supportSvc: supportSvc, logger: logger, jwtSvc: jwtSvc}, nil
}

func (s *service) Registration(ctx context.Context, dto *RegistrationDTO) (*string, error) {
	supportDto, err := s.supportSvc.CreateSupport(ctx, dto.Email, dto.Name, dto.Password)
	if err != nil {
		s.logger.Errorf("failed to save support %v", err)
		return nil, err
	}

	return &supportDto.ID, nil
}

func (s *service) Login(ctx context.Context, dto *LoginDTO) (*string, *string, error) {
	supportDto, err := s.supportSvc.GetSupportByEmail(ctx, dto.Email, true)
	if err != nil {
		s.logger.Errorf("failed to find support %v", err)
		return nil, nil, err
	}

	supportEntity, err := support.MapToEntity(supportDto)
	if err != nil {
		s.logger.Errorf("failed to conver dto %v", err)
		return nil, nil, err
	}

	cp, err := supportEntity.CheckPassword(dto.Password)
	if !cp {
		s.logger.Errorf("failed to check password %v", err)
		return nil, nil, err
	}

	accessToken, refreshToken, err := s.jwtSvc.CreateTokens(ctx, supportDto.ID, "support")
	if err != nil {
		s.logger.Errorf("failed to create jwt token %v", err)
		return nil, nil, err
	}

	//support.SetOnline()

	return accessToken, refreshToken, nil
}

func (s *service) Refresh(ctx context.Context, dto *RefreshDTO) (*string, *string, error) {
	payload, err := s.jwtSvc.ParseToken(dto.Token, false)
	if err != nil {
		s.logger.Errorf("failed parse token %v", err)
		return nil, nil, err
	}

	err = s.jwtSvc.VerifyToken(ctx, payload, false)
	if err != nil {
		s.logger.Errorf("failed to verify token %v", err)
		return nil, nil, err
	}

	accessToken, refreshToken, err := s.jwtSvc.CreateTokens(ctx, payload.Id, "support")
	if err != nil {
		s.logger.Errorf("failed to create jwt token %v", err)
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (s *service) Logout(ctx context.Context, dto *LogoutDTO) error {
	payload, err := s.jwtSvc.ParseToken(dto.Token, true)
	if err != nil {
		s.logger.Errorf("failed parse token %v", err)
		return err
	}

	err = s.jwtSvc.VerifyToken(ctx, payload, true)
	if err != nil {
		s.logger.Errorf("failed to verify token %v", err)
		return err
	}

	err = s.jwtSvc.DeleteTokens(ctx, payload)
	if err != nil {
		s.logger.Errorf("failed to delete tokens %v", err)
		return err
	}

	return nil
}
