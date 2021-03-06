package types

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/core/types/contractv2"
	"github.com/UBChainNet/UBChain/core/types/functionbody/exchange_func"
	"github.com/UBChainNet/UBChain/core/types/functionbody/tokenhub_func"
)

type IRpcTransactionBody interface {
}

type RpcTransactionHead struct {
	TxHash     string          `json:"txhash"`
	TxType     TransactionType `json:"txtype"`
	From       string          `json:"from"`
	Nonce      uint64          `json:"nonce"`
	Fees       uint64          `json:"fees"`
	Time       uint64          `json:"time"`
	Note       string          `json:"note"`
	SignScript *RpcSignScript  `json:"signscript"`
}

type RpcTransaction struct {
	TxHead *RpcTransactionHead `json:"txhead"`
	TxBody IRpcTransactionBody `json:"txbody"`
}

type RpcTransactionConfirmed struct {
	TxHead    *RpcTransactionHead `json:"txhead"`
	TxBody    IRpcTransactionBody `json:"txbody"`
	Height    uint64              `json:"height"`
	Confirmed bool                `json:"confirmed"`
}

type RpcSignScript struct {
	Signature string `json:"signature"`
	PubKey    string `json:"pubkey"`
}

func (th *RpcTransactionHead) FromBytes() []byte {
	return []byte(th.From)
}

func TranslateRpcTxToTx(rpcTx *RpcTransaction) (*Transaction, error) {
	var err error
	txHash, err := hasharry.StringToHash(rpcTx.TxHead.TxHash)
	if err != nil {
		return nil, err
	}
	signScript, err := TranslateRpcSignScriptToSignScript(rpcTx.TxHead.SignScript)
	if err != nil {
		return nil, err
	}
	var txBody ITransactionBody
	switch rpcTx.TxHead.TxType {
	case Transfer_:
		body := &RpcTransferBody{}
		bytes, err := json.Marshal(rpcTx.TxBody)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bytes, body)
		if err != nil {
			return nil, err
		}
		txBody, err = translateRpcTransferBodyToBody(body)
	case Contract_:
		body := &RpcContractBody{}
		bytes, err := json.Marshal(rpcTx.TxBody)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bytes, body)
		if err != nil {
			return nil, err
		}
		txBody, err = translateRpcContractBodyToBody(body)
	case ContractV2_:
		txBody, err = translateRpcContractV2BodyToBody(rpcTx.TxBody)
		/*case types.LoginCandidate_:
			txBody, err = translateRpcLoginBodyToBody(rpcTx.LoginBody)
		case types.LogoutCandidate:
			txBody = &types.LogoutTransactionBody{}
		case types.VoteToCandidate:
			txBody, err = translateRpcVoteBodyToBody(rpcTx.VoteBody)*/
	}
	tx := &Transaction{
		TxHead: &TransactionHead{
			TxHash:     txHash,
			TxType:     rpcTx.TxHead.TxType,
			From:       hasharry.StringToAddress(rpcTx.TxHead.From),
			Nonce:      rpcTx.TxHead.Nonce,
			Fees:       rpcTx.TxHead.Fees,
			Time:       rpcTx.TxHead.Time,
			Note:       rpcTx.TxHead.Note,
			SignScript: signScript,
		},
		TxBody: txBody,
	}
	return tx, nil
}

