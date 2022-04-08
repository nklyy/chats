package room

type DTO struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Messages *[]*RoomMessage `json:"messages"`
}
