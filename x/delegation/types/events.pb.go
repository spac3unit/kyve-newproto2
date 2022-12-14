// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kyve/delegation/v1beta1/events.proto

package types

import (
	fmt "fmt"
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

// EventDelegate is an event emitted when someone delegates to a protocol node.
type EventDelegate struct {
	// address is the account address of the delegator.
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// node is the account address of the protocol node.
	Node string `protobuf:"bytes,2,opt,name=node,proto3" json:"node,omitempty"`
	// amount ...
	Amount uint64 `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (m *EventDelegate) Reset()         { *m = EventDelegate{} }
func (m *EventDelegate) String() string { return proto.CompactTextString(m) }
func (*EventDelegate) ProtoMessage()    {}
func (*EventDelegate) Descriptor() ([]byte, []int) {
	return fileDescriptor_d01988a9108a2e89, []int{0}
}
func (m *EventDelegate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventDelegate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventDelegate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventDelegate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventDelegate.Merge(m, src)
}
func (m *EventDelegate) XXX_Size() int {
	return m.Size()
}
func (m *EventDelegate) XXX_DiscardUnknown() {
	xxx_messageInfo_EventDelegate.DiscardUnknown(m)
}

var xxx_messageInfo_EventDelegate proto.InternalMessageInfo

func (m *EventDelegate) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *EventDelegate) GetNode() string {
	if m != nil {
		return m.Node
	}
	return ""
}

func (m *EventDelegate) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

// EventUndelegate is an event emitted when someone undelegates from a protocol node.
type EventUndelegate struct {
	// address is the account address of the delegator.
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// node is the account address of the protocol node.
	Node string `protobuf:"bytes,2,opt,name=node,proto3" json:"node,omitempty"`
	// amount ...
	Amount uint64 `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (m *EventUndelegate) Reset()         { *m = EventUndelegate{} }
func (m *EventUndelegate) String() string { return proto.CompactTextString(m) }
func (*EventUndelegate) ProtoMessage()    {}
func (*EventUndelegate) Descriptor() ([]byte, []int) {
	return fileDescriptor_d01988a9108a2e89, []int{1}
}
func (m *EventUndelegate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventUndelegate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventUndelegate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventUndelegate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventUndelegate.Merge(m, src)
}
func (m *EventUndelegate) XXX_Size() int {
	return m.Size()
}
func (m *EventUndelegate) XXX_DiscardUnknown() {
	xxx_messageInfo_EventUndelegate.DiscardUnknown(m)
}

var xxx_messageInfo_EventUndelegate proto.InternalMessageInfo

func (m *EventUndelegate) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *EventUndelegate) GetNode() string {
	if m != nil {
		return m.Node
	}
	return ""
}

func (m *EventUndelegate) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

// EventRedelegate is an event emitted when someone redelegates from one protocol node to another.
type EventRedelegate struct {
	// address is the account address of the delegator.
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// from_node ...
	FromNode string `protobuf:"bytes,2,opt,name=from_node,json=fromNode,proto3" json:"from_node,omitempty"`
	// address is the account address of the new staker in the the pool
	ToNode string `protobuf:"bytes,3,opt,name=to_node,json=toNode,proto3" json:"to_node,omitempty"`
	// amount ...
	Amount uint64 `protobuf:"varint,4,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (m *EventRedelegate) Reset()         { *m = EventRedelegate{} }
func (m *EventRedelegate) String() string { return proto.CompactTextString(m) }
func (*EventRedelegate) ProtoMessage()    {}
func (*EventRedelegate) Descriptor() ([]byte, []int) {
	return fileDescriptor_d01988a9108a2e89, []int{2}
}
func (m *EventRedelegate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventRedelegate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventRedelegate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventRedelegate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventRedelegate.Merge(m, src)
}
func (m *EventRedelegate) XXX_Size() int {
	return m.Size()
}
func (m *EventRedelegate) XXX_DiscardUnknown() {
	xxx_messageInfo_EventRedelegate.DiscardUnknown(m)
}

var xxx_messageInfo_EventRedelegate proto.InternalMessageInfo

func (m *EventRedelegate) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *EventRedelegate) GetFromNode() string {
	if m != nil {
		return m.FromNode
	}
	return ""
}

func (m *EventRedelegate) GetToNode() string {
	if m != nil {
		return m.ToNode
	}
	return ""
}

func (m *EventRedelegate) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

// EventWithdrawRewards ...
type EventWithdrawRewards struct {
	// address is the account address of the delegator.
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// from_node is the account address of the protocol node the users withdraws from.
	FromNode string `protobuf:"bytes,2,opt,name=from_node,json=fromNode,proto3" json:"from_node,omitempty"`
	// amount ...
	Amount uint64 `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (m *EventWithdrawRewards) Reset()         { *m = EventWithdrawRewards{} }
func (m *EventWithdrawRewards) String() string { return proto.CompactTextString(m) }
func (*EventWithdrawRewards) ProtoMessage()    {}
func (*EventWithdrawRewards) Descriptor() ([]byte, []int) {
	return fileDescriptor_d01988a9108a2e89, []int{3}
}
func (m *EventWithdrawRewards) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventWithdrawRewards) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventWithdrawRewards.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventWithdrawRewards) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventWithdrawRewards.Merge(m, src)
}
func (m *EventWithdrawRewards) XXX_Size() int {
	return m.Size()
}
func (m *EventWithdrawRewards) XXX_DiscardUnknown() {
	xxx_messageInfo_EventWithdrawRewards.DiscardUnknown(m)
}

