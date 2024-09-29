// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.28.2
// source: world.proto

package structs

import (
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

type WorldStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *WorldStatus) Reset() {
	*x = WorldStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_world_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorldStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorldStatus) ProtoMessage() {}

func (x *WorldStatus) ProtoReflect() protoreflect.Message {
	mi := &file_world_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorldStatus.ProtoReflect.Descriptor instead.
func (*WorldStatus) Descriptor() ([]byte, []int) {
	return file_world_proto_rawDescGZIP(), []int{0}
}

type World struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string            `protobuf:"bytes,1,opt,name=Id,proto3" bson:"_id" json:"Id,omitempty"`
	Etag       string            `protobuf:"bytes,2,opt,name=Etag,proto3" json:"Etag,omitempty"`
	Name       string            `protobuf:"bytes,3,opt,name=Name,proto3" json:"Name,omitempty"`
	Labels     map[string]string `protobuf:"bytes,4,rep,name=Labels,proto3" json:"Labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Attributes map[string]int64  `protobuf:"bytes,5,rep,name=Attributes,proto3" json:"Attributes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Tick       int64             `protobuf:"varint,6,opt,name=Tick,proto3" json:"Tick,omitempty"`
	Status     *WorldStatus      `protobuf:"bytes,7,opt,name=Status,proto3" json:"Status,omitempty"`
}

func (x *World) Reset() {
	*x = World{}
	if protoimpl.UnsafeEnabled {
		mi := &file_world_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *World) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*World) ProtoMessage() {}

func (x *World) ProtoReflect() protoreflect.Message {
	mi := &file_world_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use World.ProtoReflect.Descriptor instead.
func (*World) Descriptor() ([]byte, []int) {
	return file_world_proto_rawDescGZIP(), []int{1}
}

func (x *World) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *World) GetEtag() string {
	if x != nil {
		return x.Etag
	}
	return ""
}

func (x *World) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *World) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *World) GetAttributes() map[string]int64 {
	if x != nil {
		return x.Attributes
	}
	return nil
}

func (x *World) GetTick() int64 {
	if x != nil {
		return x.Tick
	}
	return 0
}

func (x *World) GetStatus() *WorldStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

type GetWorldsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ids []string `protobuf:"bytes,1,rep,name=Ids,proto3" json:"Ids,omitempty"`
}

func (x *GetWorldsRequest) Reset() {
	*x = GetWorldsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_world_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetWorldsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetWorldsRequest) ProtoMessage() {}

func (x *GetWorldsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_world_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetWorldsRequest.ProtoReflect.Descriptor instead.
func (*GetWorldsRequest) Descriptor() ([]byte, []int) {
	return file_world_proto_rawDescGZIP(), []int{2}
}

func (x *GetWorldsRequest) GetIds() []string {
	if x != nil {
		return x.Ids
	}
	return nil
}

type GetWorldsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data  []*World `protobuf:"bytes,1,rep,name=Data,proto3" json:"Data,omitempty"`
	Error *Error   `protobuf:"bytes,2,opt,name=Error,proto3,oneof" json:"Error,omitempty"`
}

func (x *GetWorldsResponse) Reset() {
	*x = GetWorldsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_world_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetWorldsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetWorldsResponse) ProtoMessage() {}

func (x *GetWorldsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_world_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetWorldsResponse.ProtoReflect.Descriptor instead.
func (*GetWorldsResponse) Descriptor() ([]byte, []int) {
	return file_world_proto_rawDescGZIP(), []int{3}
}

func (x *GetWorldsResponse) GetData() []*World {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *GetWorldsResponse) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

type SetWorldRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *World `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *SetWorldRequest) Reset() {
	*x = SetWorldRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_world_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetWorldRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetWorldRequest) ProtoMessage() {}

func (x *SetWorldRequest) ProtoReflect() protoreflect.Message {
	mi := &file_world_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetWorldRequest.ProtoReflect.Descriptor instead.
func (*SetWorldRequest) Descriptor() ([]byte, []int) {
	return file_world_proto_rawDescGZIP(), []int{4}
}

func (x *SetWorldRequest) GetData() *World {
	if x != nil {
		return x.Data
	}
	return nil
}

type SetWorldResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Etag  string `protobuf:"bytes,1,opt,name=Etag,proto3" json:"Etag,omitempty"`
	Error *Error `protobuf:"bytes,2,opt,name=Error,proto3,oneof" json:"Error,omitempty"`
}

func (x *SetWorldResponse) Reset() {
	*x = SetWorldResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_world_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetWorldResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetWorldResponse) ProtoMessage() {}

func (x *SetWorldResponse) ProtoReflect() protoreflect.Message {
	mi := &file_world_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetWorldResponse.ProtoReflect.Descriptor instead.
func (*SetWorldResponse) Descriptor() ([]byte, []int) {
	return file_world_proto_rawDescGZIP(), []int{5}
}

func (x *SetWorldResponse) GetEtag() string {
	if x != nil {
		return x.Etag
	}
	return ""
}

func (x *SetWorldResponse) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

type DeleteWorldRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=Id,proto3" bson:"_id" json:"Id,omitempty"`
}

