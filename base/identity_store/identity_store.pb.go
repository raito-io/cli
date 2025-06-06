// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: identity_store/identity_store.proto

package identity_store

import (
	config "github.com/raito-io/cli/base/util/config"
	error1 "github.com/raito-io/cli/base/util/error"
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

// IdentityStoreSyncConfig represents the configuration that is passed from the CLI to the IdentityStoreSyncer plugin interface.
// It contains all the necessary configuration parameters for the plugin to function.
type IdentityStoreSyncConfig struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ConfigMap     *config.ConfigMap      `protobuf:"bytes,1,opt,name=config_map,json=configMap,proto3" json:"config_map,omitempty"`
	UserFile      string                 `protobuf:"bytes,2,opt,name=user_file,json=userFile,proto3" json:"user_file,omitempty"`
	GroupFile     string                 `protobuf:"bytes,3,opt,name=group_file,json=groupFile,proto3" json:"group_file,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *IdentityStoreSyncConfig) Reset() {
	*x = IdentityStoreSyncConfig{}
	mi := &file_identity_store_identity_store_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IdentityStoreSyncConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IdentityStoreSyncConfig) ProtoMessage() {}

func (x *IdentityStoreSyncConfig) ProtoReflect() protoreflect.Message {
	mi := &file_identity_store_identity_store_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IdentityStoreSyncConfig.ProtoReflect.Descriptor instead.
func (*IdentityStoreSyncConfig) Descriptor() ([]byte, []int) {
	return file_identity_store_identity_store_proto_rawDescGZIP(), []int{0}
}

func (x *IdentityStoreSyncConfig) GetConfigMap() *config.ConfigMap {
	if x != nil {
		return x.ConfigMap
	}
	return nil
}

func (x *IdentityStoreSyncConfig) GetUserFile() string {
	if x != nil {
		return x.UserFile
	}
	return ""
}

func (x *IdentityStoreSyncConfig) GetGroupFile() string {
	if x != nil {
		return x.GroupFile
	}
	return ""
}

// IdentityStoreSyncResult represents the result from the identity store sync process.
// A potential error is also modeled in here so specific errors remain intact when passed over RPC.
type IdentityStoreSyncResult struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Deprecated: Marked as deprecated in identity_store/identity_store.proto.
	Error         *error1.ErrorResult `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	UserCount     int32               `protobuf:"varint,2,opt,name=user_count,json=userCount,proto3" json:"user_count,omitempty"`
	GroupCount    int32               `protobuf:"varint,3,opt,name=group_count,json=groupCount,proto3" json:"group_count,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *IdentityStoreSyncResult) Reset() {
	*x = IdentityStoreSyncResult{}
	mi := &file_identity_store_identity_store_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IdentityStoreSyncResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IdentityStoreSyncResult) ProtoMessage() {}

func (x *IdentityStoreSyncResult) ProtoReflect() protoreflect.Message {
	mi := &file_identity_store_identity_store_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IdentityStoreSyncResult.ProtoReflect.Descriptor instead.
func (*IdentityStoreSyncResult) Descriptor() ([]byte, []int) {
	return file_identity_store_identity_store_proto_rawDescGZIP(), []int{1}
}

// Deprecated: Marked as deprecated in identity_store/identity_store.proto.
func (x *IdentityStoreSyncResult) GetError() *error1.ErrorResult {
	if x != nil {
		return x.Error
	}
	return nil
}

func (x *IdentityStoreSyncResult) GetUserCount() int32 {
	if x != nil {
		return x.UserCount
	}
	return 0
}

func (x *IdentityStoreSyncResult) GetGroupCount() int32 {
	if x != nil {
		return x.GroupCount
	}
	return 0
}

type MetaData struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Type          string                 `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Icon          string                 `protobuf:"bytes,2,opt,name=icon,proto3" json:"icon,omitempty"`
	CanBeLinked   bool                   `protobuf:"varint,3,opt,name=can_be_linked,json=canBeLinked,proto3" json:"can_be_linked,omitempty"`
	CanBeMaster   bool                   `protobuf:"varint,4,opt,name=can_be_master,json=canBeMaster,proto3" json:"can_be_master,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MetaData) Reset() {
	*x = MetaData{}
	mi := &file_identity_store_identity_store_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MetaData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetaData) ProtoMessage() {}

func (x *MetaData) ProtoReflect() protoreflect.Message {
	mi := &file_identity_store_identity_store_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetaData.ProtoReflect.Descriptor instead.
func (*MetaData) Descriptor() ([]byte, []int) {
	return file_identity_store_identity_store_proto_rawDescGZIP(), []int{2}
}

func (x *MetaData) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *MetaData) GetIcon() string {
	if x != nil {
		return x.Icon
	}
	return ""
}

func (x *MetaData) GetCanBeLinked() bool {
	if x != nil {
		return x.CanBeLinked
	}
	return false
}

func (x *MetaData) GetCanBeMaster() bool {
	if x != nil {
		return x.CanBeMaster
	}
	return false
}

var File_identity_store_identity_store_proto protoreflect.FileDescriptor

const file_identity_store_identity_store_proto_rawDesc = "" +
	"\n" +
	"#identity_store/identity_store.proto\x12\x0eidentity_store\x1a\x1bgoogle/protobuf/empty.proto\x1a\x18util/config/config.proto\x1a\x16util/error/error.proto\x1a\x1autil/version/version.proto\"\x8c\x01\n" +
	"\x17IdentityStoreSyncConfig\x125\n" +
	"\n" +
	"config_map\x18\x01 \x01(\v2\x16.util.config.ConfigMapR\tconfigMap\x12\x1b\n" +
	"\tuser_file\x18\x02 \x01(\tR\buserFile\x12\x1d\n" +
	"\n" +
	"group_file\x18\x03 \x01(\tR\tgroupFile\"\x8c\x01\n" +
	"\x17IdentityStoreSyncResult\x121\n" +
	"\x05error\x18\x01 \x01(\v2\x17.util.error.ErrorResultB\x02\x18\x01R\x05error\x12\x1d\n" +
	"\n" +
	"user_count\x18\x02 \x01(\x05R\tuserCount\x12\x1f\n" +
	"\vgroup_count\x18\x03 \x01(\x05R\n" +
	"groupCount\"z\n" +
	"\bMetaData\x12\x12\n" +
	"\x04type\x18\x01 \x01(\tR\x04type\x12\x12\n" +
	"\x04icon\x18\x02 \x01(\tR\x04icon\x12\"\n" +
	"\rcan_be_linked\x18\x03 \x01(\bR\vcanBeLinked\x12\"\n" +
	"\rcan_be_master\x18\x04 \x01(\bR\vcanBeMaster2\xa3\x02\n" +
	"\x18IdentityStoreSyncService\x12R\n" +
	"\x15CliVersionInformation\x12\x16.google.protobuf.Empty\x1a!.util.version.CliBuildInformation\x12e\n" +
	"\x11SyncIdentityStore\x12'.identity_store.IdentityStoreSyncConfig\x1a'.identity_store.IdentityStoreSyncResult\x12L\n" +
	"\x18GetIdentityStoreMetaData\x12\x16.util.config.ConfigMap\x1a\x18.identity_store.MetaDataB\xa9\x01\n" +
	"\x12com.identity_storeB\x12IdentityStoreProtoP\x01Z+github.com/raito-io/cli/base/identity_store\xa2\x02\x03IXX\xaa\x02\rIdentityStore\xca\x02\rIdentityStore\xe2\x02\x19IdentityStore\\GPBMetadata\xea\x02\rIdentityStoreb\x06proto3"

var (
	file_identity_store_identity_store_proto_rawDescOnce sync.Once
	file_identity_store_identity_store_proto_rawDescData []byte
)

func file_identity_store_identity_store_proto_rawDescGZIP() []byte {
	file_identity_store_identity_store_proto_rawDescOnce.Do(func() {
		file_identity_store_identity_store_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_identity_store_identity_store_proto_rawDesc), len(file_identity_store_identity_store_proto_rawDesc)))
	})
	return file_identity_store_identity_store_proto_rawDescData
}

var file_identity_store_identity_store_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_identity_store_identity_store_proto_goTypes = []any{
	(*IdentityStoreSyncConfig)(nil),     // 0: identity_store.IdentityStoreSyncConfig
	(*IdentityStoreSyncResult)(nil),     // 1: identity_store.IdentityStoreSyncResult
	(*MetaData)(nil),                    // 2: identity_store.MetaData
	(*config.ConfigMap)(nil),            // 3: util.config.ConfigMap
	(*error1.ErrorResult)(nil),          // 4: util.error.ErrorResult
	(*emptypb.Empty)(nil),               // 5: google.protobuf.Empty
	(*version.CliBuildInformation)(nil), // 6: util.version.CliBuildInformation
}
var file_identity_store_identity_store_proto_depIdxs = []int32{
	3, // 0: identity_store.IdentityStoreSyncConfig.config_map:type_name -> util.config.ConfigMap
	4, // 1: identity_store.IdentityStoreSyncResult.error:type_name -> util.error.ErrorResult
	5, // 2: identity_store.IdentityStoreSyncService.CliVersionInformation:input_type -> google.protobuf.Empty
	0, // 3: identity_store.IdentityStoreSyncService.SyncIdentityStore:input_type -> identity_store.IdentityStoreSyncConfig
	3, // 4: identity_store.IdentityStoreSyncService.GetIdentityStoreMetaData:input_type -> util.config.ConfigMap
	6, // 5: identity_store.IdentityStoreSyncService.CliVersionInformation:output_type -> util.version.CliBuildInformation
	1, // 6: identity_store.IdentityStoreSyncService.SyncIdentityStore:output_type -> identity_store.IdentityStoreSyncResult
	2, // 7: identity_store.IdentityStoreSyncService.GetIdentityStoreMetaData:output_type -> identity_store.MetaData
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_identity_store_identity_store_proto_init() }
func file_identity_store_identity_store_proto_init() {
	if File_identity_store_identity_store_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_identity_store_identity_store_proto_rawDesc), len(file_identity_store_identity_store_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_identity_store_identity_store_proto_goTypes,
		DependencyIndexes: file_identity_store_identity_store_proto_depIdxs,
		MessageInfos:      file_identity_store_identity_store_proto_msgTypes,
	}.Build()
	File_identity_store_identity_store_proto = out.File
	file_identity_store_identity_store_proto_goTypes = nil
	file_identity_store_identity_store_proto_depIdxs = nil
}
