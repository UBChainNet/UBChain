package exchange

import (
	"errors"
	"fmt"
	"github.com/UBChainNet/UBChain/common/encode/rlp"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/param"
	"math/big"
	"sort"
	"strings"
)

type Pledge struct {
	// 预铸币
	PreMint uint64
	// 一次性铸币量
	BlockMintAmount uint64
	// 铸币接收地址
	Receiver hasharry.Address
	// 总铸币量
	TotalSupply uint64
	// 最大供应量
	MaxSupply uint64
	// 质押成熟时间
	PledgeMatureTime uint64
	// 每天质押奖励(弃用)
	DayRewardAmount uint64
	// 开始高度
	Start uint64
	// 奖励token
	RewardToken hasharry.Address
	// 奖励token symbol
	RewardSymbol string
	// 管理员
	Admin hasharry.Address
	// 最后高度
	LastHeight uint64
	// 配对池(弃用)
	PairPool map[hasharry.Address]bool
	// 池子产币量
	PairPoolWithCount map[hasharry.Address]uint64
	// 每个pair总质押量
	PledgePair map[hasharry.Address]uint64
	// 未成熟的某天某个pair某个地址的质押量
	UnripePledge map[uint64]map[hasharry.Address]map[hasharry.Address]uint64
	// 已经成熟的pair总质押量
	MaturePair map[hasharry.Address]uint64
	// 已经成熟的pair某个地址质押量
	MaturePairAccount map[hasharry.Address]map[hasharry.Address]uint64
	// 被删除的pair的某个地址的质押量，可取回
	DeletedPairAccount map[hasharry.Address]map[hasharry.Address]uint64

	// 每天每个pair池的币量(弃用)
	PoolReward map[uint64]map[hasharry.Address]uint64
	// 每个账户每个pair池子奖励的币量
	AccountReward map[hasharry.Address]map[hasharry.Address]uint64
}

const dayBlocks =  60 * 60 * 24 / param.BlockInterval
//const dayBlocks = 60 * 2 / param.BlockInterval

func NewPledge(exchange, ReceiveAddress, admin hasharry.Address, symbol string, maxSupply,
	preMint, height uint64) *Pledge {
	return &Pledge{
		PreMint:          preMint,
		BlockMintAmount:  0,
		Receiver:         ReceiveAddress,
		TotalSupply:      preMint,
		MaxSupply:        maxSupply,
		DayRewardAmount:  0,
		RewardToken:      exchange,
		RewardSymbol:     symbol,
		PledgeMatureTime: 0,
		Admin:            admin,
		Start:            0,
		LastHeight:       0,
		PairPool:           map[hasharry.Address]bool{},
		PairPoolWithCount:  map[hasharry.Address]uint64{},
		AccountReward:      map[hasharry.Address]map[hasharry.Address]uint64{},
		PoolReward:         map[uint64]map[hasharry.Address]uint64{},
		UnripePledge:       map[uint64]map[hasharry.Address]map[hasharry.Address]uint64{},
		PledgePair:         map[hasharry.Address]uint64{},
		MaturePair:         map[hasharry.Address]uint64{},
		MaturePairAccount:  map[hasharry.Address]map[hasharry.Address]uint64{},
		DeletedPairAccount: map[hasharry.Address]map[hasharry.Address]uint64{},
	}
}

func (p *Pledge) AddPirPool(pair hasharry.Address, blockReward, height uint64) error {
	p.PairPoolWithCount[pair] = blockReward
	if addrPledge, exist := p.DeletedPairAccount[pair]; exist{
		for addr, pledge := range addrPledge{
			if err := p.In(height, addr, pair, pledge); err != nil{
				return err
			}
		}
		delete(p.DeletedPairAccount, pair)
	}
	return nil
}

func (p *Pledge) ExistPairPool(pair hasharry.Address) bool {
	_, exist := p.PairPoolWithCount[pair]
	return exist
}

func (p *Pledge) In(height uint64, address hasharry.Address, pair hasharry.Address, amount uint64) error {
	if p.Start == 0 || height <= p.Start{
		return errors.New("it hasn't started yet")
	}
	if _, exist := p.PairPoolWithCount[pair]; !exist{
		return errors.New("the pair was not found")
	}

	current := height
	pledge, exist := p.UnripePledge[current]
	if exist {
		addressPledge, exist := pledge[pair]
		if exist {
			value, exist := addressPledge[address]
			if exist {
				addressPledge[address] = value + amount
			} else {
				addressPledge[address] = amount
			}
		} else {
			pledge[pair] = map[hasharry.Address]uint64{
				address: amount,
			}
		}
	} else {
		p.UnripePledge[current] = map[hasharry.Address]map[hasharry.Address]uint64{
			pair: {
				address: amount,
			},
		}
	}
	p.PledgePair[pair] = p.PledgePair[pair] + amount
	return nil
}

