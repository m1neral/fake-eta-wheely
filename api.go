package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type FetchCarsPositionsResponse struct {
	Positions []Position
}

type FetchEtasResponse struct {
	Etas []Eta
}

type ApiRequester interface {
	FetchCarPositions(position Position, limit int) (FetchCarsPositionsResponse, error)
	FetchEtas(position Position, carsPositions []Position) (FetchEtasResponse, error)
}

type ApiRequestService struct {
	httpTimeout int
}

const (
	defaultHttpTimeout = 1 // in seconds

	carsPositionsEndpointURL = "https://dev-api.wheely.com/fake-eta/cars"
	etasEndpointURL          = "https://dev-api.wheely.com/fake-eta/predict"
)

func NewApiRequestService() *ApiRequestService {
	return &ApiRequestService{httpTimeout: defaultHttpTimeout}
}

func (s *ApiRequestService) FetchCarPositions(position Position, limit int) (FetchCarsPositionsResponse, error) {
	positions := FetchCarsPositionsResponse{}

	url, err := url.Parse(carsPositionsEndpointURL)
	if err != nil {
		return positions, err
	}

	query := url.Query()
	query.Set("limit", strconv.Itoa(limit))
	query.Set("lat", fmt.Sprintf("%f", position.Lat))
	query.Set("lng", fmt.Sprintf("%f", position.Lng))

	url.RawQuery = query.Encode()

	resp, err := http.Get(url.String())
	if err != nil {
		return positions, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return positions, err
	}

	err = json.Unmarshal(body, &positions.Positions)
	if err != nil {
		return positions, err
	}

	return positions, nil
}

type fetchEtasRequest struct {
	Target Position   `json:"target"`
	Source []Position `json:"source"`
}

func (s *ApiRequestService) FetchEtas(position Position, carsPositions []Position) (FetchEtasResponse, error) {
	etas := FetchEtasResponse{}
	requestPayload := fetchEtasRequest{Target: position, Source: carsPositions}

	requestPayloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		return etas, err
	}

	resp, err := http.Post(etasEndpointURL, "application/json; charset=utf-8", bytes.NewReader(requestPayloadBytes))
	if err != nil {
		return etas, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return etas, err
	}

	err = json.Unmarshal(body, &etas.Etas)
	if err != nil {
		return etas, err
	}

	return etas, nil
}
