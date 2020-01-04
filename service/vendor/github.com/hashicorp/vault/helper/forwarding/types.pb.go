// Code generated by protoc-gen-go.
// source: types.proto
// DO NOT EDIT!

/*
Package forwarding is a generated protocol buffer package.

It is generated from these files:
	types.proto

It has these top-level messages:
	Request
	URL
	HeaderEntry
	Response
*/
package forwarding

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

type Request struct {
	// Not used right now but reserving in case it turns out that streaming
	// makes things more economical on the gRPC side
	// uint64 id = 1;
	Method           string                  `protobuf:"bytes,2,opt,name=method" json:"method,omitempty"`
	Url              *URL                    `protobuf:"bytes,3,opt,name=url" json:"url,omitempty"`
	HeaderEntries    map[string]*HeaderEntry `protobuf:"bytes,4,rep,name=header_entries,json=headerEntries" json:"header_entries,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Body             []byte                  `protobuf:"bytes,5,opt,name=body,proto3" json:"body,omitempty"`
	Host             string                  `protobuf:"bytes,6,opt,name=host" json:"host,omitempty"`
	RemoteAddr       string                  `protobuf:"bytes,7,opt,name=remote_addr,json=remoteAddr" json:"remote_addr,omitempty"`
	PeerCertificates [][]byte                `protobuf:"bytes,8,rep,name=peer_certificates,json=peerCertificates,proto3" json:"peer_certificates,omitempty"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Request) GetUrl() *URL {
	if m != nil {
		return m.Url
	}
	return nil
}

func (m *Request) GetHeaderEntries() map[string]*HeaderEntry {
	if m != nil {
		return m.HeaderEntries
	}
	return nil
}

type URL struct {
	Scheme string `protobuf:"bytes,1,opt,name=scheme" json:"scheme,omitempty"`
	Opaque string `protobuf:"bytes,2,opt,name=opaque" json:"opaque,omitempty"`
	// This isn't needed now but might be in the future, so we'll skip the
	// number to keep the ordering in net/url
	// UserInfo user = 3;
	Host    string `protobuf:"bytes,4,opt,name=host" json:"host,omitempty"`
	Path    string `protobuf:"bytes,5,opt,name=path" json:"path,omitempty"`
	RawPath string `protobuf:"bytes,6,opt,name=raw_path,json=rawPath" json:"raw_path,omitempty"`
	// This also isn't needed right now, but we'll reserve the number
	// bool force_query = 7;
	RawQuery string `protobuf:"bytes,8,opt,name=raw_query,json=rawQuery" json:"raw_query,omitempty"`
	Fragment string `protobuf:"bytes,9,opt,name=fragment" json:"fragment,omitempty"`
}

func (m *URL) Reset()                    { *m = URL{} }
func (m *URL) String() string            { return proto.CompactTextString(m) }
func (*URL) ProtoMessage()               {}
func (*URL) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type HeaderEntry struct {
	Values []string `protobuf:"bytes,1,rep,name=values" json:"values,omitempty"`
}

