package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/UBChainNet/UBChain/common/utils"
	"github.com/UBChainNet/UBChain/config"
	"github.com/UBChainNet/UBChain/consensus"
	"github.com/UBChainNet/UBChain/core/interface"
	"github.com/UBChainNet/UBChain/core/runner"
	"github.com/UBChainNet/UBChain/crypto/certgen"
	log "github.com/UBChainNet/UBChain/log/log15"
	"github.com/UBChainNet/UBChain/p2p"
	"github.com/UBChainNet/UBChain/rpc/rpchttp"
	"github.com/UBChainNet/UBChain/rpc/rpctypes"
	"github.com/UBChainNet/UBChain/services/reqmgr"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

type Server struct {
	config     *config.RpcConfig
	api        *Api
	grpcServer *grpc.Server
	httpServer *rpchttp.RpcServer
}

func NewServer(config *config.RpcConfig, txPool _interface.ITxPool, state _interface.IAccountState, contractState _interface.IContractState,
	runner *runner.ContractRunner, consensus consensus.IConsensus, chain _interface.IBlockChain, peerManager p2p.IPeerManager,
	peers reqmgr.Peers) *Server {
	return &Server{config: config, api: NewApi(txPool, state, contractState, runner, consensus, chain, peerManager, peers)}
}

func (rs *Server) Start() error {
	var err error
	endPoint := "0.0.0.0:" + rs.config.HttpPort
	lis, err := net.Listen("tcp", ":"+rs.config.RpcPort)
	if err != nil {
		return err
	}
	rs.grpcServer, err = rs.NewServer()
	if err != nil {
		return err
	}

	RegisterGreeterServer(rs.grpcServer, rs)
	reflection.Register(rs.grpcServer)
	go func() {
		if err := rs.grpcServer.Serve(lis); err != nil {
			log.Error("Rpc startup failed!", "err", err)
			os.Exit(1)
			return
		}
	}()

	rs.httpServer, err = rpchttp.NewRPCServer(rs.config)

	// Register all the APIs exposed by the services
	for _, api := range rs.APIs() {
		if err := rs.httpServer.RegisterService("rpc", api.Service); err != nil {
			return err
		}
	}
	if err := rs.httpServer.Start([]string{endPoint}); err != nil {
		return err
	}

	if rs.config.RpcTLS {
		log.Info("Rpc startup", "port", rs.config.RpcPort, "pem", rs.config.RpcCert)
	} else {
		log.Info("Rpc startup", "port", rs.config.RpcPort)
	}
	return nil
}

func (rs *Server) APIs() []rpchttp.API {
	return []rpchttp.API{
		{
			NameSpace: "rpc",
			Service:   rs.api,
			Public:    true,
		},
	}
}

func (rs *Server) Close() {
	rs.grpcServer.Stop()
	log.Info("GRPC server closed")
}

func (rs *Server) NewServer() (*grpc.Server, error) {
	var opts []grpc.ServerOption
	var interceptor grpc.UnaryServerInterceptor
	interceptor = rs.interceptor
	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	// If tls is configured, generate tls certificate
	if rs.config.RpcTLS {
		if err := rs.generateCertFile(); err != nil {
			return nil, err
		}
		transportCredentials, err := credentials.NewServerTLSFromFile(rs.config.RpcCert, rs.config.RpcCertKey)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.Creds(transportCredentials))

	}

	// Set the maximum number of bytes received and sent
	opts = append(opts, grpc.MaxRecvMsgSize(reqmgr.MaxRequestBytes))
	opts = append(opts, grpc.MaxSendMsgSize(reqmgr.MaxRequestBytes))
	return grpc.NewServer(opts...), nil
}

func (rs *Server) SendTransaction(_ context.Context, req *Bytes) (*Response, error) {
	hash, err := rs.api.SendTransaction(string(req.Bytes))
	if err != nil {
		return NewResponse(rpctypes.RpcErrTxPool, nil, err.Error()), nil
	}
	return NewResponse(rpctypes.RpcSuccess, []byte(fmt.Sprintf("send transaction %s success", hash)), ""), nil
}

func (rs *Server) GetAccount(_ context.Context, req *Address) (*Response, error) {
	account, err := rs.api.GetAccount(req.Address)
	bytes, err := json.Marshal(account)
	if err != nil {
		return NewResponse(rpctypes.RpcErrMarshal, nil, fmt.Sprintf("%s address not exsit", req.Address)), nil
	}
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) GetPendingNonce(_ context.Context, req *Address) (*Response, error) {
	nonce, err := rs.api.GetPendingNonce(req.Address)
	if err != nil {
		return NewResponse(rpctypes.RpcErrMarshal, nil, fmt.Sprintf("%s address not exsit", req.Address)), nil
	}
	return NewResponse(rpctypes.RpcSuccess, []byte(nonce), ""), nil
}