func TranslateTxToRpcTx(tx *Transaction) (*RpcTransaction, error) {
	var err error
	rpcTx := &RpcTransaction{
		TxHead: &RpcTransactionHead{
			TxHash: tx.Hash().String(),
			TxType: tx.GetTxType(),
			From:   addressToString(tx.From()),
			Nonce:  tx.GetNonce(),
			Fees:   tx.GetFees(),
			Time:   tx.GetTime(),
			Note:   tx.GetNote(),
			SignScript: &RpcSignScript{
				Signature: hex.EncodeToString(tx.GetSignScript().Signature),
				PubKey:    hex.EncodeToString(tx.GetSignScript().PubKey),
			}},
		TxBody: nil,
	}
	switch tx.GetTxType() {
	case Transfer_:
		rpcBody := &RpcTransferBody{
			Contract:  tx.GetTxBody().GetContract().String(),
			Receivers: make([]RpcReceiver, 0),
		}
		for _, re := range tx.GetTxBody().ToAddress().ReceiverList() {
			rpcBody.Receivers = append(rpcBody.Receivers, RpcReceiver{
				Address: re.Address.String(),
				Amount:  re.Amount,
			})
		}
		rpcTx.TxBody = rpcBody
	case Contract_:
		rpcTx.TxBody = &RpcContractBody{
			Contract:    tx.GetTxBody().GetContract().String(),
			To:          tx.GetTxBody().ToAddress().ReceiverList()[0].Address.String(),
			Name:        tx.GetTxBody().GetName(),
			Abbr:        tx.GetTxBody().GetAbbr(),
			Description: tx.GetTxBody().GetDescription(),
			Increase:    tx.GetTxBody().GetIncreaseSwitch(),
			Amount:      tx.GetTxBody().GetAmount(),
		}
	case ContractV2_:
		body, ok := tx.GetTxBody().(*TxContractV2Body)
		if !ok {
			return nil, errors.New("wrong transaction body")
		}
		rpcTx.TxBody, err = translateContractV2ToRpcContractV2(body)
		if err != nil {
			return nil, err
		}
	case LoginCandidate_:
		rpcTx.TxBody = &RpcLoginTransactionBody{
			PeerId: string(tx.GetTxBody().GetPeerId()),
		}
		/*case types.LogoutCandidate:
			rpcTx.LogoutBody = &RpcLogoutTransactionBody{}
		case types.VoteToCandidate:
			rpcTx.VoteBody = &RpcVoteTransactionBody{To: tx.GetTxBody().ToAddress().String()}
		*/
	}

	return rpcTx, nil
}

