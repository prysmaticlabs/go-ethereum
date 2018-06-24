// Code generated by protoc-gen-go. DO NOT EDIT.
// source: messages.proto

package messages

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type CollationBodyRequest struct {
	ShardId              uint64   `protobuf:"varint,1,opt,name=shard_id,json=shardId,proto3" json:"shard_id,omitempty"`
	Period               uint64   `protobuf:"varint,2,opt,name=period,proto3" json:"period,omitempty"`
	ChunkRoot            []byte   `protobuf:"bytes,3,opt,name=chunk_root,json=chunkRoot,proto3" json:"chunk_root,omitempty"`
	ProposerAddress      []byte   `protobuf:"bytes,4,opt,name=proposer_address,json=proposerAddress,proto3" json:"proposer_address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CollationBodyRequest) Reset()         { *m = CollationBodyRequest{} }
func (m *CollationBodyRequest) String() string { return proto.CompactTextString(m) }
func (*CollationBodyRequest) ProtoMessage()    {}
func (*CollationBodyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_messages_3e84916602e9ec89, []int{0}
}
func (m *CollationBodyRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CollationBodyRequest.Unmarshal(m, b)
}
func (m *CollationBodyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CollationBodyRequest.Marshal(b, m, deterministic)
}
func (dst *CollationBodyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CollationBodyRequest.Merge(dst, src)
}
func (m *CollationBodyRequest) XXX_Size() int {
	return xxx_messageInfo_CollationBodyRequest.Size(m)
}
func (m *CollationBodyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CollationBodyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CollationBodyRequest proto.InternalMessageInfo

func (m *CollationBodyRequest) GetShardId() uint64 {
	if m != nil {
		return m.ShardId
	}
	return 0
}

func (m *CollationBodyRequest) GetPeriod() uint64 {
	if m != nil {
		return m.Period
	}
	return 0
}

func (m *CollationBodyRequest) GetChunkRoot() []byte {
	if m != nil {
		return m.ChunkRoot
	}
	return nil
}

func (m *CollationBodyRequest) GetProposerAddress() []byte {
	if m != nil {
		return m.ProposerAddress
	}
	return nil
}

type CollationBodyResponse struct {
	HeaderHash           []byte   `protobuf:"bytes,1,opt,name=header_hash,json=headerHash,proto3" json:"header_hash,omitempty"`
	Body                 []byte   `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CollationBodyResponse) Reset()         { *m = CollationBodyResponse{} }
