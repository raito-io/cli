// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: util/error/error.proto

package error

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

type ErrorCode int32

const (
	ErrorCode_NO_ERROR                      ErrorCode = 0
	ErrorCode_UNKNOWN_ERROR                 ErrorCode = 1
	ErrorCode_BAD_INPUT_PARAMETER_ERROR     ErrorCode = 2
	ErrorCode_MISSING_INPUT_PARAMETER_ERROR ErrorCode = 3
	ErrorCode_SOURCE_CONNECTION_ERROR       ErrorCode = 4
)

// Enum value maps for ErrorCode.
var (
	ErrorCode_name = map[int32]string{
		0: "NO_ERROR",
		1: "UNKNOWN_ERROR",
		2: "BAD_INPUT_PARAMETER_ERROR",
		3: "MISSING_INPUT_PARAMETER_ERROR",
		4: "SOURCE_CONNECTION_ERROR",
	}
	ErrorCode_value = map[string]int32{
		"NO_ERROR":                      0,
		"UNKNOWN_ERROR":                 1,
		"BAD_INPUT_PARAMETER_ERROR":     2,
		"MISSING_INPUT_PARAMETER_ERROR": 3,
		"SOURCE_CONNECTION_ERROR":       4,
	}
)

func (x ErrorCode) Enum() *ErrorCode {
	p := new(ErrorCode)
	*p = x
	return p
}

func (x ErrorCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorCode) Descriptor() protoreflect.EnumDescriptor {
	return file_util_error_error_proto_enumTypes[0].Descriptor()
}

func (ErrorCode) Type() protoreflect.EnumType {
	return &file_util_error_error_proto_enumTypes[0]
}

func (x ErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorCode.Descriptor instead.
func (ErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_util_error_error_proto_rawDescGZIP(), []int{0}
}

type ErrorResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrorMessage string    `protobuf:"bytes,1,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	ErrorCode    ErrorCode `protobuf:"varint,2,opt,name=error_code,json=errorCode,proto3,enum=util.error.ErrorCode" json:"error_code,omitempty"`
}

func (x *ErrorResult) Reset() {
	*x = ErrorResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_util_error_error_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ErrorResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrorResult) ProtoMessage() {}

func (x *ErrorResult) ProtoReflect() protoreflect.Message {
	mi := &file_util_error_error_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrorResult.ProtoReflect.Descriptor instead.
func (*ErrorResult) Descriptor() ([]byte, []int) {
	return file_util_error_error_proto_rawDescGZIP(), []int{0}
}

func (x *ErrorResult) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

func (x *ErrorResult) GetErrorCode() ErrorCode {
	if x != nil {
		return x.ErrorCode
	}
	return ErrorCode_NO_ERROR
}

var File_util_error_error_proto protoreflect.FileDescriptor

var file_util_error_error_proto_rawDesc = []byte{
	0x0a, 0x16, 0x75, 0x74, 0x69, 0x6c, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2f, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x22, 0x68, 0x0a, 0x0b, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x34, 0x0a, 0x0a, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x75,
	0x74, 0x69, 0x6c, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43,
	0x6f, 0x64, 0x65, 0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x2a, 0x8b,
	0x01, 0x0a, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x0c, 0x0a, 0x08,
	0x4e, 0x4f, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x55, 0x4e,
	0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x01, 0x12, 0x1d, 0x0a,
	0x19, 0x42, 0x41, 0x44, 0x5f, 0x49, 0x4e, 0x50, 0x55, 0x54, 0x5f, 0x50, 0x41, 0x52, 0x41, 0x4d,
	0x45, 0x54, 0x45, 0x52, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x02, 0x12, 0x21, 0x0a, 0x1d,
	0x4d, 0x49, 0x53, 0x53, 0x49, 0x4e, 0x47, 0x5f, 0x49, 0x4e, 0x50, 0x55, 0x54, 0x5f, 0x50, 0x41,
	0x52, 0x41, 0x4d, 0x45, 0x54, 0x45, 0x52, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x03, 0x12,
	0x1b, 0x0a, 0x17, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43,
	0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x04, 0x42, 0x90, 0x01, 0x0a,
	0x0e, 0x63, 0x6f, 0x6d, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x2e, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x42,
	0x0a, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x27, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x61, 0x69, 0x74, 0x6f, 0x2d,
	0x69, 0x6f, 0x2f, 0x63, 0x6c, 0x69, 0x2f, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x75, 0x74, 0x69, 0x6c,
	0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0xa2, 0x02, 0x03, 0x55, 0x45, 0x58, 0xaa, 0x02, 0x0a, 0x55,
	0x74, 0x69, 0x6c, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0xca, 0x02, 0x0b, 0x55, 0x74, 0x69, 0x6c,
	0x5c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0xe2, 0x02, 0x17, 0x55, 0x74, 0x69, 0x6c, 0x5c, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x5f, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x0b, 0x55, 0x74, 0x69, 0x6c, 0x3a, 0x3a, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_util_error_error_proto_rawDescOnce sync.Once
	file_util_error_error_proto_rawDescData = file_util_error_error_proto_rawDesc
)

func file_util_error_error_proto_rawDescGZIP() []byte {
	file_util_error_error_proto_rawDescOnce.Do(func() {
		file_util_error_error_proto_rawDescData = protoimpl.X.CompressGZIP(file_util_error_error_proto_rawDescData)
	})
	return file_util_error_error_proto_rawDescData
}

var file_util_error_error_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_util_error_error_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_util_error_error_proto_goTypes = []interface{}{
	(ErrorCode)(0),      // 0: util.error.ErrorCode
	(*ErrorResult)(nil), // 1: util.error.ErrorResult
}
var file_util_error_error_proto_depIdxs = []int32{
	0, // 0: util.error.ErrorResult.error_code:type_name -> util.error.ErrorCode
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_util_error_error_proto_init() }
func file_util_error_error_proto_init() {
	if File_util_error_error_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_util_error_error_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ErrorResult); i {
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
			RawDescriptor: file_util_error_error_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_util_error_error_proto_goTypes,
		DependencyIndexes: file_util_error_error_proto_depIdxs,
		EnumInfos:         file_util_error_error_proto_enumTypes,
		MessageInfos:      file_util_error_error_proto_msgTypes,
	}.Build()
	File_util_error_error_proto = out.File
	file_util_error_error_proto_rawDesc = nil
	file_util_error_error_proto_goTypes = nil
	file_util_error_error_proto_depIdxs = nil
}
