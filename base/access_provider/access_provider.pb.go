// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: access_provider/access_provider.proto

package access_provider

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

// AccessSyncToTarget contains all necessary configuration parameters to export Data from Raito into DS
type AccessSyncToTarget struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConfigMap *config.ConfigMap `protobuf:"bytes,1,opt,name=config_map,json=configMap,proto3" json:"config_map,omitempty"`
	// SourceFile points to the file containing the access controls that need to be pushed to the data source.
	SourceFile string `protobuf:"bytes,2,opt,name=source_file,json=sourceFile,proto3" json:"source_file,omitempty"`
	// FeedbackTargetFile points to the file where the plugin needs to export the access controls feedback to.
	FeedbackTargetFile string `protobuf:"bytes,3,opt,name=feedback_target_file,json=feedbackTargetFile,proto3" json:"feedback_target_file,omitempty"`
	Prefix             string `protobuf:"bytes,4,opt,name=prefix,proto3" json:"prefix,omitempty"`
	Test               string `protobuf:"bytes,5,opt,name=test,proto3" json:"test,omitempty"`
}

func (x *AccessSyncToTarget) Reset() {
	*x = AccessSyncToTarget{}
	if protoimpl.UnsafeEnabled {
		mi := &file_access_provider_access_provider_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessSyncToTarget) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessSyncToTarget) ProtoMessage() {}

func (x *AccessSyncToTarget) ProtoReflect() protoreflect.Message {
	mi := &file_access_provider_access_provider_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessSyncToTarget.ProtoReflect.Descriptor instead.
func (*AccessSyncToTarget) Descriptor() ([]byte, []int) {
	return file_access_provider_access_provider_proto_rawDescGZIP(), []int{0}
}

func (x *AccessSyncToTarget) GetConfigMap() *config.ConfigMap {
	if x != nil {
		return x.ConfigMap
	}
	return nil
}

func (x *AccessSyncToTarget) GetSourceFile() string {
	if x != nil {
		return x.SourceFile
	}
	return ""
}

func (x *AccessSyncToTarget) GetFeedbackTargetFile() string {
	if x != nil {
		return x.FeedbackTargetFile
	}
	return ""
}

func (x *AccessSyncToTarget) GetPrefix() string {
	if x != nil {
		return x.Prefix
	}
	return ""
}

func (x *AccessSyncToTarget) GetTest() string {
	if x != nil {
		return x.Test
	}
	return ""
}

// AccessSyncFromTarget contains all necessary configuration parameters to import Data from Raito into DS
type AccessSyncFromTarget struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConfigMap *config.ConfigMap `protobuf:"bytes,1,opt,name=config_map,json=configMap,proto3" json:"config_map,omitempty"`
	// TargetFile points to the file where the plugin needs to export the access control naming.
	TargetFile            string   `protobuf:"bytes,2,opt,name=target_file,json=targetFile,proto3" json:"target_file,omitempty"`
	Prefix                string   `protobuf:"bytes,3,opt,name=prefix,proto3" json:"prefix,omitempty"`
	LockAllWho            bool     `protobuf:"varint,4,opt,name=lock_all_who,json=lockAllWho,proto3" json:"lock_all_who,omitempty"`
	LockAllWhat           bool     `protobuf:"varint,5,opt,name=lock_all_what,json=lockAllWhat,proto3" json:"lock_all_what,omitempty"`
	LockAllNames          bool     `protobuf:"varint,6,opt,name=lock_all_names,json=lockAllNames,proto3" json:"lock_all_names,omitempty"`
	LockAllDelete         bool     `protobuf:"varint,7,opt,name=lock_all_delete,json=lockAllDelete,proto3" json:"lock_all_delete,omitempty"`
	LockAllInheritance    bool     `protobuf:"varint,8,opt,name=lock_all_inheritance,json=lockAllInheritance,proto3" json:"lock_all_inheritance,omitempty"`
	MakeNotInternalizable []string `protobuf:"bytes,9,rep,name=make_not_internalizable,json=makeNotInternalizable,proto3" json:"make_not_internalizable,omitempty"`
	LockAllOwners         bool     `protobuf:"varint,10,opt,name=lock_all_owners,json=lockAllOwners,proto3" json:"lock_all_owners,omitempty"`
}

func (x *AccessSyncFromTarget) Reset() {
	*x = AccessSyncFromTarget{}
	if protoimpl.UnsafeEnabled {
		mi := &file_access_provider_access_provider_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessSyncFromTarget) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessSyncFromTarget) ProtoMessage() {}

