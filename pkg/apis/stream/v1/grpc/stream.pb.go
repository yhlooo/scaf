// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.29.1
// source: pkg/apis/stream/v1/grpc/stream.proto

package grpc

import (
	grpc "github.com/yhlooo/scaf/pkg/apis/meta/v1/grpc"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// GetStreamRequest GetStream 请求
type GetStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GetStreamRequest) Reset() {
	*x = GetStreamRequest{}
	mi := &file_pkg_apis_stream_v1_grpc_stream_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStreamRequest) ProtoMessage() {}

func (x *GetStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_stream_v1_grpc_stream_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStreamRequest.ProtoReflect.Descriptor instead.
func (*GetStreamRequest) Descriptor() ([]byte, []int) {
	return file_pkg_apis_stream_v1_grpc_stream_proto_rawDescGZIP(), []int{0}
}

func (x *GetStreamRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// ListStreamsRequest ListStreams 请求
type ListStreamsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListStreamsRequest) Reset() {
	*x = ListStreamsRequest{}
	mi := &file_pkg_apis_stream_v1_grpc_stream_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListStreamsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListStreamsRequest) ProtoMessage() {}

func (x *ListStreamsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_stream_v1_grpc_stream_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListStreamsRequest.ProtoReflect.Descriptor instead.
func (*ListStreamsRequest) Descriptor() ([]byte, []int) {
	return file_pkg_apis_stream_v1_grpc_stream_proto_rawDescGZIP(), []int{1}
}

// DeleteStreamRequest DeleteStream 请求
type DeleteStreamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *DeleteStreamRequest) Reset() {
	*x = DeleteStreamRequest{}
	mi := &file_pkg_apis_stream_v1_grpc_stream_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteStreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteStreamRequest) ProtoMessage() {}

func (x *DeleteStreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_stream_v1_grpc_stream_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteStreamRequest.ProtoReflect.Descriptor instead.
func (*DeleteStreamRequest) Descriptor() ([]byte, []int) {
	return file_pkg_apis_stream_v1_grpc_stream_proto_rawDescGZIP(), []int{2}
}

func (x *DeleteStreamRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Package 流中传递的包
type Package struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content []byte `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *Package) Reset() {
	*x = Package{}
	mi := &file_pkg_apis_stream_v1_grpc_stream_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Package) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Package) ProtoMessage() {}

func (x *Package) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_stream_v1_grpc_stream_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Package.ProtoReflect.Descriptor instead.
func (*Package) Descriptor() ([]byte, []int) {
	return file_pkg_apis_stream_v1_grpc_stream_proto_rawDescGZIP(), []int{3}
}

func (x *Package) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

// Stream 流
type Stream struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata *grpc.ObjectMeta `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Spec     *StreamSpec      `protobuf:"bytes,2,opt,name=spec,proto3" json:"spec,omitempty"`
	Status   *StreamStatus    `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *Stream) Reset() {
	*x = Stream{}
	mi := &file_pkg_apis_stream_v1_grpc_stream_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Stream) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stream) ProtoMessage() {}

func (x *Stream) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_stream_v1_grpc_stream_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stream.ProtoReflect.Descriptor instead.
func (*Stream) Descriptor() ([]byte, []int) {
	return file_pkg_apis_stream_v1_grpc_stream_proto_rawDescGZIP(), []int{4}
}

func (x *Stream) GetMetadata() *grpc.ObjectMeta {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Stream) GetSpec() *StreamSpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

func (x *Stream) GetStatus() *StreamStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

// StreamSpec 流定义
type StreamSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 停止策略
	StopPolicy string `protobuf:"bytes,1,opt,name=stop_policy,json=stopPolicy,proto3" json:"stop_policy,omitempty"`
}

func (x *StreamSpec) Reset() {
	*x = StreamSpec{}
	mi := &file_pkg_apis_stream_v1_grpc_stream_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamSpec) ProtoMessage() {}

func (x *StreamSpec) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_stream_v1_grpc_stream_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamSpec.ProtoReflect.Descriptor instead.
func (*StreamSpec) Descriptor() ([]byte, []int) {
	return file_pkg_apis_stream_v1_grpc_stream_proto_rawDescGZIP(), []int{5}
}

func (x *StreamSpec) GetStopPolicy() string {
	if x != nil {
		return x.StopPolicy
	}
	return ""
}

// StreamStatus 流状态
type StreamStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 用于加入流的 token
	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *StreamStatus) Reset() {
	*x = StreamStatus{}
	mi := &file_pkg_apis_stream_v1_grpc_stream_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamStatus) ProtoMessage() {}

func (x *StreamStatus) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_stream_v1_grpc_stream_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamStatus.ProtoReflect.Descriptor instead.
func (*StreamStatus) Descriptor() ([]byte, []int) {
	return file_pkg_apis_stream_v1_grpc_stream_proto_rawDescGZIP(), []int{6}
}

func (x *StreamStatus) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

// StreamList 流列表
type StreamList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata *grpc.ListMeta `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Items    []*Stream      `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *StreamList) Reset() {
	*x = StreamList{}
	mi := &file_pkg_apis_stream_v1_grpc_stream_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamList) ProtoMessage() {}

func (x *StreamList) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_apis_stream_v1_grpc_stream_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamList.ProtoReflect.Descriptor instead.
func (*StreamList) Descriptor() ([]byte, []int) {
	return file_pkg_apis_stream_v1_grpc_stream_proto_rawDescGZIP(), []int{7}
}

func (x *StreamList) GetMetadata() *grpc.ListMeta {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *StreamList) GetItems() []*Stream {
	if x != nil {
		return x.Items
	}
	return nil
}

var File_pkg_apis_stream_v1_grpc_stream_proto protoreflect.FileDescriptor

var file_pkg_apis_stream_v1_grpc_stream_proto_rawDesc = []byte{
	0x0a, 0x24, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x73, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x79, 0x68, 0x6c, 0x6f, 0x6f, 0x6f, 0x2e, 0x63,
	0x6f, 0x6d, 0x2e, 0x73, 0x63, 0x61, 0x66, 0x2e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x76,
	0x31, 0x1a, 0x20, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x6d, 0x65, 0x74, 0x61,
	0x2f, 0x76, 0x31, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x6d, 0x65, 0x74, 0x61, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x22, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x6d, 0x65,
	0x74, 0x61, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x26, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22,
	0x14, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x29, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x22, 0x23, 0x0a, 0x07, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0xc5, 0x01, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x12, 0x3f, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x23, 0x2e, 0x79, 0x68, 0x6c, 0x6f, 0x6f, 0x6f, 0x2e, 0x63, 0x6f, 0x6d, 0x2e,
	0x73, 0x63, 0x61, 0x66, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x12, 0x39, 0x0a, 0x04, 0x73, 0x70, 0x65, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x25, 0x2e, 0x79, 0x68, 0x6c, 0x6f, 0x6f, 0x6f, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x63, 0x61,
	0x66, 0x2e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x53, 0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70, 0x65, 0x63, 0x12, 0x3f, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x79,
	0x68, 0x6c, 0x6f, 0x6f, 0x6f, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x63, 0x61, 0x66, 0x2e, 0x73,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x2d, 0x0a,
	0x0a, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x70, 0x65, 0x63, 0x12, 0x1f, 0x0a, 0x0b, 0x73,
	0x74, 0x6f, 0x70, 0x5f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x73, 0x74, 0x6f, 0x70, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x22, 0x24, 0x0a, 0x0c,
	0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x22, 0x84, 0x01, 0x0a, 0x0a, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4c, 0x69, 0x73,
	0x74, 0x12, 0x3d, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x79, 0x68, 0x6c, 0x6f, 0x6f, 0x6f, 0x2e, 0x63, 0x6f, 0x6d,
	0x2e, 0x73, 0x63, 0x61, 0x66, 0x2e, 0x6d, 0x65, 0x74, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x12, 0x37, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x21, 0x2e, 0x79, 0x68, 0x6c, 0x6f, 0x6f, 0x6f, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x63, 0x61,
	0x66, 0x2e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x32, 0xdf, 0x03, 0x0a, 0x07, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x73, 0x12, 0x54, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x21, 0x2e, 0x79, 0x68, 0x6c, 0x6f, 0x6f, 0x6f, 0x2e, 0x63,
	0x6f, 0x6d, 0x2e, 0x73, 0x63, 0x61, 0x66, 0x2e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x1a, 0x21, 0x2e, 0x79, 0x68, 0x6c, 0x6f, 0x6f,
	0x6f, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x63, 0x61, 0x66, 0x2e, 0x73, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x5b, 0x0a, 0x09, 0x47,
	0x65, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x2b, 0x2e, 0x79, 0x68, 0x6c, 0x6f, 0x6f,
	0x6f, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x63, 0x61, 0x66, 0x2e, 0x73, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x79, 0x68, 0x6c, 0x6f, 0x6f, 0x6f, 0x2e, 0x63,
	0x6f, 0x6d, 0x2e, 0x73, 0x63, 0x61, 0x66, 0x2e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x63, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74,
	0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x73, 0x12, 0x2d, 0x2e, 0x79, 0x68, 0x6c, 0x6f, 0x6f, 0x6f,
	0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x63, 0x61, 0x66, 0x2e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x79, 0x68, 0x6c, 0x6f, 0x6f, 0x6f, 0x2e,
	0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x63, 0x61, 0x66, 0x2e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x5f, 0x0a,
	0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x2e, 0x2e,
	0x79, 0x68, 0x6c, 0x6f, 0x6f, 0x6f, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x63, 0x61, 0x66, 0x2e,
	0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e,
	0x79, 0x68, 0x6c, 0x6f, 0x6f, 0x6f, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x63, 0x61, 0x66, 0x2e,
	0x6d, 0x65, 0x74, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x5b,
	0x0a, 0x0d, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12,
	0x22, 0x2e, 0x79, 0x68, 0x6c, 0x6f, 0x6f, 0x6f, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x73, 0x63, 0x61,
	0x66, 0x2e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x63, 0x6b,
	0x61, 0x67, 0x65, 0x1a, 0x22, 0x2e, 0x79, 0x68, 0x6c, 0x6f, 0x6f, 0x6f, 0x2e, 0x63, 0x6f, 0x6d,
	0x2e, 0x73, 0x63, 0x61, 0x66, 0x2e, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x76, 0x31, 0x2e,
	0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x28, 0x01, 0x30, 0x01, 0x42, 0x30, 0x5a, 0x2e, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x79, 0x68, 0x6c, 0x6f, 0x6f, 0x6f,
	0x2f, 0x73, 0x63, 0x61, 0x66, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x73,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_apis_stream_v1_grpc_stream_proto_rawDescOnce sync.Once
	file_pkg_apis_stream_v1_grpc_stream_proto_rawDescData = file_pkg_apis_stream_v1_grpc_stream_proto_rawDesc
)

func file_pkg_apis_stream_v1_grpc_stream_proto_rawDescGZIP() []byte {
	file_pkg_apis_stream_v1_grpc_stream_proto_rawDescOnce.Do(func() {
		file_pkg_apis_stream_v1_grpc_stream_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_apis_stream_v1_grpc_stream_proto_rawDescData)
	})
	return file_pkg_apis_stream_v1_grpc_stream_proto_rawDescData
}

var file_pkg_apis_stream_v1_grpc_stream_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_pkg_apis_stream_v1_grpc_stream_proto_goTypes = []any{
	(*GetStreamRequest)(nil),    // 0: yhlooo.com.scaf.stream.v1.GetStreamRequest
	(*ListStreamsRequest)(nil),  // 1: yhlooo.com.scaf.stream.v1.ListStreamsRequest
	(*DeleteStreamRequest)(nil), // 2: yhlooo.com.scaf.stream.v1.DeleteStreamRequest
	(*Package)(nil),             // 3: yhlooo.com.scaf.stream.v1.Package
	(*Stream)(nil),              // 4: yhlooo.com.scaf.stream.v1.Stream
	(*StreamSpec)(nil),          // 5: yhlooo.com.scaf.stream.v1.StreamSpec
	(*StreamStatus)(nil),        // 6: yhlooo.com.scaf.stream.v1.StreamStatus
	(*StreamList)(nil),          // 7: yhlooo.com.scaf.stream.v1.StreamList
	(*grpc.ObjectMeta)(nil),     // 8: yhlooo.com.scaf.meta.v1.ObjectMeta
	(*grpc.ListMeta)(nil),       // 9: yhlooo.com.scaf.meta.v1.ListMeta
	(*grpc.Status)(nil),         // 10: yhlooo.com.scaf.meta.v1.Status
}
var file_pkg_apis_stream_v1_grpc_stream_proto_depIdxs = []int32{
	8,  // 0: yhlooo.com.scaf.stream.v1.Stream.metadata:type_name -> yhlooo.com.scaf.meta.v1.ObjectMeta
	5,  // 1: yhlooo.com.scaf.stream.v1.Stream.spec:type_name -> yhlooo.com.scaf.stream.v1.StreamSpec
	6,  // 2: yhlooo.com.scaf.stream.v1.Stream.status:type_name -> yhlooo.com.scaf.stream.v1.StreamStatus
	9,  // 3: yhlooo.com.scaf.stream.v1.StreamList.metadata:type_name -> yhlooo.com.scaf.meta.v1.ListMeta
	4,  // 4: yhlooo.com.scaf.stream.v1.StreamList.items:type_name -> yhlooo.com.scaf.stream.v1.Stream
	4,  // 5: yhlooo.com.scaf.stream.v1.Streams.CreateStream:input_type -> yhlooo.com.scaf.stream.v1.Stream
	0,  // 6: yhlooo.com.scaf.stream.v1.Streams.GetStream:input_type -> yhlooo.com.scaf.stream.v1.GetStreamRequest
	1,  // 7: yhlooo.com.scaf.stream.v1.Streams.ListStreams:input_type -> yhlooo.com.scaf.stream.v1.ListStreamsRequest
	2,  // 8: yhlooo.com.scaf.stream.v1.Streams.DeleteStream:input_type -> yhlooo.com.scaf.stream.v1.DeleteStreamRequest
	3,  // 9: yhlooo.com.scaf.stream.v1.Streams.ConnectStream:input_type -> yhlooo.com.scaf.stream.v1.Package
	4,  // 10: yhlooo.com.scaf.stream.v1.Streams.CreateStream:output_type -> yhlooo.com.scaf.stream.v1.Stream
	4,  // 11: yhlooo.com.scaf.stream.v1.Streams.GetStream:output_type -> yhlooo.com.scaf.stream.v1.Stream
	7,  // 12: yhlooo.com.scaf.stream.v1.Streams.ListStreams:output_type -> yhlooo.com.scaf.stream.v1.StreamList
	10, // 13: yhlooo.com.scaf.stream.v1.Streams.DeleteStream:output_type -> yhlooo.com.scaf.meta.v1.Status
	3,  // 14: yhlooo.com.scaf.stream.v1.Streams.ConnectStream:output_type -> yhlooo.com.scaf.stream.v1.Package
	10, // [10:15] is the sub-list for method output_type
	5,  // [5:10] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_pkg_apis_stream_v1_grpc_stream_proto_init() }
func file_pkg_apis_stream_v1_grpc_stream_proto_init() {
	if File_pkg_apis_stream_v1_grpc_stream_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_apis_stream_v1_grpc_stream_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_apis_stream_v1_grpc_stream_proto_goTypes,
		DependencyIndexes: file_pkg_apis_stream_v1_grpc_stream_proto_depIdxs,
		MessageInfos:      file_pkg_apis_stream_v1_grpc_stream_proto_msgTypes,
	}.Build()
	File_pkg_apis_stream_v1_grpc_stream_proto = out.File
	file_pkg_apis_stream_v1_grpc_stream_proto_rawDesc = nil
	file_pkg_apis_stream_v1_grpc_stream_proto_goTypes = nil
	file_pkg_apis_stream_v1_grpc_stream_proto_depIdxs = nil
}