func translateRpcContractV2BodyToBody(rpcBody IRpcTransactionBody) (*TxContractV2Body, error) {
	if rpcBody == nil {
		return nil, errors.New("wrong contract transaction body")
	}
	body := &RpcContractV2TransactionBody{}
	bytes, err := json.Marshal(rpcBody)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, body)
	if err != nil {
		return nil, err
	}
	switch body.FunctionType {
	case contractv2.Exchange_Init:
		bytes, err := json.Marshal(body.Function)
		if err != nil {
			return nil, err
		}
		init := &RpcExchangeInitBody{
			Symbol: "",
			Admin:  "",
			FeeTo:  "",
		}
		err = json.Unmarshal(bytes, init)
		if err != nil {
			return nil, err
		}
		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function: &exchange_func.ExchangeInitBody{
				Admin:  hasharry.StringToAddress(init.Admin),
				FeeTo:  hasharry.StringToAddress(init.FeeTo),
				Symbol: init.Symbol,
			},
		}, nil
	case contractv2.Exchange_SetAdmin:
		bytes, err := json.Marshal(body.Function)
		if err != nil {
			return nil, err
		}
		setBody := &RpcExchangeSetAdminBody{
			Address: "",
		}
		err = json.Unmarshal(bytes, setBody)
		if err != nil {
			return nil, err
		}
		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function: &exchange_func.ExchangeAdmin{
				Address: hasharry.StringToAddress(setBody.Address),
			},
		}, nil
	case contractv2.Exchange_SetFeeTo:
		bytes, err := json.Marshal(body.Function)
		if err != nil {
			return nil, err
		}
		setBody := &RpcExchangeSetFeeToBody{
			Address: "",
		}
		err = json.Unmarshal(bytes, setBody)
		if err != nil {
			return nil, err
		}
		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function: &exchange_func.ExchangeFeeTo{
				Address: hasharry.StringToAddress(setBody.Address),
			},
		}, nil
	case contractv2.Exchange_ExactIn:
		bytes, err := json.Marshal(body.Function)
		if err != nil {
			return nil, err
		}
		inBody := &RpcExchangeExactInBody{}
		err = json.Unmarshal(bytes, inBody)
		if err != nil {
			return nil, err
		}
		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function: &exchange_func.ExactIn{
				AmountIn:     inBody.AmountIn,
				AmountOutMin: inBody.AmountOutMin,
				Path:         addrListToHashAddr(inBody.Path),
				To:           hasharry.StringToAddress(inBody.To),
				Deadline:     inBody.Deadline,
			},
		}, nil
	case contractv2.Exchange_ExactOut:
		bytes, err := json.Marshal(body.Function)
		if err != nil {
			return nil, err
		}
		outBody := &RpcExchangeExactOutBody{}
		err = json.Unmarshal(bytes, outBody)
		if err != nil {
			return nil, err
		}
		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function: &exchange_func.ExactOut{
				AmountOut:   outBody.AmountOut,
				AmountInMax: outBody.AmountInMax,
				Path:        addrListToHashAddr(outBody.Path),
				To:          hasharry.StringToAddress(outBody.To),
				Deadline:    outBody.Deadline,
			},
		}, nil
	case contractv2.Pair_AddLiquidity:
		bytes, err := json.Marshal(body.Function)
		if err != nil {
			return nil, err
		}
		createBody := &RpcExchangeAddLiquidity{}
		err = json.Unmarshal(bytes, createBody)
		if err != nil {
			return nil, err
		}
		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function: &exchange_func.ExchangeAddLiquidity{
				Exchange:       hasharry.StringToAddress(createBody.Exchange),
				TokenA:         hasharry.StringToAddress(createBody.TokenA),
				TokenB:         hasharry.StringToAddress(createBody.TokenB),
				To:             hasharry.StringToAddress(createBody.To),
				AmountADesired: createBody.AmountADesired,
				AmountBDesired: createBody.AmountBDesired,
				AmountAMin:     createBody.AmountAMin,
				AmountBMin:     createBody.AmountBMin,
				Deadline:       createBody.Deadline,
			},
		}, nil
	case contractv2.Pair_RemoveLiquidity:
		bytes, err := json.Marshal(body.Function)
		if err != nil {
			return nil, err
		}
		remove := &RpcExchangeRemoveLiquidity{}
		err = json.Unmarshal(bytes, remove)
		if err != nil {
			return nil, err
		}

		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function: &exchange_func.ExchangeRemoveLiquidity{
				Exchange:   hasharry.StringToAddress(remove.Exchange),
				TokenA:     hasharry.StringToAddress(remove.TokenA),
				TokenB:     hasharry.StringToAddress(remove.TokenB),
				To:         hasharry.StringToAddress(remove.To),
				Liquidity:  remove.Liquidity,
				AmountAMin: remove.AmountAMin,
				AmountBMin: remove.AmountBMin,
				Deadline:   remove.Deadline,
			},
		}, nil
	case contractv2.Pledge_Init:
		bytes, err = json.Marshal(body.Function)
		if err != nil {
			return nil, err
		}
		init := &RpcPledgeInit{}
		err = json.Unmarshal(bytes, init)
		if err != nil {
			return nil, err
		}

		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function: &exchange_func.PledgeInitBody{
				Exchange:  hasharry.StringToAddress(init.Exchange),
				Receiver:  hasharry.StringToAddress(init.Receiver),
				Admin:     hasharry.StringToAddress(init.Admin),
				PreMint:   init.PreMint,
				MaxSupply: init.MaxSupply,
			},
		}, nil
	case contractv2.Pledge_Start:
		bytes, err = json.Marshal(body.Function)
		if err != nil {
			return nil, err
		}
		start := &RpcPledgeStart{}
		err = json.Unmarshal(bytes, start)
		if err != nil {
			return nil, err
		}

		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function: &exchange_func.PledgeStartBody{
				BlockMintAmount:  start.BlockMintAmount,
				PledgeMatureTime: start.PledgeMatureTime,
			},
		}, nil
	case contractv2.Pledge_AddPool:
		bytes, err = json.Marshal(body.Function)
		if err != nil {
			return nil, err
		}
		addPool := &RpcPledgeAddPool{}
		err = json.Unmarshal(bytes, addPool)
		if err != nil {
			return nil, err
		}

		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function: &exchange_func.PledgeAddPoolBody{
				Pair:        hasharry.StringToAddress(addPool.Pair),
				BlockReward: addPool.BlockReward,
			},
		}, nil
	case contractv2.Pledge_RemovePool:
		bytes, err = json.Marshal(body.Function)
		if err != nil {
			return nil, err
		}
		removePool := &RpcPledgeRemovePool{}
		err = json.Unmarshal(bytes, removePool)
		if err != nil {
			return nil, err
		}

		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function: &exchange_func.PledgeRemovePoolBody{
				Pair: hasharry.StringToAddress(removePool.Pair),
			},
		}, nil
	case contractv2.Pledge_Add:
		bytes, err = json.Marshal(body.Function)
		if err != nil {
			return nil, err
		}
		add := &RpcPledgeAdd{}
		err = json.Unmarshal(bytes, add)
		if err != nil {
			return nil, err
		}

		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function: &exchange_func.PledgeAddBody{
				Pair:   hasharry.StringToAddress(add.Pair),
				Amount: add.Amount,
			},
		}, nil
	case contractv2.Pledge_Remove:
		bytes, err = json.Marshal(body.Function)
		if err != nil {
			return nil, err
		}
		remove := &RpcPledgeRemove{}
		err = json.Unmarshal(bytes, remove)
		if err != nil {
			return nil, err
		}

		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function: &exchange_func.PledgeRemoveBody{
				Pair:   hasharry.StringToAddress(remove.Pair),
				Amount: remove.Amount,
			},
		}, nil
	case contractv2.Pledge_RemoveReward:
		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function:     &exchange_func.PledgeRewardRemoveBody{},
		}, nil
	case contractv2.Pledge_Update:
		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function:     &exchange_func.PledgeUpdateBody{},
		}, nil
	case contractv2.TokenHub_init:
		bytes, err = json.Marshal(body.Function)
		if err != nil {
			return nil, err
		}
		init := &RpcTokenHubInit{}
		err = json.Unmarshal(bytes, init)
		if err != nil {
			return nil, err
		}
		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function: &tokenhub_func.TokenHubInitBody{
				Setter:  hasharry.StringToAddress(init.Setter),
				Admin:   hasharry.StringToAddress(init.Admin),
				FeeTo:   hasharry.StringToAddress(init.FeeTo),
				FeeRate: init.FeeRate,
			},
		}, nil
	case contractv2.TokenHub_Ack:
		bytes, err = json.Marshal(body.Function)
		if err != nil {
			return nil, err
		}
		ack := &RpcTokenHubAck{}
		err = json.Unmarshal(bytes, ack)
		if err != nil {
			return nil, err
		}
		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function: &tokenhub_func.TokenHubAckBody{
				Sequences: ack.Sequences,
				AckTypes:  ack.AckTypes,
				Hashes:    ack.Hashes,
			},
		}, nil
	case contractv2.TokenHub_TransferOut:
		bytes, err = json.Marshal(body.Function)
		if err != nil {
			return nil, err
		}
		tr := &RpcTokenHubTransferOut{}
		err = json.Unmarshal(bytes, tr)
		if err != nil {
			return nil, err
		}
		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function: &tokenhub_func.TokenHubTransferOutBody{
				To:     tr.To,
				Amount: tr.Amount,
			},
		}, nil
	case contractv2.TokenHub_TransferIn:
		bytes, err = json.Marshal(body.Function)
		if err != nil {
			return nil, err
		}
		tr := &RpcTokenHubTransferIn{}
		err = json.Unmarshal(bytes, tr)
		if err != nil {
			return nil, err
		}
		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function: &tokenhub_func.TokenHubTransferInBody{
				To:        hasharry.StringToAddress(tr.To),
				Amount:    tr.Amount,
				AcrossSeq: tr.AcrossSeq,
			},
		}, nil
	case contractv2.TokenHub_FinishAcross:
		bytes, err = json.Marshal(body.Function)
		if err != nil {
			return nil, err
		}
		tr := &RpcTokenHubFinishAcrossBody{}
		err = json.Unmarshal(bytes, tr)
		if err != nil {
			return nil, err
		}
		return &TxContractV2Body{
			Contract:     hasharry.StringToAddress(body.Contract),
			Type:         body.Type,
			FunctionType: body.FunctionType,
			Function: &tokenhub_func.TokenHubFinishAcrossBody{
				AcrossSeqs: tr.AcrossSeqs,
			},
		}, nil
	}
	return nil, errors.New("wrong transaction body")
}

