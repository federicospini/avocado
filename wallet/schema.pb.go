// Code generated by protoc-gen-go.
// source: github.com/jtremback/upc-core/wallet/schema.proto
// DO NOT EDIT!

/*
Package wallet is a generated protocol buffer package.

It is generated from these files:
	github.com/jtremback/upc-core/wallet/schema.proto

It has these top-level messages:
	Channel
	Account
	EscrowProvider
*/
package wallet

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import wire "github.com/jtremback/upc-core/wire"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type Channel_Phase int32

const (
	Channel_PendingOpen   Channel_Phase = 0
	Channel_Open          Channel_Phase = 1
	Channel_PendingClosed Channel_Phase = 2
	Channel_Closed        Channel_Phase = 3
)

var Channel_Phase_name = map[int32]string{
	0: "PendingOpen",
	1: "Open",
	2: "PendingClosed",
	3: "Closed",
}
var Channel_Phase_value = map[string]int32{
	"PendingOpen":   0,
	"Open":          1,
	"PendingClosed": 2,
	"Closed":        3,
}

func (x Channel_Phase) String() string {
	return proto.EnumName(Channel_Phase_name, int32(x))
}
func (Channel_Phase) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Channel struct {
	ChannelId string        `protobuf:"bytes,1,opt,name=channel_id" json:"channel_id,omitempty"`
	Phase     Channel_Phase `protobuf:"varint,2,opt,name=phase,enum=wallet.Channel_Phase" json:"phase,omitempty"`
	// This is how
	OpeningTx                *wire.OpeningTx `protobuf:"bytes,3,opt,name=opening_tx" json:"opening_tx,omitempty"`
	OpeningTxEnvelope        *wire.Envelope  `protobuf:"bytes,4,opt,name=opening_tx_envelope" json:"opening_tx_envelope,omitempty"`
	LastUpdateTx             *wire.UpdateTx  `protobuf:"bytes,5,opt,name=last_update_tx" json:"last_update_tx,omitempty"`
	LastUpdateTxEnvelope     *wire.Envelope  `protobuf:"bytes,6,opt,name=last_update_tx_envelope" json:"last_update_tx_envelope,omitempty"`
	LastFullUpdateTx         *wire.UpdateTx  `protobuf:"bytes,7,opt,name=last_full_update_tx" json:"last_full_update_tx,omitempty"`
	LastFullUpdateTxEnvelope *wire.Envelope  `protobuf:"bytes,8,opt,name=last_full_update_tx_envelope" json:"last_full_update_tx_envelope,omitempty"`
	EscrowProvider           *EscrowProvider `protobuf:"bytes,9,opt,name=escrow_provider" json:"escrow_provider,omitempty"`
	Accounts                 []*Account      `protobuf:"bytes,10,rep,name=accounts" json:"accounts,omitempty"`
	Me                       uint32          `protobuf:"varint,12,opt,name=me" json:"me,omitempty"`
	Fulfillments             [][]byte        `protobuf:"bytes,13,rep,name=fulfillments,proto3" json:"fulfillments,omitempty"`
}

func (m *Channel) Reset()                    { *m = Channel{} }
func (m *Channel) String() string            { return proto.CompactTextString(m) }
func (*Channel) ProtoMessage()               {}
func (*Channel) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Channel) GetOpeningTx() *wire.OpeningTx {
	if m != nil {
		return m.OpeningTx
	}
	return nil
}

func (m *Channel) GetOpeningTxEnvelope() *wire.Envelope {
	if m != nil {
		return m.OpeningTxEnvelope
	}
	return nil
}

func (m *Channel) GetLastUpdateTx() *wire.UpdateTx {
	if m != nil {
		return m.LastUpdateTx
	}
	return nil
}

func (m *Channel) GetLastUpdateTxEnvelope() *wire.Envelope {
	if m != nil {
		return m.LastUpdateTxEnvelope
	}
	return nil
}

func (m *Channel) GetLastFullUpdateTx() *wire.UpdateTx {
	if m != nil {
		return m.LastFullUpdateTx
	}
	return nil
}

func (m *Channel) GetLastFullUpdateTxEnvelope() *wire.Envelope {
	if m != nil {
		return m.LastFullUpdateTxEnvelope
	}
	return nil
}

func (m *Channel) GetEscrowProvider() *EscrowProvider {
	if m != nil {
		return m.EscrowProvider
	}
	return nil
}