var xxx_messageInfo_EventWithdrawRewards proto.InternalMessageInfo

func (m *EventWithdrawRewards) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *EventWithdrawRewards) GetFromNode() string {
	if m != nil {
		return m.FromNode
	}
	return ""
}

func (m *EventWithdrawRewards) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func init() {
	proto.RegisterType((*EventDelegate)(nil), "kyve.delegation.v1beta1.EventDelegate")
	proto.RegisterType((*EventUndelegate)(nil), "kyve.delegation.v1beta1.EventUndelegate")
	proto.RegisterType((*EventRedelegate)(nil), "kyve.delegation.v1beta1.EventRedelegate")
	proto.RegisterType((*EventWithdrawRewards)(nil), "kyve.delegation.v1beta1.EventWithdrawRewards")
}

func init() {
	proto.RegisterFile("kyve/delegation/v1beta1/events.proto", fileDescriptor_d01988a9108a2e89)
}

var fileDescriptor_d01988a9108a2e89 = []byte{
	// 288 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0xc9, 0xae, 0x2c, 0x4b,
	0xd5, 0x4f, 0x49, 0xcd, 0x49, 0x4d, 0x4f, 0x2c, 0xc9, 0xcc, 0xcf, 0xd3, 0x2f, 0x33, 0x4c, 0x4a,
	0x2d, 0x49, 0x34, 0xd4, 0x4f, 0x2d, 0x4b, 0xcd, 0x2b, 0x29, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x12, 0x07, 0xa9, 0xd2, 0x43, 0xa8, 0xd2, 0x83, 0xaa, 0x52, 0x0a, 0xe5, 0xe2, 0x75, 0x05,
	0x29, 0x74, 0x81, 0x48, 0xa5, 0x0a, 0x49, 0x70, 0xb1, 0x27, 0xa6, 0xa4, 0x14, 0xa5, 0x16, 0x17,
	0x4b, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0xc1, 0xb8, 0x42, 0x42, 0x5c, 0x2c, 0x79, 0xf9, 0x29,
	0xa9, 0x12, 0x4c, 0x60, 0x61, 0x30, 0x5b, 0x48, 0x8c, 0x8b, 0x2d, 0x31, 0x37, 0xbf, 0x34, 0xaf,
	0x44, 0x82, 0x59, 0x81, 0x51, 0x83, 0x25, 0x08, 0xca, 0x53, 0x0a, 0xe7, 0xe2, 0x07, 0x1b, 0x1b,
	0x9a, 0x97, 0x42, 0x5d, 0x83, 0x2b, 0xa1, 0x06, 0x07, 0xa5, 0x12, 0x61, 0xb0, 0x34, 0x17, 0x67,
	0x5a, 0x51, 0x7e, 0x6e, 0x3c, 0x92, 0xe9, 0x1c, 0x20, 0x01, 0x3f, 0x90, 0x0d, 0xe2, 0x5c, 0xec,
	0x25, 0xf9, 0x10, 0x29, 0x66, 0xb0, 0x14, 0x5b, 0x49, 0xbe, 0x1f, 0xaa, 0xd5, 0x2c, 0x28, 0x56,
	0xa7, 0x72, 0x89, 0x80, 0xad, 0x0e, 0xcf, 0x2c, 0xc9, 0x48, 0x29, 0x4a, 0x2c, 0x0f, 0x4a, 0x2d,
	0x4f, 0x2c, 0x4a, 0x29, 0x26, 0xd7, 0x7e, 0x1c, 0x3e, 0x74, 0xf2, 0x3c, 0xf1, 0x48, 0x8e, 0xf1,
	0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e,
	0xe1, 0xc6, 0x63, 0x39, 0x86, 0x28, 0xfd, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc,
	0x5c, 0x7d, 0xef, 0xc8, 0x30, 0x57, 0xbf, 0xd4, 0x92, 0xf2, 0xfc, 0xa2, 0x6c, 0xfd, 0xe4, 0x8c,
	0xc4, 0xcc, 0x3c, 0xfd, 0x0a, 0xe4, 0x44, 0x50, 0x52, 0x59, 0x90, 0x5a, 0x9c, 0xc4, 0x06, 0x8e,
	0x7c, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x73, 0xba, 0x55, 0xe1, 0x24, 0x02, 0x00, 0x00,
}

