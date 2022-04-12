package room

import (
	gerrors "errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"noname-realtime-support-chat/internal/user"
	"noname-realtime-support-chat/pkg/errors"
	"noname-realtime-support-chat/pkg/respond"
)

type Handler struct {
	roomSvc Service
}

func NewHandler(roomSvc Service) (*Handler, error) {
	if roomSvc == nil {
		return nil, gerrors.New("invalid room service")
	}

	return &Handler{roomSvc: roomSvc}, nil
}

func (h *Handler) SetupRoutes(router chi.Router) {
	router.HandleFunc("/get-room-messages", h.GetRoomMessages)
}

func (h *Handler) GetRoomMessages(w http.ResponseWriter, r *http.Request) {
	userCtxValue := r.Context().Value(contextKey("user"))
	if userCtxValue == nil {
		respond.Respond(w, http.StatusUnauthorized, errors.NewInternal("Not authenticated"))
		return
	}

	u := userCtxValue.(user.DTO)
	room, err := h.roomSvc.GetRoomWithFormatMessages(r.Context(), *u.RoomName, u.ID)
	if err != nil {
		respond.Respond(w, http.StatusInternalServerError, errors.NewInternal(err.Error()))
		return
	}

	respond.Respond(w, http.StatusOK, room)
}