func (p *Pledge) Out(address hasharry.Address, pair hasharry.Address, amount uint64) error {
	_amount := amount
	if amount == 0 {
		return nil
	}

	if addrPledge, exist := p.DeletedPairAccount[pair];exist{
		if addrPledge[address] >= amount{
			addrPledge[address] -= amount
			return nil
		}else{
			return errors.New("insufficient pledge")
		}
	}

	for _, pairPledge := range p.UnripePledge {
		addrPledge, exist := pairPledge[pair]
		if exist {
			lp := addrPledge[address]
			if lp >= amount {
				lp -= amount
				if lp == 0 {
					delete(addrPledge, address)
				} else {
					addrPledge[address] = lp
				}
			} else {
				amount -= lp
				delete(addrPledge, address)
			}
		}
	}

	//
	if amount == 0 {
		return nil
	}
	addrPledge, exist := p.MaturePairAccount[pair]
	if !exist {
		return fmt.Errorf("insufficient pledge")
	}
	lp := addrPledge[address]
	if lp >= amount {
		addrPledge[address] -= amount
		if addrPledge[address] == 0 {
			delete(addrPledge, address)
		}
	} else {
		return fmt.Errorf("insufficient pledge")
	}
	if len(addrPledge) == 0 {
		delete(p.MaturePairAccount, pair)
	}

	//
	matureTotal := p.MaturePair[pair]
	if matureTotal < amount {
		return fmt.Errorf("insufficient pledge")
	}
	matureTotal -= amount
	if matureTotal == 0 {
		delete(p.MaturePair, pair)
	} else {
		p.MaturePair[pair] = matureTotal
	}

	//
	pairTotal := p.PledgePair[pair]
	if pairTotal < _amount {
		return fmt.Errorf("insufficient pledge")
	}
	pairTotal -= _amount
	if pairTotal == 0 {
		delete(p.PledgePair, pair)
	} else {
		p.PledgePair[pair] = pairTotal
	}

	return nil
}

func (p *Pledge) SetStart(height, blockMintAmount, matureTime uint64) error {
	if p.Start == 0{
		p.Start = height
		p.LastHeight = height
	}

	p.BlockMintAmount = blockMintAmount
	p.PledgeMatureTime = matureTime
	return nil
}

func (p *Pledge) RemovePairPool(pair hasharry.Address) error {
	_, exist := p.PairPool[pair]
	if !exist{
		return errors.New("pair does not exist")
	}
	delete(p.PairPool, pair)
	delete(p.MaturePair, pair)

	for _, pairPledge := range p.UnripePledge{
		if addrPledge, exist := pairPledge[pair]; exist{
			if deleted, exist := p.DeletedPairAccount[pair]; exist{
				for addr, pledge := range addrPledge{
					deleted[addr] += pledge
				}
			}else{
				p.DeletedPairAccount[pair] = addrPledge
			}
			delete(pairPledge, pair)
		}
	}

	if addrPledge, exist := p.MaturePairAccount[pair]; exist{
		if deleted, exist := p.DeletedPairAccount[pair]; exist{
			for addr, pledge := range addrPledge{
				deleted[addr] += pledge
			}
		}else{
			p.DeletedPairAccount[pair] = addrPledge
		}
		delete(p.MaturePairAccount, pair)
	}


	return nil
}

type Reward struct {
	Token  hasharry.Address
	Amount uint64
}

func (p *Pledge) RemoveReward(address hasharry.Address) []Reward {
	rewardList := []Reward{}
	rewards, exist := p.AccountReward[address]
	if exist {
		delete(p.AccountReward, address)
		for pair, amount := range rewards {
			rewardList = append(rewardList, Reward{
				Token:  pair,
				Amount: amount,
			})
		}
		sort.Slice(rewardList, func(i, j int) bool {
			return strings.Compare(rewardList[i].Token.String(), rewardList[j].Token.String()) < 0
		})
		return rewardList
	} else {
		return rewardList
	}
}

func (p *Pledge) GetPledgeAmount(address, pair hasharry.Address) (unripeLp uint64, matureLp uint64, deleteLp uint64) {
	for _, pairPledge := range p.UnripePledge {
		addrPledge, exist := pairPledge[pair]
		if exist {
			unripeLp = addrPledge[address]
			break
		}
	}
	matureLp = 0
	addrPledge, exist := p.MaturePairAccount[pair]
	if !exist {
		return
	}
	matureLp = addrPledge[address]

	deleteLp = 0
	addrPledge, exist = p.DeletedPairAccount[pair]
	if !exist {
		return
	}
	deleteLp = addrPledge[address]
	return
}

func (p *Pledge) IsUpdate(height uint64) bool {
	if p.Start == 0 || height <= p.Start{
		return false
	}
	if p.TotalSupply >= p.MaxSupply {
		return false
	}
	if  p.LastHeight < height {
		return true
	}
	return false
}

func (p *Pledge) UpdateMature(height uint64) {
	if p.Start == 0 || height <= p.Start{
		return
	}
	current := height
	for day, pairPledge := range p.UnripePledge {
		if current-day >= p.PledgeMatureTime {
			for pair, pledge := range pairPledge {
				maturePair, exist := p.MaturePairAccount[pair]
				if exist {
					for address, amount := range pledge {
						maturePair[address] = maturePair[address] + amount
					}
				} else {
					p.MaturePairAccount[pair] = pledge
				}
				for _, amount := range pledge {
					p.MaturePair[pair] = p.MaturePair[pair] + amount
				}
			}
			delete(p.UnripePledge, day)
		}
	}
}

