package support

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Repository interface {
	GetSupportById(ctx context.Context, id string) (*Support, error)
}

type repository struct {
	db     *mongo.Client
	logger *zap.SugaredLogger
}

func NewRepository(db *mongo.Client, logger *zap.SugaredLogger) (Repository, error) {
	if db == nil {
		return nil, errors.New("invalid support database")
	}
	if logger == nil {
		return nil, errors.New("invalid logger")
	}

	return &repository{db: db, logger: logger}, nil
}

func (r *repository) GetSupportById(ctx context.Context, id string) (*Support, error) {
	var support Support

	if err := r.db.Database("Chat").Collection("support").FindOne(ctx, bson.M{"_id": id}).Decode(&support); err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Errorf("unable to find support by id '%s': %v", id, err)
			return nil, ErrNotFound
		}

		r.logger.Errorf("unable to find support due to internal error: %v; id: %s", err, id)
		return nil, err
	}

	return &support, nil
}
