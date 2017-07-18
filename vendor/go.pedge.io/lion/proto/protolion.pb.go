// Code generated by protoc-gen-go.
// source: proto/protolion.proto
// DO NOT EDIT!

package protolion

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

// Level is a logging level.
type Level int32

const (
	Level_LEVEL_DEBUG Level = 0
	Level_LEVEL_INFO  Level = 1
	Level_LEVEL_WARN  Level = 2
	Level_LEVEL_ERROR Level = 3
	Level_LEVEL_FATAL Level = 4
	Level_LEVEL_PANIC Level = 5
	Level_LEVEL_NONE  Level = 6
)

var Level_name = map[int32]string{
	0: "LEVEL_DEBUG",
	1: "LEVEL_INFO",
	2: "LEVEL_WARN",
	3: "LEVEL_ERROR",
	4: "LEVEL_FATAL",
	5: "LEVEL_PANIC",
	6: "LEVEL_NONE",
}
var Level_value = map[string]int32{
	"LEVEL_DEBUG": 0,
	"LEVEL_INFO":  1,
	"LEVEL_WARN":  2,
	"LEVEL_ERROR": 3,
	"LEVEL_FATAL": 4,
	"LEVEL_PANIC": 5,
	"LEVEL_NONE":  6,
}

func (x Level) String() string {
	return proto.EnumName(Level_name, int32(x))
}
func (Level) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Timestamp struct {
	Seconds int64 `protobuf:"varint,1,opt,name=seconds" json:"seconds,omitempty"`
	Nanos   int32 `protobuf:"varint,2,opt,name=nanos" json:"nanos,omitempty"`
}

func (m *Timestamp) Reset()                    { *m = Timestamp{} }
func (m *Timestamp) String() string            { return proto.CompactTextString(m) }
func (*Timestamp) ProtoMessage()               {}
func (*Timestamp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Timestamp) GetSeconds() int64 {
	if m != nil {
		return m.Seconds
	}
	return 0
}

func (m *Timestamp) GetNanos() int32 {
	if m != nil {
		return m.Nanos
	}
	return 0
}

// Entry is the object serialized for logging.
type Entry struct {
	// id may not be set depending on logger options
	// it is up to the user to determine if id is required
	Id        string     `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Level     Level      `protobuf:"varint,2,opt,name=level,enum=lion.Level" json:"level,omitempty"`
	Timestamp *Timestamp `protobuf:"bytes,3,opt,name=timestamp" json:"timestamp,omitempty"`
	// both context and fields may be set
	Context []*Entry_Message  `protobuf:"bytes,4,rep,name=context" json:"context,omitempty"`
	Fields  map[string]string `protobuf:"bytes,5,rep,name=fields" json:"fields,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// one of event, message, writer_output will be set
	Event        *Entry_Message `protobuf:"bytes,6,opt,name=event" json:"event,omitempty"`
	Message      string         `protobuf:"bytes,7,opt,name=message" json:"message,omitempty"`
	WriterOutput []byte         `protobuf:"bytes,8,opt,name=writer_output,json=writerOutput,proto3" json:"writer_output,omitempty"`
}

func (m *Entry) Reset()                    { *m = Entry{} }
func (m *Entry) String() string            { return proto.CompactTextString(m) }
func (*Entry) ProtoMessage()               {}
func (*Entry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Entry) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Entry) GetLevel() Level {
	if m != nil {
		return m.Level
	}
	return Level_LEVEL_DEBUG
}

func (m *Entry) GetTimestamp() *Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *Entry) GetContext() []*Entry_Message {
	if m != nil {
		return m.Context
	}
	return nil
}

func (m *Entry) GetFields() map[string]string {
	if m != nil {
		return m.Fields
	}
	return nil
}

func (m *Entry) GetEvent() *Entry_Message {
	if m != nil {
		return m.Event
	}
	return nil
}

func (m *Entry) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *Entry) GetWriterOutput() []byte {
	if m != nil {
		return m.WriterOutput
	}
	return nil
}

