// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpc.proto

package rpc

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// The request message containing the user's name.
type Bytes struct {
	Bytes                []byte   `protobuf:"bytes,1,opt,name=bytes,proto3" json:"bytes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Bytes) Reset()         { *m = Bytes{} }
func (m *Bytes) String() string { return proto.CompactTextString(m) }
func (*Bytes) ProtoMessage()    {}
func (*Bytes) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{0}
}

func (m *Bytes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Bytes.Unmarshal(m, b)
}
func (m *Bytes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Bytes.Marshal(b, m, deterministic)
}
func (m *Bytes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Bytes.Merge(m, src)
}
func (m *Bytes) XXX_Size() int {
	return xxx_messageInfo_Bytes.Size(m)
}
func (m *Bytes) XXX_DiscardUnknown() {
	xxx_messageInfo_Bytes.DiscardUnknown(m)
}

var xxx_messageInfo_Bytes proto.InternalMessageInfo

func (m *Bytes) GetBytes() []byte {
	if m != nil {
		return m.Bytes
	}
	return nil
}

type Address struct {
	Address              string   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Address) Reset()         { *m = Address{} }
func (m *Address) String() string { return proto.CompactTextString(m) }
func (*Address) ProtoMessage()    {}
func (*Address) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{1}
}

func (m *Address) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Address.Unmarshal(m, b)
}
func (m *Address) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Address.Marshal(b, m, deterministic)
}
func (m *Address) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Address.Merge(m, src)
}
func (m *Address) XXX_Size() int {
	return xxx_messageInfo_Address.Size(m)
}
func (m *Address) XXX_DiscardUnknown() {
	xxx_messageInfo_Address.DiscardUnknown(m)
}

var xxx_messageInfo_Address proto.InternalMessageInfo

func (m *Address) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type Symbol struct {
	Symbol               string   `protobuf:"bytes,1,opt,name=Symbol,proto3" json:"Symbol,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Symbol) Reset()         { *m = Symbol{} }
func (m *Symbol) String() string { return proto.CompactTextString(m) }
func (*Symbol) ProtoMessage()    {}
func (*Symbol) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{2}
}

func (m *Symbol) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Symbol.Unmarshal(m, b)
}
func (m *Symbol) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Symbol.Marshal(b, m, deterministic)
}
func (m *Symbol) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Symbol.Merge(m, src)
}
func (m *Symbol) XXX_Size() int {
	return xxx_messageInfo_Symbol.Size(m)
}
func (m *Symbol) XXX_DiscardUnknown() {
	xxx_messageInfo_Symbol.DiscardUnknown(m)
}

var xxx_messageInfo_Symbol proto.InternalMessageInfo

func (m *Symbol) GetSymbol() string {
	if m != nil {
		return m.Symbol
	}
	return ""
}

type Hash struct {
	Hash                 string   `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Hash) Reset()         { *m = Hash{} }
func (m *Hash) String() string { return proto.CompactTextString(m) }
func (*Hash) ProtoMessage()    {}
func (*Hash) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{3}
}

func (m *Hash) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Hash.Unmarshal(m, b)
}
func (m *Hash) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Hash.Marshal(b, m, deterministic)
}
func (m *Hash) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Hash.Merge(m, src)
}
func (m *Hash) XXX_Size() int {
	return xxx_messageInfo_Hash.Size(m)
}
func (m *Hash) XXX_DiscardUnknown() {
	xxx_messageInfo_Hash.DiscardUnknown(m)
}

var xxx_messageInfo_Hash proto.InternalMessageInfo

func (m *Hash) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

type Height struct {
	Height               uint64   `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Count                uint64   `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Height) Reset()         { *m = Height{} }
func (m *Height) String() string { return proto.CompactTextString(m) }
func (*Height) ProtoMessage()    {}
func (*Height) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{4}
}

func (m *Height) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Height.Unmarshal(m, b)
}
func (m *Height) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Height.Marshal(b, m, deterministic)
}
func (m *Height) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Height.Merge(m, src)
}
func (m *Height) XXX_Size() int {
	return xxx_messageInfo_Height.Size(m)
}
func (m *Height) XXX_DiscardUnknown() {
	xxx_messageInfo_Height.DiscardUnknown(m)
}

var xxx_messageInfo_Height proto.InternalMessageInfo

func (m *Height) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *Height) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

type Null struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Null) Reset()         { *m = Null{} }
func (m *Null) String() string { return proto.CompactTextString(m) }
func (*Null) ProtoMessage()    {}
func (*Null) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{5}
}

func (m *Null) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Null.Unmarshal(m, b)
}
func (m *Null) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Null.Marshal(b, m, deterministic)
}
func (m *Null) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Null.Merge(m, src)
}
func (m *Null) XXX_Size() int {
	return xxx_messageInfo_Null.Size(m)
}
func (m *Null) XXX_DiscardUnknown() {
	xxx_messageInfo_Null.DiscardUnknown(m)
}

var xxx_messageInfo_Null proto.InternalMessageInfo

type Method struct {
	Contract             string   `protobuf:"bytes,1,opt,name=contract,proto3" json:"contract,omitempty"`
	Method               string   `protobuf:"bytes,2,opt,name=method,proto3" json:"method,omitempty"`
	Params               []string `protobuf:"bytes,3,rep,name=params,proto3" json:"params,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Method) Reset()         { *m = Method{} }
func (m *Method) String() string { return proto.CompactTextString(m) }
func (*Method) ProtoMessage()    {}
func (*Method) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{6}
}

