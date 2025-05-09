// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.12.4
// source: service/service.proto

package service

import (
	api "gitlab.ozon.dev/timofey15g/homework/pkg/api"
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

type TReqAcceptOrder struct {
	state               protoimpl.MessageState `protogen:"open.v1"`
	ID                  int64                  `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	UserID              int64                  `protobuf:"varint,2,opt,name=UserID,proto3" json:"UserID,omitempty"`
	StorageDurationDays int64                  `protobuf:"varint,3,opt,name=StorageDurationDays,proto3" json:"StorageDurationDays,omitempty"`
	Weight              float64                `protobuf:"fixed64,4,opt,name=Weight,proto3" json:"Weight,omitempty"`
	Cost                string                 `protobuf:"bytes,5,opt,name=Cost,proto3" json:"Cost,omitempty"`
	Package             string                 `protobuf:"bytes,6,opt,name=Package,proto3" json:"Package,omitempty"`
	ExtraPackage        string                 `protobuf:"bytes,7,opt,name=ExtraPackage,proto3" json:"ExtraPackage,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *TReqAcceptOrder) Reset() {
	*x = TReqAcceptOrder{}
	mi := &file_service_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TReqAcceptOrder) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TReqAcceptOrder) ProtoMessage() {}

func (x *TReqAcceptOrder) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TReqAcceptOrder.ProtoReflect.Descriptor instead.
func (*TReqAcceptOrder) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{0}
}

func (x *TReqAcceptOrder) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *TReqAcceptOrder) GetUserID() int64 {
	if x != nil {
		return x.UserID
	}
	return 0
}

func (x *TReqAcceptOrder) GetStorageDurationDays() int64 {
	if x != nil {
		return x.StorageDurationDays
	}
	return 0
}

func (x *TReqAcceptOrder) GetWeight() float64 {
	if x != nil {
		return x.Weight
	}
	return 0
}

func (x *TReqAcceptOrder) GetCost() string {
	if x != nil {
		return x.Cost
	}
	return ""
}

func (x *TReqAcceptOrder) GetPackage() string {
	if x != nil {
		return x.Package
	}
	return ""
}

func (x *TReqAcceptOrder) GetExtraPackage() string {
	if x != nil {
		return x.ExtraPackage
	}
	return ""
}

type TReqIssueOrder struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ids           []int64                `protobuf:"varint,1,rep,packed,name=Ids,proto3" json:"Ids,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TReqIssueOrder) Reset() {
	*x = TReqIssueOrder{}
	mi := &file_service_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TReqIssueOrder) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TReqIssueOrder) ProtoMessage() {}

func (x *TReqIssueOrder) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TReqIssueOrder.ProtoReflect.Descriptor instead.
func (*TReqIssueOrder) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{1}
}

func (x *TReqIssueOrder) GetIds() []int64 {
	if x != nil {
		return x.Ids
	}
	return nil
}

type TReqListHistory struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Limit         int64                  `protobuf:"varint,1,opt,name=Limit,proto3" json:"Limit,omitempty"`
	Offset        int64                  `protobuf:"varint,2,opt,name=Offset,proto3" json:"Offset,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TReqListHistory) Reset() {
	*x = TReqListHistory{}
	mi := &file_service_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TReqListHistory) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TReqListHistory) ProtoMessage() {}

func (x *TReqListHistory) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TReqListHistory.ProtoReflect.Descriptor instead.
func (*TReqListHistory) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{2}
}

func (x *TReqListHistory) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *TReqListHistory) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

type TReqListOrders struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserID        int64                  `protobuf:"varint,1,opt,name=UserID,proto3" json:"UserID,omitempty"`
	Limit         int64                  `protobuf:"varint,2,opt,name=Limit,proto3" json:"Limit,omitempty"`
	CursorID      int64                  `protobuf:"varint,3,opt,name=CursorID,proto3" json:"CursorID,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TReqListOrders) Reset() {
	*x = TReqListOrders{}
	mi := &file_service_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TReqListOrders) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TReqListOrders) ProtoMessage() {}

func (x *TReqListOrders) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TReqListOrders.ProtoReflect.Descriptor instead.
func (*TReqListOrders) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{3}
}

func (x *TReqListOrders) GetUserID() int64 {
	if x != nil {
		return x.UserID
	}
	return 0
}

func (x *TReqListOrders) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *TReqListOrders) GetCursorID() int64 {
	if x != nil {
		return x.CursorID
	}
	return 0
}

type TReqListReturns struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Limit         int64                  `protobuf:"varint,1,opt,name=Limit,proto3" json:"Limit,omitempty"`
	Offset        int64                  `protobuf:"varint,2,opt,name=Offset,proto3" json:"Offset,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TReqListReturns) Reset() {
	*x = TReqListReturns{}
	mi := &file_service_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TReqListReturns) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TReqListReturns) ProtoMessage() {}

