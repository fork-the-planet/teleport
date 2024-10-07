// Copyright 2024 Gravitational, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: teleport/access_graph/v1/authorized_key.proto

package accessgraphv1

import (
	v1 "github.com/gravitational/teleport/api/gen/proto/go/teleport/header/v1"
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

// The `AuthorizedKey` message represents an authorized key entry for a specific local user.
// These authorized keys are generated by the server when a particular SSH AuthorizedKey is granted access to a user on the node.
type AuthorizedKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// metadata is the AuthorizedKey's metadata.
	Metadata *v1.Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// kind is a resource kind.
	Kind string `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	// sub_kind is an optional resource sub kind, used in some resources.
	SubKind string `protobuf:"bytes,3,opt,name=sub_kind,json=subKind,proto3" json:"sub_kind,omitempty"`
	// version is version.
	Version string `protobuf:"bytes,4,opt,name=version,proto3" json:"version,omitempty"`
	// Spec is an AuthorizedKey specification.
	Spec *AuthorizedKeySpec `protobuf:"bytes,5,opt,name=spec,proto3" json:"spec,omitempty"`
}

func (x *AuthorizedKey) Reset() {
	*x = AuthorizedKey{}
	mi := &file_teleport_access_graph_v1_authorized_key_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AuthorizedKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthorizedKey) ProtoMessage() {}

func (x *AuthorizedKey) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_access_graph_v1_authorized_key_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthorizedKey.ProtoReflect.Descriptor instead.
func (*AuthorizedKey) Descriptor() ([]byte, []int) {
	return file_teleport_access_graph_v1_authorized_key_proto_rawDescGZIP(), []int{0}
}

func (x *AuthorizedKey) GetMetadata() *v1.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *AuthorizedKey) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *AuthorizedKey) GetSubKind() string {
	if x != nil {
		return x.SubKind
	}
	return ""
}

func (x *AuthorizedKey) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *AuthorizedKey) GetSpec() *AuthorizedKeySpec {
	if x != nil {
		return x.Spec
	}
	return nil
}

// AuthorizedKeySpec is the authorized key spec.
type AuthorizedKeySpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// host_id is the node identifier and must match the credentials used.
	HostId string `protobuf:"bytes,1,opt,name=host_id,json=hostId,proto3" json:"host_id,omitempty"`
	// key_fingerprint is the SHA256 SSH public key fingerprint.
	KeyFingerprint string `protobuf:"bytes,2,opt,name=key_fingerprint,json=keyFingerprint,proto3" json:"key_fingerprint,omitempty"`
	// host_user is the user who can be accessed using the fingerprint above.
	HostUser string `protobuf:"bytes,3,opt,name=host_user,json=hostUser,proto3" json:"host_user,omitempty"`
	// key_comment is the authorized key's comment.
	// Authorized keys consist of the following space-separated fields:
	// options, keytype, base64-encoded key, comment.  The options field is optional.
	KeyComment string `protobuf:"bytes,4,opt,name=key_comment,json=keyComment,proto3" json:"key_comment,omitempty"`
	// key_type is the ssh's key type.
	KeyType string `protobuf:"bytes,5,opt,name=key_type,json=keyType,proto3" json:"key_type,omitempty"`
}

func (x *AuthorizedKeySpec) Reset() {
	*x = AuthorizedKeySpec{}
	mi := &file_teleport_access_graph_v1_authorized_key_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AuthorizedKeySpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthorizedKeySpec) ProtoMessage() {}

func (x *AuthorizedKeySpec) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_access_graph_v1_authorized_key_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthorizedKeySpec.ProtoReflect.Descriptor instead.
func (*AuthorizedKeySpec) Descriptor() ([]byte, []int) {
	return file_teleport_access_graph_v1_authorized_key_proto_rawDescGZIP(), []int{1}
}

func (x *AuthorizedKeySpec) GetHostId() string {
	if x != nil {
		return x.HostId
	}
	return ""
}

func (x *AuthorizedKeySpec) GetKeyFingerprint() string {
	if x != nil {
		return x.KeyFingerprint
	}
	return ""
}

func (x *AuthorizedKeySpec) GetHostUser() string {
	if x != nil {
		return x.HostUser
	}
	return ""
}

func (x *AuthorizedKeySpec) GetKeyComment() string {
	if x != nil {
		return x.KeyComment
	}
	return ""
}

func (x *AuthorizedKeySpec) GetKeyType() string {
	if x != nil {
		return x.KeyType
	}
	return ""
}

var File_teleport_access_graph_v1_authorized_key_proto protoreflect.FileDescriptor

var file_teleport_access_graph_v1_authorized_key_proto_rawDesc = []byte{
	0x0a, 0x2d, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x61, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x5f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x69, 0x7a, 0x65, 0x64, 0x5f, 0x6b, 0x65, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x18, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x5f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x76, 0x31, 0x1a, 0x21, 0x74, 0x65, 0x6c, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x2f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd3, 0x01, 0x0a,
	0x0d, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x64, 0x4b, 0x65, 0x79, 0x12, 0x38,
	0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x68, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08,
	0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x19, 0x0a, 0x08,
	0x73, 0x75, 0x62, 0x5f, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x73, 0x75, 0x62, 0x4b, 0x69, 0x6e, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x3f, 0x0a, 0x04, 0x73, 0x70, 0x65, 0x63, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2b, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x5f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x69, 0x7a, 0x65, 0x64, 0x4b, 0x65, 0x79, 0x53, 0x70, 0x65, 0x63, 0x52, 0x04, 0x73, 0x70,
	0x65, 0x63, 0x22, 0xae, 0x01, 0x0a, 0x11, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65,
	0x64, 0x4b, 0x65, 0x79, 0x53, 0x70, 0x65, 0x63, 0x12, 0x17, 0x0a, 0x07, 0x68, 0x6f, 0x73, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x68, 0x6f, 0x73, 0x74, 0x49,
	0x64, 0x12, 0x27, 0x0a, 0x0f, 0x6b, 0x65, 0x79, 0x5f, 0x66, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70,
	0x72, 0x69, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6b, 0x65, 0x79, 0x46,
	0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x68, 0x6f,
	0x73, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68,
	0x6f, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1f, 0x0a, 0x0b, 0x6b, 0x65, 0x79, 0x5f, 0x63,
	0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6b, 0x65,
	0x79, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6b, 0x65, 0x79, 0x54,
	0x79, 0x70, 0x65, 0x42, 0x5a, 0x5a, 0x58, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x67, 0x72, 0x61, 0x76, 0x69, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x2f,
	0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x2f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2f, 0x76,
	0x31, 0x3b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x67, 0x72, 0x61, 0x70, 0x68, 0x76, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_teleport_access_graph_v1_authorized_key_proto_rawDescOnce sync.Once
	file_teleport_access_graph_v1_authorized_key_proto_rawDescData = file_teleport_access_graph_v1_authorized_key_proto_rawDesc
)

func file_teleport_access_graph_v1_authorized_key_proto_rawDescGZIP() []byte {
	file_teleport_access_graph_v1_authorized_key_proto_rawDescOnce.Do(func() {
		file_teleport_access_graph_v1_authorized_key_proto_rawDescData = protoimpl.X.CompressGZIP(file_teleport_access_graph_v1_authorized_key_proto_rawDescData)
	})
	return file_teleport_access_graph_v1_authorized_key_proto_rawDescData
}

var file_teleport_access_graph_v1_authorized_key_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_teleport_access_graph_v1_authorized_key_proto_goTypes = []any{
	(*AuthorizedKey)(nil),     // 0: teleport.access_graph.v1.AuthorizedKey
	(*AuthorizedKeySpec)(nil), // 1: teleport.access_graph.v1.AuthorizedKeySpec
	(*v1.Metadata)(nil),       // 2: teleport.header.v1.Metadata
}
var file_teleport_access_graph_v1_authorized_key_proto_depIdxs = []int32{
	2, // 0: teleport.access_graph.v1.AuthorizedKey.metadata:type_name -> teleport.header.v1.Metadata
	1, // 1: teleport.access_graph.v1.AuthorizedKey.spec:type_name -> teleport.access_graph.v1.AuthorizedKeySpec
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_teleport_access_graph_v1_authorized_key_proto_init() }
func file_teleport_access_graph_v1_authorized_key_proto_init() {
	if File_teleport_access_graph_v1_authorized_key_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_teleport_access_graph_v1_authorized_key_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_teleport_access_graph_v1_authorized_key_proto_goTypes,
		DependencyIndexes: file_teleport_access_graph_v1_authorized_key_proto_depIdxs,
		MessageInfos:      file_teleport_access_graph_v1_authorized_key_proto_msgTypes,
	}.Build()
	File_teleport_access_graph_v1_authorized_key_proto = out.File
	file_teleport_access_graph_v1_authorized_key_proto_rawDesc = nil
	file_teleport_access_graph_v1_authorized_key_proto_goTypes = nil
	file_teleport_access_graph_v1_authorized_key_proto_depIdxs = nil
}