func (x *DeleteWorldRequest) Reset() {
	*x = DeleteWorldRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_world_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteWorldRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteWorldRequest) ProtoMessage() {}

func (x *DeleteWorldRequest) ProtoReflect() protoreflect.Message {
	mi := &file_world_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteWorldRequest.ProtoReflect.Descriptor instead.
func (*DeleteWorldRequest) Descriptor() ([]byte, []int) {
	return file_world_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteWorldRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type DeleteWorldResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error *Error `protobuf:"bytes,1,opt,name=Error,proto3,oneof" json:"Error,omitempty"`
}

func (x *DeleteWorldResponse) Reset() {
	*x = DeleteWorldResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_world_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteWorldResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteWorldResponse) ProtoMessage() {}

func (x *DeleteWorldResponse) ProtoReflect() protoreflect.Message {
	mi := &file_world_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteWorldResponse.ProtoReflect.Descriptor instead.
func (*DeleteWorldResponse) Descriptor() ([]byte, []int) {
	return file_world_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteWorldResponse) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

type ListWorldsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Limit  *uint32           `protobuf:"varint,1,opt,name=Limit,proto3,oneof" json:"Limit,omitempty"`
	Offset *uint32           `protobuf:"varint,2,opt,name=Offset,proto3,oneof" json:"Offset,omitempty"`
	Labels map[string]string `protobuf:"bytes,3,rep,name=Labels,proto3" json:"Labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *ListWorldsRequest) Reset() {
	*x = ListWorldsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_world_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListWorldsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListWorldsRequest) ProtoMessage() {}

func (x *ListWorldsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_world_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListWorldsRequest.ProtoReflect.Descriptor instead.
func (*ListWorldsRequest) Descriptor() ([]byte, []int) {
	return file_world_proto_rawDescGZIP(), []int{8}
}

func (x *ListWorldsRequest) GetLimit() uint32 {
	if x != nil && x.Limit != nil {
		return *x.Limit
	}
	return 0
}

func (x *ListWorldsRequest) GetOffset() uint32 {
	if x != nil && x.Offset != nil {
		return *x.Offset
	}
	return 0
}

func (x *ListWorldsRequest) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

type ListWorldsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data  []*World `protobuf:"bytes,1,rep,name=Data,proto3" json:"Data,omitempty"`
	Error *Error   `protobuf:"bytes,2,opt,name=Error,proto3,oneof" json:"Error,omitempty"`
}

func (x *ListWorldsResponse) Reset() {
	*x = ListWorldsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_world_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListWorldsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListWorldsResponse) ProtoMessage() {}

func (x *ListWorldsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_world_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListWorldsResponse.ProtoReflect.Descriptor instead.
func (*ListWorldsResponse) Descriptor() ([]byte, []int) {
	return file_world_proto_rawDescGZIP(), []int{9}
}

func (x *ListWorldsResponse) GetData() []*World {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *ListWorldsResponse) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

var File_world_proto protoreflect.FileDescriptor

var file_world_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0e, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x0d, 0x0a, 0x0b, 0x57,
	0x6f, 0x72, 0x6c, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0xd7, 0x02, 0x0a, 0x05, 0x57,
	0x6f, 0x72, 0x6c, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x45, 0x74, 0x61, 0x67, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x45, 0x74, 0x61, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2a, 0x0a, 0x06,
	0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x57,
	0x6f, 0x72, 0x6c, 0x64, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x06, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x36, 0x0a, 0x0a, 0x41, 0x74, 0x74, 0x72,
	0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x57,
	0x6f, 0x72, 0x6c, 0x64, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73,
	0x12, 0x12, 0x0a, 0x04, 0x54, 0x69, 0x63, 0x6b, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04,
	0x54, 0x69, 0x63, 0x6b, 0x12, 0x24, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61,
	0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3d, 0x0a, 0x0f, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75,
	0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x22, 0x24, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x6c, 0x64,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x49, 0x64, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x49, 0x64, 0x73, 0x22, 0x5c, 0x0a, 0x11, 0x47, 0x65,
	0x74, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x1a, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x06, 0x2e,
	0x57, 0x6f, 0x72, 0x6c, 0x64, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x21, 0x0a, 0x05, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x48, 0x00, 0x52, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x88, 0x01, 0x01, 0x42, 0x08,
	0x0a, 0x06, 0x5f, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x2d, 0x0a, 0x0f, 0x53, 0x65, 0x74, 0x57,
	0x6f, 0x72, 0x6c, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x04, 0x44,
	0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x57, 0x6f, 0x72, 0x6c,
	0x64, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x22, 0x53, 0x0a, 0x10, 0x53, 0x65, 0x74, 0x57, 0x6f,
	0x72, 0x6c, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x45,
	0x74, 0x61, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x45, 0x74, 0x61, 0x67, 0x12,
	0x21, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06,
	0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x00, 0x52, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x88,
	0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x24, 0x0a, 0x12,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x49, 0x64, 0x22, 0x42, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6c,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x05, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x48, 0x00, 0x52, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06,
	0x5f, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0xd3, 0x01, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x57,
	0x6f, 0x72, 0x6c, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x05,
	0x4c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x00, 0x52, 0x05, 0x4c,
	0x69, 0x6d, 0x69, 0x74, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x4f, 0x66, 0x66, 0x73, 0x65,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x01, 0x52, 0x06, 0x4f, 0x66, 0x66, 0x73, 0x65,
	0x74, 0x88, 0x01, 0x01, 0x12, 0x36, 0x0a, 0x06, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x6f, 0x72, 0x6c, 0x64,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x1a, 0x39, 0x0a, 0x0b,
	0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x4c, 0x69, 0x6d, 0x69,
	0x74, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x22, 0x5d, 0x0a, 0x12,
	0x4c, 0x69, 0x73, 0x74, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x1a, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x06, 0x2e, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x21,
	0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x00, 0x52, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x88, 0x01,
	0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x42, 0x0d, 0x5a, 0x0b, 0x70,
	0x6b, 0x67, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_world_proto_rawDescOnce sync.Once
	file_world_proto_rawDescData = file_world_proto_rawDesc
)

func file_world_proto_rawDescGZIP() []byte {
	file_world_proto_rawDescOnce.Do(func() {
		file_world_proto_rawDescData = protoimpl.X.CompressGZIP(file_world_proto_rawDescData)
	})
	return file_world_proto_rawDescData
}

var file_world_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_world_proto_goTypes = []interface{}{
	(*WorldStatus)(nil),         // 0: WorldStatus
	(*World)(nil),               // 1: World
	(*GetWorldsRequest)(nil),    // 2: GetWorldsRequest
	(*GetWorldsResponse)(nil),   // 3: GetWorldsResponse
	(*SetWorldRequest)(nil),     // 4: SetWorldRequest
	(*SetWorldResponse)(nil),    // 5: SetWorldResponse
	(*DeleteWorldRequest)(nil),  // 6: DeleteWorldRequest
	(*DeleteWorldResponse)(nil), // 7: DeleteWorldResponse
	(*ListWorldsRequest)(nil),   // 8: ListWorldsRequest
	(*ListWorldsResponse)(nil),  // 9: ListWorldsResponse
	nil,                         // 10: World.LabelsEntry
	nil,                         // 11: World.AttributesEntry
	nil,                         // 12: ListWorldsRequest.LabelsEntry
	(*Error)(nil),               // 13: Error
}
var file_world_proto_depIdxs = []int32{
	10, // 0: World.Labels:type_name -> World.LabelsEntry
	11, // 1: World.Attributes:type_name -> World.AttributesEntry
	0,  // 2: World.Status:type_name -> WorldStatus
	1,  // 3: GetWorldsResponse.Data:type_name -> World
	13, // 4: GetWorldsResponse.Error:type_name -> Error
	1,  // 5: SetWorldRequest.Data:type_name -> World
	13, // 6: SetWorldResponse.Error:type_name -> Error
	13, // 7: DeleteWorldResponse.Error:type_name -> Error
	12, // 8: ListWorldsRequest.Labels:type_name -> ListWorldsRequest.LabelsEntry
	1,  // 9: ListWorldsResponse.Data:type_name -> World
	13, // 10: ListWorldsResponse.Error:type_name -> Error
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_world_proto_init() }
func file_world_proto_init() {
	if File_world_proto != nil {
		return
	}
	file_metadata_proto_init()
	file_errors_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_world_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorldStatus); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_world_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*World); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_world_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetWorldsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_world_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetWorldsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_world_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetWorldRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_world_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetWorldResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_world_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteWorldRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_world_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteWorldResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_world_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListWorldsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_world_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListWorldsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_world_proto_msgTypes[3].OneofWrappers = []interface{}{}
	file_world_proto_msgTypes[5].OneofWrappers = []interface{}{}
	file_world_proto_msgTypes[7].OneofWrappers = []interface{}{}
	file_world_proto_msgTypes[8].OneofWrappers = []interface{}{}
	file_world_proto_msgTypes[9].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_world_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_world_proto_goTypes,
		DependencyIndexes: file_world_proto_depIdxs,
		MessageInfos:      file_world_proto_msgTypes,
	}.Build()
	File_world_proto = out.File
	file_world_proto_rawDesc = nil
	file_world_proto_goTypes = nil
	file_world_proto_depIdxs = nil
}
