package exchange

import (
	"fmt"
	"github.com/UBChainNet/UBChain/common/encode/rlp"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/param"
	"math/big"
	"sort"
	"strings"
)

type Pledge struct {
	TotalSupply uint64
	MaxSupply uint64
	DayMint uint64
	Start uint64
	Reward  hasharry.Address
	RewardSymbol string
	Admin hasharry.Address
	LastHeight uint64
	// 每个pair总质押量
	PledgePair map[hasharry.Address]uint64
	// 未成熟的某天某个pair某个地址的质押量
	UnripePledge map[uint64]map[hasharry.Address]map[hasharry.Address]uint64
	// 已经成熟的pair总质押量
	MaturePair map[hasharry.Address]uint64
	// 已经成熟的pair某个地址质押量
	MaturePairAccount map[hasharry.Address]map[hasharry.Address]uint64

	// 每天每个pair池的币量
	PoolReward map[uint64]map[hasharry.Address]uint64
	// 每个账户每个pair池子奖励的币量
	AccountReward map[hasharry.Address]map[hasharry.Address]uint64

}

//const dayBlocks =  60 * 60 * 24 / param.BlockInterval
const dayBlocks =  60 * 2 / param.BlockInterval
//const MatureTime = 10
const MatureTime = 1

func NewPledge(from, exchange hasharry.Address, symbol string, maxSupply, dayMint, height uint64)*Pledge{
	return &Pledge{
		TotalSupply:       0,
		MaxSupply:         maxSupply,
		DayMint:           dayMint,
		Reward:            exchange,
		RewardSymbol: symbol,
		Admin:			   from,
		Start:             height,
		LastHeight:        0,
		AccountReward:     map[hasharry.Address]map[hasharry.Address]uint64{},
		PoolReward:        map[uint64]map[hasharry.Address]uint64{},
		UnripePledge :     map[uint64]map[hasharry.Address]map[hasharry.Address]uint64{},
		PledgePair :       map[hasharry.Address]uint64{},
		MaturePair :       map[hasharry.Address]uint64{},
		MaturePairAccount: map[hasharry.Address]map[hasharry.Address]uint64{},
	}
}

func (p *Pledge)In(height uint64, address hasharry.Address, pair hasharry.Address, amount uint64) error{
	if p.LastHeight == 0{
		p.LastHeight = height
	}
	today := Today(height)
	pledge, exist := p.UnripePledge[today]
	if exist{
		addressPledge, exist := pledge[pair]
		if exist{
			value, exist:= addressPledge[address]
			if exist{
				addressPledge[address] = value+amount
			}else{
				addressPledge[address] = amount
			}
		}else{
			pledge[pair] = map[hasharry.Address]uint64{
				address: amount,
			}
		}
	}else{
		p.UnripePledge[today] = map[hasharry.Address]map[hasharry.Address]uint64{
			pair: {
				address: amount,
			},
		}
	}
	p.PledgePair[pair] = p.PledgePair[pair] + amount
	return nil
}

func (p *Pledge)Out(address hasharry.Address, pair hasharry.Address, amount uint64) error{
	_amount := amount
	//
	if amount == 0{
		return nil
	}
	for _, pairPledge := range p.UnripePledge{
		addrPledge, exist := pairPledge[pair]
		if exist{
			lp := addrPledge[address]
			if lp >= amount{
				lp -= amount
				if lp == 0{
					delete(addrPledge, address)
				}else{
					addrPledge[address] = lp
				}
			}else{
				amount -= lp
				delete(addrPledge, address)
			}
		}
	}

	//
	if amount == 0{
		return nil
	}
	addrPledge, exist := p.MaturePairAccount[pair]
	if !exist{
		return fmt.Errorf("insufficient pledge")
	}
	lp := addrPledge[address]
	if lp >= amount{
		addrPledge[address] -= amount
		if addrPledge[address] == 0{
			delete(addrPledge, address)
		}
	}else{
		return fmt.Errorf("insufficient pledge")
	}
	if len(addrPledge) == 0{
		delete(p.MaturePairAccount, pair)
	}

	//
	matureTotal := p.MaturePair[pair]
	if matureTotal < amount{
		return fmt.Errorf("insufficient pledge")
	}
	matureTotal -= amount
	if matureTotal == 0{
		delete(p.MaturePair, pair)
	}else{
		p.MaturePair[pair] = matureTotal
	}

	//
	pairTotal := p.PledgePair[pair]
	if pairTotal < _amount{
		return fmt.Errorf("insufficient pledge")
	}
	pairTotal -= _amount
	if pairTotal == 0{
		delete(p.PledgePair, pair)
	}else{
		p.PledgePair[pair] = pairTotal
	}


	return nil
}

