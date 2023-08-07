// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.4
// source: pb/AllMight2.proto

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

// 編譯：執行 compile_protobuf.bat
type AllMight2 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// {"comment": "Message"}
	C *Currency `protobuf:"bytes,1,opt,name=c,proto3" json:"c,omitempty"`
	// {"comment": "map<int32, int32>"}
	Mii map[int32]int32 `protobuf:"bytes,2,rep,name=mii,proto3" json:"mii,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	// {"comment": "map<int32, string>"}
	Mis map[int32]string `protobuf:"bytes,3,rep,name=mis,proto3" json:"mis,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// {"comment": "map<string, int32>"}
	Msi map[string]int32 `protobuf:"bytes,4,rep,name=msi,proto3" json:"msi,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	// {"comment": "map<string, string>"}
	Mss map[string]string `protobuf:"bytes,5,rep,name=mss,proto3" json:"mss,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// {"comment": "TINYINT", "type":"TINYINT", "size": 4, "unsigned": "true"}
	Uti int32 `protobuf:"varint,6,opt,name=uti,proto3" json:"uti,omitempty"`
	// {"comment": "SMALINT", "type":"SMALINT", "size": 6, "unsigned": "true"}
	Usi int32 `protobuf:"varint,7,opt,name=usi,proto3" json:"usi,omitempty"`
	// {"comment": "MEDIUMINT", "type":"MEDIUMINT", "size": 9, "unsigned": "true", "primary_key": "default"}
	Umi int32 `protobuf:"varint,8,opt,name=umi,proto3" json:"umi,omitempty"`
	// {"comment": "INT", "type":"INT", "size": 11, "unsigned": "true", "primary_key": "default"}
	Ui int32 `protobuf:"varint,9,opt,name=ui,proto3" json:"ui,omitempty"`
}