func (p *Pledge) UpdateMint(current uint64) uint64 {
	if p.Start == 0 || current <= p.Start{
		return 0
	}
	if p.TotalSupply >= p.MaxSupply {
		return 0
	}

	var allPoolReward, allBlockMint, blockMint, actualMint  uint64
	blockMint = p.BlockMintAmount
	blocks := current - p.LastHeight
	allBlockMint = blockMint * blocks

	for _, reward := range p.PairPoolWithCount{
		allPoolReward += reward * blocks
	}

	if allPoolReward + allBlockMint + p.TotalSupply > p.MaxSupply{
		for i := p.LastHeight + 1; i <= current; i++ {
			if p.TotalSupply >= p.MaxSupply {
				return 0
			}
			if p.TotalSupply + blockMint >= p.MaxSupply {
				blockMint = p.MaxSupply - p.TotalSupply
			}
			actualMint += blockMint
			p.TotalSupply += blockMint

			if len(p.PairPoolWithCount) != 0{
				for pairAddr, addrMature := range p.MaturePairAccount {
					pairTotalReward := p.PairPoolWithCount[pairAddr]
					if pairTotalReward != 0 && p.TotalSupply + pairTotalReward <= p.MaxSupply{
						for addr, lp := range addrMature {
							totalLp := p.MaturePair[pairAddr]
							if totalLp != 0 {
								reward := big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(int64(lp)), big.NewInt(int64(pairTotalReward))), big.NewInt(int64(totalLp))).Uint64()
								p.TotalSupply += reward
								pairReward, exist := p.AccountReward[addr]
								if exist {
									pairReward[pairAddr] = pairReward[pairAddr] + reward
								} else {
									p.AccountReward[addr] = map[hasharry.Address]uint64{
										pairAddr: reward,
									}
								}
							}
						}
					}
				}
			}
		}
	} else{
		actualMint = blocks * blockMint
		p.TotalSupply += allBlockMint

		if len(p.PairPoolWithCount) != 0{
			for pairAddr, addrMature := range p.MaturePairAccount {
				pairTotalReward := p.PairPoolWithCount[pairAddr]
				if pairTotalReward != 0{
					for addr, lp := range addrMature {
						totalLp := p.MaturePair[pairAddr]
						if totalLp != 0 {
							reward := big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(int64(lp)), big.NewInt(int64(pairTotalReward))), big.NewInt(int64(totalLp))).Uint64()
							reward = reward * blocks
							p.TotalSupply += reward
							pairReward, exist := p.AccountReward[addr]
							if exist {
								pairReward[pairAddr] = pairReward[pairAddr] + reward
							} else {
								p.AccountReward[addr] = map[hasharry.Address]uint64{
									pairAddr: reward,
								}
							}
						}
					}
				}
			}
		}
	}


	p.LastHeight = current
	return actualMint
}

func (p *Pledge) GetPledgeReward(address, pair hasharry.Address) uint64 {
	pairReward, exist := p.AccountReward[address]
	if !exist {
		return 0
	}
	return pairReward[pair]
}

func (p *Pledge) GetPledgeRewards(address hasharry.Address) map[hasharry.Address]uint64 {
	return p.AccountReward[address]
}

func (p *Pledge) GetMaturePledge(address, pair hasharry.Address) uint64 {
	addrPledge, exist := p.MaturePairAccount[pair]
	if !exist {
		return 0
	}
	return addrPledge[address]
}

func (p *Pledge) GetUnripePledge(address, pair hasharry.Address) uint64 {
	for _, unripePledge := range p.UnripePledge {
		addrPledge, exist := unripePledge[pair]
		if !exist {
			return 0
		}
		return addrPledge[address]
	}
	return 0
}

func (p *Pledge) GetDeletedPledge(address, pair hasharry.Address) uint64 {
	addrPledge, exist := p.DeletedPairAccount[pair]
	if !exist {
		return 0
	}
	return addrPledge[address]
}

func (p *Pledge) GetPledgeYields() map[hasharry.Address]float64{
	yields := map[hasharry.Address]float64{}
	for pair, blockReward := range p.PairPoolWithCount{
		pledge := p.PledgePair[pair]
		yields[pair] = float64(blockReward) / float64(pledge)
	}
	return yields
}


