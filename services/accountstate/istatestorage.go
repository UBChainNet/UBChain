package accountstate

import (
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/core/types"
)

// Storage interface for account balance information
type IAccountStorage interface {
	InitTrie(stateRoot hasharry.Hash) error
	AccountList() []string
	GetAccountState(stateKey hasharry.Address) types.IAccount
	SetAccountState(account types.IAccount)
	GetAccountBalance(stateKey hasharry.Address, contract string) uint64
	GetAccountNonce(stateKey hasharry.Address) uint64
	DeleteAccount(stateKey hasharry.Address)
	Commit() (hasharry.Hash, error)
	RootHash() hasharry.Hash
	Print()
	Close() error
}
