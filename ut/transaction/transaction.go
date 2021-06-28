package transaction

import (
	"github.com/jhdriver/UBChain/common/hasharry"
	"github.com/jhdriver/UBChain/core/types"
	"github.com/jhdriver/UBChain/param"
	"time"
)

func NewTransactionV2(from string, toMap []map[string]uint64, token string, note string, nonce uint64) *types.Transaction {
	tx := &types.Transaction{
		TxHead: &types.TransactionHead{
			TxType:     types.Transfer_,
			TxHash:     hasharry.Hash{},
			From:       hasharry.StringToAddress(from),
			Nonce:      nonce,
			Time:       uint64(time.Now().Unix()),
			Note:       note,
			SignScript: &types.SignScript{},
		},
	}
	txBody := &types.TransferBody{
		Contract:  hasharry.StringToAddress(token),
		Receivers: types.NewReceivers(),
	}
	for _, to := range toMap {
		for address, amount := range to {
			txBody.Receivers.Add(hasharry.StringToAddress(address), amount)
		}
	}
	tx.TxBody = txBody
	tx.TxHead.Fees = types.TransferFees(len(txBody.Receivers.ReceiverList()))
	tx.SetHash()
	return tx
}

func NewContract(from, to, contract string, note string, amount, nonce uint64, name, abbr string, increase bool, description string) *types.Transaction {
	tx := &types.Transaction{
		TxHead: &types.TransactionHead{
			TxType:     types.Contract_,
			TxHash:     hasharry.Hash{},
			From:       hasharry.StringToAddress(from),
			Nonce:      nonce,
			Time:       uint64(time.Now().Unix()),
			Note:       note,
			SignScript: &types.SignScript{},
			Fees:       param.TokenConsumption,
		},
		TxBody: &types.ContractBody{
			Contract:       hasharry.StringToAddress(contract),
			To:             hasharry.StringToAddress(to),
			Name:           name,
			Abbr:           abbr,
			IncreaseSwitch: increase,
			Description:    description,
			Amount:         amount,
		},
	}
	tx.SetHash()
	return tx
}
