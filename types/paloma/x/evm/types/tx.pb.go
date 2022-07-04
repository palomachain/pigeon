// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: paloma/evm/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
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

type MsgSubmitNewJob struct {
	Creator                 string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	HexSmartContractAddress string `protobuf:"bytes,2,opt,name=hexSmartContractAddress,proto3" json:"hexSmartContractAddress,omitempty"`
	HexPayload              string `protobuf:"bytes,3,opt,name=hexPayload,proto3" json:"hexPayload,omitempty"`
	Abi                     string `protobuf:"bytes,4,opt,name=abi,proto3" json:"abi,omitempty"`
	Method                  string `protobuf:"bytes,5,opt,name=method,proto3" json:"method,omitempty"`
	ChainType               string `protobuf:"bytes,6,opt,name=chainType,proto3" json:"chainType,omitempty"`
	ChainID                 string `protobuf:"bytes,7,opt,name=chainID,proto3" json:"chainID,omitempty"`
}

func (m *MsgSubmitNewJob) Reset()         { *m = MsgSubmitNewJob{} }
func (m *MsgSubmitNewJob) String() string { return proto.CompactTextString(m) }
func (*MsgSubmitNewJob) ProtoMessage()    {}
func (*MsgSubmitNewJob) Descriptor() ([]byte, []int) {
	return fileDescriptor_42076d7251ea6e49, []int{0}
}
func (m *MsgSubmitNewJob) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSubmitNewJob) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSubmitNewJob.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSubmitNewJob) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSubmitNewJob.Merge(m, src)
}
func (m *MsgSubmitNewJob) XXX_Size() int {
	return m.Size()
}
func (m *MsgSubmitNewJob) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSubmitNewJob.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSubmitNewJob proto.InternalMessageInfo

func (m *MsgSubmitNewJob) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *MsgSubmitNewJob) GetHexSmartContractAddress() string {
	if m != nil {
		return m.HexSmartContractAddress
	}
	return ""
}

func (m *MsgSubmitNewJob) GetHexPayload() string {
	if m != nil {
		return m.HexPayload
	}
	return ""
}

func (m *MsgSubmitNewJob) GetAbi() string {
	if m != nil {
		return m.Abi
	}
	return ""
}

func (m *MsgSubmitNewJob) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *MsgSubmitNewJob) GetChainType() string {
	if m != nil {
		return m.ChainType
	}
	return ""
}

func (m *MsgSubmitNewJob) GetChainID() string {
	if m != nil {
		return m.ChainID
	}
	return ""
}

type MsgSubmitNewJobResponse struct {
}

func (m *MsgSubmitNewJobResponse) Reset()         { *m = MsgSubmitNewJobResponse{} }
func (m *MsgSubmitNewJobResponse) String() string { return proto.CompactTextString(m) }
func (*MsgSubmitNewJobResponse) ProtoMessage()    {}
func (*MsgSubmitNewJobResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_42076d7251ea6e49, []int{1}
}
func (m *MsgSubmitNewJobResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSubmitNewJobResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSubmitNewJobResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSubmitNewJobResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSubmitNewJobResponse.Merge(m, src)
}
func (m *MsgSubmitNewJobResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgSubmitNewJobResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSubmitNewJobResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSubmitNewJobResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgSubmitNewJob)(nil), "palomachain.paloma.evm.MsgSubmitNewJob")
	proto.RegisterType((*MsgSubmitNewJobResponse)(nil), "palomachain.paloma.evm.MsgSubmitNewJobResponse")
}

func init() { proto.RegisterFile("paloma/evm/tx.proto", fileDescriptor_42076d7251ea6e49) }

