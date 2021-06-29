package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jhdriver/UBChain/common/hasharry"
	"github.com/jhdriver/UBChain/consensus"
	_interface "github.com/jhdriver/UBChain/core/interface"
	"github.com/jhdriver/UBChain/core/runner"
	coreTypes "github.com/jhdriver/UBChain/core/types"
	"github.com/jhdriver/UBChain/p2p"
	"github.com/jhdriver/UBChain/rpc/rpctypes"
	"github.com/jhdriver/UBChain/services/reqmgr"
	"strconv"
)

type Api struct {
	txPool        _interface.ITxPool
	accountState  _interface.IAccountState
	contractState _interface.IContractState
	runner        *runner.ContractRunner
	consensus     consensus.IConsensus
	chain         _interface.IBlockChain
	peerManager   p2p.IPeerManager
	peers         reqmgr.Peers
}

func NewApi(txPool _interface.ITxPool, state _interface.IAccountState, contractState _interface.IContractState,
	runner *runner.ContractRunner, consensus consensus.IConsensus, chain _interface.IBlockChain, peerManager p2p.IPeerManager,
	peers reqmgr.Peers) *Api {
	return &Api{
		txPool:        txPool,
		accountState:  state,
		contractState: contractState,
		consensus:     consensus,
		chain:         chain,
		peerManager:   peerManager,
		peers:         peers,
		runner:        runner,
	}
}

type ApiResponse struct {
}

func (a *Api) SendTransaction(raw string) (string, error) {
	var rpcTx *coreTypes.RpcTransaction
	if err := json.Unmarshal([]byte(raw), &rpcTx); err != nil {
		return "", err
	}
	tx, err := coreTypes.TranslateRpcTxToTx(rpcTx)
	if err != nil {
		return "", err
	}
	if err = a.txPool.Add(tx, false); err != nil {
		return "", err
	}
	return tx.Hash().String(), nil
}

func (a *Api) GetAccount(address string) (*rpctypes.Account, error) {
	addr := hasharry.StringToAddress(address)
	account := a.accountState.GetAccountState(addr)
	rpcAccount := rpctypes.TranslateAccountToRpcAccount(account.(*coreTypes.Account))
	return rpcAccount, nil
}

func (a *Api) GetTransaction(hashStr string) (*coreTypes.RpcTransactionConfirmed, error) {
	hash, err := hasharry.StringToHash(hashStr)
	if err != nil {
		return nil, errors.New("hash error")
	}
	var height uint64
	var confirmed bool
	tx, err := a.chain.GetTransaction(hash)
	if err != nil {
		tx, err = a.txPool.GetTransaction(hash)
		if err != nil {
			return nil, err
		}
		height = 0
		confirmed = false
	} else {
		index, err := a.chain.GetTransactionIndex(hash)
		if err != nil {
			return nil, fmt.Errorf("%s is not exist", hash.String())
		}
		height = index.GetHeight()
		confirmed = a.chain.GetConfirmedHeight() >= height
	}

	var rpcTx *coreTypes.RpcTransaction
	state, _ := a.chain.GetContractState(hash)
	if state != nil {
		rpcTx, _ = coreTypes.TranslateContractV2TxToRpcTx(tx.(*coreTypes.Transaction), state)
	} else {
		rpcTx, _ = coreTypes.TranslateTxToRpcTx(tx.(*coreTypes.Transaction))
	}
	rsMsg := &coreTypes.RpcTransactionConfirmed{
		TxHead:    rpcTx.TxHead,
		TxBody:    rpcTx.TxBody,
		Height:    height,
		Confirmed: confirmed,
	}
	return rsMsg, nil
}

func (a *Api) GetBlockByHash(hashStr string) (*coreTypes.RpcBlock, error) {
	hash, err := hasharry.StringToHash(hashStr)
	if err != nil {
		return nil, errors.New("hash error")
	}
	block, err := a.chain.GetBlockByHash(hash)
	if err != nil {
		return nil, err
	}
	rpcBlock, _ := coreTypes.TranslateBlockToRpcBlock(block, a.chain.GetConfirmedHeight(), a.chain.GetContractState)
	return rpcBlock, nil

}

func (a *Api) GetBlockByHeight(height uint64) (*coreTypes.RpcBlock, error) {
	block, err := a.chain.GetBlockByHeight(height)
	if err != nil {
		return nil, err
	}
	rpcBlock, _ := coreTypes.TranslateBlockToRpcBlock(block, a.chain.GetConfirmedHeight(), a.chain.GetContractState)
	return rpcBlock, nil
}

func (a *Api) GetPoolTxs() (*coreTypes.TxPool, error) {
	preparedTxs, futureTxs := a.txPool.GetAll()
	txPoolTxs, _ := coreTypes.TranslateTxsToRpcTxPool(preparedTxs, futureTxs)
	return txPoolTxs, nil
}

func (a *Api) GetCandidates() (*coreTypes.RpcCandidates, error) {
	candidates := a.consensus.GetCandidates(a.chain)
	if candidates == nil || len(candidates) == 0 {
		return nil, fmt.Errorf("no candidates")
	}
	return coreTypes.TranslateCandidatesToRpcCandidates(candidates), nil
}

func (a *Api) GetLastHeight() (string, error) {
	height := a.chain.GetLastHeight()
	sHeight := strconv.FormatUint(height, 10)
	return sHeight, nil
}

func (a *Api) GetContract(address string) (*coreTypes.RpcContract, error) {
	contract := a.contractState.GetContract(address)
	if contract == nil {
		return nil, fmt.Errorf("contract address %s is not exist", address)
	}
	return coreTypes.TranslateContractToRpcContract(contract), nil
}

func (a *Api) GetConfirmedHeight() (string, error) {
	height := a.chain.GetConfirmedHeight()
	sHeight := strconv.FormatUint(height, 10)
	return sHeight, nil
}

func (a *Api) Peers() ([]*coreTypes.NodeInfo, error) {
	peers := a.peers.PeersInfo()
	return peers, nil
}

func (a *Api) NodeInfo() (*coreTypes.NodeInfo, error) {
	node := a.peers.NodeInfo()
	return node, nil
}

func (a *Api) GetExchangePairs(address string) ([]*coreTypes.RpcPair, error) {
	pairs, err := a.runner.ExchangePair(hasharry.StringToAddress(address))
	if err != nil {
		return nil, err
	}
	return pairs, nil
}
