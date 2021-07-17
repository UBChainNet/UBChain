package _interface

import (
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/core/types"
)

// Transaction pool interface, which is used to manage the transaction pool
type ITxPool interface {
	Start() error
	Stop() error
	Add(tx types.ITransaction, isPeer bool) error
	Gets(count int, maxBytes uint64) types.Transactions
	GetAll() (types.Transactions, types.Transactions)
	Get() types.ITransaction
	GetTransaction(hash hasharry.Hash) (types.ITransaction, error)
	Remove(txs types.Transactions)
	IsExist(tx types.ITransaction) bool
}
