// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.3
// source: faction.proto

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

// Faction represents some group we would like to simulate.
// Nb. we don't assume these are the *only* factions, just that they're the
// most notable / influential / interesting.
type Faction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ethos                 *Ethos            `protobuf:"bytes,1,opt,name=Ethos,proto3" json:"Ethos,omitempty"`
	ID                    string            `protobuf:"bytes,2,opt,name=ID,proto3" json:"ID,omitempty"`
	Name                  string            `protobuf:"bytes,3,opt,name=Name,proto3" json:"Name,omitempty"`
	HomeAreaID            string            `protobuf:"bytes,4,opt,name=HomeAreaID,proto3" json:"HomeAreaID,omitempty"`                                              // where the faction is based
	HQPlotID              string            `protobuf:"bytes,5,opt,name=HQPlotID,proto3" json:"HQPlotID,omitempty"`                                                  // faction headquarters
	ActionFrequencyTicks  int64             `protobuf:"varint,6,opt,name=ActionFrequencyTicks,proto3" json:"ActionFrequencyTicks,omitempty"`                         // faction offers new jobs every X ticks
	Leadership            FactionLeadership `protobuf:"varint,7,opt,name=Leadership,proto3,enum=FactionLeadership" json:"Leadership,omitempty"`                      // how faction is run
	Structure             FactionStructure  `protobuf:"varint,8,opt,name=Structure,proto3,enum=FactionStructure" json:"Structure,omitempty"`                         // how faction is organized
	Wealth                int64             `protobuf:"varint,9,opt,name=Wealth,proto3" json:"Wealth,omitempty"`                                                     // money / liquid wealth available to spend
	Cohesion              int64             `protobuf:"varint,10,opt,name=Cohesion,proto3" json:"Cohesion,omitempty"`                                                // how well the faction sticks together
	Corruption            int64             `protobuf:"varint,11,opt,name=Corruption,proto3" json:"Corruption,omitempty"`                                            // corruption internal to the faction
	IsCovert              bool              `protobuf:"varint,12,opt,name=IsCovert,proto3" json:"IsCovert,omitempty"`                                                // is this faction a secret society?
	GovernmentID          string            `protobuf:"bytes,13,opt,name=GovernmentID,proto3" json:"GovernmentID,omitempty"`                                         // government this faction is under
	IsGovernment          bool              `protobuf:"varint,14,opt,name=IsGovernment,proto3" json:"IsGovernment,omitempty"`                                        // is this faction a government?
	ReligionID            string            `protobuf:"bytes,15,opt,name=ReligionID,proto3" json:"ReligionID,omitempty"`                                             // religion this faction is under
	IsReligion            bool              `protobuf:"varint,16,opt,name=IsReligion,proto3" json:"IsReligion,omitempty"`                                            // is this faction a religion?
	IsMemberByBirth       bool              `protobuf:"varint,17,opt,name=IsMemberByBirth,proto3" json:"IsMemberByBirth,omitempty"`                                  // if you have a parent(s) in the faction, are you a member?
	EspionageOffense      int64             `protobuf:"varint,18,opt,name=EspionageOffense,proto3" json:"EspionageOffense,omitempty"`                                // how good is this faction at spying
	EspionageDefense      int64             `protobuf:"varint,19,opt,name=EspionageDefense,proto3" json:"EspionageDefense,omitempty"`                                // how good is this faction at not being spied on
	MilitaryOffense       int64             `protobuf:"varint,20,opt,name=MilitaryOffense,proto3" json:"MilitaryOffense,omitempty"`                                  // how good is this faction at offensive military actions
	MilitaryDefense       int64             `protobuf:"varint,21,opt,name=MilitaryDefense,proto3" json:"MilitaryDefense,omitempty"`                                  // how good is this faction at defensive military actions
	ParentFactionID       string            `protobuf:"bytes,22,opt,name=ParentFactionID,proto3" json:"ParentFactionID,omitempty"`                                   // ID of parent faction if any
	ParentFactionRelation FactionRelation   `protobuf:"varint,23,opt,name=ParentFactionRelation,proto3,enum=FactionRelation" json:"ParentFactionRelation,omitempty"` // relation to parent faction
	// Numbers are best-effort estimates of the size of the faction.
	Members int64 `protobuf:"varint,24,opt,name=Members,proto3" json:"Members,omitempty"` // number of members
	Vassals int64 `protobuf:"varint,25,opt,name=Vassals,proto3" json:"Vassals,omitempty"` // number of vassals (rough guess)
	Plots   int64 `protobuf:"varint,26,opt,name=Plots,proto3" json:"Plots,omitempty"`     // number of plots owned by the faction
	Areas   int64 `protobuf:"varint,27,opt,name=Areas,proto3" json:"Areas,omitempty"`     // number of areas the faction is active in
}

