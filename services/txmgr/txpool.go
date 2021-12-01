package txmgr

import (
	"errors"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/config"
	"github.com/UBChainNet/UBChain/consensus"
	"github.com/UBChainNet/UBChain/core/interface"
	runner2 "github.com/UBChainNet/UBChain/core/runner"
	"github.com/UBChainNet/UBChain/core/types"
	"github.com/UBChainNet/UBChain/database/pooldb"
	log "github.com/UBChainNet/UBChain/log/log15"
	"github.com/UBChainNet/UBChain/p2p"
	"github.com/UBChainNet/UBChain/services/blkmgr"
	"github.com/UBChainNet/UBChain/services/txmgr/list"
	"github.com/libp2p/go-libp2p-core/peer"
	"time"
)

// Clear the expired transaction interval
const monitorTxInterval = 20

const reBroadcastInterval = 60 * 2

const txChanLength = 500

// Maximum number of transactions in the transaction pool
const maxPoolTx = 50000

const txPoolStorage = "txpool"

type lastHeightFunc func() uint64

// Manage transactions not packaged into blocks
type TxPool struct {
	accountState  _interface.IAccountState
	contractState _interface.IContractState
	consensus     consensus.IConsensus
	runner        *runner2.ContractRunner
	txs           *list.TxList
	peerManager   p2p.IPeerManager
	network       blkmgr.Network
	newStream     blkmgr.ICreateStream
	txChan        chan types.ITransaction
	recTx         chan types.ITransaction
	removeTxsCh   chan types.Transactions
	stateUpdateCh chan struct{}
	stop          chan bool
	lastHeightFunc
}

func NewTxPool(config *config.Config, accountState _interface.IAccountState, contractState _interface.IContractState,
	consensus consensus.IConsensus, peerManager p2p.IPeerManager, network blkmgr.Network, runner *runner2.ContractRunner,
	recTx chan types.ITransaction, stateUpdateCh chan struct{}, removeTxsCh chan types.Transactions,
	newStream blkmgr.ICreateStream, lastHeightFunc lastHeightFunc) *TxPool {

	return &TxPool{
		accountState:   accountState,
		contractState:  contractState,
		consensus:      consensus,
		runner:         runner,
		txs:            list.NewTxList(accountState, pooldb.NewTxPoolStorage(config.DataDir+"/"+txPoolStorage)),
		peerManager:    peerManager,
		network:        network,
		recTx:          recTx,
		removeTxsCh:    removeTxsCh,
		stateUpdateCh:  stateUpdateCh,
		newStream:      newStream,
		txChan:         make(chan types.ITransaction, txChanLength),
		stop:           make(chan bool, 1),
		lastHeightFunc: lastHeightFunc,
	}
}

// Start transaction pool
func (tp *TxPool) Start() error {
	if err := tp.txs.Load(); err != nil {
		return err
	}

	go tp.monitorTxTime()
	go tp.dealTx()
	go tp.reBroadcast()

	log.Info("Transaction pool startup successful")
	return nil
}

func (tp *TxPool) Stop() error {
	tp.stop <- true
	log.Info("Stop transaction pool")
	return tp.txs.Close()
}

// Monitor transaction time
func (tp *TxPool) monitorTxTime() {
	t := time.NewTicker(time.Second * monitorTxInterval)
	defer t.Stop()

	for range t.C {
		tp.clearExpiredTx()
	}
}

func (tp *TxPool) reBroadcast() {
	t := time.NewTicker(time.Second * reBroadcastInterval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			txs := tp.txs.GetPreparedStuck(uint64(time.Now().Unix()) - reBroadcastInterval)
			log.Info("Send stuck transaction", "txs", txs.Len())
			for _, tx := range txs {
				tp.txChan <- tx
			}
		}
	}
}

func (tp *TxPool) dealTx() {
	for {
		select {
		case _ = <-tp.stop:
			return
		case tx := <-tp.txChan:
			go tp.broadcastTx(tx)
		case tx := <-tp.recTx:
			go tp.Add(tx, true)
		case txs := <-tp.removeTxsCh:
			go tp.Remove(txs)
		case _ = <-tp.stateUpdateCh:
			go tp.txs.UpdateTxsList()
		}
	}
}

