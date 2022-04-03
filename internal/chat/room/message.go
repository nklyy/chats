package room

type Message struct {
	Action     string      `json:"action"`
	Message    interface{} `json:"message,omitempty"`
	User       string      `json:"user"`
	TargetRoom string      `json:"target_room,omitempty"`
}

type MessageResponse struct {
	Action  string      `json:"action"`
	Message interface{} `json:"message"`
	UserId  string      `json:"user_id,omitempty"`
	RoomID  string      `json:"room_id,omitempty"`
	Error   interface{} `json:"error"`
}