func (m *HeaderEntry) Reset()                    { *m = HeaderEntry{} }
func (m *HeaderEntry) String() string            { return proto.CompactTextString(m) }
func (*HeaderEntry) ProtoMessage()               {}
func (*HeaderEntry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type Response struct {
	// Not used right now but reserving in case it turns out that streaming
	// makes things more economical on the gRPC side
	// uint64 id = 1;
	StatusCode uint32 `protobuf:"varint,2,opt,name=status_code,json=statusCode" json:"status_code,omitempty"`
	Body       []byte `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
	// Added in 0.6.2 to ensure that the content-type is set appropriately, as
	// well as any other information
	HeaderEntries map[string]*HeaderEntry `protobuf:"bytes,4,rep,name=header_entries,json=headerEntries" json:"header_entries,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Response) GetHeaderEntries() map[string]*HeaderEntry {
	if m != nil {
		return m.HeaderEntries
	}
	return nil
}

func init() {
	proto.RegisterType((*Request)(nil), "forwarding.Request")
	proto.RegisterType((*URL)(nil), "forwarding.URL")
	proto.RegisterType((*HeaderEntry)(nil), "forwarding.HeaderEntry")
	proto.RegisterType((*Response)(nil), "forwarding.Response")
}

func init() { proto.RegisterFile("types.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 437 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xb4, 0x53, 0xc1, 0x6e, 0xd4, 0x30,
	0x10, 0x95, 0xd7, 0xdb, 0xdd, 0x64, 0xd2, 0x42, 0xf1, 0x01, 0x4c, 0x39, 0x10, 0x56, 0x02, 0x22,
	0x21, 0xf6, 0xb0, 0x5c, 0x10, 0x37, 0x54, 0x21, 0x71, 0x28, 0x08, 0x2c, 0xf5, 0xc0, 0x29, 0x72,
	0xd7, 0xb3, 0xcd, 0x8a, 0x26, 0xce, 0xda, 0x0e, 0xab, 0x7c, 0x16, 0xff, 0xc4, 0x89, 0xaf, 0x40,
	0xb6, 0x43, 0x1b, 0x84, 0x10, 0xa7, 0x9e, 0x76, 0xde, 0x7b, 0xb3, 0xe3, 0x79, 0x33, 0x13, 0xc8,
	0x5c, 0xdf, 0xa2, 0x5d, 0xb6, 0x46, 0x3b, 0xcd, 0x60, 0xa3, 0xcd, 0x5e, 0x1a, 0xb5, 0x6d, 0x2e,
	0x17, 0x3f, 0x26, 0x30, 0x17, 0xb8, 0xeb, 0xd0, 0x3a, 0x76, 0x1f, 0x66, 0x35, 0xba, 0x4a, 0x2b,
	0x3e, 0xc9, 0x49, 0x91, 0x8a, 0x01, 0xb1, 0x27, 0x40, 0x3b, 0x73, 0xc5, 0x69, 0x4e, 0x8a, 0x6c,
	0x75, 0x77, 0x79, 0xf3, 0xef, 0xe5, 0xb9, 0x38, 0x13, 0x5e, 0x63, 0x1f, 0xe0, 0x4e, 0x85, 0x52,
	0xa1, 0x29, 0xb1, 0x71, 0x66, 0x8b, 0x96, 0x4f, 0x73, 0x5a, 0x64, 0xab, 0x67, 0xe3, 0xec, 0xe1,
	0x9d, 0xe5, 0xfb, 0x90, 0xf9, 0x2e, 0x26, 0xfa, 0x9f, 0x5e, 0x1c, 0x55, 0x63, 0x8e, 0x31, 0x98,
	0x5e, 0x68, 0xd5, 0xf3, 0x83, 0x9c, 0x14, 0x87, 0x22, 0xc4, 0x9e, 0xab, 0xb4, 0x75, 0x7c, 0x16,
	0x7a, 0x0b, 0x31, 0x7b, 0x0c, 0x99, 0xc1, 0x5a, 0x3b, 0x2c, 0xa5, 0x52, 0x86, 0xcf, 0x83, 0x04,
	0x91, 0x7a, 0xab, 0x94, 0x61, 0x2f, 0xe0, 0x5e, 0x8b, 0x68, 0xca, 0x35, 0x1a, 0xb7, 0xdd, 0x6c,
	0xd7, 0xd2, 0xa1, 0xe5, 0x49, 0x4e, 0x8b, 0x43, 0x71, 0xec, 0x85, 0xd3, 0x11, 0x7f, 0xf2, 0x05,
	0xd8, 0xdf, 0xad, 0xb1, 0x63, 0xa0, 0x5f, 0xb1, 0xe7, 0x24, 0xd4, 0xf6, 0x21, 0x7b, 0x09, 0x07,
	0xdf, 0xe4, 0x55, 0x87, 0x61, 0x4c, 0xd9, 0xea, 0xc1, 0xd8, 0xe3, 0x4d, 0x81, 0x5e, 0xc4, 0xac,
	0x37, 0x93, 0xd7, 0x64, 0xf1, 0x9d, 0x00, 0x3d, 0x17, 0x67, 0x7e, 0xc4, 0x76, 0x5d, 0x61, 0x8d,
	0x43, 0xbd, 0x01, 0x79, 0x5e, 0xb7, 0x72, 0x37, 0xd4, 0x4c, 0xc5, 0x80, 0xae, 0x4d, 0x4f, 0x47,
	0xa6, 0x19, 0x4c, 0x5b, 0xe9, 0xaa, 0x30, 0x9c, 0x54, 0x84, 0x98, 0x3d, 0x84, 0xc4, 0xc8, 0x7d,
	0x19, 0xf8, 0x38, 0xa0, 0xb9, 0x91, 0xfb, 0x4f, 0x5e, 0x7a, 0x04, 0xa9, 0x97, 0x76, 0x1d, 0x9a,
	0x9e, 0x27, 0x41, 0xf3, 0xb9, 0x9f, 0x3d, 0x66, 0x27, 0x90, 0x6c, 0x8c, 0xbc, 0xac, 0xb1, 0x71,
	0x3c, 0x8d, 0xda, 0x6f, 0xbc, 0x78, 0x0a, 0xd9, 0xc8, 0x8d, 0x6f, 0x31, 0xf8, 0xb1, 0x9c, 0xe4,
	0xd4, 0xb7, 0x18, 0xd1, 0xe2, 0x27, 0x81, 0x44, 0xa0, 0x6d, 0x75, 0x63, 0xd1, 0x2f, 0xc4, 0x3a,
	0xe9, 0x3a, 0x5b, 0xae, 0xb5, 0x8a, 0x66, 0x8e, 0x04, 0x44, 0xea, 0x54, 0x2b, 0xbc, 0xde, 0x2c,
	0x1d, 0x6d, 0xf6, 0xe3, 0x3f, 0x8e, 0xe7, 0xf9, 0x9f, 0xc7, 0x13, 0x9f, 0xf8, 0xff, 0xf5, 0xdc,
	0xe2, 0x1e, 0x2f, 0x66, 0xe1, 0x0b, 0x7a, 0xf5, 0x2b, 0x00, 0x00, 0xff, 0xff, 0x57, 0x73, 0xdf,
	0x6b, 0x50, 0x03, 0x00, 0x00,
}
