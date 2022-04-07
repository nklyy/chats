package room

import (
	"context"
	gerrors "errors"
	"go.uber.org/zap"
	"net/http"
	"noname-realtime-support-chat/internal/user"
	"noname-realtime-support-chat/pkg/errors"
	"noname-realtime-support-chat/pkg/jwt"
	"noname-realtime-support-chat/pkg/respond"
	"strings"
)

//go:generate mockgen -source=middleware.go -destination=mocks/middleware_mock.go
type Middleware interface {
	JwtMiddleware(next http.Handler) http.Handler
}

type middleware struct {
	jwtSvc  jwt.Service
	userSvc user.Service
	logger  *zap.SugaredLogger
}

func NewMiddleware(jwtSvc jwt.Service, userSvc user.Service, logger *zap.SugaredLogger) (Middleware, error) {
	if jwtSvc == nil {
		return nil, gerrors.New("invalid jwt service")
	}
	if userSvc == nil {
		return nil, gerrors.New("invalid user service")
	}
	if logger == nil {
		return nil, gerrors.New("invalid logger")
	}

	return &middleware{jwtSvc: jwtSvc, userSvc: userSvc, logger: logger}, nil
}

func (m *middleware) JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		if len(authorization) == 0 {
			m.logger.Error("failed to get auth token")
			respond.Respond(w, errors.HTTPCode(ErrRequiredToken), ErrRequiredToken)
			return
		}

		authorizationParts := strings.Split(authorization, " ")

		if len(authorizationParts) != 2 || len(authorizationParts[1]) == 0 || authorizationParts[0] != "Bearer" {
			m.logger.Error("invalid auth token")
			respond.Respond(w, errors.HTTPCode(ErrToken), ErrToken)
			return
		}

		payload, err := m.jwtSvc.ParseToken(authorizationParts[1], true)
		if err != nil {
			m.logger.Errorf("failed to parse auth token: %v", err)
			respond.Respond(w, errors.HTTPCode(err), err)
			return
		}

		err = m.jwtSvc.VerifyToken(r.Context(), payload, true)
		if err != nil {
			m.logger.Errorf("failed to verify auth token: %v", err)
			respond.Respond(w, errors.HTTPCode(err), err)
			return
		}

		//if payload.Role == "user" {
		//	m.logger.Error("token doesn't have permission")
		//	respond.Respond(w, errors.HTTPCode(ErrTokenDoesntHavePermission), ErrTokenDoesntHavePermission)
		//	return
		//}

		u, err := m.userSvc.GetUserById(r.Context(), payload.Id, true)
		if err != nil {
			m.logger.Errorf("failed to extend expire token: %v", err)
			respond.Respond(w, errors.HTTPCode(err), err)
			return
		}

		err = m.jwtSvc.ExtendExpire(r.Context(), payload)
		if err != nil {
			m.logger.Errorf("failed to extend expire token: %v", err)
			respond.Respond(w, errors.HTTPCode(err), err)
			return
		}
		ctx := context.WithValue(r.Context(), "user", *u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
