package transaction

import (
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/core/types"
	"github.com/UBChainNet/UBChain/core/types/contractv2"
	"github.com/UBChainNet/UBChain/core/types/functionbody/tokenhub_func"
	"github.com/UBChainNet/UBChain/param"
	"time"
)

func NewTokenHubInit(from, contract, setter, admin, feeTo string, feeRate float64, nonce uint64, note string) (*types.Transaction, error) {
	tx := &types.Transaction{
		TxHead: &types.TransactionHead{
			TxType:     types.ContractV2_,
			TxHash:     hasharry.Hash{},
			From:       hasharry.StringToAddress(from),
			Nonce:      nonce,
			Time:       uint64(time.Now().Unix()),
			Note:       note,
			SignScript: &types.SignScript{},
			Fees:       param.Fees,
		},
		TxBody: &types.TxContractV2Body{
			Contract:     hasharry.StringToAddress(contract),
			Type:         contractv2.TokenHub_,
			FunctionType: contractv2.TokenHub_init,
			Function: &tokenhub_func.TokenHubInitBody{
				Setter:  hasharry.StringToAddress(setter),
				Admin:   hasharry.StringToAddress(admin),
				FeeTo:   hasharry.StringToAddress(feeTo),
				FeeRate: feeRate,
			},
		},
	}
	tx.SetHash()
	return tx, nil
}

func NewTokenHubAck(from, contract string, sequences []uint64, ackTypes []uint8 , nonce uint64, note string) (*types.Transaction, error) {
	tx := &types.Transaction{
		TxHead: &types.TransactionHead{
			TxType:     types.ContractV2_,
			TxHash:     hasharry.Hash{},
			From:       hasharry.StringToAddress(from),
			Nonce:      nonce,
			Time:       uint64(time.Now().Unix()),
			Note:       note,
			SignScript: &types.SignScript{},
			Fees:       param.Fees,
		},
		TxBody: &types.TxContractV2Body{
			Contract:     hasharry.StringToAddress(contract),
			Type:         contractv2.TokenHub_,
			FunctionType: contractv2.TokenHub_Ack,
			Function: &tokenhub_func.TokenHubAckBody{
				Sequences: sequences,
				AckTypes:  ackTypes,
			},
		},
	}
	tx.SetHash()
	return tx, nil
}