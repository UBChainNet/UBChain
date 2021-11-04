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
	_, exist := exMethods[method]
	return exist
}

type PledgeRunner struct {
	pdState      *PledgeState
	address      hasharry.Address
	tx           types.ITransaction
	contractBody *types.TxContractV2Body
	events       []*types.Event
	height       uint64
}

func NewPledgeRunner(lib *library.RunnerLibrary, tx types.ITransaction, height uint64) *PledgeRunner {
	var pd *exchange2.Pledge
	address := tx.GetTxBody().GetContract()
	exHeader := lib.GetContractV2(address.String())
	if exHeader != nil {
		pd = exHeader.Body.(*exchange2.Pledge)
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
	}
}

func (p *PledgeRunner) PreInitVerify() error{
	pd := p.pdState.library.GetContractV2(p.address.String())
	if pd != nil{
		return fmt.Errorf("pledge %s already exist", p.address.String())
	}
	initBody := p.contractBody.Function.(*exchange_func.PledgeInitBody)
	ex, err := p.pdState.library.GetExchange(initBody.Exchange)
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
	contract.Body = exchange2.NewPledge(initBody.Admin, initBody.Exchange, initBody.MaxSupply, initBody.DayMint, p.height)
	p.pdState.library.SetContractV2(contract)
}


func (p *PledgeRunner)PreAddPledgeVerify() error{
	height := p.height
	funcBody, _ := p.contractBody.Function.(*exchange_func.PledgeAddBody)
	if funcBody == nil {
		return errors.New("wrong contractV2 function")
	}
	exchange := p.pdState.body.Reward
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
	pairBody, _ := p.pdState.library.GetPair(pair)
	if !pairBody.Token0.IsEqual(param.Token){
		pairAddress, _ := PairAddress(param.Net, param.Token, pairBody.Token0, exchange)
		_, err := p.pdState.library.GetPair(hasharry.StringToAddress(pairAddress))
		if err != nil{
			return fmt.Errorf("token %s must have a pairing with %s", pairBody.Token0.String(), param.Token.String())
		}
	}
	if !pairBody.Token1.IsEqual(param.Token){
		pairAddress, _ := PairAddress(param.Net, param.Token, pairBody.Token1, exchange)
		_, err := p.pdState.library.GetPair(hasharry.StringToAddress(pairAddress))
		if err != nil{
			return fmt.Errorf("token %s must have a pairing with %s", pairBody.Token0.String(), param.Token.String())
		}
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

	if err = p.pdState.body.In(height, p.tx.From(), pair, amount);err != nil{
		ERR = err
		return
	}

	p.updatePledge(height)

	p.transferEvent(p.tx.From(), p.address, pair, amount)
	if err = p.runEvents(); err != nil {
		ERR = err
		return
	}
	p.update()
}

func (p *PledgeRunner)PreRemovePledgeVerify() error{
	height := p.height
	funcBody, _ := p.contractBody.Function.(*exchange_func.PledgeRemoveBody)
	if funcBody == nil {
		return errors.New("wrong contractV2 function")
	}
	exchange := p.pdState.body.Reward
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

	if err = p.pdState.body.Out(p.tx.From(), pair, amount);err != nil{
		ERR = err
		return
	}

	p.updatePledge(height)

	p.transferEvent(p.address, p.tx.From(), pair, amount)
	if err = p.runEvents(); err != nil {
		ERR = err
		return
	}
	p.update()
}

func (p *PledgeRunner)PreRemoveRewardVerify() error{
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

	rewards := p.pdState.body.RemoveReward(p.tx.From())
	if len(rewards) == 0{
		ERR = errors.New("no rewards")
		return
	}

	p.updatePledge(p.height)

	for _, reward := range rewards{
		p.mintEvent(p.address, reward.Token, reward.Amount)
	}

	if err = p.runEvents(); err != nil {
		ERR = err
		return
	}
	p.update()
}

func (p *PledgeRunner)PreUpdatePledgeVerify() error{
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

	p.updatePledge(p.height)

	p.update()
}

func (p *PledgeRunner)updatePledge(height uint64){
	if p.pdState.body.IsUpdate(height){
		p.pdState.body.UpdateMature(height)
		pairValue := map[hasharry.Address]uint64{}
		var totalValue uint64
		ex := p.pdState.library.GetContractV2(p.pdState.body.Reward.String())
		for pairAddr, total := range p.pdState.body.PledgePair{
			value, _ := p.getPairValue(pairAddr, total, ex)
			totalValue += value
			pairValue[pairAddr] = value
		}
		p.pdState.body.UpdateReward(p.height, totalValue, pairValue)
	}
}

func (p *PledgeRunner)getPairValue(pairAddr hasharry.Address, totalLp uint64, exchange *contractv2.ContractV2)(uint64, error){
	var token0value uint64
	var token1value uint64
	pairState, _ := NewPairStateWithExchange(p.pdState.library, pairAddr.String(), exchange)
	pairTotal, err := pairState.TotalValue(types.Amount(totalLp).ToCoin())
	if err != nil{
		return 0, err
	}
	if pairTotal.Token0 == param.Token.String(){
		token0value, _ = types.NewAmount(pairTotal.Value0)
	}else{
		pairWithMain, _ := PairAddress(param.Net, pairState.pairBody.Token0, param.Token, p.pdState.body.Reward)
		pairWithMainBody, _ := p.pdState.library.GetPair(hasharry.StringToAddress(pairWithMain))
		if pairState.pairBody.Token0.IsEqual(pairWithMainBody.Token0){
			token0value = big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(int64(pairTotal.Value0)), big.NewInt(int64(pairWithMainBody.Reserve1))),
				big.NewInt(int64(pairWithMainBody.Reserve0))).Uint64()
		}else{
			token0value = big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(int64(pairTotal.Value0)), big.NewInt(int64(pairWithMainBody.Reserve0))),
				big.NewInt(int64(pairWithMainBody.Reserve1))).Uint64()
		}
	}
	if pairTotal.Token1 == param.Token.String(){
		token1value, _ = types.NewAmount(pairTotal.Value1)
	}else{
		pairWithMain, _ := PairAddress(param.Net, pairState.pairBody.Token1, param.Token, p.pdState.body.Reward)
		pairWithMainBody, _ := p.pdState.library.GetPair(hasharry.StringToAddress(pairWithMain))
		if pairState.pairBody.Token1.IsEqual(pairWithMainBody.Token0){
			token1value = big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(int64(pairTotal.Value1)), big.NewInt(int64(pairWithMainBody.Reserve1))),
				big.NewInt(int64(pairWithMainBody.Reserve0))).Uint64()
		}else{
			token1value = big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(int64(pairTotal.Value1)), big.NewInt(int64(pairWithMainBody.Reserve0))),
				big.NewInt(int64(pairWithMainBody.Reserve1))).Uint64()
		}
	}
	return token0value + token1value, nil
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