// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: paloma/evm/common.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

type ArbitrarySmartContractCall struct {
	Method     string `protobuf:"bytes,1,opt,name=method,proto3" json:"method,omitempty"`
	Payload    []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	HexAddress string `protobuf:"bytes,3,opt,name=hexAddress,proto3" json:"hexAddress,omitempty"`
	Abi        []byte `protobuf:"bytes,4,opt,name=abi,proto3" json:"abi,omitempty"`
}

func (m *ArbitrarySmartContractCall) Reset()         { *m = ArbitrarySmartContractCall{} }
func (m *ArbitrarySmartContractCall) String() string { return proto.CompactTextString(m) }
func (*ArbitrarySmartContractCall) ProtoMessage()    {}
func (*ArbitrarySmartContractCall) Descriptor() ([]byte, []int) {
	return fileDescriptor_86b285e01d9df701, []int{0}
}
func (m *ArbitrarySmartContractCall) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ArbitrarySmartContractCall) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ArbitrarySmartContractCall.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ArbitrarySmartContractCall) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ArbitrarySmartContractCall.Merge(m, src)
}
func (m *ArbitrarySmartContractCall) XXX_Size() int {
	return m.Size()
}
func (m *ArbitrarySmartContractCall) XXX_DiscardUnknown() {
	xxx_messageInfo_ArbitrarySmartContractCall.DiscardUnknown(m)
}

var xxx_messageInfo_ArbitrarySmartContractCall proto.InternalMessageInfo

func (m *ArbitrarySmartContractCall) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *ArbitrarySmartContractCall) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *ArbitrarySmartContractCall) GetHexAddress() string {
	if m != nil {
		return m.HexAddress
	}
	return ""
}

func (m *ArbitrarySmartContractCall) GetAbi() []byte {
	if m != nil {
		return m.Abi
	}
	return nil
}

func init() {
	proto.RegisterType((*ArbitrarySmartContractCall)(nil), "palomachain.paloma.evm.ArbitrarySmartContractCall")
}

func init() { proto.RegisterFile("paloma/evm/common.proto", fileDescriptor_86b285e01d9df701) }

var fileDescriptor_86b285e01d9df701 = []byte{
	// 237 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8f, 0x31, 0x4e, 0xc3, 0x40,
	0x10, 0x45, 0xbd, 0x04, 0x05, 0xb1, 0xa2, 0x40, 0x2b, 0x14, 0x56, 0x29, 0x56, 0x11, 0x55, 0x68,
	0xbc, 0x05, 0x27, 0x08, 0xbe, 0x41, 0xe8, 0xe8, 0xc6, 0xf6, 0xca, 0xb6, 0xe4, 0xf1, 0x58, 0xeb,
	0x21, 0x8a, 0x3b, 0x8e, 0xc0, 0xb1, 0x28, 0x53, 0x52, 0x22, 0xfb, 0x22, 0x28, 0xeb, 0x20, 0xb9,
	0x7b, 0x7f, 0xe6, 0xff, 0xe2, 0xc9, 0xc7, 0x16, 0x6a, 0x42, 0xb0, 0xee, 0x80, 0x36, 0x23, 0x44,
	0x6a, 0xe2, 0xd6, 0x13, 0x93, 0x5a, 0x4d, 0x8f, 0xac, 0x84, 0xaa, 0x89, 0x27, 0x8e, 0xdd, 0x01,
	0xd7, 0x0f, 0x05, 0x15, 0x14, 0x2a, 0xf6, 0x4c, 0x53, 0xfb, 0xe9, 0x53, 0xc8, 0xf5, 0xce, 0xa7,
	0x15, 0x7b, 0xf0, 0xfd, 0x1b, 0x82, 0xe7, 0x84, 0x1a, 0xf6, 0x90, 0x71, 0x02, 0x75, 0xad, 0x56,
	0x72, 0x89, 0x8e, 0x4b, 0xca, 0xb5, 0xd8, 0x88, 0xed, 0xed, 0xfe, 0x92, 0x94, 0x96, 0x37, 0x2d,
	0xf4, 0x35, 0x41, 0xae, 0xaf, 0x36, 0x62, 0x7b, 0xb7, 0xff, 0x8f, 0xca, 0x48, 0x59, 0xba, 0xe3,
	0x2e, 0xcf, 0xbd, 0xeb, 0x3a, 0xbd, 0x08, 0xab, 0xd9, 0x45, 0xdd, 0xcb, 0x05, 0xa4, 0x95, 0xbe,
	0x0e, 0xab, 0x33, 0xbe, 0x26, 0xdf, 0x83, 0x11, 0xa7, 0xc1, 0x88, 0xdf, 0xc1, 0x88, 0xaf, 0xd1,
	0x44, 0xa7, 0xd1, 0x44, 0x3f, 0xa3, 0x89, 0xde, 0x9f, 0x8b, 0x8a, 0xcb, 0x8f, 0x34, 0xce, 0x08,
	0xed, 0xcc, 0xea, 0xc2, 0xf6, 0x18, 0xe4, 0xb9, 0x6f, 0x5d, 0x97, 0x2e, 0x83, 0xce, 0xcb, 0x5f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x44, 0xd2, 0x59, 0xaa, 0x17, 0x01, 0x00, 0x00,
}

