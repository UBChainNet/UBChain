package contractv2

import (
	"errors"
	"fmt"
	"github.com/UBChainNet/UBChain/common/encode/rlp"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/core/types/contractv2/exchange"
)

type ContractType uint
type FunctionType uint

const (
	Exchange_ ContractType = 0
	Pair_                  = 1
	Pledge_                = 2
	Token_                 = 3
	TokenHub_              = 4
)

const (
	Exchange_Init     FunctionType = 000000
	Exchange_SetAdmin              = 000001
	Exchange_SetFeeTo              = 000002
	Exchange_ExactIn               = 000003
	Exchange_ExactOut              = 000004

	Pair_AddLiquidity    = 100000
	Pair_RemoveLiquidity = 100001

	Pledge_Init         = 200000
	Pledge_Start        = 200001
	Pledge_AddPool      = 200002
	Pledge_RemovePool   = 200003
	Pledge_Add          = 200004
	Pledge_Remove       = 200005
	Pledge_RemoveReward = 200006
	Pledge_Update       = 200007

	Token_Init = 300000

	TokenHub_init = 400000
	TokenHub_Ack =  400001
)

type ContractV2 struct {
	Address    hasharry.Address
	CreateHash hasharry.Hash
	Type       ContractType
	Body       IContractV2Body
}

func (c *ContractV2) Bytes() []byte {
	rlpC := &RlpContractV2{
		Address:    c.Address,
		CreateHash: c.CreateHash,
		Type:       c.Type,
		Body:       c.Body.Bytes(),
	}
	bytes, _ := rlp.EncodeToBytes(rlpC)
	return bytes
}

func (c *ContractV2) Verify(function FunctionType, sender hasharry.Address) error {
	ex, _ := c.Body.(*exchange.Exchange)
	switch function {
	case Exchange_Init:
		return fmt.Errorf("exchange %s already exist", c.Address.String())
	case Exchange_SetAdmin:
		return ex.VerifySetter(sender)
	case Exchange_SetFeeTo:
		return ex.VerifySetter(sender)
	}

	return nil
}

type RlpContractV2 struct {
	Address    hasharry.Address
	CreateHash hasharry.Hash
	Type       ContractType
	Body       []byte
}

type IContractV2Body interface {
	Bytes() []byte
}

func DecodeContractV2(bytes []byte) (*ContractV2, error) {
	var rlpContract *RlpContractV2
	if err := rlp.DecodeBytes(bytes, &rlpContract); err != nil {
		return nil, err
	}
	var contract = &ContractV2{
		Address:    rlpContract.Address,
		CreateHash: rlpContract.CreateHash,
		Type:       rlpContract.Type,
		Body:       nil,
	}
	switch rlpContract.Type {
	case Exchange_:
		ex, err := exchange.DecodeToExchange(rlpContract.Body)
		if err != nil {
			return nil, err
		}
		contract.Body = ex
		return contract, err
	case Pair_:
		pair, err := exchange.DecodeToPair(rlpContract.Body)
		if err != nil {
			return nil, err
		}
		contract.Body = pair
		return contract, err
	case Pledge_:
		pair, err := exchange.DecodeToPledge(rlpContract.Body)
		if err != nil {
			return nil, err
		}
		contract.Body = pair
		return contract, err
	}
	return nil, errors.New("decoding failure")
}