func TranslateContractV2TxToRpcTx(tx *Transaction, state *ContractV2State) (*RpcTransaction, error) {
	var err error
	rpcTx := &RpcTransaction{
		TxHead: &RpcTransactionHead{
			TxHash: tx.Hash().String(),
			TxType: tx.GetTxType(),
			From:   addressToString(tx.From()),
			Nonce:  tx.GetNonce(),
			Fees:   tx.GetFees(),
			Time:   tx.GetTime(),
			Note:   tx.GetNote(),
			SignScript: &RpcSignScript{
				Signature: hex.EncodeToString(tx.GetSignScript().Signature),
				PubKey:    hex.EncodeToString(tx.GetSignScript().PubKey),
			}},
		TxBody: nil,
	}
	switch tx.GetTxType() {
	case ContractV2_:
		body, ok := tx.GetTxBody().(*TxContractV2Body)
		if !ok {
			return nil, errors.New("wrong transaction body")
		}
		rpcTx.TxBody, err = translateToRpcContractV2WithState(body, state)
		if err != nil {
			return nil, err
		}
	}

	return rpcTx, nil
}

func translateToRpcContractV2WithState(body *TxContractV2Body, contractState *ContractV2State) (*RpcContractV2BodyWithState, error) {
	var state *RpcContractState = &RpcContractState{
		StateCode: Contract_Wait,
		Events:    make([]*RpcEvent, 0),
		Error:     "",
	}
	if contractState != nil {
		state.StateCode = contractState.State
		state.Error = contractState.Error
		if contractState.Event != nil {
			for _, e := range contractState.Event {
				state.Events = append(state.Events, &RpcEvent{
					EventType: int(e.EventType),
					From:      e.From.String(),
					To:        e.To.String(),
					Token:     e.Token.String(),
					Amount:    e.Amount,
					Height:    e.Height,
				})
			}
		}
	}

	funcBody, err := rpcFunction(body)
	if err != nil {
		return nil, err
	}
	return &RpcContractV2BodyWithState{
		Contract:     body.Contract.String(),
		Type:         body.Type,
		FunctionType: body.FunctionType,
		Function:     funcBody,
		State:        state,
	}, nil
}

