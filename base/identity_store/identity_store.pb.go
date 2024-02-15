// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
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
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConfigMap *config.ConfigMap `protobuf:"bytes,1,opt,name=config_map,json=configMap,proto3" json:"config_map,omitempty"`
	UserFile  string            `protobuf:"bytes,2,opt,name=user_file,json=userFile,proto3" json:"user_file,omitempty"`
	GroupFile string            `protobuf:"bytes,3,opt,name=group_file,json=groupFile,proto3" json:"group_file,omitempty"`
}

func (x *IdentityStoreSyncConfig) Reset() {
	*x = IdentityStoreSyncConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_identity_store_identity_store_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IdentityStoreSyncConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IdentityStoreSyncConfig) ProtoMessage() {}

func (x *IdentityStoreSyncConfig) ProtoReflect() protoreflect.Message {
	mi := &file_identity_store_identity_store_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
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
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Deprecated: Marked as deprecated in identity_store/identity_store.proto.
	Error      *error1.ErrorResult `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	UserCount  int32               `protobuf:"varint,2,opt,name=user_count,json=userCount,proto3" json:"user_count,omitempty"`
	GroupCount int32               `protobuf:"varint,3,opt,name=group_count,json=groupCount,proto3" json:"group_count,omitempty"`
}

func (x *IdentityStoreSyncResult) Reset() {
	*x = IdentityStoreSyncResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_identity_store_identity_store_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IdentityStoreSyncResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IdentityStoreSyncResult) ProtoMessage() {}

func (x *IdentityStoreSyncResult) ProtoReflect() protoreflect.Message {
	mi := &file_identity_store_identity_store_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
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
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type        string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Icon        string `protobuf:"bytes,2,opt,name=icon,proto3" json:"icon,omitempty"`
	CanBeLinked bool   `protobuf:"varint,3,opt,name=can_be_linked,json=canBeLinked,proto3" json:"can_be_linked,omitempty"`
	CanBeMaster bool   `protobuf:"varint,4,opt,name=can_be_master,json=canBeMaster,proto3" json:"can_be_master,omitempty"`
}

func (x *MetaData) Reset() {
	*x = MetaData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_identity_store_identity_store_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetaData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetaData) ProtoMessage() {}

func (x *MetaData) ProtoReflect() protoreflect.Message {
	mi := &file_identity_store_identity_store_proto_msgTypes[2]
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

var file_identity_store_identity_store_proto_rawDesc = []byte{
	0x0a, 0x23, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x18, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x75, 0x74,
	0x69, 0x6c, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x8c, 0x01, 0x0a, 0x17, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x53, 0x74, 0x6f,
	0x72, 0x65, 0x53, 0x79, 0x6e, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x35, 0x0a, 0x0a,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x6d, 0x61, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x16, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4d, 0x61, 0x70, 0x52, 0x09, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x4d, 0x61, 0x70, 0x12, 0x1b, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x66, 0x69, 0x6c, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x46, 0x69, 0x6c, 0x65,
	0x12, 0x1d, 0x0a, 0x0a, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x46, 0x69, 0x6c, 0x65, 0x22,
	0x8c, 0x01, 0x0a, 0x17, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x53, 0x74, 0x6f, 0x72,
	0x65, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x31, 0x0a, 0x05, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x75, 0x74, 0x69,
	0x6c, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x42, 0x02, 0x18, 0x01, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x1d,
	0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x09, 0x75, 0x73, 0x65, 0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1f, 0x0a,
	0x0b, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0a, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x7a,
	0x0a, 0x08, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x69, 0x63, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x63,
	0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x0d, 0x63, 0x61, 0x6e, 0x5f, 0x62, 0x65, 0x5f, 0x6c, 0x69, 0x6e,
	0x6b, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x63, 0x61, 0x6e, 0x42, 0x65,
	0x4c, 0x69, 0x6e, 0x6b, 0x65, 0x64, 0x12, 0x22, 0x0a, 0x0d, 0x63, 0x61, 0x6e, 0x5f, 0x62, 0x65,
	0x5f, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x63,
	0x61, 0x6e, 0x42, 0x65, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x32, 0xa3, 0x02, 0x0a, 0x18, 0x49,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x53, 0x79, 0x6e, 0x63,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x52, 0x0a, 0x15, 0x43, 0x6c, 0x69, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x21, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6c, 0x69, 0x42, 0x75, 0x69, 0x6c, 0x64,
	0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x65, 0x0a, 0x11, 0x53,
	0x79, 0x6e, 0x63, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x53, 0x74, 0x6f, 0x72, 0x65,
	0x12, 0x27, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x53,
	0x79, 0x6e, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a, 0x27, 0x2e, 0x69, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x12, 0x4c, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x12, 0x16,
	0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x4d, 0x61, 0x70, 0x1a, 0x18, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61,
	0x42, 0xa9, 0x01, 0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x42, 0x12, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2b, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x61, 0x69, 0x74, 0x6f, 0x2d,
	0x69, 0x6f, 0x2f, 0x63, 0x6c, 0x69, 0x2f, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x69, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0xa2, 0x02, 0x03, 0x49, 0x58, 0x58,
	0xaa, 0x02, 0x0d, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x53, 0x74, 0x6f, 0x72, 0x65,
	0xca, 0x02, 0x0d, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x53, 0x74, 0x6f, 0x72, 0x65,
	0xe2, 0x02, 0x19, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x53, 0x74, 0x6f, 0x72, 0x65,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0d, 0x49,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_identity_store_identity_store_proto_rawDescOnce sync.Once
	file_identity_store_identity_store_proto_rawDescData = file_identity_store_identity_store_proto_rawDesc
)

func file_identity_store_identity_store_proto_rawDescGZIP() []byte {
	file_identity_store_identity_store_proto_rawDescOnce.Do(func() {
		file_identity_store_identity_store_proto_rawDescData = protoimpl.X.CompressGZIP(file_identity_store_identity_store_proto_rawDescData)
	})
	return file_identity_store_identity_store_proto_rawDescData
}

var file_identity_store_identity_store_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_identity_store_identity_store_proto_goTypes = []interface{}{
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
	if !protoimpl.UnsafeEnabled {
		file_identity_store_identity_store_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IdentityStoreSyncConfig); i {
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
		file_identity_store_identity_store_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IdentityStoreSyncResult); i {
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
		file_identity_store_identity_store_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_identity_store_identity_store_proto_rawDesc,
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
	file_identity_store_identity_store_proto_rawDesc = nil
	file_identity_store_identity_store_proto_goTypes = nil
	file_identity_store_identity_store_proto_depIdxs = nil
}