func (p *Pledge) ToRlpV2() *RlpPledgeV2 {
	rlpPledge := &RlpPledgeV2{
		PreMint:           p.PreMint,
		DayMintAmount:     p.BlockMintAmount,
		Receiver:          p.Receiver,
		TotalSupply:       p.TotalSupply,
		MaxSupply:         p.MaxSupply,
		PledgeMatureTime:  p.PledgeMatureTime,
		DayRewardAmount:   p.DayRewardAmount,
		Start:             p.Start,
		RewardToken:        p.RewardToken,
		RewardSymbol:       p.RewardSymbol,
		Admin:              p.Admin,
		LastHeight:         p.LastHeight,
		PairPool:           nil,
		PairPoolWithCount:  nil,
		PledgePair:         nil,
		UnripePledge:       nil,
		MaturePair:         nil,
		MaturePairAccount:  nil,
		DeletedPairAccount: nil,
		PoolReward:         nil,
		AccountReward:      nil,
	}
	rlpPledge.PairPool = make([]hasharry.Address, 0)
	for pair, _ := range p.PairPool {
		rlpPledge.PairPool = append(rlpPledge.PairPool, pair)
	}
	sort.Slice(rlpPledge.PairPool, func(i, j int) bool {
		return strings.Compare(rlpPledge.PairPool[i].String(), rlpPledge.PairPool[j].String()) < 0
	})

	rlpPledge.PairPoolWithCount = make([]PairPool, 0)
	for pair, reward := range p.PairPoolWithCount {
		rlpPledge.PairPoolWithCount = append(rlpPledge.PairPoolWithCount, PairPool{
			Pair:   pair,
			Reward: reward,
		})
	}
	sort.Slice(rlpPledge.PairPoolWithCount, func(i, j int) bool {
		return strings.Compare(rlpPledge.PairPoolWithCount[i].Pair.String(), rlpPledge.PairPoolWithCount[j].Pair.String()) < 0
	})

	rlpPledge.PledgePair = make([]PledgePair, 0)
	for pair, total := range p.PledgePair {
		rlpPledge.PledgePair = append(rlpPledge.PledgePair, PledgePair{
			Pair:   pair,
			Amount: total,
		})
	}
	sort.Slice(rlpPledge.PledgePair, func(i, j int) bool {
		return strings.Compare(rlpPledge.PledgePair[i].Pair.String(), rlpPledge.PledgePair[j].Pair.String()) < 0
	})

	rlpPledge.UnripePledge = make([]UnripePledge, 0)
	for day, pairPledge := range p.UnripePledge {
		for pair, addrPledge := range pairPledge {
			for addr, amount := range addrPledge {
				rlpPledge.UnripePledge = append(rlpPledge.UnripePledge, UnripePledge{
					Day:     day,
					Pair:    pair,
					Address: addr,
					Amount:  amount,
				})
			}
		}

	}
	sort.Slice(rlpPledge.UnripePledge, func(i, j int) bool {
		pairCp := strings.Compare(rlpPledge.UnripePledge[i].Pair.String(), rlpPledge.UnripePledge[j].Pair.String())
		if pairCp == 0 {
			return strings.Compare(rlpPledge.UnripePledge[i].Address.String(), rlpPledge.UnripePledge[j].Address.String()) < 0
		} else {
			return pairCp < 0
		}
	})

	rlpPledge.MaturePair = make([]PledgePair, 0)
	for pair, total := range p.MaturePair {
		rlpPledge.MaturePair = append(rlpPledge.MaturePair, PledgePair{
			Pair:   pair,
			Amount: total,
		})

	}
	sort.Slice(rlpPledge.MaturePair, func(i, j int) bool {
		return strings.Compare(rlpPledge.MaturePair[i].Pair.String(), rlpPledge.MaturePair[j].Pair.String()) < 0
	})

	rlpPledge.MaturePairAccount = make([]AccountPledge, 0)
	for pair, addrValue := range p.MaturePairAccount {
		for addr, value := range addrValue {
			rlpPledge.MaturePairAccount = append(rlpPledge.MaturePairAccount, AccountPledge{
				Address: addr,
				Pair:    pair,
				Amount:  value,
			})

		}
	}
	sort.Slice(rlpPledge.MaturePairAccount, func(i, j int) bool {
		pairCp := strings.Compare(rlpPledge.MaturePairAccount[i].Pair.String(), rlpPledge.MaturePairAccount[j].Pair.String())
		if pairCp == 0 {
			return strings.Compare(rlpPledge.MaturePairAccount[i].Address.String(), rlpPledge.MaturePairAccount[j].Address.String()) < 0
		} else {
			return pairCp < 0
		}
	})

	rlpPledge.DeletedPairAccount = make([]AccountPledge, 0)
	for pair, addrValue := range p.DeletedPairAccount {
		for addr, value := range addrValue {
			rlpPledge.DeletedPairAccount = append(rlpPledge.DeletedPairAccount, AccountPledge{
				Address: addr,
				Pair:    pair,
				Amount:  value,
			})

		}
	}
	sort.Slice(rlpPledge.DeletedPairAccount, func(i, j int) bool {
		pairCp := strings.Compare(rlpPledge.DeletedPairAccount[i].Pair.String(), rlpPledge.DeletedPairAccount[j].Pair.String())
		if pairCp == 0 {
			return strings.Compare(rlpPledge.DeletedPairAccount[i].Address.String(), rlpPledge.DeletedPairAccount[j].Address.String()) < 0
		} else {
			return pairCp < 0
		}
	})


	rlpPledge.PoolReward = make([]PoolReward, 0)
	for day, pairReward := range p.PoolReward {
		for pair, reward := range pairReward {
			rlpPledge.PoolReward = append(rlpPledge.PoolReward, PoolReward{
				Day:    day,
				Pair:   pair,
				Reward: reward,
			})

		}
	}
	sort.Slice(rlpPledge.PoolReward, func(i, j int) bool {
		if rlpPledge.PoolReward[i].Day == rlpPledge.PoolReward[j].Day {
			return strings.Compare(rlpPledge.PoolReward[i].Pair.String(), rlpPledge.PoolReward[j].Pair.String()) < 0
		} else {
			return rlpPledge.PoolReward[i].Day < rlpPledge.PoolReward[j].Day
		}
	})

	rlpPledge.AccountReward = make([]AccountReward, 0)
	for addr, pairReward := range p.AccountReward {
		for pair, reward := range pairReward {
			rlpPledge.AccountReward = append(rlpPledge.AccountReward, AccountReward{
				Address: addr,
				Pair:    pair,
				Reward:  reward,
			})

		}
	}
	sort.Slice(rlpPledge.AccountReward, func(i, j int) bool {
		addrCp := strings.Compare(rlpPledge.AccountReward[i].Address.String(), rlpPledge.AccountReward[j].Address.String())
		if addrCp == 0 {
			return strings.Compare(rlpPledge.AccountReward[i].Pair.String(), rlpPledge.AccountReward[j].Pair.String()) < 0
		} else {
			return addrCp < 0
		}
	})

	return rlpPledge
}

