// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.3
// source: family.proto

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

type Family struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ethos *Ethos `protobuf:"bytes,1,opt,name=Ethos,proto3" json:"Ethos,omitempty"`
	ID    string `protobuf:"bytes,2,opt,name=ID,proto3" json:"ID,omitempty"`
	// Race demographic assigned to children
	Race string `protobuf:"bytes,3,opt,name=Race,proto3" json:"Race,omitempty"`
	// Culture demographic assigned to children
	Culture string `protobuf:"bytes,4,opt,name=Culture,proto3" json:"Culture,omitempty"`
	// Area where the family is based (where children will be placed)
	AreaID string `protobuf:"bytes,5,opt,name=AreaID,proto3" json:"AreaID,omitempty"`
	// social class of the family
	SocialClass string `protobuf:"bytes,6,opt,name=SocialClass,proto3" json:"SocialClass,omitempty"`
	// Faction ID (if any) if this family is simulated as a major player.
	//
	// This implies the family is fairly wealthy and/or influential, probably 95% of families
	// will not have this set; which is probably a good thing and saves us lots of calculations
	// for families which don't really have the resources to act on the national / international
	// stage.
	FactionID string `protobuf:"bytes,7,opt,name=FactionID,proto3" json:"FactionID,omitempty"`
	// True while;
	// - both people are capable of bearing children
	// - both people are married or lovers (ie. willing to bear children)
	IsChildBearing bool `protobuf:"varint,8,opt,name=IsChildBearing,proto3" json:"IsChildBearing,omitempty"`
	// Represents the tick when one of the two potential parents becomes too old to bear children.
	MaxChildBearingTick int64 `protobuf:"varint,9,opt,name=MaxChildBearingTick,proto3" json:"MaxChildBearingTick,omitempty"`
	// If mother is pregnant then this is the tick when she will give birth.
	// Nb.
	//   - we always set a Family to *not* child bearing if either partner dies
	//     (as in, it cannot produce more children)
	//   - we set PregnancyEnd to 0 when the child is born or the mother dies
	//     (this means the child can be born if the father is dead)
	//
	// This saves us having to query the parents when doing calculations
	PregnancyEnd int64 `protobuf:"varint,10,opt,name=PregnancyEnd,proto3" json:"PregnancyEnd,omitempty"`
	// A family consists of a male & female and can bear children.
	// Nb. this does not imply that the couple are married ..
	MaleID   string `protobuf:"bytes,11,opt,name=MaleID,proto3" json:"MaleID,omitempty"`
	FemaleID string `protobuf:"bytes,12,opt,name=FemaleID,proto3" json:"FemaleID,omitempty"`
	// Save us looking up families (the info doesn't change anyways)
	MaGrandmaID string `protobuf:"bytes,13,opt,name=MaGrandmaID,proto3" json:"MaGrandmaID,omitempty"`
	MaGrandpaID string `protobuf:"bytes,14,opt,name=MaGrandpaID,proto3" json:"MaGrandpaID,omitempty"`
	PaGrandmaID string `protobuf:"bytes,15,opt,name=PaGrandmaID,proto3" json:"PaGrandmaID,omitempty"`
	PaGrandpaID string `protobuf:"bytes,16,opt,name=PaGrandpaID,proto3" json:"PaGrandpaID,omitempty"`
	// Number of children this family has had
	NumberofChildren int64 `protobuf:"varint,17,opt,name=NumberofChildren,proto3" json:"NumberofChildren,omitempty"`
	// The tick when the couple got married, divorced or widowed
	MarriageTick int64 `protobuf:"varint,18,opt,name=MarriageTick,proto3" json:"MarriageTick,omitempty"`
	DivorceTick  int64 `protobuf:"varint,19,opt,name=DivorceTick,proto3" json:"DivorceTick,omitempty"`
	WidowedTick  int64 `protobuf:"varint,20,opt,name=WidowedTick,proto3" json:"WidowedTick,omitempty"`
	// Random number used for blind selection
	Random int64 `protobuf:"varint,21,opt,name=Random,proto3" json:"Random,omitempty"`
}

func (x *Family) Reset() {
	*x = Family{}
	if protoimpl.UnsafeEnabled {
		mi := &file_family_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Family) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Family) ProtoMessage() {}

func (x *Family) ProtoReflect() protoreflect.Message {
	mi := &file_family_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Family.ProtoReflect.Descriptor instead.
func (*Family) Descriptor() ([]byte, []int) {
	return file_family_proto_rawDescGZIP(), []int{0}
}

func (x *Family) GetEthos() *Ethos {
	if x != nil {
		return x.Ethos
	}
	return nil
}

func (x *Family) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Family) GetRace() string {
	if x != nil {
		return x.Race
	}
	return ""
}

