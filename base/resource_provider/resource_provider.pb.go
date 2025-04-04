// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: resource_provider/resource_provider.proto

package resource_provider

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

type UpdateResourceInput struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	ConfigMap       *config.ConfigMap      `protobuf:"bytes,1,opt,name=config_map,json=configMap,proto3" json:"config_map,omitempty"`
	Domain          string                 `protobuf:"bytes,11,opt,name=domain,proto3" json:"domain,omitempty"`
	DataSourceId    string                 `protobuf:"bytes,12,opt,name=data_source_id,json=dataSourceId,proto3" json:"data_source_id,omitempty"`
	IdentityStoreId string                 `protobuf:"bytes,13,opt,name=identity_store_id,json=identityStoreId,proto3" json:"identity_store_id,omitempty"`
	UrlOverride     *string                `protobuf:"bytes,14,opt,name=url_override,json=urlOverride,proto3,oneof" json:"url_override,omitempty"`
	Credentials     *ApiCredentials        `protobuf:"bytes,101,opt,name=credentials,proto3" json:"credentials,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *UpdateResourceInput) Reset() {
	*x = UpdateResourceInput{}
	mi := &file_resource_provider_resource_provider_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateResourceInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateResourceInput) ProtoMessage() {}

func (x *UpdateResourceInput) ProtoReflect() protoreflect.Message {
	mi := &file_resource_provider_resource_provider_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateResourceInput.ProtoReflect.Descriptor instead.
func (*UpdateResourceInput) Descriptor() ([]byte, []int) {
	return file_resource_provider_resource_provider_proto_rawDescGZIP(), []int{0}
}

func (x *UpdateResourceInput) GetConfigMap() *config.ConfigMap {
	if x != nil {
		return x.ConfigMap
	}
	return nil
}

func (x *UpdateResourceInput) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *UpdateResourceInput) GetDataSourceId() string {
	if x != nil {
		return x.DataSourceId
	}
	return ""
}

func (x *UpdateResourceInput) GetIdentityStoreId() string {
	if x != nil {
		return x.IdentityStoreId
	}
	return ""
}

func (x *UpdateResourceInput) GetUrlOverride() string {
	if x != nil && x.UrlOverride != nil {
		return *x.UrlOverride
	}
	return ""
}

func (x *UpdateResourceInput) GetCredentials() *ApiCredentials {
	if x != nil {
		return x.Credentials
	}
	return nil
}

type ApiCredentials struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Username      string                 `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ApiCredentials) Reset() {
	*x = ApiCredentials{}
	mi := &file_resource_provider_resource_provider_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ApiCredentials) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApiCredentials) ProtoMessage() {}

