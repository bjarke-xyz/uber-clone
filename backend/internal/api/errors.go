package api

import (
	"errors"
	"net/http"

	"github.com/bjarke-xyz/uber-clone-backend/internal/domain"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (a *api) errorResponse(w http.ResponseWriter, _ *http.Request, status int, err error) {
	a.logger.Error("error", "error", err)
	if errors.Is(err, domain.ErrNotFound) {
		status = http.StatusNotFound
	}
	http.Error(w, err.Error(), status)
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
