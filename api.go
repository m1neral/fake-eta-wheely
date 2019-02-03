package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type FetchCarPositionsResponse struct {
	Positions []Position
}

type FetchEtasResponse struct {
	Target Position
	Source []Position
}

type ApiRequester interface {
	FetchCarPositions(position Position, limit int) (FetchCarPositionsResponse, error)
	FetchEtas(position Position, carsPositions []Position) (FetchEtasResponse, error)
}

type ApiRequestService struct {
	httpTimeout int
}

const (
	defaultHttpTimeout = 1 // in seconds
	carPositionsLimit  = 100

	carPositionsEndpointURL = "https://dev-api.wheely.com/fake-eta/cars"
	etasEndpointURL         = "https://dev-api.wheely.com/fake-eta/predict"
)

func NewApiRequestService() *ApiRequestService {
	return &ApiRequestService{httpTimeout: defaultHttpTimeout}
}

func (s *ApiRequestService) FetchCarPositions(position Position, limit int) (FetchCarPositionsResponse, error) {
	positions := FetchCarPositionsResponse{}

	url, err := url.Parse(carPositionsEndpointURL)
	if err != nil {
		return positions, err
	}

	query := url.Query()
	query.Set("limit", strconv.Itoa(carPositionsLimit))
	query.Set("lat", fmt.Sprintf("%f", lat))
	query.Set("lng", fmt.Sprintf("%f", lng))

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

func (s *ApiRequestService) FetchEtas(position Position, carsPositions []Position) (FetchEtasResponse, error) {
	return FetchEtasResponse{}, nil
}
