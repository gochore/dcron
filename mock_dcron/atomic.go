// Code generated by MockGen. DO NOT EDIT.
// Source: atomic.go
//
// Generated by this command:
//
//	mockgen -source=atomic.go -destination mock_dcron/atomic.go
//
// Package mock_dcron is a generated GoMock package.
package mock_dcron

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockAtomic is a mock of Atomic interface.
type MockAtomic struct {
	ctrl     *gomock.Controller
	recorder *MockAtomicMockRecorder
}

// MockAtomicMockRecorder is the mock recorder for MockAtomic.
type MockAtomicMockRecorder struct {
	mock *MockAtomic
}

// NewMockAtomic creates a new mock instance.
func NewMockAtomic(ctrl *gomock.Controller) *MockAtomic {
	mock := &MockAtomic{ctrl: ctrl}
	mock.recorder = &MockAtomicMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAtomic) EXPECT() *MockAtomicMockRecorder {
	return m.recorder
}

// SetIfNotExists mocks base method.
func (m *MockAtomic) SetIfNotExists(ctx context.Context, key, value string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetIfNotExists", ctx, key, value)
	ret0, _ := ret[0].(bool)
	return ret0
}

// SetIfNotExists indicates an expected call of SetIfNotExists.
func (mr *MockAtomicMockRecorder) SetIfNotExists(ctx, key, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetIfNotExists", reflect.TypeOf((*MockAtomic)(nil).SetIfNotExists), ctx, key, value)
}
