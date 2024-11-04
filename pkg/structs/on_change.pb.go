// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.19.6
// source: on_change.proto

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

type Change struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	World string  `protobuf:"bytes,1,opt,name=World,proto3" json:"World,omitempty"`
	Area  string  `protobuf:"bytes,2,opt,name=Area,proto3" json:"Area,omitempty"`
	Key   Metakey `protobuf:"varint,3,opt,name=Key,proto3,enum=Metakey" json:"Key,omitempty"`
	Id    string  `protobuf:"bytes,4,opt,name=Id,proto3" bson:"_id" json:"Id,omitempty"`
}

func (x *Change) Reset() {
	*x = Change{}
	mi := &file_on_change_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Change) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Change) ProtoMessage() {}

func (x *Change) ProtoReflect() protoreflect.Message {
	mi := &file_on_change_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Change.ProtoReflect.Descriptor instead.
func (*Change) Descriptor() ([]byte, []int) {
	return file_on_change_proto_rawDescGZIP(), []int{0}
}

func (x *Change) GetWorld() string {
	if x != nil {
		return x.World
	}
	return ""
}

func (x *Change) GetArea() string {
	if x != nil {
		return x.Area
	}
	return ""
}

func (x *Change) GetKey() Metakey {
	if x != nil {
		return x.Key
	}
	return Metakey_KeyNone
}

func (x *Change) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type OnChangeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data  *Change `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
	Queue string  `protobuf:"bytes,2,opt,name=Queue,proto3" json:"Queue,omitempty"`
}

func (x *OnChangeRequest) Reset() {
	*x = OnChangeRequest{}
	mi := &file_on_change_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OnChangeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnChangeRequest) ProtoMessage() {}

func (x *OnChangeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_on_change_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnChangeRequest.ProtoReflect.Descriptor instead.
func (*OnChangeRequest) Descriptor() ([]byte, []int) {
	return file_on_change_proto_rawDescGZIP(), []int{1}
}

func (x *OnChangeRequest) GetData() *Change {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *OnChangeRequest) GetQueue() string {
	if x != nil {
		return x.Queue
	}
	return ""
}

type OnChangeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data  *Change `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
	Ack   string  `protobuf:"bytes,2,opt,name=Ack,proto3" json:"Ack,omitempty"`
	Error *Error  `protobuf:"bytes,3,opt,name=Error,proto3,oneof" json:"Error,omitempty"`
}

func (x *OnChangeResponse) Reset() {
	*x = OnChangeResponse{}
	mi := &file_on_change_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OnChangeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OnChangeResponse) ProtoMessage() {}

func (x *OnChangeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_on_change_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OnChangeResponse.ProtoReflect.Descriptor instead.
func (*OnChangeResponse) Descriptor() ([]byte, []int) {
	return file_on_change_proto_rawDescGZIP(), []int{2}
}

func (x *OnChangeResponse) GetData() *Change {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *OnChangeResponse) GetAck() string {
	if x != nil {
		return x.Ack
	}
	return ""
}

func (x *OnChangeResponse) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

type DeferChangeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data *Change `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"` // world, key and Id are required
	// one of these is required
	ToTick *int64 `protobuf:"varint,2,opt,name=ToTick,proto3,oneof" json:"ToTick,omitempty"` // The tick to defer the change to
	ByTick *int64 `protobuf:"varint,3,opt,name=ByTick,proto3,oneof" json:"ByTick,omitempty"` // The tick to defer the change by, short hand for "ToTick=World.Tick+ByTick"
}

func (x *DeferChangeRequest) Reset() {
	*x = DeferChangeRequest{}
	mi := &file_on_change_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeferChangeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeferChangeRequest) ProtoMessage() {}

func (x *DeferChangeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_on_change_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeferChangeRequest.ProtoReflect.Descriptor instead.
func (*DeferChangeRequest) Descriptor() ([]byte, []int) {
	return file_on_change_proto_rawDescGZIP(), []int{3}
}

func (x *DeferChangeRequest) GetData() *Change {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *DeferChangeRequest) GetToTick() int64 {
	if x != nil && x.ToTick != nil {
		return *x.ToTick
	}
	return 0
}

func (x *DeferChangeRequest) GetByTick() int64 {
	if x != nil && x.ByTick != nil {
		return *x.ByTick
	}
	return 0
}

type DeferChangeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error *Error `protobuf:"bytes,2,opt,name=Error,proto3,oneof" json:"Error,omitempty"`
}

func (x *DeferChangeResponse) Reset() {
	*x = DeferChangeResponse{}
	mi := &file_on_change_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeferChangeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeferChangeResponse) ProtoMessage() {}

func (x *DeferChangeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_on_change_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeferChangeResponse.ProtoReflect.Descriptor instead.
func (*DeferChangeResponse) Descriptor() ([]byte, []int) {
	return file_on_change_proto_rawDescGZIP(), []int{4}
}

func (x *DeferChangeResponse) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

type AckRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ack []string `protobuf:"bytes,1,rep,name=Ack,proto3" json:"Ack,omitempty"`
}

func (x *AckRequest) Reset() {
	*x = AckRequest{}
	mi := &file_on_change_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AckRequest) ProtoMessage() {}

func (x *AckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_on_change_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AckRequest.ProtoReflect.Descriptor instead.
func (*AckRequest) Descriptor() ([]byte, []int) {
	return file_on_change_proto_rawDescGZIP(), []int{5}
}

func (x *AckRequest) GetAck() []string {
	if x != nil {
		return x.Ack
	}
	return nil
}

type AckResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error *Error `protobuf:"bytes,1,opt,name=Error,proto3,oneof" json:"Error,omitempty"`
}

