// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: client_storage.proto

package client_storage

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
	ErrorCode_NO_ERROR                 ErrorCode = 0
	ErrorCode_FILE_NOT_FOUND           ErrorCode = 1
	ErrorCode_FILE_ALREADY_EXISTS      ErrorCode = 2
	ErrorCode_FILE_SIZE_LIMIT_EXCEEDED ErrorCode = 3
	ErrorCode_SERVER_ERROR             ErrorCode = 4
	ErrorCode_CHECKSUM_MISMATCH        ErrorCode = 5
)

// Enum value maps for ErrorCode.
var (
	ErrorCode_name = map[int32]string{
		0: "NO_ERROR",
		1: "FILE_NOT_FOUND",
		2: "FILE_ALREADY_EXISTS",
		3: "FILE_SIZE_LIMIT_EXCEEDED",
		4: "SERVER_ERROR",
		5: "CHECKSUM_MISMATCH",
	}
	ErrorCode_value = map[string]int32{
		"NO_ERROR":                 0,
		"FILE_NOT_FOUND":           1,
		"FILE_ALREADY_EXISTS":      2,
		"FILE_SIZE_LIMIT_EXCEEDED": 3,
		"SERVER_ERROR":             4,
		"CHECKSUM_MISMATCH":        5,
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
	return file_client_storage_proto_enumTypes[0].Descriptor()
}

func (ErrorCode) Type() protoreflect.EnumType {
	return &file_client_storage_proto_enumTypes[0]
}

func (x ErrorCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorCode.Descriptor instead.
func (ErrorCode) EnumDescriptor() ([]byte, []int) {
	return file_client_storage_proto_rawDescGZIP(), []int{0}
}

// Client requests to store a file on the server
type FilePutRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileName string `protobuf:"bytes,1,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	FileSize int64  `protobuf:"varint,2,opt,name=file_size,json=fileSize,proto3" json:"file_size,omitempty"`
}

func (x *FilePutRequest) Reset() {
	*x = FilePutRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_storage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilePutRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilePutRequest) ProtoMessage() {}

func (x *FilePutRequest) ProtoReflect() protoreflect.Message {
	mi := &file_client_storage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilePutRequest.ProtoReflect.Descriptor instead.
func (*FilePutRequest) Descriptor() ([]byte, []int) {
	return file_client_storage_proto_rawDescGZIP(), []int{0}
}

func (x *FilePutRequest) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *FilePutRequest) GetFileSize() int64 {
	if x != nil {
		return x.FileSize
	}
	return 0
}

// Server response to a FilePutRequest
type FilePutResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success   bool      `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	ErrorCode ErrorCode `protobuf:"varint,2,opt,name=error_code,json=errorCode,proto3,enum=ErrorCode" json:"error_code,omitempty"`
	FileName  string    `protobuf:"bytes,3,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
}

func (x *FilePutResponse) Reset() {
	*x = FilePutResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_storage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilePutResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilePutResponse) ProtoMessage() {}

func (x *FilePutResponse) ProtoReflect() protoreflect.Message {
	mi := &file_client_storage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilePutResponse.ProtoReflect.Descriptor instead.
func (*FilePutResponse) Descriptor() ([]byte, []int) {
	return file_client_storage_proto_rawDescGZIP(), []int{1}
}

func (x *FilePutResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *FilePutResponse) GetErrorCode() ErrorCode {
	if x != nil {
		return x.ErrorCode
	}
	return ErrorCode_NO_ERROR
}

func (x *FilePutResponse) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

type FileDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileName    string   `protobuf:"bytes,1,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	MessageBody []byte   `protobuf:"bytes,2,opt,name=message_body,json=messageBody,proto3" json:"message_body,omitempty"`
	Checksum    []byte   `protobuf:"bytes,3,opt,name=checksum,proto3" json:"checksum,omitempty"`
	OtherNodes  []string `protobuf:"bytes,4,rep,name=other_nodes,json=otherNodes,proto3" json:"other_nodes,omitempty"`
}

func (x *FileDataRequest) Reset() {
	*x = FileDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_storage_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileDataRequest) ProtoMessage() {}

func (x *FileDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_client_storage_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileDataRequest.ProtoReflect.Descriptor instead.
func (*FileDataRequest) Descriptor() ([]byte, []int) {
	return file_client_storage_proto_rawDescGZIP(), []int{2}
}

func (x *FileDataRequest) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *FileDataRequest) GetMessageBody() []byte {
	if x != nil {
		return x.MessageBody
	}
	return nil
}

func (x *FileDataRequest) GetChecksum() []byte {
	if x != nil {
		return x.Checksum
	}
	return nil
}

func (x *FileDataRequest) GetOtherNodes() []string {
	if x != nil {
		return x.OtherNodes
	}
	return nil
}

// Server response to a FileDataRequest
type FileDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success   bool      `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	ErrorCode ErrorCode `protobuf:"varint,2,opt,name=error_code,json=errorCode,proto3,enum=ErrorCode" json:"error_code,omitempty"`
}

func (x *FileDataResponse) Reset() {
	*x = FileDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_storage_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileDataResponse) ProtoMessage() {}

func (x *FileDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_client_storage_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileDataResponse.ProtoReflect.Descriptor instead.
func (*FileDataResponse) Descriptor() ([]byte, []int) {
	return file_client_storage_proto_rawDescGZIP(), []int{3}
}

func (x *FileDataResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *FileDataResponse) GetErrorCode() ErrorCode {
	if x != nil {
		return x.ErrorCode
	}
	return ErrorCode_NO_ERROR
}

// Client requests to retrieve a file from the server
type FileGetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileName string `protobuf:"bytes,1,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
}

