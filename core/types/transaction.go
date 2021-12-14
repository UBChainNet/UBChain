package types

import (
	"encoding/json"
	"fmt"
	"github.com/UBChainNet/UBChain/common/encode/rlp"
	hash2 "github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/core/types/contractv2"
	"github.com/UBChainNet/UBChain/core/types/functionbody/exchange_func"
	"github.com/UBChainNet/UBChain/core/types/functionbody/tokenhub_func"
	"github.com/UBChainNet/UBChain/crypto/ecc/secp256k1"
	"github.com/UBChainNet/UBChain/crypto/hash"
	"github.com/UBChainNet/UBChain/param"
	"github.com/UBChainNet/UBChain/ut"
)

const (
	Transfer_ TransactionType = iota
	Contract_
	ContractV2_
	LoginCandidate_

	/*LogoutCandidate
	VoteToCandidate*/
)
const MaxNote = 256

const CoinBase = "coinbase"

type TransactionType uint8

type TransactionHead struct {
	TxHash     hash2.Hash
	TxType     TransactionType
	From       hash2.Address
	Nonce      uint64
	Fees       uint64
	Time       uint64
	Note       string
	SignScript *SignScript
}

type Transaction struct {
	TxHead *TransactionHead
	TxBody ITransactionBody
}

func (t *Transaction) IsCoinBase() bool {
	return t.TxHead.From.IsEqual(hash2.StringToAddress(CoinBase))
}

func (t *Transaction) Size() uint64 {
	bytes, _ := t.EncodeToBytes()
	return uint64(len(bytes))
}

func (t *Transaction) VerifyTx(height uint64) error {
	if err := t.verifyHead(height); err != nil {
		return err
	}

	if err := t.verifyBody(); err != nil {
		return err
	}
	return nil
}

func (t *Transaction) verifyHead(height uint64) error {
	if t.TxHead == nil {
		return ErrTxHead
	}

	if err := t.verifyTxType(); err != nil {
		return err
	}

	if err := t.verifyTxHash(); err != nil {
		return err
	}

	if err := t.verifyTxFrom(height); err != nil {
		return err
	}

	if err := t.verifyTxNote(); err != nil {
		return err
	}

	if err := t.verifyTxFees(); err != nil {
		return err
	}

	if err := t.verifyTxSinger(); err != nil {
		return err
	}
	return nil
}

func (t *Transaction) verifyBody() error {
	if t.TxBody == nil {
		return ErrTxBody
	}

	if err := t.verifyAmount(); err != nil {
		return err
	}

	if err := t.TxBody.VerifyBody(t.TxHead.From); err != nil {
		return err
	}
	return nil
}

func (t *Transaction) VerifyCoinBaseTx(height, sumFees uint64, miner string) error {
	if err := t.verifyTxSize(); err != nil {
		return err
	}

	if err := t.verifyCoinBaseAddress(miner); err != nil {
		return err
	}

	if err := t.verifyCoinBaseAmount(height, sumFees); err != nil {
		return err
	}

	return nil
}

func (t *Transaction) verifyTxFees() error {
	var fees uint64
	switch t.TxHead.TxType {
	case Transfer_:
		fees = TransferFees(len(t.TxBody.ToAddress().ReceiverList()))
	case Contract_:
		fees = param.TokenConsumption
	case ContractV2_:
		fees = param.Fees
	}
	if t.TxHead.Fees != fees {
		return fmt.Errorf("transaction costs %d fees", fees)
	}
	return nil
}

func (t *Transaction) verifyTxSinger() error {
	if !Verify(t.TxHead.TxHash, t.TxHead.SignScript) {
		return ErrSignature
	}

	if !VerifySigner(param.Net, t.TxHead.From, t.TxHead.SignScript.PubKey) {
		return ErrSigner
	}
	return nil
}

func (t *Transaction) verifyTxSize() error {
	// TODO change maxsize
	switch t.TxHead.TxType {
	case Transfer_:
		fallthrough
	case Contract_:
		return nil
		/*case LogoutCandidate:
			fallthrough
		case LoginCandidate_:
			fallthrough
		case VoteToCandidate:
			if t.Size() > MaxNoDataTxSize {
				return ErrTxSize
			}*/
	}
	return nil
}