func (p *Pledge) ToRlp() *RlpPledgeV2 {
	rlpPledge := &RlpPledgeV2{
		PreMint:           p.PreMint,
		DayMintAmount:     p.BlockMintAmount,
		Receiver:          p.Receiver,
		TotalSupply:       p.TotalSupply,
		MaxSupply:         p.MaxSupply,
		PledgeMatureTime:  p.PledgeMatureTime,
		DayRewardAmount:   p.DayRewardAmount,
		Start:             p.Start,
		RewardToken:        p.RewardToken,
		RewardSymbol:       p.RewardSymbol,
		Admin:              p.Admin,
		LastHeight:         p.LastHeight,
		PairPool:           nil,
		PledgePair:         nil,
		UnripePledge:       nil,
		MaturePair:         nil,
		MaturePairAccount:  nil,
		DeletedPairAccount: nil,
		PoolReward:         nil,
		AccountReward:      nil,
	}
	rlpPledge.PairPool = make([]hasharry.Address, 0)
	for pair, _ := range p.PairPool {
		rlpPledge.PairPool = append(rlpPledge.PairPool, pair)
	}
	sort.Slice(rlpPledge.PairPool, func(i, j int) bool {
		return strings.Compare(rlpPledge.PairPool[i].String(), rlpPledge.PairPool[j].String()) < 0
	})

	rlpPledge.PledgePair = make([]PledgePair, 0)
	for pair, total := range p.PledgePair {
		rlpPledge.PledgePair = append(rlpPledge.PledgePair, PledgePair{
			Pair:   pair,
			Amount: total,
		})
	}
	sort.Slice(rlpPledge.PledgePair, func(i, j int) bool {
		return strings.Compare(rlpPledge.PledgePair[i].Pair.String(), rlpPledge.PledgePair[j].Pair.String()) < 0
	})

	rlpPledge.UnripePledge = make([]UnripePledge, 0)
	for day, pairPledge := range p.UnripePledge {
		for pair, addrPledge := range pairPledge {
			for addr, amount := range addrPledge {
				rlpPledge.UnripePledge = append(rlpPledge.UnripePledge, UnripePledge{
					Day:     day,
					Pair:    pair,
					Address: addr,
					Amount:  amount,
				})
			}
		}

	}
	sort.Slice(rlpPledge.UnripePledge, func(i, j int) bool {
		pairCp := strings.Compare(rlpPledge.UnripePledge[i].Pair.String(), rlpPledge.UnripePledge[j].Pair.String())
		if pairCp == 0 {
			return strings.Compare(rlpPledge.UnripePledge[i].Address.String(), rlpPledge.UnripePledge[j].Address.String()) < 0
		} else {
			return pairCp < 0
		}
	})

	rlpPledge.MaturePair = make([]PledgePair, 0)
	for pair, total := range p.MaturePair {
		rlpPledge.MaturePair = append(rlpPledge.MaturePair, PledgePair{
			Pair:   pair,
			Amount: total,
		})

	}
	sort.Slice(rlpPledge.MaturePair, func(i, j int) bool {
		return strings.Compare(rlpPledge.MaturePair[i].Pair.String(), rlpPledge.MaturePair[j].Pair.String()) < 0
	})

	rlpPledge.MaturePairAccount = make([]AccountPledge, 0)
	for pair, addrValue := range p.MaturePairAccount {
		for addr, value := range addrValue {
			rlpPledge.MaturePairAccount = append(rlpPledge.MaturePairAccount, AccountPledge{
				Address: addr,
				Pair:    pair,
				Amount:  value,
			})

		}
	}
	sort.Slice(rlpPledge.MaturePairAccount, func(i, j int) bool {
		pairCp := strings.Compare(rlpPledge.MaturePairAccount[i].Pair.String(), rlpPledge.MaturePairAccount[j].Pair.String())
		if pairCp == 0 {
			return strings.Compare(rlpPledge.MaturePairAccount[i].Address.String(), rlpPledge.MaturePairAccount[j].Address.String()) < 0
		} else {
			return pairCp < 0
		}
	})

	rlpPledge.DeletedPairAccount = make([]AccountPledge, 0)
	for pair, addrValue := range p.DeletedPairAccount {
		for addr, value := range addrValue {
			rlpPledge.DeletedPairAccount = append(rlpPledge.DeletedPairAccount, AccountPledge{
				Address: addr,
				Pair:    pair,
				Amount:  value,
			})

		}
	}
	sort.Slice(rlpPledge.DeletedPairAccount, func(i, j int) bool {
		pairCp := strings.Compare(rlpPledge.DeletedPairAccount[i].Pair.String(), rlpPledge.DeletedPairAccount[j].Pair.String())
		if pairCp == 0 {
			return strings.Compare(rlpPledge.DeletedPairAccount[i].Address.String(), rlpPledge.DeletedPairAccount[j].Address.String()) < 0
		} else {
			return pairCp < 0
		}
	})


	rlpPledge.PoolReward = make([]PoolReward, 0)
	for day, pairReward := range p.PoolReward {
		for pair, reward := range pairReward {
			rlpPledge.PoolReward = append(rlpPledge.PoolReward, PoolReward{
				Day:    day,
				Pair:   pair,
				Reward: reward,
			})

		}
	}
	sort.Slice(rlpPledge.PoolReward, func(i, j int) bool {
		if rlpPledge.PoolReward[i].Day == rlpPledge.PoolReward[j].Day {
			return strings.Compare(rlpPledge.PoolReward[i].Pair.String(), rlpPledge.PoolReward[j].Pair.String()) < 0
		} else {
			return rlpPledge.PoolReward[i].Day < rlpPledge.PoolReward[j].Day
		}
	})

	rlpPledge.AccountReward = make([]AccountReward, 0)
	for addr, pairReward := range p.AccountReward {
		for pair, reward := range pairReward {
			rlpPledge.AccountReward = append(rlpPledge.AccountReward, AccountReward{
				Address: addr,
				Pair:    pair,
				Reward:  reward,
			})

		}
	}
	sort.Slice(rlpPledge.AccountReward, func(i, j int) bool {
		addrCp := strings.Compare(rlpPledge.AccountReward[i].Address.String(), rlpPledge.AccountReward[j].Address.String())
		if addrCp == 0 {
			return strings.Compare(rlpPledge.AccountReward[i].Pair.String(), rlpPledge.AccountReward[j].Pair.String()) < 0
		} else {
			return addrCp < 0
		}
	})

	return rlpPledge
}

