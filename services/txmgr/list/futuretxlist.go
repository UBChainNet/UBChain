package list

import (
	"fmt"
	"github.com/UBChainNet/UBChain/core/types"
)

type FutureTxList struct {
	Txs        map[string]types.ITransaction
	nonceKeMap map[string]map[uint64]string
}

func NewFutureTxList() *FutureTxList {
	return &FutureTxList{
		Txs:        make(map[string]types.ITransaction),
		nonceKeMap: make(map[string]map[uint64]string),
	}
}

func (f *FutureTxList) Put(tx types.ITransaction) error {
	if f.IsExist(tx.Hash().String()) {
		return fmt.Errorf("transation hash %s exsit", tx.Hash())
	}
	if oldTxHash := f.GetNonceKeyHash(tx.From().String(), tx.GetNonce()); oldTxHash != "" {
		oldTx := f.Txs[oldTxHash]
		if oldTx.GetFees() > tx.GetFees() {
			return fmt.Errorf("transation nonce %d exist, the fees must biger than before %d", tx.GetNonce(), oldTx.GetFees())
		}
		f.Remove(oldTx)
	}
	f.Txs[tx.Hash().String()] = tx
	nonceMap, exist := f.nonceKeMap[tx.From().String()]
	if exist{
		nonceMap[tx.GetNonce()] = tx.Hash().String()
	}else{
		f.nonceKeMap[tx.From().String()] = map[uint64]string{
			tx.GetNonce() : tx.Hash().String(),
		}
	}
	return nil
}

func (f *FutureTxList) Remove(tx types.ITransaction) {
	delete(f.Txs, tx.Hash().String())
	nonceMap, exist := f.nonceKeMap[tx.From().String()]
	if exist{
		delete(nonceMap, tx.GetNonce())
		if len(nonceMap) == 0{
			delete(f.nonceKeMap, tx.From().String())
		}
	}
}

func (f *FutureTxList) IsExist(txHash string) bool {
	_, ok := f.Txs[txHash]
	return ok
}

func (f *FutureTxList) GetTransaction(txHash string) (types.ITransaction, bool) {
	tx, ok := f.Txs[txHash]
	return tx, ok
}

func (f *FutureTxList) GetNonceKeyHash(from string, nonce uint64) string {
	nonceTx, exist := f.nonceKeMap[from]
	if exist{
		return nonceTx[nonce]
	}
	return ""
}

func (f *FutureTxList) GetAddressMaxNonce(from string) uint64 {
	nonceTx, exist := f.nonceKeMap[from]
	if exist{
		var max uint64
		for nonce, _ := range nonceTx{
			if nonce > max{
				max = nonce
			}
		}
		return max
	}
	return 0
}

func (f *FutureTxList) Len() int {
	return len(f.Txs)
}

func (f *FutureTxList) GetAll() types.Transactions {
	var all types.Transactions
	for _, tx := range f.Txs {
		all = append(all, tx)
	}
	return all
}
