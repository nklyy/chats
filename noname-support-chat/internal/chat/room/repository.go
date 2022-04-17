package room

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go
type Repository interface {
	GetRoom(ctx context.Context, filters bson.M) (*Model, error)
	CreateRoom(ctx context.Context, room *Model) (string, error)
	UpdateRoom(ctx context.Context, model *Model) error
	DeleteRoom(ctx context.Context, name string) error
}

type repository struct {
	db     *mongo.Client
	dbName string
	logger *zap.SugaredLogger
}

func NewRepository(db *mongo.Client, dbName string, logger *zap.SugaredLogger) (Repository, error) {
	if db == nil {
		return nil, errors.New("invalid rooms database")
	}
	if dbName == "" {
		return nil, errors.New("invalid database name")
	}
	if logger == nil {
		return nil, errors.New("invalid logger")
	}

	return &repository{db: db, dbName: dbName, logger: logger}, nil
}

func (r *repository) GetRoom(ctx context.Context, filters bson.M) (*Model, error) {
	var room Model

	if err := r.db.Database(r.dbName).Collection("rooms").FindOne(ctx, filters).Decode(&room); err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Errorf("unable to find room by name %v", err)
			return nil, ErrNotFound
		}

		r.logger.Errorf("unable to find room due to internal error: %v", err)
		return nil, err
	}

	return &room, nil
}

func (r *repository) CreateRoom(ctx context.Context, room *Model) (string, error) {
	//mod := mongo.IndexModel{
	//	Keys:    bson.M{"email": 1}, // index in ascending order or -1 for descending order
	//	Options: options.Index().SetUnique(true),
	//}
	//
	//_, err := r.db.Database(r.dbName).Collection("rooms").Indexes().CreateOne(ctx, mod)
	//if err != nil {
	//	r.logger.Errorf("failed to create room index: %v", err)
	//	return "", err
	//}

	_, err := r.db.Database(r.dbName).Collection("rooms").InsertOne(ctx, room)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			r.logger.Errorf("failed to insert room data to db due to duplicate error: %v", err)
			return "", ErrAlreadyExists
		}

		r.logger.Errorf("failed to insert room data to db: %v", err)
		return "", ErrFailedSaveRoom
	}

	return room.ID.Hex(), nil
}

func (r *repository) UpdateRoom(ctx context.Context, model *Model) error {
	_, err := r.db.Database(r.dbName).Collection("rooms").UpdateOne(ctx, bson.M{"name": model.Name},
		bson.D{primitive.E{Key: "$set", Value: model}})

	if err != nil {
		r.logger.Errorf("failed to update room %v", err)
		return ErrFailedUpdateRoom
	}

	return nil
}

func (r *repository) DeleteRoom(ctx context.Context, name string) error {
	_, err := r.db.Database(r.dbName).Collection("rooms").DeleteOne(ctx, bson.M{"name": name})
	if err != nil {
		r.logger.Errorf("failed to delete room %v", err)
		return ErrFailedDeleteRoom
	}

	return nil
}