func (t *Transaction) verifyCoinBaseAmount(height, amount uint64) error {
	nTx := t.TxBody.(*TransferBody)
	sumAmount := CalCoinBase(height, param.CoinHeight) + amount
	if sumAmount != nTx.GetAmount() {
		return ErrCoinBase
	}
	return nil
}

func (t *Transaction) verifyCoinBaseAddress(miner string) error {
	nTx := t.TxBody.(*TransferBody)
	for _, receiver := range nTx.Receivers.ReceiverList() {
		if !ut.CheckUBCAddress(param.Net, receiver.Address.String()) {
			return fmt.Errorf("invalid coinbase address")
		}
		if receiver.Amount > 0 {
			if receiver.Address.String() != param.MinerReward[miner] {
				return fmt.Errorf("incorrect coinbase address")
			}
		}
	}
	return nil
}

func (t *Transaction) verifyAmount() error {
	for _, to := range t.TxBody.ToAddress().ReceiverList() {
		if to.Amount < param.MinAllowedAmount {
			return fmt.Errorf("the minimum amount of the transaction must not be less than %d", param.MinAllowedAmount)
		}
	}
	return nil
}

func (t *Transaction) verifyTxFrom(height uint64) error {
	if !ut.CheckUBCAddress(param.Net, t.From().String()) {
		return ErrAddress
	}
	if height >= param.UIPBlock5 {
		_, exist := param.Blacklist[t.TxHead.From.String()]
		if exist {
			return ErrAddress
		}
	}
	return nil
}

func (t *Transaction) verifyTxType() error {
	switch t.TxHead.TxType {
	case Transfer_:
		return nil
	case Contract_:
		return nil
	case ContractV2_:
		return nil
		/*
			case VoteToCandidate:
				return nil
			case LoginCandidate_:
				return nil
			case LogoutCandidate:
				return nil*/
	}
	return ErrTxType
}

func (t *Transaction) verifyTxHash() error {
	newTx := t.copy()
	newTx.SetHash()
	if newTx.Hash().IsEqual(t.Hash()) {
		return nil
	}
	return ErrTxHash
}

func (t *Transaction) verifyTxNote() error {
	if len(t.TxHead.Note) > MaxNote {
		return fmt.Errorf("the length of the transaction note must not be greater than %d", MaxNote)
	}
	return nil
}

func (t *Transaction) EncodeToBytes() ([]byte, error) {
	return rlp.EncodeToBytes(t)
}

func (t *Transaction) SignTx(key *secp256k1.PrivateKey) error {
	var err error
	if t.TxHead.SignScript, err = Sign(key, t.TxHead.TxHash); err != nil {
		return err
	}
	return nil
}

func (t *Transaction) SetHash() error {
	t.TxHead.TxHash = hash2.Hash{}
	t.TxHead.SignScript = &SignScript{}
	rpcTx, err := TranslateTxToRpcTx(t)
	if err != nil {
		return err
	}
	mBytes, err := json.Marshal(rpcTx)
	if err != nil {
		return err
	}
	t.TxHead.TxHash = hash.Hash(mBytes)
	return nil
}

func (t *Transaction) copy() *Transaction {
	header := &TransactionHead{
		TxHash:     t.TxHead.TxHash,
		TxType:     t.TxHead.TxType,
		From:       t.TxHead.From,
		Nonce:      t.TxHead.Nonce,
		Fees:       t.TxHead.Fees,
		Time:       t.TxHead.Time,
		Note:       t.TxHead.Note,
		SignScript: t.TxHead.SignScript,
	}
	return &Transaction{
		TxHead: header,
		TxBody: t.TxBody,
	}
}

func (t *Transaction) Hash() hash2.Hash {
	return t.TxHead.TxHash
}

func (t *Transaction) From() hash2.Address {
	return t.TxHead.From
}

func (t *Transaction) GetFees() uint64 {
	return t.TxHead.Fees
}

func (t *Transaction) GetNonce() uint64 {
	return t.TxHead.Nonce
}

func (t *Transaction) GetTime() uint64 {
	return t.TxHead.Time
}

func (t *Transaction) GetNote() string {
	return t.TxHead.Note
}

func (t *Transaction) GetTxType() TransactionType {
	return t.TxHead.TxType
}

func (t *Transaction) GetSignScript() *SignScript {
	return t.TxHead.SignScript
}

func (t *Transaction) GetTxHead() *TransactionHead {
	return t.TxHead
}

