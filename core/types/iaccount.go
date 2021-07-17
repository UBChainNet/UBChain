package types

import (
	"github.com/UBChainNet/UBChain/common/hasharry"
)

// Account Status
type IAccount interface {
	GetBalance(string) uint64
	GetNonce() uint64
	Update(uint64) error
	StateKey() hasharry.Address
	IsExist() bool
	IsNeedUpdate() bool
	TransferChangeFrom(ITransaction, uint64) error
	ContractChangeFrom(ITransaction, uint64) error
	TransferChangeTo(*Receiver, hasharry.Address, uint64) error
	ContractChangeTo(*Receiver, hasharry.Address, uint64)
	FeesChange(uint64, uint64)
	ConsumptionChange(uint64, uint64)
	TransferOut(token hasharry.Address, amount, height uint64) error
	TransferIn(token hasharry.Address, amount, height uint64) error
	VerifyTxState(ITransaction) error
	VerifyNonce(uint64) error
	IsEmpty() bool
}

type IChainAddress interface {
	AddressList() []string
}
