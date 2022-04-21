package chat

import (
	gerrors "errors"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"net/http"
	"noname-support-chat/pkg/errors"
	"noname-support-chat/pkg/respond"
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
		return nil, gerrors.New("[chat_handler] invalid chat service")
	}

	return &Handler{chatSvc: chatSvc}, nil
}

func (h *Handler) SetupRoutes(router chi.Router) {
	router.HandleFunc("/chat", h.Chat)
}

func (h *Handler) Chat(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		respond.Respond(w, http.StatusInternalServerError, errors.NewInternal(err.Error()))
		return
	}

	err = h.chatSvc.Chat(r.Context(), ws)
	if err != nil {
		respond.Respond(w, http.StatusInternalServerError, errors.NewInternal(err.Error()))
		return
	}
}
