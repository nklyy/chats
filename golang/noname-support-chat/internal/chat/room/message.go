package room

import "time"

type Message struct {
	Action  string           `json:"action"`
	Message EncryptedMessage `json:"message,omitempty"`
	Token   string           `json:"token"`
}

type EncryptedMessage struct {
	Data string `json:"data" bson:"data"`
	Salt string `json:"salt" bson:"salt"`
	Iv   string `json:"iv" bson:"iv"`
}

type MessageResponse struct {
	Action  string            `json:"action"`
	Message *EncryptedMessage `json:"message,omitempty"`
	From    string            `json:"from"`
	Error   interface{}       `json:"error"`
}

type BroadcastMessage struct {
	Action   string          `json:"action"`
	Message  MessageResponse `json:"message"`
	RoomName string          `json:"roomName"`
}

type FormatMessages struct {
	To      string           `json:"to,omitempty"`
	From    string           `json:"from,omitempty"`
	Message EncryptedMessage `json:"message"`
	Time    time.Time        `json:"time"`
}

type RoomMessage struct {
	Id      string           `bson:"id"`
	Time    time.Time        `bson:"time"`
	Message EncryptedMessage `bson:"message,omitempty"`
}
