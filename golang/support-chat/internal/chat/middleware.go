package chat

import (
	"context"
	gerrors "errors"
	"net/http"
	"support-chat/internal/user"
	"support-chat/pkg/errors"
	"support-chat/pkg/jwt"
	"support-chat/pkg/respond"

	"go.uber.org/zap"
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
		return nil, gerrors.New("[chat_middleware] invalid jwt service")
	}
	if userSvc == nil {
		return nil, gerrors.New("[chat_middleware] invalid user service")
	}
	if logger == nil {
		return nil, gerrors.New("[chat_middleware] invalid logger")
	}

	return &middleware{jwtSvc: jwtSvc, userSvc: userSvc, logger: logger}, nil
}

type contextKey string

func (m *middleware) JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")

		if len(token) == 0 {
			m.logger.Error("failed to get auth token")
			respond.Respond(w, errors.HTTPCode(ErrRequiredToken), ErrRequiredToken)
			return
		}

		payload, err := m.jwtSvc.ParseToken(token, true)
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

		u, err := m.userSvc.GetUserById(r.Context(), payload.Id, true)
		if err != nil {
			m.logger.Errorf("failed to get user: %v", err)
			respond.Respond(w, errors.HTTPCode(err), err)
			return
		}

		err = m.jwtSvc.ExtendExpire(r.Context(), payload)
		if err != nil {
			m.logger.Errorf("failed to extend expire token: %v", err)
			respond.Respond(w, errors.HTTPCode(err), err)
			return
		}

		ctx := context.WithValue(r.Context(), contextKey("user"), *u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
