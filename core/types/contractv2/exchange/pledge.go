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
	//PairPool map[hasharry.Address]bool
	// 池子产币量
	PairPoolWithCount map[hasharry.Address]uint64
	// 某个地址某个池子的最后一次质押高度
	AccountLastPledge map[hasharry.Address]uint64
	// 每个pair总质押量
	PledgePair map[hasharry.Address]uint64
	// 未成熟的某天某个pair某个地址的质押量
	//UnripePledge map[uint64]map[hasharry.Address]map[hasharry.Address]uint64
	// 已经成熟的pair总质押量
	//MaturePair map[hasharry.Address]uint64
	// 已经成熟的pair某个地址质押量
	MaturePairAccount map[hasharry.Address]map[hasharry.Address]uint64
	// 被删除的pair的某个地址的质押量，可取回
	DeletedPairAccount map[hasharry.Address]map[hasharry.Address]uint64

	// 每天每个pair池的币量(弃用)
	//PoolReward map[uint64]map[hasharry.Address]uint64
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
		TotalSupply:        preMint,
		MaxSupply:          maxSupply,
		DayRewardAmount:    0,
		RewardToken:        exchange,
		RewardSymbol:       symbol,
		PledgeMatureTime:   0,
		Admin:              admin,
		Start:              0,
		LastHeight:         0,
		PairPoolWithCount:  map[hasharry.Address]uint64{},
		AccountLastPledge:  map[hasharry.Address]uint64{},
		AccountReward:      map[hasharry.Address]map[hasharry.Address]uint64{},
		PledgePair:         map[hasharry.Address]uint64{},
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

	addressPledge, exist := p.MaturePairAccount[pair]
	if exist {
		value, exist := addressPledge[address]
		if exist {
			addressPledge[address] = value + amount
		} else {
			addressPledge[address] = amount
		}
	} else {
		p.MaturePairAccount[pair] = map[hasharry.Address]uint64{
			address: amount,
		}
	}

	p.PledgePair[pair] = p.PledgePair[pair] + amount
	p.AccountLastPledge = map[hasharry.Address]uint64{address : height}
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
		}else{
			return errors.New("insufficient pledge")
		}
	}else{
		addrPledge, exist = p.MaturePairAccount[pair]
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
	_, exist := p.PairPoolWithCount[pair]
	if !exist{
		return errors.New("pair does not exist")
	}
	delete(p.PairPoolWithCount, pair)

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

func (p *Pledge) RemoveReward(address hasharry.Address, height uint64) ([]Reward, error) {
	lastHeight := p.AccountLastPledge[address]
	if lastHeight + p.PledgeMatureTime > height{
		return nil, fmt.Errorf("the last pledge has not reached the mature height, temporarily take out")
	}
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
		return rewardList, nil
	} else {
		return rewardList, nil
	}
}

