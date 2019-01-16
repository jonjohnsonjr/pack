// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/buildpack/pack/commands (interfaces: ImageFactory)

// Package mocks is a generated GoMock package.
package mocks

import (
	image "github.com/buildpack/lifecycle/image"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockImageFactory is a mock of ImageFactory interface
type MockImageFactory struct {
	ctrl     *gomock.Controller
	recorder *MockImageFactoryMockRecorder
}

// MockImageFactoryMockRecorder is the mock recorder for MockImageFactory
type MockImageFactoryMockRecorder struct {
	mock *MockImageFactory
}

// NewMockImageFactory creates a new mock instance
func NewMockImageFactory(ctrl *gomock.Controller) *MockImageFactory {
	mock := &MockImageFactory{ctrl: ctrl}
	mock.recorder = &MockImageFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockImageFactory) EXPECT() *MockImageFactoryMockRecorder {
	return m.recorder
}

// NewLocal mocks base method
func (m *MockImageFactory) NewLocal(arg0 string, arg1 bool) (image.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewLocal", arg0, arg1)
	ret0, _ := ret[0].(image.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewLocal indicates an expected call of NewLocal
func (mr *MockImageFactoryMockRecorder) NewLocal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewLocal", reflect.TypeOf((*MockImageFactory)(nil).NewLocal), arg0, arg1)
}

// NewRemote mocks base method
func (m *MockImageFactory) NewRemote(arg0 string) (image.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRemote", arg0)
	ret0, _ := ret[0].(image.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewRemote indicates an expected call of NewRemote
func (mr *MockImageFactoryMockRecorder) NewRemote(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRemote", reflect.TypeOf((*MockImageFactory)(nil).NewRemote), arg0)
}