func translateContractV2ToRpcContractV2(body *TxContractV2Body) (*RpcContractV2TransactionBody, error) {
	funcBody, err := rpcFunction(body)
	if err != nil {
		return nil, err
	}
	return &RpcContractV2TransactionBody{
		Contract:     body.Contract.String(),
		Type:         body.Type,
		FunctionType: body.FunctionType,
		Function:     funcBody,
	}, nil
}

func rpcFunction(body *TxContractV2Body) (IRCFunction, error) {
	if body == nil {
		return nil, fmt.Errorf("invalid contract body")
	}
	var function IRCFunction
	switch body.FunctionType {
	case contractv2.Exchange_Init:
		funcBody, ok := body.Function.(*exchange_func.ExchangeInitBody)
		if !ok {
			return nil, errors.New("wrong function body")
		}
		function = &RpcExchangeInitBody{
			Admin:  funcBody.Admin.String(),
			FeeTo:  funcBody.FeeTo.String(),
			Symbol: funcBody.Symbol,
		}
	case contractv2.Exchange_SetAdmin:
		funcBody, ok := body.Function.(*exchange_func.ExchangeAdmin)
		if !ok {
			return nil, errors.New("wrong function body")
		}
		function = &RpcExchangeSetAdminBody{
			Address: funcBody.Address.String(),
		}
	case contractv2.Exchange_SetFeeTo:
		funcBody, ok := body.Function.(*exchange_func.ExchangeFeeTo)
		if !ok {
			return nil, errors.New("wrong function body")
		}
		function = &RpcExchangeSetFeeToBody{
			Address: funcBody.Address.String(),
		}
	case contractv2.Exchange_ExactIn:
		funcBody, ok := body.Function.(*exchange_func.ExactIn)
		if !ok {
			return nil, errors.New("wrong function body")
		}
		function = &RpcExchangeExactInBody{
			AmountIn:     funcBody.AmountIn,
			AmountOutMin: funcBody.AmountOutMin,
			Path:         hashAddrToAddr(funcBody.Path),
			To:           funcBody.To.String(),
			Deadline:     funcBody.Deadline,
		}
	case contractv2.Exchange_ExactOut:
		funcBody, ok := body.Function.(*exchange_func.ExactOut)
		if !ok {
			return nil, errors.New("wrong function body")
		}
		function = &RpcExchangeExactOutBody{
			AmountOut:   funcBody.AmountOut,
			AmountInMax: funcBody.AmountInMax,
			Path:        hashAddrToAddr(funcBody.Path),
			To:          funcBody.To.String(),
			Deadline:    funcBody.Deadline,
		}
	case contractv2.Pair_AddLiquidity:
		funcBody, ok := body.Function.(*exchange_func.ExchangeAddLiquidity)
		if !ok {
			return nil, errors.New("wrong function body")
		}
		function = &RpcExchangeAddLiquidity{
			Exchange:       funcBody.Exchange.String(),
			TokenA:         funcBody.TokenA.String(),
			TokenB:         funcBody.TokenB.String(),
			To:             funcBody.To.String(),
			AmountADesired: funcBody.AmountADesired,
			AmountBDesired: funcBody.AmountBDesired,
			AmountAMin:     funcBody.AmountAMin,
			AmountBMin:     funcBody.AmountBMin,
			Deadline:       funcBody.Deadline,
		}
	case contractv2.Pair_RemoveLiquidity:
		funcBody, ok := body.Function.(*exchange_func.ExchangeRemoveLiquidity)
		if !ok {
			return nil, errors.New("wrong function body")
		}
		function = &RpcExchangeRemoveLiquidity{
			Exchange:   funcBody.Exchange.String(),
			TokenA:     funcBody.TokenA.String(),
			TokenB:     funcBody.TokenB.String(),
			To:         funcBody.To.String(),
			Liquidity:  funcBody.Liquidity,
			AmountAMin: funcBody.AmountAMin,
			AmountBMin: funcBody.AmountBMin,
			Deadline:   funcBody.Deadline,
		}
	case contractv2.Pledge_Init:
		funcBody, ok := body.Function.(*exchange_func.PledgeInitBody)
		if !ok {
			return nil, errors.New("wrong function body")
		}
		function = &RpcPledgeInit{
			Exchange:  funcBody.Exchange.String(),
			Receiver:  funcBody.Receiver.String(),
			Admin:     funcBody.Admin.String(),
			PreMint:   funcBody.PreMint,
			MaxSupply: funcBody.MaxSupply,
		}
	case contractv2.Pledge_Start:
		funcBody, ok := body.Function.(*exchange_func.PledgeStartBody)
		if !ok {
			return nil, errors.New("wrong function body")
		}
		function = &RpcPledgeStart{
			BlockMintAmount:  funcBody.BlockMintAmount,
			PledgeMatureTime: funcBody.PledgeMatureTime,
		}
	case contractv2.Pledge_AddPool:
		funcBody, ok := body.Function.(*exchange_func.PledgeAddPoolBody)
		if !ok {
			return nil, errors.New("wrong function body")
		}
		function = &RpcPledgeAddPool{
			Pair:        funcBody.Pair.String(),
			BlockReward: funcBody.BlockReward,
		}
	case contractv2.Pledge_RemovePool:
		funcBody, ok := body.Function.(*exchange_func.PledgeRemovePoolBody)
		if !ok {
			return nil, errors.New("wrong function body")
		}
		function = &RpcPledgeRemovePool{
			Pair: funcBody.Pair.String(),
		}
	case contractv2.Pledge_Add:
		funcBody, ok := body.Function.(*exchange_func.PledgeAddBody)
		if !ok {
			return nil, errors.New("wrong function body")
		}
		function = &RpcPledgeAdd{
			Pair:   funcBody.Pair.String(),
			Amount: funcBody.Amount,
		}
	case contractv2.Pledge_Remove:
		funcBody, ok := body.Function.(*exchange_func.PledgeRemoveBody)
		if !ok {
			return nil, errors.New("wrong function body")
		}
		function = &RpcPledgeRemove{
			Pair:   funcBody.Pair.String(),
			Amount: funcBody.Amount,
		}
	case contractv2.Pledge_RemoveReward:
		function = &RpcPledgeRewardRemove{}
	case contractv2.Pledge_Update:
		function = &RpcPledgeUpdate{}
	case contractv2.TokenHub_init:
		funcBody, ok := body.Function.(*tokenhub_func.TokenHubInitBody)
		if !ok {
			return nil, errors.New("wrong function body")
		}
		function = &RpcTokenHubInit{
			Setter:  funcBody.Setter.String(),
			Admin:   funcBody.Admin.String(),
			FeeTo:   funcBody.FeeTo.String(),
			FeeRate: funcBody.FeeRate,
		}
	case contractv2.TokenHub_Ack:
		funcBody, ok := body.Function.(*tokenhub_func.TokenHubAckBody)
		if !ok {
			return nil, errors.New("wrong function body")
		}
		function = &RpcTokenHubAck{
			Sequences: funcBody.Sequences,
			AckTypes:  funcBody.AckTypes,
			Hashes:    funcBody.Hashes,
		}
	case contractv2.TokenHub_TransferOut:
		funcBody, ok := body.Function.(*tokenhub_func.TokenHubTransferOutBody)
		if !ok {
			return nil, errors.New("wrong function body")
		}
		function = &RpcTokenHubTransferOut{
			To:     funcBody.To,
			Amount: funcBody.Amount,
		}
	case contractv2.TokenHub_TransferIn:
		funcBody, ok := body.Function.(*tokenhub_func.TokenHubTransferInBody)
		if !ok {
			return nil, errors.New("wrong function body")
		}
		function = &RpcTokenHubTransferIn{
			To:        funcBody.To.String(),
			Amount:    funcBody.Amount,
			AcrossSeq: funcBody.AcrossSeq,
		}
	case contractv2.TokenHub_FinishAcross:
		funcBody, ok := body.Function.(*tokenhub_func.TokenHubFinishAcrossBody)
		if !ok {
			return nil, errors.New("wrong function body")
		}
		function = &RpcTokenHubFinishAcrossBody{
			AcrossSeqs: funcBody.AcrossSeqs,
		}
	}
	return function, nil
}

