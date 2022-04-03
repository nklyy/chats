package old_chat

import (
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DTO struct {
}

type Client struct {
	Id         string          `json:"id"`
	Free       bool            `json:"free,omitempty"`
	Support    bool            `json:"user,omitempty"`
	Room       *Room           `json:"room"`
	Connection *websocket.Conn `json:"connection"`
	Send       chan []byte     `json:"send"`
}

type Message struct {
	Action     string      `json:"action"`
	Message    interface{} `json:"message,omitempty"`
	User       string      `json:"user"`
	TargetRoom string      `json:"target_room,omitempty"`
}

type RoomModel struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	UserId    string             `bson:"user_id"`
	SupportId string             `bson:"support_id"`
}