func (p *Pledge) Bytes() []byte {
	bytes, _ := rlp.EncodeToBytes(p.ToRlpV2())
	return bytes
}

type PledgePair struct {
	Pair   hasharry.Address
	Amount uint64
}

type PairPool struct {
	Pair   hasharry.Address
	Reward uint64
}

type UnripePledge struct {
	Day     uint64
	Pair    hasharry.Address
	Address hasharry.Address
	Amount  uint64
}

type AccountPledge struct {
	Address hasharry.Address
	Pair    hasharry.Address
	Amount  uint64
}

type PoolReward struct {
	Day    uint64
	Pair   hasharry.Address
	Reward uint64
}

type AccountReward struct {
	Address hasharry.Address
	Pair    hasharry.Address
	Reward  uint64
}

type RlpPledge struct {
	// 预铸币
	PreMint uint64
	// 一次性铸币量
	DayMintAmount uint64
	// 铸币接收地址
	Receiver hasharry.Address
	// 总铸币量
	TotalSupply uint64
	// 最大供应量
	MaxSupply uint64
	// 质押成熟时间
	PledgeMatureTime uint64
	// 每天质押奖励
	DayRewardAmount uint64
	// 开始高度
	Start uint64
	// 奖励token
	RewardToken hasharry.Address
	// 奖励token symbol
	RewardSymbol string
	// 管理员
	Admin hasharry.Address
	// 最后高度
	LastHeight uint64
	// pair 池
	PairPool []hasharry.Address
	// 每个pair总质押量
	PledgePair []PledgePair
	// 未成熟的某天某个pair某个地址的质押量
	UnripePledge []UnripePledge
	// 已经成熟的pair总质押量
	MaturePair []PledgePair
	// 已经成熟的pair某个地址质押量
	MaturePairAccount []AccountPledge
	// 被删除的pair的某个地址的质押量，可取回
	DeletedPairAccount []AccountPledge

	// 每天每个pair池的币量
	PoolReward []PoolReward
	// 每个账户每个pair池子奖励的币量
	AccountReward []AccountReward
	MintAmount    uint64
}