// Broadcast transaction
func (tp *TxPool) broadcastTx(tx types.ITransaction) {
	peers := tp.peerManager.Peers()
	for id, _ := range peers {
		if id != tp.peerManager.LocalPeerInfo().AddrInfo.ID.String() {
			peerId := new(peer.ID)
			if err := peerId.UnmarshalText([]byte(id)); err == nil {
				streamCreator := p2p.StreamCreator{PeerId: *peerId, NewStreamFunc: tp.newStream.CreateStream}
				go tp.network.SendTransaction(&streamCreator, tx)
			}
		}
	}
}

func (tp *TxPool) Add(tx types.ITransaction, isPeer bool) error {
	return tp.AddTransaction(tx, isPeer)
}

// Verify adding transactions to the transaction pool
func (tp *TxPool) AddTransaction(tx types.ITransaction, isPeer bool) error {
	log.Info("TxPool receive transaction", "hash", tx.Hash())
	if tp.IsExist(tx) {
		return errors.New("the transaction already exists")
	}

	if err := tp.verifyTx(tx); err != nil {
		return err
	}

	if tp.txs.Len() >= maxPoolTx {
		tp.txs.RemoveMinFeeTx(tx)
	}

	if err := tp.txs.Put(tx); err != nil {
		return err
	}
	log.Info("TxPool put transaction", "hash", tx.Hash())
	//if !isPeer {
	tp.txChan <- tx
	//}
	return nil
}

// Get transactions from the transaction pool
func (tp *TxPool) Gets(count int, maxSize uint64) types.Transactions {
	prepare := make(types.Transactions, 0)
	txs := tp.txs.Gets(count)
	failed := types.Transactions{}
	var txBytes uint64
	for _, tx := range txs {
		if err := tp.verifyTx(tx); err != nil {
			failed = append(failed, tx)
		} else {
			bytes, _ := tx.EncodeToBytes()
			txLength := uint64(len(bytes))
			if txBytes+txLength > maxSize {
				return txs
			}
			txBytes += uint64(len(bytes))
			prepare = append(prepare, tx)
		}
	}
	tp.Remove(failed)
	return prepare
}

// Get all transactions in the trading pool
func (tp *TxPool) GetAll() (types.Transactions, types.Transactions) {
	prepareTxs, futureTxs := tp.txs.GetAll()
	return prepareTxs, futureTxs
}

func (tp *TxPool) GetPendingNonce(address hasharry.Address) uint64 {
	nonce := tp.txs.GetPendingNonce(address)
	if nonce == 0{
		nonce, _ = tp.accountState.GetAccountNonce(address)
		return nonce
	}
	return nonce
}

// Delete transaction
func (tp *TxPool) Remove(txs types.Transactions) {
	for _, tx := range txs {
		switch tx.GetTxType {
		default:
			tp.txs.Remove(tx)
		}
	}

}

func (tp *TxPool) IsExist(tx types.ITransaction) bool {
	return tp.txs.IsExist(tx.From().String(), tx.Hash().String())
}

func (tp *TxPool) GetTransaction(hash hasharry.Hash) (types.ITransaction, error) {
	return tp.txs.GetTransaction(hash.String())
}

// Verify the transaction is legal
func (tp *TxPool) verifyTx(tx types.ITransaction) error {
	if err := tx.VerifyTx(); err != nil {
		return err
	}

	if err := tp.consensus.VerifyTx(tx); err != nil {
		return err
	}

	if err := tp.accountState.VerifyState(tx); err != nil {
		return err
	}

	if err := tp.contractState.VerifyState(tx); err != nil {
		return err
	}

	if err := tp.runner.Verify(tx, tp.lastHeightFunc()); err != nil {
		return err
	}

	return nil
}

func (tp *TxPool) clearExpiredTx() {
	timeThreshold := time.Now().Unix() - list.TxLifeTime
	tp.txs.RemoveExpiredTx(uint64(timeThreshold))
}
