package room

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"noname-realtime-support-chat/pkg/errors"
)

type Room struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Clients   map[*Client]bool
	Broadcast chan *BroadcastMessage
}

func NewRoom(name string) (*Room, error) {
	if name == "" {
		return nil, errors.WithMessage(ErrInvalidName, "should be not empty")
	}

	return &Room{
		ID:        primitive.NewObjectID(),
		Name:      name,
		Clients:   make(map[*Client]bool),
		Broadcast: make(chan *BroadcastMessage),
	}, nil
}

func (r *Room) RunRoom(redis *redis.Client) {
	go r.subscribeToRoomMessages(redis)

	for {
		//select {
		//case message := <-r.Broadcast:
		//	j, err := json.Marshal(message.Message)
		//	if err != nil {
		//		log.Println(err)
		//	}
		//	r.publishRoomMessage(redis, j, message.RoomName)
		//}
		for message := range r.Broadcast {
			j, err := json.Marshal(message.Message)
			if err != nil {
				log.Printf("failed decode broadcast message %v", err)
			}
			r.publishRoomMessage(redis, j, message.RoomName)
		}
	}
}

func (r *Room) broadcastToClientsInRoom(message []byte) {
	for client := range r.Clients {
		client.Send <- message
	}
}

func (r *Room) subscribeToRoomMessages(redis *redis.Client) {
	pubsub := redis.Subscribe(context.Background(), r.Name)

	ch := pubsub.Channel()

	for msg := range ch {
		r.broadcastToClientsInRoom([]byte(msg.Payload))
	}
}

func (r *Room) publishRoomMessage(redis *redis.Client, message []byte, roomName string) {
	err := redis.Publish(context.Background(), roomName, message).Err()

	if err != nil {
		log.Println(err)
	}
}