func (m *Method) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Method.Unmarshal(m, b)
}
func (m *Method) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Method.Marshal(b, m, deterministic)
}
func (m *Method) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Method.Merge(m, src)
}
func (m *Method) XXX_Size() int {
	return xxx_messageInfo_Method.Size(m)
}
func (m *Method) XXX_DiscardUnknown() {
	xxx_messageInfo_Method.DiscardUnknown(m)
}

var xxx_messageInfo_Method proto.InternalMessageInfo

func (m *Method) GetContract() string {
	if m != nil {
		return m.Contract
	}
	return ""
}

func (m *Method) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *Method) GetParams() []string {
	if m != nil {
		return m.Params
	}
	return nil
}

// The response message containing the greetings
type Response struct {
	Code                 int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Result               []byte   `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
	Err                  string   `protobuf:"bytes,3,opt,name=err,proto3" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_77a6da22d6a3feb1, []int{7}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Response) GetResult() []byte {
	if m != nil {
		return m.Result
	}
	return nil
}

func (m *Response) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

func init() {
	proto.RegisterType((*Bytes)(nil), "rpc.Bytes")
	proto.RegisterType((*Address)(nil), "rpc.Address")
	proto.RegisterType((*Symbol)(nil), "rpc.Symbol")
	proto.RegisterType((*Hash)(nil), "rpc.Hash")
	proto.RegisterType((*Height)(nil), "rpc.Height")
	proto.RegisterType((*Null)(nil), "rpc.Null")
	proto.RegisterType((*Method)(nil), "rpc.Method")
	proto.RegisterType((*Response)(nil), "rpc.Response")
}

func init() { proto.RegisterFile("rpc.proto", fileDescriptor_77a6da22d6a3feb1) }

