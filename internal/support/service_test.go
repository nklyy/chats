package support_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"noname-realtime-support-chat/internal/support"
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
		expect     func(*testing.T, support.Service, error)
	}{
		{
			name:       "should return service",
			repository: mock_support.NewMockRepository(controller),
			logger:     &zap.SugaredLogger{},
			salt:       &salt,
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
			expect: func(t *testing.T, s support.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid salt")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := support.NewService(tc.repository, tc.logger, tc.salt)
			tc.expect(t, svc, err)
		})
	}
}

func TestService_GetSupportById(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockRepo := mock_support.NewMockRepository(controller)
	salt := 10

	newLogger, _ := logger.NewLogger("development")
	zapLogger, _ := newLogger.SetupZapLogger()

	service, _ := support.NewService(mockRepo, zapLogger, &salt)

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

func TestService_GetSupportByEmail(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockRepo := mock_support.NewMockRepository(controller)
	salt := 10

	newLogger, _ := logger.NewLogger("development")
	zapLogger, _ := newLogger.SetupZapLogger()

	service, _ := support.NewService(mockRepo, zapLogger, &salt)

	supportEntity, _ := support.NewSupport("email", "name", "password", &salt)
	supportDTO := support.MapToDTO(supportEntity)

	tests := []struct {
		name   string
		ctx    context.Context
		email  string
		setup  func(context.Context, string)
		expect func(*testing.T, *support.DTO, error)
	}{
		{
			name:  "should return support",
			ctx:   context.Background(),
			email: "email",
			setup: func(ctx context.Context, email string) {
				mockRepo.EXPECT().GetSupportByEmail(ctx, email).Return(supportEntity, nil)
			},
			expect: func(t *testing.T, dto *support.DTO, err error) {
				assert.NotNil(t, dto)
				assert.Nil(t, err)
				assert.Equal(t, dto.Email, supportDTO.Email)
			},
		},
		{
			name:  "should return not found",
			ctx:   context.Background(),
			email: "email",
			setup: func(ctx context.Context, email string) {
				mockRepo.EXPECT().GetSupportByEmail(ctx, email).Return(nil, support.ErrNotFound)
			},
			expect: func(t *testing.T, dto *support.DTO, err error) {
				assert.Nil(t, dto)
				assert.Equal(t, support.ErrNotFound, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.email)
			s, err := service.GetSupportByEmail(tc.ctx, tc.email)
			tc.expect(t, s, err)
		})
	}
}

func TestService_CreateSupport(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockRepo := mock_support.NewMockRepository(controller)
	salt := 10

	newLogger, _ := logger.NewLogger("development")
	zapLogger, _ := newLogger.SetupZapLogger()

	service, _ := support.NewService(mockRepo, zapLogger, &salt)

	supportEntity, _ := support.NewSupport("email", "name", "password", &salt)
	supportDTO := support.MapToDTO(supportEntity)

	var emptyStr string

	tests := []struct {
		testName string
		ctx      context.Context
		email    string
		name     string
		password string
		setup    func(context.Context, string, string, string)
		expect   func(*testing.T, *support.DTO, error)
	}{
		{
			testName: "should return support",
			ctx:      context.Background(),
			email:    "email",
			name:     "name",
			password: "password",
			setup: func(ctx context.Context, email, name, password string) {
				mockRepo.EXPECT().CreateSupport(ctx, gomock.Any()).Return(supportEntity.ID.Hex(), nil)
			},
			expect: func(t *testing.T, dto *support.DTO, err error) {
				assert.NotNil(t, dto)
				assert.Nil(t, err)
				assert.Equal(t, dto.Email, supportDTO.Email)
			},
		},
		{
			testName: "should return support",
			ctx:      context.Background(),
			email:    "email",
			name:     "name",
			password: "password",
			setup: func(ctx context.Context, email, name, password string) {
				mockRepo.EXPECT().CreateSupport(ctx, gomock.Any()).Return(emptyStr, support.ErrFailedSaveSupport)
			},
			expect: func(t *testing.T, dto *support.DTO, err error) {
				assert.Empty(t, dto)
				assert.NotNil(t, err)
				assert.EqualError(t, err, support.ErrFailedSaveSupport.Error())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.email, tc.name, tc.password)
			s, err := service.CreateSupport(tc.ctx, tc.email, tc.name, tc.password)
			tc.expect(t, s, err)
		})
	}
}