func (m *CollationBodyResponse) String() string { return proto.CompactTextString(m) }
func (*CollationBodyResponse) ProtoMessage()    {}
func (*CollationBodyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_messages_3e84916602e9ec89, []int{1}
}
func (m *CollationBodyResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CollationBodyResponse.Unmarshal(m, b)
}
func (m *CollationBodyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CollationBodyResponse.Marshal(b, m, deterministic)
}
func (dst *CollationBodyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CollationBodyResponse.Merge(dst, src)
}
func (m *CollationBodyResponse) XXX_Size() int {
	return xxx_messageInfo_CollationBodyResponse.Size(m)
}
func (m *CollationBodyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CollationBodyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CollationBodyResponse proto.InternalMessageInfo

func (m *CollationBodyResponse) GetHeaderHash() []byte {
	if m != nil {
		return m.HeaderHash
	}
	return nil
}

func (m *CollationBodyResponse) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type PingRequest struct {
	Msg                  string   `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingRequest) Reset()         { *m = PingRequest{} }
func (m *PingRequest) String() string { return proto.CompactTextString(m) }
func (*PingRequest) ProtoMessage()    {}
func (*PingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_messages_3e84916602e9ec89, []int{2}
}
func (m *PingRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingRequest.Unmarshal(m, b)
}
func (m *PingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingRequest.Marshal(b, m, deterministic)
}
func (dst *PingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingRequest.Merge(dst, src)
}
func (m *PingRequest) XXX_Size() int {
	return xxx_messageInfo_PingRequest.Size(m)
}
func (m *PingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PingRequest proto.InternalMessageInfo

func (m *PingRequest) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type PingResponse struct {
	Msg                  string   `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingResponse) Reset()         { *m = PingResponse{} }
func (m *PingResponse) String() string { return proto.CompactTextString(m) }
func (*PingResponse) ProtoMessage()    {}
func (*PingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_messages_3e84916602e9ec89, []int{3}
}
func (m *PingResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingResponse.Unmarshal(m, b)
}
func (m *PingResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingResponse.Marshal(b, m, deterministic)
}
func (dst *PingResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingResponse.Merge(dst, src)
}
func (m *PingResponse) XXX_Size() int {
	return xxx_messageInfo_PingResponse.Size(m)
}
func (m *PingResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PingResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PingResponse proto.InternalMessageInfo

func (m *PingResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func init() {
	proto.RegisterType((*CollationBodyRequest)(nil), "messages.CollationBodyRequest")
	proto.RegisterType((*CollationBodyResponse)(nil), "messages.CollationBodyResponse")
	proto.RegisterType((*PingRequest)(nil), "messages.PingRequest")
	proto.RegisterType((*PingResponse)(nil), "messages.PingResponse")
}

func init() { proto.RegisterFile("messages.proto", fileDescriptor_messages_3e84916602e9ec89) }

var fileDescriptor_messages_3e84916602e9ec89 = []byte{
	// 243 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x90, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x86, 0x15, 0x1a, 0x95, 0xf6, 0x1a, 0x41, 0x65, 0x01, 0x0a, 0x03, 0x6a, 0x94, 0xa9, 0x2c,
	0x2c, 0x3c, 0x01, 0xb0, 0x80, 0xc4, 0x80, 0xfc, 0x02, 0x91, 0xcb, 0x9d, 0xe2, 0x88, 0x36, 0x67,
	0x7c, 0xee, 0xd0, 0xe7, 0xe0, 0x85, 0x11, 0x97, 0x66, 0xa1, 0xdb, 0xff, 0x7f, 0xfe, 0xad, 0xfb,
	0xef, 0xe0, 0x62, 0x47, 0x22, 0xae, 0x25, 0x79, 0x08, 0x91, 0x13, 0x9b, 0xd9, 0xe8, 0xeb, 0x9f,
	0x0c, 0xae, 0x5e, 0x78, 0xbb, 0x75, 0xa9, 0xe3, 0xfe, 0x99, 0xf1, 0x60, 0xe9, 0x7b, 0x4f, 0x92,
	0xcc, 0x2d, 0xcc, 0xc4, 0xbb, 0x88, 0x4d, 0x87, 0x65, 0x56, 0x65, 0xeb, 0xdc, 0x9e, 0xab, 0x7f,
	0x43, 0x73, 0x03, 0xd3, 0x40, 0xb1, 0x63, 0x2c, 0xcf, 0xf4, 0xe1, 0xe8, 0xcc, 0x1d, 0xc0, 0xa7,
	0xdf, 0xf7, 0x5f, 0x4d, 0x64, 0x4e, 0xe5, 0xa4, 0xca, 0xd6, 0x85, 0x9d, 0x2b, 0xb1, 0xcc, 0xc9,
	0xdc, 0xc3, 0x32, 0x44, 0x0e, 0x2c, 0x14, 0x1b, 0x87, 0x18, 0x49, 0xa4, 0xcc, 0x35, 0x74, 0x39,
	0xf2, 0xa7, 0x01, 0xd7, 0xef, 0x70, 0xfd, 0xaf, 0x94, 0x04, 0xee, 0x85, 0xcc, 0x0a, 0x16, 0x9e,
	0x1c, 0x52, 0x6c, 0xbc, 0x13, 0xaf, 0xc5, 0x0a, 0x0b, 0x03, 0x7a, 0x75, 0xe2, 0x8d, 0x81, 0x7c,
	0xc3, 0x78, 0xd0, 0x66, 0x85, 0x55, 0x5d, 0xaf, 0x60, 0xf1, 0xd1, 0xf5, 0xed, 0xb8, 0xd9, 0x12,
	0x26, 0x3b, 0x69, 0xf5, 0xef, 0xdc, 0xfe, 0xc9, 0xba, 0x82, 0x62, 0x08, 0x1c, 0xa7, 0x9c, 0x24,
	0x36, 0x53, 0xbd, 0xdb, 0xe3, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x91, 0xd2, 0x3d, 0xec, 0x49,
	0x01, 0x00, 0x00,
}
