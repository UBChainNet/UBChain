package node

import (
	"fmt"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/config"
	"github.com/UBChainNet/UBChain/consensus"
	"github.com/UBChainNet/UBChain/consensus/dpos"
	"github.com/UBChainNet/UBChain/core"
	"github.com/UBChainNet/UBChain/core/interface"
	runner2 "github.com/UBChainNet/UBChain/core/runner"
	"github.com/UBChainNet/UBChain/core/types"
	log "github.com/UBChainNet/UBChain/log/log15"
	"github.com/UBChainNet/UBChain/miner"
	"github.com/UBChainNet/UBChain/p2p"
	"github.com/UBChainNet/UBChain/param"
	"github.com/UBChainNet/UBChain/rpc"
	"github.com/UBChainNet/UBChain/services/accountstate"
	"github.com/UBChainNet/UBChain/services/blkmgr"
	"github.com/UBChainNet/UBChain/services/contractstate"
	"github.com/UBChainNet/UBChain/services/peermgr"
	"github.com/UBChainNet/UBChain/services/reqmgr"
	"github.com/UBChainNet/UBChain/services/txmgr"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)

type Node struct {
	miner *miner.Miner
	// Local p2p information
	localNode   *p2p.PeerInfo
	p2pServer   *p2p.P2pServer
	txPool      _interface.ITxPool
	peerManager p2p.IPeerManager
	blockManger *blkmgr.BlockManager
	blockChain  _interface.IBlockChain
	consensus   consensus.IConsensus
	network     blkmgr.Network
	// Node private key information
	private   *config.NodePrivate
	rpcServer *rpc.Server
}

func NewNode(cfg *config.Config) (*Node, error) {
	var err error
	node := &Node{}
	revBlkCh := make(chan *types.Block, 100)
	genBlkCh := make(chan *types.Block, 20)
	revTxCh := make(chan types.ITransaction, 50)
	minerWorkCh := make(chan bool)
	stateUpdateChan := make(chan struct{}, 50)
	removeTxsCh := make(chan types.Transactions, 100)
	secpKey, err := p2pcrypto.UnmarshalSecp256k1PrivateKey(cfg.NodePrivate.PrivateKey.Serialize())
	node.localNode = p2p.NewPeerInfo(secpKey, &peer.AddrInfo{}, nil)
	node.peerManager = peermgr.NewPeerManager(node.localNode)
	accountState, err := accountstate.NewAccountState(cfg.DataDir)
	if err != nil {
		return nil, fmt.Errorf("create account state failed! err:%s", err)
	}

	contractState, err := contractstate.NewContractState(cfg.DataDir)
	if err != nil {
		return nil, fmt.Errorf("create contract state failed! err:%s", err)
	}

	if node.consensus, err = dpos.NewDPos(cfg.DataDir, cfg.NodePrivate.Address, node); err != nil {
		return nil, fmt.Errorf("create dpos failed! err:%s", err)
	}

	runner := runner2.NewContractRunner(accountState, contractState)
	if node.blockChain, err = core.NewBlockChain(cfg.DataDir, node.consensus, stateUpdateChan, removeTxsCh, accountState, contractState, runner); err != nil {
		return nil, fmt.Errorf("create block chain failed! err:%s", err)
	}
	node.network = reqmgr.NewRequestManger(node.blockChain, revBlkCh, revTxCh, node)

	if node.p2pServer, err = p2p.NewP2pServer(cfg, node.localNode, node.peerManager, node.network); err != nil {
		return nil, fmt.Errorf("create p2p server failed! err:%s", err)
	}

	node.txPool = txmgr.NewTxPool(cfg, accountState, contractState, node.consensus, node.peerManager, node.network, runner, revTxCh, stateUpdateChan, removeTxsCh, node.p2pServer, node.blockChain.GetLastHeight)

	if err := node.consensus.Init(node.blockChain); err != nil {
		return nil, fmt.Errorf("init consensus failed! err:%s", err)
	}

	node.miner = miner.NewMiner(node.consensus, node.blockChain, node.txPool, cfg.NodePrivate.PrivateKey, cfg.NodePrivate.Address, genBlkCh, minerWorkCh)
	node.blockManger = blkmgr.NewBlockManager(node.blockChain, node.peerManager, node.network, node.consensus, revBlkCh, genBlkCh, minerWorkCh, node.p2pServer)
	node.private = cfg.NodePrivate
	rpcConfig := &config.RpcConfig{
		DataDir:  cfg.DataDir,
		RpcPort:  cfg.RpcPort,
		HttpPort: cfg.HttpPort,
		RpcTLS:   cfg.RpcTLS,
		RpcCert:  cfg.RpcCert,
		RpcPass:  cfg.RpcPass,
	}
	node.rpcServer = rpc.NewServer(rpcConfig, node.txPool, accountState, contractState, runner, node.consensus, node.blockChain, node.peerManager, node)

	if cfg.FallBackTo != config.DefaultFallBack && cfg.FallBackTo > 0 {
		if err := node.blockChain.FallBackTo(uint64(cfg.FallBackTo)); err != nil {
			return nil, err
		}
	}
	return node, nil
}

