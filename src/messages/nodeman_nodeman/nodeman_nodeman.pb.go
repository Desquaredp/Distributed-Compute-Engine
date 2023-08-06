// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: nodeman_nodeman.proto

package nodeman_nodeman

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

type NodeManagerRequest_Status int32

const (
	NodeManagerRequest_SUCCESS NodeManagerRequest_Status = 0
	NodeManagerRequest_FAILURE NodeManagerRequest_Status = 1
)

// Enum value maps for NodeManagerRequest_Status.
var (
	NodeManagerRequest_Status_name = map[int32]string{
		0: "SUCCESS",
		1: "FAILURE",
	}
	NodeManagerRequest_Status_value = map[string]int32{
		"SUCCESS": 0,
		"FAILURE": 1,
	}
)

func (x NodeManagerRequest_Status) Enum() *NodeManagerRequest_Status {
	p := new(NodeManagerRequest_Status)
	*p = x
	return p
}

func (x NodeManagerRequest_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NodeManagerRequest_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_nodeman_nodeman_proto_enumTypes[0].Descriptor()
}

func (NodeManagerRequest_Status) Type() protoreflect.EnumType {
	return &file_nodeman_nodeman_proto_enumTypes[0]
}

func (x NodeManagerRequest_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NodeManagerRequest_Status.Descriptor instead.
func (NodeManagerRequest_Status) EnumDescriptor() ([]byte, []int) {
	return file_nodeman_nodeman_proto_rawDescGZIP(), []int{0, 0}
}

type NodeManagerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Request:
	//
	//	*NodeManagerRequest_ReducerFile_
	//	*NodeManagerRequest_ReducerFileResponse_
	Request isNodeManagerRequest_Request `protobuf_oneof:"request"`
}

func (x *NodeManagerRequest) Reset() {
	*x = NodeManagerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nodeman_nodeman_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeManagerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeManagerRequest) ProtoMessage() {}

