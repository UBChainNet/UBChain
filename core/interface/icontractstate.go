package _interface

import (
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/core/types"
	"github.com/UBChainNet/UBChain/core/types/contractv2"
)

type IContractState interface {
	GetContract(contractAddr string) *types.Contract

	SetContract(contract *types.Contract)

	GetContractV2(contractAddr string) *contractv2.ContractV2

	SetContractV2(contract *contractv2.ContractV2)

	SetContractV2State(txHash string, contract *types.ContractV2State)

	GetContractV2State(hash string) *types.ContractV2State

	SetSymbol(symbol string, contract string)

	GetSymbolContract(symbol string) (string, bool)

	TokenList() []*types.Token

	VerifyState(tx types.ITransaction) error

	UpdateContract(tx types.ITransaction, blockHeight uint64)

	UpdateConfirmedHeight(height uint64)

	InitTrie(hash hasharry.Hash) error

	RootHash() hasharry.Hash

	ContractTrieCommit() (hasharry.Hash, error)

	Close() error
}
