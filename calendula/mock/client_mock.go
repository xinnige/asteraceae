// Code generated by MockGen. DO NOT EDIT.
// Source: astermisc/client.go

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	http "net/http"
	reflect "reflect"
)

// MockAsterClient is a mock of AsterClient interface
type MockAsterClient struct {
	ctrl     *gomock.Controller
	recorder *MockAsterClientMockRecorder
}

// MockAsterClientMockRecorder is the mock recorder for MockAsterClient
type MockAsterClientMockRecorder struct {
	mock *MockAsterClient
}

// NewMockAsterClient creates a new mock instance
func NewMockAsterClient(ctrl *gomock.Controller) *MockAsterClient {
	mock := &MockAsterClient{ctrl: ctrl}
	mock.recorder = &MockAsterClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAsterClient) EXPECT() *MockAsterClientMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockAsterClient) Do(arg0 *http.Request) (*http.Response, error) {
	ret := m.ctrl.Call(m, "Do", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do
func (mr *MockAsterClientMockRecorder) Do(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockAsterClient)(nil).Do), arg0)
}
