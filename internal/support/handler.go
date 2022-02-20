package support

import (
	goErr "errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"noname-realtime-support-chat/pkg/errors"
	"noname-realtime-support-chat/pkg/respond"
)

type Handler struct {
	supportSvc Service
}

func NewHandler(supportSvc Service) (*Handler, error) {
	if supportSvc == nil {
		return nil, goErr.New("invalid support service")
	}

	return &Handler{supportSvc: supportSvc}, nil
}

func (h *Handler) SetupRoutes(router chi.Router) {
	router.Get("/support/{id}", h.GetSupportById)
}

func (h *Handler) GetSupportById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	support, err := h.supportSvc.GetSupportById(r.Context(), id)
	if err != nil {
		respond.Respond(w, errors.HTTPCode(err), err)
		return
	}

	respond.Respond(w, http.StatusOK, support)
}
