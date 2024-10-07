// Copyright 2023 Gravitational, Inc
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
// source: teleport/kube/v1/kube_service.proto

package kubev1

import (
	types "github.com/gravitational/teleport/api/types"
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

// ListKubernetesResourcesRequest defines a request to retrieve resources paginated. Only
// one type of resource can be retrieved per request.
type ListKubernetesResourcesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ResourceType is the Kubernetes resource that is going to be retrieved.
	ResourceType string `protobuf:"bytes,1,opt,name=resource_type,json=resourceType,proto3" json:"resource_type,omitempty"`
	// Limit is the maximum amount of resources to retrieve.
	Limit int32 `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	// StartKey is used to start listing resources from a specific spot. It
	// should be set to the previous NextKey value if using pagination, or
	// left empty.
	StartKey string `protobuf:"bytes,3,opt,name=start_key,json=startKey,proto3" json:"start_key,omitempty"`
	// Labels is a label-based matcher if non-empty.
	Labels map[string]string `protobuf:"bytes,4,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// PredicateExpression defines boolean conditions that will be matched against the resource.
	PredicateExpression string `protobuf:"bytes,5,opt,name=predicate_expression,json=predicateExpression,proto3" json:"predicate_expression,omitempty"`
	// SearchKeywords is a list of search keywords to match against resource field values.
	SearchKeywords []string `protobuf:"bytes,6,rep,name=search_keywords,json=searchKeywords,proto3" json:"search_keywords,omitempty"`
	// SortBy describes which resource field and which direction to sort by.
	SortBy *types.SortBy `protobuf:"bytes,7,opt,name=sort_by,json=sortBy,proto3" json:"sort_by,omitempty"`
	// NeedTotalCount indicates whether or not the caller also wants the total number of resources
	// after filtering.
	NeedTotalCount bool `protobuf:"varint,8,opt,name=need_total_count,json=needTotalCount,proto3" json:"need_total_count,omitempty"`
	// UseSearchAsRoles indicates that the response should include all resources
	// the caller is able to request access to using search_as_roles
	UseSearchAsRoles bool `protobuf:"varint,9,opt,name=use_search_as_roles,json=useSearchAsRoles,proto3" json:"use_search_as_roles,omitempty"`
	// UsePreviewAsRoles indicates that the response should include all resources
	// the caller would be able to access with their preview_as_roles
	UsePreviewAsRoles bool `protobuf:"varint,11,opt,name=use_preview_as_roles,json=usePreviewAsRoles,proto3" json:"use_preview_as_roles,omitempty"`
	// TeleportCluster is the Teleport Cluster name to route the request to.
	TeleportCluster string `protobuf:"bytes,12,opt,name=teleport_cluster,json=teleportCluster,proto3" json:"teleport_cluster,omitempty"`
	// Cluster is the Kubernetes Cluster to request the resources.
	KubernetesCluster string `protobuf:"bytes,13,opt,name=kubernetes_cluster,json=kubernetesCluster,proto3" json:"kubernetes_cluster,omitempty"`
	// Namespace is the Kubernetes namespace where the resources must be located.
	// To search on every Kubernetes Namespace, do not define the value.
	KubernetesNamespace string `protobuf:"bytes,14,opt,name=kubernetes_namespace,json=kubernetesNamespace,proto3" json:"kubernetes_namespace,omitempty"`
}

func (x *ListKubernetesResourcesRequest) Reset() {
	*x = ListKubernetesResourcesRequest{}
	mi := &file_teleport_kube_v1_kube_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListKubernetesResourcesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListKubernetesResourcesRequest) ProtoMessage() {}

func (x *ListKubernetesResourcesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_kube_v1_kube_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListKubernetesResourcesRequest.ProtoReflect.Descriptor instead.
func (*ListKubernetesResourcesRequest) Descriptor() ([]byte, []int) {
	return file_teleport_kube_v1_kube_service_proto_rawDescGZIP(), []int{0}
}

func (x *ListKubernetesResourcesRequest) GetResourceType() string {
	if x != nil {
		return x.ResourceType
	}
	return ""
}

func (x *ListKubernetesResourcesRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *ListKubernetesResourcesRequest) GetStartKey() string {
	if x != nil {
		return x.StartKey
	}
	return ""
}

func (x *ListKubernetesResourcesRequest) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *ListKubernetesResourcesRequest) GetPredicateExpression() string {
	if x != nil {
		return x.PredicateExpression
	}
	return ""
}

func (x *ListKubernetesResourcesRequest) GetSearchKeywords() []string {
	if x != nil {
		return x.SearchKeywords
	}
	return nil
}

func (x *ListKubernetesResourcesRequest) GetSortBy() *types.SortBy {
	if x != nil {
		return x.SortBy
	}
	return nil
}

