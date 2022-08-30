package auth

import (
	"context"
	"errors"
	"support-chat/internal/user"
	"support-chat/pkg/jwt"

	"go.uber.org/zap"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go
type Service interface {
	Registration(ctx context.Context, dto *RegistrationDTO) (*string, error)
	Login(ctx context.Context, dto *LoginDTO) (*string, *string, error)
	Refresh(ctx context.Context, dto *RefreshDTO) (*string, *string, error)
	Logout(ctx context.Context, dto *LogoutDTO) error
	Check(ctx context.Context, dto *CheckDTO) (*CheckResponseDTO, error)
}

type service struct {
	userSvc user.Service
	jwtSvc  jwt.Service
	logger  *zap.SugaredLogger
}

func NewService(userSvc user.Service, jwtSvc jwt.Service, logger *zap.SugaredLogger) (Service, error) {
	if userSvc == nil {
		return nil, errors.New("[user_auth_service] invalid user service")
	}
	if jwtSvc == nil {
		return nil, errors.New("[user_auth_service] invalid jwt service")
	}
	if logger == nil {
		return nil, errors.New("[user_auth_service] invalid logger")
	}

	return &service{userSvc: userSvc, logger: logger, jwtSvc: jwtSvc}, nil
}

func (s *service) Registration(ctx context.Context, dto *RegistrationDTO) (*string, error) {
	userDto, err := s.userSvc.CreateUser(ctx, dto.Email, dto.Name, dto.Password)
	if err != nil {
		s.logger.Errorf("failed to save user %v", err)
		return nil, err
	}

	return &userDto.ID, nil
}

func (s *service) Login(ctx context.Context, dto *LoginDTO) (*string, *string, error) {
	userDto, err := s.userSvc.GetUserByEmail(ctx, dto.Email, true)
	if err != nil {
		s.logger.Errorf("failed to find user %v", err)
		return nil, nil, err
	}

	userEntity, err := user.MapToEntity(userDto)
	if err != nil {
		s.logger.Errorf("failed to conver dto %v", err)
		return nil, nil, err
	}

	cp, err := userEntity.CheckPassword(dto.Password)
	if !cp {
		s.logger.Errorf("failed to check password %v", err)
		return nil, nil, err
	}

	accessToken, refreshToken, err := s.jwtSvc.CreateTokens(ctx, userDto.ID, userDto.Support)
	if err != nil {
		s.logger.Errorf("failed to create jwt token %v", err)
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (s *service) Refresh(ctx context.Context, dto *RefreshDTO) (*string, *string, error) {
	payload, err := s.jwtSvc.ParseToken(dto.Token, false)
	if err != nil {
		s.logger.Errorf("failed parse token %v", err)
		return nil, nil, err
	}

	userDto, err := s.userSvc.GetUserById(ctx, payload.Id, false)
	if err != nil {
		s.logger.Errorf("failed to find user %v", err)
		return nil, nil, err
	}

	err = s.jwtSvc.VerifyToken(ctx, payload, false)
	if err != nil {
		s.logger.Errorf("failed to verify token %v", err)
		return nil, nil, err
	}

	accessToken, refreshToken, err := s.jwtSvc.CreateTokens(ctx, payload.Id, userDto.Support)
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

func (s *service) Check(ctx context.Context, dto *CheckDTO) (*CheckResponseDTO, error) {
	payload, err := s.jwtSvc.ParseToken(dto.Token, true)
	if err != nil {
		s.logger.Errorf("failed parse token %v", err)
		return nil, err
	}

	err = s.jwtSvc.VerifyToken(ctx, payload, true)
	if err != nil {
		s.logger.Errorf("failed to verify token %v", err)
		return nil, err
	}

	u, err := s.userSvc.GetUserById(ctx, payload.Id, false)
	if err != nil {
		s.logger.Errorf("failed to get user %v", err)
		return nil, err
	}

	var isRoom bool
	if u.RoomName == nil {
		isRoom = false
	} else {
		isRoom = true
	}

	return &CheckResponseDTO{
		UserId: payload.Id,
		Role:   payload.Role,
		IsRoom: isRoom,
	}, nil
}
