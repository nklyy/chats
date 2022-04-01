package user

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
	GetUserById(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, id string) (*User, error)
	CreateUser(ctx context.Context, user *User) (string, error)
}

type repository struct {
	db     *mongo.Client
	dbName string
	logger *zap.SugaredLogger
}

func NewRepository(db *mongo.Client, dbName string, logger *zap.SugaredLogger) (Repository, error) {
	if db == nil {
		return nil, errors.New("invalid user database")
	}
	if dbName == "" {
		return nil, errors.New("invalid database name")
	}
	if logger == nil {
		return nil, errors.New("invalid logger")
	}

	return &repository{db: db, dbName: dbName, logger: logger}, nil
}

func (r *repository) GetUserById(ctx context.Context, id string) (*User, error) {
	var user User

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("failed to decode id %v", err)
		return nil, ErrNotFound
	}

	if err := r.db.Database(r.dbName).Collection("user").FindOne(ctx, bson.M{"_id": objId}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Errorf("unable to find user by id '%s': %v", id, err)
			return nil, ErrNotFound
		}

		r.logger.Errorf("unable to find user due to internal error: %v; id: %s", err, id)
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	if err := r.db.Database(r.dbName).Collection("user").FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Errorf("unable to find user by email '%s': %v", email, err)
			return nil, ErrNotFound
		}

		r.logger.Errorf("unable to find user due to internal error: %v; id: %s", err, email)
		return nil, err
	}

	return &user, nil
}

func (r *repository) CreateUser(ctx context.Context, user *User) (string, error) {
	mod := mongo.IndexModel{
		Keys:    bson.M{"email": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(true),
	}

	_, err := r.db.Database(r.dbName).Collection("user").Indexes().CreateOne(ctx, mod)
	if err != nil {
		r.logger.Errorf("failed to create user index: %v", err)
		return "", err
	}

	_, err = r.db.Database(r.dbName).Collection("user").InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			r.logger.Errorf("failed to insert user data to db due to duplicate error: %v", err)
			return "", ErrAlreadyExists
		}

		r.logger.Errorf("failed to insert user data to db: %v", err)
		return "", ErrFailedSaveUser
	}

	return user.ID.Hex(), nil
}

func (r *repository) UpdateUser(ctx context.Context, user *User) error {
	_, err := r.db.Database(r.dbName).Collection("user").UpdateOne(ctx, bson.M{"email": user.Email},
		bson.D{primitive.E{Key: "$set", Value: user}})

	if err != nil {
		r.logger.Errorf("failed to update user %v", err)
		return ErrFailedUpdateUser
	}

	return nil
}
