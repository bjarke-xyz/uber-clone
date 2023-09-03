package api

import (
	"net/http"

	"github.com/bjarke-xyz/uber-clone-backend/internal/domain"
)

func (a *api) handleGetCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies := domain.GetCurrencies()
	a.respond(w, r, currencies)
}
