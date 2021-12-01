package types

import (
	"github.com/UBChainNet/UBChain/core/types/contractv2"
	"github.com/UBChainNet/UBChain/core/types/contractv2/exchange"
	"github.com/UBChainNet/UBChain/core/types/contractv2/tokenhub"
)

type RpcPair struct {
	Address              string  `json:"address"`
	Exchange             string  `json:"exchange"`
	Symbol               string  `json:"symbol"`
	Token0               string  `json:"token0"`
	Token1               string  `json:"token1"`
	Symbol0              string  `json:"symbol0"`
	Symbol1              string  `json:"symbol1"`
	Reserve0             float64 `json:"reserve0"`
	Reserve1             float64 `json:"reserve1"`
	TotalSupply          float64 `json:"totalSupply"`
	BlockTimestampLast   uint32  `json:"blockTimestampLast"`
	Price0CumulativeLast uint64  `json:"price0CumulativeLast"`
	Price1CumulativeLast uint64  `json:"price1CumulativeLast"`
	KLast                string  `json:"kLast"`
}

type RpcExchange struct {
	Address string                       `json:"address"`
	Symbol  string                       `json:"symbol"`
	FeeTo   string                       `json:"feeTo"`
	Admin   string                       `json:"admin"`
	Pair    map[string]map[string]string `json:"pair"`
}

type RpcPledge struct {
	PreMint          uint64 `json:"preMint"`
	DayMintAmount    uint64 `json:"dayMintAmount"`
	Receiver         string `json:"receiver"`
	TotalSupply      uint64 `json:"totalSupply"`
	MaxSupply        uint64 `json:"maxSupply"`
	PledgeMatureTime uint64 `json:"pledgeMaturetime"`
	DayRewardAmount  uint64 `json:"dayRewardAmount"`
	Start            uint64 `json:"start"`
	RewardToken      string `json:"rewardToken"`
	RewardSymbol     string `json:"rewardSymbol"`
	Admin            string `json:"admin"`
	LastHeight       uint64 `json:"lastHeight"`
}

type PairAddress struct {
	Key     string `json:"key"`
	Address string `json:"address"`
}

func TranslateContractV2ToRpcContractV2(contract *contractv2.ContractV2) interface{} {
	switch contract.Type {
	case contractv2.Exchange_:
		exchange, _ := contract.Body.(*exchange.Exchange)
		pair := map[string]map[string]string{}
		for token0, token1AndAddr := range exchange.Pair {
			addressMap := map[string]string{}
			for token1, address := range token1AndAddr {
				addressMap[token1.String()] = address.String()
			}
			pair[token0.String()] = addressMap
		}
		return &RpcExchange{
			Address: contract.Address.String(),
			FeeTo:   exchange.FeeTo.String(),
			Admin:   exchange.Admin.String(),
			Symbol:  exchange.Symbol,
			Pair:    pair,
		}
	case contractv2.Pair_:
		pair, _ := contract.Body.(*exchange.Pair)
		return &RpcPair{
			Address:              contract.Address.String(),
			Exchange:             pair.Exchange.String(),
			Symbol:               pair.Symbol,
			Token0:               pair.Token0.String(),
			Token1:               pair.Token1.String(),
			Symbol0:              pair.Symbol0,
			Symbol1:              pair.Symbol1,
			Reserve0:             Amount(pair.Reserve0).ToCoin(),
			Reserve1:             Amount(pair.Reserve1).ToCoin(),
			BlockTimestampLast:   pair.BlockTimestampLast,
			Price0CumulativeLast: pair.Price0CumulativeLast,
			Price1CumulativeLast: pair.Price1CumulativeLast,
			KLast:                pair.KLast.String(),
			TotalSupply:          Amount(pair.TotalSupply).ToCoin(),
		}
	case contractv2.Pledge_:
		pledge, _ := contract.Body.(*exchange.Pledge)
		return &RpcPledge{
			PreMint:          pledge.PreMint,
			DayMintAmount:    pledge.BlockMintAmount,
			Receiver:         pledge.Receiver.String(),
			TotalSupply:      pledge.TotalSupply,
			MaxSupply:        pledge.MaxSupply,
			PledgeMatureTime: pledge.PledgeMatureTime,
			DayRewardAmount:  pledge.DayRewardAmount,
			Start:            pledge.Start,
			RewardToken:      pledge.RewardToken.String(),
			RewardSymbol:     pledge.RewardSymbol,
			Admin:            pledge.Admin.String(),
			LastHeight:       pledge.LastHeight,
		}
	case contractv2.TokenHub_:
		th, _ := contract.Body.(*tokenhub.TokenHub)
		thTrs := make(map[uint64]*TokenHubTransfer, 0)
		for _, tr := range th.Transfers {
			thTrs[tr.Sequence] = &TokenHubTransfer{
				Sequence: tr.Sequence,
				From:     tr.From,
				To:       tr.To,
				Amount:   Amount(tr.Amount).ToCoin(),
				Fees:     Amount(tr.Fees).ToCoin(),
			}
		}
		unTrs := make(map[uint64]*TokenHubTransfer, 0)
		for _, tr := range th.Unconfirmed {
			unTrs[tr.Sequence] = &TokenHubTransfer{
				Sequence: tr.Sequence,
				From:     tr.From,
				To:       tr.To,
				Amount:   Amount(tr.Amount).ToCoin(),
				Fees:     Amount(tr.Fees).ToCoin(),
			}
		}
		return &RpcTokenHub{
			Address:     th.Address.String(),
			Setter:      th.Setter.String(),
			Admin:       th.Admin.String(),
			FeeTo:       th.FeeTo.String(),
			FeeRate:     th.FeeRate,
			Transfers:   thTrs,
			Unconfirmed: unTrs,
			Sequence:    th.Sequence,
		}
	}
	return nil
}

type TokenHubTransfer struct {
	Sequence uint64  `json:"sequence"`
	From     string  `json:"from"`
	To       string  `json:"to"`
	Amount   float64 `json:"amount"`
	Fees     float64 `json:"fees"`
}

type RpcTokenHub struct {
	Address     string                       `json:"address"`
	Setter      string                       `json:"setter"`
	Admin       string                       `json:"admin"`
	FeeTo       string                       `json:"feeTo"`
	FeeRate     float64                      `json:"feeRate"`
	Transfers   map[uint64]*TokenHubTransfer `json:"transfers"`
	Unconfirmed map[uint64]*TokenHubTransfer `json:"unconfirmed"`
	Sequence    uint64
}
