package api

import (
	"context"
	"net/http"
	"os"

	"github.com/bjarke-xyz/uber-clone-backend/internal/auth"
)

func (a *api) firebaseJwtVerifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idTokenStr := auth.TokenFromHeader(r)
		if idTokenStr == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		token, err := auth.ValidateToken(ctx, os.Getenv("FIREBASE_PROJECT_ID"), idTokenStr)
		if err != nil {
			a.errorResponse(w, r, http.StatusUnauthorized, err)
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

func NewContext(ctx context.Context, t *auth.AuthToken, err error) context.Context {
	ctx = context.WithValue(ctx, TokenCtxKey, t)
	ctx = context.WithValue(ctx, ErrorCtxKey, err)
	return ctx
}

func TokenFromContext(ctx context.Context) (*auth.AuthToken, error) {
	token, _ := ctx.Value(TokenCtxKey).(*auth.AuthToken)
	var err error
	err, _ = ctx.Value(ErrorCtxKey).(error)
	return token, err
}