type RlpPledgeV2 struct {
	// 预铸币
	PreMint uint64
	// 一次性铸币量
	DayMintAmount uint64
	// 铸币接收地址
	Receiver hasharry.Address
	// 总铸币量
	TotalSupply uint64
	// 最大供应量
	MaxSupply uint64
	// 质押成熟时间
	PledgeMatureTime uint64
	// 每天质押奖励
	DayRewardAmount uint64
	// 开始高度
	Start uint64
	// 奖励token
	RewardToken hasharry.Address
	// 奖励token symbol
	RewardSymbol string
	// 管理员
	Admin hasharry.Address
	// 最后高度
	LastHeight uint64
	// pair 池
	PairPool []hasharry.Address
	// pair 池的奖励数量
	PairPoolWithCount []PairPool
	// 每个pair总质押量
	PledgePair []PledgePair
	// 未成熟的某天某个pair某个地址的质押量
	UnripePledge []UnripePledge
	// 已经成熟的pair总质押量
	MaturePair []PledgePair
	// 已经成熟的pair某个地址质押量
	MaturePairAccount []AccountPledge
	// 被删除的pair的某个地址的质押量，可取回
	DeletedPairAccount []AccountPledge

	// 每天每个pair池的币量
	PoolReward []PoolReward
	// 每个账户每个pair池子奖励的币量
	AccountReward []AccountReward
	MintAmount    uint64
}

func DecodeToPledge(bytes []byte) (*Pledge, error) {
	var rlpPd *RlpPledgeV2
	err := rlp.DecodeBytes(bytes, &rlpPd)
	if err != nil{
		return DecodeV1ToPledge(bytes)
	}
	pd := &Pledge{
		PreMint:           rlpPd.PreMint,
		BlockMintAmount:   rlpPd.DayMintAmount,
		Receiver:          rlpPd.Receiver,
		TotalSupply:       rlpPd.TotalSupply,
		MaxSupply:         rlpPd.MaxSupply,
		PledgeMatureTime:  rlpPd.PledgeMatureTime,
		DayRewardAmount:   rlpPd.DayRewardAmount,
		Start:             rlpPd.Start,
		RewardToken:       rlpPd.RewardToken,
		RewardSymbol:      rlpPd.RewardSymbol,
		Admin:             rlpPd.Admin,
		LastHeight:        rlpPd.LastHeight,
		PairPool:          map[hasharry.Address]bool{},
		PairPoolWithCount: map[hasharry.Address]uint64{},
		PledgePair:        map[hasharry.Address]uint64{},
		UnripePledge:      map[uint64]map[hasharry.Address]map[hasharry.Address]uint64{},
		MaturePair:        map[hasharry.Address]uint64{},
		MaturePairAccount: map[hasharry.Address]map[hasharry.Address]uint64{},
		DeletedPairAccount:map[hasharry.Address]map[hasharry.Address]uint64{},
		PoolReward:        map[uint64]map[hasharry.Address]uint64{},
		AccountReward:     map[hasharry.Address]map[hasharry.Address]uint64{},
	}

	for _, pair := range rlpPd.PairPool {
		pd.PairPool[pair] = true
	}

	for _, pairPool := range rlpPd.PairPoolWithCount {
		pd.PairPoolWithCount[pairPool.Pair] = pairPool.Reward
	}

	for _, pledgePair := range rlpPd.PledgePair {
		pd.PledgePair[pledgePair.Pair] = pledgePair.Amount
	}


	for _, unripePledge := range rlpPd.UnripePledge {
		pairPledge, exist := pd.UnripePledge[unripePledge.Day]
		if exist {
			addrPledge, exist := pairPledge[unripePledge.Pair]
			if exist {
				addrPledge[unripePledge.Address] = unripePledge.Amount
			} else {
				pairPledge[unripePledge.Pair] = map[hasharry.Address]uint64{
					unripePledge.Address: unripePledge.Amount,
				}
			}
		} else {
			pd.UnripePledge[unripePledge.Day] = map[hasharry.Address]map[hasharry.Address]uint64{
				unripePledge.Pair: {
					unripePledge.Address: unripePledge.Amount,
				},
			}
		}
	}

	for _, maturePair := range rlpPd.MaturePair {
		pd.MaturePair[maturePair.Pair] = maturePair.Amount
	}

	for _, maturePledge := range rlpPd.MaturePairAccount {
		addrPledge, exist := pd.MaturePairAccount[maturePledge.Pair]
		if exist {
			addrPledge[maturePledge.Address] = maturePledge.Amount
		} else {
			pd.MaturePairAccount[maturePledge.Pair] = map[hasharry.Address]uint64{
				maturePledge.Address: maturePledge.Amount,
			}
		}
	}

	for _, deletedPledge := range rlpPd.DeletedPairAccount {
		addrPledge, exist := pd.DeletedPairAccount[deletedPledge.Pair]
		if exist {
			addrPledge[deletedPledge.Address] = deletedPledge.Amount
		} else {
			pd.DeletedPairAccount[deletedPledge.Pair] = map[hasharry.Address]uint64{
				deletedPledge.Address: deletedPledge.Amount,
			}
		}
	}

	for _, poolReward := range rlpPd.PoolReward {
		pairReward, exist := pd.PoolReward[poolReward.Day]
		if exist {
			pairReward[poolReward.Pair] = poolReward.Reward
		} else {
			pd.PoolReward[poolReward.Day] = map[hasharry.Address]uint64{
				poolReward.Pair: poolReward.Reward,
			}
		}
	}

	for _, accountReward := range rlpPd.AccountReward {
		pairReward, exist := pd.AccountReward[accountReward.Address]
		if exist {
			pairReward[accountReward.Pair] = accountReward.Reward
		} else {
			pd.AccountReward[accountReward.Address] = map[hasharry.Address]uint64{
				accountReward.Pair: accountReward.Reward,
			}
		}
	}

	return pd, err

}


