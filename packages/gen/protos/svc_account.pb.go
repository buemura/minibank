// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.12.4
// source: protos/svc_account.proto

package protos

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

type Account struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Balance       int32  `protobuf:"varint,2,opt,name=balance,proto3" json:"balance,omitempty"`
	OwnerName     string `protobuf:"bytes,3,opt,name=owner_name,json=ownerName,proto3" json:"owner_name,omitempty"`
	OwnerDocument string `protobuf:"bytes,4,opt,name=owner_document,json=ownerDocument,proto3" json:"owner_document,omitempty"`
	Status        string `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
	CreatedAt     string `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     string `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Account) Reset() {
	*x = Account{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_svc_account_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Account) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Account) ProtoMessage() {}

func (x *Account) ProtoReflect() protoreflect.Message {
	mi := &file_protos_svc_account_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Account.ProtoReflect.Descriptor instead.
func (*Account) Descriptor() ([]byte, []int) {
	return file_protos_svc_account_proto_rawDescGZIP(), []int{0}
}

func (x *Account) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Account) GetBalance() int32 {
	if x != nil {
		return x.Balance
	}
	return 0
}

func (x *Account) GetOwnerName() string {
	if x != nil {
		return x.OwnerName
	}
	return ""
}

func (x *Account) GetOwnerDocument() string {
	if x != nil {
		return x.OwnerDocument
	}
	return ""
}

func (x *Account) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Account) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

func (x *Account) GetUpdatedAt() string {
	if x != nil {
		return x.UpdatedAt
	}
	return ""
}

type GetAccountRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetAccountRequest) Reset() {
	*x = GetAccountRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_svc_account_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAccountRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAccountRequest) ProtoMessage() {}

func (x *GetAccountRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_svc_account_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAccountRequest.ProtoReflect.Descriptor instead.
func (*GetAccountRequest) Descriptor() ([]byte, []int) {
	return file_protos_svc_account_proto_rawDescGZIP(), []int{1}
}

func (x *GetAccountRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type CreateAccountRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OwnerName     string `protobuf:"bytes,1,opt,name=owner_name,json=ownerName,proto3" json:"owner_name,omitempty"`
	OwnerDocument string `protobuf:"bytes,2,opt,name=owner_document,json=ownerDocument,proto3" json:"owner_document,omitempty"`
}

func (x *CreateAccountRequest) Reset() {
	*x = CreateAccountRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_svc_account_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAccountRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAccountRequest) ProtoMessage() {}

func (x *CreateAccountRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_svc_account_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAccountRequest.ProtoReflect.Descriptor instead.
func (*CreateAccountRequest) Descriptor() ([]byte, []int) {
	return file_protos_svc_account_proto_rawDescGZIP(), []int{2}
}

func (x *CreateAccountRequest) GetOwnerName() string {
	if x != nil {
		return x.OwnerName
	}
	return ""
}

func (x *CreateAccountRequest) GetOwnerDocument() string {
	if x != nil {
		return x.OwnerDocument
	}
	return ""
}

type UpdateBalanceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	NewBalance int32  `protobuf:"varint,2,opt,name=new_balance,json=newBalance,proto3" json:"new_balance,omitempty"`
}

func (x *UpdateBalanceRequest) Reset() {
	*x = UpdateBalanceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_svc_account_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateBalanceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateBalanceRequest) ProtoMessage() {}

func (x *UpdateBalanceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_svc_account_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateBalanceRequest.ProtoReflect.Descriptor instead.
func (*UpdateBalanceRequest) Descriptor() ([]byte, []int) {
	return file_protos_svc_account_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateBalanceRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateBalanceRequest) GetNewBalance() int32 {
	if x != nil {
		return x.NewBalance
	}
	return 0
}

var File_protos_svc_account_proto protoreflect.FileDescriptor

var file_protos_svc_account_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x73, 0x76, 0x63, 0x5f, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcf, 0x01, 0x0a, 0x07, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65,
	0x12, 0x1d, 0x0a, 0x0a, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x25, 0x0a, 0x0e, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x5f, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x44, 0x6f,
	0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1d,
	0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a,
	0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x23, 0x0a, 0x11,
	0x47, 0x65, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x22, 0x5c, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x6f, 0x77, 0x6e,
	0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6f,
	0x77, 0x6e, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x6f, 0x77, 0x6e, 0x65,
	0x72, 0x5f, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x22,
	0x47, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x6e, 0x65, 0x77, 0x5f, 0x62,
	0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x6e, 0x65,
	0x77, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x32, 0xa6, 0x01, 0x0a, 0x0e, 0x41, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x2c, 0x0a, 0x0a, 0x67,
	0x65, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x12, 0x2e, 0x47, 0x65, 0x74, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x08, 0x2e,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x0d, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x15, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x08, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x00, 0x12, 0x32, 0x0a,
	0x0d, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x15,
	0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x08, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22,
	0x00, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x62, 0x75, 0x65, 0x6d, 0x75, 0x72, 0x61, 0x2f, 0x6d, 0x69, 0x6e, 0x69, 0x62, 0x61, 0x6e, 0x6b,
	0x2f, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x73, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_svc_account_proto_rawDescOnce sync.Once
	file_protos_svc_account_proto_rawDescData = file_protos_svc_account_proto_rawDesc
)

func file_protos_svc_account_proto_rawDescGZIP() []byte {
	file_protos_svc_account_proto_rawDescOnce.Do(func() {
		file_protos_svc_account_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_svc_account_proto_rawDescData)
	})
	return file_protos_svc_account_proto_rawDescData
}

var file_protos_svc_account_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_protos_svc_account_proto_goTypes = []any{
	(*Account)(nil),              // 0: Account
	(*GetAccountRequest)(nil),    // 1: GetAccountRequest
	(*CreateAccountRequest)(nil), // 2: CreateAccountRequest
	(*UpdateBalanceRequest)(nil), // 3: UpdateBalanceRequest
}
var file_protos_svc_account_proto_depIdxs = []int32{
	1, // 0: AccountService.getAccount:input_type -> GetAccountRequest
	2, // 1: AccountService.createAccount:input_type -> CreateAccountRequest
	3, // 2: AccountService.updateBalance:input_type -> UpdateBalanceRequest
	0, // 3: AccountService.getAccount:output_type -> Account
	0, // 4: AccountService.createAccount:output_type -> Account
	0, // 5: AccountService.updateBalance:output_type -> Account
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protos_svc_account_proto_init() }
func file_protos_svc_account_proto_init() {
	if File_protos_svc_account_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_svc_account_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Account); i {
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
		file_protos_svc_account_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*GetAccountRequest); i {
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
		file_protos_svc_account_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*CreateAccountRequest); i {
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
		file_protos_svc_account_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateBalanceRequest); i {
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
			RawDescriptor: file_protos_svc_account_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_svc_account_proto_goTypes,
		DependencyIndexes: file_protos_svc_account_proto_depIdxs,
		MessageInfos:      file_protos_svc_account_proto_msgTypes,
	}.Build()
	File_protos_svc_account_proto = out.File
	file_protos_svc_account_proto_rawDesc = nil
	file_protos_svc_account_proto_goTypes = nil
	file_protos_svc_account_proto_depIdxs = nil
}
