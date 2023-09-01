package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func (a *api) handleGetMyUser(w http.ResponseWriter, r *http.Request) {
	token, _ := TokenFromContext(r.Context())
	user, err := a.userRepo.GetByUserID(r.Context(), token.Subject)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	a.respond(w, r, user)
}

func (a *api) handleGetSimUsers(w http.ResponseWriter, r *http.Request) {
	users, err := a.userRepo.GetSimulatedUsers(r.Context())
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	a.respond(w, r, users)
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
	go a.emitUserLogEvent(userLogEvent)
	go storeUserLog(userLogEvent)
	w.WriteHeader(http.StatusAccepted)
}

func (a *api) handleGetRecentLogs(w http.ResponseWriter, r *http.Request) {
	a.respond(w, r, recentUserLogs)
}

type UserLogEvent struct {
	UserID    int64     `json:"userId"`
	Message   string    `json:"message"`
	Tag       string    `json:"tag"`
	Timestamp time.Time `json:"timestamp"`
}

var recentUserLogs = make([]UserLogEvent, 0)
var recentUserLogsLock sync.RWMutex
var maxRecentLogs = 100

func storeUserLog(event UserLogEvent) {
	recentUserLogsLock.Lock()
	defer recentUserLogsLock.Unlock()
	recentUserLogs = append([]UserLogEvent{event}, recentUserLogs...)
	if len(recentUserLogs) > maxRecentLogs {
		recentUserLogs = recentUserLogs[:maxRecentLogs]
	}
}

func (a *api) emitUserLogEvent(userLogEvent UserLogEvent) {
	sseStr, err := formatServerSentEvent("user-log", userLogEvent)
	if err != nil {
		a.logger.Error("error formatting sse event", "error", err)
		return
	}
	bytes := []byte(sseStr)
	a.broker.Notifier <- bytes
}
