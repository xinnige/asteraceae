// Code generated by MockGen. DO NOT EDIT.
// Source: api/awsapi.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	session "github.com/aws/aws-sdk-go/aws/session"
	gomock "github.com/golang/mock/gomock"
)

// MockAWSSession is a mock of AWSSession interface
type MockAWSSession struct {
	ctrl     *gomock.Controller
	recorder *MockAWSSessionMockRecorder
}

// MockAWSSessionMockRecorder is the mock recorder for MockAWSSession
type MockAWSSessionMockRecorder struct {
	mock *MockAWSSession
}

// NewMockAWSSession creates a new mock instance
func NewMockAWSSession(ctrl *gomock.Controller) *MockAWSSession {
	mock := &MockAWSSession{ctrl: ctrl}
	mock.recorder = &MockAWSSessionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAWSSession) EXPECT() *MockAWSSessionMockRecorder {
	return m.recorder
}

// NewSession mocks base method
func (m *MockAWSSession) NewSession() (*session.Session, error) {
	ret := m.ctrl.Call(m, "NewSession")
	ret0, _ := ret[0].(*session.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewSession indicates an expected call of NewSession
func (mr *MockAWSSessionMockRecorder) NewSession() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewSession", reflect.TypeOf((*MockAWSSession)(nil).NewSession))
}