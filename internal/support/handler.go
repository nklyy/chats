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

func (h *Handler) SetupAuthRoutes(router chi.Router) {
	router.Post("/support/registration", h.Registration)
	router.Post("/support/login", h.Login)
	router.Post("/support/refresh", h.Refresh)
	router.Post("/support/logout", h.Logout)
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

	respond.Respond(w, http.StatusCreated, RegistrationResponseDTO{
		Id: *supportId,
	})
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

	accessToken, refreshToken, err := h.supportSvc.Login(r.Context(), &dto)
	if err != nil {
		respond.Respond(w, errors.HTTPCode(err), err)
		return
	}

	respond.Respond(w, http.StatusOK, &LoginResponseDTO{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	})
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var dto RefreshDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respond.Respond(w, errors.HTTPCode(err), errors.NewInternal(err.Error()))
		return
	}

	if err := Validate(dto); err != nil {
		respond.Respond(w, errors.HTTPCode(err), err)
		return
	}

	accessToken, refreshToken, err := h.supportSvc.Refresh(r.Context(), &dto)
	if err != nil {
		respond.Respond(w, errors.HTTPCode(err), err)
		return
	}

	respond.Respond(w, http.StatusOK, LoginResponseDTO{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var dto LogoutDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respond.Respond(w, errors.HTTPCode(err), errors.NewInternal(err.Error()))
		return
	}

	if err := Validate(dto); err != nil {
		respond.Respond(w, errors.HTTPCode(err), err)
		return
	}

	err := h.supportSvc.Logout(r.Context(), &dto)
	if err != nil {
		respond.Respond(w, errors.HTTPCode(err), err)
		return
	}

	respond.Respond(w, http.StatusOK, "OK")
}
