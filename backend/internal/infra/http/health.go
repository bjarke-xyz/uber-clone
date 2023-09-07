package http

import (
	"context"
	"net/http"
)

func (a *api) healthCheckHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	data := map[string]string{
		"status": "OK",
	}
	return a.respond(w, r, data)
}
