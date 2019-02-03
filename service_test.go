package eta

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEtaService_FindMinEta(t *testing.T) {
	t.Run("when service is called with correct params", func(t *testing.T) {
		mCtrl := gomock.NewController(t)
		defer mCtrl.Finish()
		apiRequestMock := NewMockApiRequester(mCtrl)

		etaService := NewEtaService(apiRequestMock)

		targetPosition := Position{55.752992, 37.618333}
		sourcePositions := []Position{{55.7575429, 37.6135117}, {55.74837156167371, 37.61180107665421}, {55.7532706, 37.6076902}}
		etas := []Eta{1, 1, 1}

		apiRequestMock.EXPECT().FetchCarPositions(targetPosition, gomock.Any()).Return(sourcePositions, nil)
		apiRequestMock.EXPECT().FetchEtas(targetPosition, sourcePositions).Return(etas, nil)

		resp, err := etaService.FindMinEta(targetPosition)
		assert.Equal(t, nil, err)
		assert.Equal(t, Eta(1), resp)
	})

	t.Run("when service is called with correct critical params", func(t *testing.T) {
		mCtrl := gomock.NewController(t)
		defer mCtrl.Finish()
		apiRequestMock := NewMockApiRequester(mCtrl)

		etaService := NewEtaService(apiRequestMock)

		targetPosition := Position{-90, 180}
		sourcePositions := []Position{{55.338315, 38.194654}, {55.41528466287064, 37.901831957545426}, {55.4153892, 37.8980953}}
		etas := []Eta{16161, 16170, 16170}

		apiRequestMock.EXPECT().FetchCarPositions(targetPosition, gomock.Any()).Return(sourcePositions, nil)
		apiRequestMock.EXPECT().FetchEtas(targetPosition, sourcePositions).Return(etas, nil)

		resp, err := etaService.FindMinEta(targetPosition)
		assert.Equal(t, nil, err)
		assert.Equal(t, Eta(16161), resp)
	})

	t.Run("when cars positions service is unavailable", func(t *testing.T) {
		mCtrl := gomock.NewController(t)
		defer mCtrl.Finish()
		apiRequestMock := NewMockApiRequester(mCtrl)

		etaService := NewEtaService(apiRequestMock)

		targetPosition := Position{}
		apiRequestMock.EXPECT().FetchCarPositions(targetPosition, gomock.Any()).Return([]Position{}, errors.New("service error"))

		_, err := etaService.FindMinEta(targetPosition)
		assert.Equal(t, ErrFetchCarsPositions, err)
	})

	t.Run("when cars positions service returns empty result", func(t *testing.T) {
		mCtrl := gomock.NewController(t)
		defer mCtrl.Finish()
		apiRequestMock := NewMockApiRequester(mCtrl)

		etaService := NewEtaService(apiRequestMock)

		targetPosition := Position{}
		apiRequestMock.EXPECT().FetchCarPositions(targetPosition, gomock.Any()).Return([]Position{}, nil)

		_, err := etaService.FindMinEta(targetPosition)
		assert.Equal(t, ErrEmptyCarsPositions, err)
	})

	t.Run("when get ETAs service is unavailable", func(t *testing.T) {
		mCtrl := gomock.NewController(t)
		defer mCtrl.Finish()
		apiRequestMock := NewMockApiRequester(mCtrl)

		etaService := NewEtaService(apiRequestMock)

		targetPosition := Position{}
		sourcePositions := []Position{{0, 0}}
		etas := []Eta{}

		apiRequestMock.EXPECT().FetchCarPositions(targetPosition, gomock.Any()).Return(sourcePositions, nil)
		apiRequestMock.EXPECT().FetchEtas(targetPosition, sourcePositions).Return(etas, errors.New("service error"))

		_, err := etaService.FindMinEta(targetPosition)
		assert.Equal(t, ErrFetchEtas, err)
	})

	t.Run("when get ETAs service returns empty result", func(t *testing.T) {
		mCtrl := gomock.NewController(t)
		defer mCtrl.Finish()
		apiRequestMock := NewMockApiRequester(mCtrl)

		etaService := NewEtaService(apiRequestMock)

		targetPosition := Position{}
		sourcePositions := []Position{{0, 0}}
		etas := []Eta{}

		apiRequestMock.EXPECT().FetchCarPositions(targetPosition, gomock.Any()).Return(sourcePositions, nil)
		apiRequestMock.EXPECT().FetchEtas(targetPosition, sourcePositions).Return(etas, nil)

		_, err := etaService.FindMinEta(targetPosition)
		assert.Equal(t, ErrEmptyEtas, err)
	})
}
