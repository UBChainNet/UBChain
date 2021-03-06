package types

import (
	"github.com/UBChainNet/UBChain/core/types/contractv2"
)

type RpcContractV2TransactionBody struct {
	Contract     string                  `json:"contract"`
	Type         contractv2.ContractType `json:"type"`
	FunctionType contractv2.FunctionType `json:"functiontype"`
	Function     IRCFunction             `json:"function"`
}

type RpcContractV2BodyWithState struct {
	Contract     string                  `json:"contract"`
	Type         contractv2.ContractType `json:"type"`
	FunctionType contractv2.FunctionType `json:"functiontype"`
	Function     IRCFunction             `json:"function"`
	State        *RpcContractState       `json:"state"`
}

type RpcContractState struct {
	StateCode ContractState `json:"statecode"`
	Events    []*RpcEvent   `json:"event"`
	Error     string        `json:"error"`
}

type RpcEvent struct {
	EventType int    `json:"eventtype"`
	From      string `json:"from"`
	To        string `json:"to"`
	Token     string `json:"token"`
	Amount    uint64 `json:"amount"`
	Height    uint64 `json:"height"`
}

type IRCFunction interface {
}

type RpcExchangeInitBody struct {
	Symbol string `json:"symbol"`
	Admin  string `json:"admin"`
	FeeTo  string `json:"feeto"`
}

type RpcExchangeSetAdminBody struct {
	Address string `json:"address"`
}

type RpcExchangeSetFeeToBody struct {
	Address string `json:"address"`
}

type RpcExchangeExactInBody struct {
	AmountIn     uint64   `json:"amountin"`
	AmountOutMin uint64   `json:"amountoutmin"`
	Path         []string `json:"path"`
	To           string   `json:"to"`
	Deadline     uint64   `json:"deadline"`
}

type RpcExchangeExactOutBody struct {
	AmountOut   uint64   `json:"amountout"`
	AmountInMax uint64   `json:"amountinmax"`
	Path        []string `json:"path"`
	To          string   `json:"to"`
	Deadline    uint64   `json:"deadline"`
}

type RpcExchangeAddLiquidity struct {
	Exchange       string `json:"exchange"`
	TokenA         string `json:"tokenA"`
	TokenB         string `json:"tokenB"`
	To             string `json:"to"`
	AmountADesired uint64 `json:"amountadesired"`
	AmountBDesired uint64 `json:"amountbdesired"`
	AmountAMin     uint64 `json:"amountamin"`
	AmountBMin     uint64 `json:"amountbmin"`
	Deadline       uint64 `json:"deadline"`
}

type RpcExchangeRemoveLiquidity struct {
	Exchange   string `json:"exchange"`
	TokenA     string `json:"tokenA"`
	TokenB     string `json:"tokenB"`
	To         string `json:"to"`
	Liquidity  uint64 `json:"liquidity"`
	AmountAMin uint64 `json:"amountamin"`
	AmountBMin uint64 `json:"amountbmin"`
	Deadline   uint64 `json:"deadline"`
}

type RpcPledgeInit struct {
	Exchange  string `json:"exchange"`
	Receiver  string `json:"receiver"`
	Admin     string `json:"admin"`
	PreMint   uint64 `json:"premint"`
	MaxSupply uint64 `json:"maxsupply"`
}

type RpcPledgeStart struct {
	BlockMintAmount  uint64 `json:"blockmintamount"`
	PledgeMatureTime uint64 `json:"pledgematuretime"`
}

type RpcPledgeAddPool struct {
	Pair        string `json:"pair"`
	BlockReward uint64 `json:"blockreward"`
}

type RpcPledgeRemovePool struct {
	Pair string `json:"pair"`
}

type RpcPledgeAdd struct {
	Pair   string `json:"pair"`
	Amount uint64 `json:"amount"`
}

type RpcPledgeRemove struct {
	Pair   string `json:"pair"`
	Amount uint64 `json:"amount"`
}

type RpcPledgeRewardRemove struct {
}

type RpcPledgeUpdate struct {
}

type RpcTokenHubInit struct {
	Setter  string `json:"setter"`
	Admin   string `json:"admin"`
	FeeTo   string `json:"feeto"`
	FeeRate string `json:"feerate"`
}

type RpcTokenHubAck struct {
	Sequences []uint64 `json:"sequences"`
	AckTypes  []uint8  `json:"acktypes"`
	Hashes    []string `json:"hashes"`
}

type RpcTokenHubTransferOut struct {
	To     string `json:"to"`
	Amount uint64 `json:"amount"`
}

type RpcTokenHubTransferIn struct {
	To        string `json:"to"`
	Amount    uint64 `json:"amount"`
	AcrossSeq uint64 `json:"acrossseq"`
}

type RpcTokenHubFinishAcrossBody struct {
	AcrossSeqs []uint64 `json:"acrossseqs"`
}
