package room

import "time"

type Message struct {
	Action  string `json:"action"`
	Message string `json:"message,omitempty"`
	Token   string `json:"token"`
}

type MessageResponse struct {
	Action  string      `json:"action"`
	Message string      `json:"message"`
	From    string      `json:"from"`
	Error   interface{} `json:"error"`
}

type BroadcastMessage struct {
	Action   string          `json:"action"`
	Message  MessageResponse `json:"message"`
	RoomName string          `json:"roomName"`
}

type FormatMessages struct {
	To      string    `json:"to,omitempty"`
	From    string    `json:"from,omitempty"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

type RoomMessage struct {
	Id      string    `bson:"id"`
	Time    time.Time `bson:"time"`
	Message string    `bson:"message"`
}
