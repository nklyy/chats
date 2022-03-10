package support_test

import (
	"context"
	gjwt "github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"noname-realtime-support-chat/internal/support"
	"noname-realtime-support-chat/internal/support/jwt"
	"noname-realtime-support-chat/internal/support/jwt/mocks"
	mock_support "noname-realtime-support-chat/internal/support/mocks"
	"noname-realtime-support-chat/pkg/logger"
	"testing"
)

func TestNewService(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	salt := 10

	tests := []struct {
		name       string
		repository support.Repository
		logger     *zap.SugaredLogger
		salt       *int
		jwtSvc     jwt.Service
		expect     func(*testing.T, support.Service, error)
	}{
		{
			name:       "should return service",
			repository: mock_support.NewMockRepository(controller),
			logger:     &zap.SugaredLogger{},
			salt:       &salt,
			jwtSvc:     mock_jwt.NewMockService(controller),
			expect: func(t *testing.T, s support.Service, err error) {
				assert.NotNil(t, s)
				assert.Nil(t, err)
			},
		},
		{
			name:       "should return invalid repository",
			repository: nil,
			logger:     &zap.SugaredLogger{},
			salt:       &salt,
			jwtSvc:     mock_jwt.NewMockService(controller),
			expect: func(t *testing.T, s support.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid repository")
			},
		},
		{
			name:       "should return invalid logger",
			repository: mock_support.NewMockRepository(controller),
			logger:     nil,
			salt:       &salt,
			jwtSvc:     mock_jwt.NewMockService(controller),
			expect: func(t *testing.T, s support.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid logger")
			},
		},
		{
			name:       "should return invalid salt",
			repository: mock_support.NewMockRepository(controller),
			logger:     &zap.SugaredLogger{},
			salt:       nil,
			jwtSvc:     mock_jwt.NewMockService(controller),
			expect: func(t *testing.T, s support.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid salt")
			},
		},
		{
			name:       "should return invalid jwt service",
			repository: mock_support.NewMockRepository(controller),
			logger:     &zap.SugaredLogger{},
			salt:       &salt,
			jwtSvc:     nil,
			expect: func(t *testing.T, s support.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid jwt service")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := support.NewService(tc.repository, tc.logger, tc.salt, tc.jwtSvc)
			tc.expect(t, svc, err)
		})
	}
}

func TestService_GetSupportById(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockRepo := mock_support.NewMockRepository(controller)
	mockJwt := mock_jwt.NewMockService(controller)

	salt := 10

	newLogger, _ := logger.NewLogger("development")
	zapLogger, _ := newLogger.SetupZapLogger()

	service, _ := support.NewService(mockRepo, zapLogger, &salt, mockJwt)

	supportEntity, _ := support.NewSupport("email", "name", "password", &salt)
	supportDTO := support.MapToDTO(supportEntity)

	tests := []struct {
		name   string
		ctx    context.Context
		id     string
		setup  func(context.Context, string)
		expect func(*testing.T, *support.DTO, error)
	}{
		{
			name: "should return support",
			ctx:  context.Background(),
			id:   supportDTO.ID,
			setup: func(ctx context.Context, id string) {
				mockRepo.EXPECT().GetSupportById(ctx, id).Return(supportEntity, nil)
			},
			expect: func(t *testing.T, dto *support.DTO, err error) {
				assert.NotNil(t, dto)
				assert.Nil(t, err)
				assert.Equal(t, dto.ID, supportDTO.ID)
			},
		},
		{
			name: "should return not found",
			ctx:  context.Background(),
			id:   "incorrect_id",
			setup: func(ctx context.Context, id string) {
				mockRepo.EXPECT().GetSupportById(ctx, id).Return(nil, support.ErrNotFound)
			},
			expect: func(t *testing.T, dto *support.DTO, err error) {
				assert.Nil(t, dto)
				assert.Equal(t, support.ErrNotFound, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.id)
			s, err := service.GetSupportById(tc.ctx, tc.id)
			tc.expect(t, s, err)
		})
	}
}