type Reward struct {
	Token hasharry.Address
	Amount uint64
}

func (p *Pledge)RemoveReward(address hasharry.Address)[]Reward{
	rewardList := []Reward{}
	rewards, exist := p.AccountReward[address]
	if exist{
		delete(p.AccountReward, address)
		for pair, amount := range rewards{
			rewardList = append(rewardList, Reward{
				Token:  pair,
				Amount: amount,
			})
		}
		sort.Slice(rewardList, func(i, j int) bool {
			return strings.Compare(rewardList[i].Token.String(), rewardList[j].Token.String()) < 0
		})
		return rewardList
	}else{
		return rewardList
	}
}

func (p *Pledge)GetPledgeAmount(address, pair hasharry.Address)(unripeLp uint64, matureLp uint64){
	for _, pairPledge := range p.UnripePledge{
		addrPledge, exist := pairPledge[pair]
		if exist{
			unripeLp = addrPledge[address]
			break
		}
	}
	matureLp = 0
	addrPledge, exist := p.MaturePairAccount[pair]
	if !exist{
		return
	}
	matureLp = addrPledge[address]
	return
}


func (p *Pledge)IsUpdate(height uint64) bool {
	if p.TotalSupply >= p.MaxSupply{
		return false
	}
	if p.LastHeight == 0{
		return false
	}
	lastDay := p.LastHeight / dayBlocks
	today := Today(height)
	if today >= lastDay + 1{
		return true
	}
	return false
}

func (p *Pledge)UpdateMature(height uint64){
	today := Today(height)
	for day, pairPledge := range p.UnripePledge{
		if today - day  >= MatureTime{
			for pair, pledge := range pairPledge{
				maturePair, exist := p.MaturePairAccount[pair]
				if exist{
					for address, amount := range pledge{
						maturePair[address] = maturePair[address] + amount
					}
				}else{
					p.MaturePairAccount[pair] = pledge
				}
				for _, amount := range pledge{
					p.MaturePair[pair] = p.MaturePair[pair] + amount
				}
			}
			delete(p.UnripePledge, day)
		}
	}
}

func (p *Pledge)UpdateReward(height, totalValue uint64, pairValue map[hasharry.Address]uint64){
	lastDay := p.LastHeight / dayBlocks
	today := Today(height)
	dayMint := p.DayMint

	for i := lastDay + 1;i <= today;i++{
		if p.TotalSupply >= p.MaxSupply{
			break
		}
		if p.TotalSupply + dayMint >= p.MaxSupply{
			dayMint = p.MaxSupply - p.TotalSupply
		}
		poolReward := map[hasharry.Address]uint64{}
		p.PoolReward[i] = poolReward
		var totalReward uint64
		for pairAddr, value := range pairValue{
			reward := big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(int64(value)),  big.NewInt(int64(dayMint))) , big.NewInt(int64(totalValue))).Uint64()
			poolReward[pairAddr] = reward
			totalReward += reward
		}
		for pairAddr, addrMature := range p.MaturePairAccount{
			pairTotalReward := poolReward[pairAddr]
			for addr, lp := range addrMature{
				totalLp := p.MaturePair[pairAddr]
				if totalLp != 0{
					reward := big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(int64(lp)),  big.NewInt(int64(pairTotalReward))) , big.NewInt(int64(totalLp))).Uint64()
					pairReward, exist := p.AccountReward[addr]
					if exist{
						pairReward[pairAddr] =  pairReward[pairAddr] + reward
					}else{
						p.AccountReward[addr] = map[hasharry.Address]uint64{
							pairAddr: reward,
						}
					}
				}
			}
		}
		p.TotalSupply += totalReward
	}
	p.LastHeight = height
}



func (p *Pledge) GetPledgeReward(address, pair hasharry.Address) uint64{
	pairReward, exist := p.AccountReward[address]
	if !exist{
		return 0
	}
	return pairReward[pair]
}

func (p *Pledge) GetPledgeRewards(address hasharry.Address) map[hasharry.Address]uint64{
	return p.AccountReward[address]
}

func (p *Pledge) GetMaturePledge(address, pair hasharry.Address) uint64 {
	addrReward, exist := p.MaturePairAccount[pair]
	if !exist{
		return 0
	}
	return  addrReward[address]
}

func (p *Pledge) GetUnripePledge(address, pair hasharry.Address) uint64 {
	for _, unripePledge := range p.UnripePledge {
		addrReward, exist := unripePledge[pair]
		if !exist{
			return 0
		}
		return addrReward[address]
	}
	return 0
}