func (x *TReqListReturns) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TReqListReturns.ProtoReflect.Descriptor instead.
func (*TReqListReturns) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{4}
}

func (x *TReqListReturns) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *TReqListReturns) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

type TReqReturnOrder struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OrderID       int64                  `protobuf:"varint,1,opt,name=OrderID,proto3" json:"OrderID,omitempty"`
	UserID        int64                  `protobuf:"varint,2,opt,name=UserID,proto3" json:"UserID,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TReqReturnOrder) Reset() {
	*x = TReqReturnOrder{}
	mi := &file_service_service_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TReqReturnOrder) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TReqReturnOrder) ProtoMessage() {}

func (x *TReqReturnOrder) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TReqReturnOrder.ProtoReflect.Descriptor instead.
func (*TReqReturnOrder) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{5}
}

func (x *TReqReturnOrder) GetOrderID() int64 {
	if x != nil {
		return x.OrderID
	}
	return 0
}

func (x *TReqReturnOrder) GetUserID() int64 {
	if x != nil {
		return x.UserID
	}
	return 0
}

type TReqWithdrawOrder struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OrderID       int64                  `protobuf:"varint,1,opt,name=OrderID,proto3" json:"OrderID,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TReqWithdrawOrder) Reset() {
	*x = TReqWithdrawOrder{}
	mi := &file_service_service_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TReqWithdrawOrder) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TReqWithdrawOrder) ProtoMessage() {}

func (x *TReqWithdrawOrder) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TReqWithdrawOrder.ProtoReflect.Descriptor instead.
func (*TReqWithdrawOrder) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{6}
}

func (x *TReqWithdrawOrder) GetOrderID() int64 {
	if x != nil {
		return x.OrderID
	}
	return 0
}

type TReqRenewTask struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	TaskID        int64                  `protobuf:"varint,1,opt,name=TaskID,proto3" json:"TaskID,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TReqRenewTask) Reset() {
	*x = TReqRenewTask{}
	mi := &file_service_service_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TReqRenewTask) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TReqRenewTask) ProtoMessage() {}

func (x *TReqRenewTask) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TReqRenewTask.ProtoReflect.Descriptor instead.
func (*TReqRenewTask) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{7}
}

func (x *TReqRenewTask) GetTaskID() int64 {
	if x != nil {
		return x.TaskID
	}
	return 0
}

type TStringResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Msg           string                 `protobuf:"bytes,1,opt,name=Msg,proto3" json:"Msg,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TStringResp) Reset() {
	*x = TStringResp{}
	mi := &file_service_service_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TStringResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TStringResp) ProtoMessage() {}

func (x *TStringResp) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TStringResp.ProtoReflect.Descriptor instead.
func (*TStringResp) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{8}
}

func (x *TStringResp) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type TResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Order         *api.TOrder            `protobuf:"bytes,1,opt,name=Order,proto3" json:"Order,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TResp) Reset() {
	*x = TResp{}
	mi := &file_service_service_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TResp) ProtoMessage() {}

func (x *TResp) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TResp.ProtoReflect.Descriptor instead.
func (*TResp) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{9}
}

func (x *TResp) GetOrder() *api.TOrder {
	if x != nil {
		return x.Order
	}
	return nil
}

type TListResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Orders        []*api.TOrder          `protobuf:"bytes,1,rep,name=Orders,proto3" json:"Orders,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TListResp) Reset() {
	*x = TListResp{}
	mi := &file_service_service_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TListResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TListResp) ProtoMessage() {}

