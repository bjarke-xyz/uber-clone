package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

func (a *api) handleGetMyUser(w http.ResponseWriter, r *http.Request) {
	token, _ := TokenFromContext(r.Context())
	user, err := a.userRepo.GetByUserID(r.Context(), token.Subject)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	render.Respond(w, r, user)
}

func (a *api) handleGetSimUsers(w http.ResponseWriter, r *http.Request) {
	users, err := a.userRepo.GetSimulatedUsers(r.Context())
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	render.Respond(w, r, users)
}

type PostUserLogInput struct {
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

func (a *api) handlePostUserLog(w http.ResponseWriter, r *http.Request) {
	token, _ := TokenFromContext(r.Context())
	input := &PostUserLogInput{}
	if err := json.NewDecoder(r.Body).Decode(input); err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, fmt.Errorf("failed to decode body: %w", err))
		return
	}
	user, err := a.userRepo.GetByUserID(r.Context(), token.Subject)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	userLogEvent := UserLogEvent{
		UserID:    user.ID,
		Tag:       input.Tag,
		Message:   input.Message,
		Timestamp: time.Now().UTC(),
	}
	a.emitUserLogEvent(userLogEvent)
	w.WriteHeader(http.StatusAccepted)
}

type UserLogEvent struct {
	UserID    int64     `json:"userId"`
	Message   string    `json:"message"`
	Tag       string    `json:"tag"`
	Timestamp time.Time `json:"timestamp"`
}

func (a *api) emitUserLogEvent(userLogEvent UserLogEvent) error {
	sseStr, err := formatServerSentEvent("user-log", userLogEvent)
	if err != nil {
		return err
	}
	bytes := []byte(sseStr)
	a.broker.Notifier <- bytes
	return nil
}
