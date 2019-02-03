package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const (
	minLat = -90.0
	maxLat = 90.0
	minLng = -180.0
	maxLng = 180.0
)

var (
	ErrValidatedParams = errors.New("invalid params")
)

func MakeHandler(service *EtaService, logger *log.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		lat, lng, err := validatedParams(r)
		if err != nil {
			handleError(w, err)
			return
		}

		minETA, err := service.FindMinEta(Position{lat, lng})
		if err != nil {
			logger.Print(err)
			handleError(w, err)
			return
		}

		payloadBytes, err := json.Marshal(minETA)
		if err != nil {
			handleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(payloadBytes)
	}
}

func validatedParams(r *http.Request) (float64, float64, error) {
	var lat, lng float64

	latParam := r.URL.Query().Get("lat")
	if latParam == "" {
		return lat, lng, ErrValidatedParams
	}

	lngParam := r.URL.Query().Get("lng")
	if lngParam == "" {
		return lat, lng, ErrValidatedParams
	}

	lat, err := strconv.ParseFloat(latParam, 64)
	if err != nil {
		return lat, lng, ErrValidatedParams
	}
	if lat < minLat || lat > maxLat {
		return lat, lng, ErrValidatedParams
	}

	lng, err = strconv.ParseFloat(lngParam, 64)
	if err != nil {
		return lat, lng, ErrValidatedParams
	}
	if lng < minLng || lng > maxLng {
		return lat, lng, ErrValidatedParams
	}

	return lat, lng, nil
}

func handleError(w http.ResponseWriter, err error) {
	switch err {
	case ErrFetchCarsPositions:
		http.Error(w, "couldn't fetch car positions", http.StatusServiceUnavailable)
	case ErrFetchEtas:
		http.Error(w, "couldn't fetch etas", http.StatusServiceUnavailable)
	case ErrValidatedParams:
		errText := fmt.Sprintf(
			"Passed params incorrect: lat is required and value should be from %v to %v; lng is required and value should be from %v to %v",
			minLat, maxLat,
			minLng, maxLng,
		)
		http.Error(w, errText, http.StatusUnprocessableEntity)
	default:
		http.Error(w, "internal errros", http.StatusInternalServerError)
	}
}
