package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const (
	port = 8080
	path = "/eta"
)

type Position struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type PredictETAPayload struct {
	Target Position   `json:"target"`
	Source []Position `json:"source"`
}

func min(values []int) (int, error) {
	if len(values) == 0 {
		return 0, errors.New("cannot detect a minimum value in an empty slice")
	}

	minValue := values[0]

	for _, v := range values {
		if v < minValue {
			minValue = v
		}
	}

	return minValue, nil
}

func getCarsPositions(lat float64, lng float64) ([]Position, error) {
	limit := 100
	endpointURL := "https://dev-api.wheely.com/fake-eta/cars"
	positions := []Position{}

	url, err := url.Parse(endpointURL)
	if err != nil {
		return positions, err
	}

	query := url.Query()
	query.Set("limit", strconv.Itoa(limit))
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

	err = json.Unmarshal(body, &positions)
	if err != nil {
		return positions, err
	}

	return positions, nil
}

func getETAs(targetLat float64, targetLng float64, positions []Position) ([]int, error) {
	endpointURL := "https://dev-api.wheely.com/fake-eta/predict"
	etas := []int{}

	target := Position{Lat: targetLat, Lng: targetLng}
	payload := PredictETAPayload{Target: target, Source: positions}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return etas, err
	}

	resp, err := http.Post(endpointURL, "application/json; charset=utf-8", bytes.NewReader(payloadBytes))
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

func getMinETA(w http.ResponseWriter, r *http.Request) {
	latParam := r.URL.Query().Get("lat")
	if latParam == "" {
		http.Error(w, "lat in query is required", http.StatusUnprocessableEntity)
		return
	}

	lngParam := r.URL.Query().Get("lng")
	if lngParam == "" {
		http.Error(w, "lat in query is required", http.StatusUnprocessableEntity)
		return
	}

	lat, err := strconv.ParseFloat(latParam, 64)
	if err != nil {
		http.Error(w, "lat in query should be float64 formated", http.StatusUnprocessableEntity)
		return
	}
	if lat < -90 || lat > 90 {
		http.Error(w, "lat in query should be from -90 to 90", http.StatusUnprocessableEntity)
		return
	}

	lng, err := strconv.ParseFloat(lngParam, 64)
	if err != nil {
		http.Error(w, "lng in query should be float64 formated", http.StatusUnprocessableEntity)
		return
	}
	if lng < -180 || lng > 180 {
		http.Error(w, "lng in query should be from -180 to 180", http.StatusUnprocessableEntity)
		return
	}

	cars, err := getCarsPositions(lat, lng)
	if err != nil {
		http.Error(w, "error: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	if len(cars) == 0 {
		http.Error(w, "there are no cars by passed lat and lng", http.StatusNotFound)
		return
	}

	etas, err := getETAs(lat, lng, cars)
	if err != nil {
		http.Error(w, "error: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	if len(etas) == 0 {
		http.Error(w, "cannot predict minimal ETA by passed lat and lng", http.StatusServiceUnavailable)
		return
	}

	minETA, err := min(etas)
	if err != nil {
		http.Error(w, "error"+err.Error(), http.StatusInternalServerError)
		return
	}

	payloadBytes, err := json.Marshal(minETA)
	if err != nil {
		http.Error(w, "error"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payloadBytes)
}

func main() {
	http.HandleFunc(path, getMinETA)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
