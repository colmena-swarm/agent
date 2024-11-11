/*
 *  Copyright 2002-2025 Barcelona Supercomputing Center (www.bsc.es)
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package policy

import (
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

type InitializeOrStopRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PolicyName    string                 `protobuf:"bytes,1,opt,name=policy_name,json=policyName,proto3" json:"policy_name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InitializeOrStopRequest) Reset() {
	*x = InitializeOrStopRequest{}
	mi := &file_policy_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InitializeOrStopRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitializeOrStopRequest) ProtoMessage() {}

func (x *InitializeOrStopRequest) ProtoReflect() protoreflect.Message {
	mi := &file_policy_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitializeOrStopRequest.ProtoReflect.Descriptor instead.
func (*InitializeOrStopRequest) Descriptor() ([]byte, []int) {
	return file_policy_proto_rawDescGZIP(), []int{0}
}

func (x *InitializeOrStopRequest) GetPolicyName() string {
	if x != nil {
		return x.PolicyName
	}
	return ""
}

type InitializeOrStopResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InitializeOrStopResponse) Reset() {
	*x = InitializeOrStopResponse{}
	mi := &file_policy_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InitializeOrStopResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitializeOrStopResponse) ProtoMessage() {}

func (x *InitializeOrStopResponse) ProtoReflect() protoreflect.Message {
	mi := &file_policy_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitializeOrStopResponse.ProtoReflect.Descriptor instead.
func (*InitializeOrStopResponse) Descriptor() ([]byte, []int) {
	return file_policy_proto_rawDescGZIP(), []int{1}
}

func (x *InitializeOrStopResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *InitializeOrStopResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type DecideRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Roles         []*Role                `protobuf:"bytes,1,rep,name=roles,proto3" json:"roles,omitempty"`
	Levels        []*IndicatorLevel      `protobuf:"bytes,2,rep,name=levels,proto3" json:"levels,omitempty"`
	Resources     []*Resource            `protobuf:"bytes,3,rep,name=resources,proto3" json:"resources,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DecideRequest) Reset() {
	*x = DecideRequest{}
	mi := &file_policy_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DecideRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DecideRequest) ProtoMessage() {}

func (x *DecideRequest) ProtoReflect() protoreflect.Message {
	mi := &file_policy_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DecideRequest.ProtoReflect.Descriptor instead.
func (*DecideRequest) Descriptor() ([]byte, []int) {
	return file_policy_proto_rawDescGZIP(), []int{2}
}

func (x *DecideRequest) GetRoles() []*Role {
	if x != nil {
		return x.Roles
	}
	return nil
}

func (x *DecideRequest) GetLevels() []*IndicatorLevel {
	if x != nil {
		return x.Levels
	}
	return nil
}

func (x *DecideRequest) GetResources() []*Resource {
	if x != nil {
		return x.Resources
	}
	return nil
}

type DecideResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Decisions     map[string]bool        `protobuf:"bytes,1,rep,name=decisions,proto3" json:"decisions,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DecideResponse) Reset() {
	*x = DecideResponse{}
	mi := &file_policy_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DecideResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DecideResponse) ProtoMessage() {}

func (x *DecideResponse) ProtoReflect() protoreflect.Message {
	mi := &file_policy_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DecideResponse.ProtoReflect.Descriptor instead.
func (*DecideResponse) Descriptor() ([]byte, []int) {
	return file_policy_proto_rawDescGZIP(), []int{3}
}

func (x *DecideResponse) GetDecisions() map[string]bool {
	if x != nil {
		return x.Decisions
	}
	return nil
}

type IndicatorLevel struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	Name           string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value          float32                `protobuf:"fixed32,2,opt,name=value,proto3" json:"value,omitempty"`
	IsMet          bool                   `protobuf:"varint,3,opt,name=isMet,proto3" json:"isMet,omitempty"`
	Threshold      float32                `protobuf:"fixed32,4,opt,name=threshold,proto3" json:"threshold,omitempty"`
	AssociatedRole string                 `protobuf:"bytes,5,opt,name=associatedRole,proto3" json:"associatedRole,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *IndicatorLevel) Reset() {
	*x = IndicatorLevel{}
	mi := &file_policy_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IndicatorLevel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndicatorLevel) ProtoMessage() {}

func (x *IndicatorLevel) ProtoReflect() protoreflect.Message {
	mi := &file_policy_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndicatorLevel.ProtoReflect.Descriptor instead.
func (*IndicatorLevel) Descriptor() ([]byte, []int) {
	return file_policy_proto_rawDescGZIP(), []int{4}
}

func (x *IndicatorLevel) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *IndicatorLevel) GetValue() float32 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *IndicatorLevel) GetIsMet() bool {
	if x != nil {
		return x.IsMet
	}
	return false
}

func (x *IndicatorLevel) GetThreshold() float32 {
	if x != nil {
		return x.Threshold
	}
	return 0
}

func (x *IndicatorLevel) GetAssociatedRole() string {
	if x != nil {
		return x.AssociatedRole
	}
	return ""
}

type Role struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RoleName      string                 `protobuf:"bytes,1,opt,name=roleName,proto3" json:"roleName,omitempty"`
	IsRunning     bool                   `protobuf:"varint,2,opt,name=isRunning,proto3" json:"isRunning,omitempty"`
	Resources     []*Resource            `protobuf:"bytes,3,rep,name=resources,proto3" json:"resources,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Role) Reset() {
	*x = Role{}
	mi := &file_policy_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Role) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Role) ProtoMessage() {}

func (x *Role) ProtoReflect() protoreflect.Message {
	mi := &file_policy_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Role.ProtoReflect.Descriptor instead.
func (*Role) Descriptor() ([]byte, []int) {
	return file_policy_proto_rawDescGZIP(), []int{5}
}

func (x *Role) GetRoleName() string {
	if x != nil {
		return x.RoleName
	}
	return ""
}

func (x *Role) GetIsRunning() bool {
	if x != nil {
		return x.IsRunning
	}
	return false
}

func (x *Role) GetResources() []*Resource {
	if x != nil {
		return x.Resources
	}
	return nil
}

type Resource struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value         float32                `protobuf:"fixed32,2,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Resource) Reset() {
	*x = Resource{}
	mi := &file_policy_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Resource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Resource) ProtoMessage() {}

func (x *Resource) ProtoReflect() protoreflect.Message {
	mi := &file_policy_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Resource.ProtoReflect.Descriptor instead.
func (*Resource) Descriptor() ([]byte, []int) {
	return file_policy_proto_rawDescGZIP(), []int{6}
}

func (x *Resource) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Resource) GetValue() float32 {
	if x != nil {
		return x.Value
	}
	return 0
}

var File_policy_proto protoreflect.FileDescriptor

const file_policy_proto_rawDesc = "" +
	"\n" +
	"\fpolicy.proto\x12\x06policy\":\n" +
	"\x17InitializeOrStopRequest\x12\x1f\n" +
	"\vpolicy_name\x18\x01 \x01(\tR\n" +
	"policyName\"N\n" +
	"\x18InitializeOrStopResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\"\x93\x01\n" +
	"\rDecideRequest\x12\"\n" +
	"\x05roles\x18\x01 \x03(\v2\f.policy.RoleR\x05roles\x12.\n" +
	"\x06levels\x18\x02 \x03(\v2\x16.policy.IndicatorLevelR\x06levels\x12.\n" +
	"\tresources\x18\x03 \x03(\v2\x10.policy.ResourceR\tresources\"\x93\x01\n" +
	"\x0eDecideResponse\x12C\n" +
	"\tdecisions\x18\x01 \x03(\v2%.policy.DecideResponse.DecisionsEntryR\tdecisions\x1a<\n" +
	"\x0eDecisionsEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\bR\x05value:\x028\x01\"\x96\x01\n" +
	"\x0eIndicatorLevel\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12\x14\n" +
	"\x05value\x18\x02 \x01(\x02R\x05value\x12\x14\n" +
	"\x05isMet\x18\x03 \x01(\bR\x05isMet\x12\x1c\n" +
	"\tthreshold\x18\x04 \x01(\x02R\tthreshold\x12&\n" +
	"\x0eassociatedRole\x18\x05 \x01(\tR\x0eassociatedRole\"p\n" +
	"\x04Role\x12\x1a\n" +
	"\broleName\x18\x01 \x01(\tR\broleName\x12\x1c\n" +
	"\tisRunning\x18\x02 \x01(\bR\tisRunning\x12.\n" +
	"\tresources\x18\x03 \x03(\v2\x10.policy.ResourceR\tresources\"4\n" +
	"\bResource\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12\x14\n" +
	"\x05value\x18\x02 \x01(\x02R\x05value2\xe4\x01\n" +
	"\rPolicyService\x12O\n" +
	"\n" +
	"Initialize\x12\x1f.policy.InitializeOrStopRequest\x1a .policy.InitializeOrStopResponse\x127\n" +
	"\x06Decide\x12\x15.policy.DecideRequest\x1a\x16.policy.DecideResponse\x12I\n" +
	"\x04Stop\x12\x1f.policy.InitializeOrStopRequest\x1a .policy.InitializeOrStopResponseB\n" +
	"Z\b./policyb\x06proto3"

var (
	file_policy_proto_rawDescOnce sync.Once
	file_policy_proto_rawDescData []byte
)

func file_policy_proto_rawDescGZIP() []byte {
	file_policy_proto_rawDescOnce.Do(func() {
		file_policy_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_policy_proto_rawDesc), len(file_policy_proto_rawDesc)))
	})
	return file_policy_proto_rawDescData
}

var file_policy_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_policy_proto_goTypes = []any{
	(*InitializeOrStopRequest)(nil),  // 0: policy.InitializeOrStopRequest
	(*InitializeOrStopResponse)(nil), // 1: policy.InitializeOrStopResponse
	(*DecideRequest)(nil),            // 2: policy.DecideRequest
	(*DecideResponse)(nil),           // 3: policy.DecideResponse
	(*IndicatorLevel)(nil),           // 4: policy.IndicatorLevel
	(*Role)(nil),                     // 5: policy.Role
	(*Resource)(nil),                 // 6: policy.Resource
	nil,                              // 7: policy.DecideResponse.DecisionsEntry
}
var file_policy_proto_depIdxs = []int32{
	5, // 0: policy.DecideRequest.roles:type_name -> policy.Role
	4, // 1: policy.DecideRequest.levels:type_name -> policy.IndicatorLevel
	6, // 2: policy.DecideRequest.resources:type_name -> policy.Resource
	7, // 3: policy.DecideResponse.decisions:type_name -> policy.DecideResponse.DecisionsEntry
	6, // 4: policy.Role.resources:type_name -> policy.Resource
	0, // 5: policy.PolicyService.Initialize:input_type -> policy.InitializeOrStopRequest
	2, // 6: policy.PolicyService.Decide:input_type -> policy.DecideRequest
	0, // 7: policy.PolicyService.Stop:input_type -> policy.InitializeOrStopRequest
	1, // 8: policy.PolicyService.Initialize:output_type -> policy.InitializeOrStopResponse
	3, // 9: policy.PolicyService.Decide:output_type -> policy.DecideResponse
	1, // 10: policy.PolicyService.Stop:output_type -> policy.InitializeOrStopResponse
	8, // [8:11] is the sub-list for method output_type
	5, // [5:8] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_policy_proto_init() }
func file_policy_proto_init() {
	if File_policy_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_policy_proto_rawDesc), len(file_policy_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_policy_proto_goTypes,
		DependencyIndexes: file_policy_proto_depIdxs,
		MessageInfos:      file_policy_proto_msgTypes,
	}.Build()
	File_policy_proto = out.File
	file_policy_proto_goTypes = nil
	file_policy_proto_depIdxs = nil
}
