// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: data_source/data_source.proto

package data_source

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

// DataSourceSyncConfig represents the configuration that is passed from the CLI to the DataAccessSyncer plugin interface.
// It contains all the necessary configuration parameters for the plugin to function.
type DataSourceSyncConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConfigMap    *config.ConfigMap `protobuf:"bytes,1,opt,name=config_map,json=configMap,proto3" json:"config_map,omitempty"`
	TargetFile   string            `protobuf:"bytes,2,opt,name=target_file,json=targetFile,proto3" json:"target_file,omitempty"`
	DataSourceId string            `protobuf:"bytes,3,opt,name=data_source_id,json=dataSourceId,proto3" json:"data_source_id,omitempty"`
}

func (x *DataSourceSyncConfig) Reset() {
	*x = DataSourceSyncConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_source_data_source_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataSourceSyncConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataSourceSyncConfig) ProtoMessage() {}

func (x *DataSourceSyncConfig) ProtoReflect() protoreflect.Message {
	mi := &file_data_source_data_source_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataSourceSyncConfig.ProtoReflect.Descriptor instead.
func (*DataSourceSyncConfig) Descriptor() ([]byte, []int) {
	return file_data_source_data_source_proto_rawDescGZIP(), []int{0}
}

func (x *DataSourceSyncConfig) GetConfigMap() *config.ConfigMap {
	if x != nil {
		return x.ConfigMap
	}
	return nil
}

func (x *DataSourceSyncConfig) GetTargetFile() string {
	if x != nil {
		return x.TargetFile
	}
	return ""
}

func (x *DataSourceSyncConfig) GetDataSourceId() string {
	if x != nil {
		return x.DataSourceId
	}
	return ""
}

// DataSourceSyncResult represents the result from the data source sync process.
// A potential error is also modeled in here so specific errors remain intact when passed over RPC.
type DataSourceSyncResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Deprecated: Do not use.
	Error       *error1.ErrorResult `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	DataObjects int32               `protobuf:"varint,2,opt,name=data_objects,json=dataObjects,proto3" json:"data_objects,omitempty"`
}

func (x *DataSourceSyncResult) Reset() {
	*x = DataSourceSyncResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_source_data_source_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataSourceSyncResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataSourceSyncResult) ProtoMessage() {}

func (x *DataSourceSyncResult) ProtoReflect() protoreflect.Message {
	mi := &file_data_source_data_source_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataSourceSyncResult.ProtoReflect.Descriptor instead.
func (*DataSourceSyncResult) Descriptor() ([]byte, []int) {
	return file_data_source_data_source_proto_rawDescGZIP(), []int{1}
}

// Deprecated: Do not use.
func (x *DataSourceSyncResult) GetError() *error1.ErrorResult {
	if x != nil {
		return x.Error
	}
	return nil
}

func (x *DataSourceSyncResult) GetDataObjects() int32 {
	if x != nil {
		return x.DataObjects
	}
	return 0
}

// NoLinting ToSupport GQL schema (temporarily)
type MetaData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DataObjectTypes   []*DataObjectType `protobuf:"bytes,1,rep,name=dataObjectTypes,proto3" json:"dataObjectTypes,omitempty"`
	SupportedFeatures []string          `protobuf:"bytes,2,rep,name=supportedFeatures,proto3" json:"supportedFeatures,omitempty"`
	Type              string            `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	Icon              string            `protobuf:"bytes,4,opt,name=icon,proto3" json:"icon,omitempty"`
	UsageMetaInfo     *UsageMetaInput   `protobuf:"bytes,5,opt,name=usageMetaInfo,proto3" json:"usageMetaInfo,omitempty"`
}

func (x *MetaData) Reset() {
	*x = MetaData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_source_data_source_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetaData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetaData) ProtoMessage() {}

func (x *MetaData) ProtoReflect() protoreflect.Message {
	mi := &file_data_source_data_source_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
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
	return file_data_source_data_source_proto_rawDescGZIP(), []int{2}
}

func (x *MetaData) GetDataObjectTypes() []*DataObjectType {
	if x != nil {
		return x.DataObjectTypes
	}
	return nil
}

