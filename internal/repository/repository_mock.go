// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/interface.go

// Package mock_repository is a generated GoMock package.
package repository

import (
	types "ftm-explorer/internal/types"
	reflect "reflect"

	common "github.com/ethereum/go-ethereum/common"
	types0 "github.com/ethereum/go-ethereum/core/types"
	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of IRepository interface.
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

// GetBlockByNumber mocks base method.
func (m *MockRepository) GetBlockByNumber(arg0 uint64) (*types.Block, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockByNumber", arg0)
	ret0, _ := ret[0].(*types.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockByNumber indicates an expected call of GetBlockByNumber.
func (mr *MockRepositoryMockRecorder) GetBlockByNumber(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockByNumber", reflect.TypeOf((*MockRepository)(nil).GetBlockByNumber), arg0)
}

// GetLatestObservedBlock mocks base method.
func (m *MockRepository) GetLatestObservedBlock() *types.Block {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestObservedBlock")
	ret0, _ := ret[0].(*types.Block)
	return ret0
}

// GetLatestObservedBlock indicates an expected call of GetLatestObservedBlock.
func (mr *MockRepositoryMockRecorder) GetLatestObservedBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestObservedBlock", reflect.TypeOf((*MockRepository)(nil).GetLatestObservedBlock))
}

// GetLatestObservedBlocks mocks base method.
func (m *MockRepository) GetLatestObservedBlocks(arg0 uint) []*types.Block {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestObservedBlocks", arg0)
	ret0, _ := ret[0].([]*types.Block)
	return ret0
}

// GetLatestObservedBlocks indicates an expected call of GetLatestObservedBlocks.
func (mr *MockRepositoryMockRecorder) GetLatestObservedBlocks(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestObservedBlocks", reflect.TypeOf((*MockRepository)(nil).GetLatestObservedBlocks), arg0)
}

// GetNewHeadersChannel mocks base method.
func (m *MockRepository) GetNewHeadersChannel() <-chan *types0.Header {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNewHeadersChannel")
	ret0, _ := ret[0].(<-chan *types0.Header)
	return ret0
}

// GetNewHeadersChannel indicates an expected call of GetNewHeadersChannel.
func (mr *MockRepositoryMockRecorder) GetNewHeadersChannel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNewHeadersChannel", reflect.TypeOf((*MockRepository)(nil).GetNewHeadersChannel))
}

// GetTransactionByHash mocks base method.
func (m *MockRepository) GetTransactionByHash(arg0 common.Hash) (*types.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionByHash", arg0)
	ret0, _ := ret[0].(*types.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactionByHash indicates an expected call of GetTransactionByHash.
func (mr *MockRepositoryMockRecorder) GetTransactionByHash(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionByHash", reflect.TypeOf((*MockRepository)(nil).GetTransactionByHash), arg0)
}

// UpdateLatestObservedBlock mocks base method.
func (m *MockRepository) UpdateLatestObservedBlock(arg0 *types.Block) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateLatestObservedBlock", arg0)
}

// UpdateLatestObservedBlock indicates an expected call of UpdateLatestObservedBlock.
func (mr *MockRepositoryMockRecorder) UpdateLatestObservedBlock(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLatestObservedBlock", reflect.TypeOf((*MockRepository)(nil).UpdateLatestObservedBlock), arg0)
}
