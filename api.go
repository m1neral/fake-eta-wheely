package eta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type ApiRequester interface {
	FetchCarPositions(position Position, limit int) ([]Position, error)
	FetchEtas(position Position, carsPositions []Position) ([]Eta, error)
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

func (s *ApiRequestService) FetchCarPositions(position Position, limit int) ([]Position, error) {
	positions := []Position{}

	parsedUrl, err := url.Parse(carsPositionsEndpointURL)
	if err != nil {
		return positions, err
	}

	query := parsedUrl.Query()
	query.Set("limit", strconv.Itoa(limit))
	query.Set("lat", fmt.Sprintf("%f", position.Lat))
	query.Set("lng", fmt.Sprintf("%f", position.Lng))

	parsedUrl.RawQuery = query.Encode()

	resp, err := http.Get(parsedUrl.String())
	if err != nil {
		return positions, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return positions, err
	}

	err = json.Unmarshal(body, &positions)
	if err != nil {
		return positions, err
	}

	return positions, nil
}

type fetchEtasRequest struct {
	Target Position   `json:"target"`
	Source []Position `json:"source"`
}

func (s *ApiRequestService) FetchEtas(position Position, carsPositions []Position) ([]Eta, error) {
	etas := []Eta{}
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

	err = json.Unmarshal(body, &etas)
	if err != nil {
		return etas, err
	}

	return etas, nil
}
