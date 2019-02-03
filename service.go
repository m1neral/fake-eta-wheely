package eta

import (
	"errors"
)

type EtaService struct {
	limitCars  int
	apiRequest ApiRequester
}

var (
	ErrFetchCarsPositions = errors.New("couldn't fetch cars positions")
	ErrEmptyCarsPositions = errors.New("empty cars positions")
	ErrFetchEtas          = errors.New("couldn't fetch etas")
	ErrEmptyEtas          = errors.New("empty etas")
	ErrInternal           = errors.New("couldn't detect min ETA")
)

var (
	defaultLimitCars = 100
)

func NewEtaService(apiRequest ApiRequester) *EtaService {
	return &EtaService{limitCars: defaultLimitCars, apiRequest: apiRequest}
}

func (s *EtaService) FindMinEta(targetPosition Position) (Eta, error) {
	positions, err := s.getCarPositions(targetPosition)
	if err != nil {
		return 0, err
	}
	if len(positions) == 0 {
		return 0, ErrEmptyCarsPositions
	}

	etas, err := s.getEtas(targetPosition, positions)
	if err != nil {
		return 0, err
	}
	if len(etas) == 0 {
		return 0, ErrEmptyEtas
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
		return resp, ErrFetchEtas
	}
	return resp, nil
}

func (s *EtaService) getCarPositions(position Position) ([]Position, error) {
	resp, err := s.apiRequest.FetchCarPositions(position, s.limitCars)
	if err != nil {
		return resp, ErrFetchCarsPositions
	}

	return resp, err
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
