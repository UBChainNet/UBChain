package exchange_runner

import (
	"errors"
	"fmt"
	"github.com/UBChainNet/UBChain/common/codec"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/core/runner/library"
	"github.com/UBChainNet/UBChain/core/types"
	"github.com/UBChainNet/UBChain/core/types/contractv2"
	exchange2 "github.com/UBChainNet/UBChain/core/types/contractv2/exchange"
	"github.com/UBChainNet/UBChain/core/types/functionbody/exchange_func"
	"github.com/UBChainNet/UBChain/crypto/base58"
	"github.com/UBChainNet/UBChain/param"
	"github.com/UBChainNet/UBChain/ut"
	"math/big"
)

type PledgeState struct {
	library    *library.RunnerLibrary
	header   *contractv2.ContractV2
	body     *exchange2.Pledge
}

func NewPledgeState(runnerLibrary *library.RunnerLibrary, pdAddress string) (*PledgeState, error) {
	pdHeader := runnerLibrary.GetContractV2(pdAddress)
	if pdHeader == nil {
		return nil, fmt.Errorf("pledge %s does not exist", pdAddress)
	}
	pdBody, _ := pdHeader.Body.(*exchange2.Pledge)
	return &PledgeState{
		header:  pdHeader,
		body:    pdBody,
		library: runnerLibrary,
	}, nil
}

func (ps *PledgeState)Methods() map[string]*MethodInfo{
	return pledgeMethods
}

func (ps *PledgeState) MethodExist(method string) bool{
	_, exist := pledgeMethods[method]
	return exist
}


type PledgeReward struct {
	Token string  `json:"token"`
	Symbol string `json:"symbol"`
	Amount float64 `json:"amount"`
	Pair 	string `json:"pair"`
}

func (ps *PledgeState)GetPledgeReward(address string, pair string) *PledgeReward{
	reward := ps.body.GetPledgeReward(hasharry.StringToAddress(address), hasharry.StringToAddress(pair))
	return &PledgeReward{
		Token:  ps.body.RewardToken.String(),
		Symbol: ps.body.RewardSymbol,
		Amount: types.Amount(reward).ToCoin(),
		Pair:   pair,
	}
}

func (ps *PledgeState)GetPledgeRewards(address string) []*PledgeReward{
	var rewards []*PledgeReward
	rewardMap := ps.body.GetPledgeRewards(hasharry.StringToAddress(address))
	if rewardMap != nil{
		for pair, reward := range rewardMap{
			rewards = append(rewards, &PledgeReward{
				Token:  ps.body.RewardToken.String(),
				Symbol: ps.body.RewardSymbol,
				Amount: types.Amount(reward).ToCoin(),
				Pair:   pair.String(),
			})
		}
	}
	return rewards
}

type PledgeValue struct {
	MaturePledge float64 `json:"maturePledge"`
	UnripePledge float64 `json:"unripePledge"`
	Pair string `json:"pair"`
}

func (ps *PledgeState)GetPledge(address string, pair string) *PledgeValue{
	maturePledge := ps.body.GetMaturePledge(hasharry.StringToAddress(address), hasharry.StringToAddress(pair))
	unripePledge := ps.body.GetUnripePledge(hasharry.StringToAddress(address), hasharry.StringToAddress(pair))
	return &PledgeValue{
		MaturePledge:  types.Amount(maturePledge).ToCoin(),
		UnripePledge:  types.Amount(unripePledge).ToCoin(),
		Pair:   pair,
	}
}

func (ps *PledgeState)GetPledges(address string) []*PledgeValue{
	var pledges []*PledgeValue
	for pair, _ := range ps.body.PledgePair {
		pledge := ps.GetPledge(address, pair.String())
		if pledge.UnripePledge != 0 || pledge.MaturePledge != 0{
			pledges = append(pledges, pledge)
		}
	}
	return pledges
}

type Pool struct {
	Pair string 	`json:"pair"`
	Weight float64 `json:"weight"`
}

