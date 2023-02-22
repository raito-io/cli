// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: data_usage/data_usage.proto

package data_usage

import (
	config "github.com/raito-io/cli/base/util/config"
	error1 "github.com/raito-io/cli/base/util/error"
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

// DataUsageSyncConfig represents the configuration that is passed from the CLI to the DataUsageSyncer plugin interface.
// It contains all the necessary configuration parameters for the plugin to function.
type DataUsageSyncConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConfigMap  *config.ConfigMap `protobuf:"bytes,1,opt,name=config_map,json=configMap,proto3" json:"config_map,omitempty"`
	TargetFile string            `protobuf:"bytes,2,opt,name=target_file,json=targetFile,proto3" json:"target_file,omitempty"`
}

func (x *DataUsageSyncConfig) Reset() {
	*x = DataUsageSyncConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_usage_data_usage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataUsageSyncConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataUsageSyncConfig) ProtoMessage() {}

func (x *DataUsageSyncConfig) ProtoReflect() protoreflect.Message {
	mi := &file_data_usage_data_usage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataUsageSyncConfig.ProtoReflect.Descriptor instead.
func (*DataUsageSyncConfig) Descriptor() ([]byte, []int) {
	return file_data_usage_data_usage_proto_rawDescGZIP(), []int{0}
}

func (x *DataUsageSyncConfig) GetConfigMap() *config.ConfigMap {
	if x != nil {
		return x.ConfigMap
	}
	return nil
}

func (x *DataUsageSyncConfig) GetTargetFile() string {
	if x != nil {
		return x.TargetFile
	}
	return ""
}

// DataUsageSyncResult represents the result from the data usage sync process.
// A potential error is also modeled in here so specific errors remain intact when passed over RPC.
type DataUsageSyncResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error *error1.ErrorResult `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *DataUsageSyncResult) Reset() {
	*x = DataUsageSyncResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_usage_data_usage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataUsageSyncResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataUsageSyncResult) ProtoMessage() {}

func (x *DataUsageSyncResult) ProtoReflect() protoreflect.Message {
	mi := &file_data_usage_data_usage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataUsageSyncResult.ProtoReflect.Descriptor instead.
func (*DataUsageSyncResult) Descriptor() ([]byte, []int) {
	return file_data_usage_data_usage_proto_rawDescGZIP(), []int{1}
}

func (x *DataUsageSyncResult) GetError() *error1.ErrorResult {
	if x != nil {
		return x.Error
	}
	return nil
}

var File_data_usage_data_usage_proto protoreflect.FileDescriptor

var file_data_usage_data_usage_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2f, 0x64, 0x61, 0x74,
	0x61, 0x5f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x64,
	0x61, 0x74, 0x61, 0x5f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x18, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x16, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2f, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6d, 0x0a, 0x13, 0x44, 0x61, 0x74, 0x61, 0x55, 0x73, 0x61, 0x67,
	0x65, 0x53, 0x79, 0x6e, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x35, 0x0a, 0x0a, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x6d, 0x61, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x4d, 0x61, 0x70, 0x52, 0x09, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4d,
	0x61, 0x70, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x66, 0x69, 0x6c,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x46,
	0x69, 0x6c, 0x65, 0x22, 0x44, 0x0a, 0x13, 0x44, 0x61, 0x74, 0x61, 0x55, 0x73, 0x61, 0x67, 0x65,
	0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x2d, 0x0a, 0x05, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x75, 0x74, 0x69, 0x6c,
	0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x32, 0xbd, 0x01, 0x0a, 0x14, 0x44, 0x61,
	0x74, 0x61, 0x55, 0x73, 0x61, 0x67, 0x65, 0x53, 0x79, 0x6e, 0x63, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x52, 0x0a, 0x15, 0x43, 0x6c, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x21, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x2e, 0x43, 0x6c, 0x69, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x72,
	0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x51, 0x0a, 0x0d, 0x53, 0x79, 0x6e, 0x63, 0x44, 0x61,
	0x74, 0x61, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1f, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x75,
	0x73, 0x61, 0x67, 0x65, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x55, 0x73, 0x61, 0x67, 0x65, 0x53, 0x79,
	0x6e, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a, 0x1f, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x5f,
	0x75, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x55, 0x73, 0x61, 0x67, 0x65, 0x53,
	0x79, 0x6e, 0x63, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42, 0x8d, 0x01, 0x0a, 0x0e, 0x63, 0x6f,
	0x6d, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x42, 0x0e, 0x44, 0x61,
	0x74, 0x61, 0x55, 0x73, 0x61, 0x67, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x27,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x61, 0x69, 0x74, 0x6f,
	0x2d, 0x69, 0x6f, 0x2f, 0x63, 0x6c, 0x69, 0x2f, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x64, 0x61, 0x74,
	0x61, 0x5f, 0x75, 0x73, 0x61, 0x67, 0x65, 0xa2, 0x02, 0x03, 0x44, 0x58, 0x58, 0xaa, 0x02, 0x09,
	0x44, 0x61, 0x74, 0x61, 0x55, 0x73, 0x61, 0x67, 0x65, 0xca, 0x02, 0x09, 0x44, 0x61, 0x74, 0x61,
	0x55, 0x73, 0x61, 0x67, 0x65, 0xe2, 0x02, 0x15, 0x44, 0x61, 0x74, 0x61, 0x55, 0x73, 0x61, 0x67,
	0x65, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x09,
	0x44, 0x61, 0x74, 0x61, 0x55, 0x73, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_data_usage_data_usage_proto_rawDescOnce sync.Once
	file_data_usage_data_usage_proto_rawDescData = file_data_usage_data_usage_proto_rawDesc
)

func file_data_usage_data_usage_proto_rawDescGZIP() []byte {
	file_data_usage_data_usage_proto_rawDescOnce.Do(func() {
		file_data_usage_data_usage_proto_rawDescData = protoimpl.X.CompressGZIP(file_data_usage_data_usage_proto_rawDescData)
	})
	return file_data_usage_data_usage_proto_rawDescData
}

var file_data_usage_data_usage_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_data_usage_data_usage_proto_goTypes = []interface{}{
	(*DataUsageSyncConfig)(nil),         // 0: data_usage.DataUsageSyncConfig
	(*DataUsageSyncResult)(nil),         // 1: data_usage.DataUsageSyncResult
	(*config.ConfigMap)(nil),            // 2: util.config.ConfigMap
	(*error1.ErrorResult)(nil),          // 3: util.error.ErrorResult
	(*emptypb.Empty)(nil),               // 4: google.protobuf.Empty
	(*version.CliBuildInformation)(nil), // 5: util.version.CliBuildInformation
}
var file_data_usage_data_usage_proto_depIdxs = []int32{
	2, // 0: data_usage.DataUsageSyncConfig.config_map:type_name -> util.config.ConfigMap
	3, // 1: data_usage.DataUsageSyncResult.error:type_name -> util.error.ErrorResult
	4, // 2: data_usage.DataUsageSyncService.CliVersionInformation:input_type -> google.protobuf.Empty
	0, // 3: data_usage.DataUsageSyncService.SyncDataUsage:input_type -> data_usage.DataUsageSyncConfig
	5, // 4: data_usage.DataUsageSyncService.CliVersionInformation:output_type -> util.version.CliBuildInformation
	1, // 5: data_usage.DataUsageSyncService.SyncDataUsage:output_type -> data_usage.DataUsageSyncResult
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_data_usage_data_usage_proto_init() }
func file_data_usage_data_usage_proto_init() {
	if File_data_usage_data_usage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_data_usage_data_usage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataUsageSyncConfig); i {
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
		file_data_usage_data_usage_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataUsageSyncResult); i {
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
			RawDescriptor: file_data_usage_data_usage_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_data_usage_data_usage_proto_goTypes,
		DependencyIndexes: file_data_usage_data_usage_proto_depIdxs,
		MessageInfos:      file_data_usage_data_usage_proto_msgTypes,
	}.Build()
	File_data_usage_data_usage_proto = out.File
	file_data_usage_data_usage_proto_rawDesc = nil
	file_data_usage_data_usage_proto_goTypes = nil
	file_data_usage_data_usage_proto_depIdxs = nil
}
