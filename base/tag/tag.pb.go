// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: tag/tag.proto

package tag

import (
	config "github.com/raito-io/cli/base/util/config"
	version "github.com/raito-io/cli/base/util/version"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TagSyncConfig struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	ConfigMap       *config.ConfigMap      `protobuf:"bytes,1,opt,name=config_map,json=configMap,proto3" json:"config_map,omitempty"`
	TargetFile      string                 `protobuf:"bytes,2,opt,name=target_file,json=targetFile,proto3" json:"target_file,omitempty"`
	DataSourceId    string                 `protobuf:"bytes,3,opt,name=data_source_id,json=dataSourceId,proto3" json:"data_source_id,omitempty"`
	IdentityStoreId string                 `protobuf:"bytes,4,opt,name=identity_store_id,json=identityStoreId,proto3" json:"identity_store_id,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *TagSyncConfig) Reset() {
	*x = TagSyncConfig{}
	mi := &file_tag_tag_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TagSyncConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TagSyncConfig) ProtoMessage() {}

func (x *TagSyncConfig) ProtoReflect() protoreflect.Message {
	mi := &file_tag_tag_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TagSyncConfig.ProtoReflect.Descriptor instead.
func (*TagSyncConfig) Descriptor() ([]byte, []int) {
	return file_tag_tag_proto_rawDescGZIP(), []int{0}
}

func (x *TagSyncConfig) GetConfigMap() *config.ConfigMap {
	if x != nil {
		return x.ConfigMap
	}
	return nil
}

func (x *TagSyncConfig) GetTargetFile() string {
	if x != nil {
		return x.TargetFile
	}
	return ""
}

func (x *TagSyncConfig) GetDataSourceId() string {
	if x != nil {
		return x.DataSourceId
	}
	return ""
}

func (x *TagSyncConfig) GetIdentityStoreId() string {
	if x != nil {
		return x.IdentityStoreId
	}
	return ""
}

type TagSyncResult struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	Tags            int32                  `protobuf:"varint,1,opt,name=tags,proto3" json:"tags,omitempty"`
	TagSourcesScope []string               `protobuf:"bytes,2,rep,name=tag_sources_scope,json=tagSourcesScope,proto3" json:"tag_sources_scope,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *TagSyncResult) Reset() {
	*x = TagSyncResult{}
	mi := &file_tag_tag_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TagSyncResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TagSyncResult) ProtoMessage() {}

