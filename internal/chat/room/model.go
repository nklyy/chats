package room

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type RoomMessage struct {
	Id      string    `bson:"id"`
	Time    time.Time `bson:"time"`
	Message string    `bson:"message"`
}

type Model struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Messages *[]*RoomMessage    `bson:"messages"`
}