var fileDescriptor_77a6da22d6a3feb1 = []byte{
	// 497 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0xed, 0x8a, 0x1a, 0x31,
	0x14, 0xed, 0x56, 0x9d, 0xd5, 0xab, 0xfb, 0x41, 0x5a, 0x8a, 0x08, 0x05, 0x99, 0x85, 0x76, 0x4b,
	0x17, 0x2b, 0xbb, 0xd0, 0xff, 0xeb, 0xfe, 0x18, 0x0b, 0x5b, 0x91, 0x59, 0x5f, 0x20, 0x66, 0xee,
	0x3a, 0xb2, 0x63, 0x32, 0x24, 0x11, 0xea, 0x1b, 0xf6, 0xb1, 0x4a, 0x6e, 0x32, 0x45, 0xb0, 0x1d,
	0xf7, 0xdf, 0x39, 0xc9, 0x39, 0xc9, 0x99, 0x9c, 0xab, 0xd0, 0xd1, 0xa5, 0x18, 0x95, 0x5a, 0x59,
	0xc5, 0x1a, 0xba, 0x14, 0xf1, 0x47, 0x68, 0x4d, 0x76, 0x16, 0x0d, 0x7b, 0x0f, 0xad, 0xa5, 0x03,
	0xfd, 0x93, 0xe1, 0xc9, 0x75, 0x2f, 0xf5, 0x24, 0xbe, 0x82, 0xd3, 0xfb, 0x2c, 0xd3, 0x68, 0x0c,
	0xeb, 0xc3, 0x29, 0xf7, 0x90, 0x24, 0x9d, 0xb4, 0xa2, 0xf1, 0x10, 0xa2, 0xa7, 0xdd, 0x66, 0xa9,
	0x0a, 0xf6, 0xa1, 0x42, 0x41, 0x12, 0x58, 0x3c, 0x80, 0xe6, 0x94, 0x9b, 0x9c, 0x31, 0x68, 0xe6,
	0xdc, 0xe4, 0x61, 0x97, 0x70, 0xfc, 0x1d, 0xa2, 0x29, 0xae, 0x57, 0xb9, 0x75, 0xee, 0x9c, 0x10,
	0xed, 0x37, 0xd3, 0xc0, 0x5c, 0x34, 0xa1, 0xb6, 0xd2, 0xf6, 0xdf, 0xd2, 0xb2, 0x27, 0x71, 0x04,
	0xcd, 0xd9, 0xb6, 0x28, 0xe2, 0x05, 0x44, 0x3f, 0xd1, 0xe6, 0x2a, 0x63, 0x03, 0x68, 0x0b, 0x25,
	0xad, 0xe6, 0xc2, 0x86, 0x1b, 0xfe, 0x72, 0x77, 0xf6, 0x86, 0x54, 0x74, 0x48, 0x27, 0x0d, 0xcc,
	0xad, 0x97, 0x5c, 0xf3, 0x8d, 0xe9, 0x37, 0x86, 0x0d, 0xb7, 0xee, 0x59, 0x3c, 0x85, 0x76, 0x8a,
	0xa6, 0x54, 0xd2, 0xa0, 0x4b, 0x2d, 0x54, 0x86, 0x74, 0x66, 0x2b, 0x25, 0xec, 0x7c, 0x1a, 0xcd,
	0xb6, 0xf0, 0xa1, 0x7a, 0x69, 0x60, 0xec, 0x12, 0x1a, 0xa8, 0x75, 0xbf, 0x41, 0x97, 0x38, 0x78,
	0xfb, 0x3b, 0x82, 0xd3, 0x44, 0x23, 0x5a, 0xd4, 0x6c, 0x04, 0x17, 0x4f, 0x28, 0xb3, 0x85, 0xe6,
	0xd2, 0x70, 0x61, 0xd7, 0x4a, 0x32, 0x18, 0xb9, 0x46, 0xa8, 0x83, 0xc1, 0x19, 0xe1, 0xea, 0xde,
	0xf8, 0x0d, 0xfb, 0x0a, 0x90, 0xa0, 0xbd, 0x17, 0xf4, 0xc5, 0xac, 0x47, 0xdb, 0xa1, 0x8f, 0x43,
	0xf1, 0x18, 0x2e, 0x12, 0xb4, 0x73, 0x94, 0xd9, 0x5a, 0xae, 0x66, 0x4a, 0x0a, 0x3c, 0xe6, 0xb8,
	0x81, 0xf3, 0x04, 0xed, 0x7e, 0x9a, 0x0e, 0x49, 0x5c, 0x57, 0xff, 0x53, 0x4f, 0x0a, 0x25, 0x5e,
	0x26, 0x3b, 0xaa, 0xb3, 0x4e, 0x3d, 0x86, 0xcb, 0x3d, 0xb5, 0x2f, 0xb2, 0xeb, 0xf5, 0x44, 0x0e,
	0x1d, 0xdf, 0x28, 0x7f, 0x70, 0xa4, 0x5c, 0xae, 0xf0, 0x88, 0xe1, 0x9a, 0x5e, 0x67, 0xae, 0x54,
	0xb1, 0xf8, 0x65, 0x42, 0x18, 0x37, 0x12, 0xff, 0x7a, 0xc7, 0xb3, 0x04, 0xed, 0x23, 0x37, 0x36,
	0x24, 0xa9, 0x13, 0xdf, 0x40, 0x37, 0x41, 0xfb, 0x50, 0x4d, 0xce, 0x91, 0x37, 0xbc, 0x83, 0x77,
	0x7b, 0xea, 0xc9, 0x2e, 0xfc, 0x12, 0x7c, 0x72, 0x4f, 0x0e, 0x4d, 0xb7, 0xc0, 0x5c, 0xaf, 0xfe,
	0xcc, 0x57, 0x7a, 0x46, 0x70, 0x5e, 0xdd, 0x12, 0xe6, 0xdd, 0xeb, 0x3d, 0x39, 0xd4, 0x7f, 0x86,
	0xce, 0x42, 0xbd, 0xa0, 0x7c, 0x5c, 0x9b, 0xfa, 0xef, 0xfd, 0x02, 0xdd, 0x30, 0x61, 0x47, 0xa5,
	0x63, 0xca, 0xfd, 0xa0, 0xe4, 0xf3, 0x5a, 0x6f, 0x30, 0x7b, 0xc5, 0x63, 0x5e, 0x41, 0x6b, 0x8e,
	0xa8, 0xeb, 0xeb, 0xf9, 0x04, 0xed, 0x99, 0xca, 0xf0, 0x87, 0x7c, 0x56, 0x75, 0xba, 0x65, 0x44,
	0x7f, 0x5c, 0x77, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0xeb, 0x30, 0xa0, 0x33, 0xc5, 0x04, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GreeterClient is the client API for Greeter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GreeterClient interface {
	// Sends a greeting
	SendTransaction(ctx context.Context, in *Bytes, opts ...grpc.CallOption) (*Response, error)
	GetAccount(ctx context.Context, in *Address, opts ...grpc.CallOption) (*Response, error)
	GetPendingNonce(ctx context.Context, in *Address, opts ...grpc.CallOption) (*Response, error)
	GetTransaction(ctx context.Context, in *Hash, opts ...grpc.CallOption) (*Response, error)
	GetBlockByHash(ctx context.Context, in *Hash, opts ...grpc.CallOption) (*Response, error)
	GetBlockByHeight(ctx context.Context, in *Height, opts ...grpc.CallOption) (*Response, error)
	GetBlockByRange(ctx context.Context, in *Height, opts ...grpc.CallOption) (*Response, error)
	GetPoolTxs(ctx context.Context, in *Null, opts ...grpc.CallOption) (*Response, error)
	GetLastHeight(ctx context.Context, in *Null, opts ...grpc.CallOption) (*Response, error)
	GetContract(ctx context.Context, in *Address, opts ...grpc.CallOption) (*Response, error)
	GetContractBySymbol(ctx context.Context, in *Symbol, opts ...grpc.CallOption) (*Response, error)
	GetAddressBySymbol(ctx context.Context, in *Symbol, opts ...grpc.CallOption) (*Response, error)
	ContractMethod(ctx context.Context, in *Method, opts ...grpc.CallOption) (*Response, error)
	TokenList(ctx context.Context, in *Null, opts ...grpc.CallOption) (*Response, error)
	AccountList(ctx context.Context, in *Null, opts ...grpc.CallOption) (*Response, error)
	GetConfirmedHeight(ctx context.Context, in *Null, opts ...grpc.CallOption) (*Response, error)
	Peers(ctx context.Context, in *Null, opts ...grpc.CallOption) (*Response, error)
	NodeInfo(ctx context.Context, in *Null, opts ...grpc.CallOption) (*Response, error)
}

type greeterClient struct {
	cc *grpc.ClientConn
}

func NewGreeterClient(cc *grpc.ClientConn) GreeterClient {
	return &greeterClient{cc}
}

func (c *greeterClient) SendTransaction(ctx context.Context, in *Bytes, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/rpc.Greeter/SendTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) GetAccount(ctx context.Context, in *Address, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/rpc.Greeter/GetAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) GetPendingNonce(ctx context.Context, in *Address, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/rpc.Greeter/GetPendingNonce", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) GetTransaction(ctx context.Context, in *Hash, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/rpc.Greeter/GetTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) GetBlockByHash(ctx context.Context, in *Hash, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/rpc.Greeter/GetBlockByHash", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) GetBlockByHeight(ctx context.Context, in *Height, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/rpc.Greeter/GetBlockByHeight", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) GetBlockByRange(ctx context.Context, in *Height, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/rpc.Greeter/GetBlockByRange", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) GetPoolTxs(ctx context.Context, in *Null, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/rpc.Greeter/GetPoolTxs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) GetLastHeight(ctx context.Context, in *Null, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/rpc.Greeter/GetLastHeight", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) GetContract(ctx context.Context, in *Address, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/rpc.Greeter/GetContract", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) GetContractBySymbol(ctx context.Context, in *Symbol, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/rpc.Greeter/GetContractBySymbol", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) GetAddressBySymbol(ctx context.Context, in *Symbol, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/rpc.Greeter/GetAddressBySymbol", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) ContractMethod(ctx context.Context, in *Method, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/rpc.Greeter/ContractMethod", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) TokenList(ctx context.Context, in *Null, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/rpc.Greeter/TokenList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) AccountList(ctx context.Context, in *Null, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/rpc.Greeter/AccountList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) GetConfirmedHeight(ctx context.Context, in *Null, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/rpc.Greeter/GetConfirmedHeight", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) Peers(ctx context.Context, in *Null, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/rpc.Greeter/Peers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) NodeInfo(ctx context.Context, in *Null, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/rpc.Greeter/NodeInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GreeterServer is the server API for Greeter service.
type GreeterServer interface {
	// Sends a greeting
	SendTransaction(context.Context, *Bytes) (*Response, error)
	GetAccount(context.Context, *Address) (*Response, error)
	GetPendingNonce(context.Context, *Address) (*Response, error)
	GetTransaction(context.Context, *Hash) (*Response, error)
	GetBlockByHash(context.Context, *Hash) (*Response, error)
	GetBlockByHeight(context.Context, *Height) (*Response, error)
	GetBlockByRange(context.Context, *Height) (*Response, error)
	GetPoolTxs(context.Context, *Null) (*Response, error)
	GetLastHeight(context.Context, *Null) (*Response, error)
	GetContract(context.Context, *Address) (*Response, error)
	GetContractBySymbol(context.Context, *Symbol) (*Response, error)
	GetAddressBySymbol(context.Context, *Symbol) (*Response, error)
	ContractMethod(context.Context, *Method) (*Response, error)
	TokenList(context.Context, *Null) (*Response, error)
	AccountList(context.Context, *Null) (*Response, error)
	GetConfirmedHeight(context.Context, *Null) (*Response, error)
	Peers(context.Context, *Null) (*Response, error)
	NodeInfo(context.Context, *Null) (*Response, error)
}

// UnimplementedGreeterServer can be embedded to have forward compatible implementations.
type UnimplementedGreeterServer struct {
}

func (*UnimplementedGreeterServer) SendTransaction(ctx context.Context, req *Bytes) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendTransaction not implemented")
}
func (*UnimplementedGreeterServer) GetAccount(ctx context.Context, req *Address) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccount not implemented")
}
func (*UnimplementedGreeterServer) GetPendingNonce(ctx context.Context, req *Address) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPendingNonce not implemented")
}
func (*UnimplementedGreeterServer) GetTransaction(ctx context.Context, req *Hash) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransaction not implemented")
}
func (*UnimplementedGreeterServer) GetBlockByHash(ctx context.Context, req *Hash) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockByHash not implemented")
}
func (*UnimplementedGreeterServer) GetBlockByHeight(ctx context.Context, req *Height) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockByHeight not implemented")
}
func (*UnimplementedGreeterServer) GetBlockByRange(ctx context.Context, req *Height) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockByRange not implemented")
}
func (*UnimplementedGreeterServer) GetPoolTxs(ctx context.Context, req *Null) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPoolTxs not implemented")
}
func (*UnimplementedGreeterServer) GetLastHeight(ctx context.Context, req *Null) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLastHeight not implemented")
}
func (*UnimplementedGreeterServer) GetContract(ctx context.Context, req *Address) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetContract not implemented")
}
func (*UnimplementedGreeterServer) GetContractBySymbol(ctx context.Context, req *Symbol) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetContractBySymbol not implemented")
}
func (*UnimplementedGreeterServer) GetAddressBySymbol(ctx context.Context, req *Symbol) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAddressBySymbol not implemented")
}
func (*UnimplementedGreeterServer) ContractMethod(ctx context.Context, req *Method) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ContractMethod not implemented")
}
func (*UnimplementedGreeterServer) TokenList(ctx context.Context, req *Null) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TokenList not implemented")
}
func (*UnimplementedGreeterServer) AccountList(ctx context.Context, req *Null) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AccountList not implemented")
}
func (*UnimplementedGreeterServer) GetConfirmedHeight(ctx context.Context, req *Null) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfirmedHeight not implemented")
}
func (*UnimplementedGreeterServer) Peers(ctx context.Context, req *Null) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Peers not implemented")
}
func (*UnimplementedGreeterServer) NodeInfo(ctx context.Context, req *Null) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NodeInfo not implemented")
}

