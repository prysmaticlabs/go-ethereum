// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/prysmaticlabs/prysm/proto/beacon/rpc/v1 (interfaces: BeaconServiceClient,BeaconService_LatestBeaconBlockClient,BeaconService_LatestCrystallizedStateClient)

package internal

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	empty "github.com/golang/protobuf/ptypes/empty"
	v1 "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	v10 "github.com/prysmaticlabs/prysm/proto/beacon/rpc/v1"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)

// MockBeaconServiceClient is a mock of BeaconServiceClient interface
type MockBeaconServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockBeaconServiceClientMockRecorder
}

// MockBeaconServiceClientMockRecorder is the mock recorder for MockBeaconServiceClient
type MockBeaconServiceClientMockRecorder struct {
	mock *MockBeaconServiceClient
}

// NewMockBeaconServiceClient creates a new mock instance
func NewMockBeaconServiceClient(ctrl *gomock.Controller) *MockBeaconServiceClient {
	mock := &MockBeaconServiceClient{ctrl: ctrl}
	mock.recorder = &MockBeaconServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBeaconServiceClient) EXPECT() *MockBeaconServiceClientMockRecorder {
	return m.recorder
}

// CanonicalHeadAndState mocks base method
func (m *MockBeaconServiceClient) CanonicalHeadAndState(arg0 context.Context, arg1 *empty.Empty, arg2 ...grpc.CallOption) (*v10.CanonicalResponse, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CanonicalHeadAndState", varargs...)
	ret0, _ := ret[0].(*v10.CanonicalResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CanonicalHeadAndState indicates an expected call of CanonicalHeadAndState
func (mr *MockBeaconServiceClientMockRecorder) CanonicalHeadAndState(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CanonicalHeadAndState", reflect.TypeOf((*MockBeaconServiceClient)(nil).CanonicalHeadAndState), varargs...)
}

// FetchShuffledValidatorIndices mocks base method
func (m *MockBeaconServiceClient) FetchShuffledValidatorIndices(arg0 context.Context, arg1 *v10.ShuffleRequest, arg2 ...grpc.CallOption) (*v10.ShuffleResponse, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FetchShuffledValidatorIndices", varargs...)
	ret0, _ := ret[0].(*v10.ShuffleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchShuffledValidatorIndices indicates an expected call of FetchShuffledValidatorIndices
func (mr *MockBeaconServiceClientMockRecorder) FetchShuffledValidatorIndices(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchShuffledValidatorIndices", reflect.TypeOf((*MockBeaconServiceClient)(nil).FetchShuffledValidatorIndices), varargs...)
}

// LatestBeaconBlock mocks base method
func (m *MockBeaconServiceClient) LatestBeaconBlock(arg0 context.Context, arg1 *empty.Empty, arg2 ...grpc.CallOption) (v10.BeaconService_LatestBeaconBlockClient, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "LatestBeaconBlock", varargs...)
	ret0, _ := ret[0].(v10.BeaconService_LatestBeaconBlockClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LatestBeaconBlock indicates an expected call of LatestBeaconBlock
func (mr *MockBeaconServiceClientMockRecorder) LatestBeaconBlock(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LatestBeaconBlock", reflect.TypeOf((*MockBeaconServiceClient)(nil).LatestBeaconBlock), varargs...)
}

// LatestCrystallizedState mocks base method
func (m *MockBeaconServiceClient) LatestCrystallizedState(arg0 context.Context, arg1 *empty.Empty, arg2 ...grpc.CallOption) (v10.BeaconService_LatestCrystallizedStateClient, error) {
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "LatestCrystallizedState", varargs...)
	ret0, _ := ret[0].(v10.BeaconService_LatestCrystallizedStateClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LatestCrystallizedState indicates an expected call of LatestCrystallizedState
func (mr *MockBeaconServiceClientMockRecorder) LatestCrystallizedState(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LatestCrystallizedState", reflect.TypeOf((*MockBeaconServiceClient)(nil).LatestCrystallizedState), varargs...)
}

// MockBeaconService_LatestBeaconBlockClient is a mock of BeaconService_LatestBeaconBlockClient interface
type MockBeaconService_LatestBeaconBlockClient struct {
	ctrl     *gomock.Controller
	recorder *MockBeaconService_LatestBeaconBlockClientMockRecorder
}

// MockBeaconService_LatestBeaconBlockClientMockRecorder is the mock recorder for MockBeaconService_LatestBeaconBlockClient
type MockBeaconService_LatestBeaconBlockClientMockRecorder struct {
	mock *MockBeaconService_LatestBeaconBlockClient
}

// NewMockBeaconService_LatestBeaconBlockClient creates a new mock instance
func NewMockBeaconService_LatestBeaconBlockClient(ctrl *gomock.Controller) *MockBeaconService_LatestBeaconBlockClient {
	mock := &MockBeaconService_LatestBeaconBlockClient{ctrl: ctrl}
	mock.recorder = &MockBeaconService_LatestBeaconBlockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBeaconService_LatestBeaconBlockClient) EXPECT() *MockBeaconService_LatestBeaconBlockClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method
func (m *MockBeaconService_LatestBeaconBlockClient) CloseSend() error {
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend
func (mr *MockBeaconService_LatestBeaconBlockClientMockRecorder) CloseSend() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockBeaconService_LatestBeaconBlockClient)(nil).CloseSend))
}

// Context mocks base method
func (m *MockBeaconService_LatestBeaconBlockClient) Context() context.Context {
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockBeaconService_LatestBeaconBlockClientMockRecorder) Context() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockBeaconService_LatestBeaconBlockClient)(nil).Context))
}

