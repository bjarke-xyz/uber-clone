package api

import (
	"net/http"

	"github.com/go-chi/render"
)

func (a *api) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "OK",
	}
	_, err := a.userRepo.GetByID(r.Context(), 1)
	if err != nil {
		a.logger.Error("failed to get user", "error", err)
	}
	render.Respond(w, r, data)
}