var fileDescriptor_42076d7251ea6e49 = []byte{
	// 303 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x51, 0xcf, 0x4a, 0x02, 0x41,
	0x18, 0x77, 0xb3, 0x14, 0x3f, 0x82, 0x62, 0x02, 0x9d, 0x22, 0x86, 0xf0, 0x52, 0x5d, 0x66, 0xa1,
	0x2e, 0x5d, 0xcb, 0x2e, 0x05, 0x46, 0x68, 0xa7, 0x6e, 0xb3, 0xbb, 0x1f, 0xce, 0x82, 0xe3, 0x2c,
	0x33, 0xa3, 0xad, 0x6f, 0xd1, 0x63, 0x75, 0xf4, 0xd8, 0x31, 0xf4, 0xd0, 0x6b, 0x84, 0xb3, 0x2b,
	0x99, 0x14, 0x74, 0xfb, 0xfd, 0xe3, 0xfb, 0x0b, 0x07, 0x99, 0x18, 0x6a, 0x25, 0x42, 0x9c, 0xa8,
	0xd0, 0xe5, 0x3c, 0x33, 0xda, 0x69, 0xd2, 0x2c, 0xc4, 0x58, 0x8a, 0x74, 0xc4, 0x0b, 0xcc, 0x71,
	0xa2, 0xda, 0x9f, 0x01, 0xec, 0x75, 0xed, 0xa0, 0x3f, 0x8e, 0x54, 0xea, 0x1e, 0xf0, 0xe5, 0x5e,
	0x47, 0x84, 0x42, 0x3d, 0x36, 0x28, 0x9c, 0x36, 0x34, 0x38, 0x09, 0xce, 0x1a, 0xbd, 0x15, 0x25,
	0x57, 0xd0, 0x92, 0x98, 0xf7, 0x95, 0x30, 0xae, 0xa3, 0x47, 0xce, 0x88, 0xd8, 0x5d, 0x27, 0x89,
	0x41, 0x6b, 0xe9, 0x96, 0x4f, 0xfe, 0x65, 0x13, 0x06, 0x20, 0x31, 0x7f, 0x14, 0xd3, 0xa1, 0x16,
	0x09, 0xad, 0xfa, 0xf0, 0x9a, 0x42, 0xf6, 0xa1, 0x2a, 0xa2, 0x94, 0x6e, 0x7b, 0x63, 0x09, 0x49,
	0x13, 0x6a, 0x0a, 0x9d, 0xd4, 0x09, 0xdd, 0xf1, 0x62, 0xc9, 0xc8, 0x31, 0x34, 0xfc, 0x16, 0x4f,
	0xd3, 0x0c, 0x69, 0xcd, 0x5b, 0xdf, 0x82, 0x9f, 0x7d, 0x49, 0xee, 0x6e, 0x69, 0xbd, 0x9c, 0xbd,
	0xa0, 0xed, 0x43, 0x68, 0x6d, 0x2c, 0xda, 0x43, 0x9b, 0xe9, 0x91, 0xc5, 0x0b, 0x0d, 0xd5, 0xae,
	0x1d, 0x10, 0x09, 0xbb, 0x3f, 0xee, 0x70, 0xca, 0x7f, 0x3f, 0x1a, 0xdf, 0xa8, 0x73, 0x14, 0xfe,
	0x33, 0xb8, 0x6a, 0x78, 0xd3, 0x79, 0x9b, 0xb3, 0x60, 0x36, 0x67, 0xc1, 0xc7, 0x9c, 0x05, 0xaf,
	0x0b, 0x56, 0x99, 0x2d, 0x58, 0xe5, 0x7d, 0xc1, 0x2a, 0xcf, 0xe7, 0x83, 0xd4, 0xc9, 0x71, 0xc4,
	0x63, 0xad, 0xc2, 0xb5, 0xa2, 0x25, 0x0e, 0xf3, 0xe2, 0xab, 0xd3, 0x0c, 0x6d, 0x54, 0xf3, 0x9f,
	0xbd, 0xfc, 0x0a, 0x00, 0x00, 0xff, 0xff, 0x02, 0x51, 0x75, 0x2d, 0xf0, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	SubmitNewJob(ctx context.Context, in *MsgSubmitNewJob, opts ...grpc.CallOption) (*MsgSubmitNewJobResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) SubmitNewJob(ctx context.Context, in *MsgSubmitNewJob, opts ...grpc.CallOption) (*MsgSubmitNewJobResponse, error) {
	out := new(MsgSubmitNewJobResponse)
	err := c.cc.Invoke(ctx, "/palomachain.paloma.evm.Msg/SubmitNewJob", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	SubmitNewJob(context.Context, *MsgSubmitNewJob) (*MsgSubmitNewJobResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) SubmitNewJob(ctx context.Context, req *MsgSubmitNewJob) (*MsgSubmitNewJobResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitNewJob not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_SubmitNewJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSubmitNewJob)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SubmitNewJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/palomachain.paloma.evm.Msg/SubmitNewJob",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SubmitNewJob(ctx, req.(*MsgSubmitNewJob))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "palomachain.paloma.evm.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmitNewJob",
			Handler:    _Msg_SubmitNewJob_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "paloma/evm/tx.proto",
}

func (m *MsgSubmitNewJob) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSubmitNewJob) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSubmitNewJob) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ChainID) > 0 {
		i -= len(m.ChainID)
		copy(dAtA[i:], m.ChainID)
		i = encodeVarintTx(dAtA, i, uint64(len(m.ChainID)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.ChainType) > 0 {
		i -= len(m.ChainType)
		copy(dAtA[i:], m.ChainType)
		i = encodeVarintTx(dAtA, i, uint64(len(m.ChainType)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.Method) > 0 {
		i -= len(m.Method)
		copy(dAtA[i:], m.Method)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Method)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Abi) > 0 {
		i -= len(m.Abi)
		copy(dAtA[i:], m.Abi)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Abi)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.HexPayload) > 0 {
		i -= len(m.HexPayload)
		copy(dAtA[i:], m.HexPayload)
		i = encodeVarintTx(dAtA, i, uint64(len(m.HexPayload)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.HexSmartContractAddress) > 0 {
		i -= len(m.HexSmartContractAddress)
		copy(dAtA[i:], m.HexSmartContractAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.HexSmartContractAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgSubmitNewJobResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSubmitNewJobResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSubmitNewJobResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgSubmitNewJob) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.HexSmartContractAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.HexPayload)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Abi)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Method)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.ChainType)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.ChainID)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgSubmitNewJobResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgSubmitNewJob) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgSubmitNewJob: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSubmitNewJob: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HexSmartContractAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HexSmartContractAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HexPayload", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HexPayload = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Abi", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Abi = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Method", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Method = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainType", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChainType = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChainID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgSubmitNewJobResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgSubmitNewJobResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSubmitNewJobResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)