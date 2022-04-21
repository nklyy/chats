package auth

import (
	"encoding/json"
	goErr "errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"noname-support-chat/pkg/errors"
	"noname-support-chat/pkg/respond"
)

type Handler struct {
	authSvc Service
}

func NewHandler(authSvc Service) (*Handler, error) {
	if authSvc == nil {
		return nil, goErr.New("[chat_auth_handler] invalid auth service")
	}

	return &Handler{authSvc: authSvc}, nil
}

func (h *Handler) SetupRoutes(router chi.Router) {
	router.Post("/registration", h.Registration)
	router.Post("/login", h.Login)
	router.Post("/refresh", h.Refresh)
	router.Post("/logout", h.Logout)
	router.Post("/check", h.Check)
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

	supportId, err := h.authSvc.Registration(r.Context(), &dto)
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

	accessToken, refreshToken, err := h.authSvc.Login(r.Context(), &dto)
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

	accessToken, refreshToken, err := h.authSvc.Refresh(r.Context(), &dto)
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

	err := h.authSvc.Logout(r.Context(), &dto)
	if err != nil {
		respond.Respond(w, errors.HTTPCode(err), err)
		return
	}

	respond.Respond(w, http.StatusOK, "OK")
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	var dto CheckDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respond.Respond(w, errors.HTTPCode(err), errors.NewInternal(err.Error()))
		return
	}

	if err := Validate(dto); err != nil {
		respond.Respond(w, errors.HTTPCode(err), err)
		return
	}

	check, err := h.authSvc.Check(r.Context(), &dto)
	if err != nil {
		respond.Respond(w, errors.HTTPCode(err), err)
		return
	}

	respond.Respond(w, http.StatusOK, check)
}