func TranslateRpcSignScriptToSignScript(rpcSignScript *RpcSignScript) (*SignScript, error) {
	if rpcSignScript == nil {
		return nil, ErrNoSignature
	}
	if rpcSignScript.Signature == "" || rpcSignScript.PubKey == "" {
		return nil, ErrWrongSignature
	}
	signature, err := hex.DecodeString(rpcSignScript.Signature)
	if err != nil {
		return nil, err
	}
	pubKey, err := hex.DecodeString(rpcSignScript.PubKey)
	if err != nil {
		return nil, err
	}
	return &SignScript{
		Signature: signature,
		PubKey:    pubKey,
	}, nil
}

func translateRpcTransferBodyToBody(rpcBody *RpcTransferBody) (*TransferBody, error) {
	if rpcBody == nil {
		return nil, errors.New("wrong transaction body")
	}
	txBody := &TransferBody{
		Contract:  hasharry.StringToAddress(rpcBody.Contract),
		Receivers: NewReceivers(),
	}
	for _, re := range rpcBody.Receivers {
		txBody.Receivers.Add(hasharry.StringToAddress(re.Address), re.Amount)
	}
	return txBody, nil
}

func translateRpcContractBodyToBody(rpcBody *RpcContractBody) (*ContractBody, error) {
	if rpcBody == nil {
		return nil, errors.New("wrong contract transaction body")
	}

	return &ContractBody{
		Contract:       hasharry.StringToAddress(rpcBody.Contract),
		To:             hasharry.StringToAddress(rpcBody.To),
		Abbr:           rpcBody.Abbr,
		IncreaseSwitch: rpcBody.Increase,
		Name:           rpcBody.Name,
		Description:    rpcBody.Description,
		Amount:         rpcBody.Amount,
	}, nil
}

