// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.27.3
// source: api.proto

package structs

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_api_proto protoreflect.FileDescriptor

var file_api_proto_rawDesc = []byte{
	0x0a, 0x09, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b, 0x77, 0x6f, 0x72,
	0x6c, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x66, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0f, 0x6f, 0x6e, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x9b, 0x05, 0x0a, 0x03, 0x41, 0x50, 0x49, 0x12, 0x2f, 0x0a,
	0x06, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x73, 0x12, 0x11, 0x2e, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72,
	0x6c, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x47, 0x65, 0x74,
	0x57, 0x6f, 0x72, 0x6c, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f,
	0x0a, 0x08, 0x53, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x12, 0x10, 0x2e, 0x53, 0x65, 0x74,
	0x57, 0x6f, 0x72, 0x6c, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x53,
	0x65, 0x74, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x35, 0x0a, 0x0a, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x73, 0x12, 0x12, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x13, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x57, 0x6f, 0x72, 0x6c, 0x64, 0x12, 0x13, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x57, 0x6f,
	0x72, 0x6c, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x2f, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x12, 0x11, 0x2e, 0x47, 0x65, 0x74,
	0x41, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e,
	0x47, 0x65, 0x74, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x32, 0x0a, 0x09, 0x53, 0x65, 0x74, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x12, 0x11,
	0x2e, 0x53, 0x65, 0x74, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x12, 0x2e, 0x53, 0x65, 0x74, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x0a, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x63, 0x74,
	0x6f, 0x72, 0x73, 0x12, 0x12, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x63,
	0x74, 0x6f, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x0b,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x13, 0x2e, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x14, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x13, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a,
	0x0a, 0x53, 0x65, 0x74, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x13, 0x2e, 0x53, 0x65,
	0x74, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x14, 0x2e, 0x53, 0x65, 0x74, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74, 0x46, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x14, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x46, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x15, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x46, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x0d, 0x5a, 0x0b, 0x70, 0x6b, 0x67, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63,
	0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_api_proto_goTypes = []interface{}{
	(*GetWorldsRequest)(nil),      // 0: GetWorldsRequest
	(*SetWorldRequest)(nil),       // 1: SetWorldRequest
	(*ListWorldsRequest)(nil),     // 2: ListWorldsRequest
	(*DeleteWorldRequest)(nil),    // 3: DeleteWorldRequest
	(*GetActorsRequest)(nil),      // 4: GetActorsRequest
	(*SetActorsRequest)(nil),      // 5: SetActorsRequest
	(*ListActorsRequest)(nil),     // 6: ListActorsRequest
	(*DeleteActorRequest)(nil),    // 7: DeleteActorRequest
	(*GetFactionsRequest)(nil),    // 8: GetFactionsRequest
	(*SetFactionsRequest)(nil),    // 9: SetFactionsRequest
	(*ListFactionsRequest)(nil),   // 10: ListFactionsRequest
	(*DeleteFactionRequest)(nil),  // 11: DeleteFactionRequest
	(*GetWorldsResponse)(nil),     // 12: GetWorldsResponse
	(*SetWorldResponse)(nil),      // 13: SetWorldResponse
	(*ListWorldsResponse)(nil),    // 14: ListWorldsResponse
	(*DeleteWorldResponse)(nil),   // 15: DeleteWorldResponse
	(*GetActorsResponse)(nil),     // 16: GetActorsResponse
	(*SetActorsResponse)(nil),     // 17: SetActorsResponse
	(*ListActorsResponse)(nil),    // 18: ListActorsResponse
	(*DeleteActorResponse)(nil),   // 19: DeleteActorResponse
	(*GetFactionsResponse)(nil),   // 20: GetFactionsResponse
	(*SetFactionsResponse)(nil),   // 21: SetFactionsResponse
	(*ListFactionsResponse)(nil),  // 22: ListFactionsResponse
	(*DeleteFactionResponse)(nil), // 23: DeleteFactionResponse
}
var file_api_proto_depIdxs = []int32{
	0,  // 0: API.Worlds:input_type -> GetWorldsRequest
	1,  // 1: API.SetWorld:input_type -> SetWorldRequest
	2,  // 2: API.ListWorlds:input_type -> ListWorldsRequest
	3,  // 3: API.DeleteWorld:input_type -> DeleteWorldRequest
	4,  // 4: API.Actors:input_type -> GetActorsRequest
	5,  // 5: API.SetActors:input_type -> SetActorsRequest
	6,  // 6: API.ListActors:input_type -> ListActorsRequest
	7,  // 7: API.DeleteActor:input_type -> DeleteActorRequest
	8,  // 8: API.Factions:input_type -> GetFactionsRequest
	9,  // 9: API.SetFaction:input_type -> SetFactionsRequest
	10, // 10: API.ListFactions:input_type -> ListFactionsRequest
	11, // 11: API.DeleteFaction:input_type -> DeleteFactionRequest
	12, // 12: API.Worlds:output_type -> GetWorldsResponse
	13, // 13: API.SetWorld:output_type -> SetWorldResponse
	14, // 14: API.ListWorlds:output_type -> ListWorldsResponse
	15, // 15: API.DeleteWorld:output_type -> DeleteWorldResponse
	16, // 16: API.Actors:output_type -> GetActorsResponse
	17, // 17: API.SetActors:output_type -> SetActorsResponse
	18, // 18: API.ListActors:output_type -> ListActorsResponse
	19, // 19: API.DeleteActor:output_type -> DeleteActorResponse
	20, // 20: API.Factions:output_type -> GetFactionsResponse
	21, // 21: API.SetFaction:output_type -> SetFactionsResponse
	22, // 22: API.ListFactions:output_type -> ListFactionsResponse
	23, // 23: API.DeleteFaction:output_type -> DeleteFactionResponse
	12, // [12:24] is the sub-list for method output_type
	0,  // [0:12] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_api_proto_init() }
func file_api_proto_init() {
	if File_api_proto != nil {
		return
	}
	file_world_proto_init()
	file_actor_proto_init()
	file_faction_proto_init()
	file_on_change_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_goTypes,
		DependencyIndexes: file_api_proto_depIdxs,
	}.Build()
	File_api_proto = out.File
	file_api_proto_rawDesc = nil
	file_api_proto_goTypes = nil
	file_api_proto_depIdxs = nil
}
