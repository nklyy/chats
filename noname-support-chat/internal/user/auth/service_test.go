package auth_test

import (
	"context"
	gjwt "github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"noname-realtime-support-chat/internal/user"
	"noname-realtime-support-chat/internal/user/auth"
	mock_user "noname-realtime-support-chat/internal/user/mocks"
	"noname-realtime-support-chat/pkg/jwt"
	mock_jwt "noname-realtime-support-chat/pkg/jwt/mocks"
	"noname-realtime-support-chat/pkg/logger"

	"testing"
)

func TestNewService(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tests := []struct {
		name    string
		userSvc user.Service
		logger  *zap.SugaredLogger
		jwtSvc  jwt.Service
		expect  func(*testing.T, auth.Service, error)
	}{
		{
			name:    "should return service",
			userSvc: mock_user.NewMockService(controller),
			logger:  &zap.SugaredLogger{},
			jwtSvc:  mock_jwt.NewMockService(controller),
			expect: func(t *testing.T, service auth.Service, err error) {
				assert.NotNil(t, service)
				assert.Nil(t, err)
			},
		},
		{
			name:    "should return invalid user service",
			userSvc: nil,
			logger:  &zap.SugaredLogger{},
			jwtSvc:  mock_jwt.NewMockService(controller),
			expect: func(t *testing.T, service auth.Service, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid user service")
			},
		},
		{
			name:    "should return invalid logger",
			userSvc: mock_user.NewMockService(controller),
			logger:  nil,
			jwtSvc:  mock_jwt.NewMockService(controller),
			expect: func(t *testing.T, service auth.Service, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid logger")
			},
		},
		{
			name:    "should return invalid jwt service",
			userSvc: mock_user.NewMockService(controller),
			logger:  &zap.SugaredLogger{},
			jwtSvc:  nil,
			expect: func(t *testing.T, service auth.Service, err error) {
				assert.Nil(t, service)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid jwt service")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := auth.NewService(tc.userSvc, tc.logger, tc.jwtSvc)
			tc.expect(t, svc, err)
		})
	}
}

func TestService_Registration(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockUserSvc := mock_user.NewMockService(controller)
	mockJwt := mock_jwt.NewMockService(controller)

	salt := 10

	newLogger, _ := logger.NewLogger("development")
	zapLogger, _ := newLogger.SetupZapLogger()

	service, _ := auth.NewService(mockUserSvc, zapLogger, mockJwt)

	userEntity, _ := user.NewUser("email", "name", "password", &salt)
	userDto := user.MapToDTO(userEntity)

	tests := []struct {
		name   string
		ctx    context.Context
		dto    *auth.RegistrationDTO
		setup  func(context.Context, *auth.RegistrationDTO)
		expect func(*testing.T, *string, error)
	}{
		{
			name: "should return registered user id",
			ctx:  context.Background(),
			dto: &auth.RegistrationDTO{
				Email:    "email",
				Name:     "name",
				Password: "password",
			},
			setup: func(ctx context.Context, dto *auth.RegistrationDTO) {
				mockUserSvc.EXPECT().CreateUser(ctx, dto.Email, dto.Name, dto.Password).Return(userDto, nil)
			},
			expect: func(t *testing.T, s *string, err error) {
				assert.NotNil(t, s)
				assert.Nil(t, err)
				assert.Equal(t, userEntity.ID.Hex(), *s)
			},
		},
		{
			name: "should return failed to create user",
			ctx:  context.Background(),
			dto: &auth.RegistrationDTO{
				Email:    "email",
				Name:     "name",
				Password: "password",
			},
			setup: func(ctx context.Context, dto *auth.RegistrationDTO) {
				mockUserSvc.EXPECT().CreateUser(ctx, dto.Email, dto.Name, dto.Password).Return(nil, user.ErrFailedSaveUser)
			},
			expect: func(t *testing.T, s *string, err error) {
				assert.Empty(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, user.ErrFailedSaveUser.Error())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.dto)
			s, err := service.Registration(tc.ctx, tc.dto)
			tc.expect(t, s, err)
		})
	}
}