func (p *Pledge) Bytes()[]byte {
	rlpPledge := &RlpPledge{
		TotalSupply:       p.TotalSupply,
		MaxSupply:         p.MaxSupply,
		DayMint:           p.DayMint,
		Start:             p.Start,
		Reward:            p.Reward,
		RewardSymbol:     p.RewardSymbol ,
		Admin:             p.Admin,
		LastHeight:        p.LastHeight,
		PledgePair:        nil,
		UnripePledge:      nil,
		MaturePair:        nil,
		MaturePairAccount: nil,
		PoolReward:        nil,
		AccountReward:     nil,
	}
	rlpPledge.PledgePair = make([]PledgePair, 0)
	for pair, total := range p.PledgePair{
		rlpPledge.PledgePair = append(rlpPledge.PledgePair, PledgePair{
			Pair:   pair,
			Amount: total,
		})
	}
	sort.Slice(rlpPledge.PledgePair , func(i, j int) bool {
		return strings.Compare(rlpPledge.PledgePair[i].Pair.String(), rlpPledge.PledgePair[j].Pair.String()) < 0
	})

	rlpPledge.UnripePledge = make([]UnripePledge, 0)
	for day, pairPledge := range p.UnripePledge{
		for pair, addrPledge := range pairPledge{
			for addr, amount := range addrPledge{
				rlpPledge.UnripePledge = append(rlpPledge.UnripePledge, UnripePledge{
					Day:    day,
					Pair:   pair,
					Address: addr,
					Amount: amount,
				})
			}
		}

	}
	sort.Slice(rlpPledge.UnripePledge , func(i, j int) bool {
		pairCp := strings.Compare(rlpPledge.UnripePledge[i].Pair.String(), rlpPledge.UnripePledge[j].Pair.String())
		if  pairCp == 0{
			return strings.Compare(rlpPledge.UnripePledge[i].Address.String(), rlpPledge.UnripePledge[j].Address.String()) < 0
		}else{
			return pairCp < 0
		}
	})

	rlpPledge.MaturePair = make([]PledgePair, 0)
	for pair, total := range p.MaturePair{
		rlpPledge.MaturePair = append(rlpPledge.MaturePair, PledgePair{
			Pair:   pair,
			Amount: total,
		})

	}
	sort.Slice(rlpPledge.MaturePair , func(i, j int) bool {
		return strings.Compare(rlpPledge.MaturePair[i].Pair.String(), rlpPledge.MaturePair[j].Pair.String()) < 0
	})

	rlpPledge.MaturePairAccount = make([]MaturePledge, 0)
	for pair, addrValue := range p.MaturePairAccount{
		for addr, value := range addrValue{
			rlpPledge.MaturePairAccount = append(rlpPledge.MaturePairAccount, MaturePledge{
				Address: addr,
				Pair:    pair,
				Amount:  value,
			})

		}
	}
	sort.Slice(rlpPledge.MaturePairAccount , func(i, j int) bool {
		pairCp := strings.Compare(rlpPledge.MaturePairAccount[i].Pair.String(), rlpPledge.MaturePairAccount[j].Pair.String())
		if  pairCp == 0{
			return strings.Compare(rlpPledge.MaturePairAccount[i].Address.String(), rlpPledge.MaturePairAccount[j].Address.String()) < 0
		}else{
			return pairCp < 0
		}
	})

	rlpPledge.PoolReward = make([]PoolReward, 0)
	for day, pairReward := range p.PoolReward{
		for pair, reward := range pairReward{
			rlpPledge.PoolReward = append(rlpPledge.PoolReward, PoolReward{
				Day:    day,
				Pair:   pair,
				Reward: reward,
			})

		}
	}
	sort.Slice(rlpPledge.PoolReward , func(i, j int) bool {
		if rlpPledge.PoolReward[i].Day == rlpPledge.PoolReward[j].Day{
			return strings.Compare(rlpPledge.PoolReward[i].Pair.String(), rlpPledge.PoolReward[j].Pair.String()) < 0
		}else{
			return rlpPledge.PoolReward[i].Day < rlpPledge.PoolReward[j].Day
		}
	})

	rlpPledge.AccountReward = make([]AccountReward, 0)
	for addr, pairReward := range p.AccountReward{
		for pair, reward := range pairReward{
			rlpPledge.AccountReward = append(rlpPledge.AccountReward, AccountReward{
				Address:    addr,
				Pair:   pair,
				Reward: reward,
			})

		}
	}
	sort.Slice(rlpPledge.AccountReward , func(i, j int) bool {
		addrCp := strings.Compare(rlpPledge.AccountReward[i].Address.String(), rlpPledge.AccountReward[j].Address.String())
		if  addrCp == 0{
			return strings.Compare(rlpPledge.AccountReward[i].Pair.String(), rlpPledge.AccountReward[j].Pair.String()) < 0
		}else{
			return addrCp < 0
		}
	})

	bytes, _ := rlp.EncodeToBytes(rlpPledge)
	return bytes
}