func (x *MetaData) GetSupportedFeatures() []string {
	if x != nil {
		return x.SupportedFeatures
	}
	return nil
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

func (x *MetaData) GetUsageMetaInfo() *UsageMetaInput {
	if x != nil {
		return x.UsageMetaInfo
	}
	return nil
}

type DataObjectType struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string                      `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Type        string                      `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Label       string                      `protobuf:"bytes,3,opt,name=label,proto3" json:"label,omitempty"`
	Icon        string                      `protobuf:"bytes,4,opt,name=icon,proto3" json:"icon,omitempty"`
	Permissions []*DataObjectTypePermission `protobuf:"bytes,5,rep,name=permissions,proto3" json:"permissions,omitempty"`
	Actions     []*DataObjectTypeAction     `protobuf:"bytes,6,rep,name=actions,proto3" json:"actions,omitempty"`
	Children    []string                    `protobuf:"bytes,7,rep,name=children,proto3" json:"children,omitempty"`
}

func (x *DataObjectType) Reset() {
	*x = DataObjectType{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_source_data_source_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataObjectType) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataObjectType) ProtoMessage() {}

func (x *DataObjectType) ProtoReflect() protoreflect.Message {
	mi := &file_data_source_data_source_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataObjectType.ProtoReflect.Descriptor instead.
func (*DataObjectType) Descriptor() ([]byte, []int) {
	return file_data_source_data_source_proto_rawDescGZIP(), []int{3}
}

func (x *DataObjectType) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DataObjectType) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *DataObjectType) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

func (x *DataObjectType) GetIcon() string {
	if x != nil {
		return x.Icon
	}
	return ""
}

func (x *DataObjectType) GetPermissions() []*DataObjectTypePermission {
	if x != nil {
		return x.Permissions
	}
	return nil
}

func (x *DataObjectType) GetActions() []*DataObjectTypeAction {
	if x != nil {
		return x.Actions
	}
	return nil
}

func (x *DataObjectType) GetChildren() []string {
	if x != nil {
		return x.Children
	}
	return nil
}

// NoLinting ToSupport GQL schema (temporarily)
type DataObjectTypePermission struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Permission             string   `protobuf:"bytes,1,opt,name=permission,proto3" json:"permission,omitempty"`
	GlobalPermissions      []string `protobuf:"bytes,2,rep,name=globalPermissions,proto3" json:"globalPermissions,omitempty"`
	Description            string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Action                 []string `protobuf:"bytes,4,rep,name=action,proto3" json:"action,omitempty"`
	UsageGlobalPermissions []string `protobuf:"bytes,5,rep,name=usageGlobalPermissions,proto3" json:"usageGlobalPermissions,omitempty"`
	CannotBeGranted        bool     `protobuf:"varint,6,opt,name=cannotBeGranted,proto3" json:"cannotBeGranted,omitempty"`
}

func (x *DataObjectTypePermission) Reset() {
	*x = DataObjectTypePermission{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_source_data_source_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataObjectTypePermission) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataObjectTypePermission) ProtoMessage() {}

func (x *DataObjectTypePermission) ProtoReflect() protoreflect.Message {
	mi := &file_data_source_data_source_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataObjectTypePermission.ProtoReflect.Descriptor instead.
func (*DataObjectTypePermission) Descriptor() ([]byte, []int) {
	return file_data_source_data_source_proto_rawDescGZIP(), []int{4}
}

func (x *DataObjectTypePermission) GetPermission() string {
	if x != nil {
		return x.Permission
	}
	return ""
}

func (x *DataObjectTypePermission) GetGlobalPermissions() []string {
	if x != nil {
		return x.GlobalPermissions
	}
	return nil
}

func (x *DataObjectTypePermission) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *DataObjectTypePermission) GetAction() []string {
	if x != nil {
		return x.Action
	}
	return nil
}

func (x *DataObjectTypePermission) GetUsageGlobalPermissions() []string {
	if x != nil {
		return x.UsageGlobalPermissions
	}
	return nil
}

func (x *DataObjectTypePermission) GetCannotBeGranted() bool {
	if x != nil {
		return x.CannotBeGranted
	}
	return false
}

// NoLinting ToSupport GQL schema (temporarily)
type DataObjectTypeAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Action        string   `protobuf:"bytes,1,opt,name=action,proto3" json:"action,omitempty"`
	GlobalActions []string `protobuf:"bytes,2,rep,name=globalActions,proto3" json:"globalActions,omitempty"`
}

