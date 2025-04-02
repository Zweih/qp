// go:generate protoc --go_out=. protobuf/pkginfo.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.1
// source: protobuf/pkginfo.proto

package protobuf

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RelationOp int32

const (
	RelationOp_NONE          RelationOp = 0
	RelationOp_EQUAL         RelationOp = 1
	RelationOp_LESS          RelationOp = 2
	RelationOp_LESS_EQUAL    RelationOp = 3
	RelationOp_GREATER       RelationOp = 4
	RelationOp_GREATER_EQUAL RelationOp = 5
)

// Enum value maps for RelationOp.
var (
	RelationOp_name = map[int32]string{
		0: "NONE",
		1: "EQUAL",
		2: "LESS",
		3: "LESS_EQUAL",
		4: "GREATER",
		5: "GREATER_EQUAL",
	}
	RelationOp_value = map[string]int32{
		"NONE":          0,
		"EQUAL":         1,
		"LESS":          2,
		"LESS_EQUAL":    3,
		"GREATER":       4,
		"GREATER_EQUAL": 5,
	}
)

func (x RelationOp) Enum() *RelationOp {
	p := new(RelationOp)
	*p = x
	return p
}

func (x RelationOp) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RelationOp) Descriptor() protoreflect.EnumDescriptor {
	return file_protobuf_pkginfo_proto_enumTypes[0].Descriptor()
}

func (RelationOp) Type() protoreflect.EnumType {
	return &file_protobuf_pkginfo_proto_enumTypes[0]
}

func (x RelationOp) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RelationOp.Descriptor instead.
func (RelationOp) EnumDescriptor() ([]byte, []int) {
	return file_protobuf_pkginfo_proto_rawDescGZIP(), []int{0}
}

type Relation struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version       string                 `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Operator      RelationOp             `protobuf:"varint,3,opt,name=operator,proto3,enum=pkginfo.RelationOp" json:"operator,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Relation) Reset() {
	*x = Relation{}
	mi := &file_protobuf_pkginfo_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Relation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Relation) ProtoMessage() {}

func (x *Relation) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_pkginfo_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Relation.ProtoReflect.Descriptor instead.
func (*Relation) Descriptor() ([]byte, []int) {
	return file_protobuf_pkginfo_proto_rawDescGZIP(), []int{0}
}

func (x *Relation) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Relation) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Relation) GetOperator() RelationOp {
	if x != nil {
		return x.Operator
	}
	return RelationOp_NONE
}

type PkgInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Timestamp     int64                  `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Size          int64                  `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	Name          string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Reason        string                 `protobuf:"bytes,4,opt,name=reason,proto3" json:"reason,omitempty"`
	Version       string                 `protobuf:"bytes,5,opt,name=version,proto3" json:"version,omitempty"`
	Arch          string                 `protobuf:"bytes,6,opt,name=arch,proto3" json:"arch,omitempty"`
	License       string                 `protobuf:"bytes,7,opt,name=license,proto3" json:"license,omitempty"`
	Url           string                 `protobuf:"bytes,8,opt,name=url,proto3" json:"url,omitempty"`
	Description   string                 `protobuf:"bytes,13,opt,name=description,proto3" json:"description,omitempty"`
	PkgBase       string                 `protobuf:"bytes,14,opt,name=pkg_base,json=pkgBase,proto3" json:"pkg_base,omitempty"`
	Depends       []*Relation            `protobuf:"bytes,9,rep,name=depends,proto3" json:"depends,omitempty"`
	RequiredBy    []*Relation            `protobuf:"bytes,10,rep,name=required_by,json=requiredBy,proto3" json:"required_by,omitempty"`
	Provides      []*Relation            `protobuf:"bytes,11,rep,name=provides,proto3" json:"provides,omitempty"`
	Conflicts     []*Relation            `protobuf:"bytes,12,rep,name=conflicts,proto3" json:"conflicts,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PkgInfo) Reset() {
	*x = PkgInfo{}
	mi := &file_protobuf_pkginfo_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PkgInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PkgInfo) ProtoMessage() {}

func (x *PkgInfo) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_pkginfo_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PkgInfo.ProtoReflect.Descriptor instead.
func (*PkgInfo) Descriptor() ([]byte, []int) {
	return file_protobuf_pkginfo_proto_rawDescGZIP(), []int{1}
}

func (x *PkgInfo) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *PkgInfo) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *PkgInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PkgInfo) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

func (x *PkgInfo) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *PkgInfo) GetArch() string {
	if x != nil {
		return x.Arch
	}
	return ""
}

func (x *PkgInfo) GetLicense() string {
	if x != nil {
		return x.License
	}
	return ""
}

func (x *PkgInfo) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *PkgInfo) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *PkgInfo) GetPkgBase() string {
	if x != nil {
		return x.PkgBase
	}
	return ""
}

func (x *PkgInfo) GetDepends() []*Relation {
	if x != nil {
		return x.Depends
	}
	return nil
}

func (x *PkgInfo) GetRequiredBy() []*Relation {
	if x != nil {
		return x.RequiredBy
	}
	return nil
}

func (x *PkgInfo) GetProvides() []*Relation {
	if x != nil {
		return x.Provides
	}
	return nil
}

func (x *PkgInfo) GetConflicts() []*Relation {
	if x != nil {
		return x.Conflicts
	}
	return nil
}

type CachedPkgs struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	LastModified  int64                  `protobuf:"varint,1,opt,name=last_modified,json=lastModified,proto3" json:"last_modified,omitempty"`
	Pkgs          []*PkgInfo             `protobuf:"bytes,2,rep,name=pkgs,proto3" json:"pkgs,omitempty"`
	Version       int32                  `protobuf:"varint,3,opt,name=version,proto3" json:"version,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CachedPkgs) Reset() {
	*x = CachedPkgs{}
	mi := &file_protobuf_pkginfo_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CachedPkgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CachedPkgs) ProtoMessage() {}

