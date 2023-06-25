// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.2
// source: proto/match_evaluator.proto

package match_evaluator

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

type MatchDetail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ids    []string `protobuf:"bytes,1,rep,name=ids,proto3" json:"ids,omitempty"`      // 组内详情
	Score  int64    `protobuf:"varint,2,opt,name=score,proto3" json:"score,omitempty"` // 组的评分
	GameId string   `protobuf:"bytes,3,opt,name=gameId,proto3" json:"gameId,omitempty"`
}

func (x *MatchDetail) Reset() {
	*x = MatchDetail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_match_evaluator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MatchDetail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MatchDetail) ProtoMessage() {}

func (x *MatchDetail) ProtoReflect() protoreflect.Message {
	mi := &file_proto_match_evaluator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MatchDetail.ProtoReflect.Descriptor instead.
func (*MatchDetail) Descriptor() ([]byte, []int) {
	return file_proto_match_evaluator_proto_rawDescGZIP(), []int{0}
}

func (x *MatchDetail) GetIds() []string {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *MatchDetail) GetScore() int64 {
	if x != nil {
		return x.Score
	}
	return 0
}

func (x *MatchDetail) GetGameId() string {
	if x != nil {
		return x.GameId
	}
	return ""
}

type ToEvalReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Details            []*MatchDetail `protobuf:"bytes,1,rep,name=details,proto3" json:"details,omitempty"`                        // 成组结果
	TaskId             string         `protobuf:"bytes,2,opt,name=taskId,proto3" json:"taskId,omitempty"`                          // 任务ID
	SubTaskId          string         `protobuf:"bytes,3,opt,name=subTaskId,proto3" json:"subTaskId,omitempty"`                    // 子任务ID
	GameId             string         `protobuf:"bytes,4,opt,name=gameId,proto3" json:"gameId,omitempty"`                          // 游戏ID
	SubType            int64          `protobuf:"varint,5,opt,name=subType,proto3" json:"subType,omitempty"`                       // poolID
	EvalGroupId        string         `protobuf:"bytes,6,opt,name=evalGroupId,proto3" json:"evalGroupId,omitempty"`                // 后置优化组的ID
	EvalGroupTaskCount int64          `protobuf:"varint,7,opt,name=evalGroupTaskCount,proto3" json:"evalGroupTaskCount,omitempty"` // 后置优化组内子任务数
	EvalGroupSubId     int64          `protobuf:"varint,8,opt,name=evalGroupSubId,proto3" json:"evalGroupSubId,omitempty"`         // 后置优化组内子任务ID
	Version            int64          `protobuf:"varint,9,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *ToEvalReq) Reset() {
	*x = ToEvalReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_match_evaluator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ToEvalReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ToEvalReq) ProtoMessage() {}

func (x *ToEvalReq) ProtoReflect() protoreflect.Message {
	mi := &file_proto_match_evaluator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ToEvalReq.ProtoReflect.Descriptor instead.
func (*ToEvalReq) Descriptor() ([]byte, []int) {
	return file_proto_match_evaluator_proto_rawDescGZIP(), []int{1}
}

func (x *ToEvalReq) GetDetails() []*MatchDetail {
	if x != nil {
		return x.Details
	}
	return nil
}

func (x *ToEvalReq) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

func (x *ToEvalReq) GetSubTaskId() string {
	if x != nil {
		return x.SubTaskId
	}
	return ""
}

func (x *ToEvalReq) GetGameId() string {
	if x != nil {
		return x.GameId
	}
	return ""
}

func (x *ToEvalReq) GetSubType() int64 {
	if x != nil {
		return x.SubType
	}
	return 0
}

func (x *ToEvalReq) GetEvalGroupId() string {
	if x != nil {
		return x.EvalGroupId
	}
	return ""
}

func (x *ToEvalReq) GetEvalGroupTaskCount() int64 {
	if x != nil {
		return x.EvalGroupTaskCount
	}
	return 0
}

func (x *ToEvalReq) GetEvalGroupSubId() int64 {
	if x != nil {
		return x.EvalGroupSubId
	}
	return 0
}

func (x *ToEvalReq) GetVersion() int64 {
	if x != nil {
		return x.Version
	}
	return 0
}

type ToEvalRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Err  string `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *ToEvalRsp) Reset() {
	*x = ToEvalRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_match_evaluator_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ToEvalRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ToEvalRsp) ProtoMessage() {}

func (x *ToEvalRsp) ProtoReflect() protoreflect.Message {
	mi := &file_proto_match_evaluator_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ToEvalRsp.ProtoReflect.Descriptor instead.
func (*ToEvalRsp) Descriptor() ([]byte, []int) {
	return file_proto_match_evaluator_proto_rawDescGZIP(), []int{2}
}

func (x *ToEvalRsp) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *ToEvalRsp) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

var File_proto_match_evaluator_proto protoreflect.FileDescriptor

