package chat

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Handler struct {
}

func NewHandler() (*Handler, error) {

	return &Handler{}, nil
}

func (h *Handler) SetupRoutes(router chi.Router) {
	router.HandleFunc("/chat", h.Chat)
}

func (h *Handler) Chat(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ws.RemoteAddr())
}

//curl  --include \
//--no-buffer \
//--header "Connection: Upgrade" \
//--header "Upgrade: websocket" \
//--header "Host: localhost:5000" \
//--header "Origin: http://localhost:5000" \
//--header "Sec-WebSocket-Key: SGVsbG8sIHdvcmxkIQ==" \
//--header "Sec-WebSocket-Version: 13" \
//http://localhost:5000/api/v1/chat
