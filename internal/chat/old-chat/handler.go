package old_chat

import (
	gerrors "errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"net/http"
	"noname-realtime-support-chat/pkg/errors"
	"noname-realtime-support-chat/pkg/respond"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Handler struct {
	chatSvc Service
}

func NewHandler(chatSvc Service) (*Handler, error) {
	if chatSvc == nil {
		return nil, gerrors.New("invalid chat service")
	}

	return &Handler{chatSvc: chatSvc}, nil
}

func (h *Handler) SetupRoutes(router chi.Router) {
	router.HandleFunc("/chat", h.Chat)
	router.Get("/get-user", h.GetUser)
}

func (h *Handler) Chat(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		respond.Respond(w, http.StatusInternalServerError, errors.NewInternal(err.Error()))
		return
	}

	err = h.chatSvc.Chat(ws, r.URL.Query().Get("token"), r.URL.Query().Get("userId"), r.URL.Query().Get("roomId"))
	if err != nil {
		respond.Respond(w, http.StatusInternalServerError, errors.NewInternal(err.Error()))
		return
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.chatSvc.GetUser()
	if err != nil {
		respond.Respond(w, http.StatusInternalServerError, errors.NewInternal(err.Error()))
		return
	}

	if user == nil {
		respond.Respond(w, http.StatusInternalServerError, errors.NewInternal("No user yet"))
		return
	}
	fmt.Println("GET USER", user)

	respond.Respond(w, http.StatusOK, map[string]string{"userId": user.Id, "roomId": user.Room.Name})
}
