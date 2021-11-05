package exchange_func

import (
	"errors"
	"fmt"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/param"
	"github.com/UBChainNet/UBChain/ut"
)

const minPledge = 10000

type PledgeInitBody struct {
	Exchange  hasharry.Address
	Admin     hasharry.Address
	MaxSupply uint64
	DayMint   uint64
}

func (p *PledgeInitBody) Verify() error {
	if ok := ut.IsValidContractAddress(param.Net, p.Exchange.String()); !ok {
		return errors.New("wrong exchange address")
	}
	if ok := ut.CheckUBCAddress(param.Net, p.Admin.String()); !ok {
		return errors.New("wrong admin address")
	}
	if p.MaxSupply > param.MaxContractCoin{
		return fmt.Errorf("max supply cannot exceed %d", param.MaxContractCoin)
	}
	if p.DayMint > p.MaxSupply {
		return fmt.Errorf("no more than totalsupply coins can be minted per day")
	}
	return nil
}

type PledgeAddBody struct {
	Pair hasharry.Address
	Amount uint64
}

func (p *PledgeAddBody) Verify() error {
	if ok := ut.IsContractV2Address(param.Net, p.Pair.String()); !ok {
		return errors.New("wrong Pair address")
	}
	if p.Amount < minPledge{
		return fmt.Errorf("the amount pledged shall not be less than %d", minPledge)
	}
	return nil
}

type PledgeRemoveBody struct {
	Pair hasharry.Address
	Amount uint64
}

func (p *PledgeRemoveBody) Verify() error {
	if ok := ut.IsContractV2Address(param.Net, p.Pair.String()); !ok {
		return errors.New("wrong Pair address")
	}
	if p.Amount < minPledge{
		return fmt.Errorf("the amount pledged shall not be less than %d", minPledge)
	}
	return nil
}


type PledgeRewardRemoveBody struct {
}

func (p *PledgeRewardRemoveBody) Verify() error {
	return nil
}

type PledgeUpdateBody struct {
}

func (p *PledgeUpdateBody) Verify() error {
	return nil
}