func (x *TListResp) ProtoReflect() protoreflect.Message {
	mi := &file_service_service_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TListResp.ProtoReflect.Descriptor instead.
func (*TListResp) Descriptor() ([]byte, []int) {
	return file_service_service_proto_rawDescGZIP(), []int{10}
}

func (x *TListResp) GetOrders() []*api.TOrder {
	if x != nil {
		return x.Orders
	}
	return nil
}

var File_service_service_proto protoreflect.FileDescriptor

const file_service_service_proto_rawDesc = "" +
	"\n" +
	"\x15service/service.proto\x12\bNService\x1a\rapi/api.proto\"\xd5\x01\n" +
	"\x0fTReqAcceptOrder\x12\x0e\n" +
	"\x02ID\x18\x01 \x01(\x03R\x02ID\x12\x16\n" +
	"\x06UserID\x18\x02 \x01(\x03R\x06UserID\x120\n" +
	"\x13StorageDurationDays\x18\x03 \x01(\x03R\x13StorageDurationDays\x12\x16\n" +
	"\x06Weight\x18\x04 \x01(\x01R\x06Weight\x12\x12\n" +
	"\x04Cost\x18\x05 \x01(\tR\x04Cost\x12\x18\n" +
	"\aPackage\x18\x06 \x01(\tR\aPackage\x12\"\n" +
	"\fExtraPackage\x18\a \x01(\tR\fExtraPackage\"\"\n" +
	"\x0eTReqIssueOrder\x12\x10\n" +
	"\x03Ids\x18\x01 \x03(\x03R\x03Ids\"?\n" +
	"\x0fTReqListHistory\x12\x14\n" +
	"\x05Limit\x18\x01 \x01(\x03R\x05Limit\x12\x16\n" +
	"\x06Offset\x18\x02 \x01(\x03R\x06Offset\"Z\n" +
	"\x0eTReqListOrders\x12\x16\n" +
	"\x06UserID\x18\x01 \x01(\x03R\x06UserID\x12\x14\n" +
	"\x05Limit\x18\x02 \x01(\x03R\x05Limit\x12\x1a\n" +
	"\bCursorID\x18\x03 \x01(\x03R\bCursorID\"?\n" +
	"\x0fTReqListReturns\x12\x14\n" +
	"\x05Limit\x18\x01 \x01(\x03R\x05Limit\x12\x16\n" +
	"\x06Offset\x18\x02 \x01(\x03R\x06Offset\"C\n" +
	"\x0fTReqReturnOrder\x12\x18\n" +
	"\aOrderID\x18\x01 \x01(\x03R\aOrderID\x12\x16\n" +
	"\x06UserID\x18\x02 \x01(\x03R\x06UserID\"-\n" +
	"\x11TReqWithdrawOrder\x12\x18\n" +
	"\aOrderID\x18\x01 \x01(\x03R\aOrderID\"'\n" +
	"\rTReqRenewTask\x12\x16\n" +
	"\x06TaskID\x18\x01 \x01(\x03R\x06TaskID\"\x1f\n" +
	"\vTStringResp\x12\x10\n" +
	"\x03Msg\x18\x01 \x01(\tR\x03Msg\"+\n" +
	"\x05TResp\x12\"\n" +
	"\x05Order\x18\x01 \x01(\v2\f.NApi.TOrderR\x05Order\"1\n" +
	"\tTListResp\x12$\n" +
	"\x06Orders\x18\x01 \x03(\v2\f.NApi.TOrderR\x06Orders2\x9c\x04\n" +
	"\fOrderService\x12A\n" +
	"\vCreateOrder\x12\x19.NService.TReqAcceptOrder\x1a\x15.NService.TStringResp\"\x00\x12?\n" +
	"\n" +
	"IssueOrder\x12\x18.NService.TReqIssueOrder\x1a\x15.NService.TStringResp\"\x00\x12?\n" +
	"\vListHistory\x12\x19.NService.TReqListHistory\x1a\x13.NService.TListResp\"\x00\x12=\n" +
	"\n" +
	"ListOrders\x12\x18.NService.TReqListOrders\x1a\x13.NService.TListResp\"\x00\x12?\n" +
	"\vListReturns\x12\x19.NService.TReqListReturns\x1a\x13.NService.TListResp\"\x00\x12A\n" +
	"\vReturnOrder\x12\x19.NService.TReqReturnOrder\x1a\x15.NService.TStringResp\"\x00\x12E\n" +
	"\rWithdrawOrder\x12\x1b.NService.TReqWithdrawOrder\x1a\x15.NService.TStringResp\"\x00\x12=\n" +
	"\tRenewTask\x12\x17.NService.TReqRenewTask\x1a\x15.NService.TStringResp\"\x00B9Z7gitlab.ozon.dev/timofey15g/homework/pkg/service;serviceb\x06proto3"

var (
	file_service_service_proto_rawDescOnce sync.Once
	file_service_service_proto_rawDescData []byte
)

func file_service_service_proto_rawDescGZIP() []byte {
	file_service_service_proto_rawDescOnce.Do(func() {
		file_service_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_service_service_proto_rawDesc), len(file_service_service_proto_rawDesc)))
	})
	return file_service_service_proto_rawDescData
}