func TestService_Login(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockUserSvc := mock_user.NewMockService(controller)
	mockJwt := mock_jwt.NewMockService(controller)

	salt := 10

	newLogger, _ := logger.NewLogger("development")
	zapLogger, _ := newLogger.SetupZapLogger()

	service, _ := auth.NewService(mockUserSvc, zapLogger, mockJwt)

	userEntity, _ := user.NewUser("email", "name", "password", &salt)
	userDto := user.MapToDTO(userEntity)

	tokenAccess := "tokenAccess"
	tokenRefresh := "tokenRefresh"
	var emptyStr string

	tests := []struct {
		name         string
		ctx          context.Context
		dto          *auth.LoginDTO
		withPassword bool
		setup        func(context.Context, *auth.LoginDTO, bool)
		expect       func(*testing.T, *string, *string, error)
	}{
		{
			name: "should return jwt token",
			ctx:  context.Background(),
			dto: &auth.LoginDTO{
				Email:    "email",
				Password: "password",
			},
			withPassword: true,
			setup: func(ctx context.Context, dto *auth.LoginDTO, withPassword bool) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email, withPassword).Return(userDto, nil)
				mockJwt.EXPECT().CreateTokens(ctx, userDto.ID, false).Return(&tokenAccess, &tokenRefresh, nil)
			},
			expect: func(t *testing.T, a *string, r *string, err error) {
				assert.NotNil(t, a)
				assert.NotNil(t, r)
				assert.Nil(t, err)
				assert.Equal(t, *a, tokenAccess)
				assert.Equal(t, *r, tokenRefresh)
			},
		},
		{
			name: "should return failed to find user",
			ctx:  context.Background(),
			dto: &auth.LoginDTO{
				Email:    "email",
				Password: "password",
			},
			withPassword: true,
			setup: func(ctx context.Context, dto *auth.LoginDTO, withPassword bool) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email, withPassword).Return(nil, user.ErrNotFound)
			},
			expect: func(t *testing.T, a *string, r *string, err error) {
				assert.Empty(t, a)
				assert.Empty(t, r)
				assert.NotNil(t, err)
				assert.EqualError(t, err, user.ErrNotFound.Error())
			},
		},
		{
			name: "should return failed to create jwt token",
			ctx:  context.Background(),
			dto: &auth.LoginDTO{
				Email:    "email",
				Password: "password",
			},
			withPassword: true,
			setup: func(ctx context.Context, dto *auth.LoginDTO, withPassword bool) {
				mockUserSvc.EXPECT().GetUserByEmail(ctx, dto.Email, withPassword).Return(userDto, nil)
				mockJwt.EXPECT().CreateTokens(ctx, userDto.ID, false).Return(&emptyStr, &emptyStr, jwt.ErrFailedCreateTokens)
			},
			expect: func(t *testing.T, a *string, r *string, err error) {
				assert.Empty(t, a)
				assert.Empty(t, r)
				assert.NotNil(t, err)
				assert.EqualError(t, err, jwt.ErrFailedCreateTokens.Error())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.dto, tc.withPassword)
			a, r, err := service.Login(tc.ctx, tc.dto)
			tc.expect(t, a, r, err)
		})
	}
}

func TestService_Refresh(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockUserSvc := mock_user.NewMockService(controller)
	mockJwt := mock_jwt.NewMockService(controller)

	newLogger, _ := logger.NewLogger("development")
	zapLogger, _ := newLogger.SetupZapLogger()

	service, _ := auth.NewService(mockUserSvc, zapLogger, mockJwt)

	payload := jwt.Payload{
		Id:             "id",
		Role:           "role",
		Uid:            "uid",
		StandardClaims: gjwt.StandardClaims{},
	}
	tokenAccess := "tokenAccess"
	tokenRefresh := "tokenRefresh"

	salt := 10
	userEntity, _ := user.NewUser("email", "name", "password", &salt)
	userDto := user.MapToDTO(userEntity)

	var emptyStr string

	tests := []struct {
		name   string
		ctx    context.Context
		dto    *auth.RefreshDTO
		setup  func(context.Context, *auth.RefreshDTO)
		expect func(*testing.T, *string, *string, error)
	}{
		{
			name: "should return tokens",
			ctx:  context.Background(),
			dto: &auth.RefreshDTO{
				Token: "token",
			},
			setup: func(ctx context.Context, dto *auth.RefreshDTO) {
				mockJwt.EXPECT().ParseToken(dto.Token, false).Return(&payload, nil)
				mockUserSvc.EXPECT().GetUserById(ctx, payload.Id, false).Return(userDto, nil)
				mockJwt.EXPECT().VerifyToken(ctx, &payload, false).Return(nil)
				mockJwt.EXPECT().CreateTokens(ctx, payload.Id, userDto.Support).Return(&tokenAccess, &tokenRefresh, nil)
			},
			expect: func(t *testing.T, a *string, r *string, err error) {
				assert.NotNil(t, a)
				assert.NotNil(t, r)
				assert.Nil(t, err)
				assert.Equal(t, *a, tokenAccess)
				assert.Equal(t, *r, tokenRefresh)
			},
		},
		{
			name: "should return failed parse token",
			ctx:  context.Background(),
			dto: &auth.RefreshDTO{
				Token: "token",
			},
			setup: func(ctx context.Context, dto *auth.RefreshDTO) {
				mockJwt.EXPECT().ParseToken(dto.Token, false).Return(nil, jwt.ErrToken)
			},
			expect: func(t *testing.T, a *string, r *string, err error) {
				assert.Empty(t, a)
				assert.Empty(t, r)
				assert.NotNil(t, err)
				assert.EqualError(t, err, jwt.ErrToken.Error())
			},
		},
		{
			name: "should return failed to find user",
			ctx:  context.Background(),
			dto: &auth.RefreshDTO{
				Token: "token",
			},
			setup: func(ctx context.Context, dto *auth.RefreshDTO) {
				mockJwt.EXPECT().ParseToken(dto.Token, false).Return(&payload, nil)
				mockUserSvc.EXPECT().GetUserById(ctx, payload.Id, false).Return(nil, user.ErrNotFound)
			},
			expect: func(t *testing.T, a *string, r *string, err error) {
				assert.Empty(t, a)
				assert.Empty(t, r)
				assert.NotNil(t, err)
				assert.EqualError(t, err, user.ErrNotFound.Error())
			},
		},
		{
			name: "should return failed to verify token",
			ctx:  context.Background(),
			dto: &auth.RefreshDTO{
				Token: "token",
			},
			setup: func(ctx context.Context, dto *auth.RefreshDTO) {
				mockJwt.EXPECT().ParseToken(dto.Token, false).Return(&payload, nil)
				mockUserSvc.EXPECT().GetUserById(ctx, payload.Id, false).Return(userDto, nil)
				mockJwt.EXPECT().VerifyToken(ctx, &payload, false).Return(jwt.ErrNotFound)
			},
			expect: func(t *testing.T, a *string, r *string, err error) {
				assert.Empty(t, a)
				assert.Empty(t, r)
				assert.NotNil(t, err)
				assert.EqualError(t, err, jwt.ErrNotFound.Error())
			},
		},
		{
			name: "should return failed to create jwt token",
			ctx:  context.Background(),
			dto: &auth.RefreshDTO{
				Token: "token",
			},
			setup: func(ctx context.Context, dto *auth.RefreshDTO) {
				mockJwt.EXPECT().ParseToken(dto.Token, false).Return(&payload, nil)
				mockUserSvc.EXPECT().GetUserById(ctx, payload.Id, false).Return(userDto, nil)
				mockJwt.EXPECT().VerifyToken(ctx, &payload, false).Return(nil)
				mockJwt.EXPECT().CreateTokens(ctx, payload.Id, userDto.Support).Return(&emptyStr, &emptyStr, jwt.ErrFailedCreateTokens)
			},
			expect: func(t *testing.T, a *string, r *string, err error) {
				assert.Empty(t, a)
				assert.Empty(t, r)
				assert.NotNil(t, err)
				assert.EqualError(t, err, jwt.ErrFailedCreateTokens.Error())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.dto)
			a, r, err := service.Refresh(tc.ctx, tc.dto)
			tc.expect(t, a, r, err)
		})
	}
}

