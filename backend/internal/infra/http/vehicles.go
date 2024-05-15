package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bjarke-xyz/uber-clone-backend/internal/core/vehicles"
)

func (a *api) getVehiclesHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	token, _ := TokenFromContext(r.Context())
	vehicleList, err := a.vehicleService.GetVehicles(ctx, token.Subject)
	if err != nil {
		return err
	}
	return a.respond(w, r, vehicleList)
}

type GetSimulatedVehiclesResponse struct {
	Vehicles  []vehicles.Vehicle         `json:"vehicles"`
	Positions []vehicles.VehiclePosition `json:"positions"`
}

func (a *api) handleGetSimulatedVehicles(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vehicleList, err := a.vehicleService.GetSimulatedVehicles(ctx)
	if err != nil {
		return err
	}
	return a.respond(w, r, vehicleList)
}

func (a *api) updateVehiclePositionHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	token, _ := TokenFromContext(r.Context())
	vehicleId, err := urlParamInt(r, "vehicleID")
	if err != nil {
		return err
	}
	input := &vehicles.UpdateVehiclePositionInput{}
	if err := decodeBody(r.Body, input); err != nil {
		return err
	}
	go func() {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().UTC().Add(1*time.Second))
		err := func() error {
			err := a.vehicleService.UpdateVehiclePosition(ctx, token.Subject, vehicleId, input)
			if err != nil {
				cancel()
				return err
			}
			return nil
		}()
		if err != nil {
			a.logger.Error("updateVehiclePositionHandler had error", "error", err)
		}
	}()

	return a.respondStatus(w, r, http.StatusAccepted, nil)
}

func (a *api) handleEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Each connection registers its own message channel with the Broker's connections registry
	messageChan := make(chan []byte)
	a.broker.newClients <- messageChan
	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		a.broker.closingClients <- messageChan
	}()
	go func() {
		<-r.Context().Done()
		a.broker.closingClients <- messageChan
	}()

	for bytes := range messageChan {
		sseStr := string(bytes)
		_, err := fmt.Fprint(w, sseStr)
		if err != nil {
			a.logger.Error("failed to fprintf sse event", "error", err)
			break
		}
		flusher.Flush()
	}
}

func (a *api) pubsubSubscribeVehicle(ctx context.Context) {
	go func() {
		ch := a.pubSub.SubscribeBytes(vehicles.TopicPositionUpdate)
		for {
			select {
			case msg := <-ch:
				event := vehicles.VehiclePosition{}
				err := json.Unmarshal(msg, &event)
				if err != nil {
					a.logger.Error("failed to unmarshal VehiclePosition", "error", err)
				} else {
					a.emitPositionUpdateEvent(event)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (a *api) emitPositionUpdateEvent(vehiclePos vehicles.VehiclePosition) error {
	sseStr, err := formatServerSentEvent(vehicles.TopicPositionUpdate, vehiclePos)
	if err != nil {
		return err
	}
	bytes := []byte(sseStr)
	a.broker.Notifier <- bytes
	return nil
}

func formatServerSentEvent(event string, data any) (string, error) {
	m := map[string]any{
		"data": data,
	}
	buff := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buff)
	err := encoder.Encode(m)
	if err != nil {
		return "", err
	}
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("event: %s\n", event))
	sb.WriteString(fmt.Sprintf("data: %v\n\n", buff.String()))
	return sb.String(), nil
}
