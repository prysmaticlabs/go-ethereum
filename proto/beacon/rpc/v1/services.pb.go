// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/beacon/rpc/v1/services.proto

package ethereum_beacon_rpc_v1

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
	return fileDescriptor_services_f2b7e4f7fe6852c7, []int{0}
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
	return fileDescriptor_services_f2b7e4f7fe6852c7, []int{1}
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
	return fileDescriptor_services_f2b7e4f7fe6852c7, []int{2}
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
	return fileDescriptor_services_f2b7e4f7fe6852c7, []int{3}
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
	return fileDescriptor_services_f2b7e4f7fe6852c7, []int{4}
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
	return fileDescriptor_services_f2b7e4f7fe6852c7, []int{5}
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

type PublicKey struct {
	PublicKey            uint64   `protobuf:"varint,1,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PublicKey) Reset()         { *m = PublicKey{} }
func (m *PublicKey) String() string { return proto.CompactTextString(m) }
func (*PublicKey) ProtoMessage()    {}
func (*PublicKey) Descriptor() ([]byte, []int) {
	return fileDescriptor_services_f2b7e4f7fe6852c7, []int{6}
}
func (m *PublicKey) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PublicKey.Unmarshal(m, b)
}
func (m *PublicKey) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PublicKey.Marshal(b, m, deterministic)
}
func (dst *PublicKey) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PublicKey.Merge(dst, src)
}
func (m *PublicKey) XXX_Size() int {
	return xxx_messageInfo_PublicKey.Size(m)
}
func (m *PublicKey) XXX_DiscardUnknown() {
	xxx_messageInfo_PublicKey.DiscardUnknown(m)
}

var xxx_messageInfo_PublicKey proto.InternalMessageInfo

func (m *PublicKey) GetPublicKey() uint64 {
	if m != nil {
		return m.PublicKey
	}
	return 0
}

type ShardIDResponse struct {
	ShardId              uint64   `protobuf:"varint,1,opt,name=shard_id,json=shardId,proto3" json:"shard_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShardIDResponse) Reset()         { *m = ShardIDResponse{} }
func (m *ShardIDResponse) String() string { return proto.CompactTextString(m) }
func (*ShardIDResponse) ProtoMessage()    {}
func (*ShardIDResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_services_f2b7e4f7fe6852c7, []int{7}
}
func (m *ShardIDResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShardIDResponse.Unmarshal(m, b)
}
func (m *ShardIDResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShardIDResponse.Marshal(b, m, deterministic)
}
func (dst *ShardIDResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShardIDResponse.Merge(dst, src)
}
func (m *ShardIDResponse) XXX_Size() int {
	return xxx_messageInfo_ShardIDResponse.Size(m)
}
func (m *ShardIDResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ShardIDResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ShardIDResponse proto.InternalMessageInfo

func (m *ShardIDResponse) GetShardId() uint64 {
	if m != nil {
		return m.ShardId
	}
	return 0
}

type IndexResponse struct {
	Index                uint32   `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IndexResponse) Reset()         { *m = IndexResponse{} }
func (m *IndexResponse) String() string { return proto.CompactTextString(m) }
func (*IndexResponse) ProtoMessage()    {}
func (*IndexResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_services_f2b7e4f7fe6852c7, []int{8}
}
func (m *IndexResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IndexResponse.Unmarshal(m, b)
}
func (m *IndexResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IndexResponse.Marshal(b, m, deterministic)
}
func (dst *IndexResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IndexResponse.Merge(dst, src)
}
func (m *IndexResponse) XXX_Size() int {
	return xxx_messageInfo_IndexResponse.Size(m)
}
func (m *IndexResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_IndexResponse.DiscardUnknown(m)
}

var xxx_messageInfo_IndexResponse proto.InternalMessageInfo

func (m *IndexResponse) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

type SlotResponse struct {
	Slot                 uint64   `protobuf:"varint,1,opt,name=slot,proto3" json:"slot,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SlotResponse) Reset()         { *m = SlotResponse{} }
func (m *SlotResponse) String() string { return proto.CompactTextString(m) }
func (*SlotResponse) ProtoMessage()    {}
func (*SlotResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_services_f2b7e4f7fe6852c7, []int{9}
}
func (m *SlotResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SlotResponse.Unmarshal(m, b)
}
func (m *SlotResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SlotResponse.Marshal(b, m, deterministic)
}
func (dst *SlotResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SlotResponse.Merge(dst, src)
}
func (m *SlotResponse) XXX_Size() int {
	return xxx_messageInfo_SlotResponse.Size(m)
}
func (m *SlotResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SlotResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SlotResponse proto.InternalMessageInfo

func (m *SlotResponse) GetSlot() uint64 {
	if m != nil {
		return m.Slot
	}
	return 0
}

func init() {
	proto.RegisterType((*ShuffleRequest)(nil), "ethereum.beacon.rpc.v1.ShuffleRequest")
	proto.RegisterType((*ShuffleResponse)(nil), "ethereum.beacon.rpc.v1.ShuffleResponse")
	proto.RegisterType((*ProposeRequest)(nil), "ethereum.beacon.rpc.v1.ProposeRequest")
	proto.RegisterType((*ProposeResponse)(nil), "ethereum.beacon.rpc.v1.ProposeResponse")
	proto.RegisterType((*AttestRequest)(nil), "ethereum.beacon.rpc.v1.AttestRequest")
	proto.RegisterType((*AttestResponse)(nil), "ethereum.beacon.rpc.v1.AttestResponse")
	proto.RegisterType((*PublicKey)(nil), "ethereum.beacon.rpc.v1.PublicKey")
	proto.RegisterType((*ShardIDResponse)(nil), "ethereum.beacon.rpc.v1.ShardIDResponse")
	proto.RegisterType((*IndexResponse)(nil), "ethereum.beacon.rpc.v1.IndexResponse")
	proto.RegisterType((*SlotResponse)(nil), "ethereum.beacon.rpc.v1.SlotResponse")
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

func (c *beaconServiceClient) LatestAttestation(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (BeaconService_LatestAttestationClient, error) {
	stream, err := c.cc.NewStream(ctx, &_BeaconService_serviceDesc.Streams[2], "/ethereum.beacon.rpc.v1.BeaconService/LatestAttestation", opts...)
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
	Recv() (*v1.AggregatedAttestation, error)
	grpc.ClientStream
}

type beaconServiceLatestAttestationClient struct {
	grpc.ClientStream
}

func (x *beaconServiceLatestAttestationClient) Recv() (*v1.AggregatedAttestation, error) {
	m := new(v1.AggregatedAttestation)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BeaconServiceServer is the server API for BeaconService service.
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
	Send(*v1.AggregatedAttestation) error
	grpc.ServerStream
}

type beaconServiceLatestAttestationServer struct {
	grpc.ServerStream
}

func (x *beaconServiceLatestAttestationServer) Send(m *v1.AggregatedAttestation) error {
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

// AttesterServiceClient is the client API for AttesterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
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
	err := c.cc.Invoke(ctx, "/ethereum.beacon.rpc.v1.AttesterService/AttestHead", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AttesterServiceServer is the server API for AttesterService service.
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

// ValidatorServiceClient is the client API for ValidatorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ValidatorServiceClient interface {
	GetValidatorShardID(ctx context.Context, in *PublicKey, opts ...grpc.CallOption) (*ShardIDResponse, error)
	GetValidatorIndex(ctx context.Context, in *PublicKey, opts ...grpc.CallOption) (*IndexResponse, error)
	GetValidatorSlot(ctx context.Context, in *PublicKey, opts ...grpc.CallOption) (*SlotResponse, error)
}

type validatorServiceClient struct {
	cc *grpc.ClientConn
}

func NewValidatorServiceClient(cc *grpc.ClientConn) ValidatorServiceClient {
	return &validatorServiceClient{cc}
}

func (c *validatorServiceClient) GetValidatorShardID(ctx context.Context, in *PublicKey, opts ...grpc.CallOption) (*ShardIDResponse, error) {
	out := new(ShardIDResponse)
	err := c.cc.Invoke(ctx, "/ethereum.beacon.rpc.v1.ValidatorService/GetValidatorShardID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *validatorServiceClient) GetValidatorIndex(ctx context.Context, in *PublicKey, opts ...grpc.CallOption) (*IndexResponse, error) {
	out := new(IndexResponse)
	err := c.cc.Invoke(ctx, "/ethereum.beacon.rpc.v1.ValidatorService/GetValidatorIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *validatorServiceClient) GetValidatorSlot(ctx context.Context, in *PublicKey, opts ...grpc.CallOption) (*SlotResponse, error) {
	out := new(SlotResponse)
	err := c.cc.Invoke(ctx, "/ethereum.beacon.rpc.v1.ValidatorService/GetValidatorSlot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ValidatorServiceServer is the server API for ValidatorService service.
type ValidatorServiceServer interface {
	GetValidatorShardID(context.Context, *PublicKey) (*ShardIDResponse, error)
	GetValidatorIndex(context.Context, *PublicKey) (*IndexResponse, error)
	GetValidatorSlot(context.Context, *PublicKey) (*SlotResponse, error)
}

func RegisterValidatorServiceServer(s *grpc.Server, srv ValidatorServiceServer) {
	s.RegisterService(&_ValidatorService_serviceDesc, srv)
}

func _ValidatorService_GetValidatorShardID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublicKey)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValidatorServiceServer).GetValidatorShardID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ethereum.beacon.rpc.v1.ValidatorService/GetValidatorShardID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValidatorServiceServer).GetValidatorShardID(ctx, req.(*PublicKey))
	}
	return interceptor(ctx, in, info, handler)
}

func _ValidatorService_GetValidatorIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublicKey)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValidatorServiceServer).GetValidatorIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ethereum.beacon.rpc.v1.ValidatorService/GetValidatorIndex",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValidatorServiceServer).GetValidatorIndex(ctx, req.(*PublicKey))
	}
	return interceptor(ctx, in, info, handler)
}

func _ValidatorService_GetValidatorSlot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublicKey)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValidatorServiceServer).GetValidatorSlot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ethereum.beacon.rpc.v1.ValidatorService/GetValidatorSlot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValidatorServiceServer).GetValidatorSlot(ctx, req.(*PublicKey))
	}
	return interceptor(ctx, in, info, handler)
}

