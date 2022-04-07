package room

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