func RegisterGreeterServer(s *grpc.Server, srv GreeterServer) {
	s.RegisterService(&_Greeter_serviceDesc, srv)
}

func _Greeter_SendTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Bytes)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).SendTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/SendTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).SendTransaction(ctx, req.(*Bytes))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_GetAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Address)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).GetAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/GetAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).GetAccount(ctx, req.(*Address))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_GetPendingNonce_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Address)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).GetPendingNonce(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/GetPendingNonce",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).GetPendingNonce(ctx, req.(*Address))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_GetTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Hash)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).GetTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/GetTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).GetTransaction(ctx, req.(*Hash))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_GetBlockByHash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Hash)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).GetBlockByHash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/GetBlockByHash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).GetBlockByHash(ctx, req.(*Hash))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_GetBlockByHeight_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Height)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).GetBlockByHeight(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/GetBlockByHeight",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).GetBlockByHeight(ctx, req.(*Height))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_GetBlockByRange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Height)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).GetBlockByRange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/GetBlockByRange",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).GetBlockByRange(ctx, req.(*Height))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_GetPoolTxs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Null)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).GetPoolTxs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/GetPoolTxs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).GetPoolTxs(ctx, req.(*Null))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_GetLastHeight_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Null)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).GetLastHeight(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/GetLastHeight",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).GetLastHeight(ctx, req.(*Null))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_GetContract_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Address)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).GetContract(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/GetContract",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).GetContract(ctx, req.(*Address))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_GetContractBySymbol_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Symbol)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).GetContractBySymbol(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/GetContractBySymbol",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).GetContractBySymbol(ctx, req.(*Symbol))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_GetAddressBySymbol_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Symbol)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).GetAddressBySymbol(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/GetAddressBySymbol",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).GetAddressBySymbol(ctx, req.(*Symbol))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_ContractMethod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Method)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).ContractMethod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/ContractMethod",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).ContractMethod(ctx, req.(*Method))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_TokenList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Null)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).TokenList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/TokenList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).TokenList(ctx, req.(*Null))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_AccountList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Null)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).AccountList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/AccountList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).AccountList(ctx, req.(*Null))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_GetConfirmedHeight_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Null)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).GetConfirmedHeight(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/GetConfirmedHeight",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).GetConfirmedHeight(ctx, req.(*Null))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_Peers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Null)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).Peers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/Peers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).Peers(ctx, req.(*Null))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_NodeInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Null)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).NodeInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/NodeInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).NodeInfo(ctx, req.(*Null))
	}
	return interceptor(ctx, in, info, handler)
}

