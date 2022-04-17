package user_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"noname-realtime-support-chat/internal/user"
	mock_user "noname-realtime-support-chat/internal/user/mocks"
	"noname-realtime-support-chat/pkg/logger"

	"testing"
)

func TestNewService(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	salt := 10

	tests := []struct {
		name       string
		repository user.Repository
		logger     *zap.SugaredLogger
		salt       *int
		expect     func(*testing.T, user.Service, error)
	}{
		{
			name:       "should return service",
			repository: mock_user.NewMockRepository(controller),
			logger:     &zap.SugaredLogger{},
			salt:       &salt,
			expect: func(t *testing.T, s user.Service, err error) {
				assert.NotNil(t, s)
				assert.Nil(t, err)
			},
		},
		{
			name:       "should return invalid repository",
			repository: nil,
			logger:     &zap.SugaredLogger{},
			salt:       &salt,
			expect: func(t *testing.T, s user.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid repository")
			},
		},
		{
			name:       "should return invalid logger",
			repository: mock_user.NewMockRepository(controller),
			logger:     nil,
			salt:       &salt,
			expect: func(t *testing.T, s user.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid logger")
			},
		},
		{
			name:       "should return invalid salt",
			repository: mock_user.NewMockRepository(controller),
			logger:     &zap.SugaredLogger{},
			salt:       nil,
			expect: func(t *testing.T, s user.Service, err error) {
				assert.Nil(t, s)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "invalid salt")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := user.NewService(tc.repository, tc.logger, tc.salt)
			tc.expect(t, svc, err)
		})
	}
}

func TestService_GetUserById(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockRepo := mock_user.NewMockRepository(controller)
	salt := 10

	newLogger, _ := logger.NewLogger("development")
	zapLogger, _ := newLogger.SetupZapLogger()

	service, _ := user.NewService(mockRepo, zapLogger, &salt)

	userEntity, _ := user.NewUser("email", "name", "password", &salt)
	userDTO := user.MapToDTO(userEntity)

	tests := []struct {
		name         string
		ctx          context.Context
		id           string
		withPassword bool
		setup        func(context.Context, string)
		expect       func(*testing.T, *user.DTO, error)
	}{
		{
			name:         "should return user",
			ctx:          context.Background(),
			id:           userDTO.ID,
			withPassword: false,
			setup: func(ctx context.Context, id string) {
				objId, _ := primitive.ObjectIDFromHex(id)
				mockRepo.EXPECT().GetUser(ctx, bson.M{"_id": objId}).Return(userEntity, nil)
			},
			expect: func(t *testing.T, dto *user.DTO, err error) {
				assert.NotNil(t, dto)
				assert.Nil(t, err)
				assert.Equal(t, dto.ID, userDTO.ID)
			},
		},
		{
			name:         "should return not found",
			ctx:          context.Background(),
			id:           userDTO.ID,
			withPassword: false,
			setup: func(ctx context.Context, id string) {
				objId, _ := primitive.ObjectIDFromHex(id)
				mockRepo.EXPECT().GetUser(ctx, bson.M{"_id": objId}).Return(nil, user.ErrNotFound)
			},
			expect: func(t *testing.T, dto *user.DTO, err error) {
				assert.Nil(t, dto)
				assert.Equal(t, user.ErrNotFound, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.id)
			s, err := service.GetUserById(tc.ctx, tc.id, tc.withPassword)
			tc.expect(t, s, err)
		})
	}
}

func TestService_GetUserByEmail(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockRepo := mock_user.NewMockRepository(controller)
	salt := 10

	newLogger, _ := logger.NewLogger("development")
	zapLogger, _ := newLogger.SetupZapLogger()

	service, _ := user.NewService(mockRepo, zapLogger, &salt)

	userEntity, _ := user.NewUser("email", "name", "password", &salt)
	userDTO := user.MapToDTO(userEntity)

	tests := []struct {
		name         string
		ctx          context.Context
		email        string
		withPassword bool
		setup        func(context.Context, string)
		expect       func(*testing.T, *user.DTO, error)
	}{
		{
			name:         "should return user",
			ctx:          context.Background(),
			email:        "email",
			withPassword: false,
			setup: func(ctx context.Context, email string) {
				mockRepo.EXPECT().GetUser(ctx, bson.M{"email": email}).Return(userEntity, nil)
			},
			expect: func(t *testing.T, dto *user.DTO, err error) {
				assert.NotNil(t, dto)
				assert.Nil(t, err)
				assert.Equal(t, dto.Email, userDTO.Email)
			},
		},
		{
			name:         "should return not found",
			ctx:          context.Background(),
			email:        "email",
			withPassword: false,
			setup: func(ctx context.Context, email string) {
				mockRepo.EXPECT().GetUser(ctx, bson.M{"email": email}).Return(nil, user.ErrNotFound)
			},
			expect: func(t *testing.T, dto *user.DTO, err error) {
				assert.Nil(t, dto)
				assert.Equal(t, user.ErrNotFound, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.email)
			s, err := service.GetUserByEmail(tc.ctx, tc.email, tc.withPassword)
			tc.expect(t, s, err)
		})
	}
}

func TestService_CreateUser(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockRepo := mock_user.NewMockRepository(controller)
	salt := 10

	newLogger, _ := logger.NewLogger("development")
	zapLogger, _ := newLogger.SetupZapLogger()

	service, _ := user.NewService(mockRepo, zapLogger, &salt)

	userEntity, _ := user.NewUser("email", "name", "password", &salt)
	userDTO := user.MapToDTO(userEntity)

	var emptyStr string

	tests := []struct {
		testName string
		ctx      context.Context
		email    string
		name     string
		password string
		setup    func(context.Context, string, string, string)
		expect   func(*testing.T, *user.DTO, error)
	}{
		{
			testName: "should return user",
			ctx:      context.Background(),
			email:    "email",
			name:     "name",
			password: "password",
			setup: func(ctx context.Context, email, name, password string) {
				mockRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(userEntity.ID.Hex(), nil)
			},
			expect: func(t *testing.T, dto *user.DTO, err error) {
				assert.NotNil(t, dto)
				assert.Nil(t, err)
				assert.Equal(t, dto.Email, userDTO.Email)
			},
		},
		{
			testName: "should return user",
			ctx:      context.Background(),
			email:    "email",
			name:     "name",
			password: "password",
			setup: func(ctx context.Context, email, name, password string) {
				mockRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(emptyStr, user.ErrFailedSaveUser)
			},
			expect: func(t *testing.T, dto *user.DTO, err error) {
				assert.Empty(t, dto)
				assert.NotNil(t, err)
				assert.EqualError(t, err, user.ErrFailedSaveUser.Error())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.email, tc.name, tc.password)
			s, err := service.CreateUser(tc.ctx, tc.email, tc.name, tc.password)
			tc.expect(t, s, err)
		})
	}
}
