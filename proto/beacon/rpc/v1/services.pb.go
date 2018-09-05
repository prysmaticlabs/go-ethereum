// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/beacon/rpc/v1/services.proto

package ethereum_beacon_rpc_v1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import empty "github.com/golang/protobuf/ptypes/empty"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"
import v1 "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
import v11 "github.com/prysmaticlabs/prysm/proto/sharding/p2p/v1"

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

type GenesisTimeAndStateResponse struct {
	GenesisTimestamp        *timestamp.Timestamp  `protobuf:"bytes,1,opt,name=genesis_timestamp,json=genesisTimestamp,proto3" json:"genesis_timestamp,omitempty"`
	LatestCrystallizedState *v1.CrystallizedState `protobuf:"bytes,2,opt,name=latest_crystallized_state,json=latestCrystallizedState,proto3" json:"latest_crystallized_state,omitempty"`
	XXX_NoUnkeyedLiteral    struct{}              `json:"-"`
	XXX_unrecognized        []byte                `json:"-"`
	XXX_sizecache           int32                 `json:"-"`
}

func (m *GenesisTimeAndStateResponse) Reset()         { *m = GenesisTimeAndStateResponse{} }
func (m *GenesisTimeAndStateResponse) String() string { return proto.CompactTextString(m) }
func (*GenesisTimeAndStateResponse) ProtoMessage()    {}
func (*GenesisTimeAndStateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_services_81edfbdcc1a1e9b3, []int{0}
}
func (m *GenesisTimeAndStateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenesisTimeAndStateResponse.Unmarshal(m, b)
}
func (m *GenesisTimeAndStateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenesisTimeAndStateResponse.Marshal(b, m, deterministic)
}
func (dst *GenesisTimeAndStateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisTimeAndStateResponse.Merge(dst, src)
}
func (m *GenesisTimeAndStateResponse) XXX_Size() int {
	return xxx_messageInfo_GenesisTimeAndStateResponse.Size(m)
}
func (m *GenesisTimeAndStateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisTimeAndStateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisTimeAndStateResponse proto.InternalMessageInfo

func (m *GenesisTimeAndStateResponse) GetGenesisTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.GenesisTimestamp
	}
	return nil
}

func (m *GenesisTimeAndStateResponse) GetLatestCrystallizedState() *v1.CrystallizedState {
	if m != nil {
		return m.LatestCrystallizedState
	}
	return nil
}

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
	return fileDescriptor_services_81edfbdcc1a1e9b3, []int{1}
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
	ShuffledValidatorIndices []uint64 `protobuf:"varint,1,rep,packed,name=shuffled_validator_indices,json=shuffledValidatorIndices,proto3" json:"shuffled_validator_indices,omitempty"`
	CutoffIndices            []uint64 `protobuf:"varint,2,rep,packed,name=cutoff_indices,json=cutoffIndices,proto3" json:"cutoff_indices,omitempty"`
	AssignedAttestationSlots []uint64 `protobuf:"varint,3,rep,packed,name=assigned_attestation_slots,json=assignedAttestationSlots,proto3" json:"assigned_attestation_slots,omitempty"`
	XXX_NoUnkeyedLiteral     struct{} `json:"-"`
	XXX_unrecognized         []byte   `json:"-"`
	XXX_sizecache            int32    `json:"-"`
}