var file_service_service_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_service_service_proto_goTypes = []any{
	(*TReqAcceptOrder)(nil),   // 0: NService.TReqAcceptOrder
	(*TReqIssueOrder)(nil),    // 1: NService.TReqIssueOrder
	(*TReqListHistory)(nil),   // 2: NService.TReqListHistory
	(*TReqListOrders)(nil),    // 3: NService.TReqListOrders
	(*TReqListReturns)(nil),   // 4: NService.TReqListReturns
	(*TReqReturnOrder)(nil),   // 5: NService.TReqReturnOrder
	(*TReqWithdrawOrder)(nil), // 6: NService.TReqWithdrawOrder
	(*TReqRenewTask)(nil),     // 7: NService.TReqRenewTask
	(*TStringResp)(nil),       // 8: NService.TStringResp
	(*TResp)(nil),             // 9: NService.TResp
	(*TListResp)(nil),         // 10: NService.TListResp
	(*api.TOrder)(nil),        // 11: NApi.TOrder
}
var file_service_service_proto_depIdxs = []int32{
	11, // 0: NService.TResp.Order:type_name -> NApi.TOrder
	11, // 1: NService.TListResp.Orders:type_name -> NApi.TOrder
	0,  // 2: NService.OrderService.CreateOrder:input_type -> NService.TReqAcceptOrder
	1,  // 3: NService.OrderService.IssueOrder:input_type -> NService.TReqIssueOrder
	2,  // 4: NService.OrderService.ListHistory:input_type -> NService.TReqListHistory
	3,  // 5: NService.OrderService.ListOrders:input_type -> NService.TReqListOrders
	4,  // 6: NService.OrderService.ListReturns:input_type -> NService.TReqListReturns
	5,  // 7: NService.OrderService.ReturnOrder:input_type -> NService.TReqReturnOrder
	6,  // 8: NService.OrderService.WithdrawOrder:input_type -> NService.TReqWithdrawOrder
	7,  // 9: NService.OrderService.RenewTask:input_type -> NService.TReqRenewTask
	8,  // 10: NService.OrderService.CreateOrder:output_type -> NService.TStringResp
	8,  // 11: NService.OrderService.IssueOrder:output_type -> NService.TStringResp
	10, // 12: NService.OrderService.ListHistory:output_type -> NService.TListResp
	10, // 13: NService.OrderService.ListOrders:output_type -> NService.TListResp
	10, // 14: NService.OrderService.ListReturns:output_type -> NService.TListResp
	8,  // 15: NService.OrderService.ReturnOrder:output_type -> NService.TStringResp
	8,  // 16: NService.OrderService.WithdrawOrder:output_type -> NService.TStringResp
	8,  // 17: NService.OrderService.RenewTask:output_type -> NService.TStringResp
	10, // [10:18] is the sub-list for method output_type
	2,  // [2:10] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_service_service_proto_init() }
func file_service_service_proto_init() {
	if File_service_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_service_service_proto_rawDesc), len(file_service_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_service_proto_goTypes,
		DependencyIndexes: file_service_service_proto_depIdxs,
		MessageInfos:      file_service_service_proto_msgTypes,
	}.Build()
	File_service_service_proto = out.File
	file_service_service_proto_goTypes = nil
	file_service_service_proto_depIdxs = nil
}
