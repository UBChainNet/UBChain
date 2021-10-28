package exchange

import (
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/core/runner/library"
	"github.com/UBChainNet/UBChain/param"
)

type Pledge struct {
	library library.RunnerLibrary
	Start uint64
	Reward  hasharry.Address
	Symbol  string
	LastHeight uint64
	PledgeAddress map[string]map[hasharry.Address]uint64
	MinerPool map[uint64]map[hasharry.Address]float64
}

const dayBlocks =  60 * 60 * 24 / param.BlockInterval

func (p *Pledge)Update(height uint64){
	if height < p.Start{
		return
	}
	if p.LastHeight == 0{
		p.LastHeight = p.Start
	}
	lastDay := p.LastHeight / dayBlocks
	nextHeight := (lastDay + 1) * dayBlocks
	if height < nextHeight{
		return
	}
	curDay := height / dayBlocks
	pool, exist := p.MinerPool[curDay]
	if exist{
		return
	}
	exchange, err := p.library.GetExchange(p.Reward)
	if err != nil{
		return
	}
	for _, token1AndPair := range exchange.Pair{
		for _, pairAddr := range token1AndPair{

		}
	}
	p.MinerPool[curDay] =
	for i := lastDay+1;i <= curDay;i++{

	}
}

func (p *Pledge)day