// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package repository is a generated GoMock package.
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

// FetchDiskSizePer100MTxs mocks base method.
func (m *MockRepository) FetchDiskSizePer100MTxs() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchDiskSizePer100MTxs")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchDiskSizePer100MTxs indicates an expected call of FetchDiskSizePer100MTxs.
func (mr *MockRepositoryMockRecorder) FetchDiskSizePer100MTxs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchDiskSizePer100MTxs", reflect.TypeOf((*MockRepository)(nil).FetchDiskSizePer100MTxs))
}

// FetchNumberOfAccounts mocks base method.
func (m *MockRepository) FetchNumberOfAccounts() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchNumberOfAccounts")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchNumberOfAccounts indicates an expected call of FetchNumberOfAccounts.
func (mr *MockRepositoryMockRecorder) FetchNumberOfAccounts() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchNumberOfAccounts", reflect.TypeOf((*MockRepository)(nil).FetchNumberOfAccounts))
}

// FetchTimeToFinality mocks base method.
func (m *MockRepository) FetchTimeToFinality() (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchTimeToFinality")
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchTimeToFinality indicates an expected call of FetchTimeToFinality.
func (mr *MockRepositoryMockRecorder) FetchTimeToFinality() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchTimeToFinality", reflect.TypeOf((*MockRepository)(nil).FetchTimeToFinality))
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

// GetDiskSizePer100MTxs mocks base method.
func (m *MockRepository) GetDiskSizePer100MTxs() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDiskSizePer100MTxs")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetDiskSizePer100MTxs indicates an expected call of GetDiskSizePer100MTxs.
func (mr *MockRepositoryMockRecorder) GetDiskSizePer100MTxs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDiskSizePer100MTxs", reflect.TypeOf((*MockRepository)(nil).GetDiskSizePer100MTxs))
}

// GetGasUsedAggByTimestamp mocks base method.
func (m *MockRepository) GetGasUsedAggByTimestamp(arg0 types.AggResolution, arg1 uint, arg2 *uint64) ([]types.HexUintTick, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGasUsedAggByTimestamp", arg0, arg1, arg2)
	ret0, _ := ret[0].([]types.HexUintTick)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGasUsedAggByTimestamp indicates an expected call of GetGasUsedAggByTimestamp.
func (mr *MockRepositoryMockRecorder) GetGasUsedAggByTimestamp(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGasUsedAggByTimestamp", reflect.TypeOf((*MockRepository)(nil).GetGasUsedAggByTimestamp), arg0, arg1, arg2)
}

// GetGasUsedPer10Secs mocks base method.
func (m *MockRepository) GetGasUsedPer10Secs() []types.HexUintTick {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGasUsedPer10Secs")
	ret0, _ := ret[0].([]types.HexUintTick)
	return ret0
}

// GetGasUsedPer10Secs indicates an expected call of GetGasUsedPer10Secs.
func (mr *MockRepositoryMockRecorder) GetGasUsedPer10Secs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGasUsedPer10Secs", reflect.TypeOf((*MockRepository)(nil).GetGasUsedPer10Secs))
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

// GetNumberOfAccounts mocks base method.
func (m *MockRepository) GetNumberOfAccounts() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNumberOfAccounts")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetNumberOfAccounts indicates an expected call of GetNumberOfAccounts.
func (mr *MockRepositoryMockRecorder) GetNumberOfAccounts() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNumberOfAccounts", reflect.TypeOf((*MockRepository)(nil).GetNumberOfAccounts))
}

// GetNumberOfValidators mocks base method.
func (m *MockRepository) GetNumberOfValidators() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNumberOfValidators")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNumberOfValidators indicates an expected call of GetNumberOfValidators.
func (mr *MockRepositoryMockRecorder) GetNumberOfValidators() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNumberOfValidators", reflect.TypeOf((*MockRepository)(nil).GetNumberOfValidators))
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

// GetTrxCount mocks base method.
func (m *MockRepository) GetTrxCount() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrxCount")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrxCount indicates an expected call of GetTrxCount.
func (mr *MockRepositoryMockRecorder) GetTrxCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrxCount", reflect.TypeOf((*MockRepository)(nil).GetTrxCount))
}