func (x *DataObjectTypeAction) Reset() {
	*x = DataObjectTypeAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_source_data_source_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataObjectTypeAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataObjectTypeAction) ProtoMessage() {}

func (x *DataObjectTypeAction) ProtoReflect() protoreflect.Message {
	mi := &file_data_source_data_source_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataObjectTypeAction.ProtoReflect.Descriptor instead.
func (*DataObjectTypeAction) Descriptor() ([]byte, []int) {
	return file_data_source_data_source_proto_rawDescGZIP(), []int{5}
}

func (x *DataObjectTypeAction) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *DataObjectTypeAction) GetGlobalActions() []string {
	if x != nil {
		return x.GlobalActions
	}
	return nil
}

type UsageMetaInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DefaultLevel string                  `protobuf:"bytes,1,opt,name=defaultLevel,proto3" json:"defaultLevel,omitempty"`
	Levels       []*UsageMetaInputDetail `protobuf:"bytes,2,rep,name=levels,proto3" json:"levels,omitempty"`
}

func (x *UsageMetaInput) Reset() {
	*x = UsageMetaInput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_source_data_source_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UsageMetaInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UsageMetaInput) ProtoMessage() {}

func (x *UsageMetaInput) ProtoReflect() protoreflect.Message {
	mi := &file_data_source_data_source_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UsageMetaInput.ProtoReflect.Descriptor instead.
func (*UsageMetaInput) Descriptor() ([]byte, []int) {
	return file_data_source_data_source_proto_rawDescGZIP(), []int{6}
}

func (x *UsageMetaInput) GetDefaultLevel() string {
	if x != nil {
		return x.DefaultLevel
	}
	return ""
}

func (x *UsageMetaInput) GetLevels() []*UsageMetaInputDetail {
	if x != nil {
		return x.Levels
	}
	return nil
}

type UsageMetaInputDetail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name            string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	DataObjectTypes []string `protobuf:"bytes,2,rep,name=dataObjectTypes,proto3" json:"dataObjectTypes,omitempty"`
}

func (x *UsageMetaInputDetail) Reset() {
	*x = UsageMetaInputDetail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_data_source_data_source_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UsageMetaInputDetail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UsageMetaInputDetail) ProtoMessage() {}

func (x *UsageMetaInputDetail) ProtoReflect() protoreflect.Message {
	mi := &file_data_source_data_source_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UsageMetaInputDetail.ProtoReflect.Descriptor instead.
func (*UsageMetaInputDetail) Descriptor() ([]byte, []int) {
	return file_data_source_data_source_proto_rawDescGZIP(), []int{7}
}

func (x *UsageMetaInputDetail) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UsageMetaInputDetail) GetDataObjectTypes() []string {
	if x != nil {
		return x.DataObjectTypes
	}
	return nil
}

var File_data_source_data_source_proto protoreflect.FileDescriptor