func (x *Faction) Reset() {
	*x = Faction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_faction_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Faction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Faction) ProtoMessage() {}

func (x *Faction) ProtoReflect() protoreflect.Message {
	mi := &file_faction_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Faction.ProtoReflect.Descriptor instead.
func (*Faction) Descriptor() ([]byte, []int) {
	return file_faction_proto_rawDescGZIP(), []int{0}
}

func (x *Faction) GetEthos() *Ethos {
	if x != nil {
		return x.Ethos
	}
	return nil
}

func (x *Faction) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Faction) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Faction) GetHomeAreaID() string {
	if x != nil {
		return x.HomeAreaID
	}
	return ""
}

func (x *Faction) GetHQPlotID() string {
	if x != nil {
		return x.HQPlotID
	}
	return ""
}

func (x *Faction) GetActionFrequencyTicks() int64 {
	if x != nil {
		return x.ActionFrequencyTicks
	}
	return 0
}

func (x *Faction) GetLeadership() FactionLeadership {
	if x != nil {
		return x.Leadership
	}
	return FactionLeadership_Single
}

func (x *Faction) GetStructure() FactionStructure {
	if x != nil {
		return x.Structure
	}
	return FactionStructure_Pyramid
}

func (x *Faction) GetWealth() int64 {
	if x != nil {
		return x.Wealth
	}
	return 0
}

func (x *Faction) GetCohesion() int64 {
	if x != nil {
		return x.Cohesion
	}
	return 0
}

func (x *Faction) GetCorruption() int64 {
	if x != nil {
		return x.Corruption
	}
	return 0
}

func (x *Faction) GetIsCovert() bool {
	if x != nil {
		return x.IsCovert
	}
	return false
}

func (x *Faction) GetGovernmentID() string {
	if x != nil {
		return x.GovernmentID
	}
	return ""
}

func (x *Faction) GetIsGovernment() bool {
	if x != nil {
		return x.IsGovernment
	}
	return false
}

func (x *Faction) GetReligionID() string {
	if x != nil {
		return x.ReligionID
	}
	return ""
}

func (x *Faction) GetIsReligion() bool {
	if x != nil {
		return x.IsReligion
	}
	return false
}

func (x *Faction) GetIsMemberByBirth() bool {
	if x != nil {
		return x.IsMemberByBirth
	}
	return false
}

func (x *Faction) GetEspionageOffense() int64 {
	if x != nil {
		return x.EspionageOffense
	}
	return 0
}

func (x *Faction) GetEspionageDefense() int64 {
	if x != nil {
		return x.EspionageDefense
	}
	return 0
}

func (x *Faction) GetMilitaryOffense() int64 {
	if x != nil {
		return x.MilitaryOffense
	}
	return 0
}

func (x *Faction) GetMilitaryDefense() int64 {
	if x != nil {
		return x.MilitaryDefense
	}
	return 0
}

func (x *Faction) GetParentFactionID() string {
	if x != nil {
		return x.ParentFactionID
	}
	return ""
}

func (x *Faction) GetParentFactionRelation() FactionRelation {
	if x != nil {
		return x.ParentFactionRelation
	}
	return FactionRelation_Tributary
}

func (x *Faction) GetMembers() int64 {
	if x != nil {
		return x.Members
	}
	return 0
}

func (x *Faction) GetVassals() int64 {
	if x != nil {
		return x.Vassals
	}
	return 0
}

func (x *Faction) GetPlots() int64 {
	if x != nil {
		return x.Plots
	}
	return 0
}

func (x *Faction) GetAreas() int64 {
	if x != nil {
		return x.Areas
	}
	return 0
}