// Message represents a serialized protobuf message.
// The name is the name registered with lion.
type Entry_Message struct {
	Encoding string `protobuf:"bytes,1,opt,name=encoding" json:"encoding,omitempty"`
	Name     string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Value    []byte `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *Entry_Message) Reset()                    { *m = Entry_Message{} }
func (m *Entry_Message) String() string            { return proto.CompactTextString(m) }
func (*Entry_Message) ProtoMessage()               {}
func (*Entry_Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

func (m *Entry_Message) GetEncoding() string {
	if m != nil {
		return m.Encoding
	}
	return ""
}

func (m *Entry_Message) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Entry_Message) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func init() {
	proto.RegisterType((*Timestamp)(nil), "lion.Timestamp")
	proto.RegisterType((*Entry)(nil), "lion.Entry")
	proto.RegisterType((*Entry_Message)(nil), "lion.Entry.Message")
	proto.RegisterEnum("lion.Level", Level_name, Level_value)
}

func init() { proto.RegisterFile("proto/protolion.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 418 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x52, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xc5, 0x76, 0x9c, 0xd4, 0xe3, 0x90, 0xae, 0x06, 0x10, 0xab, 0x9c, 0x4c, 0xb9, 0x18, 0xa4,
	0xa6, 0x52, 0xb8, 0xf0, 0x71, 0x4a, 0xc1, 0x41, 0x95, 0x82, 0x8d, 0x56, 0x05, 0x24, 0x2e, 0x95,
	0x89, 0x97, 0xca, 0xc2, 0xd9, 0x8d, 0xe2, 0x4d, 0xa0, 0xe2, 0xc7, 0xf1, 0xd7, 0x90, 0x67, 0xe3,
	0x38, 0x87, 0x5e, 0xac, 0xf7, 0xde, 0xbe, 0x79, 0xb3, 0x3b, 0x63, 0x78, 0xb2, 0xde, 0x68, 0xa3,
	0x2f, 0xe8, 0x5b, 0x95, 0x5a, 0x4d, 0x08, 0x61, 0xaf, 0xc1, 0x67, 0xef, 0x20, 0xb8, 0x2e, 0x57,
	0xb2, 0x36, 0xf9, 0x6a, 0x8d, 0x1c, 0x06, 0xb5, 0x5c, 0x6a, 0x55, 0xd4, 0xdc, 0x89, 0x9c, 0xd8,
	0x13, 0x2d, 0xc5, 0xc7, 0xe0, 0xab, 0x5c, 0xe9, 0x9a, 0xbb, 0x91, 0x13, 0xfb, 0xc2, 0x92, 0xb3,
	0x7f, 0x1e, 0xf8, 0x89, 0x32, 0x9b, 0x3b, 0x1c, 0x81, 0x5b, 0x16, 0x54, 0x14, 0x08, 0xb7, 0x2c,
	0xf0, 0x19, 0xf8, 0x95, 0xdc, 0xc9, 0x8a, 0xfc, 0xa3, 0x69, 0x38, 0xa1, 0xc6, 0x8b, 0x46, 0x12,
	0xf6, 0x04, 0xcf, 0x21, 0x30, 0x6d, 0x67, 0xee, 0x45, 0x4e, 0x1c, 0x4e, 0x4f, 0xad, 0xed, 0x70,
	0x21, 0xd1, 0x39, 0xf0, 0x1c, 0x06, 0x4b, 0xad, 0x8c, 0xfc, 0x63, 0x78, 0x2f, 0xf2, 0xe2, 0x70,
	0xfa, 0xc8, 0x9a, 0xa9, 0xff, 0xe4, 0x93, 0xac, 0xeb, 0xfc, 0x56, 0x8a, 0xd6, 0x83, 0x17, 0xd0,
	0xff, 0x59, 0xca, 0xaa, 0xa8, 0xb9, 0x4f, 0xee, 0xa7, 0xc7, 0xee, 0x39, 0x9d, 0x10, 0x16, 0x7b,
	0x1b, 0xbe, 0x00, 0x5f, 0xee, 0xa4, 0x32, 0xbc, 0x4f, 0x57, 0xb9, 0x37, 0xdd, 0x3a, 0x9a, 0x31,
	0xad, 0xac, 0xc2, 0x07, 0xf4, 0xe2, 0x96, 0xe2, 0x73, 0x78, 0xf8, 0x7b, 0x53, 0x1a, 0xb9, 0xb9,
	0xd1, 0x5b, 0xb3, 0xde, 0x1a, 0x7e, 0x12, 0x39, 0xf1, 0x50, 0x0c, 0xad, 0x98, 0x91, 0x36, 0xce,
	0x60, 0xb0, 0x0f, 0xc4, 0x31, 0x9c, 0x48, 0xb5, 0xd4, 0x45, 0xa9, 0x6e, 0xf7, 0xc3, 0x3b, 0x70,
	0x44, 0xe8, 0xa9, 0x7c, 0x25, 0x69, 0x82, 0x81, 0x20, 0xdc, 0xac, 0x61, 0x97, 0x57, 0x5b, 0x49,
	0xf3, 0x1a, 0x0a, 0x4b, 0xc6, 0x6f, 0x20, 0x3c, 0x7a, 0x11, 0x32, 0xf0, 0x7e, 0xc9, 0xbb, 0x7d,
	0x5e, 0x03, 0xbb, 0x32, 0x9b, 0x65, 0xc9, 0x5b, 0xf7, 0xb5, 0xf3, 0xf2, 0x2f, 0xf8, 0xb4, 0x14,
	0x3c, 0x85, 0x70, 0x91, 0x7c, 0x4d, 0x16, 0x37, 0x1f, 0x92, 0xcb, 0x2f, 0x1f, 0xd9, 0x03, 0x1c,
	0x01, 0x58, 0xe1, 0x2a, 0x9d, 0x67, 0xcc, 0xe9, 0xf8, 0xb7, 0x99, 0x48, 0x99, 0xdb, 0x15, 0x24,
	0x42, 0x64, 0x82, 0x79, 0x9d, 0x30, 0x9f, 0x5d, 0xcf, 0x16, 0xac, 0xd7, 0x09, 0x9f, 0x67, 0xe9,
	0xd5, 0x7b, 0xe6, 0x77, 0x11, 0x69, 0x96, 0x26, 0xac, 0x7f, 0x19, 0x7e, 0x0f, 0x0e, 0x3f, 0xe5,
	0x8f, 0x3e, 0xc1, 0x57, 0xff, 0x03, 0x00, 0x00, 0xff, 0xff, 0x74, 0x2f, 0xc7, 0x8c, 0xae, 0x02,
	0x00, 0x00,
}
