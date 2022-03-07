package support

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go
type Repository interface {
	GetSupportById(ctx context.Context, id string) (*Support, error)
	GetSupportByEmail(ctx context.Context, id string) (*Support, error)
	CreateSupport(ctx context.Context, support *Support) (string, error)
}

type repository struct {
	db     *mongo.Client
	dbName string
	logger *zap.SugaredLogger
}

func NewRepository(db *mongo.Client, dbName string, logger *zap.SugaredLogger) (Repository, error) {
	if db == nil {
		return nil, errors.New("invalid support database")
	}
	if dbName == "" {
		return nil, errors.New("invalid database name")
	}
	if logger == nil {
		return nil, errors.New("invalid logger")
	}

	return &repository{db: db, dbName: dbName, logger: logger}, nil
}

func (r *repository) GetSupportById(ctx context.Context, id string) (*Support, error) {
	var support Support

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("failed to decode id %v", err)
		return nil, ErrNotFound
	}

	if err := r.db.Database(r.dbName).Collection("support").FindOne(ctx, bson.M{"_id": objId}).Decode(&support); err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Errorf("unable to find support by id '%s': %v", id, err)
			return nil, ErrNotFound
		}

		r.logger.Errorf("unable to find support due to internal error: %v; id: %s", err, id)
		return nil, err
	}

	return &support, nil
}

func (r *repository) GetSupportByEmail(ctx context.Context, email string) (*Support, error) {
	var support Support

	if err := r.db.Database(r.dbName).Collection("support").FindOne(ctx, bson.M{"email": email}).Decode(&support); err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Errorf("unable to find support by email '%s': %v", email, err)
			return nil, ErrNotFound
		}

		r.logger.Errorf("unable to find support due to internal error: %v; id: %s", err, email)
		return nil, err
	}

	return &support, nil
}

func (r *repository) CreateSupport(ctx context.Context, support *Support) (string, error) {
	mod := mongo.IndexModel{
		Keys:    bson.M{"email": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(true),
	}

	_, err := r.db.Database(r.dbName).Collection("support").Indexes().CreateOne(ctx, mod)
	if err != nil {
		r.logger.Errorf("failed to create support index: %v", err)
		return "", err
	}

	_, err = r.db.Database(r.dbName).Collection("support").InsertOne(ctx, support)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			r.logger.Errorf("failed to insert support data to db due to duplicate error: %v", err)
			return "", ErrAlreadyExists
		}

		r.logger.Errorf("failed to insert support data to db: %v", err)
		return "", err
	}

	return support.ID.Hex(), nil
}

func (r *repository) UpdateSupport(ctx context.Context, support *Support) error {
	_, err := r.db.Database(r.dbName).Collection("support").UpdateOne(ctx, bson.M{"email": support.Email},
		bson.D{primitive.E{Key: "$set", Value: support}})

	if err != nil {
		r.logger.Errorf("failed to update support %v", err)
		return err
	}

	return nil
}
