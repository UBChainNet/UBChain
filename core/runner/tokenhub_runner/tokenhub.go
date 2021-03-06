package tokenhub_runner

import (
	"errors"
	"fmt"
	"github.com/UBChainNet/UBChain/common/codec"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/core/runner/library"
	"github.com/UBChainNet/UBChain/core/runner/method"
	"github.com/UBChainNet/UBChain/core/types"
	"github.com/UBChainNet/UBChain/core/types/contractv2"
	"github.com/UBChainNet/UBChain/core/types/contractv2/tokenhub"
	"github.com/UBChainNet/UBChain/core/types/functionbody/tokenhub_func"
	"github.com/UBChainNet/UBChain/crypto/base58"
	"github.com/UBChainNet/UBChain/param"
	"github.com/UBChainNet/UBChain/ut"
)

type TokenHubState struct {
	library *library.RunnerLibrary
	header  *contractv2.ContractV2
	body    *tokenhub.TokenHub
}

func NewTokenHubState(runnerLibrary *library.RunnerLibrary, thAddress string) (*TokenHubState, error) {
	pdHeader := runnerLibrary.GetContractV2(thAddress)
	if pdHeader == nil {
		return nil, fmt.Errorf("tokenhub %s does not exist", thAddress)
	}
	thBody, _ := pdHeader.Body.(*tokenhub.TokenHub)
	return &TokenHubState{
		header:  pdHeader,
		body:    thBody,
		library: runnerLibrary,
	}, nil
}

func (ts *TokenHubState) Methods() map[string]*method.MethodInfo {
	return method.TokenHubMethods
}

func (ts *TokenHubState) MethodExist(mth string) bool {
	_, exist := method.TokenHubMethods[mth]
	return exist
}

type TokenHubRunner struct {
	thState      *TokenHubState
	address      hasharry.Address
	tx           types.ITransaction
	contractBody *types.TxContractV2Body
	events       []*types.Event
	height       uint64
}

func NewTokenHubRunner(lib *library.RunnerLibrary, tx types.ITransaction, height uint64) (*TokenHubRunner, error) {
	var th *tokenhub.TokenHub
	var ok bool
	address := tx.GetTxBody().GetContract()
	thHeader := lib.GetContractV2(address.String())
	if thHeader != nil {
		th, ok = thHeader.Body.(*tokenhub.TokenHub)
		if !ok {
			return nil, errors.New("invalid contract type")
		}
	}

	contractBody := tx.GetTxBody().(*types.TxContractV2Body)
	return &TokenHubRunner{
		thState: &TokenHubState{
			header:  thHeader,
			body:    th,
			library: lib,
		},
		address:      address,
		tx:           tx,
		contractBody: contractBody,
		events:       make([]*types.Event, 0),
		height:       height,
	}, nil
}

func (t *TokenHubRunner) PreInitVerify() error {
	if t.thState.body != nil{
		if t.thState.body.Setter != t.tx.From() {
			return errors.New("forbidden")
		}
	}else{
		address, err := TokenHubAddress(param.Net, t.tx.From().String(), t.tx.GetNonce())
		if err != nil{
			return err
		}
		if address != t.address.String(){
			return errors.New("invalid contract address")
		}
	}
	return nil
}