// FactionSummary is a high level overview of a faction, including related tuples
// (weights), information on current faction leadership (Ranks) and any research (ResearchProgress).
type FactionSummary struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Faction *Faction `protobuf:"bytes,1,opt,name=Faction,proto3" json:"Faction,omitempty"`
	// Amassed research
	ResearchProgress map[string]int64 `protobuf:"bytes,2,rep,name=ResearchProgress,proto3" json:"ResearchProgress,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	// weights
	Professions map[string]int64 `protobuf:"bytes,3,rep,name=Professions,proto3" json:"Professions,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Actions     map[string]int64 `protobuf:"bytes,4,rep,name=Actions,proto3" json:"Actions,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Research    map[string]int64 `protobuf:"bytes,5,rep,name=Research,proto3" json:"Research,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Trust       map[string]int64 `protobuf:"bytes,6,rep,name=Trust,proto3" json:"Trust,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	// counts of people of each rank
	Ranks *DemographicRankSpread `protobuf:"bytes,7,opt,name=Ranks,proto3" json:"Ranks,omitempty"`
}

func (x *FactionSummary) Reset() {
	*x = FactionSummary{}
	if protoimpl.UnsafeEnabled {
		mi := &file_faction_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FactionSummary) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FactionSummary) ProtoMessage() {}

