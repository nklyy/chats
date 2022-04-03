package old_chat

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
)

//const welcomeMessage = "%s joined the room"

type Room struct {
	Name      string `json:"name"`
	clients   map[*Client]bool
	broadcast chan *Message
}

// NewRoom creates a new Room
func NewRoom(name string) *Room {
	return &Room{
		Name:      name,
		clients:   make(map[*Client]bool),
		broadcast: make(chan *Message),
	}
}

// RunRoom runs our room, accepting various requests
func (room *Room) RunRoom(redis *redis.Client) {
	for {
		select {
		case message := <-room.broadcast:
			j, err := json.Marshal(message.Message)
			if err != nil {
				log.Println(err)
			}
			room.publishRoomMessage(redis, j, message.TargetRoom)
		}
	}
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.Send <- message
	}
}

func (room *Room) subscribeToRoomMessages(redis *redis.Client, roomName string) {
	pubsub := redis.Subscribe(context.Background(), roomName)

	ch := pubsub.Channel()

	for msg := range ch {
		room.broadcastToClientsInRoom([]byte(msg.Payload))
	}
}

func (room *Room) publishRoomMessage(redis *redis.Client, message []byte, roomName string) {
	err := redis.Publish(context.Background(), roomName, message).Err()

	if err != nil {
		log.Println(err)
	}
}

func (room *Room) GetName() string {
	return room.Name
}