func (m *ArbitrarySmartContractCall) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ArbitrarySmartContractCall) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ArbitrarySmartContractCall) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Abi) > 0 {
		i -= len(m.Abi)
		copy(dAtA[i:], m.Abi)
		i = encodeVarintCommon(dAtA, i, uint64(len(m.Abi)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.HexAddress) > 0 {
		i -= len(m.HexAddress)
		copy(dAtA[i:], m.HexAddress)
		i = encodeVarintCommon(dAtA, i, uint64(len(m.HexAddress)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Payload) > 0 {
		i -= len(m.Payload)
		copy(dAtA[i:], m.Payload)
		i = encodeVarintCommon(dAtA, i, uint64(len(m.Payload)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Method) > 0 {
		i -= len(m.Method)
		copy(dAtA[i:], m.Method)
		i = encodeVarintCommon(dAtA, i, uint64(len(m.Method)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintCommon(dAtA []byte, offset int, v uint64) int {
	offset -= sovCommon(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ArbitrarySmartContractCall) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Method)
	if l > 0 {
		n += 1 + l + sovCommon(uint64(l))
	}
	l = len(m.Payload)
	if l > 0 {
		n += 1 + l + sovCommon(uint64(l))
	}
	l = len(m.HexAddress)
	if l > 0 {
		n += 1 + l + sovCommon(uint64(l))
	}
	l = len(m.Abi)
	if l > 0 {
		n += 1 + l + sovCommon(uint64(l))
	}
	return n
}

func sovCommon(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCommon(x uint64) (n int) {
	return sovCommon(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ArbitrarySmartContractCall) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommon
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
			return fmt.Errorf("proto: ArbitrarySmartContractCall: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ArbitrarySmartContractCall: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Method", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
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
				return ErrInvalidLengthCommon
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommon
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Method = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Payload", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthCommon
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthCommon
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Payload = append(m.Payload[:0], dAtA[iNdEx:postIndex]...)
			if m.Payload == nil {
				m.Payload = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HexAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
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
				return ErrInvalidLengthCommon
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommon
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HexAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Abi", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommon
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthCommon
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthCommon
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Abi = append(m.Abi[:0], dAtA[iNdEx:postIndex]...)
			if m.Abi == nil {
				m.Abi = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommon(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommon
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
func skipCommon(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCommon
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
					return 0, ErrIntOverflowCommon
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
					return 0, ErrIntOverflowCommon
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
				return 0, ErrInvalidLengthCommon
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupCommon
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthCommon
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthCommon        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCommon          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupCommon = fmt.Errorf("proto: unexpected end of group")
)