func (x *ApiCredentials) ProtoReflect() protoreflect.Message {
	mi := &file_resource_provider_resource_provider_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApiCredentials.ProtoReflect.Descriptor instead.
func (*ApiCredentials) Descriptor() ([]byte, []int) {
	return file_resource_provider_resource_provider_proto_rawDescGZIP(), []int{1}
}

func (x *ApiCredentials) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *ApiCredentials) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type UpdateResourceResult struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	AddedObjects   int32                  `protobuf:"varint,1,opt,name=added_objects,json=addedObjects,proto3" json:"added_objects,omitempty"`
	UpdatedObjects int32                  `protobuf:"varint,2,opt,name=updated_objects,json=updatedObjects,proto3" json:"updated_objects,omitempty"`
	DeletedObjects int32                  `protobuf:"varint,3,opt,name=deleted_objects,json=deletedObjects,proto3" json:"deleted_objects,omitempty"`
	Failures       int32                  `protobuf:"varint,4,opt,name=failures,proto3" json:"failures,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *UpdateResourceResult) Reset() {
	*x = UpdateResourceResult{}
	mi := &file_resource_provider_resource_provider_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateResourceResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateResourceResult) ProtoMessage() {}

func (x *UpdateResourceResult) ProtoReflect() protoreflect.Message {
	mi := &file_resource_provider_resource_provider_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateResourceResult.ProtoReflect.Descriptor instead.
func (*UpdateResourceResult) Descriptor() ([]byte, []int) {
	return file_resource_provider_resource_provider_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateResourceResult) GetAddedObjects() int32 {
	if x != nil {
		return x.AddedObjects
	}
	return 0
}

func (x *UpdateResourceResult) GetUpdatedObjects() int32 {
	if x != nil {
		return x.UpdatedObjects
	}
	return 0
}

func (x *UpdateResourceResult) GetDeletedObjects() int32 {
	if x != nil {
		return x.DeletedObjects
	}
	return 0
}

func (x *UpdateResourceResult) GetFailures() int32 {
	if x != nil {
		return x.Failures
	}
	return 0
}

var File_resource_provider_resource_provider_proto protoreflect.FileDescriptor

const file_resource_provider_resource_provider_proto_rawDesc = "" +
	"\n" +
	")resource_provider/resource_provider.proto\x12\x11resource_provider\x1a\x1bgoogle/protobuf/empty.proto\x1a\x18util/config/config.proto\x1a\x1autil/version/version.proto\"\xb4\x02\n" +
	"\x13UpdateResourceInput\x125\n" +
	"\n" +
	"config_map\x18\x01 \x01(\v2\x16.util.config.ConfigMapR\tconfigMap\x12\x16\n" +
	"\x06domain\x18\v \x01(\tR\x06domain\x12$\n" +
	"\x0edata_source_id\x18\f \x01(\tR\fdataSourceId\x12*\n" +
	"\x11identity_store_id\x18\r \x01(\tR\x0fidentityStoreId\x12&\n" +
	"\furl_override\x18\x0e \x01(\tH\x00R\vurlOverride\x88\x01\x01\x12C\n" +
	"\vcredentials\x18e \x01(\v2!.resource_provider.ApiCredentialsR\vcredentialsB\x0f\n" +
	"\r_url_override\"H\n" +
	"\x0eApiCredentials\x12\x1a\n" +
	"\busername\x18\x01 \x01(\tR\busername\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\"\xa9\x01\n" +
	"\x14UpdateResourceResult\x12#\n" +
	"\radded_objects\x18\x01 \x01(\x05R\faddedObjects\x12'\n" +
	"\x0fupdated_objects\x18\x02 \x01(\x05R\x0eupdatedObjects\x12'\n" +
	"\x0fdeleted_objects\x18\x03 \x01(\x05R\x0edeletedObjects\x12\x1a\n" +
	"\bfailures\x18\x04 \x01(\x05R\bfailures2\xd1\x01\n" +
	"\x17ResourceProviderService\x12R\n" +
	"\x15CliVersionInformation\x12\x16.google.protobuf.Empty\x1a!.util.version.CliBuildInformation\x12b\n" +
	"\x0fUpdateResources\x12&.resource_provider.UpdateResourceInput\x1a'.resource_provider.UpdateResourceResultB\xbe\x01\n" +
	"\x15com.resource_providerB\x15ResourceProviderProtoP\x01Z.github.com/raito-io/cli/base/resource_provider\xa2\x02\x03RXX\xaa\x02\x10ResourceProvider\xca\x02\x10ResourceProvider\xe2\x02\x1cResourceProvider\\GPBMetadata\xea\x02\x10ResourceProviderb\x06proto3"

var (
	file_resource_provider_resource_provider_proto_rawDescOnce sync.Once
	file_resource_provider_resource_provider_proto_rawDescData []byte
)

func file_resource_provider_resource_provider_proto_rawDescGZIP() []byte {
	file_resource_provider_resource_provider_proto_rawDescOnce.Do(func() {
		file_resource_provider_resource_provider_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_resource_provider_resource_provider_proto_rawDesc), len(file_resource_provider_resource_provider_proto_rawDesc)))
	})
	return file_resource_provider_resource_provider_proto_rawDescData
}

var file_resource_provider_resource_provider_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_resource_provider_resource_provider_proto_goTypes = []any{
	(*UpdateResourceInput)(nil),         // 0: resource_provider.UpdateResourceInput
	(*ApiCredentials)(nil),              // 1: resource_provider.ApiCredentials
	(*UpdateResourceResult)(nil),        // 2: resource_provider.UpdateResourceResult
	(*config.ConfigMap)(nil),            // 3: util.config.ConfigMap
	(*emptypb.Empty)(nil),               // 4: google.protobuf.Empty
	(*version.CliBuildInformation)(nil), // 5: util.version.CliBuildInformation
}
var file_resource_provider_resource_provider_proto_depIdxs = []int32{
	3, // 0: resource_provider.UpdateResourceInput.config_map:type_name -> util.config.ConfigMap
	1, // 1: resource_provider.UpdateResourceInput.credentials:type_name -> resource_provider.ApiCredentials
	4, // 2: resource_provider.ResourceProviderService.CliVersionInformation:input_type -> google.protobuf.Empty
	0, // 3: resource_provider.ResourceProviderService.UpdateResources:input_type -> resource_provider.UpdateResourceInput
	5, // 4: resource_provider.ResourceProviderService.CliVersionInformation:output_type -> util.version.CliBuildInformation
	2, // 5: resource_provider.ResourceProviderService.UpdateResources:output_type -> resource_provider.UpdateResourceResult
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_resource_provider_resource_provider_proto_init() }
func file_resource_provider_resource_provider_proto_init() {
	if File_resource_provider_resource_provider_proto != nil {
		return
	}
	file_resource_provider_resource_provider_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_resource_provider_resource_provider_proto_rawDesc), len(file_resource_provider_resource_provider_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_resource_provider_resource_provider_proto_goTypes,
		DependencyIndexes: file_resource_provider_resource_provider_proto_depIdxs,
		MessageInfos:      file_resource_provider_resource_provider_proto_msgTypes,
	}.Build()
	File_resource_provider_resource_provider_proto = out.File
	file_resource_provider_resource_provider_proto_goTypes = nil
	file_resource_provider_resource_provider_proto_depIdxs = nil
}
