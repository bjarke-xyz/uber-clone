package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/bjarke-xyz/auth/pkg/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) firebaseJwtVerifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idTokenStr := tokenFromHeader(r)
		if idTokenStr == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		validateTokenRequest := jwt.ValidateTokenRequest{
			Token:    idTokenStr,
			Audience: a.cfg.FirebaseProjectId,
		}
		token, err := jwt.ValidateToken(ctx, validateTokenRequest)
		if err != nil {
			a.logger.Warn("error validating token", "error", err)
			if status.Code(err) == codes.Unavailable {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
			return
		}

		ctx = NewContext(ctx, token, err)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func NewContext(ctx context.Context, t jwt.AuthToken, err error) context.Context {
	ctx = context.WithValue(ctx, TokenCtxKey, t)
	ctx = context.WithValue(ctx, ErrorCtxKey, err)
	return ctx
}

func TokenFromContext(ctx context.Context) (jwt.AuthToken, error) {
	token, _ := ctx.Value(TokenCtxKey).(jwt.AuthToken)
	var err error
	err, _ = ctx.Value(ErrorCtxKey).(error)
	return token, err
}

func tokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return authHeader[7:]
	}
	return ""
}