var _Greeter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.Greeter",
	HandlerType: (*GreeterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendTransaction",
			Handler:    _Greeter_SendTransaction_Handler,
		},
		{
			MethodName: "GetAccount",
			Handler:    _Greeter_GetAccount_Handler,
		},
		{
			MethodName: "GetPendingNonce",
			Handler:    _Greeter_GetPendingNonce_Handler,
		},
		{
			MethodName: "GetTransaction",
			Handler:    _Greeter_GetTransaction_Handler,
		},
		{
			MethodName: "GetBlockByHash",
			Handler:    _Greeter_GetBlockByHash_Handler,
		},
		{
			MethodName: "GetBlockByHeight",
			Handler:    _Greeter_GetBlockByHeight_Handler,
		},
		{
			MethodName: "GetBlockByRange",
			Handler:    _Greeter_GetBlockByRange_Handler,
		},
		{
			MethodName: "GetPoolTxs",
			Handler:    _Greeter_GetPoolTxs_Handler,
		},
		{
			MethodName: "GetLastHeight",
			Handler:    _Greeter_GetLastHeight_Handler,
		},
		{
			MethodName: "GetContract",
			Handler:    _Greeter_GetContract_Handler,
		},
		{
			MethodName: "GetContractBySymbol",
			Handler:    _Greeter_GetContractBySymbol_Handler,
		},
		{
			MethodName: "GetAddressBySymbol",
			Handler:    _Greeter_GetAddressBySymbol_Handler,
		},
		{
			MethodName: "ContractMethod",
			Handler:    _Greeter_ContractMethod_Handler,
		},
		{
			MethodName: "TokenList",
			Handler:    _Greeter_TokenList_Handler,
		},
		{
			MethodName: "AccountList",
			Handler:    _Greeter_AccountList_Handler,
		},
		{
			MethodName: "GetConfirmedHeight",
			Handler:    _Greeter_GetConfirmedHeight_Handler,
		},
		{
			MethodName: "Peers",
			Handler:    _Greeter_Peers_Handler,
		},
		{
			MethodName: "NodeInfo",
			Handler:    _Greeter_NodeInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc.proto",
}
