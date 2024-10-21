// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: util/config/config.proto

package config

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

type ConfigMap struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Parameters map[string]string `protobuf:"bytes,1,rep,name=parameters,proto3" json:"parameters,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *ConfigMap) Reset() {
	*x = ConfigMap{}
	if protoimpl.UnsafeEnabled {
		mi := &file_util_config_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigMap) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigMap) ProtoMessage() {}

func (x *ConfigMap) ProtoReflect() protoreflect.Message {
	mi := &file_util_config_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigMap.ProtoReflect.Descriptor instead.
func (*ConfigMap) Descriptor() ([]byte, []int) {
	return file_util_config_config_proto_rawDescGZIP(), []int{0}
}

func (x *ConfigMap) GetParameters() map[string]string {
	if x != nil {
		return x.Parameters
	}
	return nil
}

var File_util_config_config_proto protoreflect.FileDescriptor

var file_util_config_config_proto_rawDesc = []byte{
	0x0a, 0x18, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x75, 0x74, 0x69, 0x6c,
	0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x92, 0x01, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x4d, 0x61, 0x70, 0x12, 0x46, 0x0a, 0x0a, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74,
	0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x75, 0x74, 0x69, 0x6c,
	0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4d, 0x61,
	0x70, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x0a, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x1a, 0x3d, 0x0a,
	0x0f, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x95, 0x01, 0x0a,
	0x0f, 0x63, 0x6f, 0x6d, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x42, 0x0b, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x61, 0x69, 0x74,
	0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x63, 0x6c, 0x69, 0x2f, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x75, 0x74,
	0x69, 0x6c, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0xa2, 0x02, 0x03, 0x55, 0x43, 0x58, 0xaa,
	0x02, 0x0b, 0x55, 0x74, 0x69, 0x6c, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0xca, 0x02, 0x0b,
	0x55, 0x74, 0x69, 0x6c, 0x5c, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0xe2, 0x02, 0x17, 0x55, 0x74,
	0x69, 0x6c, 0x5c, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0c, 0x55, 0x74, 0x69, 0x6c, 0x3a, 0x3a, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_util_config_config_proto_rawDescOnce sync.Once
	file_util_config_config_proto_rawDescData = file_util_config_config_proto_rawDesc
)

func file_util_config_config_proto_rawDescGZIP() []byte {
	file_util_config_config_proto_rawDescOnce.Do(func() {
		file_util_config_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_util_config_config_proto_rawDescData)
	})
	return file_util_config_config_proto_rawDescData
}

var file_util_config_config_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_util_config_config_proto_goTypes = []any{
	(*ConfigMap)(nil), // 0: util.config.ConfigMap
	nil,               // 1: util.config.ConfigMap.ParametersEntry
}
var file_util_config_config_proto_depIdxs = []int32{
	1, // 0: util.config.ConfigMap.parameters:type_name -> util.config.ConfigMap.ParametersEntry
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_util_config_config_proto_init() }
func file_util_config_config_proto_init() {
	if File_util_config_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_util_config_config_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ConfigMap); i {
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
			RawDescriptor: file_util_config_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_util_config_config_proto_goTypes,
		DependencyIndexes: file_util_config_config_proto_depIdxs,
		MessageInfos:      file_util_config_config_proto_msgTypes,
	}.Build()
	File_util_config_config_proto = out.File
	file_util_config_config_proto_rawDesc = nil
	file_util_config_config_proto_goTypes = nil
	file_util_config_config_proto_depIdxs = nil
}
