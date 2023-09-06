// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.2
// source: internal/auth/auth.proto

package auth

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AuthToken struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthTime float64               `protobuf:"fixed64,1,opt,name=authTime,proto3" json:"authTime,omitempty"`
	Issuer   string                `protobuf:"bytes,2,opt,name=issuer,proto3" json:"issuer,omitempty"`
	Audience string                `protobuf:"bytes,3,opt,name=audience,proto3" json:"audience,omitempty"`
	Expires  float64               `protobuf:"fixed64,4,opt,name=expires,proto3" json:"expires,omitempty"`
	IssuedAt float64               `protobuf:"fixed64,5,opt,name=issuedAt,proto3" json:"issuedAt,omitempty"`
	Subject  string                `protobuf:"bytes,6,opt,name=subject,proto3" json:"subject,omitempty"`
	UID      string                `protobuf:"bytes,7,opt,name=UID,proto3" json:"UID,omitempty"`
	Claims   map[string]*anypb.Any `protobuf:"bytes,8,rep,name=claims,proto3" json:"claims,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *AuthToken) Reset() {
	*x = AuthToken{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_auth_auth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthToken) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthToken) ProtoMessage() {}

func (x *AuthToken) ProtoReflect() protoreflect.Message {
	mi := &file_internal_auth_auth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthToken.ProtoReflect.Descriptor instead.
func (*AuthToken) Descriptor() ([]byte, []int) {
	return file_internal_auth_auth_proto_rawDescGZIP(), []int{0}
}

func (x *AuthToken) GetAuthTime() float64 {
	if x != nil {
		return x.AuthTime
	}
	return 0
}

func (x *AuthToken) GetIssuer() string {
	if x != nil {
		return x.Issuer
	}
	return ""
}

func (x *AuthToken) GetAudience() string {
	if x != nil {
		return x.Audience
	}
	return ""
}

func (x *AuthToken) GetExpires() float64 {
	if x != nil {
		return x.Expires
	}
	return 0
}

func (x *AuthToken) GetIssuedAt() float64 {
	if x != nil {
		return x.IssuedAt
	}
	return 0
}

func (x *AuthToken) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *AuthToken) GetUID() string {
	if x != nil {
		return x.UID
	}
	return ""
}

func (x *AuthToken) GetClaims() map[string]*anypb.Any {
	if x != nil {
		return x.Claims
	}
	return nil
}

type ValidateTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token    string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Audience string `protobuf:"bytes,2,opt,name=audience,proto3" json:"audience,omitempty"`
}

func (x *ValidateTokenRequest) Reset() {
	*x = ValidateTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_auth_auth_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidateTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateTokenRequest) ProtoMessage() {}

func (x *ValidateTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_auth_auth_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateTokenRequest.ProtoReflect.Descriptor instead.
func (*ValidateTokenRequest) Descriptor() ([]byte, []int) {
	return file_internal_auth_auth_proto_rawDescGZIP(), []int{1}
}

func (x *ValidateTokenRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *ValidateTokenRequest) GetAudience() string {
	if x != nil {
		return x.Audience
	}
	return ""
}

var File_internal_auth_auth_proto protoreflect.FileDescriptor

var file_internal_auth_auth_proto_rawDesc = []byte{
	0x0a, 0x18, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f,
	0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x61, 0x75, 0x74, 0x68,
	0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc3, 0x02, 0x0a, 0x09,
	0x41, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x61, 0x75, 0x74,
	0x68, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x61, 0x75, 0x74,
	0x68, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x12, 0x1a, 0x0a,
	0x08, 0x61, 0x75, 0x64, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x61, 0x75, 0x64, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x78, 0x70,
	0x69, 0x72, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x07, 0x65, 0x78, 0x70, 0x69,
	0x72, 0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x73, 0x73, 0x75, 0x65, 0x64, 0x41, 0x74, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x69, 0x73, 0x73, 0x75, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x18, 0x0a, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x55, 0x49, 0x44,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x55, 0x49, 0x44, 0x12, 0x33, 0x0a, 0x06, 0x63,
	0x6c, 0x61, 0x69, 0x6d, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x61, 0x75,
	0x74, 0x68, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x43, 0x6c, 0x61,
	0x69, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x63, 0x6c, 0x61, 0x69, 0x6d, 0x73,
	0x1a, 0x4f, 0x0a, 0x0b, 0x43, 0x6c, 0x61, 0x69, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x2a, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x22, 0x48, 0x0a, 0x14, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12,
	0x1a, 0x0a, 0x08, 0x61, 0x75, 0x64, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x61, 0x75, 0x64, 0x69, 0x65, 0x6e, 0x63, 0x65, 0x32, 0x46, 0x0a, 0x04, 0x41,
	0x75, 0x74, 0x68, 0x12, 0x3e, 0x0a, 0x0d, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1a, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x56, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x0f, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x22, 0x00, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x62, 0x6a, 0x61, 0x72, 0x6b, 0x65, 0x2d, 0x78, 0x79, 0x7a, 0x2f, 0x61, 0x75, 0x74,
	0x68, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_auth_auth_proto_rawDescOnce sync.Once
	file_internal_auth_auth_proto_rawDescData = file_internal_auth_auth_proto_rawDesc
)

func file_internal_auth_auth_proto_rawDescGZIP() []byte {
	file_internal_auth_auth_proto_rawDescOnce.Do(func() {
		file_internal_auth_auth_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_auth_auth_proto_rawDescData)
	})
	return file_internal_auth_auth_proto_rawDescData
}

var file_internal_auth_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_internal_auth_auth_proto_goTypes = []interface{}{
	(*AuthToken)(nil),            // 0: auth.AuthToken
	(*ValidateTokenRequest)(nil), // 1: auth.ValidateTokenRequest
	nil,                          // 2: auth.AuthToken.ClaimsEntry
	(*anypb.Any)(nil),            // 3: google.protobuf.Any
}
var file_internal_auth_auth_proto_depIdxs = []int32{
	2, // 0: auth.AuthToken.claims:type_name -> auth.AuthToken.ClaimsEntry
	3, // 1: auth.AuthToken.ClaimsEntry.value:type_name -> google.protobuf.Any
	1, // 2: auth.Auth.ValidateToken:input_type -> auth.ValidateTokenRequest
	0, // 3: auth.Auth.ValidateToken:output_type -> auth.AuthToken
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_internal_auth_auth_proto_init() }
func file_internal_auth_auth_proto_init() {
	if File_internal_auth_auth_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_auth_auth_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthToken); i {
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
		file_internal_auth_auth_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValidateTokenRequest); i {
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
			RawDescriptor: file_internal_auth_auth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_auth_auth_proto_goTypes,
		DependencyIndexes: file_internal_auth_auth_proto_depIdxs,
		MessageInfos:      file_internal_auth_auth_proto_msgTypes,
	}.Build()
	File_internal_auth_auth_proto = out.File
	file_internal_auth_auth_proto_rawDesc = nil
	file_internal_auth_auth_proto_goTypes = nil
	file_internal_auth_auth_proto_depIdxs = nil
}