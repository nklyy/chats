package room_test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"noname-support-chat/internal/chat/room"
	"testing"
)

func TestNewRepository(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tests := []struct {
		name   string
		db     *mongo.Client
		dbName string
		logger *zap.SugaredLogger
		expect func(*testing.T, room.Repository, error)
	}{
		{
			name:   "should return repository",
			db:     &mongo.Client{},
			dbName: "Chat",
			logger: &zap.SugaredLogger{},
			expect: func(t *testing.T, r room.Repository, err error) {
				assert.NotNil(t, r)
				assert.Nil(t, err)
			},
		},
		{
			name:   "should return invalid database",
			db:     nil,
			dbName: "Chat",
			logger: &zap.SugaredLogger{},
			expect: func(t *testing.T, r room.Repository, err error) {
				assert.Nil(t, r)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[chat_room_repository] invalid rooms database")
			},
		},
		{
			name:   "should return invalid database name",
			db:     &mongo.Client{},
			dbName: "",
			logger: &zap.SugaredLogger{},
			expect: func(t *testing.T, r room.Repository, err error) {
				assert.Nil(t, r)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[chat_room_repository] invalid database name")
			},
		},
		{
			name:   "should return invalid logger",
			db:     &mongo.Client{},
			dbName: "Chat",
			logger: nil,
			expect: func(t *testing.T, r room.Repository, err error) {
				assert.Nil(t, r)
				assert.NotNil(t, err)
				assert.EqualError(t, err, "[chat_room_repository] invalid logger")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := room.NewRepository(tc.db, tc.dbName, tc.logger)
			tc.expect(t, svc, err)
		})
	}
}
