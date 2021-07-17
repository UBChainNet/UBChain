package list

import (
	"fmt"
	"github.com/UBChainNet/UBChain/core/interface"
	"github.com/UBChainNet/UBChain/core/types"
	"sync"
	"time"
)

// Maximum transaction survival time
const TxLifeTime = 60 * 60 * 3

// List of transactions in the transaction pool
type TxList struct {
	// For transactions with a nonce value that is too large,
	// when a transaction with nonce -1 is packaged, the
	// transaction will be moved to the list of ready
	// transactions, ready to be packaged.
	futureTxs *FutureTxList

	// Ready to be packaged as a block transaction list.
	preparedTxs *TxSortedMap
	storage     ITxPoolStorage
	state       _interface.IAccountState
	mutex       sync.RWMutex
}

type ITxPoolStorage interface {
	Open() error
	LoadFutureTxs() *FutureTxList
	LoadPreparesTxs() *TxSortedMap
	SaveFutureTxs(*FutureTxList)
	SavePreparesTxs(*TxSortedMap)
	Close() error
}

func NewTxList(state _interface.IAccountState, storage ITxPoolStorage) *TxList {
	return &TxList{
		preparedTxs: NewTxSortedMap(),
		futureTxs:   NewFutureTxList(),
		storage:     storage,
		state:       state,
	}
}

func (t *TxList) Load() error {
	if err := t.storage.Open(); err != nil {
		return err
	}
	t.futureTxs = t.storage.LoadFutureTxs()
	t.preparedTxs = t.storage.LoadPreparesTxs()
	timeThreshold := uint64(time.Now().Unix() - TxLifeTime)
	t.RemoveExpiredTx(timeThreshold)
	t.UpdateTxsList()
	return nil
}

func (t *TxList) Close() error {
	t.storage.SaveFutureTxs(t.futureTxs)
	t.storage.SavePreparesTxs(t.preparedTxs)
	return t.storage.Close()
}

func (t *TxList) Len() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.futureTxs.Len() + t.preparedTxs.Len()
}

// Add a new transaction. If there is already a transaction with
// the same nonce value, the transaction fee for the new transaction
// needs to be greater than the transaction fee for the existing
// transaction, otherwise add returns an error. If the nonce value
// of the new transaction is greater than the nonce of the existing
// transaction, add To the list of future transactions.
func (t *TxList) Put(tx types.ITransaction) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	from := tx.From().String()
	nonce, _ := t.state.GetAccountNonce(tx.From())
	if nonce == tx.GetNonce()-1 {
		oldTx := t.preparedTxs.GetByAddress(from)
		if oldTx != nil {
			if oldTx.GetNonce() == tx.GetNonce() && oldTx.GetFees() < tx.GetFees() {
				t.preparedTxs.Remove(oldTx)
			} else if oldTx.GetNonce() < tx.GetNonce() {
				t.preparedTxs.Remove(oldTx)
			} else if oldTx.GetNonce() == tx.GetNonce() {
				return fmt.Errorf("the same nonce %d transaction already exists, so if you want to replace the nonce transaction, add a fee", tx.GetNonce())
			} else {
				return types.ErrTxNonceRepeat
			}
		}
		t.preparedTxs.Put(tx)
	} else if nonce >= tx.GetNonce() {
		return types.ErrTxNonceRepeat
	} else {
		return t.futureTxs.Put(tx)
	}
	return nil
}

//
func (t *TxList) RemoveMinFeeTx(newTx types.ITransaction) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.preparedTxs.PopMin(newTx.GetFees())
}

func (t *TxList) GetPreparedStuck(sec uint64) types.Transactions {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.preparedTxs.GetStuck(sec)
}

func (t *TxList) Gets(count int) types.Transactions {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.preparedTxs.Gets(count)
}

func (t *TxList) GetAll() (types.Transactions, types.Transactions) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	preparedTxs := t.preparedTxs.GetAll()
	futureTxs := t.futureTxs.GetAll()
	return preparedTxs, futureTxs
}

func (t *TxList) IsExist(from string, txHash string) bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if !t.preparedTxs.IsExist(txHash) {
		return t.futureTxs.IsExist(txHash)
	}
	return true
}

func (t *TxList) GetTransaction(hash string) (types.ITransaction, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	tx, isExist := t.preparedTxs.GetTransaction(hash)
	if !isExist {
		tx, isExist = t.futureTxs.GetTransaction(hash)
		if isExist {
			return tx, nil
		} else {
			return nil, fmt.Errorf("%s does not exist", hash)
		}
	}
	return tx, nil
}

func (t *TxList) UpdateTxsList() {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	t.preparedTxs.RemoveExecuted(t.state)

	for _, tx := range t.futureTxs.Txs {
		nonce, _ := t.state.GetAccountNonce(tx.From())
		if nonce < tx.GetNonce()-1 {
			continue
		}
		if nonce == tx.GetNonce()-1 {
			t.preparedTxs.Put(tx)
		}
		t.futureTxs.Remove(tx)
	}
}

func (t *TxList) RemoveExpiredTx(timeThreshold uint64) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.preparedTxs.RemoveExpiredTx(timeThreshold)

	for _, tx := range t.futureTxs.Txs {
		if tx.GetTime() <= timeThreshold {
			t.futureTxs.Remove(tx)
		}
	}
}

func (t *TxList) Remove(tx types.ITransaction) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.futureTxs.Remove(tx)
	t.preparedTxs.Remove(tx)
}
