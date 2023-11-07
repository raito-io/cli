// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
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
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// PluginInfo represents the information about a plugin.
type PluginInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string           `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description string           `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Version     *version.SemVer  `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	Parameters  []*ParameterInfo `protobuf:"bytes,4,rep,name=parameters,proto3" json:"parameters,omitempty"`
}

func (x *PluginInfo) Reset() {
	*x = PluginInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_util_plugin_info_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PluginInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PluginInfo) ProtoMessage() {}

func (x *PluginInfo) ProtoReflect() protoreflect.Message {
	mi := &file_util_plugin_info_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
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

// ParameterInfo contains the information about a parameter.
// This is used to inform the CLI user what command-line parameters are expected explicitly for this target (plugin).
type ParameterInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Mandatory   bool   `protobuf:"varint,3,opt,name=mandatory,proto3" json:"mandatory,omitempty"`
}

func (x *ParameterInfo) Reset() {
	*x = ParameterInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_util_plugin_info_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ParameterInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ParameterInfo) ProtoMessage() {}

func (x *ParameterInfo) ProtoReflect() protoreflect.Message {
	mi := &file_util_plugin_info_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
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

var file_util_plugin_info_proto_rawDesc = []byte{
	0x0a, 0x16, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x69, 0x6e,
	0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1a, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xae,
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
	0x6e, 0x66, 0x6f, 0x52, 0x0a, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x22,
	0x63, 0x0a, 0x0d, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x6d, 0x61, 0x6e, 0x64, 0x61, 0x74,
	0x6f, 0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x6d, 0x61, 0x6e, 0x64, 0x61,
	0x74, 0x6f, 0x72, 0x79, 0x32, 0x49, 0x0a, 0x0b, 0x49, 0x6e, 0x66, 0x6f, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x3a, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x17, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x70, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0x2e, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x42,
	0x93, 0x01, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x70, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x42, 0x09, 0x49, 0x6e, 0x66, 0x6f, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x61, 0x69,
	0x74, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x63, 0x6c, 0x69, 0x2f, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x75,
	0x74, 0x69, 0x6c, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0xa2, 0x02, 0x03, 0x55, 0x50, 0x58,
	0xaa, 0x02, 0x0b, 0x55, 0x74, 0x69, 0x6c, 0x2e, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0xca, 0x02,
	0x0b, 0x55, 0x74, 0x69, 0x6c, 0x5c, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0xe2, 0x02, 0x17, 0x55,
	0x74, 0x69, 0x6c, 0x5c, 0x50, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0c, 0x55, 0x74, 0x69, 0x6c, 0x3a, 0x3a, 0x50,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_util_plugin_info_proto_rawDescOnce sync.Once
	file_util_plugin_info_proto_rawDescData = file_util_plugin_info_proto_rawDesc
)

func file_util_plugin_info_proto_rawDescGZIP() []byte {
	file_util_plugin_info_proto_rawDescOnce.Do(func() {
		file_util_plugin_info_proto_rawDescData = protoimpl.X.CompressGZIP(file_util_plugin_info_proto_rawDescData)
	})
	return file_util_plugin_info_proto_rawDescData
}

var file_util_plugin_info_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_util_plugin_info_proto_goTypes = []interface{}{
	(*PluginInfo)(nil),     // 0: util.plugin.PluginInfo
	(*ParameterInfo)(nil),  // 1: util.plugin.ParameterInfo
	(*version.SemVer)(nil), // 2: util.version.SemVer
	(*emptypb.Empty)(nil),  // 3: google.protobuf.Empty
}
var file_util_plugin_info_proto_depIdxs = []int32{
	2, // 0: util.plugin.PluginInfo.version:type_name -> util.version.SemVer
	1, // 1: util.plugin.PluginInfo.parameters:type_name -> util.plugin.ParameterInfo
	3, // 2: util.plugin.InfoService.GetInfo:input_type -> google.protobuf.Empty
	0, // 3: util.plugin.InfoService.GetInfo:output_type -> util.plugin.PluginInfo
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_util_plugin_info_proto_init() }
func file_util_plugin_info_proto_init() {
	if File_util_plugin_info_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_util_plugin_info_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PluginInfo); i {
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
		file_util_plugin_info_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ParameterInfo); i {
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
			RawDescriptor: file_util_plugin_info_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_util_plugin_info_proto_goTypes,
		DependencyIndexes: file_util_plugin_info_proto_depIdxs,
		MessageInfos:      file_util_plugin_info_proto_msgTypes,
	}.Build()
	File_util_plugin_info_proto = out.File
	file_util_plugin_info_proto_rawDesc = nil
	file_util_plugin_info_proto_goTypes = nil
	file_util_plugin_info_proto_depIdxs = nil
}