// GetTrxCountAggByTimestamp mocks base method.
func (m *MockRepository) GetTrxCountAggByTimestamp(arg0 types.AggResolution, arg1 uint, arg2 *uint64) ([]types.HexUintTick, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrxCountAggByTimestamp", arg0, arg1, arg2)
	ret0, _ := ret[0].([]types.HexUintTick)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrxCountAggByTimestamp indicates an expected call of GetTrxCountAggByTimestamp.
func (mr *MockRepositoryMockRecorder) GetTrxCountAggByTimestamp(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrxCountAggByTimestamp", reflect.TypeOf((*MockRepository)(nil).GetTrxCountAggByTimestamp), arg0, arg1, arg2)
}

// GetTxCountPer10Secs mocks base method.
func (m *MockRepository) GetTxCountPer10Secs() []types.HexUintTick {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTxCountPer10Secs")
	ret0, _ := ret[0].([]types.HexUintTick)
	return ret0
}

// GetTxCountPer10Secs indicates an expected call of GetTxCountPer10Secs.
func (mr *MockRepositoryMockRecorder) GetTxCountPer10Secs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTxCountPer10Secs", reflect.TypeOf((*MockRepository)(nil).GetTxCountPer10Secs))
}

// IncrementTrxCount mocks base method.
func (m *MockRepository) IncrementTrxCount(arg0 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrementTrxCount", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrementTrxCount indicates an expected call of IncrementTrxCount.
func (mr *MockRepositoryMockRecorder) IncrementTrxCount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrementTrxCount", reflect.TypeOf((*MockRepository)(nil).IncrementTrxCount), arg0)
}

// SetDiskSizePer100MTxs mocks base method.
func (m *MockRepository) SetDiskSizePer100MTxs(arg0 uint64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetDiskSizePer100MTxs", arg0)
}

// SetDiskSizePer100MTxs indicates an expected call of SetDiskSizePer100MTxs.
func (mr *MockRepositoryMockRecorder) SetDiskSizePer100MTxs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDiskSizePer100MTxs", reflect.TypeOf((*MockRepository)(nil).SetDiskSizePer100MTxs), arg0)
}

// SetGasUsedPer10Secs mocks base method.
func (m *MockRepository) SetGasUsedPer10Secs(arg0 []types.HexUintTick) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetGasUsedPer10Secs", arg0)
}

// SetGasUsedPer10Secs indicates an expected call of SetGasUsedPer10Secs.
func (mr *MockRepositoryMockRecorder) SetGasUsedPer10Secs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetGasUsedPer10Secs", reflect.TypeOf((*MockRepository)(nil).SetGasUsedPer10Secs), arg0)
}

// SetNumberOfAccounts mocks base method.
func (m *MockRepository) SetNumberOfAccounts(arg0 uint64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetNumberOfAccounts", arg0)
}

// SetNumberOfAccounts indicates an expected call of SetNumberOfAccounts.
func (mr *MockRepositoryMockRecorder) SetNumberOfAccounts(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNumberOfAccounts", reflect.TypeOf((*MockRepository)(nil).SetNumberOfAccounts), arg0)
}

// SetTxCountPer10Secs mocks base method.
func (m *MockRepository) SetTxCountPer10Secs(arg0 []types.HexUintTick) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTxCountPer10Secs", arg0)
}

// SetTxCountPer10Secs indicates an expected call of SetTxCountPer10Secs.
func (mr *MockRepositoryMockRecorder) SetTxCountPer10Secs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTxCountPer10Secs", reflect.TypeOf((*MockRepository)(nil).SetTxCountPer10Secs), arg0)
}

// UpdateLatestObservedBlock mocks base method.
func (m *MockRepository) UpdateLatestObservedBlock(arg0 *types.Block) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLatestObservedBlock", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateLatestObservedBlock indicates an expected call of UpdateLatestObservedBlock.
func (mr *MockRepositoryMockRecorder) UpdateLatestObservedBlock(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLatestObservedBlock", reflect.TypeOf((*MockRepository)(nil).UpdateLatestObservedBlock), arg0)
}
