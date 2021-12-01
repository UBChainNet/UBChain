package tokenhub

import (
	"errors"
	"fmt"
	"github.com/UBChainNet/UBChain/common/encode/rlp"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/param"
	"github.com/ethereum/go-ethereum/common"
	"sort"
)

const MinTransferAmount = 100000000
const MaxTransferLength = 10000

type AckType uint8

const(
	Send AckType = 1
	Confirmed  = 2
	Failed = 3
)

type Transfer struct {
	Sequence uint64 `json:"sequence"`
	From string `json:"from"`
	To   string `json:"to"`
	Amount uint64 `json:"amount"`
	Fees  uint64 `json:"fees"`
}

type TokenHub struct {
	Address hasharry.Address
	Setter hasharry.Address
	Admin hasharry.Address
	FeeTo hasharry.Address
	FeeRate  float64
	Transfers map[uint64]*Transfer
	Unconfirmed map[uint64]*Transfer
	Sequence uint64
}

func NewTokenHub(address, setter, admin, feeTo hasharry.Address, feeRate float64)*TokenHub{
	return &TokenHub{
		Address: address,
		Setter:     setter,
		Admin:        admin,
		FeeTo:        feeTo,
		FeeRate:      feeRate,
		Transfers: make(map[uint64]*Transfer),
		Unconfirmed:  make(map[uint64]*Transfer),
	}
}

func (t *TokenHub)SetSetter(from, setter hasharry.Address) error {
	if !from.IsEqual(t.Setter){
		return errors.New("forbidden")
	}
	t.Setter = setter
	return nil
}

func (t *TokenHub)SetAdmin(from, admin hasharry.Address) error {
	if !from.IsEqual(t.Setter){
		return errors.New("forbidden")
	}
	t.Setter = admin
	return nil
}

func (t *TokenHub)SetFeeTo(from, feeTo hasharry.Address) error {
	if !from.IsEqual(t.Setter) {
		return errors.New("forbidden")
	}
	t.FeeTo = feeTo
	return nil
}

func (t *TokenHub)SetFeeRate(from hasharry.Address, rate float64) error {
	if !from.IsEqual(t.Setter) {
		return errors.New("forbidden")
	}
	t.FeeRate = rate
	return nil
}

func (t *TokenHub)Transfer(from hasharry.Address, to string, amount uint64) ([]*TransferEvent, error){
	if !from.IsEqual(t.Admin){
		return nil, errors.New("forbidden")
	}
	if amount < MinTransferAmount{
		return nil, fmt.Errorf("the minimum allowable transfer amount is %d", MinTransferAmount)
	}
	if !common.IsHexAddress(to){
		return nil, errors.New("incorrect to address")
	}
	if len(t.Transfers) >= MaxTransferLength{
		return nil, errors.New("too many transfers, please wait")
	}
	fees := uint64(float64(amount) * t.FeeRate)
	if fees >= amount{
		return nil, errors.New("the transfer fee is insufficient")
	}

	t.Sequence++
	t.Transfers[t.Sequence] = &Transfer{
		Sequence: t.Sequence,
		From:   from.String(),
		To:     to,
		Amount: amount - fees,
		Fees:   fees,
	}
	var events []*TransferEvent
	events = append(events, &TransferEvent{
		From:   from,
		To:     t.Address,
		Token:  param.Token,
		Amount: amount,
	})
	return events, nil
}