func (ps *PledgeState)GetPairPool() []*Pool{
	var poolList  = make([]*Pool, 0)
	pairValue := map[hasharry.Address]uint64{}
	var totalValue uint64
	ex := ps.library.GetContractV2(ps.body.RewardToken.String())
	for pairAddr, total := range ps.body.MaturePair{
		value, _ := ps.getPairValue(pairAddr, total, ex)
		totalValue += value
		pairValue[pairAddr] = value
	}
	for pairAddr, _ := range ps.body.PairPool{
		value := pairValue[pairAddr]
		if totalValue != 0{
			poolList = append(poolList, &Pool{
				Pair:   pairAddr.String(),
				Weight: float64(value) / float64(totalValue),
			})
		}else{
			poolList = append(poolList, &Pool{
				Pair:   pairAddr.String(),
				Weight: 0,
			})
		}

	}
	return poolList
}


func (ps *PledgeState)getPairValue(pairAddr hasharry.Address, totalLp uint64, exchange *contractv2.ContractV2)(uint64, error){
	var token0value uint64
	var token1value uint64
	pairState, _ := NewPairStateWithExchange(ps.library, pairAddr.String(), exchange)
	pairTotal, err := pairState.TotalValue(types.Amount(totalLp).ToCoin())
	if err != nil{
		return 0, err
	}
	pairToken0, _ := types.NewAmount(pairTotal.Value0)
	pairToken1, _ := types.NewAmount(pairTotal.Value1)
	if pairTotal.Token0 == param.Token.String() || pairTotal.Token1 == param.Token.String() {
		token0value = pairToken0
		token1value = pairToken1
	}else{
		token0WithMain, _ := PairAddress(param.Net, pairState.pairBody.Token0, param.Token, ps.body.RewardToken)
		token0WithMainBody, _ := ps.library.GetPair(hasharry.StringToAddress(token0WithMain))
		if pairState.pairBody.Token0.IsEqual(token0WithMainBody.Token0){
			token0value = big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(int64(pairToken0)), big.NewInt(int64(token0WithMainBody.Reserve1))),
				big.NewInt(int64(token0WithMainBody.Reserve0))).Uint64()
		}else{
			token0value = big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(int64(pairToken0)), big.NewInt(int64(token0WithMainBody.Reserve0))),
				big.NewInt(int64(token0WithMainBody.Reserve1))).Uint64()
		}

		token1WithMain, _ := PairAddress(param.Net, pairState.pairBody.Token1, param.Token, ps.body.RewardToken)
		token1WithMainBody, _ := ps.library.GetPair(hasharry.StringToAddress(token1WithMain))
		if pairState.pairBody.Token1.IsEqual(token1WithMainBody.Token0){
			token1value = big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(int64(pairToken1)), big.NewInt(int64(token1WithMainBody.Reserve1))),
				big.NewInt(int64(token1WithMainBody.Reserve0))).Uint64()
		}else{
			token1value = big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(int64(pairToken1)), big.NewInt(int64(token1WithMainBody.Reserve0))),
				big.NewInt(int64(token1WithMainBody.Reserve1))).Uint64()
		}
	}

	return token0value + token1value, nil
}

type PledgeRunner struct {
	pdState      *PledgeState
	address      hasharry.Address
	tx           types.ITransaction
	contractBody *types.TxContractV2Body
	events       []*types.Event
	height       uint64
}

func NewPledgeRunner(lib *library.RunnerLibrary, tx types.ITransaction, height uint64) (*PledgeRunner, error) {
	var pd *exchange2.Pledge
	var ok bool
	address := tx.GetTxBody().GetContract()
	exHeader := lib.GetContractV2(address.String())
	if exHeader != nil {
		pd, ok  = exHeader.Body.(*exchange2.Pledge)
		if !ok{
			return nil, errors.New("invalid contract type")
		}
	}

	contractBody := tx.GetTxBody().(*types.TxContractV2Body)
	return &PledgeRunner{
		pdState: &PledgeState{
			header:  exHeader,
			body:    pd,
			library: lib,
		},
		address:      address,
		tx:           tx,
		contractBody: contractBody,
		events:       make([]*types.Event, 0),
		height:       height,
	}, nil
}

func (p *PledgeRunner) PreInitVerify() error{
	pd := p.pdState.library.GetContractV2(p.address.String())
	if pd != nil{
		return fmt.Errorf("pledge %s already exist", p.address.String())
	}
	initBody := p.contractBody.Function.(*exchange_func.PledgeInitBody)
	ex, err := p.pdState.library.GetExchange(initBody.Exchange)
	if err != nil{
		return err
	}
	if !ex.Admin.IsEqual(p.tx.From()){
		return fmt.Errorf("forbidden")
	}
	return err
}