func (m *ShuffleResponse) Reset()         { *m = ShuffleResponse{} }
func (m *ShuffleResponse) String() string { return proto.CompactTextString(m) }
func (*ShuffleResponse) ProtoMessage()    {}
func (*ShuffleResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_services_81edfbdcc1a1e9b3, []int{2}
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
	SlotNumber              uint64               `protobuf:"varint,2,opt,name=slot_number,json=slotNumber,proto3" json:"slot_number,omitempty"`
	RandaoReveal            []byte               `protobuf:"bytes,3,opt,name=randao_reveal,json=randaoReveal,proto3" json:"randao_reveal,omitempty"`
	AttestationBitmask      []byte               `protobuf:"bytes,4,opt,name=attestation_bitmask,json=attestationBitmask,proto3" json:"attestation_bitmask,omitempty"`
	AttestationAggregateSig []uint32             `protobuf:"varint,5,rep,packed,name=attestation_aggregate_sig,json=attestationAggregateSig,proto3" json:"attestation_aggregate_sig,omitempty"`
	Timestamp               *timestamp.Timestamp `protobuf:"bytes,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral    struct{}             `json:"-"`
	XXX_unrecognized        []byte               `json:"-"`
	XXX_sizecache           int32                `json:"-"`
}

func (m *ProposeRequest) Reset()         { *m = ProposeRequest{} }
func (m *ProposeRequest) String() string { return proto.CompactTextString(m) }
func (*ProposeRequest) ProtoMessage()    {}
func (*ProposeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_services_81edfbdcc1a1e9b3, []int{3}
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
	return fileDescriptor_services_81edfbdcc1a1e9b3, []int{4}
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

type SignRequest struct {
	BlockHash            []byte         `protobuf:"bytes,1,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`
	Signature            *v11.Signature `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *SignRequest) Reset()         { *m = SignRequest{} }
func (m *SignRequest) String() string { return proto.CompactTextString(m) }
func (*SignRequest) ProtoMessage()    {}
func (*SignRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_services_81edfbdcc1a1e9b3, []int{5}
}
func (m *SignRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignRequest.Unmarshal(m, b)
}
func (m *SignRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignRequest.Marshal(b, m, deterministic)
}
func (dst *SignRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignRequest.Merge(dst, src)
}
func (m *SignRequest) XXX_Size() int {
	return xxx_messageInfo_SignRequest.Size(m)
}
func (m *SignRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SignRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SignRequest proto.InternalMessageInfo

func (m *SignRequest) GetBlockHash() []byte {
	if m != nil {
		return m.BlockHash
	}
	return nil
}

func (m *SignRequest) GetSignature() *v11.Signature {
	if m != nil {
		return m.Signature
	}
	return nil
}

type SignResponse struct {
	Signed               bool     `protobuf:"varint,1,opt,name=signed,proto3" json:"signed,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignResponse) Reset()         { *m = SignResponse{} }
func (m *SignResponse) String() string { return proto.CompactTextString(m) }
func (*SignResponse) ProtoMessage()    {}
func (*SignResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_services_81edfbdcc1a1e9b3, []int{6}
}
func (m *SignResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignResponse.Unmarshal(m, b)
}
func (m *SignResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignResponse.Marshal(b, m, deterministic)
}
func (dst *SignResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignResponse.Merge(dst, src)
}
func (m *SignResponse) XXX_Size() int {
	return xxx_messageInfo_SignResponse.Size(m)
}
func (m *SignResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SignResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SignResponse proto.InternalMessageInfo

func (m *SignResponse) GetSigned() bool {
	if m != nil {
		return m.Signed
	}
	return false
}

func init() {
	proto.RegisterType((*GenesisTimeAndStateResponse)(nil), "ethereum.beacon.rpc.v1.GenesisTimeAndStateResponse")
	proto.RegisterType((*ShuffleRequest)(nil), "ethereum.beacon.rpc.v1.ShuffleRequest")
	proto.RegisterType((*ShuffleResponse)(nil), "ethereum.beacon.rpc.v1.ShuffleResponse")
	proto.RegisterType((*ProposeRequest)(nil), "ethereum.beacon.rpc.v1.ProposeRequest")
	proto.RegisterType((*ProposeResponse)(nil), "ethereum.beacon.rpc.v1.ProposeResponse")
	proto.RegisterType((*SignRequest)(nil), "ethereum.beacon.rpc.v1.SignRequest")
	proto.RegisterType((*SignResponse)(nil), "ethereum.beacon.rpc.v1.SignResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// BeaconServiceClient is the client API for BeaconService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BeaconServiceClient interface {
	GenesisTimeAndCanonicalState(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GenesisTimeAndStateResponse, error)
	LatestBeaconBlock(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (BeaconService_LatestBeaconBlockClient, error)
	LatestCrystallizedState(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (BeaconService_LatestCrystallizedStateClient, error)
	FetchShuffledValidatorIndices(ctx context.Context, in *ShuffleRequest, opts ...grpc.CallOption) (*ShuffleResponse, error)
}

type beaconServiceClient struct {
	cc *grpc.ClientConn
}

func NewBeaconServiceClient(cc *grpc.ClientConn) BeaconServiceClient {
	return &beaconServiceClient{cc}
}

func (c *beaconServiceClient) GenesisTimeAndCanonicalState(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GenesisTimeAndStateResponse, error) {
	out := new(GenesisTimeAndStateResponse)
	err := c.cc.Invoke(ctx, "/ethereum.beacon.rpc.v1.BeaconService/GenesisTimeAndCanonicalState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *beaconServiceClient) LatestBeaconBlock(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (BeaconService_LatestBeaconBlockClient, error) {
	stream, err := c.cc.NewStream(ctx, &_BeaconService_serviceDesc.Streams[0], "/ethereum.beacon.rpc.v1.BeaconService/LatestBeaconBlock", opts...)
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
	stream, err := c.cc.NewStream(ctx, &_BeaconService_serviceDesc.Streams[1], "/ethereum.beacon.rpc.v1.BeaconService/LatestCrystallizedState", opts...)
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
	err := c.cc.Invoke(ctx, "/ethereum.beacon.rpc.v1.BeaconService/FetchShuffledValidatorIndices", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BeaconServiceServer is the server API for BeaconService service.
type BeaconServiceServer interface {
	GenesisTimeAndCanonicalState(context.Context, *empty.Empty) (*GenesisTimeAndStateResponse, error)
	LatestBeaconBlock(*empty.Empty, BeaconService_LatestBeaconBlockServer) error
	LatestCrystallizedState(*empty.Empty, BeaconService_LatestCrystallizedStateServer) error
	FetchShuffledValidatorIndices(context.Context, *ShuffleRequest) (*ShuffleResponse, error)
}

func RegisterBeaconServiceServer(s *grpc.Server, srv BeaconServiceServer) {
	s.RegisterService(&_BeaconService_serviceDesc, srv)
}

func _BeaconService_GenesisTimeAndCanonicalState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BeaconServiceServer).GenesisTimeAndCanonicalState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ethereum.beacon.rpc.v1.BeaconService/GenesisTimeAndCanonicalState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BeaconServiceServer).GenesisTimeAndCanonicalState(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
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

var _BeaconService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ethereum.beacon.rpc.v1.BeaconService",
	HandlerType: (*BeaconServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GenesisTimeAndCanonicalState",
			Handler:    _BeaconService_GenesisTimeAndCanonicalState_Handler,
		},
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
	},
	Metadata: "proto/beacon/rpc/v1/services.proto",
}

// AttesterServiceClient is the client API for AttesterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AttesterServiceClient interface {
	SignBlock(ctx context.Context, in *SignRequest, opts ...grpc.CallOption) (*SignResponse, error)
}

type attesterServiceClient struct {
	cc *grpc.ClientConn
}

func NewAttesterServiceClient(cc *grpc.ClientConn) AttesterServiceClient {
	return &attesterServiceClient{cc}
}

func (c *attesterServiceClient) SignBlock(ctx context.Context, in *SignRequest, opts ...grpc.CallOption) (*SignResponse, error) {
	out := new(SignResponse)
	err := c.cc.Invoke(ctx, "/ethereum.beacon.rpc.v1.AttesterService/SignBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AttesterServiceServer is the server API for AttesterService service.
type AttesterServiceServer interface {
	SignBlock(context.Context, *SignRequest) (*SignResponse, error)
}

func RegisterAttesterServiceServer(s *grpc.Server, srv AttesterServiceServer) {
	s.RegisterService(&_AttesterService_serviceDesc, srv)
}

func _AttesterService_SignBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AttesterServiceServer).SignBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ethereum.beacon.rpc.v1.AttesterService/SignBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AttesterServiceServer).SignBlock(ctx, req.(*SignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AttesterService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ethereum.beacon.rpc.v1.AttesterService",
	HandlerType: (*AttesterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SignBlock",
			Handler:    _AttesterService_SignBlock_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/beacon/rpc/v1/services.proto",
}

// ProposerServiceClient is the client API for ProposerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
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
	err := c.cc.Invoke(ctx, "/ethereum.beacon.rpc.v1.ProposerService/ProposeBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProposerServiceServer is the server API for ProposerService service.
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
	proto.RegisterFile("proto/beacon/rpc/v1/services.proto", fileDescriptor_services_81edfbdcc1a1e9b3)
}

var fileDescriptor_services_81edfbdcc1a1e9b3 = []byte{
	// 733 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0x5f, 0x4f, 0xdb, 0x48,
	0x10, 0x57, 0x08, 0x87, 0x2e, 0x43, 0x02, 0x87, 0x4f, 0x47, 0x8c, 0x39, 0x04, 0x4a, 0x38, 0x8e,
	0xbe, 0xd8, 0x10, 0xa4, 0xaa, 0xaa, 0x78, 0x09, 0xa8, 0x85, 0x4a, 0xa8, 0xaa, 0x6c, 0xc4, 0x63,
	0xad, 0x8d, 0x3d, 0x71, 0x56, 0xd8, 0x5e, 0x77, 0x77, 0x13, 0x89, 0x7e, 0x98, 0x3e, 0xf4, 0x33,
	0xf4, 0x6b, 0xf4, 0x3b, 0x55, 0xde, 0x8d, 0x13, 0x27, 0x60, 0xe8, 0xe3, 0xce, 0xfc, 0xe6, 0x37,
	0x3b, 0xbf, 0xf9, 0x03, 0x9d, 0x8c, 0x33, 0xc9, 0x9c, 0x01, 0x92, 0x80, 0xa5, 0x0e, 0xcf, 0x02,
	0x67, 0x72, 0xea, 0x08, 0xe4, 0x13, 0x1a, 0xa0, 0xb0, 0x95, 0xd3, 0xd8, 0x46, 0x39, 0x42, 0x8e,
	0xe3, 0xc4, 0xd6, 0x30, 0x9b, 0x67, 0x81, 0x3d, 0x39, 0xb5, 0x0e, 0x75, 0xac, 0x18, 0x11, 0x1e,
	0xd2, 0x34, 0x72, 0xb2, 0x5e, 0x96, 0x47, 0x27, 0x28, 0x04, 0x89, 0x8a, 0x68, 0x6b, 0x31, 0xc3,
	0xd3, 0x98, 0xdd, 0x88, 0xb1, 0x28, 0x46, 0x47, 0xbd, 0x06, 0xe3, 0xa1, 0x83, 0x49, 0x26, 0x1f,
	0xa6, 0xce, 0xfd, 0x65, 0xa7, 0xa4, 0x09, 0x0a, 0x49, 0x92, 0x4c, 0x03, 0x3a, 0x3f, 0x6b, 0xb0,
	0x7b, 0x85, 0x29, 0x0a, 0x2a, 0x6e, 0x69, 0x82, 0xfd, 0x34, 0xf4, 0x24, 0x91, 0xe8, 0xa2, 0xc8,
	0x58, 0x2a, 0xd0, 0xb8, 0x82, 0xad, 0x48, 0xbb, 0xfd, 0x59, 0xa8, 0x59, 0x3b, 0xa8, 0x1d, 0xaf,
	0xf7, 0x2c, 0x5b, 0x93, 0xdb, 0x05, 0xb9, 0x7d, 0x5b, 0x20, 0xdc, 0xbf, 0xa2, 0x39, 0xa7, 0xb2,
	0x18, 0x08, 0x3b, 0x31, 0x91, 0x28, 0xa4, 0x1f, 0xf0, 0x07, 0x21, 0x49, 0x1c, 0xd3, 0xaf, 0x18,
	0xfa, 0x22, 0xcf, 0x66, 0xae, 0x28, 0xc2, 0x57, 0xf6, 0xb2, 0x58, 0x59, 0x2f, 0xb3, 0x27, 0xa7,
	0xf6, 0x65, 0x29, 0x42, 0x7f, 0xaf, 0xad, 0xb9, 0x1e, 0x39, 0x3a, 0xd7, 0xb0, 0xe1, 0x8d, 0xc6,
	0xc3, 0x61, 0x8c, 0x2e, 0x7e, 0x19, 0xa3, 0x90, 0xc6, 0x6b, 0x68, 0x3f, 0xce, 0xe8, 0x8f, 0x88,
	0x18, 0xa9, 0x3a, 0x9a, 0xee, 0x3f, 0xc1, 0x32, 0xcb, 0x35, 0x11, 0xa3, 0xce, 0x8f, 0x1a, 0x6c,
	0xce, 0xa8, 0xa6, 0x6a, 0x9c, 0x83, 0x25, 0xb4, 0x29, 0xf4, 0x27, 0x24, 0xa6, 0x21, 0x91, 0x8c,
	0xfb, 0x34, 0x0d, 0xf3, 0x8e, 0x9b, 0xb5, 0x83, 0xfa, 0xf1, 0xaa, 0x6b, 0x16, 0x88, 0xbb, 0x02,
	0xf0, 0x41, 0xfb, 0x8d, 0xff, 0x60, 0x23, 0x18, 0x4b, 0x36, 0x1c, 0xce, 0x22, 0x56, 0x54, 0x44,
	0x4b, 0x5b, 0x0b, 0xd8, 0x39, 0x58, 0x44, 0x08, 0x1a, 0xa5, 0x18, 0xfa, 0x44, 0xe6, 0x65, 0x12,
	0x49, 0x59, 0xea, 0x8b, 0x98, 0x49, 0x61, 0xd6, 0x75, 0x92, 0x02, 0xd1, 0x9f, 0x03, 0xbc, 0xdc,
	0xdf, 0xf9, 0xb6, 0x02, 0x1b, 0x9f, 0x38, 0xcb, 0x98, 0x98, 0x29, 0xb0, 0x0f, 0xeb, 0x19, 0xe1,
	0x98, 0xca, 0x72, 0xd5, 0xa0, 0x4d, 0x79, 0xa9, 0x39, 0x20, 0x27, 0xf7, 0xd3, 0x71, 0x32, 0x40,
	0xae, 0xba, 0xb1, 0xea, 0x42, 0x6e, 0xfa, 0xa8, 0x2c, 0x46, 0x17, 0x5a, 0x9c, 0xa4, 0x21, 0x61,
	0x3e, 0xc7, 0x09, 0x92, 0xd8, 0xac, 0x2b, 0x8e, 0xa6, 0x36, 0xba, 0xca, 0x66, 0x38, 0xf0, 0x77,
	0xf9, 0xbb, 0x03, 0x2a, 0x13, 0x22, 0xee, 0xcd, 0x55, 0x05, 0x35, 0x4a, 0xae, 0x0b, 0xed, 0x31,
	0xde, 0xc2, 0x4e, 0x39, 0x80, 0x44, 0x11, 0xc7, 0x28, 0x6f, 0x8e, 0xa0, 0x91, 0xf9, 0xc7, 0x41,
	0xfd, 0xb8, 0xe5, 0xb6, 0x4b, 0x80, 0x7e, 0xe1, 0xf7, 0x68, 0x64, 0xbc, 0x81, 0xc6, 0x7c, 0x1e,
	0xd7, 0x5e, 0x9c, 0xc7, 0x39, 0xb8, 0x73, 0x02, 0x9b, 0x33, 0x7d, 0xa6, 0x6d, 0xdd, 0x03, 0x18,
	0xc4, 0x2c, 0xb8, 0x2f, 0xeb, 0xd3, 0x50, 0x16, 0x35, 0x09, 0x0c, 0xd6, 0x3d, 0x1a, 0xa5, 0x85,
	0x9c, 0xcf, 0xa3, 0x8d, 0x3e, 0x34, 0xf2, 0xd6, 0x10, 0x39, 0xe6, 0xc5, 0x60, 0x77, 0xe7, 0x83,
	0x5d, 0x2c, 0x7c, 0x31, 0xda, 0x5e, 0x01, 0x75, 0xe7, 0x51, 0x9d, 0x23, 0x68, 0xea, 0x84, 0xd3,
	0xff, 0x6d, 0xc3, 0x9a, 0xee, 0xb6, 0xca, 0xf6, 0xa7, 0x3b, 0x7d, 0xf5, 0xbe, 0xd7, 0xa1, 0x75,
	0xa1, 0x36, 0xc5, 0xd3, 0x57, 0xc7, 0xb8, 0x87, 0x7f, 0x17, 0xb7, 0xf9, 0x92, 0xa4, 0x2c, 0xa5,
	0x01, 0x89, 0xd5, 0x60, 0x1b, 0xdb, 0x8f, 0x34, 0x7a, 0x97, 0x5f, 0x0b, 0xeb, 0xcc, 0x7e, 0xfa,
	0x4e, 0xd9, 0xcf, 0xdd, 0x06, 0x17, 0xb6, 0x6e, 0xd4, 0x1a, 0xea, 0x3f, 0x5c, 0xe4, 0x12, 0x54,
	0x66, 0xe8, 0x56, 0x2d, 0x77, 0x29, 0xf8, 0xa4, 0x66, 0x7c, 0x86, 0xf6, 0xcd, 0xd3, 0xab, 0x5d,
	0xc9, 0xfc, 0xfb, 0x67, 0xe3, 0xa4, 0x66, 0x64, 0xb0, 0xf7, 0x1e, 0x65, 0x30, 0xf2, 0xaa, 0x96,
	0xf4, 0xa8, 0x4a, 0x89, 0xc5, 0xb3, 0x62, 0xfd, 0xff, 0x22, 0x4e, 0xab, 0xd4, 0xa3, 0xb0, 0xa9,
	0x97, 0x14, 0x79, 0xd1, 0xa5, 0x3b, 0x68, 0xe4, 0xfd, 0xd5, 0x82, 0x75, 0x2b, 0x89, 0xe6, 0x33,
	0x67, 0x1d, 0x3e, 0x0f, 0x9a, 0xa6, 0xe2, 0xb3, 0xd1, 0x9e, 0xa5, 0xf2, 0xa1, 0x39, 0x35, 0xe9,
	0x6c, 0x95, 0xe5, 0x2d, 0xde, 0x8c, 0xea, 0xf2, 0x96, 0x76, 0x67, 0xb0, 0xa6, 0xba, 0x71, 0xf6,
	0x2b, 0x00, 0x00, 0xff, 0xff, 0x3f, 0x51, 0x6b, 0x8d, 0x0d, 0x07, 0x00, 0x00,
}
