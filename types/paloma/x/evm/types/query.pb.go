// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: paloma/evm/query.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type QueryGetValsetByIDRequest struct {
	ValsetID uint64 `protobuf:"varint,1,opt,name=valsetID,proto3" json:"valsetID,omitempty"`
	ChainID  string `protobuf:"bytes,2,opt,name=chainID,proto3" json:"chainID,omitempty"`
}

func (m *QueryGetValsetByIDRequest) Reset()         { *m = QueryGetValsetByIDRequest{} }
func (m *QueryGetValsetByIDRequest) String() string { return proto.CompactTextString(m) }
func (*QueryGetValsetByIDRequest) ProtoMessage()    {}
func (*QueryGetValsetByIDRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_eb5cae5990a5e532, []int{0}
}
func (m *QueryGetValsetByIDRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetValsetByIDRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetValsetByIDRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetValsetByIDRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetValsetByIDRequest.Merge(m, src)
}
func (m *QueryGetValsetByIDRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetValsetByIDRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetValsetByIDRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetValsetByIDRequest proto.InternalMessageInfo

func (m *QueryGetValsetByIDRequest) GetValsetID() uint64 {
	if m != nil {
		return m.ValsetID
	}
	return 0
}

func (m *QueryGetValsetByIDRequest) GetChainID() string {
	if m != nil {
		return m.ChainID
	}
	return ""
}

type QueryGetValsetByIDResponse struct {
	Valset *Valset `protobuf:"bytes,1,opt,name=valset,proto3" json:"valset,omitempty"`
}

func (m *QueryGetValsetByIDResponse) Reset()         { *m = QueryGetValsetByIDResponse{} }
func (m *QueryGetValsetByIDResponse) String() string { return proto.CompactTextString(m) }
func (*QueryGetValsetByIDResponse) ProtoMessage()    {}
func (*QueryGetValsetByIDResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_eb5cae5990a5e532, []int{1}
}
func (m *QueryGetValsetByIDResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetValsetByIDResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetValsetByIDResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetValsetByIDResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetValsetByIDResponse.Merge(m, src)
}
func (m *QueryGetValsetByIDResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetValsetByIDResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetValsetByIDResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetValsetByIDResponse proto.InternalMessageInfo

func (m *QueryGetValsetByIDResponse) GetValset() *Valset {
	if m != nil {
		return m.Valset
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryGetValsetByIDRequest)(nil), "palomachain.paloma.evm.QueryGetValsetByIDRequest")
	proto.RegisterType((*QueryGetValsetByIDResponse)(nil), "palomachain.paloma.evm.QueryGetValsetByIDResponse")
}

func init() { proto.RegisterFile("paloma/evm/query.proto", fileDescriptor_eb5cae5990a5e532) }

var fileDescriptor_eb5cae5990a5e532 = []byte{
	// 360 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0x41, 0x4b, 0xeb, 0x40,
	0x10, 0xc7, 0xbb, 0xe5, 0xbd, 0xbe, 0xf7, 0xf6, 0xe1, 0x25, 0x8a, 0xd4, 0x20, 0xa1, 0xf4, 0x54,
	0x3d, 0x64, 0x69, 0x95, 0x9e, 0x3c, 0xd5, 0x80, 0xe4, 0xd8, 0x20, 0x1e, 0xbc, 0x94, 0x4d, 0x1d,
	0xd2, 0x40, 0xb3, 0x93, 0x76, 0x37, 0xc1, 0x20, 0x5e, 0xfc, 0x04, 0x82, 0x5f, 0xc5, 0xb3, 0x67,
	0x8f, 0x05, 0x2f, 0x1e, 0xa5, 0xf5, 0x83, 0x48, 0x77, 0xa3, 0x54, 0x68, 0x0f, 0xde, 0x66, 0x26,
	0xff, 0xff, 0x2f, 0xff, 0xd9, 0xa1, 0xbb, 0x29, 0x1f, 0x63, 0xc2, 0x19, 0xe4, 0x09, 0x9b, 0x64,
	0x30, 0x2d, 0xdc, 0x74, 0x8a, 0x0a, 0xad, 0x72, 0x3e, 0x1c, 0xf1, 0x58, 0xb8, 0xa6, 0x76, 0x21,
	0x4f, 0xec, 0x9d, 0x08, 0x23, 0xd4, 0x12, 0xb6, 0xac, 0x8c, 0xda, 0xde, 0x8f, 0x10, 0xa3, 0x31,
	0x30, 0x9e, 0xc6, 0x8c, 0x0b, 0x81, 0x8a, 0xab, 0x18, 0x85, 0x2c, 0xbf, 0x1e, 0x0e, 0x51, 0x26,
	0x28, 0x59, 0xc8, 0x25, 0x98, 0x9f, 0xb0, 0xbc, 0x1d, 0x82, 0xe2, 0x6d, 0x96, 0xf2, 0x28, 0x16,
	0x5a, 0x5c, 0x6a, 0xb7, 0x97, 0x41, 0x54, 0x36, 0x15, 0x52, 0xa1, 0x00, 0x33, 0x6c, 0xf6, 0xe9,
	0x5e, 0x7f, 0x69, 0x3b, 0x03, 0x75, 0xc1, 0xc7, 0x12, 0x54, 0xaf, 0xf0, 0xbd, 0x00, 0x26, 0x19,
	0x48, 0x65, 0xd9, 0xf4, 0x6f, 0xae, 0x87, 0xbe, 0x57, 0x27, 0x0d, 0xd2, 0xfa, 0x15, 0x7c, 0xf5,
	0x56, 0x9d, 0xfe, 0xd1, 0x1b, 0xf8, 0x5e, 0xbd, 0xda, 0x20, 0xad, 0x7f, 0xc1, 0x67, 0xdb, 0x3c,
	0xa7, 0xf6, 0x3a, 0xa4, 0x4c, 0x51, 0x48, 0xb0, 0xba, 0xb4, 0x66, 0x18, 0x9a, 0xf8, 0xbf, 0xe3,
	0xb8, 0xeb, 0x9f, 0xc3, 0x35, 0xde, 0xa0, 0x54, 0x77, 0x9e, 0x08, 0xfd, 0xad, 0xb1, 0xd6, 0x23,
	0xa1, 0x5b, 0xdf, 0xd8, 0x56, 0x7b, 0x13, 0x63, 0xe3, 0x6a, 0x76, 0xe7, 0x27, 0x16, 0x13, 0xbd,
	0x79, 0x72, 0xf7, 0xf2, 0xfe, 0x50, 0xed, 0x5a, 0xc7, 0x6c, 0xc5, 0xcb, 0x56, 0xae, 0x1c, 0x81,
	0x1a, 0x98, 0xb8, 0x83, 0xb0, 0x18, 0xc4, 0x57, 0xec, 0xc6, 0x74, 0xca, 0xf7, 0x6e, 0x7b, 0xa7,
	0xcf, 0x73, 0x87, 0xcc, 0xe6, 0x0e, 0x79, 0x9b, 0x3b, 0xe4, 0x7e, 0xe1, 0x54, 0x66, 0x0b, 0xa7,
	0xf2, 0xba, 0x70, 0x2a, 0x97, 0x07, 0x51, 0xac, 0x46, 0x59, 0xe8, 0x0e, 0x31, 0x59, 0x47, 0xbe,
	0xd6, 0x6c, 0x55, 0xa4, 0x20, 0xc3, 0x9a, 0xbe, 0xda, 0xd1, 0x47, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xbb, 0x06, 0xe4, 0xf1, 0x5c, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Queries a list of GetValsetByID items.
	GetValsetByID(ctx context.Context, in *QueryGetValsetByIDRequest, opts ...grpc.CallOption) (*QueryGetValsetByIDResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) GetValsetByID(ctx context.Context, in *QueryGetValsetByIDRequest, opts ...grpc.CallOption) (*QueryGetValsetByIDResponse, error) {
	out := new(QueryGetValsetByIDResponse)
	err := c.cc.Invoke(ctx, "/palomachain.paloma.evm.Query/GetValsetByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Queries a list of GetValsetByID items.
	GetValsetByID(context.Context, *QueryGetValsetByIDRequest) (*QueryGetValsetByIDResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) GetValsetByID(ctx context.Context, req *QueryGetValsetByIDRequest) (*QueryGetValsetByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetValsetByID not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_GetValsetByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetValsetByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetValsetByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/palomachain.paloma.evm.Query/GetValsetByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetValsetByID(ctx, req.(*QueryGetValsetByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "palomachain.paloma.evm.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetValsetByID",
			Handler:    _Query_GetValsetByID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "paloma/evm/query.proto",
}

func (m *QueryGetValsetByIDRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetValsetByIDRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetValsetByIDRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ChainID) > 0 {
		i -= len(m.ChainID)
		copy(dAtA[i:], m.ChainID)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.ChainID)))
		i--
		dAtA[i] = 0x12
	}
	if m.ValsetID != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.ValsetID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryGetValsetByIDResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetValsetByIDResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetValsetByIDResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Valset != nil {
		{
			size, err := m.Valset.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryGetValsetByIDRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ValsetID != 0 {
		n += 1 + sovQuery(uint64(m.ValsetID))
	}
	l = len(m.ChainID)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryGetValsetByIDResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Valset != nil {
		l = m.Valset.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryGetValsetByIDRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryGetValsetByIDRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetValsetByIDRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValsetID", wireType)
			}
			m.ValsetID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ValsetID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChainID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryGetValsetByIDResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryGetValsetByIDResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetValsetByIDResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Valset", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Valset == nil {
				m.Valset = &Valset{}
			}
			if err := m.Valset.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)