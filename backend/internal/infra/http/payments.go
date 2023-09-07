package http

import (
	"context"
	"net/http"
)

func (a *api) handleGetCurrencies(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	currencies := a.paymentsService.GetCurrencies()
	return a.respond(w, r, currencies)
}