var file_data_source_data_source_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2f, 0x64, 0x61,
	0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0b, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x18, 0x75, 0x74, 0x69, 0x6c, 0x2f,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2f,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x75, 0x74, 0x69,
	0x6c, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x94, 0x01, 0x0a, 0x14, 0x44, 0x61, 0x74, 0x61,
	0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x79, 0x6e, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x12, 0x35, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x6d, 0x61, 0x70, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4d, 0x61, 0x70, 0x52, 0x09, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x4d, 0x61, 0x70, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x24, 0x0a, 0x0e, 0x64, 0x61, 0x74, 0x61,
	0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x64, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x22, 0x6c,
	0x0a, 0x14, 0x44, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x79, 0x6e, 0x63,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x31, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42, 0x02,
	0x18, 0x01, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x61, 0x74,
	0x61, 0x5f, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0b, 0x64, 0x61, 0x74, 0x61, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x22, 0xea, 0x01, 0x0a,
	0x08, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x12, 0x45, 0x0a, 0x0f, 0x64, 0x61, 0x74,
	0x61, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x54, 0x79, 0x70, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x2e, 0x44, 0x61, 0x74, 0x61, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x0f, 0x64, 0x61, 0x74, 0x61, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x54, 0x79, 0x70, 0x65, 0x73,
	0x12, 0x2c, 0x0a, 0x11, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x46, 0x65, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x11, 0x73, 0x75, 0x70,
	0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x12, 0x12,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x63, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x69, 0x63, 0x6f, 0x6e, 0x12, 0x41, 0x0a, 0x0d, 0x75, 0x73, 0x61, 0x67, 0x65, 0x4d,
	0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e,
	0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x55, 0x73, 0x61, 0x67,
	0x65, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x52, 0x0d, 0x75, 0x73, 0x61, 0x67,
	0x65, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x84, 0x02, 0x0a, 0x0e, 0x44, 0x61,
	0x74, 0x61, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x63,
	0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x63, 0x6f, 0x6e, 0x12, 0x47,
	0x0a, 0x0b, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x05, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x54, 0x79, 0x70, 0x65,
	0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x70, 0x65, 0x72, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x3b, 0x0a, 0x07, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x5f,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x4f, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x54, 0x79, 0x70, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e,
	0x18, 0x07, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x63, 0x68, 0x69, 0x6c, 0x64, 0x72, 0x65, 0x6e,
	0x22, 0x84, 0x02, 0x0a, 0x18, 0x44, 0x61, 0x74, 0x61, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1e, 0x0a,
	0x0a, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x2c, 0x0a,
	0x11, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x11, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c,
	0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a,
	0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x36, 0x0a, 0x16, 0x75, 0x73, 0x61, 0x67, 0x65, 0x47, 0x6c,
	0x6f, 0x62, 0x61, 0x6c, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x16, 0x75, 0x73, 0x61, 0x67, 0x65, 0x47, 0x6c, 0x6f, 0x62,
	0x61, 0x6c, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x28, 0x0a,
	0x0f, 0x63, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x42, 0x65, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x64,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0f, 0x63, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x42, 0x65,
	0x47, 0x72, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x22, 0x54, 0x0a, 0x14, 0x44, 0x61, 0x74, 0x61, 0x4f,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x54, 0x79, 0x70, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x16, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x24, 0x0a, 0x0d, 0x67, 0x6c, 0x6f, 0x62, 0x61,
	0x6c, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0d,
	0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x6f, 0x0a,
	0x0e, 0x55, 0x73, 0x61, 0x67, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x12,
	0x22, 0x0a, 0x0c, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x4c, 0x65,
	0x76, 0x65, 0x6c, 0x12, 0x39, 0x0a, 0x06, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x2e, 0x55, 0x73, 0x61, 0x67, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x70, 0x75, 0x74,
	0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x06, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x73, 0x22, 0x54,
	0x0a, 0x14, 0x55, 0x73, 0x61, 0x67, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x49, 0x6e, 0x70, 0x75, 0x74,
	0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x28, 0x0a, 0x0f, 0x64, 0x61,
	0x74, 0x61, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x54, 0x79, 0x70, 0x65, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x0f, 0x64, 0x61, 0x74, 0x61, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x73, 0x32, 0x8b, 0x02, 0x0a, 0x15, 0x44, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x53, 0x79, 0x6e, 0x63, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x52,
	0x0a, 0x15, 0x43, 0x6c, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f,
	0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a,
	0x21, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x43,
	0x6c, 0x69, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x56, 0x0a, 0x0e, 0x53, 0x79, 0x6e, 0x63, 0x44, 0x61, 0x74, 0x61, 0x53, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x12, 0x21, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x79, 0x6e,
	0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a, 0x21, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x46, 0x0a, 0x15, 0x47, 0x65,
	0x74, 0x44, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x44,
	0x61, 0x74, 0x61, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x15, 0x2e, 0x64, 0x61,
	0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61,
	0x74, 0x61, 0x42, 0x94, 0x01, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x5f,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x42, 0x0f, 0x44, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x61, 0x69, 0x74, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x63,
	0x6c, 0x69, 0x2f, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0xa2, 0x02, 0x03, 0x44, 0x58, 0x58, 0xaa, 0x02, 0x0a, 0x44, 0x61, 0x74, 0x61,
	0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0xca, 0x02, 0x0a, 0x44, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0xe2, 0x02, 0x16, 0x44, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0a, 0x44,
	0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_data_source_data_source_proto_rawDescOnce sync.Once
	file_data_source_data_source_proto_rawDescData = file_data_source_data_source_proto_rawDesc
)

