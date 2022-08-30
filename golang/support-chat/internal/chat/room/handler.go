package room

import (
	gerrors "errors"
	"net/http"
	"support-chat/internal/user"
	"support-chat/pkg/errors"
	"support-chat/pkg/respond"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	roomSvc Service
}

func NewHandler(roomSvc Service) (*Handler, error) {
	if roomSvc == nil {
		return nil, gerrors.New("[chat_room_handler] invalid room service")
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