func (t *Transaction) GetTxBody() ITransactionBody {
	return t.TxBody
}

func (t *Transaction) TranslateToRlpTransaction() *RlpTransaction {
	rlpTx := &RlpTransaction{}
	rlpTx.TxHead = t.TxHead
	switch t.GetTxType() {
	case ContractV2_:
		body, _ := t.TxBody.(*TxContractV2Body)
		rlpC := &RlpContract{
			TxHead: t.TxHead,
			TxBody: RlpContractBody{
				Contract:     body.Contract,
				Type:         body.Type,
				FunctionType: body.FunctionType,
				Function:     nil,
			},
		}
		switch body.FunctionType {
		case contractv2.Exchange_Init:
			function, _ := body.Function.(*exchange_func.ExchangeInitBody)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		case contractv2.Exchange_SetAdmin:
			function, _ := body.Function.(*exchange_func.ExchangeAdmin)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		case contractv2.Exchange_SetFeeTo:
			function, _ := body.Function.(*exchange_func.ExchangeFeeTo)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		case contractv2.Exchange_ExactIn:
			function, _ := body.Function.(*exchange_func.ExactIn)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		case contractv2.Exchange_ExactOut:
			function, _ := body.Function.(*exchange_func.ExactOut)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		case contractv2.Pair_AddLiquidity:
			function, _ := body.Function.(*exchange_func.ExchangeAddLiquidity)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		case contractv2.Pair_RemoveLiquidity:
			function, _ := body.Function.(*exchange_func.ExchangeRemoveLiquidity)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		case contractv2.Pledge_Init:
			function, _ := body.Function.(*exchange_func.PledgeInitBody)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		case contractv2.Pledge_Start:
			function, _ := body.Function.(*exchange_func.PledgeStartBody)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		case contractv2.Pledge_AddPool:
			function, _ := body.Function.(*exchange_func.PledgeAddPoolBody)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		case contractv2.Pledge_RemovePool:
			function, _ := body.Function.(*exchange_func.PledgeRemovePoolBody)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		case contractv2.Pledge_Add:
			function, _ := body.Function.(*exchange_func.PledgeAddBody)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		case contractv2.Pledge_Remove:
			function, _ := body.Function.(*exchange_func.PledgeRemoveBody)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		case contractv2.Pledge_RemoveReward:
			function, _ := body.Function.(*exchange_func.PledgeRewardRemoveBody)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		case contractv2.Pledge_Update:
			function, _ := body.Function.(*exchange_func.PledgeUpdateBody)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		case contractv2.TokenHub_init:
			function, _ := body.Function.(*tokenhub_func.TokenHubInitBody)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		case contractv2.TokenHub_Ack:
			function, _ := body.Function.(*tokenhub_func.TokenHubAckBody)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		case contractv2.TokenHub_TransferOut:
			function, _ := body.Function.(*tokenhub_func.TokenHubTransferOutBody)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		case contractv2.TokenHub_TransferIn:
			function, _ := body.Function.(*tokenhub_func.TokenHubTransferInBody)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		case contractv2.TokenHub_FinishAcross:
			function, _ := body.Function.(*tokenhub_func.TokenHubFinishAcrossBody)
			bytes, _ := rlp.EncodeToBytes(function)
			rlpC.TxBody.Function = bytes
		}
		rlpTx.TxBody, _ = rlp.EncodeToBytes(rlpC.TxBody)
	default:
		rlpTx.TxBody, _ = rlp.EncodeToBytes(t.TxBody)
	}
	return rlpTx
}

type TxLocation struct {
	TxRoot  hash2.Hash
	TxIndex uint32
	Height  uint64
}

func (t *TxLocation) GetHeight() uint64 {
	return t.Height
}

func CalCoinBase(height, startHeight uint64) uint64 {
	if height < startHeight {
		return 0
	}

	count := height - startHeight
	heightRange := count / ((uint64(3600*24) / param.BlockInterval) * 365 * 10)
	switch heightRange {
	case 0:
		return 451864536
	case 1:
		return 677796804
	case 2:
		return 903729072
	case 3:
		return 1129661339
	case 4:
		return 1355593607
	default:
		return 0
	}

}

func TransferFees(receiverCount int) uint64 {
	return param.Fees * uint64(receiverCount)
}