func (p *PledgeRunner) Init() {
	var ERR error
	state := &types.ContractV2State{State: types.Contract_Success}
	defer func() {
		if ERR != nil {
			state.State = types.Contract_Failed
			state.Error = ERR.Error()
		} else {
			state.Event = p.events
		}
		p.pdState.library.SetContractV2State(p.tx.Hash().String(), state)
	}()

	contract := &contractv2.ContractV2{
		Address:    p.contractBody.Contract,
		CreateHash: p.tx.Hash(),
		Type:       p.contractBody.Type,
		Body:       nil,
	}

	initBody := p.contractBody.Function.(*exchange_func.PledgeInitBody)
	ex, _ := p.pdState.library.GetExchange(initBody.Exchange)
	pledgeData := exchange2.NewPledge(
		initBody.Exchange, initBody.Receiver, initBody.Admin, ex.Symbol, initBody.MaxSupply,
		initBody.DayMintAmount, initBody.PreMint, initBody.DayRewardAmount, initBody.PledgeMatureTime, p.height,
		)
	contract.Body = pledgeData

	if pledgeData.PreMint != 0{
		p.mintEvent(pledgeData.Receiver, pledgeData.RewardToken, pledgeData.PreMint)
		if err := p.runEvents(); err != nil {
			ERR = err
			return
		}
	}
	p.pdState.library.SetContractV2(contract)
}

func (p *PledgeRunner)PreAddPairPoolVerify() error{
	if p.pdState.body == nil{
		return errors.New("pledge contract does not exist")
	}
	height := p.height
	funcBody, _ := p.contractBody.Function.(*exchange_func.PledgeAddPoolBody)
	if funcBody == nil {
		return errors.New("wrong contractV2 function")
	}
	exchange := p.pdState.body.RewardToken
	pair := funcBody.Pair

	if !p.pdState.body.Admin.IsEqual(p.tx.From()){
		return errors.New("forbidden")
	}

	if height < p.pdState.body.Start{
		return errors.New("invalid height")
	}
	exBody, err := p.pdState.library.GetExchange(exchange)
	if err != nil{
		return errors.New("invalid exchange")
	}
	if !exBody.PairExist(pair){
		return errors.New("invalid pair")
	}
	if p.pdState.body.ExistPairPool(pair){
		return errors.New("the pair already exists")
	}
	pairBody, _ := p.pdState.library.GetPair(pair)
	if !pairBody.Token0.IsEqual(param.Token){
		pairAddress, _ := PairAddress(param.Net, param.Token, pairBody.Token0, exchange)
		_, err = p.pdState.library.GetPair(hasharry.StringToAddress(pairAddress))
		if err != nil{
			return fmt.Errorf("token %s must have a pairing with %s", pairBody.Token0.String(), param.Token.String())
		}
	}
	if !pairBody.Token1.IsEqual(param.Token){
		pairAddress, _ := PairAddress(param.Net, param.Token, pairBody.Token1, exchange)
		_, err = p.pdState.library.GetPair(hasharry.StringToAddress(pairAddress))
		if err != nil{
			return fmt.Errorf("token %s must have a pairing with %s", pairBody.Token0.String(), param.Token.String())
		}
	}
	return nil
}

func (p *PledgeRunner)AddPairPool(){
	var ERR error
	var err error
	state := &types.ContractV2State{State: types.Contract_Success}
	defer func() {
		if ERR != nil {
			state.State = types.Contract_Failed
			state.Error = ERR.Error()
		} else {
			state.Event = p.events
		}
		p.pdState.library.SetContractV2State(p.tx.Hash().String(), state)
	}()
	height := p.height
	funcBody, _ := p.contractBody.Function.(*exchange_func.PledgeAddPoolBody)
	pair := funcBody.Pair

	if p.pdState.body.ExistPairPool(pair){
		ERR = errors.New("the pair already exists")
		return
	}

	mintAmount := p.updatePledge(height)

	if mintAmount != 0{
		p.mintEvent(p.pdState.body.Receiver, p.pdState.body.RewardToken, mintAmount)
	}

	p.pdState.body.AddPirPool(pair)

	if err = p.runEvents(); err != nil {
		ERR = err
		return
	}
	p.update()
}


