package api

import (
	"context"
	"net/http"

	"github.com/bjarke-xyz/uber-clone-backend/internal/core/users"
)

func (a *api) handleGetMyUser(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	token, _ := TokenFromContext(r.Context())
	user, err := a.userService.GetUserByID(ctx, token.Subject)
	if err != nil {
		return err
	}
	return a.respond(w, r, user)
}

func (a *api) handleGetSimUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	users, err := a.userService.GetSimulatedUsers(ctx)
	if err != nil {
		return err
	}
	return a.respond(w, r, users)
}

func (a *api) handlePostUserLog(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	token, _ := TokenFromContext(r.Context())
	input := &users.PostUserLogInput{}
	if err := decodeBody(r.Body, input); err != nil {
		return err
	}
	userLogEvent, err := a.userService.AddUserLog(ctx, token.Subject, input)
	if err != nil {
		return err
	}
	// TODO: Move to event listener
	go a.emitUserLogEvent(userLogEvent)
	return a.respondStatus(w, r, http.StatusAccepted, nil)
}

func (a *api) handleGetRecentLogs(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	recentUserLogs, err := a.userService.GetRecentLogs(ctx)
	if err != nil {
		return err
	}
	return a.respond(w, r, recentUserLogs)
}

func (a *api) emitUserLogEvent(userLogEvent users.UserLogEvent) {
	sseStr, err := formatServerSentEvent("user-log", userLogEvent)
	if err != nil {
		a.logger.Error("error formatting sse event", "error", err)
		return
	}
	bytes := []byte(sseStr)
	a.broker.Notifier <- bytes
}
