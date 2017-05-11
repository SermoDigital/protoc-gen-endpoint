// Code generated by protoc-gen-go.
// source: annotations.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	annotations.proto
	endpoint.proto

It has these top-level messages:
	API
	Endpoint
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

var E_ActionPrefix = &proto1.ExtensionDesc{
	ExtendedType:  (*google_protobuf.FileOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         2131600,
	Name:          "proto.action_prefix",
	Tag:           "bytes,2131600,opt,name=action_prefix,json=actionPrefix",
	Filename:      "annotations.proto",
}

var E_Endpoint = &proto1.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MethodOptions)(nil),
	ExtensionType: (*Endpoint)(nil),
	Field:         2131610,
	Name:          "proto.endpoint",
	Tag:           "bytes,2131610,opt,name=endpoint",
	Filename:      "annotations.proto",
}

func init() {
	proto1.RegisterExtension(E_ActionPrefix)
	proto1.RegisterExtension(E_Endpoint)
}

func init() { proto1.RegisterFile("annotations.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 172 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4c, 0xcc, 0xcb, 0xcb,
	0x2f, 0x49, 0x2c, 0xc9, 0xcc, 0xcf, 0x2b, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05,
	0x53, 0x52, 0x7c, 0xa9, 0x79, 0x29, 0x05, 0xf9, 0x99, 0x79, 0x25, 0x10, 0x61, 0x29, 0x85, 0xf4,
	0xfc, 0xfc, 0xf4, 0x9c, 0x54, 0x7d, 0x30, 0x2f, 0xa9, 0x34, 0x4d, 0x3f, 0x25, 0xb5, 0x38, 0xb9,
	0x28, 0xb3, 0xa0, 0x24, 0xbf, 0x08, 0xa2, 0xc2, 0xca, 0x85, 0x8b, 0x37, 0x31, 0x19, 0x64, 0x52,
	0x7c, 0x41, 0x51, 0x6a, 0x5a, 0x66, 0x85, 0x90, 0x8c, 0x1e, 0x44, 0x8f, 0x1e, 0x4c, 0x8f, 0x9e,
	0x5b, 0x66, 0x4e, 0xaa, 0x7f, 0x01, 0xd8, 0x36, 0x89, 0x09, 0xbd, 0x4d, 0x8c, 0x0a, 0x8c, 0x1a,
	0x9c, 0x41, 0x3c, 0x10, 0x5d, 0x01, 0x60, 0x4d, 0x56, 0x7e, 0x5c, 0x1c, 0x30, 0x9b, 0x85, 0xe4,
	0x30, 0x0c, 0xf0, 0x4d, 0x2d, 0xc9, 0xc8, 0x4f, 0x81, 0x19, 0x31, 0x0b, 0x62, 0x04, 0xb7, 0x11,
	0x3f, 0x44, 0x85, 0x9e, 0x2b, 0x54, 0x67, 0x10, 0xdc, 0x8c, 0x24, 0x36, 0xb0, 0x8c, 0x31, 0x20,
	0x00, 0x00, 0xff, 0xff, 0x85, 0x14, 0x99, 0xdf, 0xea, 0x00, 0x00, 0x00,
}
