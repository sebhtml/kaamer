// Code generated by protoc-gen-go. DO NOT EDIT.
// source: kstats.proto

package kvstore

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type KStats struct {
	NumberOfProteins     uint64   `protobuf:"varint,1,opt,name=NumberOfProteins,proto3" json:"NumberOfProteins,omitempty"`
	NumberOfAA           uint64   `protobuf:"varint,2,opt,name=NumberOfAA,proto3" json:"NumberOfAA,omitempty"`
	NumberOfKmers        uint64   `protobuf:"varint,4,opt,name=NumberOfKmers,proto3" json:"NumberOfKmers,omitempty"`
	NumberOfKCombSets    uint64   `protobuf:"varint,5,opt,name=NumberOfKCombSets,proto3" json:"NumberOfKCombSets,omitempty"`
	Features             []string `protobuf:"bytes,6,rep,name=Features,proto3" json:"Features,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KStats) Reset()         { *m = KStats{} }
func (m *KStats) String() string { return proto.CompactTextString(m) }
func (*KStats) ProtoMessage()    {}
func (*KStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_69d4a9d99f3c1d26, []int{0}
}

func (m *KStats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KStats.Unmarshal(m, b)
}
func (m *KStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KStats.Marshal(b, m, deterministic)
}
func (m *KStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KStats.Merge(m, src)
}
func (m *KStats) XXX_Size() int {
	return xxx_messageInfo_KStats.Size(m)
}
func (m *KStats) XXX_DiscardUnknown() {
	xxx_messageInfo_KStats.DiscardUnknown(m)
}

var xxx_messageInfo_KStats proto.InternalMessageInfo

func (m *KStats) GetNumberOfProteins() uint64 {
	if m != nil {
		return m.NumberOfProteins
	}
	return 0
}

func (m *KStats) GetNumberOfAA() uint64 {
	if m != nil {
		return m.NumberOfAA
	}
	return 0
}

func (m *KStats) GetNumberOfKmers() uint64 {
	if m != nil {
		return m.NumberOfKmers
	}
	return 0
}

func (m *KStats) GetNumberOfKCombSets() uint64 {
	if m != nil {
		return m.NumberOfKCombSets
	}
	return 0
}

func (m *KStats) GetFeatures() []string {
	if m != nil {
		return m.Features
	}
	return nil
}

func init() {
	proto.RegisterType((*KStats)(nil), "kvstore.KStats")
}

func init() { proto.RegisterFile("kstats.proto", fileDescriptor_69d4a9d99f3c1d26) }

var fileDescriptor_69d4a9d99f3c1d26 = []byte{
	// 162 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xc9, 0x2e, 0x2e, 0x49,
	0x2c, 0x29, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xcf, 0x2e, 0x2b, 0x2e, 0xc9, 0x2f,
	0x4a, 0x55, 0x3a, 0xc2, 0xc8, 0xc5, 0xe6, 0x1d, 0x0c, 0x92, 0x11, 0xd2, 0xe2, 0x12, 0xf0, 0x2b,
	0xcd, 0x4d, 0x4a, 0x2d, 0xf2, 0x4f, 0x0b, 0x28, 0xca, 0x2f, 0x49, 0xcd, 0xcc, 0x2b, 0x96, 0x60,
	0x54, 0x60, 0xd4, 0x60, 0x09, 0xc2, 0x10, 0x17, 0x92, 0xe3, 0xe2, 0x82, 0x89, 0x39, 0x3a, 0x4a,
	0x30, 0x81, 0x55, 0x21, 0x89, 0x08, 0xa9, 0x70, 0xf1, 0xc2, 0x78, 0xde, 0xb9, 0xa9, 0x45, 0xc5,
	0x12, 0x2c, 0x60, 0x25, 0xa8, 0x82, 0x42, 0x3a, 0x5c, 0x82, 0x70, 0x01, 0xe7, 0xfc, 0xdc, 0xa4,
	0xe0, 0xd4, 0x92, 0x62, 0x09, 0x56, 0xb0, 0x4a, 0x4c, 0x09, 0x21, 0x29, 0x2e, 0x0e, 0xb7, 0xd4,
	0xc4, 0x92, 0xd2, 0xa2, 0xd4, 0x62, 0x09, 0x36, 0x05, 0x66, 0x0d, 0xce, 0x20, 0x38, 0x3f, 0x89,
	0x0d, 0xec, 0x2d, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd7, 0x22, 0x74, 0x33, 0xe6, 0x00,
	0x00, 0x00,
}