func TestService_Logout(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockUserSvc := mock_user.NewMockService(controller)
	mockJwt := mock_jwt.NewMockService(controller)

	newLogger, _ := logger.NewLogger("development")
	zapLogger, _ := newLogger.SetupZapLogger()

	service, _ := auth.NewService(mockUserSvc, zapLogger, mockJwt)

	payload := jwt.Payload{
		Id:             "id",
		Role:           "role",
		Uid:            "uid",
		StandardClaims: gjwt.StandardClaims{},
	}

	tests := []struct {
		name   string
		ctx    context.Context
		dto    *auth.LogoutDTO
		setup  func(context.Context, *auth.LogoutDTO)
		expect func(*testing.T, error)
	}{
		{
			name: "should logout user",
			ctx:  context.Background(),
			dto: &auth.LogoutDTO{
				Token: "token",
			},
			setup: func(ctx context.Context, dto *auth.LogoutDTO) {
				mockJwt.EXPECT().ParseToken(dto.Token, true).Return(&payload, nil)
				mockJwt.EXPECT().VerifyToken(ctx, &payload, true).Return(nil)
				mockJwt.EXPECT().DeleteTokens(ctx, &payload).Return(nil)
			},
			expect: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "should failed parse token",
			ctx:  context.Background(),
			dto: &auth.LogoutDTO{
				Token: "token",
			},
			setup: func(ctx context.Context, dto *auth.LogoutDTO) {
				mockJwt.EXPECT().ParseToken(dto.Token, true).Return(nil, jwt.ErrToken)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.EqualError(t, err, jwt.ErrToken.Error())
			},
		},
		{
			name: "should failed to verify token",
			ctx:  context.Background(),
			dto: &auth.LogoutDTO{
				Token: "token",
			},
			setup: func(ctx context.Context, dto *auth.LogoutDTO) {
				mockJwt.EXPECT().ParseToken(dto.Token, true).Return(&payload, nil)
				mockJwt.EXPECT().VerifyToken(ctx, &payload, true).Return(jwt.ErrNotFound)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.EqualError(t, err, jwt.ErrNotFound.Error())
			},
		},
		{
			name: "should failed to delete tokens",
			ctx:  context.Background(),
			dto: &auth.LogoutDTO{
				Token: "token",
			},
			setup: func(ctx context.Context, dto *auth.LogoutDTO) {
				mockJwt.EXPECT().ParseToken(dto.Token, true).Return(&payload, nil)
				mockJwt.EXPECT().VerifyToken(ctx, &payload, true).Return(nil)
				mockJwt.EXPECT().DeleteTokens(ctx, &payload).Return(jwt.ErrFailedDeleteToken)
			},
			expect: func(t *testing.T, err error) {
				assert.NotNil(t, err)
				assert.EqualError(t, err, jwt.ErrFailedDeleteToken.Error())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.dto)
			err := service.Logout(tc.ctx, tc.dto)
			tc.expect(t, err)
		})
	}
}
