package token

import "errors"

type Token struct {
	Admin string
	Symbol string
	MaxSupply uint64
	PreSupply uint64
	Start uint64
	Final uint64
	LastMint uint64
	LastHeight uint64
	Interval uint64
	EachAmount uint64
}

func NewToken(admin string, symbol string, maxSupply, preSupply, start, interval, eachAmount uint64)(*Token, error){
	if interval == 0 || maxSupply == 0{
		return nil, errors.New("the minting interval and maximum supply cannot be 0")
	}
	if maxSupply < preSupply{
		return nil, errors.New("the maximum supply cannot be less than the pre-supply")
	}
	mintCount := (maxSupply - preSupply) / eachAmount
	if (maxSupply - preSupply) % eachAmount != 0{
		mintCount++
	}
	final := mintCount * interval + start
	return &Token{
		Admin:       admin,
		Symbol:      symbol,
		MaxSupply:   maxSupply,
		PreSupply:   preSupply,
		Start:       start,
		LastMint:    start,
		LastHeight:  start,
		Final: 		 final,
		Interval:    interval,
		EachAmount:  eachAmount,
	}, nil
}

func (t *Token)Mint(from string, height uint64)(uint64, error){
	if t.Admin != from{
		return 0, errors.New("forbidden")
	}
	if height <= t.LastHeight{
		return 0, errors.New("it's not mint time")
	}
	if t.LastHeight >= t.Final{
		return 0, errors.New("over maximum supply")
	}
	if t.LastHeight >= t.LastMint && height <  t.LastMint + t.Interval{
		return 0, errors.New("it's not mint time")
	}
	t.LastHeight = height
	if height >= t.Final{
		height = t.Final
	}
	var mintAmount uint64
	mintCount := (height - t.LastMint) / t.Interval
	if height >= t.Final && mintCount > 0 && (t.MaxSupply - t.PreSupply) % t.EachAmount != 0{
		mintAmount = (mintCount - 1) * t.EachAmount + (t.MaxSupply - t.PreSupply) % t.EachAmount
	}else{
		mintAmount = mintCount * t.EachAmount
	}
	t.LastMint = t.LastMint +  mintCount * t.Interval
	return mintAmount, nil
}

func (t *Token)TotalSupply(height uint64)uint64{
	if height >= t.Final{
		return t.MaxSupply
	}
	return height - t.Start / t.Interval * t.EachAmount
}

func (t *Token)