func (x *FileGetRequest) Reset() {
	*x = FileGetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_storage_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileGetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileGetRequest) ProtoMessage() {}

func (x *FileGetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_client_storage_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileGetRequest.ProtoReflect.Descriptor instead.
func (*FileGetRequest) Descriptor() ([]byte, []int) {
	return file_client_storage_proto_rawDescGZIP(), []int{4}
}

func (x *FileGetRequest) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

// Server response to a FileGetRequest
type FileGetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success     bool      `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	FileSize    int64     `protobuf:"varint,2,opt,name=file_size,json=fileSize,proto3" json:"file_size,omitempty"`
	Checksum    []byte    `protobuf:"bytes,3,opt,name=checksum,proto3" json:"checksum,omitempty"`
	MessageBody []byte    `protobuf:"bytes,4,opt,name=message_body,json=messageBody,proto3" json:"message_body,omitempty"`
	ErrorCode   ErrorCode `protobuf:"varint,5,opt,name=error_code,json=errorCode,proto3,enum=ErrorCode" json:"error_code,omitempty"`
}

func (x *FileGetResponse) Reset() {
	*x = FileGetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_storage_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileGetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileGetResponse) ProtoMessage() {}

func (x *FileGetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_client_storage_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileGetResponse.ProtoReflect.Descriptor instead.
func (*FileGetResponse) Descriptor() ([]byte, []int) {
	return file_client_storage_proto_rawDescGZIP(), []int{5}
}

func (x *FileGetResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *FileGetResponse) GetFileSize() int64 {
	if x != nil {
		return x.FileSize
	}
	return 0
}

func (x *FileGetResponse) GetChecksum() []byte {
	if x != nil {
		return x.Checksum
	}
	return nil
}

func (x *FileGetResponse) GetMessageBody() []byte {
	if x != nil {
		return x.MessageBody
	}
	return nil
}

func (x *FileGetResponse) GetErrorCode() ErrorCode {
	if x != nil {
		return x.ErrorCode
	}
	return ErrorCode_NO_ERROR
}

// Server notifies the client that the file transfer is complete
type FileTransferComplete struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success   bool      `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	ErrorCode ErrorCode `protobuf:"varint,2,opt,name=error_code,json=errorCode,proto3,enum=ErrorCode" json:"error_code,omitempty"`
}

func (x *FileTransferComplete) Reset() {
	*x = FileTransferComplete{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_storage_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileTransferComplete) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileTransferComplete) ProtoMessage() {}

func (x *FileTransferComplete) ProtoReflect() protoreflect.Message {
	mi := &file_client_storage_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileTransferComplete.ProtoReflect.Descriptor instead.
func (*FileTransferComplete) Descriptor() ([]byte, []int) {
	return file_client_storage_proto_rawDescGZIP(), []int{6}
}

func (x *FileTransferComplete) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *FileTransferComplete) GetErrorCode() ErrorCode {
	if x != nil {
		return x.ErrorCode
	}
	return ErrorCode_NO_ERROR
}

// A wrapper message for client requests
type ClientRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Request:
	//
	//	*ClientRequest_FilePutRequest
	//	*ClientRequest_FileDataRequest
	//	*ClientRequest_FileGetRequest
	Request isClientRequest_Request `protobuf_oneof:"request"`
}

