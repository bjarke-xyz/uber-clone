package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bjarke-xyz/uber-clone-backend/internal/core/rides"
)

const orsBaseUrl string = "https://api.openrouteservice.org"

type OpenRouteServiceClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewOpenRouteServiceClient(apiKey string) rides.RouteServiceClient {
	httpClient := &http.Client{
		Timeout: time.Minute,
	}
	return &OpenRouteServiceClient{
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (o *OpenRouteServiceClient) GetDirections(locations [][]float64) (*rides.ORSDirections, error) {
	reqBody := make(map[string]any, 0)
	reqBody["coordinates"] = locations
	reqBody["maneuvers"] = true
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	reqBodyReader := bytes.NewReader(reqBodyBytes)
	req, err := http.NewRequest("POST", orsBaseUrl+"/v2/directions/driving-car", reqBodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", o.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		respStr := string(respBytes)
		return nil, fmt.Errorf("got error response from OSR: status=%v body=%v", resp.StatusCode, respStr)
	}
	directions := &rides.ORSDirections{}
	err = json.Unmarshal(respBytes, directions)
	if err != nil {
		return nil, err
	}
	return directions, nil
}
