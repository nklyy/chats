package user_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"noname-realtime-support-chat/internal/user"
	mock_user "noname-realtime-support-chat/internal/user/mocks"
	"noname-realtime-support-chat/pkg/logger"
	"testing"
)

func TestNewService(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	salt := "salt"

	tests := []struct {
		name       string
		repository user.Repository
		logger     *zap.SugaredLogger
		salt       string
		expect     func(*testing.T, user.Service, error)
	}{
		{
			name:       "should return service",
			repository: mock_user.NewMockRepository(controller),
			logger:     &zap.SugaredLogger{},
			salt:       salt,
			expect: func(t *testing.T, s user.Service, err error) {
				assert.NotNil(t, s)
				assert.Nil(t, err)
			},
		},
		{
			name:       "should return invalid repository",
			repository: nil,
			logger:     &zap.SugaredLogger{},
			salt:       salt,
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
			salt:       salt,
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
			salt:       "",
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
	salt := "salt"
	ipAddr := "127.0.0.1"

	newLogger, _ := logger.NewLogger("development")
	zapLogger, _ := newLogger.SetupZapLogger()

	service, _ := user.NewService(mockRepo, zapLogger, salt)

	userEntity, _ := user.NewUser(ipAddr, salt)
	userDTO := user.MapToDTO(userEntity)

	tests := []struct {
		name     string
		ctx      context.Context
		hashedId string
		setup    func(context.Context, string)
		expect   func(*testing.T, *user.DTO, error)
	}{
		{
			name:     "should return user",
			ctx:      context.Background(),
			hashedId: userDTO.IpAddress,
			setup: func(ctx context.Context, hashedIp string) {
				mockRepo.EXPECT().GetUser(ctx, bson.M{"ip_address": hashedIp}).Return(userEntity, nil)
			},
			expect: func(t *testing.T, dto *user.DTO, err error) {
				assert.NotNil(t, dto)
				assert.Nil(t, err)
				assert.Equal(t, dto.ID, userDTO.ID)
			},
		},
		{
			name:     "should return not found",
			ctx:      context.Background(),
			hashedId: "incorrect",
			setup: func(ctx context.Context, hashedIp string) {
				mockRepo.EXPECT().GetUser(ctx, bson.M{"ip_address": hashedIp}).Return(nil, user.ErrNotFound)
			},
			expect: func(t *testing.T, dto *user.DTO, err error) {
				assert.Nil(t, dto)
				assert.Equal(t, user.ErrNotFound, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(tc.ctx, tc.hashedId)
			s, err := service.GetUserByIp(tc.ctx, tc.hashedId)
			tc.expect(t, s, err)
		})
	}
}

func TestService_CreateUser(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockRepo := mock_user.NewMockRepository(controller)
	salt := "salt"
	ipAddr := "127.0.0.1"

	newLogger, _ := logger.NewLogger("development")
	zapLogger, _ := newLogger.SetupZapLogger()

	service, _ := user.NewService(mockRepo, zapLogger, salt)

	userEntity, _ := user.NewUser(ipAddr, salt)

	var emptyStr string

	tests := []struct {
		name   string
		ctx    context.Context
		ipAddr string
		setup  func(context.Context, string)
		expect func(*testing.T, *user.DTO, error)
	}{
		{
			name:   "should return user",
			ctx:    context.Background(),
			ipAddr: ipAddr,
			setup: func(ctx context.Context, ip string) {
				mockRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(userEntity.ID.Hex(), nil)
			},
			expect: func(t *testing.T, dto *user.DTO, err error) {
				assert.NotNil(t, dto)
				assert.Nil(t, err)
			},
		},
		{
			name:   "should return user",
			ctx:    context.Background(),
			ipAddr: ipAddr,
			setup: func(ctx context.Context, ip string) {
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
			tc.setup(tc.ctx, tc.ipAddr)
			s, err := service.CreateUser(tc.ctx, tc.ipAddr)
			tc.expect(t, s, err)
		})
	}
}