func (x *NodeManagerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_nodeman_nodeman_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeManagerRequest.ProtoReflect.Descriptor instead.
func (*NodeManagerRequest) Descriptor() ([]byte, []int) {
	return file_nodeman_nodeman_proto_rawDescGZIP(), []int{0}
}

func (m *NodeManagerRequest) GetRequest() isNodeManagerRequest_Request {
	if m != nil {
		return m.Request
	}
	return nil
}

func (x *NodeManagerRequest) GetReducerFile() *NodeManagerRequest_ReducerFile {
	if x, ok := x.GetRequest().(*NodeManagerRequest_ReducerFile_); ok {
		return x.ReducerFile
	}
	return nil
}

func (x *NodeManagerRequest) GetReducerFileResponse() *NodeManagerRequest_ReducerFileResponse {
	if x, ok := x.GetRequest().(*NodeManagerRequest_ReducerFileResponse_); ok {
		return x.ReducerFileResponse
	}
	return nil
}

type isNodeManagerRequest_Request interface {
	isNodeManagerRequest_Request()
}

type NodeManagerRequest_ReducerFile_ struct {
	ReducerFile *NodeManagerRequest_ReducerFile `protobuf:"bytes,1,opt,name=reducer_file,json=reducerFile,proto3,oneof"`
}

type NodeManagerRequest_ReducerFileResponse_ struct {
	ReducerFileResponse *NodeManagerRequest_ReducerFileResponse `protobuf:"bytes,2,opt,name=reducer_file_response,json=reducerFileResponse,proto3,oneof"`
}

func (*NodeManagerRequest_ReducerFile_) isNodeManagerRequest_Request() {}

func (*NodeManagerRequest_ReducerFileResponse_) isNodeManagerRequest_Request() {}

type NodeManagerRequest_ReducerFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileName string `protobuf:"bytes,1,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	FileData []byte `protobuf:"bytes,2,opt,name=file_data,json=fileData,proto3" json:"file_data,omitempty"`
	Checksum []byte `protobuf:"bytes,3,opt,name=checksum,proto3" json:"checksum,omitempty"`
}

func (x *NodeManagerRequest_ReducerFile) Reset() {
	*x = NodeManagerRequest_ReducerFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nodeman_nodeman_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeManagerRequest_ReducerFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeManagerRequest_ReducerFile) ProtoMessage() {}

func (x *NodeManagerRequest_ReducerFile) ProtoReflect() protoreflect.Message {
	mi := &file_nodeman_nodeman_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeManagerRequest_ReducerFile.ProtoReflect.Descriptor instead.
func (*NodeManagerRequest_ReducerFile) Descriptor() ([]byte, []int) {
	return file_nodeman_nodeman_proto_rawDescGZIP(), []int{0, 0}
}

func (x *NodeManagerRequest_ReducerFile) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *NodeManagerRequest_ReducerFile) GetFileData() []byte {
	if x != nil {
		return x.FileData
	}
	return nil
}

func (x *NodeManagerRequest_ReducerFile) GetChecksum() []byte {
	if x != nil {
		return x.Checksum
	}
	return nil
}

type NodeManagerRequest_ReducerFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status NodeManagerRequest_Status `protobuf:"varint,1,opt,name=status,proto3,enum=NodeManagerRequest_Status" json:"status,omitempty"`
}

func (x *NodeManagerRequest_ReducerFileResponse) Reset() {
	*x = NodeManagerRequest_ReducerFileResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nodeman_nodeman_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeManagerRequest_ReducerFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeManagerRequest_ReducerFileResponse) ProtoMessage() {}

func (x *NodeManagerRequest_ReducerFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_nodeman_nodeman_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeManagerRequest_ReducerFileResponse.ProtoReflect.Descriptor instead.
func (*NodeManagerRequest_ReducerFileResponse) Descriptor() ([]byte, []int) {
	return file_nodeman_nodeman_proto_rawDescGZIP(), []int{0, 1}
}

func (x *NodeManagerRequest_ReducerFileResponse) GetStatus() NodeManagerRequest_Status {
	if x != nil {
		return x.Status
	}
	return NodeManagerRequest_SUCCESS
}

var File_nodeman_nodeman_proto protoreflect.FileDescriptor

var file_nodeman_nodeman_proto_rawDesc = []byte{
	0x0a, 0x15, 0x6e, 0x6f, 0x64, 0x65, 0x6d, 0x61, 0x6e, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x6d, 0x61,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x98, 0x03, 0x0a, 0x12, 0x4e, 0x6f, 0x64, 0x65,
	0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x44,
	0x0a, 0x0c, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x72, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x4d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65,
	0x72, 0x46, 0x69, 0x6c, 0x65, 0x48, 0x00, 0x52, 0x0b, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x72,
	0x46, 0x69, 0x6c, 0x65, 0x12, 0x5d, 0x0a, 0x15, 0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x72, 0x5f,
	0x66, 0x69, 0x6c, 0x65, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x72,
	0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x00, 0x52, 0x13,
	0x72, 0x65, 0x64, 0x75, 0x63, 0x65, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x1a, 0x63, 0x0a, 0x0b, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x72, 0x46, 0x69,
	0x6c, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08,
	0x63, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08,
	0x63, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x1a, 0x49, 0x0a, 0x13, 0x52, 0x65, 0x64, 0x75,
	0x63, 0x65, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x32, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x1a, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x22, 0x22, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a,
	0x07, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x46, 0x41,
	0x49, 0x4c, 0x55, 0x52, 0x45, 0x10, 0x01, 0x42, 0x09, 0x0a, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x42, 0x1c, 0x5a, 0x1a, 0x2e, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x6d, 0x61, 0x6e, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x6d, 0x61, 0x6e,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_nodeman_nodeman_proto_rawDescOnce sync.Once
	file_nodeman_nodeman_proto_rawDescData = file_nodeman_nodeman_proto_rawDesc
)

func file_nodeman_nodeman_proto_rawDescGZIP() []byte {
	file_nodeman_nodeman_proto_rawDescOnce.Do(func() {
		file_nodeman_nodeman_proto_rawDescData = protoimpl.X.CompressGZIP(file_nodeman_nodeman_proto_rawDescData)
	})
	return file_nodeman_nodeman_proto_rawDescData
}

var file_nodeman_nodeman_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_nodeman_nodeman_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_nodeman_nodeman_proto_goTypes = []interface{}{
	(NodeManagerRequest_Status)(0),                 // 0: NodeManagerRequest.Status
	(*NodeManagerRequest)(nil),                     // 1: NodeManagerRequest
	(*NodeManagerRequest_ReducerFile)(nil),         // 2: NodeManagerRequest.ReducerFile
	(*NodeManagerRequest_ReducerFileResponse)(nil), // 3: NodeManagerRequest.ReducerFileResponse
}
var file_nodeman_nodeman_proto_depIdxs = []int32{
	2, // 0: NodeManagerRequest.reducer_file:type_name -> NodeManagerRequest.ReducerFile
	3, // 1: NodeManagerRequest.reducer_file_response:type_name -> NodeManagerRequest.ReducerFileResponse
	0, // 2: NodeManagerRequest.ReducerFileResponse.status:type_name -> NodeManagerRequest.Status
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_nodeman_nodeman_proto_init() }
func file_nodeman_nodeman_proto_init() {
	if File_nodeman_nodeman_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_nodeman_nodeman_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeManagerRequest); i {
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
		file_nodeman_nodeman_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeManagerRequest_ReducerFile); i {
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
		file_nodeman_nodeman_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeManagerRequest_ReducerFileResponse); i {
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
	file_nodeman_nodeman_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*NodeManagerRequest_ReducerFile_)(nil),
		(*NodeManagerRequest_ReducerFileResponse_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_nodeman_nodeman_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_nodeman_nodeman_proto_goTypes,
		DependencyIndexes: file_nodeman_nodeman_proto_depIdxs,
		EnumInfos:         file_nodeman_nodeman_proto_enumTypes,
		MessageInfos:      file_nodeman_nodeman_proto_msgTypes,
	}.Build()
	File_nodeman_nodeman_proto = out.File
	file_nodeman_nodeman_proto_rawDesc = nil
	file_nodeman_nodeman_proto_goTypes = nil
	file_nodeman_nodeman_proto_depIdxs = nil
}