func (t *TokenHub)AckTransfer(from hasharry.Address, ackData map[uint64]AckType)([]*TransferEvent, error) {
	if !from.IsEqual(t.Admin){
		return nil, errors.New("forbidden")
	}
	if len(ackData) == 0{
		return nil, errors.New("invalid ack data")
	}
	var events  = []*TransferEvent{}
	for sequence, ack := range ackData{
		switch ack {
		case Send:
			transfer, exist := t.Transfers[sequence]
			if exist{
				delete(t.Transfers, sequence)
			}else{
				return nil, fmt.Errorf("ack transafer sequence %d does not exist", sequence)
			}
			t.Unconfirmed[sequence] = transfer
		case Confirmed:
			transfer, exist := t.Unconfirmed[sequence]
			if exist{
				delete(t.Unconfirmed, sequence)
			}else{
				return nil, fmt.Errorf("ack unconfirmed sequence %d does not exist", sequence)
			}
			events = append(events, &TransferEvent{
				From:   t.Address,
				To:     t.FeeTo,
				Token:  param.Token,
				Amount: transfer.Fees,
			})
		case Failed:
			transfer, exist := t.Transfers[sequence]
			if exist{
				delete(t.Transfers, sequence)
			}else{
				transfer, exist = t.Unconfirmed[sequence]
				if exist{
					delete(t.Unconfirmed, sequence)
				}else{
					return nil, fmt.Errorf("ack failed transfer sequence %d does not exist", sequence)
				}
			}
			events = append(events, &TransferEvent{
				From:   t.Address,
				To:     hasharry.StringToAddress(transfer.From),
				Token:  param.Token,
				Amount: transfer.Amount + transfer.Fees,
			})
		default:
			return nil, fmt.Errorf("invalid ack type %d", ack)
		}
	}
	return events, nil
}

func (t *TokenHub) ToRlp() *RlpTokenHub {
	rlpTh := &RlpTokenHub{
		Address:     t.Address.String(),
		Setter:      t.Setter.String(),
		Admin:       t.Admin.String(),
		FeeTo:       t.FeeTo.String(),
		FeeRate:     t.FeeRate,
		Transfers:   make([]*Transfer, 0),
		Unconfirmed: make([]*Transfer, 0),
		Sequence:    0,
	}
	for _, transfer := range t.Transfers {
		rlpTh.Transfers = append(rlpTh.Transfers, transfer)
	}
	for _, transfer := range t.Unconfirmed {
		rlpTh.Unconfirmed = append(rlpTh.Unconfirmed, transfer)
	}

	sort.Slice(rlpTh.Transfers, func(i, j int) bool {
		return rlpTh.Transfers[i].Sequence < rlpTh.Transfers[j].Sequence
	})

	sort.Slice(rlpTh.Unconfirmed, func(i, j int) bool {
		return rlpTh.Unconfirmed[i].Sequence < rlpTh.Unconfirmed[j].Sequence
	})
	return rlpTh
}

func (t *TokenHub)Bytes() []byte {
	bytes, _ := rlp.EncodeToBytes(t.ToRlp())
	return bytes
}

type TransferEvent struct {
	From      hasharry.Address
	To        hasharry.Address
	Token     hasharry.Address
	Amount    uint64
}


type RlpTokenHub struct {
	Address string
	Setter string
	Admin string
	FeeTo string
	FeeRate  float64
	Transfers []*Transfer
	Unconfirmed []*Transfer
	Sequence uint64
}

func DecodeToTokenHub(bytes []byte) (*TokenHub, error) {
	var rlpTh *RlpTokenHub
	rlp.DecodeBytes(bytes, &rlpTh)
	th := &TokenHub{
		Address:     hasharry.StringToAddress(rlpTh.Address),
		Setter:      hasharry.StringToAddress(rlpTh.Setter),
		Admin:       hasharry.StringToAddress(rlpTh.Admin),
		FeeTo:       hasharry.StringToAddress(rlpTh.FeeTo),
		FeeRate:     rlpTh.FeeRate,
		Transfers:   make(map[uint64]*Transfer),
		Unconfirmed: make(map[uint64]*Transfer),
		Sequence:    rlpTh.Sequence,
	}

	for _, transfer := range rlpTh.Transfers {
		th.Transfers[transfer.Sequence] = transfer
	}
	for _, transfer := range rlpTh.Unconfirmed {
		th.Unconfirmed[transfer.Sequence] = transfer
	}
	return th, nil
}