func (x *FactionSummary) ProtoReflect() protoreflect.Message {
	mi := &file_faction_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FactionSummary.ProtoReflect.Descriptor instead.
func (*FactionSummary) Descriptor() ([]byte, []int) {
	return file_faction_proto_rawDescGZIP(), []int{1}
}

func (x *FactionSummary) GetFaction() *Faction {
	if x != nil {
		return x.Faction
	}
	return nil
}

func (x *FactionSummary) GetResearchProgress() map[string]int64 {
	if x != nil {
		return x.ResearchProgress
	}
	return nil
}

func (x *FactionSummary) GetProfessions() map[string]int64 {
	if x != nil {
		return x.Professions
	}
	return nil
}

func (x *FactionSummary) GetActions() map[string]int64 {
	if x != nil {
		return x.Actions
	}
	return nil
}

func (x *FactionSummary) GetResearch() map[string]int64 {
	if x != nil {
		return x.Research
	}
	return nil
}

func (x *FactionSummary) GetTrust() map[string]int64 {
	if x != nil {
		return x.Trust
	}
	return nil
}

func (x *FactionSummary) GetRanks() *DemographicRankSpread {
	if x != nil {
		return x.Ranks
	}
	return nil
}

var File_faction_proto protoreflect.FileDescriptor

var file_faction_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x66, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x0b, 0x65, 0x74, 0x68, 0x6f, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x64, 0x65,
	0x6d, 0x6f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x69, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x18, 0x66, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x73, 0x68, 0x69, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x66, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x66, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x6c,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc0, 0x07, 0x0a, 0x07,
	0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x05, 0x45, 0x74, 0x68, 0x6f, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x45, 0x74, 0x68, 0x6f, 0x73, 0x52, 0x05,
	0x45, 0x74, 0x68, 0x6f, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x48, 0x6f, 0x6d,
	0x65, 0x41, 0x72, 0x65, 0x61, 0x49, 0x44, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x48,
	0x6f, 0x6d, 0x65, 0x41, 0x72, 0x65, 0x61, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x48, 0x51, 0x50,
	0x6c, 0x6f, 0x74, 0x49, 0x44, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x48, 0x51, 0x50,
	0x6c, 0x6f, 0x74, 0x49, 0x44, 0x12, 0x32, 0x0a, 0x14, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x46,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x79, 0x54, 0x69, 0x63, 0x6b, 0x73, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x14, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x6e, 0x63, 0x79, 0x54, 0x69, 0x63, 0x6b, 0x73, 0x12, 0x32, 0x0a, 0x0a, 0x4c, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e,
	0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x68, 0x69,
	0x70, 0x52, 0x0a, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x12, 0x2f, 0x0a,
	0x09, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x11, 0x2e, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74,
	0x75, 0x72, 0x65, 0x52, 0x09, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x57, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x57, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x6f, 0x68, 0x65, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x43, 0x6f, 0x68, 0x65, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x43, 0x6f, 0x72, 0x72, 0x75, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x0b, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x43, 0x6f, 0x72, 0x72, 0x75, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x49, 0x73, 0x43, 0x6f, 0x76, 0x65, 0x72, 0x74, 0x18, 0x0c,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x49, 0x73, 0x43, 0x6f, 0x76, 0x65, 0x72, 0x74, 0x12, 0x22,
	0x0a, 0x0c, 0x47, 0x6f, 0x76, 0x65, 0x72, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x0d,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x47, 0x6f, 0x76, 0x65, 0x72, 0x6e, 0x6d, 0x65, 0x6e, 0x74,
	0x49, 0x44, 0x12, 0x22, 0x0a, 0x0c, 0x49, 0x73, 0x47, 0x6f, 0x76, 0x65, 0x72, 0x6e, 0x6d, 0x65,
	0x6e, 0x74, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x49, 0x73, 0x47, 0x6f, 0x76, 0x65,
	0x72, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x52, 0x65, 0x6c, 0x69, 0x67, 0x69,
	0x6f, 0x6e, 0x49, 0x44, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x52, 0x65, 0x6c, 0x69,
	0x67, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x12, 0x1e, 0x0a, 0x0a, 0x49, 0x73, 0x52, 0x65, 0x6c, 0x69,
	0x67, 0x69, 0x6f, 0x6e, 0x18, 0x10, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x49, 0x73, 0x52, 0x65,
	0x6c, 0x69, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x0a, 0x0f, 0x49, 0x73, 0x4d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x42, 0x79, 0x42, 0x69, 0x72, 0x74, 0x68, 0x18, 0x11, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0f, 0x49, 0x73, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x42, 0x79, 0x42, 0x69, 0x72, 0x74, 0x68,
	0x12, 0x2a, 0x0a, 0x10, 0x45, 0x73, 0x70, 0x69, 0x6f, 0x6e, 0x61, 0x67, 0x65, 0x4f, 0x66, 0x66,
	0x65, 0x6e, 0x73, 0x65, 0x18, 0x12, 0x20, 0x01, 0x28, 0x03, 0x52, 0x10, 0x45, 0x73, 0x70, 0x69,
	0x6f, 0x6e, 0x61, 0x67, 0x65, 0x4f, 0x66, 0x66, 0x65, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x10,
	0x45, 0x73, 0x70, 0x69, 0x6f, 0x6e, 0x61, 0x67, 0x65, 0x44, 0x65, 0x66, 0x65, 0x6e, 0x73, 0x65,
	0x18, 0x13, 0x20, 0x01, 0x28, 0x03, 0x52, 0x10, 0x45, 0x73, 0x70, 0x69, 0x6f, 0x6e, 0x61, 0x67,
	0x65, 0x44, 0x65, 0x66, 0x65, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x0f, 0x4d, 0x69, 0x6c, 0x69,
	0x74, 0x61, 0x72, 0x79, 0x4f, 0x66, 0x66, 0x65, 0x6e, 0x73, 0x65, 0x18, 0x14, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0f, 0x4d, 0x69, 0x6c, 0x69, 0x74, 0x61, 0x72, 0x79, 0x4f, 0x66, 0x66, 0x65, 0x6e,
	0x73, 0x65, 0x12, 0x28, 0x0a, 0x0f, 0x4d, 0x69, 0x6c, 0x69, 0x74, 0x61, 0x72, 0x79, 0x44, 0x65,
	0x66, 0x65, 0x6e, 0x73, 0x65, 0x18, 0x15, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0f, 0x4d, 0x69, 0x6c,
	0x69, 0x74, 0x61, 0x72, 0x79, 0x44, 0x65, 0x66, 0x65, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x0f,
	0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x18,
	0x16, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x46, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x12, 0x46, 0x0a, 0x15, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74,
	0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x17, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x15, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x46,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18,
	0x0a, 0x07, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x18, 0x18, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x07, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x56, 0x61, 0x73, 0x73,
	0x61, 0x6c, 0x73, 0x18, 0x19, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x56, 0x61, 0x73, 0x73, 0x61,
	0x6c, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x50, 0x6c, 0x6f, 0x74, 0x73, 0x18, 0x1a, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x05, 0x50, 0x6c, 0x6f, 0x74, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x41, 0x72, 0x65, 0x61,
	0x73, 0x18, 0x1b, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x41, 0x72, 0x65, 0x61, 0x73, 0x22, 0xd6,
	0x05, 0x0a, 0x0e, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72,
	0x79, 0x12, 0x22, 0x0a, 0x07, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x08, 0x2e, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x46, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x51, 0x0a, 0x10, 0x52, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x25, 0x2e, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79,
	0x2e, 0x52, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x10, 0x52, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x50, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x12, 0x42, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x66,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e,
	0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x2e, 0x50,
	0x72, 0x6f, 0x66, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x0b, 0x50, 0x72, 0x6f, 0x66, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x36, 0x0a, 0x07,
	0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e,
	0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x2e, 0x41,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x41, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x39, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x2e, 0x52, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x52, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x12,
	0x30, 0x0a, 0x05, 0x54, 0x72, 0x75, 0x73, 0x74, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x2e,
	0x54, 0x72, 0x75, 0x73, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x54, 0x72, 0x75, 0x73,
	0x74, 0x12, 0x2c, 0x0a, 0x05, 0x52, 0x61, 0x6e, 0x6b, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x16, 0x2e, 0x44, 0x65, 0x6d, 0x6f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x69, 0x63, 0x52, 0x61,
	0x6e, 0x6b, 0x53, 0x70, 0x72, 0x65, 0x61, 0x64, 0x52, 0x05, 0x52, 0x61, 0x6e, 0x6b, 0x73, 0x1a,
	0x43, 0x0a, 0x15, 0x52, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x50, 0x72, 0x6f, 0x67, 0x72,
	0x65, 0x73, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3e, 0x0a, 0x10, 0x50, 0x72, 0x6f, 0x66, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3a, 0x0a, 0x0c, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x1a, 0x3b, 0x0a, 0x0d, 0x52, 0x65, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x38, 0x0a,
	0x0a, 0x54, 0x72, 0x75, 0x73, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x0d, 0x5a, 0x0b, 0x70, 0x6b, 0x67, 0x2f, 0x73,
	0x74, 0x72, 0x75, 0x63, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_faction_proto_rawDescOnce sync.Once
	file_faction_proto_rawDescData = file_faction_proto_rawDesc
)

