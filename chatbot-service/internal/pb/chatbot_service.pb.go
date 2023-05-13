// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: chatbot_service.proto

package pb

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

// Chat
type ChatRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Messages []*ChatRequest_ChatMessage `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"`
}

func (x *ChatRequest) Reset() {
	*x = ChatRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chatbot_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatRequest) ProtoMessage() {}

func (x *ChatRequest) ProtoReflect() protoreflect.Message {
	mi := &file_chatbot_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatRequest.ProtoReflect.Descriptor instead.
func (*ChatRequest) Descriptor() ([]byte, []int) {
	return file_chatbot_service_proto_rawDescGZIP(), []int{0}
}

func (x *ChatRequest) GetMessages() []*ChatRequest_ChatMessage {
	if x != nil {
		return x.Messages
	}
	return nil
}

type ChatResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Event string `protobuf:"bytes,1,opt,name=event,proto3" json:"event,omitempty"`
	Data  string `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *ChatResponse) Reset() {
	*x = ChatResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chatbot_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatResponse) ProtoMessage() {}

func (x *ChatResponse) ProtoReflect() protoreflect.Message {
	mi := &file_chatbot_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatResponse.ProtoReflect.Descriptor instead.
func (*ChatResponse) Descriptor() ([]byte, []int) {
	return file_chatbot_service_proto_rawDescGZIP(), []int{1}
}

func (x *ChatResponse) GetEvent() string {
	if x != nil {
		return x.Event
	}
	return ""
}

func (x *ChatResponse) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

type ChatRequest_ChatMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Role    string `protobuf:"bytes,1,opt,name=role,proto3" json:"role,omitempty"`
	Content string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *ChatRequest_ChatMessage) Reset() {
	*x = ChatRequest_ChatMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chatbot_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatRequest_ChatMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatRequest_ChatMessage) ProtoMessage() {}

func (x *ChatRequest_ChatMessage) ProtoReflect() protoreflect.Message {
	mi := &file_chatbot_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatRequest_ChatMessage.ProtoReflect.Descriptor instead.
func (*ChatRequest_ChatMessage) Descriptor() ([]byte, []int) {
	return file_chatbot_service_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ChatRequest_ChatMessage) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *ChatRequest_ChatMessage) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

var File_chatbot_service_proto protoreflect.FileDescriptor

var file_chatbot_service_proto_rawDesc = []byte{
	0x0a, 0x15, 0x63, 0x68, 0x61, 0x74, 0x62, 0x6f, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x22, 0x83, 0x01, 0x0a, 0x0b,
	0x43, 0x68, 0x61, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x37, 0x0a, 0x08, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e,
	0x70, 0x62, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x43,
	0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x73, 0x1a, 0x3b, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x22, 0x38, 0x0a, 0x0c, 0x43, 0x68, 0x61, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0x38, 0x0a, 0x07, 0x43,
	0x68, 0x61, 0x74, 0x62, 0x6f, 0x74, 0x12, 0x2d, 0x0a, 0x04, 0x43, 0x68, 0x61, 0x74, 0x12, 0x0f,
	0x2e, 0x70, 0x62, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x10, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x30, 0x01, 0x42, 0x3f, 0x5a, 0x3d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x42, 0x61, 0x6e, 0x61, 0x6e, 0x61, 0x2d, 0x42, 0x6f, 0x61, 0x74, 0x2f,
	0x74, 0x65, 0x72, 0x72, 0x79, 0x6d, 0x69, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x62,
	0x6f, 0x74, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_chatbot_service_proto_rawDescOnce sync.Once
	file_chatbot_service_proto_rawDescData = file_chatbot_service_proto_rawDesc
)

func file_chatbot_service_proto_rawDescGZIP() []byte {
	file_chatbot_service_proto_rawDescOnce.Do(func() {
		file_chatbot_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_chatbot_service_proto_rawDescData)
	})
	return file_chatbot_service_proto_rawDescData
}

var file_chatbot_service_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_chatbot_service_proto_goTypes = []interface{}{
	(*ChatRequest)(nil),             // 0: pb.ChatRequest
	(*ChatResponse)(nil),            // 1: pb.ChatResponse
	(*ChatRequest_ChatMessage)(nil), // 2: pb.ChatRequest.ChatMessage
}
var file_chatbot_service_proto_depIdxs = []int32{
	2, // 0: pb.ChatRequest.messages:type_name -> pb.ChatRequest.ChatMessage
	0, // 1: pb.Chatbot.Chat:input_type -> pb.ChatRequest
	1, // 2: pb.Chatbot.Chat:output_type -> pb.ChatResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_chatbot_service_proto_init() }
func file_chatbot_service_proto_init() {
	if File_chatbot_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_chatbot_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatRequest); i {
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
		file_chatbot_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatResponse); i {
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
		file_chatbot_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatRequest_ChatMessage); i {
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
			RawDescriptor: file_chatbot_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_chatbot_service_proto_goTypes,
		DependencyIndexes: file_chatbot_service_proto_depIdxs,
		MessageInfos:      file_chatbot_service_proto_msgTypes,
	}.Build()
	File_chatbot_service_proto = out.File
	file_chatbot_service_proto_rawDesc = nil
	file_chatbot_service_proto_goTypes = nil
	file_chatbot_service_proto_depIdxs = nil
}
