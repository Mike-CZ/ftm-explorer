// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package meta_fetcher is a generated GoMock package.
package meta_fetcher

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMetaFetcher is a mock of IMetaFetcher interface.
type MockMetaFetcher struct {
	ctrl     *gomock.Controller
	recorder *MockMetaFetcherMockRecorder
}

// MockMetaFetcherMockRecorder is the mock recorder for MockMetaFetcher.
type MockMetaFetcherMockRecorder struct {
	mock *MockMetaFetcher
}

// NewMockMetaFetcher creates a new mock instance.
func NewMockMetaFetcher(ctrl *gomock.Controller) *MockMetaFetcher {
	mock := &MockMetaFetcher{ctrl: ctrl}
	mock.recorder = &MockMetaFetcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetaFetcher) EXPECT() *MockMetaFetcherMockRecorder {
	return m.recorder
}

// DiskSizePer100MTxs mocks base method.
func (m *MockMetaFetcher) DiskSizePer100MTxs() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DiskSizePer100MTxs")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DiskSizePer100MTxs indicates an expected call of DiskSizePer100MTxs.
func (mr *MockMetaFetcherMockRecorder) DiskSizePer100MTxs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DiskSizePer100MTxs", reflect.TypeOf((*MockMetaFetcher)(nil).DiskSizePer100MTxs))
}

// DiskSizePrunedPer100MTxs mocks base method.
func (m *MockMetaFetcher) DiskSizePrunedPer100MTxs() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DiskSizePrunedPer100MTxs")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DiskSizePrunedPer100MTxs indicates an expected call of DiskSizePrunedPer100MTxs.
func (mr *MockMetaFetcherMockRecorder) DiskSizePrunedPer100MTxs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DiskSizePrunedPer100MTxs", reflect.TypeOf((*MockMetaFetcher)(nil).DiskSizePrunedPer100MTxs))
}

// IsIdleStatus mocks base method.
func (m *MockMetaFetcher) IsIdleStatus() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsIdleStatus")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsIdleStatus indicates an expected call of IsIdleStatus.
func (mr *MockMetaFetcherMockRecorder) IsIdleStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsIdleStatus", reflect.TypeOf((*MockMetaFetcher)(nil).IsIdleStatus))
}

// NumberOfAccounts mocks base method.
func (m *MockMetaFetcher) NumberOfAccounts() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NumberOfAccounts")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NumberOfAccounts indicates an expected call of NumberOfAccounts.
func (mr *MockMetaFetcherMockRecorder) NumberOfAccounts() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NumberOfAccounts", reflect.TypeOf((*MockMetaFetcher)(nil).NumberOfAccounts))
}

// TimeToFinality mocks base method.
func (m *MockMetaFetcher) TimeToFinality() (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TimeToFinality")
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TimeToFinality indicates an expected call of TimeToFinality.
func (mr *MockMetaFetcherMockRecorder) TimeToFinality() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TimeToFinality", reflect.TypeOf((*MockMetaFetcher)(nil).TimeToFinality))
}
