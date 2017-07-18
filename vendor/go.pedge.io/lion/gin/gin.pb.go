// Code generated by protoc-gen-go.
// source: gin/gin.proto
// DO NOT EDIT!

package ginlion

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "go.pedge.io/pb/go/google/protobuf"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Call struct {
	Method      string            `protobuf:"bytes,1,opt,name=method" json:"method,omitempty"`
	Path        string            `protobuf:"bytes,2,opt,name=path" json:"path,omitempty"`
	Query       map[string]string `protobuf:"bytes,3,rep,name=query" json:"query,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	UserAgent   string            `protobuf:"bytes,4,opt,name=user_agent,json=userAgent" json:"user_agent,omitempty"`
	RequestForm map[string]string `protobuf:"bytes,5,rep,name=request_form,json=requestForm" json:"request_form,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	ClientIp    string            `protobuf:"bytes,6,opt,name=client_ip,json=clientIp" json:"client_ip,omitempty"`
	StatusCode  uint32            `protobuf:"varint,7,opt,name=status_code,json=statusCode" json:"status_code,omitempty"`
	Duration    string            `protobuf:"bytes,8,opt,name=duration" json:"duration,omitempty"`
	Error       []string          `protobuf:"bytes,9,rep,name=error" json:"error,omitempty"`
}

func (m *Call) Reset()                    { *m = Call{} }
func (m *Call) String() string            { return proto.CompactTextString(m) }
func (*Call) ProtoMessage()               {}
func (*Call) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Call) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *Call) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Call) GetQuery() map[string]string {
	if m != nil {
		return m.Query
	}
	return nil
}

func (m *Call) GetUserAgent() string {
	if m != nil {
		return m.UserAgent
	}
	return ""
}

func (m *Call) GetRequestForm() map[string]string {
	if m != nil {
		return m.RequestForm
	}
	return nil
}

func (m *Call) GetClientIp() string {
	if m != nil {
		return m.ClientIp
	}
	return ""
}

func (m *Call) GetStatusCode() uint32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *Call) GetDuration() string {
	if m != nil {
		return m.Duration
	}
	return ""
}

func (m *Call) GetError() []string {
	if m != nil {
		return m.Error
	}
	return nil
}

func init() {
	proto.RegisterType((*Call)(nil), "lion.gin.Call")
}

func init() { proto.RegisterFile("gin/gin.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 324 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0x3f, 0x6b, 0xf3, 0x30,
	0x10, 0xc6, 0x71, 0x9c, 0x3f, 0xf6, 0xe5, 0x0d, 0x04, 0xf1, 0x52, 0x54, 0x97, 0x36, 0xa6, 0x93,
	0x27, 0x1b, 0xda, 0x25, 0x74, 0x28, 0x34, 0xa1, 0x85, 0x8e, 0xf5, 0xd8, 0xc5, 0x38, 0xf1, 0xc5,
	0x11, 0x75, 0x24, 0x47, 0x96, 0x0a, 0xf9, 0x56, 0xfd, 0x88, 0x45, 0x52, 0x42, 0x20, 0x5b, 0x37,
	0x3d, 0xbf, 0xbb, 0x7b, 0x8e, 0x7b, 0x04, 0x93, 0x9a, 0xf1, 0xac, 0x66, 0x3c, 0x6d, 0xa5, 0x50,
	0x82, 0x04, 0x0d, 0x13, 0x3c, 0xad, 0x19, 0x8f, 0xee, 0x6a, 0x21, 0xea, 0x06, 0x33, 0xcb, 0x57,
	0x7a, 0x93, 0x55, 0x5a, 0x96, 0xca, 0x54, 0x2d, 0xb9, 0xff, 0xf1, 0xa1, 0xbf, 0x2c, 0x9b, 0x86,
	0x5c, 0xc1, 0x70, 0x87, 0x6a, 0x2b, 0x2a, 0xea, 0xc5, 0x5e, 0x12, 0xe6, 0x47, 0x45, 0x08, 0xf4,
	0xdb, 0x52, 0x6d, 0x69, 0xcf, 0x52, 0xfb, 0x26, 0x19, 0x0c, 0xf6, 0x1a, 0xe5, 0x81, 0xfa, 0xb1,
	0x9f, 0x8c, 0x1f, 0xae, 0xd3, 0xd3, 0xba, 0xd4, 0x58, 0xa5, 0x1f, 0xa6, 0xf6, 0xca, 0x95, 0x3c,
	0xe4, 0xae, 0x8f, 0xdc, 0x02, 0xe8, 0x0e, 0x65, 0x51, 0xd6, 0xc8, 0x15, 0xed, 0x5b, 0xab, 0xd0,
	0x90, 0x17, 0x03, 0xc8, 0x02, 0xfe, 0x49, 0xdc, 0x6b, 0xec, 0x54, 0xb1, 0x11, 0x72, 0x47, 0x07,
	0xd6, 0x76, 0x76, 0x61, 0x9b, 0xbb, 0x96, 0x37, 0x21, 0x77, 0xce, 0x7c, 0x2c, 0xcf, 0x84, 0xdc,
	0x40, 0xb8, 0x6e, 0x18, 0x72, 0x55, 0xb0, 0x96, 0x0e, 0xed, 0x86, 0xc0, 0x81, 0xf7, 0x96, 0xcc,
	0x60, 0xdc, 0xa9, 0x52, 0xe9, 0xae, 0x58, 0x8b, 0x0a, 0xe9, 0x28, 0xf6, 0x92, 0x49, 0x0e, 0x0e,
	0x2d, 0x45, 0x85, 0x24, 0x82, 0xe0, 0x14, 0x0c, 0x0d, 0xdc, 0xf0, 0x49, 0x93, 0xff, 0x30, 0x40,
	0x29, 0x85, 0xa4, 0x61, 0xec, 0x27, 0x61, 0xee, 0x44, 0x34, 0x07, 0x38, 0xdf, 0x49, 0xa6, 0xe0,
	0x7f, 0xe1, 0xe1, 0x18, 0x9d, 0x79, 0x9a, 0xa9, 0xef, 0xb2, 0xd1, 0x78, 0x0c, 0xce, 0x89, 0xa7,
	0xde, 0xdc, 0x8b, 0x9e, 0x61, 0x7a, 0x79, 0xca, 0x5f, 0xe6, 0x17, 0xe1, 0xe7, 0xa8, 0x66, 0xdc,
	0x64, 0xb3, 0x1a, 0xda, 0x4f, 0x7c, 0xfc, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x41, 0xff, 0x60, 0x6b,
	0xff, 0x01, 0x00, 0x00,
}
