package auth

import (
	"context"
	"time"

	"github.com/tquach/prestige-api/platform/logger"

	"net/http"
)

// DefaultHandler wraps a middleware handler.
type DefaultHandler struct {
	tokenSvc TokenService
	logger   logger.Logger
}

// Handler returns a new middleware handler that verifies a token for validity.
func (d *DefaultHandler) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "OPTIONS" {
			next.ServeHTTP(w, req)
			return
		}

		token, err := d.tokenSvc.VerifyToken(req)
		if err != nil {
			d.logger.Errorf("Failed to verify token: %s (%s)", err.Error(), req.RemoteAddr)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if token.Valid {
			d.logger.Debugf("Valid token %s", req.RemoteAddr)
			d.logger.Debugf("token header: %s, claims: %v, sign: %s", token.Header, token.Claims, token.Signature)

			ts := token.Claims["exp"].(float64)
			exp := time.Unix(int64(ts), 0)
			if err != nil {
				d.logger.Warningf("Invalid claim for exp %s", token.Claims["exp"])
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			d.logger.Debugf("Token valid and expires at %s", exp)
			if time.Now().After(exp) {
				d.logger.Warningf("Token has expired %s", token.Claims["exp"])
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			uid := int(token.Claims["uid"].(float64))
			d.logger.Debugf("Found userID in JWT token: %d", uid)
			ctx := context.WithValue(req.Context(), "userID", uid)
			next.ServeHTTP(w, req.WithContext(ctx))
		} else {
			d.logger.Warningf("Invalid token from %s", req.RemoteAddr)
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}

// SecuredHandler constructs a DefaultHandler instance and initializes required properties.
func SecuredHandler(tokenSvc TokenService) *DefaultHandler {
	return &DefaultHandler{
		tokenSvc: tokenSvc,
		logger:   logger.New("auth"),
	}
}