func file_data_source_data_source_proto_rawDescGZIP() []byte {
	file_data_source_data_source_proto_rawDescOnce.Do(func() {
		file_data_source_data_source_proto_rawDescData = protoimpl.X.CompressGZIP(file_data_source_data_source_proto_rawDescData)
	})
	return file_data_source_data_source_proto_rawDescData
}

var file_data_source_data_source_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_data_source_data_source_proto_goTypes = []interface{}{
	(*DataSourceSyncConfig)(nil),        // 0: data_source.DataSourceSyncConfig
	(*DataSourceSyncResult)(nil),        // 1: data_source.DataSourceSyncResult
	(*MetaData)(nil),                    // 2: data_source.MetaData
	(*DataObjectType)(nil),              // 3: data_source.DataObjectType
	(*DataObjectTypePermission)(nil),    // 4: data_source.DataObjectTypePermission
	(*DataObjectTypeAction)(nil),        // 5: data_source.DataObjectTypeAction
	(*UsageMetaInput)(nil),              // 6: data_source.UsageMetaInput
	(*UsageMetaInputDetail)(nil),        // 7: data_source.UsageMetaInputDetail
	(*config.ConfigMap)(nil),            // 8: util.config.ConfigMap
	(*error1.ErrorResult)(nil),          // 9: util.error.ErrorResult
	(*emptypb.Empty)(nil),               // 10: google.protobuf.Empty
	(*version.CliBuildInformation)(nil), // 11: util.version.CliBuildInformation
}
var file_data_source_data_source_proto_depIdxs = []int32{
	8,  // 0: data_source.DataSourceSyncConfig.config_map:type_name -> util.config.ConfigMap
	9,  // 1: data_source.DataSourceSyncResult.error:type_name -> util.error.ErrorResult
	3,  // 2: data_source.MetaData.dataObjectTypes:type_name -> data_source.DataObjectType
	6,  // 3: data_source.MetaData.usageMetaInfo:type_name -> data_source.UsageMetaInput
	4,  // 4: data_source.DataObjectType.permissions:type_name -> data_source.DataObjectTypePermission
	5,  // 5: data_source.DataObjectType.actions:type_name -> data_source.DataObjectTypeAction
	7,  // 6: data_source.UsageMetaInput.levels:type_name -> data_source.UsageMetaInputDetail
	10, // 7: data_source.DataSourceSyncService.CliVersionInformation:input_type -> google.protobuf.Empty
	0,  // 8: data_source.DataSourceSyncService.SyncDataSource:input_type -> data_source.DataSourceSyncConfig
	10, // 9: data_source.DataSourceSyncService.GetDataSourceMetaData:input_type -> google.protobuf.Empty
	11, // 10: data_source.DataSourceSyncService.CliVersionInformation:output_type -> util.version.CliBuildInformation
	1,  // 11: data_source.DataSourceSyncService.SyncDataSource:output_type -> data_source.DataSourceSyncResult
	2,  // 12: data_source.DataSourceSyncService.GetDataSourceMetaData:output_type -> data_source.MetaData
	10, // [10:13] is the sub-list for method output_type
	7,  // [7:10] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_data_source_data_source_proto_init() }
func file_data_source_data_source_proto_init() {
	if File_data_source_data_source_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_data_source_data_source_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataSourceSyncConfig); i {
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
		file_data_source_data_source_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataSourceSyncResult); i {
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
		file_data_source_data_source_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MetaData); i {
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
		file_data_source_data_source_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataObjectType); i {
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
		file_data_source_data_source_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataObjectTypePermission); i {
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
		file_data_source_data_source_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataObjectTypeAction); i {
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
		file_data_source_data_source_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UsageMetaInput); i {
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
		file_data_source_data_source_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UsageMetaInputDetail); i {
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
			RawDescriptor: file_data_source_data_source_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_data_source_data_source_proto_goTypes,
		DependencyIndexes: file_data_source_data_source_proto_depIdxs,
		MessageInfos:      file_data_source_data_source_proto_msgTypes,
	}.Build()
	File_data_source_data_source_proto = out.File
	file_data_source_data_source_proto_rawDesc = nil
	file_data_source_data_source_proto_goTypes = nil
	file_data_source_data_source_proto_depIdxs = nil
}
