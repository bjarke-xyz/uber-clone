package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bjarke-xyz/uber-clone-backend/internal/domain"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation"
)

func (a *api) getVehiclesHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := TokenFromContext(r.Context())
	user, err := a.userRepo.GetByUserID(r.Context(), token.Subject)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	vehicles, err := a.vehicleRepo.GetByOwnerId(r.Context(), user.ID)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	err = a.enrichWithPositions(r.Context(), vehicles)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	a.respond(w, r, vehicles)
}

type GetSimulatedVehiclesResponse struct {
	Vehicles  []domain.Vehicle         `json:"vehicles"`
	Positions []domain.VehiclePosition `json:"positions"`
}

func (a *api) enrichWithPositions(ctx context.Context, vehicles []domain.Vehicle) error {
	vehicleIds := make([]int64, len(vehicles))
	for i, v := range vehicles {
		vehicleIds[i] = v.ID
	}
	vehiclePositions, err := a.vehicleRepo.GetVehiclePositions(ctx, vehicleIds)
	if err != nil {
		return err
	}
	vehiclePosMap := make(map[int64]domain.VehiclePosition)
	for _, p := range vehiclePositions {
		vehiclePosMap[p.VehicleID] = p
	}
	for i, vehicle := range vehicles {
		pos, ok := vehiclePosMap[vehicle.ID]
		if ok {
			vehicles[i].LastRecordedPosition = &pos
		}
	}
	return nil
}

func (a *api) handleGetSimulatedVehicles(w http.ResponseWriter, r *http.Request) {
	vehicles, err := a.vehicleRepo.GetSimulatedVehicles(r.Context())
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	err = a.enrichWithPositions(r.Context(), vehicles)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	a.respond(w, r, vehicles)
}

func (a *api) updateVehiclePositionHandler(w http.ResponseWriter, r *http.Request) {
	token, _ := TokenFromContext(r.Context())
	input := &UpdateVehiclePositionInput{}
	if err := json.NewDecoder(r.Body).Decode(input); err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, fmt.Errorf("failed to decode body: %w", err))
		return
	}
	err := input.Validate()
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	user, err := a.userRepo.GetByUserID(r.Context(), token.Subject)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, fmt.Errorf("failed to get user: %w", err))
		return
	}
	vehicleIdStr := chi.URLParam(r, "vehicleID")
	vehicleId, err := strconv.Atoi(vehicleIdStr)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, fmt.Errorf("failed to parse vehicle id: %w", err))
		return
	}
	vehicle, err := a.vehicleRepo.GetByIdAndOwnerId(r.Context(), int64(vehicleId), user.ID)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, fmt.Errorf("failed to get vehicle: %w", err))
		return
	}
	err = a.vehicleRepo.UpdatePosition(r.Context(), vehicle.ID, input.Lat, input.Lng)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, fmt.Errorf("failed to update position: %w", err))
		return
	}

	a.emitPositionUpdateEvent(domain.VehiclePosition{
		VehicleID:  vehicle.ID,
		Lat:        input.Lat,
		Lng:        input.Lng,
		RecordedAt: time.Now(),
		Bearing:    input.Bearing,
		Speed:      input.Speed,
	})

	w.WriteHeader(http.StatusAccepted)
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

type UpdateVehiclePositionInput struct {
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Bearing float32 `json:"bearing"`
	Speed   float32 `json:"speed"`
}

func (a *api) emitPositionUpdateEvent(vehiclePos domain.VehiclePosition) error {
	sseStr, err := formatServerSentEvent("position-update", vehiclePos)
	if err != nil {
		return err
	}
	bytes := []byte(sseStr)
	a.broker.Notifier <- bytes
	return nil
}

func (i *UpdateVehiclePositionInput) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.Lat, validation.Required),
		validation.Field(&i.Lng, validation.Required),
	)
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
