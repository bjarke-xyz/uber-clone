package ws

import (
	"log/slog"
	"net/http"

	"github.com/olahol/melody"
)

type WsManager struct {
	m      *melody.Melody
	logger slog.Logger
}

func NewWsManager(logger slog.Logger) *WsManager {
	m := melody.New()
	wm := &WsManager{
		m: m,
	}
	m.HandleConnect(wm.handleConnect)
	return wm
}

func (wm *WsManager) HandleRequest(w http.ResponseWriter, r *http.Request) {
	wm.m.HandleRequest(w, r)
}

func (wm *WsManager) handleConnect(s *melody.Session) {
	// TODO:
}