func (t *TokenHubRunner) Init() {
	var ERR error
	state := &types.ContractV2State{State: types.Contract_Success}
	defer func() {
		if ERR != nil {
			state.State = types.Contract_Failed
			state.Error = ERR.Error()
		} else {
			state.Event = t.events
		}
		t.thState.library.SetContractV2State(t.tx.Hash().String(), state)
	}()

	initBody := t.contractBody.Function.(*tokenhub_func.TokenHubInitBody)
	if t.thState.header == nil {
		t.thState.header = &contractv2.ContractV2{
			Address:    t.contractBody.Contract,
			CreateHash: t.tx.Hash(),
			Type:       t.contractBody.Type,
			Body: tokenhub.NewTokenHub(
				t.contractBody.Contract,
				initBody.Setter,
				initBody.Admin,
				initBody.FeeTo,
				initBody.FeeRate,
			),
		}
	} else {
		if err := t.thState.body.SetSetter(t.tx.From(), initBody.Setter); err != nil {
			ERR = err
			return
		}
		if err := t.thState.body.SetAdmin(t.tx.From(), initBody.Admin); err != nil {
			ERR = err
			return
		}
		if err := t.thState.body.SetFeeTo(t.tx.From(), initBody.FeeTo); err != nil {
			ERR = err
			return
		}
		if err := t.thState.body.SetFeeRate(t.tx.From(), initBody.FeeRate); err != nil {
			ERR = err
			return
		}
		t.thState.header.Body = t.thState.body
	}

	t.thState.library.SetContractV2(t.thState.header)
}

func (t *TokenHubRunner) PreAckVerify() error {
	if t.thState.body == nil{
		return fmt.Errorf("tokenhub %s does not exist", t.address.String())
	}
	ackBody := t.contractBody.Function.(*tokenhub_func.TokenHubAckBody)
	ackType := make(map[uint64]tokenhub.AckType)
	ackHash := make(map[uint64]string)
	for i, sequence := range ackBody.Sequences {
		ackType[sequence] = tokenhub.AckType(ackBody.AckTypes[i])
		ackHash[sequence] = ackBody.Hashes[i]
	}
	_, err := t.thState.body.AckTransfer(t.tx.From(), ackType, ackHash)
	return err
}

func (t *TokenHubRunner) AckTransfer() {
	var ERR error
	state := &types.ContractV2State{State: types.Contract_Success}
	defer func() {
		if ERR != nil {
			state.State = types.Contract_Failed
			state.Error = ERR.Error()
		} else {
			state.Event = t.events
		}
		t.thState.library.SetContractV2State(t.tx.Hash().String(), state)
	}()

	ackBody := t.contractBody.Function.(*tokenhub_func.TokenHubAckBody)
	ackData := make(map[uint64]tokenhub.AckType)
	ackHash := make(map[uint64]string)
	for i, sequence := range ackBody.Sequences {
		ackData[sequence] = tokenhub.AckType(ackBody.AckTypes[i])
		ackHash[sequence] = ackBody.Hashes[i]
	}
	transfers, err := t.thState.body.AckTransfer(t.tx.From(), ackData, ackHash)
	if err != nil {
		ERR = err
		return
	}
	for _, transfer := range transfers {
		t.transferEvent(transfer.From, transfer.To, transfer.Token, transfer.Amount)
	}
	if err = t.runEvents(); err != nil {
		ERR = err
		return
	}
	t.update()
}

func (t *TokenHubRunner) PreTransferOutVerify() error {
	if t.thState.body == nil{
		return fmt.Errorf("tokenhub %s does not exist", t.address.String())
	}
	trBody := t.contractBody.Function.(*tokenhub_func.TokenHubTransferOutBody)
	balance := t.thState.library.GetBalance(t.tx.From(), param.Token)
	if trBody.Amount > balance{
		return fmt.Errorf("insufficient balance %s", t.tx.From().String())
	}
	return nil
}

func (t *TokenHubRunner) TransferOut() {
	var ERR error
	state := &types.ContractV2State{State: types.Contract_Success}
	defer func() {
		if ERR != nil {
			state.State = types.Contract_Failed
			state.Error = ERR.Error()
		} else {
			state.Event = t.events
		}
		t.thState.library.SetContractV2State(t.tx.Hash().String(), state)
	}()
	trBody := t.contractBody.Function.(*tokenhub_func.TokenHubTransferOutBody)
	transfers, err := t.thState.body.TransferOut(t.tx.From(), trBody.To, trBody.Amount)
	if err != nil {
		ERR = err
		return
	}
	for _, transfer := range transfers {
		t.transferEvent(transfer.From, transfer.To, transfer.Token, transfer.Amount)
	}
	if err = t.runEvents(); err != nil {
		ERR = err
		return
	}
	t.update()
}

