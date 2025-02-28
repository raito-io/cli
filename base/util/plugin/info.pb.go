// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: util/plugin/info.proto

package plugin

import (
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

type PluginType int32

const (
	// UNKNOWN plugin type. Avoid this type.
	// If type is unknown CLI will fallback on FULL_DS_SYNC
	PluginType_PLUGIN_TYPE_UNKNOWN PluginType = 0
	// FULL_DS_SYNC execute data source sync, identity store sync, access provider sync and data usage sync.
	// This type should be used for most data sources such as snowflake, bigquery, databricks, and so on.
	// A websocket may be initialized for this type.
	PluginType_PLUGIN_TYPE_FULL_DS_SYNC PluginType = 1
	// IS_SYNC execute only an identity store sync.
	// This type should be used for syncing identity providers (such as okta).
	PluginType_PLUGIN_TYPE_IS_SYNC PluginType = 2
	// TAG_SYNC execute only a tag sync.
	// This type should be used for syncing tags on external sources (such as catalogs).
	PluginType_PLUGIN_TYPE_TAG_SYNC PluginType = 3
	// AC_PROVIDER execute only an access provider sync.
	// This type should be used for plugins that provide Raito Cloud Resources.
	PluginType_PLUGIN_TYPE_RESOURCE_PROVIDER PluginType = 4
)

// Enum value maps for PluginType.
var (
	PluginType_name = map[int32]string{
		0: "PLUGIN_TYPE_UNKNOWN",
		1: "PLUGIN_TYPE_FULL_DS_SYNC",
		2: "PLUGIN_TYPE_IS_SYNC",
		3: "PLUGIN_TYPE_TAG_SYNC",
		4: "PLUGIN_TYPE_RESOURCE_PROVIDER",
	}
	PluginType_value = map[string]int32{
		"PLUGIN_TYPE_UNKNOWN":           0,
		"PLUGIN_TYPE_FULL_DS_SYNC":      1,
		"PLUGIN_TYPE_IS_SYNC":           2,
		"PLUGIN_TYPE_TAG_SYNC":          3,
		"PLUGIN_TYPE_RESOURCE_PROVIDER": 4,
	}
)

func (x PluginType) Enum() *PluginType {
	p := new(PluginType)
	*p = x
	return p
}

func (x PluginType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PluginType) Descriptor() protoreflect.EnumDescriptor {
	return file_util_plugin_info_proto_enumTypes[0].Descriptor()
}

func (PluginType) Type() protoreflect.EnumType {
	return &file_util_plugin_info_proto_enumTypes[0]
}

func (x PluginType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PluginType.Descriptor instead.
func (PluginType) EnumDescriptor() ([]byte, []int) {
	return file_util_plugin_info_proto_rawDescGZIP(), []int{0}
}

// PluginInfo represents the information about a plugin.
type PluginInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Version       *version.SemVer        `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	Parameters    []*ParameterInfo       `protobuf:"bytes,4,rep,name=parameters,proto3" json:"parameters,omitempty"`
	TagSource     string                 `protobuf:"bytes,5,opt,name=tag_source,json=tagSource,proto3" json:"tag_source,omitempty"`
	Type          []PluginType           `protobuf:"varint,6,rep,packed,name=type,proto3,enum=util.plugin.PluginType" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PluginInfo) Reset() {
	*x = PluginInfo{}
	mi := &file_util_plugin_info_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PluginInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PluginInfo) ProtoMessage() {}

func (x *PluginInfo) ProtoReflect() protoreflect.Message {
	mi := &file_util_plugin_info_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PluginInfo.ProtoReflect.Descriptor instead.
func (*PluginInfo) Descriptor() ([]byte, []int) {
	return file_util_plugin_info_proto_rawDescGZIP(), []int{0}
}

func (x *PluginInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PluginInfo) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *PluginInfo) GetVersion() *version.SemVer {
	if x != nil {
		return x.Version
	}
	return nil
}

func (x *PluginInfo) GetParameters() []*ParameterInfo {
	if x != nil {
		return x.Parameters
	}
	return nil
}

func (x *PluginInfo) GetTagSource() string {
	if x != nil {
		return x.TagSource
	}
	return ""
}

func (x *PluginInfo) GetType() []PluginType {
	if x != nil {
		return x.Type
	}
	return nil
}

// ParameterInfo contains the information about a parameter.
// This is used to inform the CLI user what command-line parameters are expected explicitly for this target (plugin).
type ParameterInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Mandatory     bool                   `protobuf:"varint,3,opt,name=mandatory,proto3" json:"mandatory,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ParameterInfo) Reset() {
	*x = ParameterInfo{}
	mi := &file_util_plugin_info_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ParameterInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ParameterInfo) ProtoMessage() {}

func (x *ParameterInfo) ProtoReflect() protoreflect.Message {
	mi := &file_util_plugin_info_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ParameterInfo.ProtoReflect.Descriptor instead.
func (*ParameterInfo) Descriptor() ([]byte, []int) {
	return file_util_plugin_info_proto_rawDescGZIP(), []int{1}
}

func (x *ParameterInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ParameterInfo) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ParameterInfo) GetMandatory() bool {
	if x != nil {
		return x.Mandatory
	}
	return false
}

var File_util_plugin_info_proto protoreflect.FileDescriptor

var file_util_plugin_info_proto_rawDesc = string([]byte{
	0x0a, 0x16, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x69, 0x6e,
	0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1a, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xfa,
	0x01, 0x0a, 0x0a, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x2e, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x65, 0x6d, 0x56, 0x65, 0x72, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x3a, 0x0a, 0x0a, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72,
	0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x0a, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x12,
	0x1d, 0x0a, 0x0a, 0x74, 0x61, 0x67, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x61, 0x67, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x2b,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x75,
	0x74, 0x69, 0x6c, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x50, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x63, 0x0a, 0x0d, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x6d, 0x61, 0x6e, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x79, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x6d, 0x61, 0x6e, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x79,
	0x2a, 0x99, 0x01, 0x0a, 0x0a, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x17, 0x0a, 0x13, 0x50, 0x4c, 0x55, 0x47, 0x49, 0x4e, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55,
	0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x1c, 0x0a, 0x18, 0x50, 0x4c, 0x55, 0x47,
	0x49, 0x4e, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x46, 0x55, 0x4c, 0x4c, 0x5f, 0x44, 0x53, 0x5f,
	0x53, 0x59, 0x4e, 0x43, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x50, 0x4c, 0x55, 0x47, 0x49, 0x4e,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x49, 0x53, 0x5f, 0x53, 0x59, 0x4e, 0x43, 0x10, 0x02, 0x12,
	0x18, 0x0a, 0x14, 0x50, 0x4c, 0x55, 0x47, 0x49, 0x4e, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x54,
	0x41, 0x47, 0x5f, 0x53, 0x59, 0x4e, 0x43, 0x10, 0x03, 0x12, 0x21, 0x0a, 0x1d, 0x50, 0x4c, 0x55,
	0x47, 0x49, 0x4e, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43,
	0x45, 0x5f, 0x50, 0x52, 0x4f, 0x56, 0x49, 0x44, 0x45, 0x52, 0x10, 0x04, 0x32, 0x49, 0x0a, 0x0b,
	0x49, 0x6e, 0x66, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3a, 0x0a, 0x07, 0x47,
	0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x17,
	0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x50, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x42, 0x93, 0x01, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x2e,
	0x75, 0x74, 0x69, 0x6c, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x42, 0x09, 0x49, 0x6e, 0x66,
	0x6f, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x61, 0x69, 0x74, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x63, 0x6c,
	0x69, 0x2f, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x70, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0xa2, 0x02, 0x03, 0x55, 0x50, 0x58, 0xaa, 0x02, 0x0b, 0x55, 0x74, 0x69, 0x6c, 0x2e,
	0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0xca, 0x02, 0x0b, 0x55, 0x74, 0x69, 0x6c, 0x5c, 0x50, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0xe2, 0x02, 0x17, 0x55, 0x74, 0x69, 0x6c, 0x5c, 0x50, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x0c, 0x55, 0x74, 0x69, 0x6c, 0x3a, 0x3a, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_util_plugin_info_proto_rawDescOnce sync.Once
	file_util_plugin_info_proto_rawDescData []byte
)

func file_util_plugin_info_proto_rawDescGZIP() []byte {
	file_util_plugin_info_proto_rawDescOnce.Do(func() {
		file_util_plugin_info_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_util_plugin_info_proto_rawDesc), len(file_util_plugin_info_proto_rawDesc)))
	})
	return file_util_plugin_info_proto_rawDescData
}

var file_util_plugin_info_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_util_plugin_info_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_util_plugin_info_proto_goTypes = []any{
	(PluginType)(0),        // 0: util.plugin.PluginType
	(*PluginInfo)(nil),     // 1: util.plugin.PluginInfo
	(*ParameterInfo)(nil),  // 2: util.plugin.ParameterInfo
	(*version.SemVer)(nil), // 3: util.version.SemVer
	(*emptypb.Empty)(nil),  // 4: google.protobuf.Empty
}
var file_util_plugin_info_proto_depIdxs = []int32{
	3, // 0: util.plugin.PluginInfo.version:type_name -> util.version.SemVer
	2, // 1: util.plugin.PluginInfo.parameters:type_name -> util.plugin.ParameterInfo
	0, // 2: util.plugin.PluginInfo.type:type_name -> util.plugin.PluginType
	4, // 3: util.plugin.InfoService.GetInfo:input_type -> google.protobuf.Empty
	1, // 4: util.plugin.InfoService.GetInfo:output_type -> util.plugin.PluginInfo
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_util_plugin_info_proto_init() }
func file_util_plugin_info_proto_init() {
	if File_util_plugin_info_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_util_plugin_info_proto_rawDesc), len(file_util_plugin_info_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_util_plugin_info_proto_goTypes,
		DependencyIndexes: file_util_plugin_info_proto_depIdxs,
		EnumInfos:         file_util_plugin_info_proto_enumTypes,
		MessageInfos:      file_util_plugin_info_proto_msgTypes,
	}.Build()
	File_util_plugin_info_proto = out.File
	file_util_plugin_info_proto_goTypes = nil
	file_util_plugin_info_proto_depIdxs = nil
}