func (m *Channel) GetAccounts() []*Account {
	if m != nil {
		return m.Accounts
	}
	return nil
}

type Account struct {
	Name           string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Pubkey         []byte `protobuf:"bytes,2,opt,name=pubkey,proto3" json:"pubkey,omitempty"`
	Privkey        []byte `protobuf:"bytes,3,opt,name=privkey,proto3" json:"privkey,omitempty"`
	Address        string `protobuf:"bytes,4,opt,name=address" json:"address,omitempty"`
	EscrowProvider string `protobuf:"bytes,5,opt,name=escrow_provider" json:"escrow_provider,omitempty"`
}

func (m *Account) Reset()                    { *m = Account{} }
func (m *Account) String() string            { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()               {}
func (*Account) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type EscrowProvider struct {
	Name    string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Pubkey  []byte `protobuf:"bytes,2,opt,name=pubkey,proto3" json:"pubkey,omitempty"`
	Privkey []byte `protobuf:"bytes,3,opt,name=privkey,proto3" json:"privkey,omitempty"`
	Address string `protobuf:"bytes,4,opt,name=address" json:"address,omitempty"`
}

func (m *EscrowProvider) Reset()                    { *m = EscrowProvider{} }
func (m *EscrowProvider) String() string            { return proto.CompactTextString(m) }
func (*EscrowProvider) ProtoMessage()               {}
func (*EscrowProvider) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func init() {
	proto.RegisterType((*Channel)(nil), "wallet.Channel")
	proto.RegisterType((*Account)(nil), "wallet.Account")
	proto.RegisterType((*EscrowProvider)(nil), "wallet.EscrowProvider")
	proto.RegisterEnum("wallet.Channel_Phase", Channel_Phase_name, Channel_Phase_value)
}

var fileDescriptor0 = []byte{
	// 497 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xb4, 0x53, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0x25, 0x1f, 0x76, 0xe2, 0xa9, 0xe3, 0xa4, 0x5b, 0x28, 0x16, 0x02, 0xa9, 0xf2, 0x05, 0x24,
	0x54, 0x47, 0x14, 0xae, 0x80, 0xaa, 0x2a, 0x48, 0x5c, 0x20, 0x5a, 0xc1, 0xd9, 0xda, 0xd8, 0xd3,
	0xc6, 0x74, 0xfd, 0x21, 0x7f, 0xa4, 0xe5, 0x6f, 0xf0, 0x63, 0xf8, 0x7d, 0xac, 0x77, 0xd7, 0x4e,
	0x52, 0x22, 0x71, 0xe2, 0x12, 0xcd, 0xcc, 0x9b, 0x37, 0xef, 0x65, 0x67, 0x0c, 0x6f, 0x6e, 0xe2,
	0x6a, 0x5d, 0xaf, 0xfc, 0x30, 0x4b, 0xe6, 0x3f, 0xaa, 0x02, 0x93, 0x15, 0x0b, 0x6f, 0xe7, 0x75,
	0x1e, 0x9e, 0x87, 0x59, 0x81, 0xf3, 0x3b, 0xc6, 0x39, 0x56, 0xf3, 0x32, 0x5c, 0x63, 0xc2, 0xfc,
	0xbc, 0xc8, 0xaa, 0x8c, 0x98, 0xaa, 0xf8, 0xec, 0xfc, 0x1f, 0xd4, 0x58, 0xff, 0x28, 0x9a, 0xf7,
	0xdb, 0x80, 0xd1, 0xd5, 0x9a, 0xa5, 0x29, 0x72, 0xf2, 0x02, 0x20, 0x54, 0x61, 0x10, 0x47, 0x6e,
	0xef, 0xac, 0xf7, 0xca, 0xa2, 0x96, 0xae, 0x7c, 0x8e, 0xc8, 0x6b, 0x30, 0xf2, 0x35, 0x2b, 0xd1,
	0xed, 0x0b, 0xc4, 0xb9, 0x78, 0xe2, 0x2b, 0x45, 0x5f, 0xd3, 0xfd, 0x65, 0x03, 0x52, 0xd5, 0x43,
	0x7c, 0x80, 0x2c, 0xc7, 0x34, 0x4e, 0x6f, 0x82, 0xea, 0xde, 0x1d, 0x08, 0xc6, 0xd1, 0xc5, 0xd4,
	0x97, 0xc2, 0x5f, 0x55, 0xfd, 0xdb, 0x3d, 0xb5, 0xb2, 0x36, 0x24, 0x1f, 0xe0, 0x64, 0xdb, 0x1f,
	0x60, 0xba, 0x41, 0x2e, 0x72, 0x77, 0x28, 0x89, 0x8e, 0x22, 0x2e, 0x74, 0x95, 0x1e, 0x77, 0xbc,
	0xb6, 0x44, 0xde, 0x81, 0xc3, 0x59, 0x59, 0x05, 0x75, 0x1e, 0xb1, 0x0a, 0x1b, 0x4d, 0x63, 0x97,
	0xfa, 0x5d, 0x96, 0x85, 0xa4, 0xdd, 0x74, 0xb5, 0x19, 0x59, 0xc0, 0xd3, 0x7d, 0xd6, 0x56, 0xd9,
	0x3c, 0xa8, 0xfc, 0x78, 0x97, 0xde, 0x89, 0xbf, 0x87, 0x13, 0x39, 0xe6, 0xba, 0xe6, 0x7c, 0xc7,
	0xc1, 0xe8, 0xa0, 0x83, 0x59, 0xd3, 0xfa, 0x49, 0x74, 0x76, 0x2e, 0xbe, 0xc0, 0xf3, 0x03, 0xf4,
	0xad, 0x95, 0xf1, 0x41, 0x2b, 0xee, 0xc3, 0x39, 0x9d, 0x9d, 0x8f, 0x30, 0xc5, 0x32, 0x2c, 0xb2,
	0xbb, 0x40, 0xec, 0x78, 0x13, 0x47, 0x58, 0xb8, 0x96, 0x1c, 0x71, 0xda, 0xae, 0x6c, 0x21, 0xe1,
	0xa5, 0x46, 0xa9, 0x83, 0x7b, 0xb9, 0xd8, 0xf4, 0x98, 0x85, 0x61, 0x56, 0xa7, 0x55, 0xe9, 0xc2,
	0xd9, 0x40, 0xad, 0x4e, 0x31, 0x2f, 0x55, 0x9d, 0x76, 0x0d, 0xc4, 0x81, 0x7e, 0x82, 0xae, 0x2d,
	0x04, 0x26, 0x54, 0x44, 0xc4, 0x03, 0x5b, 0xfc, 0x91, 0xeb, 0x98, 0xf3, 0x04, 0x9b, 0x01, 0x13,
	0x31, 0xc0, 0xa6, 0x7b, 0x35, 0xef, 0x12, 0x0c, 0x79, 0x2d, 0x64, 0x0a, 0x47, 0x4b, 0x4c, 0x23,
	0xb1, 0xcb, 0xe6, 0x2a, 0x66, 0x8f, 0xc8, 0x18, 0x86, 0x32, 0xea, 0x91, 0x63, 0x98, 0x68, 0xe8,
	0x8a, 0x67, 0x25, 0x46, 0xb3, 0x3e, 0x01, 0x30, 0x75, 0x3c, 0xf0, 0x7e, 0xf5, 0x60, 0xa4, 0xcd,
	0x10, 0x02, 0xc3, 0x94, 0x09, 0x13, 0xea, 0x64, 0x65, 0x4c, 0x4e, 0xc1, 0xcc, 0xeb, 0xd5, 0x2d,
	0xfe, 0x94, 0xe7, 0x6a, 0x53, 0x9d, 0x11, 0x17, 0x46, 0x79, 0x11, 0x6f, 0x1a, 0x60, 0x20, 0x81,
	0x36, 0x6d, 0x10, 0x16, 0x45, 0x05, 0x96, 0xa5, 0x3c, 0x3b, 0x8b, 0xb6, 0x29, 0x79, 0xf9, 0xf7,
	0x83, 0x1a, 0xb2, 0xe3, 0xc1, 0xc3, 0x79, 0x39, 0x38, 0xfb, 0x4f, 0xfb, 0xbf, 0xad, 0xad, 0x4c,
	0xf9, 0x19, 0xbf, 0xfd, 0x13, 0x00, 0x00, 0xff, 0xff, 0xd5, 0xb3, 0x80, 0x79, 0x32, 0x04, 0x00,
	0x00,
}