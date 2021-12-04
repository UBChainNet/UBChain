package types

import (
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/crypto/ecc/secp256k1"
)

type ITransaction interface {
	Size() uint64
	IsCoinBase() bool
	VerifyTx() error
	VerifyCoinBaseTx(height, sumFees uint64, miner string) error
	EncodeToBytes() ([]byte, error)
	SignTx(key *secp256k1.PrivateKey) error
	SetHash() error
	TranslateToRlpTransaction() *RlpTransaction

	Hash() hasharry.Hash
	From() hasharry.Address
	GetFees() uint64
	GetNonce() uint64
	GetTime() uint64
	GetTxType() TransactionType
	GetSignScript() *SignScript
	GetTxHead() *TransactionHead
	GetTxBody() ITransactionBody
}

type ITransactionBody interface {
	GetAmount() uint64
	GetContract() hasharry.Address
	GetName() string
	GetAbbr() string
	GetDescription() string
	GetIncreaseSwitch() bool
	ToAddress() *Receivers
	GetPeerId() []byte
	VerifyBody(from hasharry.Address) error
}

type ITransactionIndex interface {
	GetHeight() uint64
}