func (x *ListKubernetesResourcesRequest) GetNeedTotalCount() bool {
	if x != nil {
		return x.NeedTotalCount
	}
	return false
}

func (x *ListKubernetesResourcesRequest) GetUseSearchAsRoles() bool {
	if x != nil {
		return x.UseSearchAsRoles
	}
	return false
}

func (x *ListKubernetesResourcesRequest) GetUsePreviewAsRoles() bool {
	if x != nil {
		return x.UsePreviewAsRoles
	}
	return false
}

func (x *ListKubernetesResourcesRequest) GetTeleportCluster() string {
	if x != nil {
		return x.TeleportCluster
	}
	return ""
}

func (x *ListKubernetesResourcesRequest) GetKubernetesCluster() string {
	if x != nil {
		return x.KubernetesCluster
	}
	return ""
}

func (x *ListKubernetesResourcesRequest) GetKubernetesNamespace() string {
	if x != nil {
		return x.KubernetesNamespace
	}
	return ""
}

// ListKubernetesResourcesResponse is the response of ListKubernetesResources method.
type ListKubernetesResourcesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Resources is a list of resource.
	Resources []*types.KubernetesResourceV1 `protobuf:"bytes,1,rep,name=resources,proto3" json:"resources,omitempty"`
	// NextKey is the next Key to use as StartKey in a ListResourcesRequest to
	// continue retrieving pages of resource. If NextKey is empty, there are no
	// more pages.
	NextKey string `protobuf:"bytes,2,opt,name=next_key,json=nextKey,proto3" json:"next_key,omitempty"`
	// TotalCount is the total number of resources available after filter, if any.
	TotalCount int32 `protobuf:"varint,3,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
}

func (x *ListKubernetesResourcesResponse) Reset() {
	*x = ListKubernetesResourcesResponse{}
	mi := &file_teleport_kube_v1_kube_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListKubernetesResourcesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListKubernetesResourcesResponse) ProtoMessage() {}

func (x *ListKubernetesResourcesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_teleport_kube_v1_kube_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListKubernetesResourcesResponse.ProtoReflect.Descriptor instead.
func (*ListKubernetesResourcesResponse) Descriptor() ([]byte, []int) {
	return file_teleport_kube_v1_kube_service_proto_rawDescGZIP(), []int{1}
}

func (x *ListKubernetesResourcesResponse) GetResources() []*types.KubernetesResourceV1 {
	if x != nil {
		return x.Resources
	}
	return nil
}

func (x *ListKubernetesResourcesResponse) GetNextKey() string {
	if x != nil {
		return x.NextKey
	}
	return ""
}

func (x *ListKubernetesResourcesResponse) GetTotalCount() int32 {
	if x != nil {
		return x.TotalCount
	}
	return 0
}

var File_teleport_kube_v1_kube_service_proto protoreflect.FileDescriptor

var file_teleport_kube_v1_kube_service_proto_rawDesc = []byte{
	0x0a, 0x23, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x6b, 0x75, 0x62, 0x65, 0x2f,
	0x76, 0x31, 0x2f, 0x6b, 0x75, 0x62, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e,
	0x6b, 0x75, 0x62, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x21, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x2f, 0x6c, 0x65, 0x67, 0x61, 0x63, 0x79, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa4, 0x05, 0x0a, 0x1e, 0x4c,
	0x69, 0x73, 0x74, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a,
	0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x54, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18,
	0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x3c, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x2e, 0x6b, 0x75, 0x62, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4b, 0x75, 0x62,
	0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x31, 0x0a, 0x14, 0x70,
	0x72, 0x65, 0x64, 0x69, 0x63, 0x61, 0x74, 0x65, 0x5f, 0x65, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x70, 0x72, 0x65, 0x64, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x45, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x27,
	0x0a, 0x0f, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x5f, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64,
	0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x4b,
	0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x26, 0x0a, 0x07, 0x73, 0x6f, 0x72, 0x74, 0x5f,
	0x62, 0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x2e, 0x53, 0x6f, 0x72, 0x74, 0x42, 0x79, 0x52, 0x06, 0x73, 0x6f, 0x72, 0x74, 0x42, 0x79, 0x12,
	0x28, 0x0a, 0x10, 0x6e, 0x65, 0x65, 0x64, 0x5f, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x6e, 0x65, 0x65, 0x64, 0x54,
	0x6f, 0x74, 0x61, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2d, 0x0a, 0x13, 0x75, 0x73, 0x65,
	0x5f, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x5f, 0x61, 0x73, 0x5f, 0x72, 0x6f, 0x6c, 0x65, 0x73,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x08, 0x52, 0x10, 0x75, 0x73, 0x65, 0x53, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x41, 0x73, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x12, 0x2f, 0x0a, 0x14, 0x75, 0x73, 0x65, 0x5f,
	0x70, 0x72, 0x65, 0x76, 0x69, 0x65, 0x77, 0x5f, 0x61, 0x73, 0x5f, 0x72, 0x6f, 0x6c, 0x65, 0x73,
	0x18, 0x0b, 0x20, 0x01, 0x28, 0x08, 0x52, 0x11, 0x75, 0x73, 0x65, 0x50, 0x72, 0x65, 0x76, 0x69,
	0x65, 0x77, 0x41, 0x73, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x12, 0x29, 0x0a, 0x10, 0x74, 0x65, 0x6c,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0f, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x43, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x12, 0x2d, 0x0a, 0x12, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74,
	0x65, 0x73, 0x5f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x11, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x43, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x12, 0x31, 0x0a, 0x14, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65,
	0x73, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x13, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x4e, 0x61, 0x6d,
	0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x22, 0x98, 0x01, 0x0a, 0x1f, 0x4c, 0x69, 0x73, 0x74, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e,
	0x65, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x2e, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x56, 0x31, 0x52, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x12, 0x19, 0x0a, 0x08, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x6e, 0x65, 0x78, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x32, 0x8d, 0x01, 0x0a,
	0x0b, 0x4b, 0x75, 0x62, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x7e, 0x0a, 0x17,
	0x4c, 0x69, 0x73, 0x74, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x52, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x12, 0x30, 0x2e, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4b,
	0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x74, 0x65, 0x6c, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x4b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x4c, 0x5a, 0x4a,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x72, 0x61, 0x76, 0x69,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x67, 0x6f, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x6b, 0x75, 0x62, 0x65,
	0x2f, 0x76, 0x31, 0x3b, 0x6b, 0x75, 0x62, 0x65, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_teleport_kube_v1_kube_service_proto_rawDescOnce sync.Once
	file_teleport_kube_v1_kube_service_proto_rawDescData = file_teleport_kube_v1_kube_service_proto_rawDesc
)

func file_teleport_kube_v1_kube_service_proto_rawDescGZIP() []byte {
	file_teleport_kube_v1_kube_service_proto_rawDescOnce.Do(func() {
		file_teleport_kube_v1_kube_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_teleport_kube_v1_kube_service_proto_rawDescData)
	})
	return file_teleport_kube_v1_kube_service_proto_rawDescData
}

var file_teleport_kube_v1_kube_service_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_teleport_kube_v1_kube_service_proto_goTypes = []any{
	(*ListKubernetesResourcesRequest)(nil),  // 0: teleport.kube.v1.ListKubernetesResourcesRequest
	(*ListKubernetesResourcesResponse)(nil), // 1: teleport.kube.v1.ListKubernetesResourcesResponse
	nil,                                     // 2: teleport.kube.v1.ListKubernetesResourcesRequest.LabelsEntry
	(*types.SortBy)(nil),                    // 3: types.SortBy
	(*types.KubernetesResourceV1)(nil),      // 4: types.KubernetesResourceV1
}
var file_teleport_kube_v1_kube_service_proto_depIdxs = []int32{
	2, // 0: teleport.kube.v1.ListKubernetesResourcesRequest.labels:type_name -> teleport.kube.v1.ListKubernetesResourcesRequest.LabelsEntry
	3, // 1: teleport.kube.v1.ListKubernetesResourcesRequest.sort_by:type_name -> types.SortBy
	4, // 2: teleport.kube.v1.ListKubernetesResourcesResponse.resources:type_name -> types.KubernetesResourceV1
	0, // 3: teleport.kube.v1.KubeService.ListKubernetesResources:input_type -> teleport.kube.v1.ListKubernetesResourcesRequest
	1, // 4: teleport.kube.v1.KubeService.ListKubernetesResources:output_type -> teleport.kube.v1.ListKubernetesResourcesResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_teleport_kube_v1_kube_service_proto_init() }
func file_teleport_kube_v1_kube_service_proto_init() {
	if File_teleport_kube_v1_kube_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_teleport_kube_v1_kube_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_teleport_kube_v1_kube_service_proto_goTypes,
		DependencyIndexes: file_teleport_kube_v1_kube_service_proto_depIdxs,
		MessageInfos:      file_teleport_kube_v1_kube_service_proto_msgTypes,
	}.Build()
	File_teleport_kube_v1_kube_service_proto = out.File
	file_teleport_kube_v1_kube_service_proto_rawDesc = nil
	file_teleport_kube_v1_kube_service_proto_goTypes = nil
	file_teleport_kube_v1_kube_service_proto_depIdxs = nil
}
