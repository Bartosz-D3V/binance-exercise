// Code generated by MockGen. DO NOT EDIT.
// Source: ./app/transaction/repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	transaction "github.com/Bartosz-D3V/binance-exercise/app/transaction"
	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetAll mocks base method.
func (m *MockRepository) GetAll(ctx context.Context) []transaction.LogEntry {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]transaction.LogEntry)
	return ret0
}

// GetAll indicates an expected call of GetAll.
func (mr *MockRepositoryMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockRepository)(nil).GetAll), ctx)
}

// Save mocks base method.
func (m *MockRepository) Save(ctx context.Context, log transaction.LogEntry) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Save", ctx, log)
}

// Save indicates an expected call of Save.
func (mr *MockRepositoryMockRecorder) Save(ctx, log interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockRepository)(nil).Save), ctx, log)
}