type PledgePair struct {
	Pair hasharry.Address
	Amount uint64
}

type UnripePledge struct {
	Day uint64
	Pair hasharry.Address
	Address hasharry.Address
	Amount uint64
}

type MaturePledge struct {
	Address hasharry.Address
	Pair hasharry.Address
	Amount uint64
}

type PoolReward struct {
	Day uint64
	Pair hasharry.Address
	Reward uint64
}

type AccountReward struct {
	Address hasharry.Address
	Pair hasharry.Address
	Reward uint64
}

type RlpPledge struct {
	TotalSupply uint64
	MaxSupply uint64
	DayMint uint64
	Start uint64
	Reward  hasharry.Address
	RewardSymbol string
	Admin hasharry.Address
	LastHeight uint64
	// 每个pair总质押量
	PledgePair []PledgePair
	// 未成熟的某天某个pair某个地址的质押量
	UnripePledge []UnripePledge
	// 已经成熟的pair总质押量
	MaturePair []PledgePair
	// 已经成熟的pair某个地址质押量
	MaturePairAccount []MaturePledge

	// 每天每个pair池的币量
	PoolReward []PoolReward
	// 每个账户每个pair池子奖励的币量
	AccountReward []AccountReward
}

func DecodeToPledge(bytes []byte) (*Pledge, error) {
	var rlpPd *RlpPledge
	err := rlp.DecodeBytes(bytes, &rlpPd)
	pd := &Pledge{
		TotalSupply:       rlpPd.TotalSupply,
		MaxSupply:         rlpPd.MaxSupply,
		DayMint:           rlpPd.DayMint,
		Start:             rlpPd.Start,
		Reward:            rlpPd.Reward,
		RewardSymbol:      rlpPd.RewardSymbol,
		Admin:             rlpPd.Admin,
		LastHeight:        rlpPd.LastHeight,
		PledgePair:        map[hasharry.Address]uint64{},
		UnripePledge: 	   map[uint64]map[hasharry.Address]map[hasharry.Address]uint64{},
		MaturePair: 	   map[hasharry.Address]uint64{},
		MaturePairAccount: map[hasharry.Address]map[hasharry.Address]uint64{},
		PoolReward: 	   map[uint64]map[hasharry.Address]uint64{},
		AccountReward: 	   map[hasharry.Address]map[hasharry.Address]uint64{},
	}
	for _, pledgePair := range rlpPd.PledgePair{
		pd.PledgePair[pledgePair.Pair] = pledgePair.Amount
	}
	for _, unripePledge := range rlpPd.UnripePledge{
		pairPledge, exist := pd.UnripePledge[unripePledge.Day]
		if exist{
			addrPledge, exist := pairPledge[unripePledge.Pair]
			if exist{
				addrPledge[unripePledge.Address] = unripePledge.Amount
			}else{
				pairPledge[unripePledge.Pair] = map[hasharry.Address]uint64{
					unripePledge.Address : unripePledge.Amount,
				}
			}
		}else{
			pd.UnripePledge[unripePledge.Day] = map[hasharry.Address]map[hasharry.Address]uint64{
				unripePledge.Pair : {
					unripePledge.Address : unripePledge.Amount,
				},
			}
		}
	}

	for _, maturePair := range rlpPd.MaturePair{
		pd.MaturePair[maturePair.Pair] = maturePair.Amount
	}

	for _, maturePledge := range rlpPd.MaturePairAccount{
		addrPledge, exist := pd.MaturePairAccount[maturePledge.Pair]
		if exist{
			addrPledge[maturePledge.Address] =  maturePledge.Amount
		}else{
			pd.MaturePairAccount[maturePledge.Pair] = map[hasharry.Address]uint64{
				maturePledge.Address : maturePledge.Amount,
			}
		}
	}

	for _, poolReward := range rlpPd.PoolReward{
		pairReward, exist := pd.PoolReward[poolReward.Day]
		if exist{
			pairReward[poolReward.Pair] = poolReward.Reward
		}else{
			pd.PoolReward[poolReward.Day] = map[hasharry.Address]uint64{
				poolReward.Pair : poolReward.Reward,
			}
		}
	}

	for _, accountReward := range rlpPd.AccountReward{
		pairReward, exist := pd.AccountReward[accountReward.Address]
		if exist{
			pairReward[accountReward.Pair] =  accountReward.Reward
		}else{
			pd.AccountReward[accountReward.Address] = map[hasharry.Address]uint64{
				accountReward.Pair : accountReward.Reward,
			}
		}
	}

	return pd, err
}

func Today(height uint64)uint64{
	return  height / dayBlocks
}