func file_faction_proto_rawDescGZIP() []byte {
	file_faction_proto_rawDescOnce.Do(func() {
		file_faction_proto_rawDescData = protoimpl.X.CompressGZIP(file_faction_proto_rawDescData)
	})
	return file_faction_proto_rawDescData
}

var file_faction_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_faction_proto_goTypes = []interface{}{
	(*Faction)(nil),               // 0: Faction
	(*FactionSummary)(nil),        // 1: FactionSummary
	nil,                           // 2: FactionSummary.ResearchProgressEntry
	nil,                           // 3: FactionSummary.ProfessionsEntry
	nil,                           // 4: FactionSummary.ActionsEntry
	nil,                           // 5: FactionSummary.ResearchEntry
	nil,                           // 6: FactionSummary.TrustEntry
	(*Ethos)(nil),                 // 7: Ethos
	(FactionLeadership)(0),        // 8: FactionLeadership
	(FactionStructure)(0),         // 9: FactionStructure
	(FactionRelation)(0),          // 10: FactionRelation
	(*DemographicRankSpread)(nil), // 11: DemographicRankSpread
}
var file_faction_proto_depIdxs = []int32{
	7,  // 0: Faction.Ethos:type_name -> Ethos
	8,  // 1: Faction.Leadership:type_name -> FactionLeadership
	9,  // 2: Faction.Structure:type_name -> FactionStructure
	10, // 3: Faction.ParentFactionRelation:type_name -> FactionRelation
	0,  // 4: FactionSummary.Faction:type_name -> Faction
	2,  // 5: FactionSummary.ResearchProgress:type_name -> FactionSummary.ResearchProgressEntry
	3,  // 6: FactionSummary.Professions:type_name -> FactionSummary.ProfessionsEntry
	4,  // 7: FactionSummary.Actions:type_name -> FactionSummary.ActionsEntry
	5,  // 8: FactionSummary.Research:type_name -> FactionSummary.ResearchEntry
	6,  // 9: FactionSummary.Trust:type_name -> FactionSummary.TrustEntry
	11, // 10: FactionSummary.Ranks:type_name -> DemographicRankSpread
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_faction_proto_init() }
func file_faction_proto_init() {
	if File_faction_proto != nil {
		return
	}
	file_ethos_proto_init()
	file_demographics_proto_init()
	file_faction_leadership_proto_init()
	file_faction_structure_proto_init()
	file_faction_relation_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_faction_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Faction); i {
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
		file_faction_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FactionSummary); i {
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
			RawDescriptor: file_faction_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_faction_proto_goTypes,
		DependencyIndexes: file_faction_proto_depIdxs,
		MessageInfos:      file_faction_proto_msgTypes,
	}.Build()
	File_faction_proto = out.File
	file_faction_proto_rawDesc = nil
	file_faction_proto_goTypes = nil
	file_faction_proto_depIdxs = nil
}