func TestService_Registration(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockRepo := mock_support.NewMockRepository(controller)
	mockJwt := mock_jwt.NewMockService(controller)

	salt := 10

	newLogger, _ := logger.NewLogger("development")
	zapLogger, _ := newLogger.SetupZapLogger()

	service, _ := support.NewService(mockRepo, zapLogger, &salt, mockJwt)

	supportEntity, _ := support.NewSupport("email", "name", "password", &salt)
	var emptyStr string

	tests := []struct {
		name   string
		ctx    context.Context
		dto    *support.RegistrationDTO
		setup  func(context.Context, *support.RegistrationDTO)
		expect func(*testing.T, *string, error)
	}{
		{
			name: "should return registered support id",
			ctx:  context.Background(),
			dto: &support.RegistrationDTO{
				Email:    "email",
				Name:     "name",
				Password: "password",
			},
			setup: func(ctx context.Context, dto *support.RegistrationDTO) {
				mockRepo.EXPECT().CreateSupport(ctx, gomock.Any()).Return(supportEntity.ID.Hex(), nil)
			},
			expect: func(t *testing.T, s *string, err error) {
				assert.NotNil(t, s)
				assert.Nil(t, err)
				assert.Equal(t, supportEntity.ID.Hex(), *s)
			},
		},
		{
			name: "should return failed to create support",
			ctx:  context.Background(),
			dto: &support.RegistrationDTO{
				Email:    "email",
				Name:     "name",
				Password: "password",
			},
			setup: func(ctx context.Context, dto *support.RegistrationDTO) {
				mockRepo.EXPECT().CreateSupport(ctx, gomock.Any()).Return(emptyStr, support.ErrFailedSaveSupport)
			},
			expect: func(t *testing.T, s *string, err error) {
				assert.Empty(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, support.ErrFailedSaveSupport.Error())
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

	mockRepo := mock_support.NewMockRepository(controller)
	mockJwt := mock_jwt.NewMockService(controller)

	salt := 10

	newLogger, _ := logger.NewLogger("development")
	zapLogger, _ := newLogger.SetupZapLogger()

	service, _ := support.NewService(mockRepo, zapLogger, &salt, mockJwt)

	supportEntity, _ := support.NewSupport("email", "name", "password", &salt)
	tokenAccess := "tokenAccess"
	tokenRefresh := "tokenRefresh"
	var emptyStr string

	tests := []struct {
		name   string
		ctx    context.Context
		dto    *support.LoginDTO
		setup  func(context.Context, *support.LoginDTO)
		expect func(*testing.T, *string, *string, error)
	}{
		{
			name: "should return jwt token",
			ctx:  context.Background(),
			dto: &support.LoginDTO{
				Email:    "email",
				Password: "password",
			},
			setup: func(ctx context.Context, dto *support.LoginDTO) {
				mockRepo.EXPECT().GetSupportByEmail(ctx, dto.Email).Return(supportEntity, nil)
				mockJwt.EXPECT().CreateTokens(ctx, supportEntity.ID.Hex(), "support").Return(&tokenAccess, &tokenRefresh, nil)
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
			name: "should return failed to find support",
			ctx:  context.Background(),
			dto: &support.LoginDTO{
				Email:    "email",
				Password: "password",
			},
			setup: func(ctx context.Context, dto *support.LoginDTO) {
				mockRepo.EXPECT().GetSupportByEmail(ctx, dto.Email).Return(nil, support.ErrNotFound)
			},
			expect: func(t *testing.T, a *string, r *string, err error) {
				assert.Empty(t, a)
				assert.Empty(t, r)
				assert.NotNil(t, err)
				assert.EqualError(t, err, support.ErrNotFound.Error())
			},
		},
		{
			name: "should return failed to create jwt token",
			ctx:  context.Background(),
			dto: &support.LoginDTO{
				Email:    "email",
				Password: "password",
			},
			setup: func(ctx context.Context, dto *support.LoginDTO) {
				mockRepo.EXPECT().GetSupportByEmail(ctx, dto.Email).Return(supportEntity, nil)
				mockJwt.EXPECT().CreateTokens(ctx, supportEntity.ID.Hex(), "support").Return(&emptyStr, &emptyStr, jwt.ErrFailedCreateTokens)
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
			a, r, err := service.Login(tc.ctx, tc.dto)
			tc.expect(t, a, r, err)
		})
	}
}

func TestService_Refresh(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockRepo := mock_support.NewMockRepository(controller)
	mockJwt := mock_jwt.NewMockService(controller)

	salt := 10

	newLogger, _ := logger.NewLogger("development")
	zapLogger, _ := newLogger.SetupZapLogger()

	service, _ := support.NewService(mockRepo, zapLogger, &salt, mockJwt)

	payload := jwt.Payload{
		Id:             "id",
		Role:           "role",
		Uid:            "uid",
		StandardClaims: gjwt.StandardClaims{},
	}
	tokenAccess := "tokenAccess"
	tokenRefresh := "tokenRefresh"
	var emptyStr string

	tests := []struct {
		name   string
		ctx    context.Context
		dto    *support.RefreshDTO
		setup  func(context.Context, *support.RefreshDTO)
		expect func(*testing.T, *string, *string, error)
	}{
		{
			name: "should return tokens",
			ctx:  context.Background(),
			dto: &support.RefreshDTO{
				Token: "token",
			},
			setup: func(ctx context.Context, dto *support.RefreshDTO) {
				mockJwt.EXPECT().ParseToken(dto.Token, false).Return(&payload, nil)
				mockJwt.EXPECT().VerifyToken(ctx, &payload, false).Return(nil)
				mockJwt.EXPECT().CreateTokens(ctx, payload.Id, "support").Return(&tokenAccess, &tokenRefresh, nil)
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
			dto: &support.RefreshDTO{
				Token: "token",
			},
			setup: func(ctx context.Context, dto *support.RefreshDTO) {
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
			name: "should return failed to verify token",
			ctx:  context.Background(),
			dto: &support.RefreshDTO{
				Token: "token",
			},
			setup: func(ctx context.Context, dto *support.RefreshDTO) {
				mockJwt.EXPECT().ParseToken(dto.Token, false).Return(&payload, nil)
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
			dto: &support.RefreshDTO{
				Token: "token",
			},
			setup: func(ctx context.Context, dto *support.RefreshDTO) {
				mockJwt.EXPECT().ParseToken(dto.Token, false).Return(&payload, nil)
				mockJwt.EXPECT().VerifyToken(ctx, &payload, false).Return(nil)
				mockJwt.EXPECT().CreateTokens(ctx, payload.Id, "support").Return(&emptyStr, &emptyStr, jwt.ErrFailedCreateTokens)
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

	mockRepo := mock_support.NewMockRepository(controller)
	mockJwt := mock_jwt.NewMockService(controller)

	salt := 10

	newLogger, _ := logger.NewLogger("development")
	zapLogger, _ := newLogger.SetupZapLogger()

	service, _ := support.NewService(mockRepo, zapLogger, &salt, mockJwt)

	payload := jwt.Payload{
		Id:             "id",
		Role:           "role",
		Uid:            "uid",
		StandardClaims: gjwt.StandardClaims{},
	}

	tests := []struct {
		name   string
		ctx    context.Context
		dto    *support.LogoutDTO
		setup  func(context.Context, *support.LogoutDTO)
		expect func(*testing.T, error)
	}{
		{
			name: "should logout support",
			ctx:  context.Background(),
			dto: &support.LogoutDTO{
				Token: "token",
			},
			setup: func(ctx context.Context, dto *support.LogoutDTO) {
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
			dto: &support.LogoutDTO{
				Token: "token",
			},
			setup: func(ctx context.Context, dto *support.LogoutDTO) {
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
			dto: &support.LogoutDTO{
				Token: "token",
			},
			setup: func(ctx context.Context, dto *support.LogoutDTO) {
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
			dto: &support.LogoutDTO{
				Token: "token",
			},
			setup: func(ctx context.Context, dto *support.LogoutDTO) {
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