func (p *Pledge) GetPledgeAmount(address, pair hasharry.Address) (matureLp uint64, deleteLp uint64) {
	matureLp = 0
	addrPledge, exist := p.MaturePairAccount[pair]
	if exist {
		matureLp = addrPledge[address]
	}

	deleteLp = 0
	addrPledge, exist = p.DeletedPairAccount[pair]
	if exist {
		deleteLp = addrPledge[address]
	}
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
							totalLp := p.PledgePair[pairAddr]
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
						totalLp := p.PledgePair[pairAddr]
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

func (p *Pledge) GetDeletedPledge(address, pair hasharry.Address) uint64 {
	addrPledge, exist := p.DeletedPairAccount[pair]
	if !exist {
		return 0
	}
	return addrPledge[address]
}

type PoolInfo struct {
	Address string
	YieldRate float64
	TotalPledge uint64
	TotalReward uint64

}

func (p *Pledge) GetPoolInfo(pair hasharry.Address) *PoolInfo {
	blockReward, exist := p.PairPoolWithCount[pair]
	if !exist{
		return &PoolInfo{}
	}
	pledge := p.PledgePair[pair]
	var yields float64
	if pledge != 0{
		yields = float64(blockReward) / float64(pledge)
	}else{
		yields = float64(blockReward) / 1e8
	}
	return &PoolInfo{
		Address:     pair.String(),
		YieldRate:   yields,
		TotalPledge: pledge,
		TotalReward: blockReward,
	}
}

func (p *Pledge) GetPoolInfos() []PoolInfo {
	infos := make([]PoolInfo, 0)
	for pair, blockReward := range p.PairPoolWithCount{
		pledge := p.PledgePair[pair]
		var yields float64
		if pledge != 0{
			yields = float64(blockReward) / float64(pledge)
		}else{
			yields = float64(blockReward) / 1e8
		}
		infos = append(infos, PoolInfo{
			Address:     pair.String(),
			YieldRate:   yields,
			TotalPledge: pledge,
			TotalReward: blockReward,
		})

	}
	return infos
}

func (p *Pledge) GetPledgeYields() map[hasharry.Address]float64{
	yields := map[hasharry.Address]float64{}
	for pair, blockReward := range p.PairPoolWithCount{
		pledge := p.PledgePair[pair]
		if pledge != 0{
			yields[pair] = float64(blockReward) / float64(pledge)
		}else{
			yields[pair] = float64(blockReward) / 1e8
		}
	}
	return yields
}

func (p *Pledge) ToRlp() *RlpPledge {
	rlpPledge := &RlpPledge{
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
	}


	rlpPledge.PairPool = make([]hasharry.Address, 0)

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


	rlpPledge.MaturePair = make([]PledgePair, 0)

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
	}

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

	rlpPledge.AccountLastPledge = make([]PledgeHeight, 0)
	for pair, height := range p.AccountLastPledge {
		rlpPledge.AccountLastPledge = append(rlpPledge.AccountLastPledge, PledgeHeight{
			Pair:   pair,
			Height: height,
		})
	}
	sort.Slice(rlpPledge.AccountLastPledge, func(i, j int) bool {
		return strings.Compare(rlpPledge.AccountLastPledge[i].Pair.String(), rlpPledge.AccountLastPledge[j].Pair.String()) < 0
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

func (p *Pledge)GetSymbol() string{
	return ""
}

func (p *Pledge) Bytes() []byte {
	var bytes []byte
	if p.Start > param.UIPBlock_3{
		bytes, _ = rlp.EncodeToBytes(p.ToRlpV2())
	}else{
		bytes, _ = rlp.EncodeToBytes(p.ToRlp())
	}
	return bytes
}

type PledgePair struct {
	Pair   hasharry.Address
	Amount uint64
}

type PledgeHeight struct {
	Pair   hasharry.Address
	Height uint64
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
	// 弃用
	PairPool []hasharry.Address
	// 弃用
	UnripePledge []UnripePledge
	// 弃用
	MaturePair []PledgePair
	PoolReward []PoolReward

	// pair 池的奖励数量
	PairPoolWithCount []PairPool
	// 每个pair总质押量
	PledgePair []PledgePair
	// 最后一次质押高度
	AccountLastPledge []PledgeHeight
	// 已经成熟的pair某个地址质押量
	MaturePairAccount []AccountPledge
	// 被删除的pair的某个地址的质押量，可取回
	DeletedPairAccount []AccountPledge
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
		PairPoolWithCount: map[hasharry.Address]uint64{},
		PledgePair:        map[hasharry.Address]uint64{},
		AccountLastPledge: map[hasharry.Address]uint64{},
		MaturePairAccount: map[hasharry.Address]map[hasharry.Address]uint64{},
		DeletedPairAccount:map[hasharry.Address]map[hasharry.Address]uint64{},
		AccountReward:     map[hasharry.Address]map[hasharry.Address]uint64{},
	}


	for _, pairPool := range rlpPd.PairPoolWithCount {
		pd.PairPoolWithCount[pairPool.Pair] = pairPool.Reward
	}

	for _, pledgePair := range rlpPd.PledgePair {
		pd.PledgePair[pledgePair.Pair] = pledgePair.Amount
	}

	for _, pledgeHeight := range rlpPd.AccountLastPledge {
		pd.AccountLastPledge[pledgeHeight.Pair] = pledgeHeight.Height
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
		PairPoolWithCount: map[hasharry.Address]uint64{},
		PledgePair:        map[hasharry.Address]uint64{},
		AccountLastPledge: map[hasharry.Address]uint64{},
		MaturePairAccount: map[hasharry.Address]map[hasharry.Address]uint64{},
		DeletedPairAccount:map[hasharry.Address]map[hasharry.Address]uint64{},
		AccountReward:     map[hasharry.Address]map[hasharry.Address]uint64{},
	}

	for _, pledgePair := range rlpPd.PledgePair {
		pd.PledgePair[pledgePair.Pair] = pledgePair.Amount
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