func (x *Family) GetCulture() string {
	if x != nil {
		return x.Culture
	}
	return ""
}

func (x *Family) GetAreaID() string {
	if x != nil {
		return x.AreaID
	}
	return ""
}

func (x *Family) GetSocialClass() string {
	if x != nil {
		return x.SocialClass
	}
	return ""
}

func (x *Family) GetFactionID() string {
	if x != nil {
		return x.FactionID
	}
	return ""
}

func (x *Family) GetIsChildBearing() bool {
	if x != nil {
		return x.IsChildBearing
	}
	return false
}

func (x *Family) GetMaxChildBearingTick() int64 {
	if x != nil {
		return x.MaxChildBearingTick
	}
	return 0
}

func (x *Family) GetPregnancyEnd() int64 {
	if x != nil {
		return x.PregnancyEnd
	}
	return 0
}

func (x *Family) GetMaleID() string {
	if x != nil {
		return x.MaleID
	}
	return ""
}

func (x *Family) GetFemaleID() string {
	if x != nil {
		return x.FemaleID
	}
	return ""
}

func (x *Family) GetMaGrandmaID() string {
	if x != nil {
		return x.MaGrandmaID
	}
	return ""
}

func (x *Family) GetMaGrandpaID() string {
	if x != nil {
		return x.MaGrandpaID
	}
	return ""
}

func (x *Family) GetPaGrandmaID() string {
	if x != nil {
		return x.PaGrandmaID
	}
	return ""
}

func (x *Family) GetPaGrandpaID() string {
	if x != nil {
		return x.PaGrandpaID
	}
	return ""
}

func (x *Family) GetNumberofChildren() int64 {
	if x != nil {
		return x.NumberofChildren
	}
	return 0
}

func (x *Family) GetMarriageTick() int64 {
	if x != nil {
		return x.MarriageTick
	}
	return 0
}

func (x *Family) GetDivorceTick() int64 {
	if x != nil {
		return x.DivorceTick
	}
	return 0
}

func (x *Family) GetWidowedTick() int64 {
	if x != nil {
		return x.WidowedTick
	}
	return 0
}

func (x *Family) GetRandom() int64 {
	if x != nil {
		return x.Random
	}
	return 0
}

var File_family_proto protoreflect.FileDescriptor

