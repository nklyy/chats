package support

import (
	"encoding/json"
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
	router.Post("/support/registration", h.Registration)
	router.Post("/support/login", h.Login)
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

func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	var dto RegistrationDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respond.Respond(w, errors.HTTPCode(err), errors.NewInternal(err.Error()))
		return
	}

	if err := Validate(dto); err != nil {
		respond.Respond(w, errors.HTTPCode(err), err)
		return
	}

	supportId, err := h.supportSvc.Registration(r.Context(), &dto)
	if err != nil {
		respond.Respond(w, errors.HTTPCode(err), err)
		return
	}

	respond.Respond(w, http.StatusCreated, map[string]string{"id": *supportId})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var dto LoginDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respond.Respond(w, errors.HTTPCode(err), errors.NewInternal(err.Error()))
		return
	}

	if err := Validate(dto); err != nil {
		respond.Respond(w, errors.HTTPCode(err), err)
		return
	}

	token, err := h.supportSvc.Login(r.Context(), &dto)
	if err != nil {
		respond.Respond(w, errors.HTTPCode(err), err)
		return
	}

	respond.Respond(w, http.StatusOK, map[string]string{"token": *token})
}
