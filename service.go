package main

import (
	"errors"
)

type EtaService struct {
	limitCars  int
	apiRequest ApiRequester
}

var (
	ErrFetchCarsPositions = errors.New("couldn't fetch car positions")
	ErrFetchEtas          = errors.New("couldn't fetch etas")
	ErrInternal           = errors.New("coldn't detect min ETA")
)

var (
	defaultLimitCars = 100
)

func NewEtaService(apiRequest ApiRequester) *EtaService {
	return &EtaService{limitCars: defaultLimitCars, apiRequest: apiRequest}
}

func (s *EtaService) FindMinEta(targetPosition Position) (Eta, error) {
	postions, err := s.getCarPositions(targetPosition)
	if err != nil {
		return 0, err
	}

	etas, err := s.getEtas(targetPosition, postions)
	if err != nil {
		return 0, err
	}

	minEta, err := s.min(etas)
	if err != nil {
		return 0, ErrInternal
	}

	return minEta, nil
}

func (s *EtaService) getEtas(position Position, carsPositions []Position) ([]Eta, error) {
	resp, err := s.apiRequest.FetchEtas(position, carsPositions)
	if err != nil {
		return resp.Etas, ErrFetchEtas
	}
	return resp.Etas, nil
}

func (s *EtaService) getCarPositions(position Position) ([]Position, error) {
	resp, err := s.apiRequest.FetchCarPositions(position, s.limitCars)
	if err != nil {
		return resp.Positions, ErrFetchCarsPositions
	}

	return resp.Positions, err
}

func (s *EtaService) min(values []Eta) (Eta, error) {
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