func (x *ClientRequest) Reset() {
	*x = ClientRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_storage_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientRequest) ProtoMessage() {}

func (x *ClientRequest) ProtoReflect() protoreflect.Message {
	mi := &file_client_storage_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientRequest.ProtoReflect.Descriptor instead.
func (*ClientRequest) Descriptor() ([]byte, []int) {
	return file_client_storage_proto_rawDescGZIP(), []int{7}
}

func (m *ClientRequest) GetRequest() isClientRequest_Request {
	if m != nil {
		return m.Request
	}
	return nil
}

func (x *ClientRequest) GetFilePutRequest() *FilePutRequest {
	if x, ok := x.GetRequest().(*ClientRequest_FilePutRequest); ok {
		return x.FilePutRequest
	}
	return nil
}

func (x *ClientRequest) GetFileDataRequest() *FileDataRequest {
	if x, ok := x.GetRequest().(*ClientRequest_FileDataRequest); ok {
		return x.FileDataRequest
	}
	return nil
}

func (x *ClientRequest) GetFileGetRequest() *FileGetRequest {
	if x, ok := x.GetRequest().(*ClientRequest_FileGetRequest); ok {
		return x.FileGetRequest
	}
	return nil
}

type isClientRequest_Request interface {
	isClientRequest_Request()
}

type ClientRequest_FilePutRequest struct {
	FilePutRequest *FilePutRequest `protobuf:"bytes,1,opt,name=file_put_request,json=filePutRequest,proto3,oneof"`
}

type ClientRequest_FileDataRequest struct {
	FileDataRequest *FileDataRequest `protobuf:"bytes,2,opt,name=file_data_request,json=fileDataRequest,proto3,oneof"`
}

type ClientRequest_FileGetRequest struct {
	FileGetRequest *FileGetRequest `protobuf:"bytes,3,opt,name=file_get_request,json=fileGetRequest,proto3,oneof"`
}

func (*ClientRequest_FilePutRequest) isClientRequest_Request() {}

func (*ClientRequest_FileDataRequest) isClientRequest_Request() {}

func (*ClientRequest_FileGetRequest) isClientRequest_Request() {}

// A wrapper message for server responses
type ServerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Response:
	//
	//	*ServerResponse_FilePutResponse
	//	*ServerResponse_FileDataResponse
	//	*ServerResponse_FileGetResponse
	//	*ServerResponse_FileTransferComplete
	Response isServerResponse_Response `protobuf_oneof:"response"`
}

func (x *ServerResponse) Reset() {
	*x = ServerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_storage_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerResponse) ProtoMessage() {}