func (n *Node) Start() error {
	log.Info("System info", "version", param.Version, "network", param.Net, "p2p net", param.UniqueNetWork)
	if err := n.p2pServer.Start(); err != nil {
		log.Error("P2p startup failure", "err", err)
		return err
	}

	if err := n.blockManger.Start(); err != nil {
		log.Error("Block manager startup failure", "err", err)
		return err
	}

	if err := n.miner.Start(); err != nil {
		log.Error("Miner startup failure", "err", err)
		return err
	}

	if err := n.txPool.Start(); err != nil {
		log.Error("Transaction pool startup failure", "err", err)
		return err
	}

	if err := n.rpcServer.Start(); err != nil {
		log.Error("Rpc startup failure", "err", err)
		return err
	}

	go n.peerManager.Check()

	return nil
}

func (n *Node) Stop() {
	log.Info("Ready to quit")
	if err := n.blockManger.Stop(); err != nil {
		log.Error("Stop block manager failed!", "error", err)
	}
	if err := n.miner.Stop(); err != nil {
		log.Error("Stop miner failed!", "error", err)
	}
	if err := n.txPool.Stop(); err != nil {
		log.Error("Stop tx pool failed!", "error", err)
	}
	n.rpcServer.Close()
	if err := n.blockChain.CloseStorage(); err != nil {
		log.Error("Close storage failed!", "error", err)
	}
	if err := n.p2pServer.Stop(); err != nil {
		log.Error("Stop p2p failed!", "error", err)
	}
}

func (n *Node) SignHash(hash hasharry.Hash) (*types.SignScript, error) {
	return types.Sign(n.private.PrivateKey, hash)
}

func (n *Node) NodeInfo() *types.NodeInfo {
	return &types.NodeInfo{
		Version:     param.StringifySingleLine(),
		Net:         param.Net,
		P2pId:       n.p2pServer.ID(),
		P2pAddr:     n.p2pServer.Addr(),
		Connections: n.peerManager.Count(),
		Height:      n.blockChain.GetLastHeight(),
		Confirmed:   n.blockChain.GetConfirmedHeight(),
	}
}

func (n *Node) PeersInfo() []*types.NodeInfo {
	peerNodeInfos := make([]*types.NodeInfo, 0)
	peers := n.peerManager.Peers()
	for _, peer := range peers {
		streamCreator := p2p.StreamCreator{PeerId: peer.PeerId, NewStreamFunc: peer.NewStreamFunc}
		if nodeInfo, err := n.network.GetNodeInfo(&streamCreator); err != nil {
			continue
		} else {
			peerNodeInfos = append(peerNodeInfos, nodeInfo)
		}
	}
	return peerNodeInfos
}
