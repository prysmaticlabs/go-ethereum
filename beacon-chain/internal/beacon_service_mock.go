// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/prysmaticlabs/prysm/proto/beacon/rpc/v1 (interfaces: BeaconServiceServer,BeaconService_LatestCrystallizedStateServer,BeaconService_LatestAttestationServer)

// Package mock_v1 is a generated GoMock package.
package mock_v1

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	empty "github.com/golang/protobuf/ptypes/empty"
	v1 "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	v10 "github.com/prysmaticlabs/prysm/proto/beacon/rpc/v1"
	metadata "google.golang.org/grpc/metadata"
	reflect "reflect"
)

// MockBeaconServiceServer is a mock of BeaconServiceServer interface
type MockBeaconServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockBeaconServiceServerMockRecorder
}

// MockBeaconServiceServerMockRecorder is the mock recorder for MockBeaconServiceServer
type MockBeaconServiceServerMockRecorder struct {
	mock *MockBeaconServiceServer
}

// NewMockBeaconServiceServer creates a new mock instance
func NewMockBeaconServiceServer(ctrl *gomock.Controller) *MockBeaconServiceServer {
	mock := &MockBeaconServiceServer{ctrl: ctrl}
	mock.recorder = &MockBeaconServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBeaconServiceServer) EXPECT() *MockBeaconServiceServerMockRecorder {
	return m.recorder
}

// CanonicalHead mocks base method
func (m *MockBeaconServiceServer) CanonicalHead(arg0 context.Context, arg1 *empty.Empty) (*v1.BeaconBlock, error) {
	ret := m.ctrl.Call(m, "CanonicalHead", arg0, arg1)
	ret0, _ := ret[0].(*v1.BeaconBlock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CanonicalHead indicates an expected call of CanonicalHead
func (mr *MockBeaconServiceServerMockRecorder) CanonicalHead(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CanonicalHead", reflect.TypeOf((*MockBeaconServiceServer)(nil).CanonicalHead), arg0, arg1)
}

// GenesisTimeAndCanonicalState mocks base method
func (m *MockBeaconServiceServer) GenesisTimeAndCanonicalState(arg0 context.Context, arg1 *empty.Empty) (*v10.GenesisTimeAndStateResponse, error) {
	ret := m.ctrl.Call(m, "GenesisTimeAndCanonicalState", arg0, arg1)
	ret0, _ := ret[0].(*v10.GenesisTimeAndStateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenesisTimeAndCanonicalState indicates an expected call of GenesisTimeAndCanonicalState
func (mr *MockBeaconServiceServerMockRecorder) GenesisTimeAndCanonicalState(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenesisTimeAndCanonicalState", reflect.TypeOf((*MockBeaconServiceServer)(nil).GenesisTimeAndCanonicalState), arg0, arg1)
}

// LatestAttestation mocks base method
func (m *MockBeaconServiceServer) LatestAttestation(arg0 *empty.Empty, arg1 v10.BeaconService_LatestAttestationServer) error {
	ret := m.ctrl.Call(m, "LatestAttestation", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// LatestAttestation indicates an expected call of LatestAttestation
func (mr *MockBeaconServiceServerMockRecorder) LatestAttestation(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LatestAttestation", reflect.TypeOf((*MockBeaconServiceServer)(nil).LatestAttestation), arg0, arg1)
}

// LatestCrystallizedState mocks base method
func (m *MockBeaconServiceServer) LatestCrystallizedState(arg0 *empty.Empty, arg1 v10.BeaconService_LatestCrystallizedStateServer) error {
	ret := m.ctrl.Call(m, "LatestCrystallizedState", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// LatestCrystallizedState indicates an expected call of LatestCrystallizedState
func (mr *MockBeaconServiceServerMockRecorder) LatestCrystallizedState(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LatestCrystallizedState", reflect.TypeOf((*MockBeaconServiceServer)(nil).LatestCrystallizedState), arg0, arg1)
}

// MockBeaconService_LatestCrystallizedStateServer is a mock of BeaconService_LatestCrystallizedStateServer interface
type MockBeaconService_LatestCrystallizedStateServer struct {
	ctrl     *gomock.Controller
	recorder *MockBeaconService_LatestCrystallizedStateServerMockRecorder
}

// MockBeaconService_LatestCrystallizedStateServerMockRecorder is the mock recorder for MockBeaconService_LatestCrystallizedStateServer
type MockBeaconService_LatestCrystallizedStateServerMockRecorder struct {
	mock *MockBeaconService_LatestCrystallizedStateServer
}

// NewMockBeaconService_LatestCrystallizedStateServer creates a new mock instance
func NewMockBeaconService_LatestCrystallizedStateServer(ctrl *gomock.Controller) *MockBeaconService_LatestCrystallizedStateServer {
	mock := &MockBeaconService_LatestCrystallizedStateServer{ctrl: ctrl}
	mock.recorder = &MockBeaconService_LatestCrystallizedStateServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBeaconService_LatestCrystallizedStateServer) EXPECT() *MockBeaconService_LatestCrystallizedStateServerMockRecorder {
	return m.recorder
}

// Context mocks base method
func (m *MockBeaconService_LatestCrystallizedStateServer) Context() context.Context {
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockBeaconService_LatestCrystallizedStateServerMockRecorder) Context() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockBeaconService_LatestCrystallizedStateServer)(nil).Context))
}

// RecvMsg mocks base method
func (m *MockBeaconService_LatestCrystallizedStateServer) RecvMsg(arg0 interface{}) error {
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg
func (mr *MockBeaconService_LatestCrystallizedStateServerMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockBeaconService_LatestCrystallizedStateServer)(nil).RecvMsg), arg0)
}

// Send mocks base method
func (m *MockBeaconService_LatestCrystallizedStateServer) Send(arg0 *v1.CrystallizedState) error {
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send
func (mr *MockBeaconService_LatestCrystallizedStateServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockBeaconService_LatestCrystallizedStateServer)(nil).Send), arg0)
}