func (p *PledgeRunner)PreAddPledgeVerify() error{
	if p.pdState.body == nil{
		return errors.New("pledge contract does not exist")
	}
	height := p.height
	funcBody, _ := p.contractBody.Function.(*exchange_func.PledgeAddBody)
	if funcBody == nil {
		return errors.New("wrong contractV2 function")
	}
	exchange := p.pdState.body.RewardToken
	pair := funcBody.Pair
	amount := funcBody.Amount

	if height < p.pdState.body.Start{
		return errors.New("invalid height")
	}
	exBody, err := p.pdState.library.GetExchange(exchange)
	if err != nil{
		return errors.New("invalid exchange")
	}
	if !exBody.PairExist(pair){
		return errors.New("invalid pair")
	}
	if !p.pdState.body.ExistPairPool(pair){
		return errors.New("the pair was not found")
	}

	balance := p.pdState.library.GetBalance(p.tx.From(), pair)
	if balance < amount{
		return fmt.Errorf("insufficient balance %s", p.tx.From().String())
	}
	return nil
}

func (p *PledgeRunner)AddPledge() {
	var ERR error
	var err error
	state := &types.ContractV2State{State: types.Contract_Success}
	defer func() {
		if ERR != nil {
			state.State = types.Contract_Failed
			state.Error = ERR.Error()
		} else {
			state.Event = p.events
		}
		p.pdState.library.SetContractV2State(p.tx.Hash().String(), state)
	}()
	height := p.height
	funcBody, _ := p.contractBody.Function.(*exchange_func.PledgeAddBody)
	pair := funcBody.Pair
	amount := funcBody.Amount
	if !p.pdState.body.ExistPairPool(pair){
		ERR = errors.New("the pair was not found")
		return
	}
	mintAmount := p.updatePledge(height)

	if mintAmount != 0{
		p.mintEvent(p.pdState.body.Receiver, p.pdState.body.RewardToken, mintAmount)
	}

	if err = p.pdState.body.In(height, p.tx.From(), pair, amount);err != nil{
		ERR = err
		return
	}


	p.transferEvent(p.tx.From(), p.address, pair, amount)
	if err = p.runEvents(); err != nil {
		ERR = err
		return
	}
	p.update()
}

func (p *PledgeRunner)PreRemovePledgeVerify() error{
	if p.pdState.body == nil{
		return errors.New("pledge contract does not exist")
	}
	height := p.height
	funcBody, _ := p.contractBody.Function.(*exchange_func.PledgeRemoveBody)
	if funcBody == nil {
		return errors.New("wrong contractV2 function")
	}
	exchange := p.pdState.body.RewardToken
	pair := funcBody.Pair
	amount := funcBody.Amount

	if height < p.pdState.body.Start{
		return errors.New("invalid height")
	}
	exBody, err := p.pdState.library.GetExchange(exchange)
	if err != nil{
		return errors.New("invalid exchange")
	}
	if !exBody.PairExist(pair){
		return errors.New("invalid pair")
	}
	unripe, mature :=  p.pdState.body.GetPledgeAmount(p.tx.From(), pair)
	if amount > unripe + mature{
		return errors.New("insufficient pledge")
	}
	balance := p.pdState.library.GetBalance(p.address, pair)
	if balance < amount{
		return fmt.Errorf("insufficient balance %s", p.tx.From().String())
	}
	return nil
}

func (p *PledgeRunner)RemovePledge() {
	var ERR error
	var err error
	state := &types.ContractV2State{State: types.Contract_Success}
	defer func() {
		if ERR != nil {
			state.State = types.Contract_Failed
			state.Error = ERR.Error()
		} else {
			state.Event = p.events
		}
		p.pdState.library.SetContractV2State(p.tx.Hash().String(), state)
	}()

	height := p.height
	funcBody, _ := p.contractBody.Function.(*exchange_func.PledgeRemoveBody)
	pair := funcBody.Pair
	amount := funcBody.Amount

	mintAmount := p.updatePledge(height)

	if mintAmount != 0{
		p.mintEvent(p.pdState.body.Receiver, p.pdState.body.RewardToken, mintAmount)
	}

	if err = p.pdState.body.Out(p.tx.From(), pair, amount);err != nil{
		ERR = err
		return
	}

	p.transferEvent(p.address, p.tx.From(), pair, amount)
	if err = p.runEvents(); err != nil {
		ERR = err
		return
	}
	p.update()
}