func (m *EventDelegate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventDelegate) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventDelegate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Amount != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Amount))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Node) > 0 {
		i -= len(m.Node)
		copy(dAtA[i:], m.Node)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Node)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *EventUndelegate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventUndelegate) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventUndelegate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Amount != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Amount))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Node) > 0 {
		i -= len(m.Node)
		copy(dAtA[i:], m.Node)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Node)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *EventRedelegate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventRedelegate) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventRedelegate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Amount != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Amount))
		i--
		dAtA[i] = 0x20
	}
	if len(m.ToNode) > 0 {
		i -= len(m.ToNode)
		copy(dAtA[i:], m.ToNode)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.ToNode)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.FromNode) > 0 {
		i -= len(m.FromNode)
		copy(dAtA[i:], m.FromNode)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.FromNode)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *EventWithdrawRewards) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventWithdrawRewards) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventWithdrawRewards) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Amount != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Amount))
		i--
		dAtA[i] = 0x18
	}
	if len(m.FromNode) > 0 {
		i -= len(m.FromNode)
		copy(dAtA[i:], m.FromNode)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.FromNode)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintEvents(dAtA []byte, offset int, v uint64) int {
	offset -= sovEvents(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *EventDelegate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.Node)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	if m.Amount != 0 {
		n += 1 + sovEvents(uint64(m.Amount))
	}
	return n
}

func (m *EventUndelegate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.Node)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	if m.Amount != 0 {
		n += 1 + sovEvents(uint64(m.Amount))
	}
	return n
}

func (m *EventRedelegate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.FromNode)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.ToNode)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	if m.Amount != 0 {
		n += 1 + sovEvents(uint64(m.Amount))
	}
	return n
}

func (m *EventWithdrawRewards) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.FromNode)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	if m.Amount != 0 {
		n += 1 + sovEvents(uint64(m.Amount))
	}
	return n
}

func sovEvents(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvents(x uint64) (n int) {
	return sovEvents(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *EventDelegate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
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
			return fmt.Errorf("proto: EventDelegate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventDelegate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Node", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Node = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			m.Amount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Amount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvents
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
func (m *EventUndelegate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
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
			return fmt.Errorf("proto: EventUndelegate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventUndelegate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Node", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Node = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			m.Amount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Amount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvents
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
func (m *EventRedelegate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
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
			return fmt.Errorf("proto: EventRedelegate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventRedelegate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FromNode", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FromNode = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ToNode", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ToNode = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			m.Amount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Amount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvents
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
func (m *EventWithdrawRewards) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
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
			return fmt.Errorf("proto: EventWithdrawRewards: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventWithdrawRewards: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FromNode", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FromNode = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			m.Amount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Amount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvents
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
func skipEvents(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvents
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
					return 0, ErrIntOverflowEvents
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
					return 0, ErrIntOverflowEvents
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
				return 0, ErrInvalidLengthEvents
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEvents
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEvents
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEvents        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvents          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEvents = fmt.Errorf("proto: unexpected end of group")
)