func (x *AllMight2) Reset() {
	*x = AllMight2{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_AllMight2_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllMight2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllMight2) ProtoMessage() {}

func (x *AllMight2) ProtoReflect() protoreflect.Message {
	mi := &file_pb_AllMight2_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllMight2.ProtoReflect.Descriptor instead.
func (*AllMight2) Descriptor() ([]byte, []int) {
	return file_pb_AllMight2_proto_rawDescGZIP(), []int{0}
}

func (x *AllMight2) GetC() *Currency {
	if x != nil {
		return x.C
	}
	return nil
}

func (x *AllMight2) GetMii() map[int32]int32 {
	if x != nil {
		return x.Mii
	}
	return nil
}

func (x *AllMight2) GetMis() map[int32]string {
	if x != nil {
		return x.Mis
	}
	return nil
}

func (x *AllMight2) GetMsi() map[string]int32 {
	if x != nil {
		return x.Msi
	}
	return nil
}

func (x *AllMight2) GetMss() map[string]string {
	if x != nil {
		return x.Mss
	}
	return nil
}

func (x *AllMight2) GetUti() int32 {
	if x != nil {
		return x.Uti
	}
	return 0
}

func (x *AllMight2) GetUsi() int32 {
	if x != nil {
		return x.Usi
	}
	return 0
}

func (x *AllMight2) GetUmi() int32 {
	if x != nil {
		return x.Umi
	}
	return 0
}

func (x *AllMight2) GetUi() int32 {
	if x != nil {
		return x.Ui
	}
	return 0
}

var File_pb_AllMight2_proto protoreflect.FileDescriptor

var file_pb_AllMight2_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x62, 0x2f, 0x41, 0x6c, 0x6c, 0x4d, 0x69, 0x67, 0x68, 0x74, 0x32, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x70, 0x62, 0x2f, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe6, 0x03, 0x0a, 0x09, 0x41, 0x6c, 0x6c, 0x4d,
	0x69, 0x67, 0x68, 0x74, 0x32, 0x12, 0x17, 0x0a, 0x01, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x09, 0x2e, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x52, 0x01, 0x63, 0x12, 0x25,
	0x0a, 0x03, 0x6d, 0x69, 0x69, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x41, 0x6c,
	0x6c, 0x4d, 0x69, 0x67, 0x68, 0x74, 0x32, 0x2e, 0x4d, 0x69, 0x69, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x03, 0x6d, 0x69, 0x69, 0x12, 0x25, 0x0a, 0x03, 0x6d, 0x69, 0x73, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x13, 0x2e, 0x41, 0x6c, 0x6c, 0x4d, 0x69, 0x67, 0x68, 0x74, 0x32, 0x2e, 0x4d,
	0x69, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x03, 0x6d, 0x69, 0x73, 0x12, 0x25, 0x0a, 0x03,
	0x6d, 0x73, 0x69, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x41, 0x6c, 0x6c, 0x4d,
	0x69, 0x67, 0x68, 0x74, 0x32, 0x2e, 0x4d, 0x73, 0x69, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x03,
	0x6d, 0x73, 0x69, 0x12, 0x25, 0x0a, 0x03, 0x6d, 0x73, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x41, 0x6c, 0x6c, 0x4d, 0x69, 0x67, 0x68, 0x74, 0x32, 0x2e, 0x4d, 0x73, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x03, 0x6d, 0x73, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x74,
	0x69, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x75, 0x74, 0x69, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x73, 0x69, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x75, 0x73, 0x69, 0x12, 0x10,
	0x0a, 0x03, 0x75, 0x6d, 0x69, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x75, 0x6d, 0x69,
	0x12, 0x0e, 0x0a, 0x02, 0x75, 0x69, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x75, 0x69,
	0x1a, 0x36, 0x0a, 0x08, 0x4d, 0x69, 0x69, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x36, 0x0a, 0x08, 0x4d, 0x69, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x1a, 0x36, 0x0a, 0x08, 0x4d, 0x73, 0x69, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x36, 0x0a, 0x08, 0x4d, 0x73, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x42, 0x08, 0x5a, 0x06, 0x2e, 0x3b, 0x70, 0x62, 0x67, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_pb_AllMight2_proto_rawDescOnce sync.Once
	file_pb_AllMight2_proto_rawDescData = file_pb_AllMight2_proto_rawDesc
)

func file_pb_AllMight2_proto_rawDescGZIP() []byte {
	file_pb_AllMight2_proto_rawDescOnce.Do(func() {
		file_pb_AllMight2_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_AllMight2_proto_rawDescData)
	})
	return file_pb_AllMight2_proto_rawDescData
}

var file_pb_AllMight2_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pb_AllMight2_proto_goTypes = []interface{}{
	(*AllMight2)(nil), // 0: AllMight2
	nil,               // 1: AllMight2.MiiEntry
	nil,               // 2: AllMight2.MisEntry
	nil,               // 3: AllMight2.MsiEntry
	nil,               // 4: AllMight2.MssEntry
	(*Currency)(nil),  // 5: Currency
}
var file_pb_AllMight2_proto_depIdxs = []int32{
	5, // 0: AllMight2.c:type_name -> Currency
	1, // 1: AllMight2.mii:type_name -> AllMight2.MiiEntry
	2, // 2: AllMight2.mis:type_name -> AllMight2.MisEntry
	3, // 3: AllMight2.msi:type_name -> AllMight2.MsiEntry
	4, // 4: AllMight2.mss:type_name -> AllMight2.MssEntry
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_pb_AllMight2_proto_init() }
func file_pb_AllMight2_proto_init() {
	if File_pb_AllMight2_proto != nil {
		return
	}
	file_pb_Currency_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_pb_AllMight2_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllMight2); i {
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
			RawDescriptor: file_pb_AllMight2_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pb_AllMight2_proto_goTypes,
		DependencyIndexes: file_pb_AllMight2_proto_depIdxs,
		MessageInfos:      file_pb_AllMight2_proto_msgTypes,
	}.Build()
	File_pb_AllMight2_proto = out.File
	file_pb_AllMight2_proto_rawDesc = nil
	file_pb_AllMight2_proto_goTypes = nil
	file_pb_AllMight2_proto_depIdxs = nil
}