func (rs *Server) GetTransaction(ctx context.Context, req *Hash) (*Response, error) {
	tx, err := rs.api.GetTransaction(req.Hash)
	if err != nil {
		return NewResponse(rpctypes.RpcErrBlockChain, nil, err.Error()), nil
	}
	bytes, _ := json.Marshal(tx)
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) GetBlockByHash(ctx context.Context, req *Hash) (*Response, error) {
	block, err := rs.api.GetBlockByHash(req.Hash)
	if err != nil {
		return NewResponse(rpctypes.RpcErrBlockChain, nil, err.Error()), nil
	}
	bytes, _ := json.Marshal(block)
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) GetBlockByHeight(_ context.Context, req *Height) (*Response, error) {
	block, err := rs.api.GetBlockByHeight(req.Height)
	if err != nil {
		return NewResponse(rpctypes.RpcErrBlockChain, nil, err.Error()), nil
	}
	bytes, _ := json.Marshal(block)
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) GetBlockByRange(_ context.Context, req *Height) (*Response, error) {
	block, err := rs.api.GetBlockByRange(req.Height, req.Count)
	if err != nil {
		return NewResponse(rpctypes.RpcErrBlockChain, nil, err.Error()), nil
	}
	bytes, _ := json.Marshal(block)
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) GetPoolTxs(context.Context, *Null) (*Response, error) {
	txPool, err := rs.api.GetPoolTxs()
	if err != nil {
		return NewResponse(rpctypes.RpcErrTxPool, nil, err.Error()), nil
	}
	bytes, _ := json.Marshal(txPool)
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) GetCandidates(context.Context, *Null) (*Response, error) {
	cans, err := rs.api.GetCandidates()
	if err != nil {
		return NewResponse(rpctypes.RpcErrDPos, nil, err.Error()), nil
	}
	bytes, _ := json.Marshal(cans)
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) GetLastHeight(context.Context, *Null) (*Response, error) {
	height, err := rs.api.GetLastHeight()
	if err != nil {
		return NewResponse(rpctypes.RpcErrBlockChain, nil, err.Error()), nil
	}
	return NewResponse(rpctypes.RpcSuccess, []byte(height), ""), nil
}

func (rs *Server) GetContract(ctx context.Context, req *Address) (*Response, error) {
	contract, err := rs.api.GetContract(req.Address)
	if err != nil {
		return NewResponse(rpctypes.RpcErrContract, nil, err.Error()), nil
	}
	bytes, _ := json.Marshal(contract)
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) GetContractBySymbol(ctx context.Context, req *Symbol) (*Response, error) {
	contract, err := rs.api.GetContractBySymbol(req.Symbol)
	if err != nil {
		return NewResponse(rpctypes.RpcErrContract, nil, err.Error()), nil
	}
	bytes, _ := json.Marshal(contract)
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) GetAddressBySymbol(ctx context.Context, req *Symbol) (*Response, error) {
	address, err := rs.api.GetAddressBySymbol(req.Symbol)
	if err != nil {
		return NewResponse(rpctypes.RpcErrContract, nil, err.Error()), nil
	}
	bytes, _ := json.Marshal(address)
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) TokenList(ctx context.Context, req *Null) (*Response, error) {
	list, err := rs.api.TokenList()
	if err != nil {
		return NewResponse(rpctypes.RpcErrContract, nil, err.Error()), nil
	}
	bytes, _ := json.Marshal(list)
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) AccountList(ctx context.Context, req *Null) (*Response, error) {
	list, err := rs.api.AccountList()
	if err != nil {
		return NewResponse(rpctypes.RpcErrContract, nil, err.Error()), nil
	}
	bytes, _ := json.Marshal(list)
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) ContractMethod(ctx context.Context, req *Method) (*Response, error) {
	result, err := rs.api.ContractMethod(req.Contract, req.Method, req.Params)
	if err != nil {
		return NewResponse(rpctypes.RpcErrContract, nil, err.Error()), nil
	}
	bytes, _ := json.Marshal(result)
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) GetConfirmedHeight(context.Context, *Null) (*Response, error) {
	height, err := rs.api.GetConfirmedHeight()
	if err != nil {
		return NewResponse(rpctypes.RpcErrBlockChain, nil, err.Error()), nil
	}
	return NewResponse(rpctypes.RpcSuccess, []byte(height), ""), nil
}

func (rs *Server) Peers(context.Context, *Null) (*Response, error) {
	peers, err := rs.api.Peers()
	if err != nil {
		return NewResponse(rpctypes.RpcErrBlockChain, nil, err.Error()), nil
	}
	bytes, _ := json.Marshal(peers)
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func (rs *Server) NodeInfo(context.Context, *Null) (*Response, error) {
	nodeInfo, err := rs.api.NodeInfo()
	if err != nil {
		return NewResponse(rpctypes.RpcErrBlockChain, nil, err.Error()), nil
	}
	bytes, _ := json.Marshal(nodeInfo)
	return NewResponse(rpctypes.RpcSuccess, bytes, ""), nil
}

func NewResponse(code int32, result []byte, err string) *Response {
	return &Response{Code: code, Result: result, Err: err}
}

// Authenticate rpc users
func (rs *Server) auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("no token authentication information")
	}
	var (
		password string
	)

	if val, ok := md["password"]; ok {
		password = val[0]
	}

	if password != rs.config.RpcPass {
		return fmt.Errorf("the token authentication information is invalid:password=%s", password)
	}
	return nil
}

func (rs *Server) interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	err = rs.auth(ctx)
	if err != nil {
		return
	}
	return handler(ctx, req)
}

func (rs *Server) generateCertFile() error {
	if rs.config.RpcCert == "" {
		rs.config.RpcCert = rs.config.DataDir + "/server.pem"
	}
	if rs.config.RpcCertKey == "" {
		rs.config.RpcCertKey = rs.config.DataDir + "/server.key"
	}
	if !utils.IsExist(rs.config.RpcCert) || !utils.IsExist(rs.config.RpcCertKey) {
		return certgen.GenCertPair(rs.config.RpcCert, rs.config.RpcCertKey)
	}
	return nil
}