func DecodeV1ToPledge(bytes []byte) (*Pledge, error) {
	var rlpPd *RlpPledge
	err := rlp.DecodeBytes(bytes, &rlpPd)
	if err != nil{
		return nil, err
	}
	pd := &Pledge{
		PreMint:          rlpPd.PreMint,
		BlockMintAmount:  rlpPd.DayMintAmount,
		Receiver:         rlpPd.Receiver,
		TotalSupply:      rlpPd.TotalSupply,
		MaxSupply:        rlpPd.MaxSupply,
		PledgeMatureTime: rlpPd.PledgeMatureTime,
		DayRewardAmount:  rlpPd.DayRewardAmount,
		Start:            rlpPd.Start,
		RewardToken:      rlpPd.RewardToken,
		RewardSymbol:     rlpPd.RewardSymbol,
		Admin:            rlpPd.Admin,
		LastHeight:       rlpPd.LastHeight,
		PairPool:          map[hasharry.Address]bool{},
		PairPoolWithCount: map[hasharry.Address]uint64{},
		PledgePair:        map[hasharry.Address]uint64{},
		UnripePledge:      map[uint64]map[hasharry.Address]map[hasharry.Address]uint64{},
		MaturePair:        map[hasharry.Address]uint64{},
		MaturePairAccount: map[hasharry.Address]map[hasharry.Address]uint64{},
		DeletedPairAccount:map[hasharry.Address]map[hasharry.Address]uint64{},
		PoolReward:        map[uint64]map[hasharry.Address]uint64{},
		AccountReward:     map[hasharry.Address]map[hasharry.Address]uint64{},
	}

	for _, pair := range rlpPd.PairPool {
		pd.PairPool[pair] = true
	}

	for _, pledgePair := range rlpPd.PledgePair {
		pd.PledgePair[pledgePair.Pair] = pledgePair.Amount
	}

	for _, unripePledge := range rlpPd.UnripePledge {
		pairPledge, exist := pd.UnripePledge[unripePledge.Day]
		if exist {
			addrPledge, exist := pairPledge[unripePledge.Pair]
			if exist {
				addrPledge[unripePledge.Address] = unripePledge.Amount
			} else {
				pairPledge[unripePledge.Pair] = map[hasharry.Address]uint64{
					unripePledge.Address: unripePledge.Amount,
				}
			}
		} else {
			pd.UnripePledge[unripePledge.Day] = map[hasharry.Address]map[hasharry.Address]uint64{
				unripePledge.Pair: {
					unripePledge.Address: unripePledge.Amount,
				},
			}
		}
	}

	for _, maturePair := range rlpPd.MaturePair {
		pd.MaturePair[maturePair.Pair] = maturePair.Amount
	}

	for _, maturePledge := range rlpPd.MaturePairAccount {
		addrPledge, exist := pd.MaturePairAccount[maturePledge.Pair]
		if exist {
			addrPledge[maturePledge.Address] = maturePledge.Amount
		} else {
			pd.MaturePairAccount[maturePledge.Pair] = map[hasharry.Address]uint64{
				maturePledge.Address: maturePledge.Amount,
			}
		}
	}

	for _, deletedPledge := range rlpPd.DeletedPairAccount {
		addrPledge, exist := pd.DeletedPairAccount[deletedPledge.Pair]
		if exist {
			addrPledge[deletedPledge.Address] = deletedPledge.Amount
		} else {
			pd.DeletedPairAccount[deletedPledge.Pair] = map[hasharry.Address]uint64{
				deletedPledge.Address: deletedPledge.Amount,
			}
		}
	}

	for _, poolReward := range rlpPd.PoolReward {
		pairReward, exist := pd.PoolReward[poolReward.Day]
		if exist {
			pairReward[poolReward.Pair] = poolReward.Reward
		} else {
			pd.PoolReward[poolReward.Day] = map[hasharry.Address]uint64{
				poolReward.Pair: poolReward.Reward,
			}
		}
	}

	for _, accountReward := range rlpPd.AccountReward {
		pairReward, exist := pd.AccountReward[accountReward.Address]
		if exist {
			pairReward[accountReward.Pair] = accountReward.Reward
		} else {
			pd.AccountReward[accountReward.Address] = map[hasharry.Address]uint64{
				accountReward.Pair: accountReward.Reward,
			}
		}
	}

	return pd, err
}