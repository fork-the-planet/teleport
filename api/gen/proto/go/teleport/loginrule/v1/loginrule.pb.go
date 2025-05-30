// Copyright 2022 Gravitational, Inc
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
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: teleport/loginrule/v1/loginrule.proto

package loginrulev1

import (
	types "github.com/gravitational/teleport/api/types"
	wrappers "github.com/gravitational/teleport/api/types/wrappers"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

// LoginRule is a resource to configure rules and logic which should run during
// Teleport user login.
type LoginRule struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Metadata is resource metadata.
	Metadata *types.Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// Version is the resource version.
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	// Priority is the priority of the login rule relative to other login rules
	// in the same cluster. Login rules with a lower numbered priority will be
	// evaluated first.
	Priority int32 `protobuf:"varint,3,opt,name=priority,proto3" json:"priority,omitempty"`
	// TraitsMap is a map of trait keys to lists of predicate expressions which
	// should evaluate to the desired values for that trait.
	TraitsMap map[string]*wrappers.StringValues `protobuf:"bytes,4,rep,name=traits_map,json=traitsMap,proto3" json:"traits_map,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// TraitsExpression is a predicate expression which should return the
	// desired traits for the user upon login.
	TraitsExpression string `protobuf:"bytes,5,opt,name=traits_expression,json=traitsExpression,proto3" json:"traits_expression,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *LoginRule) Reset() {
	*x = LoginRule{}
	mi := &file_teleport_loginrule_v1_loginrule_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginRule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRule) ProtoMessage() {}

func (x *LoginRule) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_loginrule_v1_loginrule_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRule.ProtoReflect.Descriptor instead.
func (*LoginRule) Descriptor() ([]byte, []int) {
	return file_teleport_loginrule_v1_loginrule_proto_rawDescGZIP(), []int{0}
}

func (x *LoginRule) GetMetadata() *types.Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *LoginRule) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *LoginRule) GetPriority() int32 {
	if x != nil {
		return x.Priority
	}
	return 0
}

func (x *LoginRule) GetTraitsMap() map[string]*wrappers.StringValues {
	if x != nil {
		return x.TraitsMap
	}
	return nil
}

func (x *LoginRule) GetTraitsExpression() string {
	if x != nil {
		return x.TraitsExpression
	}
	return ""
}

var File_teleport_loginrule_v1_loginrule_proto protoreflect.FileDescriptor

const file_teleport_loginrule_v1_loginrule_proto_rawDesc = "" +
	"\n" +
	"%teleport/loginrule/v1/loginrule.proto\x12\x15teleport.loginrule.v1\x1a!teleport/legacy/types/types.proto\x1a-teleport/legacy/types/wrappers/wrappers.proto\"\xc1\x02\n" +
	"\tLoginRule\x12+\n" +
	"\bmetadata\x18\x01 \x01(\v2\x0f.types.MetadataR\bmetadata\x12\x18\n" +
	"\aversion\x18\x02 \x01(\tR\aversion\x12\x1a\n" +
	"\bpriority\x18\x03 \x01(\x05R\bpriority\x12N\n" +
	"\n" +
	"traits_map\x18\x04 \x03(\v2/.teleport.loginrule.v1.LoginRule.TraitsMapEntryR\ttraitsMap\x12+\n" +
	"\x11traits_expression\x18\x05 \x01(\tR\x10traitsExpression\x1aT\n" +
	"\x0eTraitsMapEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12,\n" +
	"\x05value\x18\x02 \x01(\v2\x16.wrappers.StringValuesR\x05value:\x028\x01BVZTgithub.com/gravitational/teleport/api/gen/proto/go/teleport/loginrule/v1;loginrulev1b\x06proto3"

var (
	file_teleport_loginrule_v1_loginrule_proto_rawDescOnce sync.Once
	file_teleport_loginrule_v1_loginrule_proto_rawDescData []byte
)

func file_teleport_loginrule_v1_loginrule_proto_rawDescGZIP() []byte {
	file_teleport_loginrule_v1_loginrule_proto_rawDescOnce.Do(func() {
		file_teleport_loginrule_v1_loginrule_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_teleport_loginrule_v1_loginrule_proto_rawDesc), len(file_teleport_loginrule_v1_loginrule_proto_rawDesc)))
	})
	return file_teleport_loginrule_v1_loginrule_proto_rawDescData
}

var file_teleport_loginrule_v1_loginrule_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_teleport_loginrule_v1_loginrule_proto_goTypes = []any{
	(*LoginRule)(nil),             // 0: teleport.loginrule.v1.LoginRule
	nil,                           // 1: teleport.loginrule.v1.LoginRule.TraitsMapEntry
	(*types.Metadata)(nil),        // 2: types.Metadata
	(*wrappers.StringValues)(nil), // 3: wrappers.StringValues
}
var file_teleport_loginrule_v1_loginrule_proto_depIdxs = []int32{
	2, // 0: teleport.loginrule.v1.LoginRule.metadata:type_name -> types.Metadata
	1, // 1: teleport.loginrule.v1.LoginRule.traits_map:type_name -> teleport.loginrule.v1.LoginRule.TraitsMapEntry
	3, // 2: teleport.loginrule.v1.LoginRule.TraitsMapEntry.value:type_name -> wrappers.StringValues
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_teleport_loginrule_v1_loginrule_proto_init() }
func file_teleport_loginrule_v1_loginrule_proto_init() {
	if File_teleport_loginrule_v1_loginrule_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_teleport_loginrule_v1_loginrule_proto_rawDesc), len(file_teleport_loginrule_v1_loginrule_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_teleport_loginrule_v1_loginrule_proto_goTypes,
		DependencyIndexes: file_teleport_loginrule_v1_loginrule_proto_depIdxs,
		MessageInfos:      file_teleport_loginrule_v1_loginrule_proto_msgTypes,
	}.Build()
	File_teleport_loginrule_v1_loginrule_proto = out.File
	file_teleport_loginrule_v1_loginrule_proto_goTypes = nil
	file_teleport_loginrule_v1_loginrule_proto_depIdxs = nil
}
