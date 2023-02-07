// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.4
// source: pb/gold.proto

package pbgo

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

type Gold struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// {"comment": "玩家 ID", "type": "BIGINT", "size": 20, "default": "0", "primary_key": "default"}
	Uid int64 `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	// {"comment": "遊戲幣", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Gold_0 int64 `protobuf:"varint,2,opt,name=gold_0,json=gold0,proto3" json:"gold_0,omitempty"`
	// {"comment": "綠寶石", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Gold_1 int64 `protobuf:"varint,3,opt,name=gold_1,json=gold1,proto3" json:"gold_1,omitempty"`
	// {"comment": "普通鑰匙", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Gold_2 int64 `protobuf:"varint,4,opt,name=gold_2,json=gold2,proto3" json:"gold_2,omitempty"`
	// {"comment": "黑曜石鑰匙", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Gold_3 int64 `protobuf:"varint,5,opt,name=gold_3,json=gold3,proto3" json:"gold_3,omitempty"`
	// {"comment": "保留", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Gold_4 int64 `protobuf:"varint,6,opt,name=gold_4,json=gold4,proto3" json:"gold_4,omitempty"`
	// {"comment": "保留", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Gold_5 int64 `protobuf:"varint,7,opt,name=gold_5,json=gold5,proto3" json:"gold_5,omitempty"`
	// {"comment": "保留", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Gold_6 int64 `protobuf:"varint,8,opt,name=gold_6,json=gold6,proto3" json:"gold_6,omitempty"`
	// {"comment": "保留", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Gold_7 int64 `protobuf:"varint,9,opt,name=gold_7,json=gold7,proto3" json:"gold_7,omitempty"`
	// {"comment": "保留", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Gold_8 int64 `protobuf:"varint,10,opt,name=gold_8,json=gold8,proto3" json:"gold_8,omitempty"`
	// {"comment": "保留", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Gold_9 int64 `protobuf:"varint,11,opt,name=gold_9,json=gold9,proto3" json:"gold_9,omitempty"`
	// {"comment": "保留", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Bonus_0 int64 `protobuf:"varint,12,opt,name=bonus_0,json=bonus0,proto3" json:"bonus_0,omitempty"`
	// {"comment": "保留", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Bonus_1 int64 `protobuf:"varint,13,opt,name=bonus_1,json=bonus1,proto3" json:"bonus_1,omitempty"`
	// {"comment": "保留", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Bonus_2 int64 `protobuf:"varint,14,opt,name=bonus_2,json=bonus2,proto3" json:"bonus_2,omitempty"`
	// {"comment": "保留", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Bonus_3 int64 `protobuf:"varint,15,opt,name=bonus_3,json=bonus3,proto3" json:"bonus_3,omitempty"`
	// {"comment": "保留", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Bonus_4 int64 `protobuf:"varint,16,opt,name=bonus_4,json=bonus4,proto3" json:"bonus_4,omitempty"`
	// {"comment": "保留", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Bonus_5 int64 `protobuf:"varint,17,opt,name=bonus_5,json=bonus5,proto3" json:"bonus_5,omitempty"`
	// {"comment": "保留", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Bonus_6 int64 `protobuf:"varint,18,opt,name=bonus_6,json=bonus6,proto3" json:"bonus_6,omitempty"`
	// {"comment": "保留", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Bonus_7 int64 `protobuf:"varint,19,opt,name=bonus_7,json=bonus7,proto3" json:"bonus_7,omitempty"`
	// {"comment": "保留", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Bonus_8 int64 `protobuf:"varint,20,opt,name=bonus_8,json=bonus8,proto3" json:"bonus_8,omitempty"`
	// {"comment": "保留", "type": "BIGINT", "size": 20, "default": "0", "can_null": "false"}
	Bonus_9 int64 `protobuf:"varint,21,opt,name=bonus_9,json=bonus9,proto3" json:"bonus_9,omitempty"`
}

func (x *Gold) Reset() {
	*x = Gold{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_gold_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Gold) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Gold) ProtoMessage() {}

func (x *Gold) ProtoReflect() protoreflect.Message {
	mi := &file_pb_gold_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Gold.ProtoReflect.Descriptor instead.
func (*Gold) Descriptor() ([]byte, []int) {
	return file_pb_gold_proto_rawDescGZIP(), []int{0}
}

func (x *Gold) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *Gold) GetGold_0() int64 {
	if x != nil {
		return x.Gold_0
	}
	return 0
}

func (x *Gold) GetGold_1() int64 {
	if x != nil {
		return x.Gold_1
	}
	return 0
}

func (x *Gold) GetGold_2() int64 {
	if x != nil {
		return x.Gold_2
	}
	return 0
}

func (x *Gold) GetGold_3() int64 {
	if x != nil {
		return x.Gold_3
	}
	return 0
}

func (x *Gold) GetGold_4() int64 {
	if x != nil {
		return x.Gold_4
	}
	return 0
}

func (x *Gold) GetGold_5() int64 {
	if x != nil {
		return x.Gold_5
	}
	return 0
}

func (x *Gold) GetGold_6() int64 {
	if x != nil {
		return x.Gold_6
	}
	return 0
}

func (x *Gold) GetGold_7() int64 {
	if x != nil {
		return x.Gold_7
	}
	return 0
}

func (x *Gold) GetGold_8() int64 {
	if x != nil {
		return x.Gold_8
	}
	return 0
}

func (x *Gold) GetGold_9() int64 {
	if x != nil {
		return x.Gold_9
	}
	return 0
}

func (x *Gold) GetBonus_0() int64 {
	if x != nil {
		return x.Bonus_0
	}
	return 0
}

func (x *Gold) GetBonus_1() int64 {
	if x != nil {
		return x.Bonus_1
	}
	return 0
}

func (x *Gold) GetBonus_2() int64 {
	if x != nil {
		return x.Bonus_2
	}
	return 0
}

func (x *Gold) GetBonus_3() int64 {
	if x != nil {
		return x.Bonus_3
	}
	return 0
}

func (x *Gold) GetBonus_4() int64 {
	if x != nil {
		return x.Bonus_4
	}
	return 0
}

func (x *Gold) GetBonus_5() int64 {
	if x != nil {
		return x.Bonus_5
	}
	return 0
}

func (x *Gold) GetBonus_6() int64 {
	if x != nil {
		return x.Bonus_6
	}
	return 0
}

func (x *Gold) GetBonus_7() int64 {
	if x != nil {
		return x.Bonus_7
	}
	return 0
}

func (x *Gold) GetBonus_8() int64 {
	if x != nil {
		return x.Bonus_8
	}
	return 0
}

func (x *Gold) GetBonus_9() int64 {
	if x != nil {
		return x.Bonus_9
	}
	return 0
}

var File_pb_gold_proto protoreflect.FileDescriptor

var file_pb_gold_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x70, 0x62, 0x2f, 0x67, 0x6f, 0x6c, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xf8, 0x03, 0x0a, 0x04, 0x67, 0x6f, 0x6c, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x15, 0x0a, 0x06, 0x67, 0x6f,
	0x6c, 0x64, 0x5f, 0x30, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x67, 0x6f, 0x6c, 0x64,
	0x30, 0x12, 0x15, 0x0a, 0x06, 0x67, 0x6f, 0x6c, 0x64, 0x5f, 0x31, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x05, 0x67, 0x6f, 0x6c, 0x64, 0x31, 0x12, 0x15, 0x0a, 0x06, 0x67, 0x6f, 0x6c, 0x64,
	0x5f, 0x32, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x67, 0x6f, 0x6c, 0x64, 0x32, 0x12,
	0x15, 0x0a, 0x06, 0x67, 0x6f, 0x6c, 0x64, 0x5f, 0x33, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x05, 0x67, 0x6f, 0x6c, 0x64, 0x33, 0x12, 0x15, 0x0a, 0x06, 0x67, 0x6f, 0x6c, 0x64, 0x5f, 0x34,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x67, 0x6f, 0x6c, 0x64, 0x34, 0x12, 0x15, 0x0a,
	0x06, 0x67, 0x6f, 0x6c, 0x64, 0x5f, 0x35, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x67,
	0x6f, 0x6c, 0x64, 0x35, 0x12, 0x15, 0x0a, 0x06, 0x67, 0x6f, 0x6c, 0x64, 0x5f, 0x36, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x67, 0x6f, 0x6c, 0x64, 0x36, 0x12, 0x15, 0x0a, 0x06, 0x67,
	0x6f, 0x6c, 0x64, 0x5f, 0x37, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x67, 0x6f, 0x6c,
	0x64, 0x37, 0x12, 0x15, 0x0a, 0x06, 0x67, 0x6f, 0x6c, 0x64, 0x5f, 0x38, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x67, 0x6f, 0x6c, 0x64, 0x38, 0x12, 0x15, 0x0a, 0x06, 0x67, 0x6f, 0x6c,
	0x64, 0x5f, 0x39, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x67, 0x6f, 0x6c, 0x64, 0x39,
	0x12, 0x17, 0x0a, 0x07, 0x62, 0x6f, 0x6e, 0x75, 0x73, 0x5f, 0x30, 0x18, 0x0c, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x62, 0x6f, 0x6e, 0x75, 0x73, 0x30, 0x12, 0x17, 0x0a, 0x07, 0x62, 0x6f, 0x6e,
	0x75, 0x73, 0x5f, 0x31, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x62, 0x6f, 0x6e, 0x75,
	0x73, 0x31, 0x12, 0x17, 0x0a, 0x07, 0x62, 0x6f, 0x6e, 0x75, 0x73, 0x5f, 0x32, 0x18, 0x0e, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x06, 0x62, 0x6f, 0x6e, 0x75, 0x73, 0x32, 0x12, 0x17, 0x0a, 0x07, 0x62,
	0x6f, 0x6e, 0x75, 0x73, 0x5f, 0x33, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x62, 0x6f,
	0x6e, 0x75, 0x73, 0x33, 0x12, 0x17, 0x0a, 0x07, 0x62, 0x6f, 0x6e, 0x75, 0x73, 0x5f, 0x34, 0x18,
	0x10, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x62, 0x6f, 0x6e, 0x75, 0x73, 0x34, 0x12, 0x17, 0x0a,
	0x07, 0x62, 0x6f, 0x6e, 0x75, 0x73, 0x5f, 0x35, 0x18, 0x11, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x62, 0x6f, 0x6e, 0x75, 0x73, 0x35, 0x12, 0x17, 0x0a, 0x07, 0x62, 0x6f, 0x6e, 0x75, 0x73, 0x5f,
	0x36, 0x18, 0x12, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x62, 0x6f, 0x6e, 0x75, 0x73, 0x36, 0x12,
	0x17, 0x0a, 0x07, 0x62, 0x6f, 0x6e, 0x75, 0x73, 0x5f, 0x37, 0x18, 0x13, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x62, 0x6f, 0x6e, 0x75, 0x73, 0x37, 0x12, 0x17, 0x0a, 0x07, 0x62, 0x6f, 0x6e, 0x75,
	0x73, 0x5f, 0x38, 0x18, 0x14, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x62, 0x6f, 0x6e, 0x75, 0x73,
	0x38, 0x12, 0x17, 0x0a, 0x07, 0x62, 0x6f, 0x6e, 0x75, 0x73, 0x5f, 0x39, 0x18, 0x15, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x06, 0x62, 0x6f, 0x6e, 0x75, 0x73, 0x39, 0x42, 0x08, 0x5a, 0x06, 0x2e, 0x3b,
	0x70, 0x62, 0x67, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_gold_proto_rawDescOnce sync.Once
	file_pb_gold_proto_rawDescData = file_pb_gold_proto_rawDesc
)

func file_pb_gold_proto_rawDescGZIP() []byte {
	file_pb_gold_proto_rawDescOnce.Do(func() {
		file_pb_gold_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_gold_proto_rawDescData)
	})
	return file_pb_gold_proto_rawDescData
}

var file_pb_gold_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_pb_gold_proto_goTypes = []interface{}{
	(*Gold)(nil), // 0: gold
}
var file_pb_gold_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pb_gold_proto_init() }
func file_pb_gold_proto_init() {
	if File_pb_gold_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_gold_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Gold); i {
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
			RawDescriptor: file_pb_gold_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pb_gold_proto_goTypes,
		DependencyIndexes: file_pb_gold_proto_depIdxs,
		MessageInfos:      file_pb_gold_proto_msgTypes,
	}.Build()
	File_pb_gold_proto = out.File
	file_pb_gold_proto_rawDesc = nil
	file_pb_gold_proto_goTypes = nil
	file_pb_gold_proto_depIdxs = nil
}