// SendHeader mocks base method
func (m *MockBeaconService_LatestCrystallizedStateServer) SendHeader(arg0 metadata.MD) error {
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader
func (mr *MockBeaconService_LatestCrystallizedStateServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockBeaconService_LatestCrystallizedStateServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method
func (m *MockBeaconService_LatestCrystallizedStateServer) SendMsg(arg0 interface{}) error {
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg
func (mr *MockBeaconService_LatestCrystallizedStateServerMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockBeaconService_LatestCrystallizedStateServer)(nil).SendMsg), arg0)
}

// SetHeader mocks base method
func (m *MockBeaconService_LatestCrystallizedStateServer) SetHeader(arg0 metadata.MD) error {
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader
func (mr *MockBeaconService_LatestCrystallizedStateServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockBeaconService_LatestCrystallizedStateServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method
func (m *MockBeaconService_LatestCrystallizedStateServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer
func (mr *MockBeaconService_LatestCrystallizedStateServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockBeaconService_LatestCrystallizedStateServer)(nil).SetTrailer), arg0)
}

// MockBeaconService_LatestAttestationServer is a mock of BeaconService_LatestAttestationServer interface
type MockBeaconService_LatestAttestationServer struct {
	ctrl     *gomock.Controller
	recorder *MockBeaconService_LatestAttestationServerMockRecorder
}

// MockBeaconService_LatestAttestationServerMockRecorder is the mock recorder for MockBeaconService_LatestAttestationServer
type MockBeaconService_LatestAttestationServerMockRecorder struct {
	mock *MockBeaconService_LatestAttestationServer
}

// NewMockBeaconService_LatestAttestationServer creates a new mock instance
func NewMockBeaconService_LatestAttestationServer(ctrl *gomock.Controller) *MockBeaconService_LatestAttestationServer {
	mock := &MockBeaconService_LatestAttestationServer{ctrl: ctrl}
	mock.recorder = &MockBeaconService_LatestAttestationServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBeaconService_LatestAttestationServer) EXPECT() *MockBeaconService_LatestAttestationServerMockRecorder {
	return m.recorder
}

// Context mocks base method
func (m *MockBeaconService_LatestAttestationServer) Context() context.Context {
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockBeaconService_LatestAttestationServerMockRecorder) Context() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockBeaconService_LatestAttestationServer)(nil).Context))
}

// RecvMsg mocks base method
func (m *MockBeaconService_LatestAttestationServer) RecvMsg(arg0 interface{}) error {
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg
func (mr *MockBeaconService_LatestAttestationServerMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockBeaconService_LatestAttestationServer)(nil).RecvMsg), arg0)
}

// Send mocks base method
func (m *MockBeaconService_LatestAttestationServer) Send(arg0 *v1.AggregatedAttestation) error {
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send
func (mr *MockBeaconService_LatestAttestationServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockBeaconService_LatestAttestationServer)(nil).Send), arg0)
}

// SendHeader mocks base method
func (m *MockBeaconService_LatestAttestationServer) SendHeader(arg0 metadata.MD) error {
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader
func (mr *MockBeaconService_LatestAttestationServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockBeaconService_LatestAttestationServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method
func (m *MockBeaconService_LatestAttestationServer) SendMsg(arg0 interface{}) error {
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg
func (mr *MockBeaconService_LatestAttestationServerMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockBeaconService_LatestAttestationServer)(nil).SendMsg), arg0)
}

// SetHeader mocks base method
func (m *MockBeaconService_LatestAttestationServer) SetHeader(arg0 metadata.MD) error {
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader
func (mr *MockBeaconService_LatestAttestationServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockBeaconService_LatestAttestationServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method
func (m *MockBeaconService_LatestAttestationServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer
func (mr *MockBeaconService_LatestAttestationServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockBeaconService_LatestAttestationServer)(nil).SetTrailer), arg0)
}