func translateRpcLoginBodyToBody(rpcBody *RpcLoginTransactionBody) (*LoginTransactionBody, error) {
	if rpcBody == nil {
		return nil, errors.New("wrong transaction body")
	}
	loginTx := &LoginTransactionBody{}
	copy(loginTx.PeerId[:], rpcBody.PeerIdBytes())
	return loginTx, nil
}

func translateRpcVoteBodyToBody(rpcBody *RpcVoteTransactionBody) (*VoteTransactionBody, error) {
	if rpcBody == nil {
		return nil, errors.New("wrong transaction body")
	}

	return &VoteTransactionBody{To: hasharry.StringToAddress(rpcBody.To)}, nil
}

func addressToString(address hasharry.Address) string {
	if address.IsEqual(hasharry.StringToAddress(CoinBase)) {
		return CoinBase
	}
	return address.String()
}

func addrListToHashAddr(addrList []string) []hasharry.Address {
	hashList := make([]hasharry.Address, len(addrList))
	for i, addr := range addrList {
		hashList[i] = hasharry.StringToAddress(addr)
	}
	return hashList
}

func hashAddrToAddr(hashList []hasharry.Address) []string {
	addrList := make([]string, len(hashList))
	for i, hash := range hashList {
		addrList[i] = hash.String()
	}
	return addrList
}
