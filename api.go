package main

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
	return FetchCarPositionsResponse{}, nil
}

func (s *ApiRequestService) FetchEtas(position Position, carsPositions []Position) (FetchEtasResponse, error) {
	return FetchEtasResponse{}, nil
}