// Header mocks base method
func (m *MockBeaconService_LatestBeaconBlockClient) Header() (metadata.MD, error) {
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header
func (mr *MockBeaconService_LatestBeaconBlockClientMockRecorder) Header() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockBeaconService_LatestBeaconBlockClient)(nil).Header))
}

// Recv mocks base method
func (m *MockBeaconService_LatestBeaconBlockClient) Recv() (*v1.BeaconBlock, error) {
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*v1.BeaconBlock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv
func (mr *MockBeaconService_LatestBeaconBlockClientMockRecorder) Recv() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockBeaconService_LatestBeaconBlockClient)(nil).Recv))
}

// RecvMsg mocks base method
func (m *MockBeaconService_LatestBeaconBlockClient) RecvMsg(arg0 interface{}) error {
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg
func (mr *MockBeaconService_LatestBeaconBlockClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockBeaconService_LatestBeaconBlockClient)(nil).RecvMsg), arg0)
}

// SendMsg mocks base method
func (m *MockBeaconService_LatestBeaconBlockClient) SendMsg(arg0 interface{}) error {
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg
func (mr *MockBeaconService_LatestBeaconBlockClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockBeaconService_LatestBeaconBlockClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method
func (m *MockBeaconService_LatestBeaconBlockClient) Trailer() metadata.MD {
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer
func (mr *MockBeaconService_LatestBeaconBlockClientMockRecorder) Trailer() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockBeaconService_LatestBeaconBlockClient)(nil).Trailer))
}

// MockBeaconService_LatestCrystallizedStateClient is a mock of BeaconService_LatestCrystallizedStateClient interface
type MockBeaconService_LatestCrystallizedStateClient struct {
	ctrl     *gomock.Controller
	recorder *MockBeaconService_LatestCrystallizedStateClientMockRecorder
}

// MockBeaconService_LatestCrystallizedStateClientMockRecorder is the mock recorder for MockBeaconService_LatestCrystallizedStateClient
type MockBeaconService_LatestCrystallizedStateClientMockRecorder struct {
	mock *MockBeaconService_LatestCrystallizedStateClient
}

// NewMockBeaconService_LatestCrystallizedStateClient creates a new mock instance
func NewMockBeaconService_LatestCrystallizedStateClient(ctrl *gomock.Controller) *MockBeaconService_LatestCrystallizedStateClient {
	mock := &MockBeaconService_LatestCrystallizedStateClient{ctrl: ctrl}
	mock.recorder = &MockBeaconService_LatestCrystallizedStateClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBeaconService_LatestCrystallizedStateClient) EXPECT() *MockBeaconService_LatestCrystallizedStateClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method
func (m *MockBeaconService_LatestCrystallizedStateClient) CloseSend() error {
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend
func (mr *MockBeaconService_LatestCrystallizedStateClientMockRecorder) CloseSend() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockBeaconService_LatestCrystallizedStateClient)(nil).CloseSend))
}

// Context mocks base method
func (m *MockBeaconService_LatestCrystallizedStateClient) Context() context.Context {
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockBeaconService_LatestCrystallizedStateClientMockRecorder) Context() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockBeaconService_LatestCrystallizedStateClient)(nil).Context))
}

// Header mocks base method
func (m *MockBeaconService_LatestCrystallizedStateClient) Header() (metadata.MD, error) {
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header
func (mr *MockBeaconService_LatestCrystallizedStateClientMockRecorder) Header() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockBeaconService_LatestCrystallizedStateClient)(nil).Header))
}

// Recv mocks base method
func (m *MockBeaconService_LatestCrystallizedStateClient) Recv() (*v1.CrystallizedState, error) {
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*v1.CrystallizedState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv
func (mr *MockBeaconService_LatestCrystallizedStateClientMockRecorder) Recv() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockBeaconService_LatestCrystallizedStateClient)(nil).Recv))
}

// RecvMsg mocks base method
func (m *MockBeaconService_LatestCrystallizedStateClient) RecvMsg(arg0 interface{}) error {
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg
func (mr *MockBeaconService_LatestCrystallizedStateClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockBeaconService_LatestCrystallizedStateClient)(nil).RecvMsg), arg0)
}

// SendMsg mocks base method
func (m *MockBeaconService_LatestCrystallizedStateClient) SendMsg(arg0 interface{}) error {
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg
func (mr *MockBeaconService_LatestCrystallizedStateClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockBeaconService_LatestCrystallizedStateClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method
func (m *MockBeaconService_LatestCrystallizedStateClient) Trailer() metadata.MD {
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer
func (mr *MockBeaconService_LatestCrystallizedStateClientMockRecorder) Trailer() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockBeaconService_LatestCrystallizedStateClient)(nil).Trailer))
}