func (x *CachedPkgs) ProtoReflect() protoreflect.Message {
	mi := &file_protobuf_pkginfo_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CachedPkgs.ProtoReflect.Descriptor instead.
func (*CachedPkgs) Descriptor() ([]byte, []int) {
	return file_protobuf_pkginfo_proto_rawDescGZIP(), []int{2}
}

func (x *CachedPkgs) GetLastModified() int64 {
	if x != nil {
		return x.LastModified
	}
	return 0
}

func (x *CachedPkgs) GetPkgs() []*PkgInfo {
	if x != nil {
		return x.Pkgs
	}
	return nil
}

func (x *CachedPkgs) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

var File_protobuf_pkginfo_proto protoreflect.FileDescriptor

const file_protobuf_pkginfo_proto_rawDesc = "" +
	"\n" +
	"\x16protobuf/pkginfo.proto\x12\apkginfo\"i\n" +
	"\bRelation\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12\x18\n" +
	"\aversion\x18\x02 \x01(\tR\aversion\x12/\n" +
	"\boperator\x18\x03 \x01(\x0e2\x13.pkginfo.RelationOpR\boperator\"\xbf\x03\n" +
	"\aPkgInfo\x12\x1c\n" +
	"\ttimestamp\x18\x01 \x01(\x03R\ttimestamp\x12\x12\n" +
	"\x04size\x18\x02 \x01(\x03R\x04size\x12\x12\n" +
	"\x04name\x18\x03 \x01(\tR\x04name\x12\x16\n" +
	"\x06reason\x18\x04 \x01(\tR\x06reason\x12\x18\n" +
	"\aversion\x18\x05 \x01(\tR\aversion\x12\x12\n" +
	"\x04arch\x18\x06 \x01(\tR\x04arch\x12\x18\n" +
	"\alicense\x18\a \x01(\tR\alicense\x12\x10\n" +
	"\x03url\x18\b \x01(\tR\x03url\x12 \n" +
	"\vdescription\x18\r \x01(\tR\vdescription\x12\x19\n" +
	"\bpkg_base\x18\x0e \x01(\tR\apkgBase\x12+\n" +
	"\adepends\x18\t \x03(\v2\x11.pkginfo.RelationR\adepends\x122\n" +
	"\vrequired_by\x18\n" +
	" \x03(\v2\x11.pkginfo.RelationR\n" +
	"requiredBy\x12-\n" +
	"\bprovides\x18\v \x03(\v2\x11.pkginfo.RelationR\bprovides\x12/\n" +
	"\tconflicts\x18\f \x03(\v2\x11.pkginfo.RelationR\tconflicts\"q\n" +
	"\n" +
	"CachedPkgs\x12#\n" +
	"\rlast_modified\x18\x01 \x01(\x03R\flastModified\x12$\n" +
	"\x04pkgs\x18\x02 \x03(\v2\x10.pkginfo.PkgInfoR\x04pkgs\x12\x18\n" +
	"\aversion\x18\x03 \x01(\x05R\aversion*[\n" +
	"\n" +
	"RelationOp\x12\b\n" +
	"\x04NONE\x10\x00\x12\t\n" +
	"\x05EQUAL\x10\x01\x12\b\n" +
	"\x04LESS\x10\x02\x12\x0e\n" +
	"\n" +
	"LESS_EQUAL\x10\x03\x12\v\n" +
	"\aGREATER\x10\x04\x12\x11\n" +
	"\rGREATER_EQUAL\x10\x05B\x1cZ\x1ainternal/protobuf;protobufb\x06proto3"

var (
	file_protobuf_pkginfo_proto_rawDescOnce sync.Once
	file_protobuf_pkginfo_proto_rawDescData []byte
)

func file_protobuf_pkginfo_proto_rawDescGZIP() []byte {
	file_protobuf_pkginfo_proto_rawDescOnce.Do(func() {
		file_protobuf_pkginfo_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_protobuf_pkginfo_proto_rawDesc), len(file_protobuf_pkginfo_proto_rawDesc)))
	})
	return file_protobuf_pkginfo_proto_rawDescData
}

var (
	file_protobuf_pkginfo_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
	file_protobuf_pkginfo_proto_msgTypes  = make([]protoimpl.MessageInfo, 3)
	file_protobuf_pkginfo_proto_goTypes   = []any{
		(RelationOp)(0),    // 0: pkginfo.RelationOp
		(*Relation)(nil),   // 1: pkginfo.Relation
		(*PkgInfo)(nil),    // 2: pkginfo.PkgInfo
		(*CachedPkgs)(nil), // 3: pkginfo.CachedPkgs
	}
)
var file_protobuf_pkginfo_proto_depIdxs = []int32{
	0, // 0: pkginfo.Relation.operator:type_name -> pkginfo.RelationOp
	1, // 1: pkginfo.PkgInfo.depends:type_name -> pkginfo.Relation
	1, // 2: pkginfo.PkgInfo.required_by:type_name -> pkginfo.Relation
	1, // 3: pkginfo.PkgInfo.provides:type_name -> pkginfo.Relation
	1, // 4: pkginfo.PkgInfo.conflicts:type_name -> pkginfo.Relation
	2, // 5: pkginfo.CachedPkgs.pkgs:type_name -> pkginfo.PkgInfo
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_protobuf_pkginfo_proto_init() }
func file_protobuf_pkginfo_proto_init() {
	if File_protobuf_pkginfo_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_protobuf_pkginfo_proto_rawDesc), len(file_protobuf_pkginfo_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protobuf_pkginfo_proto_goTypes,
		DependencyIndexes: file_protobuf_pkginfo_proto_depIdxs,
		EnumInfos:         file_protobuf_pkginfo_proto_enumTypes,
		MessageInfos:      file_protobuf_pkginfo_proto_msgTypes,
	}.Build()
	File_protobuf_pkginfo_proto = out.File
	file_protobuf_pkginfo_proto_goTypes = nil
	file_protobuf_pkginfo_proto_depIdxs = nil
}