func (x *AccessSyncFromTarget) ProtoReflect() protoreflect.Message {
	mi := &file_access_provider_access_provider_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessSyncFromTarget.ProtoReflect.Descriptor instead.
func (*AccessSyncFromTarget) Descriptor() ([]byte, []int) {
	return file_access_provider_access_provider_proto_rawDescGZIP(), []int{1}
}

func (x *AccessSyncFromTarget) GetConfigMap() *config.ConfigMap {
	if x != nil {
		return x.ConfigMap
	}
	return nil
}

func (x *AccessSyncFromTarget) GetTargetFile() string {
	if x != nil {
		return x.TargetFile
	}
	return ""
}

func (x *AccessSyncFromTarget) GetPrefix() string {
	if x != nil {
		return x.Prefix
	}
	return ""
}

func (x *AccessSyncFromTarget) GetLockAllWho() bool {
	if x != nil {
		return x.LockAllWho
	}
	return false
}

func (x *AccessSyncFromTarget) GetLockAllWhat() bool {
	if x != nil {
		return x.LockAllWhat
	}
	return false
}

func (x *AccessSyncFromTarget) GetLockAllNames() bool {
	if x != nil {
		return x.LockAllNames
	}
	return false
}

func (x *AccessSyncFromTarget) GetLockAllDelete() bool {
	if x != nil {
		return x.LockAllDelete
	}
	return false
}

func (x *AccessSyncFromTarget) GetLockAllInheritance() bool {
	if x != nil {
		return x.LockAllInheritance
	}
	return false
}

func (x *AccessSyncFromTarget) GetMakeNotInternalizable() []string {
	if x != nil {
		return x.MakeNotInternalizable
	}
	return nil
}

func (x *AccessSyncFromTarget) GetLockAllOwners() bool {
	if x != nil {
		return x.LockAllOwners
	}
	return false
}

// AccessSyncResult represents the result from the data access sync process.
// A potential error is also modeled in here so specific errors remain intact when passed over RPC.
type AccessSyncResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Deprecated: Marked as deprecated in access_provider/access_provider.proto.
	Error               *error1.ErrorResult `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	AccessProviderCount int32               `protobuf:"varint,2,opt,name=access_provider_count,json=accessProviderCount,proto3" json:"access_provider_count,omitempty"`
}

func (x *AccessSyncResult) Reset() {
	*x = AccessSyncResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_access_provider_access_provider_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessSyncResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessSyncResult) ProtoMessage() {}

func (x *AccessSyncResult) ProtoReflect() protoreflect.Message {
	mi := &file_access_provider_access_provider_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessSyncResult.ProtoReflect.Descriptor instead.
func (*AccessSyncResult) Descriptor() ([]byte, []int) {
	return file_access_provider_access_provider_proto_rawDescGZIP(), []int{2}
}

// Deprecated: Marked as deprecated in access_provider/access_provider.proto.
func (x *AccessSyncResult) GetError() *error1.ErrorResult {
	if x != nil {
		return x.Error
	}
	return nil
}

func (x *AccessSyncResult) GetAccessProviderCount() int32 {
	if x != nil {
		return x.AccessProviderCount
	}
	return 0
}

// AccessSyncConfig gives us information on how the CLI can sync access providers
type AccessSyncConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// SupportPartialSync if true, syncing only out of sync access providers is allowed
	SupportPartialSync bool `protobuf:"varint,1,opt,name=support_partial_sync,json=supportPartialSync,proto3" json:"support_partial_sync,omitempty"`
}

func (x *AccessSyncConfig) Reset() {
	*x = AccessSyncConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_access_provider_access_provider_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessSyncConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessSyncConfig) ProtoMessage() {}

func (x *AccessSyncConfig) ProtoReflect() protoreflect.Message {
	mi := &file_access_provider_access_provider_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessSyncConfig.ProtoReflect.Descriptor instead.
func (*AccessSyncConfig) Descriptor() ([]byte, []int) {
	return file_access_provider_access_provider_proto_rawDescGZIP(), []int{3}
}

func (x *AccessSyncConfig) GetSupportPartialSync() bool {
	if x != nil {
		return x.SupportPartialSync
	}
	return false
}

var File_access_provider_access_provider_proto protoreflect.FileDescriptor

var file_access_provider_access_provider_proto_rawDesc = []byte{
	0x0a, 0x25, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x2f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f,
	0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x18, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x16, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2f, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xca, 0x01, 0x0a, 0x12, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x53, 0x79,
	0x6e, 0x63, 0x54, 0x6f, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x35, 0x0a, 0x0a, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x5f, 0x6d, 0x61, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16,
	0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x4d, 0x61, 0x70, 0x52, 0x09, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4d, 0x61,
	0x70, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x66, 0x69, 0x6c, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x46, 0x69,
	0x6c, 0x65, 0x12, 0x30, 0x0a, 0x14, 0x66, 0x65, 0x65, 0x64, 0x62, 0x61, 0x63, 0x6b, 0x5f, 0x74,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x12, 0x66, 0x65, 0x65, 0x64, 0x62, 0x61, 0x63, 0x6b, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x46, 0x69, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x65, 0x73, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x73, 0x74,
	0x22, 0xac, 0x03, 0x0a, 0x14, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x53, 0x79, 0x6e, 0x63, 0x46,
	0x72, 0x6f, 0x6d, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x35, 0x0a, 0x0a, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x5f, 0x6d, 0x61, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x75, 0x74, 0x69, 0x6c, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x4d, 0x61, 0x70, 0x52, 0x09, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4d, 0x61, 0x70,
	0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x46, 0x69, 0x6c,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x12, 0x20, 0x0a, 0x0c, 0x6c, 0x6f, 0x63,
	0x6b, 0x5f, 0x61, 0x6c, 0x6c, 0x5f, 0x77, 0x68, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0a, 0x6c, 0x6f, 0x63, 0x6b, 0x41, 0x6c, 0x6c, 0x57, 0x68, 0x6f, 0x12, 0x22, 0x0a, 0x0d, 0x6c,
	0x6f, 0x63, 0x6b, 0x5f, 0x61, 0x6c, 0x6c, 0x5f, 0x77, 0x68, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x0b, 0x6c, 0x6f, 0x63, 0x6b, 0x41, 0x6c, 0x6c, 0x57, 0x68, 0x61, 0x74, 0x12,
	0x24, 0x0a, 0x0e, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x61, 0x6c, 0x6c, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x6c, 0x6f, 0x63, 0x6b, 0x41, 0x6c, 0x6c,
	0x4e, 0x61, 0x6d, 0x65, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x61, 0x6c,
	0x6c, 0x5f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d,
	0x6c, 0x6f, 0x63, 0x6b, 0x41, 0x6c, 0x6c, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x30, 0x0a,
	0x14, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x61, 0x6c, 0x6c, 0x5f, 0x69, 0x6e, 0x68, 0x65, 0x72, 0x69,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x12, 0x6c, 0x6f, 0x63,
	0x6b, 0x41, 0x6c, 0x6c, 0x49, 0x6e, 0x68, 0x65, 0x72, 0x69, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x12,
	0x36, 0x0a, 0x17, 0x6d, 0x61, 0x6b, 0x65, 0x5f, 0x6e, 0x6f, 0x74, 0x5f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x09, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x15, 0x6d, 0x61, 0x6b, 0x65, 0x4e, 0x6f, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x69, 0x7a, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x26, 0x0a, 0x0f, 0x6c, 0x6f, 0x63, 0x6b, 0x5f,
	0x61, 0x6c, 0x6c, 0x5f, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0d, 0x6c, 0x6f, 0x63, 0x6b, 0x41, 0x6c, 0x6c, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x73, 0x22,
	0x79, 0x0a, 0x10, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x12, 0x31, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42, 0x02, 0x18, 0x01, 0x52,
	0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x32, 0x0a, 0x15, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x13, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x72, 0x6f,
	0x76, 0x69, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x50, 0x0a, 0x10, 0x41, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x53, 0x79, 0x6e, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x30,
	0x0a, 0x14, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x70, 0x61, 0x72, 0x74, 0x69, 0x61,
	0x6c, 0x5f, 0x73, 0x79, 0x6e, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x12, 0x73, 0x75,
	0x70, 0x70, 0x6f, 0x72, 0x74, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x53, 0x79, 0x6e, 0x63,
	0x4a, 0x04, 0x08, 0x02, 0x10, 0x03, 0x4a, 0x04, 0x08, 0x03, 0x10, 0x04, 0x32, 0xec, 0x02, 0x0a,
	0x19, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x53,
	0x79, 0x6e, 0x63, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x52, 0x0a, 0x15, 0x43, 0x6c,
	0x69, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x21, 0x2e, 0x75, 0x74,
	0x69, 0x6c, 0x2e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6c, 0x69, 0x42, 0x75,
	0x69, 0x6c, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x5a,
	0x0a, 0x0e, 0x53, 0x79, 0x6e, 0x63, 0x46, 0x72, 0x6f, 0x6d, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x12, 0x25, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64,
	0x65, 0x72, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x53, 0x79, 0x6e, 0x63, 0x46, 0x72, 0x6f,
	0x6d, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x1a, 0x21, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x56, 0x0a, 0x0c, 0x53, 0x79,
	0x6e, 0x63, 0x54, 0x6f, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x23, 0x2e, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x41, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x53, 0x79, 0x6e, 0x63, 0x54, 0x6f, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x1a,
	0x21, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x12, 0x47, 0x0a, 0x0a, 0x53, 0x79, 0x6e, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x21, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x53, 0x79, 0x6e, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0xb0, 0x01, 0x0a, 0x13,
	0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x72, 0x42, 0x13, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x61, 0x69, 0x74, 0x6f, 0x2d, 0x69, 0x6f, 0x2f,
	0x63, 0x6c, 0x69, 0x2f, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f,
	0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0xa2, 0x02, 0x03, 0x41, 0x58, 0x58, 0xaa, 0x02,
	0x0e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0xca,
	0x02, 0x0e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0xe2, 0x02, 0x1a, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0e,
	0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_access_provider_access_provider_proto_rawDescOnce sync.Once
	file_access_provider_access_provider_proto_rawDescData = file_access_provider_access_provider_proto_rawDesc
)

func file_access_provider_access_provider_proto_rawDescGZIP() []byte {
	file_access_provider_access_provider_proto_rawDescOnce.Do(func() {
		file_access_provider_access_provider_proto_rawDescData = protoimpl.X.CompressGZIP(file_access_provider_access_provider_proto_rawDescData)
	})
	return file_access_provider_access_provider_proto_rawDescData
}

var file_access_provider_access_provider_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_access_provider_access_provider_proto_goTypes = []interface{}{
	(*AccessSyncToTarget)(nil),          // 0: access_provider.AccessSyncToTarget
	(*AccessSyncFromTarget)(nil),        // 1: access_provider.AccessSyncFromTarget
	(*AccessSyncResult)(nil),            // 2: access_provider.AccessSyncResult
	(*AccessSyncConfig)(nil),            // 3: access_provider.AccessSyncConfig
	(*config.ConfigMap)(nil),            // 4: util.config.ConfigMap
	(*error1.ErrorResult)(nil),          // 5: util.error.ErrorResult
	(*emptypb.Empty)(nil),               // 6: google.protobuf.Empty
	(*version.CliBuildInformation)(nil), // 7: util.version.CliBuildInformation
}
var file_access_provider_access_provider_proto_depIdxs = []int32{
	4, // 0: access_provider.AccessSyncToTarget.config_map:type_name -> util.config.ConfigMap
	4, // 1: access_provider.AccessSyncFromTarget.config_map:type_name -> util.config.ConfigMap
	5, // 2: access_provider.AccessSyncResult.error:type_name -> util.error.ErrorResult
	6, // 3: access_provider.AccessProviderSyncService.CliVersionInformation:input_type -> google.protobuf.Empty
	1, // 4: access_provider.AccessProviderSyncService.SyncFromTarget:input_type -> access_provider.AccessSyncFromTarget
	0, // 5: access_provider.AccessProviderSyncService.SyncToTarget:input_type -> access_provider.AccessSyncToTarget
	6, // 6: access_provider.AccessProviderSyncService.SyncConfig:input_type -> google.protobuf.Empty
	7, // 7: access_provider.AccessProviderSyncService.CliVersionInformation:output_type -> util.version.CliBuildInformation
	2, // 8: access_provider.AccessProviderSyncService.SyncFromTarget:output_type -> access_provider.AccessSyncResult
	2, // 9: access_provider.AccessProviderSyncService.SyncToTarget:output_type -> access_provider.AccessSyncResult
	3, // 10: access_provider.AccessProviderSyncService.SyncConfig:output_type -> access_provider.AccessSyncConfig
	7, // [7:11] is the sub-list for method output_type
	3, // [3:7] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_access_provider_access_provider_proto_init() }
func file_access_provider_access_provider_proto_init() {
	if File_access_provider_access_provider_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_access_provider_access_provider_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessSyncToTarget); i {
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
		file_access_provider_access_provider_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessSyncFromTarget); i {
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
		file_access_provider_access_provider_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessSyncResult); i {
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
		file_access_provider_access_provider_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessSyncConfig); i {
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
			RawDescriptor: file_access_provider_access_provider_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_access_provider_access_provider_proto_goTypes,
		DependencyIndexes: file_access_provider_access_provider_proto_depIdxs,
		MessageInfos:      file_access_provider_access_provider_proto_msgTypes,
	}.Build()
	File_access_provider_access_provider_proto = out.File
	file_access_provider_access_provider_proto_rawDesc = nil
	file_access_provider_access_provider_proto_goTypes = nil
	file_access_provider_access_provider_proto_depIdxs = nil
}
