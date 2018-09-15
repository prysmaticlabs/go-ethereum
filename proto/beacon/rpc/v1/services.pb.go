// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/beacon/rpc/v1/services.proto

package v1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import empty "github.com/golang/protobuf/ptypes/empty"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"
import v1 "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ShuffleRequest struct {
	CrystallizedStateHash []byte   `protobuf:"bytes,1,opt,name=crystallized_state_hash,json=crystallizedStateHash,proto3" json:"crystallized_state_hash,omitempty"`
	XXX_NoUnkeyedLiteral  struct{} `json:"-"`
	XXX_unrecognized      []byte   `json:"-"`
	XXX_sizecache         int32    `json:"-"`
}

func (m *ShuffleRequest) Reset()         { *m = ShuffleRequest{} }
func (m *ShuffleRequest) String() string { return proto.CompactTextString(m) }
func (*ShuffleRequest) ProtoMessage()    {}
func (*ShuffleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_services_0da4be7efeaeb4c6, []int{0}
}
func (m *ShuffleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShuffleRequest.Unmarshal(m, b)
}
func (m *ShuffleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShuffleRequest.Marshal(b, m, deterministic)
}
func (dst *ShuffleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShuffleRequest.Merge(dst, src)
}
func (m *ShuffleRequest) XXX_Size() int {
	return xxx_messageInfo_ShuffleRequest.Size(m)
}
func (m *ShuffleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ShuffleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ShuffleRequest proto.InternalMessageInfo

func (m *ShuffleRequest) GetCrystallizedStateHash() []byte {
	if m != nil {
		return m.CrystallizedStateHash
	}
	return nil
}

type ShuffleResponse struct {
	ShuffledValidatorIndices []uint64 `protobuf:"varint,1,rep,packed,name=shuffled_validator_indices,json=shuffledValidatorIndices" json:"shuffled_validator_indices,omitempty"`
	CutoffIndices            []uint64 `protobuf:"varint,2,rep,packed,name=cutoff_indices,json=cutoffIndices" json:"cutoff_indices,omitempty"`
	AssignedAttestationSlots []uint64 `protobuf:"varint,3,rep,packed,name=assigned_attestation_slots,json=assignedAttestationSlots" json:"assigned_attestation_slots,omitempty"`
	XXX_NoUnkeyedLiteral     struct{} `json:"-"`
	XXX_unrecognized         []byte   `json:"-"`
	XXX_sizecache            int32    `json:"-"`
}

func (m *ShuffleResponse) Reset()         { *m = ShuffleResponse{} }
func (m *ShuffleResponse) String() string { return proto.CompactTextString(m) }
func (*ShuffleResponse) ProtoMessage()    {}
func (*ShuffleResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_services_0da4be7efeaeb4c6, []int{1}
}
func (m *ShuffleResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShuffleResponse.Unmarshal(m, b)
}
func (m *ShuffleResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShuffleResponse.Marshal(b, m, deterministic)
}
func (dst *ShuffleResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShuffleResponse.Merge(dst, src)
}
func (m *ShuffleResponse) XXX_Size() int {
	return xxx_messageInfo_ShuffleResponse.Size(m)
}
func (m *ShuffleResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ShuffleResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ShuffleResponse proto.InternalMessageInfo

func (m *ShuffleResponse) GetShuffledValidatorIndices() []uint64 {
	if m != nil {
		return m.ShuffledValidatorIndices
	}
	return nil
}

func (m *ShuffleResponse) GetCutoffIndices() []uint64 {
	if m != nil {
		return m.CutoffIndices
	}
	return nil
}

func (m *ShuffleResponse) GetAssignedAttestationSlots() []uint64 {
	if m != nil {
		return m.AssignedAttestationSlots
	}
	return nil
}

type ProposeRequest struct {
	ParentHash              []byte               `protobuf:"bytes,1,opt,name=parent_hash,json=parentHash,proto3" json:"parent_hash,omitempty"`
	SlotNumber              uint64               `protobuf:"varint,2,opt,name=slot_number,json=slotNumber" json:"slot_number,omitempty"`
	RandaoReveal            []byte               `protobuf:"bytes,3,opt,name=randao_reveal,json=randaoReveal,proto3" json:"randao_reveal,omitempty"`
	AttestationBitmask      []byte               `protobuf:"bytes,4,opt,name=attestation_bitmask,json=attestationBitmask,proto3" json:"attestation_bitmask,omitempty"`
	AttestationAggregateSig []uint32             `protobuf:"varint,5,rep,packed,name=attestation_aggregate_sig,json=attestationAggregateSig" json:"attestation_aggregate_sig,omitempty"`
	Timestamp               *timestamp.Timestamp `protobuf:"bytes,6,opt,name=timestamp" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral    struct{}             `json:"-"`
	XXX_unrecognized        []byte               `json:"-"`
	XXX_sizecache           int32                `json:"-"`
}

func (m *ProposeRequest) Reset()         { *m = ProposeRequest{} }
func (m *ProposeRequest) String() string { return proto.CompactTextString(m) }
func (*ProposeRequest) ProtoMessage()    {}
func (*ProposeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_services_0da4be7efeaeb4c6, []int{2}
}
func (m *ProposeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProposeRequest.Unmarshal(m, b)
}
func (m *ProposeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProposeRequest.Marshal(b, m, deterministic)
}
func (dst *ProposeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProposeRequest.Merge(dst, src)
}
func (m *ProposeRequest) XXX_Size() int {
	return xxx_messageInfo_ProposeRequest.Size(m)
}
func (m *ProposeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ProposeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProposeRequest proto.InternalMessageInfo

func (m *ProposeRequest) GetParentHash() []byte {
	if m != nil {
		return m.ParentHash
	}
	return nil
}

func (m *ProposeRequest) GetSlotNumber() uint64 {
	if m != nil {
		return m.SlotNumber
	}
	return 0
}

func (m *ProposeRequest) GetRandaoReveal() []byte {
	if m != nil {
		return m.RandaoReveal
	}
	return nil
}

func (m *ProposeRequest) GetAttestationBitmask() []byte {
	if m != nil {
		return m.AttestationBitmask
	}
	return nil
}

func (m *ProposeRequest) GetAttestationAggregateSig() []uint32 {
	if m != nil {
		return m.AttestationAggregateSig
	}
	return nil
}

func (m *ProposeRequest) GetTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

type ProposeResponse struct {
	BlockHash            []byte   `protobuf:"bytes,1,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProposeResponse) Reset()         { *m = ProposeResponse{} }
func (m *ProposeResponse) String() string { return proto.CompactTextString(m) }
func (*ProposeResponse) ProtoMessage()    {}
func (*ProposeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_services_0da4be7efeaeb4c6, []int{3}
}
func (m *ProposeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProposeResponse.Unmarshal(m, b)
}
func (m *ProposeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProposeResponse.Marshal(b, m, deterministic)
}
func (dst *ProposeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProposeResponse.Merge(dst, src)
}
func (m *ProposeResponse) XXX_Size() int {
	return xxx_messageInfo_ProposeResponse.Size(m)
}
func (m *ProposeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ProposeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ProposeResponse proto.InternalMessageInfo

func (m *ProposeResponse) GetBlockHash() []byte {
	if m != nil {
		return m.BlockHash
	}
	return nil
}

type AttestRequest struct {
	Attestation          *v1.AggregatedAttestation `protobuf:"bytes,1,opt,name=attestation,proto3" json:"attestation,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *AttestRequest) Reset()         { *m = AttestRequest{} }
func (m *AttestRequest) String() string { return proto.CompactTextString(m) }
func (*AttestRequest) ProtoMessage()    {}
func (*AttestRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_services_0da4be7efeaeb4c6, []int{4}
}
func (m *AttestRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AttestRequest.Unmarshal(m, b)
}
func (m *AttestRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AttestRequest.Marshal(b, m, deterministic)
}
func (dst *AttestRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AttestRequest.Merge(dst, src)
}
func (m *AttestRequest) XXX_Size() int {
	return xxx_messageInfo_AttestRequest.Size(m)
}
func (m *AttestRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AttestRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AttestRequest proto.InternalMessageInfo

func (m *AttestRequest) GetAttestation() *v1.AggregatedAttestation {
	if m != nil {
		return m.Attestation
	}
	return nil
}

type AttestResponse struct {
	AttestationHash      []byte   `protobuf:"bytes,1,opt,name=attestation_hash,json=attestationHash,proto3" json:"attestation_hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AttestResponse) Reset()         { *m = AttestResponse{} }
func (m *AttestResponse) String() string { return proto.CompactTextString(m) }
func (*AttestResponse) ProtoMessage()    {}
func (*AttestResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_services_0da4be7efeaeb4c6, []int{5}
}
func (m *AttestResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AttestResponse.Unmarshal(m, b)
}
func (m *AttestResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AttestResponse.Marshal(b, m, deterministic)
}
func (dst *AttestResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AttestResponse.Merge(dst, src)
}
func (m *AttestResponse) XXX_Size() int {
	return xxx_messageInfo_AttestResponse.Size(m)
}
func (m *AttestResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AttestResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AttestResponse proto.InternalMessageInfo

func (m *AttestResponse) GetAttestationHash() []byte {
	if m != nil {
		return m.AttestationHash
	}
	return nil
}

func init() {
	proto.RegisterType((*ShuffleRequest)(nil), "ethereum.beacon.rpc.v1.ShuffleRequest")
	proto.RegisterType((*ShuffleResponse)(nil), "ethereum.beacon.rpc.v1.ShuffleResponse")
	proto.RegisterType((*ProposeRequest)(nil), "ethereum.beacon.rpc.v1.ProposeRequest")
	proto.RegisterType((*ProposeResponse)(nil), "ethereum.beacon.rpc.v1.ProposeResponse")
	proto.RegisterType((*AttestRequest)(nil), "ethereum.beacon.rpc.v1.AttestRequest")
	proto.RegisterType((*AttestResponse)(nil), "ethereum.beacon.rpc.v1.AttestResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for BeaconService service

type BeaconServiceClient interface {
	LatestBeaconBlock(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (BeaconService_LatestBeaconBlockClient, error)
	LatestCrystallizedState(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (BeaconService_LatestCrystallizedStateClient, error)
	FetchShuffledValidatorIndices(ctx context.Context, in *ShuffleRequest, opts ...grpc.CallOption) (*ShuffleResponse, error)
	LatestAttestation(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (BeaconService_LatestAttestationClient, error)
}

type beaconServiceClient struct {
	cc *grpc.ClientConn
}

func NewBeaconServiceClient(cc *grpc.ClientConn) BeaconServiceClient {
	return &beaconServiceClient{cc}
}

func (c *beaconServiceClient) LatestBeaconBlock(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (BeaconService_LatestBeaconBlockClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_BeaconService_serviceDesc.Streams[0], c.cc, "/ethereum.beacon.rpc.v1.BeaconService/LatestBeaconBlock", opts...)
	if err != nil {
		return nil, err
	}
	x := &beaconServiceLatestBeaconBlockClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BeaconService_LatestBeaconBlockClient interface {
	Recv() (*v1.BeaconBlock, error)
	grpc.ClientStream
}

type beaconServiceLatestBeaconBlockClient struct {
	grpc.ClientStream
}

func (x *beaconServiceLatestBeaconBlockClient) Recv() (*v1.BeaconBlock, error) {
	m := new(v1.BeaconBlock)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *beaconServiceClient) LatestCrystallizedState(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (BeaconService_LatestCrystallizedStateClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_BeaconService_serviceDesc.Streams[1], c.cc, "/ethereum.beacon.rpc.v1.BeaconService/LatestCrystallizedState", opts...)
	if err != nil {
		return nil, err
	}
	x := &beaconServiceLatestCrystallizedStateClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BeaconService_LatestCrystallizedStateClient interface {
	Recv() (*v1.CrystallizedState, error)
	grpc.ClientStream
}

type beaconServiceLatestCrystallizedStateClient struct {
	grpc.ClientStream
}

func (x *beaconServiceLatestCrystallizedStateClient) Recv() (*v1.CrystallizedState, error) {
	m := new(v1.CrystallizedState)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *beaconServiceClient) FetchShuffledValidatorIndices(ctx context.Context, in *ShuffleRequest, opts ...grpc.CallOption) (*ShuffleResponse, error) {
	out := new(ShuffleResponse)
	err := grpc.Invoke(ctx, "/ethereum.beacon.rpc.v1.BeaconService/FetchShuffledValidatorIndices", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *beaconServiceClient) LatestAttestation(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (BeaconService_LatestAttestationClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_BeaconService_serviceDesc.Streams[2], c.cc, "/ethereum.beacon.rpc.v1.BeaconService/LatestAttestation", opts...)
	if err != nil {
		return nil, err
	}
	x := &beaconServiceLatestAttestationClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BeaconService_LatestAttestationClient interface {
	Recv() (*v1.AttestationRecord, error)
	grpc.ClientStream
}

type beaconServiceLatestAttestationClient struct {
	grpc.ClientStream
}

func (x *beaconServiceLatestAttestationClient) Recv() (*v1.AttestationRecord, error) {
	m := new(v1.AttestationRecord)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for BeaconService service

type BeaconServiceServer interface {
	LatestBeaconBlock(*empty.Empty, BeaconService_LatestBeaconBlockServer) error
	LatestCrystallizedState(*empty.Empty, BeaconService_LatestCrystallizedStateServer) error
	FetchShuffledValidatorIndices(context.Context, *ShuffleRequest) (*ShuffleResponse, error)
	LatestAttestation(*empty.Empty, BeaconService_LatestAttestationServer) error
}

func RegisterBeaconServiceServer(s *grpc.Server, srv BeaconServiceServer) {
	s.RegisterService(&_BeaconService_serviceDesc, srv)
}

func _BeaconService_LatestBeaconBlock_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(empty.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BeaconServiceServer).LatestBeaconBlock(m, &beaconServiceLatestBeaconBlockServer{stream})
}

type BeaconService_LatestBeaconBlockServer interface {
	Send(*v1.BeaconBlock) error
	grpc.ServerStream
}

type beaconServiceLatestBeaconBlockServer struct {
	grpc.ServerStream
}

func (x *beaconServiceLatestBeaconBlockServer) Send(m *v1.BeaconBlock) error {
	return x.ServerStream.SendMsg(m)
}

func _BeaconService_LatestCrystallizedState_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(empty.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BeaconServiceServer).LatestCrystallizedState(m, &beaconServiceLatestCrystallizedStateServer{stream})
}

type BeaconService_LatestCrystallizedStateServer interface {
	Send(*v1.CrystallizedState) error
	grpc.ServerStream
}

type beaconServiceLatestCrystallizedStateServer struct {
	grpc.ServerStream
}

func (x *beaconServiceLatestCrystallizedStateServer) Send(m *v1.CrystallizedState) error {
	return x.ServerStream.SendMsg(m)
}

func _BeaconService_FetchShuffledValidatorIndices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShuffleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BeaconServiceServer).FetchShuffledValidatorIndices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ethereum.beacon.rpc.v1.BeaconService/FetchShuffledValidatorIndices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BeaconServiceServer).FetchShuffledValidatorIndices(ctx, req.(*ShuffleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BeaconService_LatestAttestation_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(empty.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BeaconServiceServer).LatestAttestation(m, &beaconServiceLatestAttestationServer{stream})
}

type BeaconService_LatestAttestationServer interface {
	Send(*v1.AttestationRecord) error
	grpc.ServerStream
}

type beaconServiceLatestAttestationServer struct {
	grpc.ServerStream
}

func (x *beaconServiceLatestAttestationServer) Send(m *v1.AttestationRecord) error {
	return x.ServerStream.SendMsg(m)
}

var _BeaconService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ethereum.beacon.rpc.v1.BeaconService",
	HandlerType: (*BeaconServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FetchShuffledValidatorIndices",
			Handler:    _BeaconService_FetchShuffledValidatorIndices_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "LatestBeaconBlock",
			Handler:       _BeaconService_LatestBeaconBlock_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "LatestCrystallizedState",
			Handler:       _BeaconService_LatestCrystallizedState_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "LatestAttestation",
			Handler:       _BeaconService_LatestAttestation_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/beacon/rpc/v1/services.proto",
}

// Client API for AttesterService service

type AttesterServiceClient interface {
	AttestHead(ctx context.Context, in *AttestRequest, opts ...grpc.CallOption) (*AttestResponse, error)
}

type attesterServiceClient struct {
	cc *grpc.ClientConn
}

func NewAttesterServiceClient(cc *grpc.ClientConn) AttesterServiceClient {
	return &attesterServiceClient{cc}
}

func (c *attesterServiceClient) AttestHead(ctx context.Context, in *AttestRequest, opts ...grpc.CallOption) (*AttestResponse, error) {
	out := new(AttestResponse)
	err := grpc.Invoke(ctx, "/ethereum.beacon.rpc.v1.AttesterService/AttestHead", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AttesterService service

type AttesterServiceServer interface {
	AttestHead(context.Context, *AttestRequest) (*AttestResponse, error)
}

func RegisterAttesterServiceServer(s *grpc.Server, srv AttesterServiceServer) {
	s.RegisterService(&_AttesterService_serviceDesc, srv)
}

func _AttesterService_AttestHead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AttestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AttesterServiceServer).AttestHead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ethereum.beacon.rpc.v1.AttesterService/AttestHead",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AttesterServiceServer).AttestHead(ctx, req.(*AttestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AttesterService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ethereum.beacon.rpc.v1.AttesterService",
	HandlerType: (*AttesterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AttestHead",
			Handler:    _AttesterService_AttestHead_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/beacon/rpc/v1/services.proto",
}

// Client API for ProposerService service

type ProposerServiceClient interface {
	ProposeBlock(ctx context.Context, in *ProposeRequest, opts ...grpc.CallOption) (*ProposeResponse, error)
}

type proposerServiceClient struct {
	cc *grpc.ClientConn
}

func NewProposerServiceClient(cc *grpc.ClientConn) ProposerServiceClient {
	return &proposerServiceClient{cc}
}

func (c *proposerServiceClient) ProposeBlock(ctx context.Context, in *ProposeRequest, opts ...grpc.CallOption) (*ProposeResponse, error) {
	out := new(ProposeResponse)
	err := grpc.Invoke(ctx, "/ethereum.beacon.rpc.v1.ProposerService/ProposeBlock", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ProposerService service

type ProposerServiceServer interface {
	ProposeBlock(context.Context, *ProposeRequest) (*ProposeResponse, error)
}

func RegisterProposerServiceServer(s *grpc.Server, srv ProposerServiceServer) {
	s.RegisterService(&_ProposerService_serviceDesc, srv)
}

func _ProposerService_ProposeBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProposeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProposerServiceServer).ProposeBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ethereum.beacon.rpc.v1.ProposerService/ProposeBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProposerServiceServer).ProposeBlock(ctx, req.(*ProposeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProposerService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ethereum.beacon.rpc.v1.ProposerService",
	HandlerType: (*ProposerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProposeBlock",
			Handler:    _ProposerService_ProposeBlock_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/beacon/rpc/v1/services.proto",
}

func init() {
	proto.RegisterFile("proto/beacon/rpc/v1/services.proto", fileDescriptor_services_0da4be7efeaeb4c6)
}

var fileDescriptor_services_0da4be7efeaeb4c6 = []byte{
	// 652 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xd1, 0x6e, 0xd3, 0x4a,
	0x10, 0x95, 0x9b, 0xde, 0x4a, 0x9d, 0x34, 0xc9, 0x65, 0x11, 0xad, 0x31, 0xaa, 0x1a, 0xa5, 0x6a,
	0x49, 0x5f, 0xec, 0xd4, 0x48, 0x08, 0x41, 0x5f, 0x5a, 0x04, 0x2a, 0x02, 0x21, 0xe4, 0x20, 0x84,
	0x04, 0xc2, 0xda, 0xd8, 0x1b, 0xc7, 0xaa, 0xed, 0x5d, 0x76, 0x37, 0x91, 0xca, 0x4f, 0xf0, 0x07,
	0x7c, 0x05, 0x1f, 0x88, 0xd6, 0x1b, 0xbb, 0x9b, 0x14, 0xab, 0xe2, 0xd1, 0xe7, 0x9c, 0x99, 0x9d,
	0x39, 0x33, 0x63, 0x18, 0x30, 0x4e, 0x25, 0xf5, 0x26, 0x04, 0x47, 0xb4, 0xf0, 0x38, 0x8b, 0xbc,
	0xc5, 0xa9, 0x27, 0x08, 0x5f, 0xa4, 0x11, 0x11, 0x6e, 0x49, 0xa2, 0x5d, 0x22, 0x67, 0x84, 0x93,
	0x79, 0xee, 0x6a, 0x99, 0xcb, 0x59, 0xe4, 0x2e, 0x4e, 0x9d, 0xd5, 0x58, 0xe6, 0x33, 0x15, 0x9b,
	0x13, 0x21, 0x70, 0x52, 0xc5, 0x3a, 0x8f, 0x12, 0x4a, 0x93, 0x8c, 0x78, 0xe5, 0xd7, 0x64, 0x3e,
	0xf5, 0x48, 0xce, 0xe4, 0xf5, 0x92, 0x3c, 0x58, 0x27, 0x65, 0x9a, 0x13, 0x21, 0x71, 0xce, 0xb4,
	0x60, 0x70, 0x09, 0xdd, 0xf1, 0x6c, 0x3e, 0x9d, 0x66, 0x24, 0x20, 0xdf, 0xe7, 0x44, 0x48, 0xf4,
	0x14, 0xf6, 0x22, 0x7e, 0x2d, 0x24, 0xce, 0xb2, 0xf4, 0x07, 0x89, 0x43, 0x21, 0xb1, 0x24, 0xe1,
	0x0c, 0x8b, 0x99, 0x6d, 0xf5, 0xad, 0xe1, 0x4e, 0xf0, 0xc0, 0xa4, 0xc7, 0x8a, 0xbd, 0xc4, 0x62,
	0x36, 0xf8, 0x6d, 0x41, 0xaf, 0x4e, 0x25, 0x18, 0x2d, 0x04, 0x41, 0x67, 0xe0, 0x08, 0x0d, 0xc5,
	0xe1, 0x02, 0x67, 0x69, 0x8c, 0x25, 0xe5, 0x61, 0x5a, 0xc4, 0xaa, 0x77, 0xdb, 0xea, 0xb7, 0x86,
	0x9b, 0x81, 0x5d, 0x29, 0x3e, 0x55, 0x82, 0x37, 0x9a, 0x47, 0x47, 0xd0, 0x8d, 0xe6, 0x92, 0x4e,
	0xa7, 0x75, 0xc4, 0x46, 0x19, 0xd1, 0xd1, 0x68, 0x25, 0x3b, 0x03, 0x07, 0x0b, 0x91, 0x26, 0x05,
	0x89, 0x43, 0x2c, 0xa5, 0x6a, 0x4f, 0xa6, 0xb4, 0x08, 0x45, 0x46, 0xa5, 0xb0, 0x5b, 0xfa, 0x91,
	0x4a, 0x71, 0x7e, 0x23, 0x18, 0x2b, 0x7e, 0xf0, 0x6b, 0x03, 0xba, 0x1f, 0x38, 0x65, 0x54, 0xd4,
	0x0e, 0x1c, 0x40, 0x9b, 0x61, 0x4e, 0x0a, 0x69, 0x76, 0x0d, 0x1a, 0x52, 0xad, 0x2a, 0x81, 0x4a,
	0x1e, 0x16, 0xf3, 0x7c, 0x42, 0xb8, 0xbd, 0xd1, 0xb7, 0x86, 0x9b, 0x01, 0x28, 0xe8, 0x7d, 0x89,
	0xa0, 0x43, 0xe8, 0x70, 0x5c, 0xc4, 0x98, 0x86, 0x9c, 0x2c, 0x08, 0xce, 0xec, 0x56, 0x99, 0x63,
	0x47, 0x83, 0x41, 0x89, 0x21, 0x0f, 0xee, 0x9b, 0xe5, 0x4e, 0x52, 0x99, 0x63, 0x71, 0x65, 0x6f,
	0x96, 0x52, 0x64, 0x50, 0x17, 0x9a, 0x41, 0xcf, 0xe1, 0xa1, 0x19, 0x80, 0x93, 0x84, 0x93, 0x44,
	0x0d, 0x47, 0xa4, 0x89, 0xfd, 0x5f, 0xbf, 0x35, 0xec, 0x04, 0x7b, 0x86, 0xe0, 0xbc, 0xe2, 0xc7,
	0x69, 0x82, 0x9e, 0xc1, 0x76, 0x3d, 0x7a, 0x7b, 0xab, 0x6f, 0x0d, 0xdb, 0xbe, 0xe3, 0xea, 0xe5,
	0x70, 0xab, 0xe5, 0x70, 0x3f, 0x56, 0x8a, 0xe0, 0x46, 0x3c, 0x18, 0x41, 0xaf, 0xf6, 0x67, 0x39,
	0xd6, 0x7d, 0x80, 0x49, 0x46, 0xa3, 0x2b, 0xd3, 0x9f, 0xed, 0x12, 0x29, 0x37, 0xe1, 0x2b, 0x74,
	0xb4, 0xcd, 0x95, 0xa1, 0x6f, 0xa1, 0x6d, 0xd4, 0x55, 0x06, 0xb4, 0xfd, 0x13, 0x77, 0x7d, 0xe9,
	0x99, 0xcf, 0xdc, 0xc5, 0xa9, 0x6b, 0x8c, 0x28, 0x20, 0x11, 0xe5, 0x71, 0x60, 0x46, 0x0f, 0x5e,
	0x40, 0xb7, 0xca, 0xbe, 0x2c, 0xe7, 0x04, 0xfe, 0x37, 0x7d, 0x31, 0x8a, 0xea, 0x19, 0xb8, 0x2a,
	0xcd, 0xff, 0xd9, 0x82, 0xce, 0x45, 0xf9, 0xda, 0x58, 0x5f, 0x20, 0x0a, 0xe0, 0xde, 0x3b, 0xac,
	0x44, 0x1a, 0xbe, 0x50, 0x5d, 0xa0, 0xdd, 0x5b, 0xd6, 0xbc, 0x52, 0x47, 0xe5, 0x1c, 0x36, 0xd5,
	0x6c, 0x04, 0x8f, 0x2c, 0xf4, 0x0d, 0xf6, 0x74, 0xce, 0x97, 0xeb, 0x97, 0xd2, 0x98, 0xb9, 0xd1,
	0x8d, 0x5b, 0x29, 0x46, 0x16, 0x62, 0xb0, 0xff, 0x9a, 0xc8, 0x68, 0x36, 0x6e, 0xba, 0x9c, 0x63,
	0xf7, 0xef, 0x3f, 0x14, 0x77, 0xf5, 0xd6, 0x9d, 0xc7, 0x77, 0xea, 0x96, 0x16, 0x7f, 0xae, 0x5c,
	0x32, 0x86, 0xf3, 0xef, 0xbd, 0xdc, 0x9a, 0xec, 0xc8, 0xf2, 0x0b, 0xe8, 0x69, 0x98, 0xf0, 0x6a,
	0x24, 0x5f, 0x00, 0x34, 0x74, 0x49, 0x70, 0x8c, 0x8e, 0x9a, 0x6a, 0x5c, 0xd9, 0x31, 0xe7, 0xf8,
	0x2e, 0x99, 0xee, 0xc4, 0xe7, 0xf5, 0x3a, 0xd7, 0xef, 0x85, 0xb0, 0xb3, 0x84, 0xf4, 0xf4, 0x1b,
	0x53, 0xad, 0xfe, 0x27, 0x9a, 0xdd, 0x5b, 0xbb, 0x97, 0xc9, 0x56, 0x69, 0xd0, 0x93, 0x3f, 0x01,
	0x00, 0x00, 0xff, 0xff, 0x6d, 0xd3, 0xb2, 0x63, 0x0b, 0x06, 0x00, 0x00,
}