var file_family_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x66, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b,
	0x65, 0x74, 0x68, 0x6f, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa2, 0x05, 0x0a, 0x06,
	0x46, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x12, 0x1c, 0x0a, 0x05, 0x45, 0x74, 0x68, 0x6f, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x45, 0x74, 0x68, 0x6f, 0x73, 0x52, 0x05, 0x45,
	0x74, 0x68, 0x6f, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x52, 0x61, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x52, 0x61, 0x63, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x75, 0x6c, 0x74,
	0x75, 0x72, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x43, 0x75, 0x6c, 0x74, 0x75,
	0x72, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x41, 0x72, 0x65, 0x61, 0x49, 0x44, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x41, 0x72, 0x65, 0x61, 0x49, 0x44, 0x12, 0x20, 0x0a, 0x0b, 0x53, 0x6f,
	0x63, 0x69, 0x61, 0x6c, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x53, 0x6f, 0x63, 0x69, 0x61, 0x6c, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x12, 0x1c, 0x0a, 0x09,
	0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x12, 0x26, 0x0a, 0x0e, 0x49, 0x73,
	0x43, 0x68, 0x69, 0x6c, 0x64, 0x42, 0x65, 0x61, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x0e, 0x49, 0x73, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x42, 0x65, 0x61, 0x72, 0x69,
	0x6e, 0x67, 0x12, 0x30, 0x0a, 0x13, 0x4d, 0x61, 0x78, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x42, 0x65,
	0x61, 0x72, 0x69, 0x6e, 0x67, 0x54, 0x69, 0x63, 0x6b, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x13, 0x4d, 0x61, 0x78, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x42, 0x65, 0x61, 0x72, 0x69, 0x6e, 0x67,
	0x54, 0x69, 0x63, 0x6b, 0x12, 0x22, 0x0a, 0x0c, 0x50, 0x72, 0x65, 0x67, 0x6e, 0x61, 0x6e, 0x63,
	0x79, 0x45, 0x6e, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x50, 0x72, 0x65, 0x67,
	0x6e, 0x61, 0x6e, 0x63, 0x79, 0x45, 0x6e, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x4d, 0x61, 0x6c, 0x65,
	0x49, 0x44, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4d, 0x61, 0x6c, 0x65, 0x49, 0x44,
	0x12, 0x1a, 0x0a, 0x08, 0x46, 0x65, 0x6d, 0x61, 0x6c, 0x65, 0x49, 0x44, 0x18, 0x0c, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x46, 0x65, 0x6d, 0x61, 0x6c, 0x65, 0x49, 0x44, 0x12, 0x20, 0x0a, 0x0b,
	0x4d, 0x61, 0x47, 0x72, 0x61, 0x6e, 0x64, 0x6d, 0x61, 0x49, 0x44, 0x18, 0x0d, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x4d, 0x61, 0x47, 0x72, 0x61, 0x6e, 0x64, 0x6d, 0x61, 0x49, 0x44, 0x12, 0x20,
	0x0a, 0x0b, 0x4d, 0x61, 0x47, 0x72, 0x61, 0x6e, 0x64, 0x70, 0x61, 0x49, 0x44, 0x18, 0x0e, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x4d, 0x61, 0x47, 0x72, 0x61, 0x6e, 0x64, 0x70, 0x61, 0x49, 0x44,
	0x12, 0x20, 0x0a, 0x0b, 0x50, 0x61, 0x47, 0x72, 0x61, 0x6e, 0x64, 0x6d, 0x61, 0x49, 0x44, 0x18,
	0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x50, 0x61, 0x47, 0x72, 0x61, 0x6e, 0x64, 0x6d, 0x61,
	0x49, 0x44, 0x12, 0x20, 0x0a, 0x0b, 0x50, 0x61, 0x47, 0x72, 0x61, 0x6e, 0x64, 0x70, 0x61, 0x49,
	0x44, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x50, 0x61, 0x47, 0x72, 0x61, 0x6e, 0x64,
	0x70, 0x61, 0x49, 0x44, 0x12, 0x2a, 0x0a, 0x10, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x6f, 0x66,
	0x43, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e, 0x18, 0x11, 0x20, 0x01, 0x28, 0x03, 0x52, 0x10,
	0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x6f, 0x66, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e,
	0x12, 0x22, 0x0a, 0x0c, 0x4d, 0x61, 0x72, 0x72, 0x69, 0x61, 0x67, 0x65, 0x54, 0x69, 0x63, 0x6b,
	0x18, 0x12, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x4d, 0x61, 0x72, 0x72, 0x69, 0x61, 0x67, 0x65,
	0x54, 0x69, 0x63, 0x6b, 0x12, 0x20, 0x0a, 0x0b, 0x44, 0x69, 0x76, 0x6f, 0x72, 0x63, 0x65, 0x54,
	0x69, 0x63, 0x6b, 0x18, 0x13, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x44, 0x69, 0x76, 0x6f, 0x72,
	0x63, 0x65, 0x54, 0x69, 0x63, 0x6b, 0x12, 0x20, 0x0a, 0x0b, 0x57, 0x69, 0x64, 0x6f, 0x77, 0x65,
	0x64, 0x54, 0x69, 0x63, 0x6b, 0x18, 0x14, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x57, 0x69, 0x64,
	0x6f, 0x77, 0x65, 0x64, 0x54, 0x69, 0x63, 0x6b, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x61, 0x6e, 0x64,
	0x6f, 0x6d, 0x18, 0x15, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d,
	0x42, 0x0d, 0x5a, 0x0b, 0x70, 0x6b, 0x67, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_family_proto_rawDescOnce sync.Once
	file_family_proto_rawDescData = file_family_proto_rawDesc
)

func file_family_proto_rawDescGZIP() []byte {
	file_family_proto_rawDescOnce.Do(func() {
		file_family_proto_rawDescData = protoimpl.X.CompressGZIP(file_family_proto_rawDescData)
	})
	return file_family_proto_rawDescData
}

var file_family_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_family_proto_goTypes = []interface{}{
	(*Family)(nil), // 0: Family
	(*Ethos)(nil),  // 1: Ethos
}
var file_family_proto_depIdxs = []int32{
	1, // 0: Family.Ethos:type_name -> Ethos
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_family_proto_init() }
func file_family_proto_init() {
	if File_family_proto != nil {
		return
	}
	file_ethos_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_family_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Family); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_family_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_family_proto_goTypes,
		DependencyIndexes: file_family_proto_depIdxs,
		MessageInfos:      file_family_proto_msgTypes,
	}.Build()
	File_family_proto = out.File
	file_family_proto_rawDesc = nil
	file_family_proto_goTypes = nil
	file_family_proto_depIdxs = nil
}