func (x *AckResponse) Reset() {
	*x = AckResponse{}
	mi := &file_on_change_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AckResponse) ProtoMessage() {}

func (x *AckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_on_change_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AckResponse.ProtoReflect.Descriptor instead.
func (*AckResponse) Descriptor() ([]byte, []int) {
	return file_on_change_proto_rawDescGZIP(), []int{6}
}

func (x *AckResponse) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

var File_on_change_proto protoreflect.FileDescriptor

var file_on_change_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x6f, 0x6e, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x0e, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x5e, 0x0a, 0x06, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x57, 0x6f, 0x72,
	0x6c, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x41, 0x72, 0x65, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x41,
	0x72, 0x65, 0x61, 0x12, 0x1a, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x08, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x6b, 0x65, 0x79, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x12,
	0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x64, 0x22,
	0x44, 0x0a, 0x0f, 0x4f, 0x6e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1b, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x07, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x14, 0x0a, 0x05, 0x51, 0x75, 0x65, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x51, 0x75, 0x65, 0x75, 0x65, 0x22, 0x6e, 0x0a, 0x10, 0x4f, 0x6e, 0x43, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x04, 0x44, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x10, 0x0a, 0x03, 0x41, 0x63, 0x6b, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x41, 0x63, 0x6b, 0x12, 0x21, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48,
	0x00, 0x52, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x81, 0x01, 0x0a, 0x12, 0x44, 0x65, 0x66, 0x65, 0x72, 0x43,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x04,
	0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x43, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x0a, 0x06, 0x54, 0x6f, 0x54,
	0x69, 0x63, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x06, 0x54, 0x6f, 0x54,
	0x69, 0x63, 0x6b, 0x88, 0x01, 0x01, 0x12, 0x1b, 0x0a, 0x06, 0x42, 0x79, 0x54, 0x69, 0x63, 0x6b,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x48, 0x01, 0x52, 0x06, 0x42, 0x79, 0x54, 0x69, 0x63, 0x6b,
	0x88, 0x01, 0x01, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x54, 0x6f, 0x54, 0x69, 0x63, 0x6b, 0x42, 0x09,
	0x0a, 0x07, 0x5f, 0x42, 0x79, 0x54, 0x69, 0x63, 0x6b, 0x22, 0x42, 0x0a, 0x13, 0x44, 0x65, 0x66,
	0x65, 0x72, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x21, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x06, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x00, 0x52, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x1e, 0x0a,
	0x0a, 0x41, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x41,
	0x63, 0x6b, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x41, 0x63, 0x6b, 0x22, 0x3a, 0x0a,
	0x0b, 0x41, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x05,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x48, 0x00, 0x52, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x88, 0x01, 0x01, 0x42,
	0x08, 0x0a, 0x06, 0x5f, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x42, 0x0d, 0x5a, 0x0b, 0x70, 0x6b, 0x67,
	0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_on_change_proto_rawDescOnce sync.Once
	file_on_change_proto_rawDescData = file_on_change_proto_rawDesc
)

func file_on_change_proto_rawDescGZIP() []byte {
	file_on_change_proto_rawDescOnce.Do(func() {
		file_on_change_proto_rawDescData = protoimpl.X.CompressGZIP(file_on_change_proto_rawDescData)
	})
	return file_on_change_proto_rawDescData
}

var file_on_change_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_on_change_proto_goTypes = []any{
	(*Change)(nil),              // 0: Change
	(*OnChangeRequest)(nil),     // 1: OnChangeRequest
	(*OnChangeResponse)(nil),    // 2: OnChangeResponse
	(*DeferChangeRequest)(nil),  // 3: DeferChangeRequest
	(*DeferChangeResponse)(nil), // 4: DeferChangeResponse
	(*AckRequest)(nil),          // 5: AckRequest
	(*AckResponse)(nil),         // 6: AckResponse
	(Metakey)(0),                // 7: Metakey
	(*Error)(nil),               // 8: Error
}
var file_on_change_proto_depIdxs = []int32{
	7, // 0: Change.Key:type_name -> Metakey
	0, // 1: OnChangeRequest.Data:type_name -> Change
	0, // 2: OnChangeResponse.Data:type_name -> Change
	8, // 3: OnChangeResponse.Error:type_name -> Error
	0, // 4: DeferChangeRequest.Data:type_name -> Change
	8, // 5: DeferChangeResponse.Error:type_name -> Error
	8, // 6: AckResponse.Error:type_name -> Error
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_on_change_proto_init() }
func file_on_change_proto_init() {
	if File_on_change_proto != nil {
		return
	}
	file_metadata_proto_init()
	file_errors_proto_init()
	file_on_change_proto_msgTypes[2].OneofWrappers = []any{}
	file_on_change_proto_msgTypes[3].OneofWrappers = []any{}
	file_on_change_proto_msgTypes[4].OneofWrappers = []any{}
	file_on_change_proto_msgTypes[6].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_on_change_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_on_change_proto_goTypes,
		DependencyIndexes: file_on_change_proto_depIdxs,
		MessageInfos:      file_on_change_proto_msgTypes,
	}.Build()
	File_on_change_proto = out.File
	file_on_change_proto_rawDesc = nil
	file_on_change_proto_goTypes = nil
	file_on_change_proto_depIdxs = nil
}
