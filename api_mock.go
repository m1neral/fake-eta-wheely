// Code generated by MockGen. DO NOT EDIT.
// Source: api.go

// Package eta is a generated GoMock package.
package eta

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockApiRequester is a mock of ApiRequester interface
type MockApiRequester struct {
	ctrl     *gomock.Controller
	recorder *MockApiRequesterMockRecorder
}

// MockApiRequesterMockRecorder is the mock recorder for MockApiRequester
type MockApiRequesterMockRecorder struct {
	mock *MockApiRequester
}

// NewMockApiRequester creates a new mock instance
func NewMockApiRequester(ctrl *gomock.Controller) *MockApiRequester {
	mock := &MockApiRequester{ctrl: ctrl}
	mock.recorder = &MockApiRequesterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockApiRequester) EXPECT() *MockApiRequesterMockRecorder {
	return m.recorder
}

// FetchCarPositions mocks base method
func (m *MockApiRequester) FetchCarPositions(position Position, limit int) ([]Position, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchCarPositions", position, limit)
	ret0, _ := ret[0].([]Position)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchCarPositions indicates an expected call of FetchCarPositions
func (mr *MockApiRequesterMockRecorder) FetchCarPositions(position, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchCarPositions", reflect.TypeOf((*MockApiRequester)(nil).FetchCarPositions), position, limit)
}

// FetchEtas mocks base method
func (m *MockApiRequester) FetchEtas(position Position, carsPositions []Position) ([]Eta, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchEtas", position, carsPositions)
	ret0, _ := ret[0].([]Eta)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchEtas indicates an expected call of FetchEtas
func (mr *MockApiRequesterMockRecorder) FetchEtas(position, carsPositions interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchEtas", reflect.TypeOf((*MockApiRequester)(nil).FetchEtas), position, carsPositions)
}