var _ValidatorService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ethereum.beacon.rpc.v1.ValidatorService",
	HandlerType: (*ValidatorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetValidatorShardID",
			Handler:    _ValidatorService_GetValidatorShardID_Handler,
		},
		{
			MethodName: "GetValidatorIndex",
			Handler:    _ValidatorService_GetValidatorIndex_Handler,
		},
		{
			MethodName: "GetValidatorSlot",
			Handler:    _ValidatorService_GetValidatorSlot_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/beacon/rpc/v1/services.proto",
}

func init() {
	proto.RegisterFile("proto/beacon/rpc/v1/services.proto", fileDescriptor_services_f2b7e4f7fe6852c7)
}

var fileDescriptor_services_f2b7e4f7fe6852c7 = []byte{
	// 795 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xdd, 0x6e, 0xeb, 0x44,
	0x10, 0x96, 0x93, 0x9c, 0x03, 0x99, 0xfc, 0xf5, 0xec, 0x81, 0x53, 0xd7, 0xa8, 0x6a, 0x70, 0x69,
	0x49, 0x11, 0x38, 0xa9, 0x91, 0x10, 0x82, 0xde, 0xb4, 0xfc, 0x35, 0x02, 0x41, 0xe5, 0x20, 0x6e,
	0x0a, 0x98, 0x8d, 0xbd, 0x71, 0xac, 0x3a, 0xb6, 0xd9, 0xdd, 0x44, 0x84, 0xd7, 0xe0, 0x9e, 0x77,
	0x40, 0xe2, 0x01, 0xd1, 0xee, 0xc6, 0xee, 0x26, 0xc5, 0xb4, 0xdc, 0xd9, 0xdf, 0x7c, 0xf3, 0xf7,
	0xcd, 0xcc, 0x82, 0x9d, 0xd3, 0x8c, 0x67, 0xc3, 0x29, 0xc1, 0x41, 0x96, 0x0e, 0x69, 0x1e, 0x0c,
	0x57, 0xe7, 0x43, 0x46, 0xe8, 0x2a, 0x0e, 0x08, 0x73, 0xa4, 0x11, 0xbd, 0x22, 0x7c, 0x4e, 0x28,
	0x59, 0x2e, 0x1c, 0x45, 0x73, 0x68, 0x1e, 0x38, 0xab, 0x73, 0x6b, 0xdb, 0x37, 0x77, 0x73, 0xe1,
	0xbb, 0x20, 0x8c, 0xe1, 0xa8, 0xf0, 0xb5, 0xde, 0x8a, 0xb2, 0x2c, 0x4a, 0xc8, 0x50, 0xfe, 0x4d,
	0x97, 0xb3, 0x21, 0x59, 0xe4, 0x7c, 0xbd, 0x31, 0x1e, 0xed, 0x1a, 0x79, 0xbc, 0x20, 0x8c, 0xe3,
	0x45, 0xae, 0x08, 0xf6, 0x35, 0x74, 0x27, 0xf3, 0xe5, 0x6c, 0x96, 0x10, 0x8f, 0xfc, 0xba, 0x24,
	0x8c, 0xa3, 0x8f, 0x60, 0x3f, 0xa0, 0x6b, 0xc6, 0x71, 0x92, 0xc4, 0xbf, 0x93, 0xd0, 0x67, 0x1c,
	0x73, 0xe2, 0xcf, 0x31, 0x9b, 0x9b, 0x46, 0xdf, 0x18, 0xb4, 0xbd, 0x37, 0x75, 0xf3, 0x44, 0x58,
	0xaf, 0x31, 0x9b, 0xdb, 0x7f, 0x1b, 0xd0, 0x2b, 0x43, 0xb1, 0x3c, 0x4b, 0x19, 0x41, 0x17, 0x60,
	0x31, 0x05, 0x85, 0xfe, 0x0a, 0x27, 0x71, 0x88, 0x79, 0x46, 0xfd, 0x38, 0x0d, 0x45, 0xef, 0xa6,
	0xd1, 0xaf, 0x0f, 0x1a, 0x9e, 0x59, 0x30, 0x7e, 0x28, 0x08, 0x63, 0x65, 0x47, 0x27, 0xd0, 0x0d,
	0x96, 0x3c, 0x9b, 0xcd, 0x4a, 0x8f, 0x9a, 0xf4, 0xe8, 0x28, 0xb4, 0xa0, 0x5d, 0x80, 0x85, 0x19,
	0x8b, 0xa3, 0x94, 0x84, 0x3e, 0xe6, 0x5c, 0xb4, 0xc7, 0xe3, 0x2c, 0xf5, 0x59, 0x92, 0x71, 0x66,
	0xd6, 0x55, 0x92, 0x82, 0x71, 0x79, 0x4f, 0x98, 0x08, 0xbb, 0xfd, 0x67, 0x0d, 0xba, 0x37, 0x34,
	0xcb, 0x33, 0x56, 0x2a, 0x70, 0x04, 0xad, 0x1c, 0x53, 0x92, 0x72, 0xbd, 0x6b, 0x50, 0x90, 0x68,
	0x55, 0x10, 0x44, 0x70, 0x3f, 0x5d, 0x2e, 0xa6, 0x84, 0x9a, 0xb5, 0xbe, 0x31, 0x68, 0x78, 0x20,
	0xa0, 0x6f, 0x25, 0x82, 0x8e, 0xa1, 0x43, 0x71, 0x1a, 0xe2, 0xcc, 0xa7, 0x64, 0x45, 0x70, 0x62,
	0xd6, 0x65, 0x8c, 0xb6, 0x02, 0x3d, 0x89, 0xa1, 0x21, 0xbc, 0xd4, 0xcb, 0x9d, 0xc6, 0x7c, 0x81,
	0xd9, 0x9d, 0xd9, 0x90, 0x54, 0xa4, 0x99, 0xae, 0x94, 0x05, 0x7d, 0x02, 0x07, 0xba, 0x03, 0x8e,
	0x22, 0x4a, 0x22, 0x31, 0x1c, 0x16, 0x47, 0xe6, 0xb3, 0x7e, 0x7d, 0xd0, 0xf1, 0xf6, 0x35, 0xc2,
	0x65, 0x61, 0x9f, 0xc4, 0x11, 0xfa, 0x18, 0x9a, 0xe5, 0xe8, 0xcd, 0xe7, 0x7d, 0x63, 0xd0, 0x72,
	0x2d, 0x47, 0x2d, 0x87, 0x53, 0x2c, 0x87, 0xf3, 0x7d, 0xc1, 0xf0, 0xee, 0xc9, 0xf6, 0x08, 0x7a,
	0xa5, 0x3e, 0x9b, 0xb1, 0x1e, 0x02, 0x4c, 0x93, 0x2c, 0xb8, 0xd3, 0xf5, 0x69, 0x4a, 0x44, 0x6e,
	0xc2, 0x2f, 0xd0, 0x51, 0x32, 0x17, 0x82, 0x7e, 0x07, 0x2d, 0xad, 0x2e, 0xe9, 0xd0, 0x72, 0x3f,
	0x70, 0x76, 0x97, 0x3e, 0x77, 0x73, 0x67, 0x75, 0xee, 0x94, 0x75, 0xeb, 0xc3, 0xf2, 0xf4, 0x08,
	0xf6, 0xa7, 0xd0, 0x2d, 0x32, 0x6c, 0x4a, 0x3a, 0x83, 0x3d, 0x5d, 0x1b, 0xad, 0xb0, 0x9e, 0x86,
	0xcb, 0xf2, 0xde, 0x83, 0xe6, 0xcd, 0x72, 0x9a, 0xc4, 0xc1, 0xd7, 0x64, 0x2d, 0x5a, 0xc9, 0xe5,
	0x8f, 0x7f, 0x47, 0xd6, 0xd2, 0xa3, 0xe1, 0x35, 0xf3, 0xc2, 0x6c, 0xbf, 0x2f, 0x76, 0x1a, 0xd3,
	0x70, 0xfc, 0x79, 0x99, 0xe9, 0x00, 0x5e, 0x67, 0x02, 0xf2, 0xe3, 0x70, 0xc3, 0x7f, 0x4d, 0xfe,
	0x8f, 0x43, 0xfb, 0x04, 0x3a, 0xe3, 0x34, 0x24, 0xbf, 0x95, 0xdc, 0x37, 0xe0, 0x59, 0x2c, 0x00,
	0x49, 0xec, 0x78, 0xea, 0xc7, 0xb6, 0xa1, 0x2d, 0x76, 0xaf, 0x64, 0x21, 0x68, 0x88, 0xdd, 0xd9,
	0x44, 0x93, 0xdf, 0xee, 0x1f, 0x75, 0xe8, 0x5c, 0x49, 0x59, 0x26, 0xea, 0xa9, 0x40, 0x1e, 0xbc,
	0xf8, 0x06, 0x8b, 0x4e, 0x14, 0x7c, 0x25, 0xe4, 0x46, 0xaf, 0x1e, 0xcc, 0xf0, 0x0b, 0x71, 0xfd,
	0xd6, 0x71, 0x95, 0xb8, 0x9a, 0xf3, 0xc8, 0x40, 0x3f, 0xc3, 0xbe, 0x8a, 0xf9, 0xd9, 0xee, 0x49,
	0x57, 0x46, 0x3e, 0xab, 0x8a, 0xfc, 0x20, 0xc4, 0xc8, 0x40, 0x39, 0x1c, 0x7e, 0x49, 0x78, 0x30,
	0x9f, 0x54, 0x9d, 0xf8, 0xa9, 0xf3, 0xef, 0x2f, 0x9f, 0xb3, 0xfd, 0x28, 0x59, 0xef, 0x3e, 0xca,
	0xdb, 0x68, 0xf9, 0x63, 0xa1, 0x92, 0xb6, 0x3b, 0x95, 0xbd, 0xfc, 0xbf, 0x15, 0x1c, 0x19, 0x6e,
	0x0a, 0x3d, 0x05, 0x10, 0x5a, 0x8c, 0xe5, 0x16, 0x40, 0x41, 0xd7, 0x04, 0x87, 0xe8, 0xa4, 0xaa,
	0xce, 0xad, 0x83, 0xb0, 0x4e, 0x1f, 0xa3, 0xa9, 0x6e, 0x5c, 0x5a, 0xde, 0x5e, 0x99, 0xcf, 0x87,
	0xf6, 0x06, 0x52, 0x1b, 0x50, 0x19, 0x6a, 0xfb, 0x51, 0xab, 0x56, 0x70, 0xe7, 0xb8, 0xdd, 0xbf,
	0x6a, 0xb0, 0x57, 0xce, 0xa9, 0xc8, 0x8a, 0xe1, 0xe5, 0x57, 0x84, 0xdf, 0xc3, 0xea, 0x26, 0xd0,
	0xdb, 0x95, 0x41, 0x8b, 0x0b, 0xfa, 0xaf, 0xc9, 0x6d, 0xdf, 0xd5, 0x4f, 0xf0, 0x42, 0x4f, 0x21,
	0x0f, 0xe9, 0x29, 0x09, 0x2a, 0x25, 0xdf, 0x3e, 0xc5, 0x5b, 0xd8, 0xdb, 0xea, 0x20, 0xc9, 0xf8,
	0x53, 0xa2, 0xbf, 0x53, 0x59, 0xbe, 0x76, 0xc1, 0xd3, 0xe7, 0x72, 0xb1, 0x3e, 0xfc, 0x27, 0x00,
	0x00, 0xff, 0xff, 0x8e, 0x2a, 0x92, 0xe1, 0xec, 0x07, 0x00, 0x00,
}
