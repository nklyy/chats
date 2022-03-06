package support

import (
	gerrors "errors"
	"go.uber.org/zap"
	"net/http"
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
	jwtSvc jwt.Service
	logger *zap.SugaredLogger
}

func NewMiddleware(jwtSvc jwt.Service, logger *zap.SugaredLogger) (Middleware, error) {
	if jwtSvc == nil {
		return nil, gerrors.New("invalid jwt service")
	}
	if logger == nil {
		return nil, gerrors.New("invalid logger")
	}

	return &middleware{jwtSvc: jwtSvc, logger: logger}, nil
}

func (m *middleware) JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		if len(authorization) == 0 {
			m.logger.Error("failed to get auth token")
			respond.Respond(w, errors.HTTPCode(ErrRequiredAuthorizationToken), ErrRequiredAuthorizationToken)
			return
		}

		authorizationParts := strings.Split(authorization, " ")

		if len(authorizationParts) != 2 || len(authorizationParts[1]) == 0 || authorizationParts[0] != "Bearer" {
			m.logger.Error("invalid auth token")
			respond.Respond(w, errors.HTTPCode(ErrAuthorizationToken), ErrAuthorizationToken)
			return
		}

		payload, err := m.jwtSvc.VerifyJWT(authorizationParts[1])
		if err != nil {
			m.logger.Errorf("failed to verify auth token: %v", err)
			respond.Respond(w, errors.HTTPCode(ErrFailedVerifyAuthorizationToken), ErrFailedVerifyAuthorizationToken)
			return
		}

		if payload.Role != "support" {
			m.logger.Error("token doesn't have permission")
			respond.Respond(w, errors.HTTPCode(ErrTokenDoesntHavePermission), ErrTokenDoesntHavePermission)
			return
		}

		next.ServeHTTP(w, r)
	})
}