var file_proto_match_evaluator_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x65, 0x76,
	0x61, 0x6c, 0x75, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x6d,
	0x61, 0x74, 0x63, 0x68, 0x5f, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x6f, 0x72, 0x22, 0x4d,
	0x0a, 0x0b, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x10, 0x0a,
	0x03, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x69, 0x64, 0x73, 0x12,
	0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05,
	0x73, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x61, 0x6d, 0x65, 0x49, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x67, 0x61, 0x6d, 0x65, 0x49, 0x64, 0x22, 0xbf, 0x02,
	0x0a, 0x09, 0x54, 0x6f, 0x45, 0x76, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x12, 0x36, 0x0a, 0x07, 0x64,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6d,
	0x61, 0x74, 0x63, 0x68, 0x5f, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x4d,
	0x61, 0x74, 0x63, 0x68, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x07, 0x64, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x73,
	0x75, 0x62, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x73, 0x75, 0x62, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x61, 0x6d,
	0x65, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x67, 0x61, 0x6d, 0x65, 0x49,
	0x64, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x62, 0x54, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x07, 0x73, 0x75, 0x62, 0x54, 0x79, 0x70, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x65,
	0x76, 0x61, 0x6c, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x65, 0x76, 0x61, 0x6c, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x12, 0x2e, 0x0a,
	0x12, 0x65, 0x76, 0x61, 0x6c, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x12, 0x65, 0x76, 0x61, 0x6c, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x54, 0x61, 0x73, 0x6b, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x26, 0x0a,
	0x0e, 0x65, 0x76, 0x61, 0x6c, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x53, 0x75, 0x62, 0x49, 0x64, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x65, 0x76, 0x61, 0x6c, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x53, 0x75, 0x62, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22,
	0x31, 0x0a, 0x09, 0x54, 0x6f, 0x45, 0x76, 0x61, 0x6c, 0x52, 0x73, 0x70, 0x12, 0x12, 0x0a, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65,
	0x72, 0x72, 0x32, 0x55, 0x0a, 0x0f, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x65, 0x76, 0x61, 0x6c,
	0x75, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x42, 0x0a, 0x06, 0x54, 0x6f, 0x45, 0x76, 0x61, 0x6c, 0x12,
	0x1a, 0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x6f,
	0x72, 0x2e, 0x54, 0x6f, 0x45, 0x76, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x1a, 0x1a, 0x2e, 0x6d, 0x61,
	0x74, 0x63, 0x68, 0x5f, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x54, 0x6f,
	0x45, 0x76, 0x61, 0x6c, 0x52, 0x73, 0x70, 0x22, 0x00, 0x42, 0x19, 0x5a, 0x17, 0x2e, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x65, 0x76, 0x61, 0x6c, 0x75,
	0x61, 0x74, 0x6f, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_match_evaluator_proto_rawDescOnce sync.Once
	file_proto_match_evaluator_proto_rawDescData = file_proto_match_evaluator_proto_rawDesc
)

func file_proto_match_evaluator_proto_rawDescGZIP() []byte {
	file_proto_match_evaluator_proto_rawDescOnce.Do(func() {
		file_proto_match_evaluator_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_match_evaluator_proto_rawDescData)
	})
	return file_proto_match_evaluator_proto_rawDescData
}

var file_proto_match_evaluator_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_match_evaluator_proto_goTypes = []interface{}{
	(*MatchDetail)(nil), // 0: match_evaluator.MatchDetail
	(*ToEvalReq)(nil),   // 1: match_evaluator.ToEvalReq
	(*ToEvalRsp)(nil),   // 2: match_evaluator.ToEvalRsp
}
var file_proto_match_evaluator_proto_depIdxs = []int32{
	0, // 0: match_evaluator.ToEvalReq.details:type_name -> match_evaluator.MatchDetail
	1, // 1: match_evaluator.Match_evaluator.ToEval:input_type -> match_evaluator.ToEvalReq
	2, // 2: match_evaluator.Match_evaluator.ToEval:output_type -> match_evaluator.ToEvalRsp
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_match_evaluator_proto_init() }
func file_proto_match_evaluator_proto_init() {
	if File_proto_match_evaluator_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_match_evaluator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MatchDetail); i {
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
		file_proto_match_evaluator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ToEvalReq); i {
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
		file_proto_match_evaluator_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ToEvalRsp); i {
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
			RawDescriptor: file_proto_match_evaluator_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_match_evaluator_proto_goTypes,
		DependencyIndexes: file_proto_match_evaluator_proto_depIdxs,
		MessageInfos:      file_proto_match_evaluator_proto_msgTypes,
	}.Build()
	File_proto_match_evaluator_proto = out.File
	file_proto_match_evaluator_proto_rawDesc = nil
	file_proto_match_evaluator_proto_goTypes = nil
	file_proto_match_evaluator_proto_depIdxs = nil
}
