package room

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Messages *[]*RoomMessage    `bson:"messages"`
}