func (p *PledgeRunner)PreRemoveRewardVerify() error{
	if p.pdState.body == nil{
		return errors.New("pledge contract does not exist")
	}
	p.updatePledge(p.height)
	rewards := p.pdState.body.RemoveReward(p.tx.From())
	if len(rewards) == 0{
		return errors.New("no rewards")
	}
	return nil
}

func (p *PledgeRunner)RemoveReward(){
	var ERR error
	var err error
	state := &types.ContractV2State{State: types.Contract_Success}
	defer func() {
		if ERR != nil {
			state.State = types.Contract_Failed
			state.Error = ERR.Error()
		} else {
			state.Event = p.events
		}
		p.pdState.library.SetContractV2State(p.tx.Hash().String(), state)
	}()

	mintAmount := p.updatePledge(p.height)

	if mintAmount != 0{
		p.mintEvent(p.pdState.body.Receiver, p.pdState.body.RewardToken, mintAmount)
	}

	rewards := p.pdState.body.RemoveReward(p.tx.From())
	if len(rewards) == 0{
		ERR = errors.New("no rewards")
		return
	}

	for _, reward := range rewards{
		p.mintEvent(p.tx.From(), p.pdState.body.RewardToken, reward.Amount)
	}

	if err = p.runEvents(); err != nil {
		ERR = err
		return
	}
	p.update()
}

func (p *PledgeRunner)PreUpdatePledgeVerify() error{
	if p.pdState.body == nil{
		return errors.New("pledge contract does not exist")
	}
	if !p.tx.From().IsEqual(p.pdState.body.Admin){
		return errors.New("forbidden")
	}
	return nil
}

func (p *PledgeRunner)UpdatePledge() {
	var ERR error
	state := &types.ContractV2State{State: types.Contract_Success}
	defer func() {
		if ERR != nil {
			state.State = types.Contract_Failed
			state.Error = ERR.Error()
		} else {
			state.Event = p.events
		}
		p.pdState.library.SetContractV2State(p.tx.Hash().String(), state)
	}()

	if !p.tx.From().IsEqual(p.pdState.body.Admin){
		ERR =  errors.New("forbidden")
		return
	}

	mintAmount := p.updatePledge(p.height)
	if mintAmount != 0{
		p.mintEvent(p.pdState.body.Receiver, p.pdState.body.RewardToken, mintAmount)
	}

	if err := p.runEvents(); err != nil {
		ERR = err
		return
	}
	p.update()
}

func (p *PledgeRunner)updatePledge(height uint64) uint64{
	var mintAmount uint64
	if p.pdState.body.IsUpdate(height){
		p.pdState.body.UpdateMature(height)
		pairValue := map[hasharry.Address]uint64{}
		var totalValue uint64
		ex := p.pdState.library.GetContractV2(p.pdState.body.RewardToken.String())
		for pairAddr, total := range p.pdState.body.MaturePair{
			value, _ := p.pdState.getPairValue(pairAddr, total, ex)
			totalValue += value
			pairValue[pairAddr] = value
		}
		mintAmount = p.pdState.body.UpdateMint(p.height, totalValue, pairValue)
	}
	return mintAmount
}


func (p *PledgeRunner) update() {
	p.pdState.header.Body = p.pdState.body
	p.pdState.library.SetContractV2(p.pdState.header)
}

func (p *PledgeRunner) mintEvent(to, token hasharry.Address, amount uint64) {
	p.events = append(p.events, &types.Event{
		EventType: types.Event_Mint,
		From:      hasharry.StringToAddress("mint"),
		To:        to,
		Token:     token,
		Amount:    amount,
		Height:    p.height,
	})
}


func (p *PledgeRunner) transferEvent(from, to, token hasharry.Address, amount uint64) {
	p.events = append(p.events, &types.Event{
		EventType: types.Event_Transfer,
		From:      from,
		To:        to,
		Token:     token,
		Amount:    amount,
		Height:    p.height,
	})
}

func (p *PledgeRunner) runEvents() error {
	for _, event := range p.events {
		if err := p.pdState.library.PreRunEvent(event); err != nil {
			return err
		}
	}
	for _, event := range p.events {
		p.pdState.library.RunEvent(event)
	}
	return nil
}

func PledgeAddress(net, from string, nonce uint64) (string, error) {
	bytes := make([]byte, 0)
	nonceBytes := codec.Uint64toBytes(nonce)
	bytes = append(base58.Decode(from), nonceBytes...)
	return ut.GenerateContractV2Address(net, bytes)
}