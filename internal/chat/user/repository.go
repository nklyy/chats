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
	GetUser(ctx context.Context, filters bson.M) (*User, error)
	GetUsers(ctx context.Context, filters bson.M) ([]*User, error)
	CreateUser(ctx context.Context, user *User) (string, error)
	UpdateUser(ctx context.Context, user *User) error
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

func (r *repository) GetUser(ctx context.Context, filters bson.M) (*User, error) {
	var user User

	if err := r.db.Database(r.dbName).Collection("users").FindOne(ctx, filters).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Errorf("unable to find user by fingerprint : %v", err)
			return nil, ErrNotFound
		}

		r.logger.Errorf("unable to find user due to internal error: %v", err)
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetUsers(ctx context.Context, filters bson.M) ([]*User, error) {
	var users []*User

	cursor, err := r.db.Database(r.dbName).Collection("users").Find(ctx, filters)
	if err != nil {
		r.logger.Errorf("failed to get users: %v", err)
		return nil, ErrFailedFindFreeUsers
	}

	if err = cursor.All(ctx, &users); err != nil {
		r.logger.Errorf("failed to get users: %v", err)
		return nil, ErrFailedFindFreeUsers
	}

	return users, nil
}

func (r *repository) CreateUser(ctx context.Context, user *User) (string, error) {
	mod := mongo.IndexModel{
		Keys:    bson.M{"fingerprint": 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(true),
	}

	_, err := r.db.Database(r.dbName).Collection("users").Indexes().CreateOne(ctx, mod)
	if err != nil {
		r.logger.Errorf("failed to create user index: %v", err)
		return "", err
	}

	_, err = r.db.Database(r.dbName).Collection("users").InsertOne(ctx, user)
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
	_, err := r.db.Database(r.dbName).Collection("users").UpdateOne(ctx, bson.M{"fingerprint": user.Fingerprint},
		bson.D{primitive.E{Key: "$set", Value: user}})

	if err != nil {
		r.logger.Errorf("failed to update user %v", err)
		return ErrFailedUpdateUser
	}

	return nil
}