func (x *TagSyncResult) ProtoReflect() protoreflect.Message {
	mi := &file_tag_tag_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TagSyncResult.ProtoReflect.Descriptor instead.
func (*TagSyncResult) Descriptor() ([]byte, []int) {
	return file_tag_tag_proto_rawDescGZIP(), []int{1}
}

func (x *TagSyncResult) GetTags() int32 {
	if x != nil {
		return x.Tags
	}
	return 0
}

func (x *TagSyncResult) GetTagSourcesScope() []string {
	if x != nil {
		return x.TagSourcesScope
	}
	return nil
}

var File_tag_tag_proto protoreflect.FileDescriptor

var file_tag_tag_proto_rawDesc = string([]byte{
	0x0a, 0x0d, 0x74, 0x61, 0x67, 0x2f, 0x74, 0x61, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x03, 0x74, 0x61, 0x67, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x18, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x75, 0x74, 0x69,
	0x6c, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb9, 0x01, 0x0a, 0x0d, 0x54, 0x61, 0x67, 0x53,
	0x79, 0x6e, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x35, 0x0a, 0x0a, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x5f, 0x6d, 0x61, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x75, 0x74, 0x69, 0x6c, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x4d, 0x61, 0x70, 0x52, 0x09, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4d, 0x61, 0x70,
	0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x46, 0x69, 0x6c,
	0x65, 0x12, 0x24, 0x0a, 0x0e, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x64, 0x61, 0x74, 0x61, 0x53,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x2a, 0x0a, 0x11, 0x69, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x53, 0x74, 0x6f, 0x72,
	0x65, 0x49, 0x64, 0x22, 0x4f, 0x0a, 0x0d, 0x54, 0x61, 0x67, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12, 0x2a, 0x0a, 0x11, 0x74, 0x61, 0x67, 0x5f,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x5f, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x0f, 0x74, 0x61, 0x67, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x53,
	0x63, 0x6f, 0x70, 0x65, 0x32, 0x98, 0x01, 0x0a, 0x0e, 0x54, 0x61, 0x67, 0x53, 0x79, 0x6e, 0x63,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x52, 0x0a, 0x15, 0x43, 0x6c, 0x69, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x21, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6c, 0x69, 0x42, 0x75, 0x69, 0x6c, 0x64,
	0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x32, 0x0a, 0x08, 0x53,
	0x79, 0x6e, 0x63, 0x54, 0x61, 0x67, 0x73, 0x12, 0x12, 0x2e, 0x74, 0x61, 0x67, 0x2e, 0x54, 0x61,
	0x67, 0x53, 0x79, 0x6e, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a, 0x12, 0x2e, 0x74, 0x61,
	0x67, 0x2e, 0x54, 0x61, 0x67, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42,
	0x61, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x61, 0x67, 0x42, 0x08, 0x54, 0x61, 0x67, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x20, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x72, 0x61, 0x69, 0x74, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x63, 0x6c, 0x69, 0x2f,
	0x62, 0x61, 0x73, 0x65, 0x2f, 0x74, 0x61, 0x67, 0xa2, 0x02, 0x03, 0x54, 0x58, 0x58, 0xaa, 0x02,
	0x03, 0x54, 0x61, 0x67, 0xca, 0x02, 0x03, 0x54, 0x61, 0x67, 0xe2, 0x02, 0x0f, 0x54, 0x61, 0x67,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x03, 0x54,
	0x61, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_tag_tag_proto_rawDescOnce sync.Once
	file_tag_tag_proto_rawDescData []byte
)

func file_tag_tag_proto_rawDescGZIP() []byte {
	file_tag_tag_proto_rawDescOnce.Do(func() {
		file_tag_tag_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_tag_tag_proto_rawDesc), len(file_tag_tag_proto_rawDesc)))
	})
	return file_tag_tag_proto_rawDescData
}

var file_tag_tag_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_tag_tag_proto_goTypes = []any{
	(*TagSyncConfig)(nil),               // 0: tag.TagSyncConfig
	(*TagSyncResult)(nil),               // 1: tag.TagSyncResult
	(*config.ConfigMap)(nil),            // 2: util.config.ConfigMap
	(*emptypb.Empty)(nil),               // 3: google.protobuf.Empty
	(*version.CliBuildInformation)(nil), // 4: util.version.CliBuildInformation
}
var file_tag_tag_proto_depIdxs = []int32{
	2, // 0: tag.TagSyncConfig.config_map:type_name -> util.config.ConfigMap
	3, // 1: tag.TagSyncService.CliVersionInformation:input_type -> google.protobuf.Empty
	0, // 2: tag.TagSyncService.SyncTags:input_type -> tag.TagSyncConfig
	4, // 3: tag.TagSyncService.CliVersionInformation:output_type -> util.version.CliBuildInformation
	1, // 4: tag.TagSyncService.SyncTags:output_type -> tag.TagSyncResult
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_tag_tag_proto_init() }
func file_tag_tag_proto_init() {
	if File_tag_tag_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_tag_tag_proto_rawDesc), len(file_tag_tag_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tag_tag_proto_goTypes,
		DependencyIndexes: file_tag_tag_proto_depIdxs,
		MessageInfos:      file_tag_tag_proto_msgTypes,
	}.Build()
	File_tag_tag_proto = out.File
	file_tag_tag_proto_goTypes = nil
	file_tag_tag_proto_depIdxs = nil
}
