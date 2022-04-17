package user

import (
	goErr "errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"noname-support-chat/pkg/errors"
	"noname-support-chat/pkg/respond"
)

type Handler struct {
	userSvc Service
}

func NewHandler(userSvc Service) (*Handler, error) {
	if userSvc == nil {
		return nil, goErr.New("invalid user service")
	}

	return &Handler{userSvc: userSvc}, nil
}

func (h *Handler) SetupRoutes(router chi.Router) {
	router.Get("/user/{id}", h.GetUserById)
	router.Get("/free-user", h.GetFreeUser)
}

func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.userSvc.GetUserById(r.Context(), id, false)
	if err != nil {
		respond.Respond(w, errors.HTTPCode(err), err)
		return
	}

	respond.Respond(w, http.StatusOK, user)
}

func (h *Handler) GetFreeUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.userSvc.GetFreeUser(r.Context())
	if err != nil {
		respond.Respond(w, errors.HTTPCode(err), err)
		return
	}

	respond.Respond(w, http.StatusOK, user)
}