func (t *TokenHubRunner) PreTransferInVerify() error {
	if t.thState.body == nil{
		return fmt.Errorf("tokenhub %s does not exist", t.address.String())
	}
	if t.tx.From() != t.thState.body.Admin{
		return errors.New("forbidden")
	}
	trBody := t.contractBody.Function.(*tokenhub_func.TokenHubTransferInBody)
	balance := t.thState.library.GetBalance(t.thState.body.Address, param.Token)
	if trBody.Amount > balance{
		return fmt.Errorf("insufficient balance %s", t.thState.body.Address.String())
	}
	_, exist := t.thState.body.AcrossSeqs[trBody.AcrossSeq]
	if exist{
		return fmt.Errorf("%d across seq has send", trBody.AcrossSeq)
	}
	return nil
}

func (t *TokenHubRunner) TransferIn() {
	var ERR error
	state := &types.ContractV2State{State: types.Contract_Success}
	defer func() {
		if ERR != nil {
			state.State = types.Contract_Failed
			state.Error = ERR.Error()
		} else {
			state.Event = t.events
		}
		t.thState.library.SetContractV2State(t.tx.Hash().String(), state)
	}()

	trBody := t.contractBody.Function.(*tokenhub_func.TokenHubTransferInBody)
	transfers, err := t.thState.body.TransferIn(t.tx.From(), trBody.To, trBody.Amount, trBody.AcrossSeq, t.tx.Hash().String())
	if err != nil {
		ERR = err
		return
	}
	for _, transfer := range transfers {
		t.transferEvent(transfer.From, transfer.To, transfer.Token, transfer.Amount)
	}
	if err = t.runEvents(); err != nil {
		ERR = err
		return
	}
	t.update()
}

func (t *TokenHubRunner) PreFinishAcrossVerify() error {
	if t.thState.body == nil{
		return fmt.Errorf("tokenhub %s does not exist", t.address.String())
	}
	if t.tx.From() != t.thState.body.Admin{
		return errors.New("forbidden")
	}
	return nil
}

func (t *TokenHubRunner) FinishAcross() {
	var ERR error
	state := &types.ContractV2State{State: types.Contract_Success}
	defer func() {
		if ERR != nil {
			state.State = types.Contract_Failed
			state.Error = ERR.Error()
		} else {
			state.Event = t.events
		}
		t.thState.library.SetContractV2State(t.tx.Hash().String(), state)
	}()

	trBody := t.contractBody.Function.(*tokenhub_func.TokenHubFinishAcrossBody)
	err := t.thState.body.AcrossFinished(t.tx.From(), trBody.AcrossSeqs)
	if err != nil {
		ERR = err
		return
	}
	t.update()
}

func (t *TokenHubRunner) update() {
	t.thState.header.Body = t.thState.body
	t.thState.library.SetContractV2(t.thState.header)
}

func (t *TokenHubRunner) transferEvent(from, to, token hasharry.Address, amount uint64) {
	t.events = append(t.events, &types.Event{
		EventType: types.Event_Transfer,
		From:      from,
		To:        to,
		Token:     token,
		Amount:    amount,
		Height:    t.height,
	})
}

func (t *TokenHubRunner) runEvents() error {
	for _, event := range t.events {
		if err := t.thState.library.PreRunEvent(event); err != nil {
			return err
		}
	}
	for _, event := range t.events {
		t.thState.library.RunEvent(event)
	}
	return nil
}

func TokenHubAddress(net, from string, nonce uint64) (string, error) {
	bytes := make([]byte, 0)
	nonceBytes := codec.Uint64toBytes(nonce)
	bytes = append(base58.Decode(from), nonceBytes...)
	return ut.GenerateContractV2Address(net, bytes)
}