func (x *ServerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_client_storage_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerResponse.ProtoReflect.Descriptor instead.
func (*ServerResponse) Descriptor() ([]byte, []int) {
	return file_client_storage_proto_rawDescGZIP(), []int{8}
}

func (m *ServerResponse) GetResponse() isServerResponse_Response {
	if m != nil {
		return m.Response
	}
	return nil
}

func (x *ServerResponse) GetFilePutResponse() *FilePutResponse {
	if x, ok := x.GetResponse().(*ServerResponse_FilePutResponse); ok {
		return x.FilePutResponse
	}
	return nil
}

func (x *ServerResponse) GetFileDataResponse() *FileDataResponse {
	if x, ok := x.GetResponse().(*ServerResponse_FileDataResponse); ok {
		return x.FileDataResponse
	}
	return nil
}

func (x *ServerResponse) GetFileGetResponse() *FileGetResponse {
	if x, ok := x.GetResponse().(*ServerResponse_FileGetResponse); ok {
		return x.FileGetResponse
	}
	return nil
}

func (x *ServerResponse) GetFileTransferComplete() *FileTransferComplete {
	if x, ok := x.GetResponse().(*ServerResponse_FileTransferComplete); ok {
		return x.FileTransferComplete
	}
	return nil
}

type isServerResponse_Response interface {
	isServerResponse_Response()
}

type ServerResponse_FilePutResponse struct {
	FilePutResponse *FilePutResponse `protobuf:"bytes,1,opt,name=file_put_response,json=filePutResponse,proto3,oneof"`
}

type ServerResponse_FileDataResponse struct {
	FileDataResponse *FileDataResponse `protobuf:"bytes,2,opt,name=file_data_response,json=fileDataResponse,proto3,oneof"`
}

type ServerResponse_FileGetResponse struct {
	FileGetResponse *FileGetResponse `protobuf:"bytes,3,opt,name=file_get_response,json=fileGetResponse,proto3,oneof"`
}

type ServerResponse_FileTransferComplete struct {
	FileTransferComplete *FileTransferComplete `protobuf:"bytes,5,opt,name=file_transfer_complete,json=fileTransferComplete,proto3,oneof"`
}

func (*ServerResponse_FilePutResponse) isServerResponse_Response() {}

func (*ServerResponse_FileDataResponse) isServerResponse_Response() {}

func (*ServerResponse_FileGetResponse) isServerResponse_Response() {}

func (*ServerResponse_FileTransferComplete) isServerResponse_Response() {}

var File_client_storage_proto protoreflect.FileDescriptor

var file_client_storage_proto_rawDesc = []byte{
	0x0a, 0x14, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4a, 0x0a, 0x0e, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x75,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x73, 0x69,
	0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x53, 0x69,
	0x7a, 0x65, 0x22, 0x73, 0x0a, 0x0f, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x75, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12,
	0x29, 0x0a, 0x0a, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x0a, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x52,
	0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69,
	0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x8e, 0x01, 0x0a, 0x0f, 0x46, 0x69, 0x6c, 0x65,
	0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x66,
	0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x5f, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x63,
	0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x63,
	0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x12, 0x1f, 0x0a, 0x0b, 0x6f, 0x74, 0x68, 0x65, 0x72,
	0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x6f, 0x74,
	0x68, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x22, 0x57, 0x0a, 0x10, 0x46, 0x69, 0x6c, 0x65,
	0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x29, 0x0a, 0x0a, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f,
	0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0a, 0x2e, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64,
	0x65, 0x22, 0x2d, 0x0a, 0x0e, 0x46, 0x69, 0x6c, 0x65, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x22, 0xb2, 0x01, 0x0a, 0x0f, 0x46, 0x69, 0x6c, 0x65, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x1b,
	0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x63,
	0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x63,
	0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x5f, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x29, 0x0a, 0x0a, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0a,
	0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x09, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x5b, 0x0a, 0x14, 0x46, 0x69, 0x6c, 0x65, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x66, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x29, 0x0a, 0x0a, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0a, 0x2e, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f,
	0x64, 0x65, 0x22, 0xd4, 0x01, 0x0a, 0x0d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x3b, 0x0a, 0x10, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x70, 0x75, 0x74,
	0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x75, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48,
	0x00, 0x52, 0x0e, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x75, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x3e, 0x0a, 0x11, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x00,
	0x52, 0x0f, 0x66, 0x69, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x3b, 0x0a, 0x10, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x67, 0x65, 0x74, 0x5f, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x46, 0x69,
	0x6c, 0x65, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52, 0x0e,
	0x66, 0x69, 0x6c, 0x65, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x09,
	0x0a, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0xae, 0x02, 0x0a, 0x0e, 0x53, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x11,
	0x66, 0x69, 0x6c, 0x65, 0x5f, 0x70, 0x75, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x75,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x00, 0x52, 0x0f, 0x66, 0x69, 0x6c,
	0x65, 0x50, 0x75, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x12,
	0x66, 0x69, 0x6c, 0x65, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x44,
	0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x00, 0x52, 0x10, 0x66,
	0x69, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x3e, 0x0a, 0x11, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x67, 0x65, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x46, 0x69, 0x6c,
	0x65, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x00, 0x52, 0x0f,
	0x66, 0x69, 0x6c, 0x65, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x4d, 0x0a, 0x16, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72,
	0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x15, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x43, 0x6f,
	0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x48, 0x00, 0x52, 0x14, 0x66, 0x69, 0x6c, 0x65, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x42, 0x0a,
	0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2a, 0x8d, 0x01, 0x0a, 0x09, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x0c, 0x0a, 0x08, 0x4e, 0x4f, 0x5f, 0x45,
	0x52, 0x52, 0x4f, 0x52, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x46, 0x49, 0x4c, 0x45, 0x5f, 0x4e,
	0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x46, 0x49,
	0x4c, 0x45, 0x5f, 0x41, 0x4c, 0x52, 0x45, 0x41, 0x44, 0x59, 0x5f, 0x45, 0x58, 0x49, 0x53, 0x54,
	0x53, 0x10, 0x02, 0x12, 0x1c, 0x0a, 0x18, 0x46, 0x49, 0x4c, 0x45, 0x5f, 0x53, 0x49, 0x5a, 0x45,
	0x5f, 0x4c, 0x49, 0x4d, 0x49, 0x54, 0x5f, 0x45, 0x58, 0x43, 0x45, 0x45, 0x44, 0x45, 0x44, 0x10,
	0x03, 0x12, 0x10, 0x0a, 0x0c, 0x53, 0x45, 0x52, 0x56, 0x45, 0x52, 0x5f, 0x45, 0x52, 0x52, 0x4f,
	0x52, 0x10, 0x04, 0x12, 0x15, 0x0a, 0x11, 0x43, 0x48, 0x45, 0x43, 0x4b, 0x53, 0x55, 0x4d, 0x5f,
	0x4d, 0x49, 0x53, 0x4d, 0x41, 0x54, 0x43, 0x48, 0x10, 0x05, 0x42, 0x1b, 0x5a, 0x19, 0x2e, 0x2f,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f,
	0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_client_storage_proto_rawDescOnce sync.Once
	file_client_storage_proto_rawDescData = file_client_storage_proto_rawDesc
)

func file_client_storage_proto_rawDescGZIP() []byte {
	file_client_storage_proto_rawDescOnce.Do(func() {
		file_client_storage_proto_rawDescData = protoimpl.X.CompressGZIP(file_client_storage_proto_rawDescData)
	})
	return file_client_storage_proto_rawDescData
}

var file_client_storage_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_client_storage_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_client_storage_proto_goTypes = []interface{}{
	(ErrorCode)(0),               // 0: ErrorCode
	(*FilePutRequest)(nil),       // 1: FilePutRequest
	(*FilePutResponse)(nil),      // 2: FilePutResponse
	(*FileDataRequest)(nil),      // 3: FileDataRequest
	(*FileDataResponse)(nil),     // 4: FileDataResponse
	(*FileGetRequest)(nil),       // 5: FileGetRequest
	(*FileGetResponse)(nil),      // 6: FileGetResponse
	(*FileTransferComplete)(nil), // 7: FileTransferComplete
	(*ClientRequest)(nil),        // 8: ClientRequest
	(*ServerResponse)(nil),       // 9: ServerResponse
}
var file_client_storage_proto_depIdxs = []int32{
	0,  // 0: FilePutResponse.error_code:type_name -> ErrorCode
	0,  // 1: FileDataResponse.error_code:type_name -> ErrorCode
	0,  // 2: FileGetResponse.error_code:type_name -> ErrorCode
	0,  // 3: FileTransferComplete.error_code:type_name -> ErrorCode
	1,  // 4: ClientRequest.file_put_request:type_name -> FilePutRequest
	3,  // 5: ClientRequest.file_data_request:type_name -> FileDataRequest
	5,  // 6: ClientRequest.file_get_request:type_name -> FileGetRequest
	2,  // 7: ServerResponse.file_put_response:type_name -> FilePutResponse
	4,  // 8: ServerResponse.file_data_response:type_name -> FileDataResponse
	6,  // 9: ServerResponse.file_get_response:type_name -> FileGetResponse
	7,  // 10: ServerResponse.file_transfer_complete:type_name -> FileTransferComplete
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_client_storage_proto_init() }
func file_client_storage_proto_init() {
	if File_client_storage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_client_storage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilePutRequest); i {
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
		file_client_storage_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilePutResponse); i {
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
		file_client_storage_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileDataRequest); i {
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
		file_client_storage_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileDataResponse); i {
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
		file_client_storage_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileGetRequest); i {
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
		file_client_storage_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileGetResponse); i {
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
		file_client_storage_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileTransferComplete); i {
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
		file_client_storage_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientRequest); i {
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
		file_client_storage_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerResponse); i {
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
	file_client_storage_proto_msgTypes[7].OneofWrappers = []interface{}{
		(*ClientRequest_FilePutRequest)(nil),
		(*ClientRequest_FileDataRequest)(nil),
		(*ClientRequest_FileGetRequest)(nil),
	}
	file_client_storage_proto_msgTypes[8].OneofWrappers = []interface{}{
		(*ServerResponse_FilePutResponse)(nil),
		(*ServerResponse_FileDataResponse)(nil),
		(*ServerResponse_FileGetResponse)(nil),
		(*ServerResponse_FileTransferComplete)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_client_storage_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_client_storage_proto_goTypes,
		DependencyIndexes: file_client_storage_proto_depIdxs,
		EnumInfos:         file_client_storage_proto_enumTypes,
		MessageInfos:      file_client_storage_proto_msgTypes,
	}.Build()
	File_client_storage_proto = out.File
	file_client_storage_proto_rawDesc = nil
	file_client_storage_proto_goTypes = nil
	file_client_storage_proto_depIdxs = nil
}