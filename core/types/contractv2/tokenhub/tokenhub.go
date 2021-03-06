package tokenhub

import (
	"errors"
	"fmt"
	"github.com/UBChainNet/UBChain/common/encode/rlp"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/param"
	"github.com/ethereum/go-ethereum/common"
	"sort"
	"strconv"
)

const MinTransferAmount = 100000000
const MaxTransferLength = 10000

type AckType uint8

const (
	Send      AckType = 1
	Confirmed AckType = 2
	Failed    AckType = 3
	Error     AckType = 4
)

type Transfer struct {
	Sequence uint64 `json:"sequence"`
	From     string `json:"from"`
	To       string `json:"to"`
	Amount   uint64 `json:"amount"`
	Fees     uint64 `json:"fees"`
	Hash     string `json:"hash"`
}

type TokenHub struct {
	Address         hasharry.Address
	Setter          hasharry.Address
	Admin           hasharry.Address
	FeeTo           hasharry.Address
	FeeRate         string
	TransferOuts    map[uint64]*Transfer
	UnconfirmedOuts map[uint64]*Transfer
	AcrossSeqs      map[uint64]string
	FinishSeq       map[uint64]bool
	ContinueSeq     uint64
	Sequence        uint64
	OutAmount       uint64
	InAmount        uint64
}

func NewTokenHub(address, setter, admin, feeTo hasharry.Address, feeRate string) *TokenHub {
	return &TokenHub{
		Address:         address,
		Setter:          setter,
		Admin:           admin,
		FeeTo:           feeTo,
		FeeRate:         feeRate,
		TransferOuts:    make(map[uint64]*Transfer),
		UnconfirmedOuts: make(map[uint64]*Transfer),
		AcrossSeqs:      make(map[uint64]string),
		FinishSeq:       make(map[uint64]bool),
	}
}

func (t *TokenHub) SetSetter(from, setter hasharry.Address) error {
	if !from.IsEqual(t.Setter) {
		return errors.New("forbidden")
	}
	t.Setter = setter
	return nil
}

func (t *TokenHub) SetAdmin(from, admin hasharry.Address) error {
	if !from.IsEqual(t.Setter) {
		return errors.New("forbidden")
	}
	t.Admin = admin
	return nil
}

func (t *TokenHub) SetFeeTo(from, feeTo hasharry.Address) error {
	if !from.IsEqual(t.Setter) {
		return errors.New("forbidden")
	}
	t.FeeTo = feeTo
	return nil
}

func (t *TokenHub) SetFeeRate(from hasharry.Address, rate string) error {
	if !from.IsEqual(t.Setter) {
		return errors.New("forbidden")
	}
	t.FeeRate = rate
	return nil
}

func (t *TokenHub) TransferIn(from hasharry.Address, to hasharry.Address, amount, acrossSeq uint64, msgHash string) ([]*TransferEvent, error) {
	if !t.Admin.IsEqual(from) {
		return nil, errors.New("forbidden")
	}
	if _, exist := t.AcrossSeqs[acrossSeq]; exist {
		return nil, fmt.Errorf("%d across transfer has send", acrossSeq)
	}
	feeRate, _ := strconv.ParseFloat(t.FeeRate, 64)
	fees := uint64(float64(amount) * feeRate)
	if fees >= amount {
		return nil, errors.New("the transfer fee is insufficient")
	}

	var events []*TransferEvent
	events = append(events, &TransferEvent{
		From:   t.Address,
		To:     to,
		Token:  param.Token,
		Amount: amount - fees,
	})
	if fees != 0 {
		events = append(events, &TransferEvent{
			From:   t.Address,
			To:     t.FeeTo,
			Token:  param.Token,
			Amount: fees,
		})
	}
	t.AcrossSeqs[acrossSeq] = msgHash
	t.InAmount += amount
	return events, nil
}

func (t *TokenHub) AcrossFinished(from hasharry.Address, acrossSeq []uint64) error {
	if !t.Admin.IsEqual(from) {
		return errors.New("forbidden")
	}
	for _, seq := range acrossSeq {
		t.FinishSeq[seq] = true
		delete(t.AcrossSeqs, seq)
	}
	tmp := make([]uint64, 0)
	for seq := range t.FinishSeq {
		tmp = append(tmp, seq)
	}
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i] < tmp[j]
	})
	for _, seq := range tmp {
		if t.ContinueSeq >= seq {
			delete(t.FinishSeq, seq)
		} else if t.ContinueSeq+1 == seq {
			t.ContinueSeq++
			delete(t.FinishSeq, seq)
		} else {
			break
		}
	}
	return nil
}

func (t *TokenHub) TransferOut(from hasharry.Address, to string, amount uint64) ([]*TransferEvent, error) {
	if !common.IsHexAddress(to) {
		return nil, errors.New("incorrect to address")
	}
	if len(t.TransferOuts) >= MaxTransferLength {
		return nil, errors.New("too many transfers, please wait")
	}
	feeRate, _ := strconv.ParseFloat(t.FeeRate, 64)
	fees := uint64(float64(amount) * feeRate)
	if fees >= amount {
		return nil, errors.New("the transfer fee is insufficient")
	}

	t.Sequence++
	t.TransferOuts[t.Sequence] = &Transfer{
		Sequence: t.Sequence,
		From:     from.String(),
		To:       to,
		Amount:   amount - fees,
		Fees:     fees,
	}
	var events []*TransferEvent
	events = append(events, &TransferEvent{
		From:   from,
		To:     t.Address,
		Token:  param.Token,
		Amount: amount,
	})
	t.OutAmount += amount - fees
	return events, nil
}

func (t *TokenHub) AckTransfer(from hasharry.Address, ackData map[uint64]AckType, ackHash map[uint64]string) ([]*TransferEvent, error) {
	if !from.IsEqual(t.Admin) {
		return nil, errors.New("forbidden")
	}
	if len(ackData) == 0 {
		return nil, errors.New("invalid ack data")
	}
	var events = []*TransferEvent{}
	for sequence, ack := range ackData {
		switch ack {
		case Send:
			transfer, exist := t.TransferOuts[sequence]
			if exist {
				delete(t.TransferOuts, sequence)
			} else {
				return nil, fmt.Errorf("ack transafer sequence %d does not exist", sequence)
			}
			transfer.Hash = ackHash[sequence]
			t.UnconfirmedOuts[sequence] = transfer
		case Confirmed:
			transfer, exist := t.UnconfirmedOuts[sequence]
			if exist {
				delete(t.UnconfirmedOuts, sequence)
			} else {
				transfer, exist = t.TransferOuts[sequence]
				if exist {
					delete(t.TransferOuts, sequence)
				} else {
					return nil, fmt.Errorf("ack unconfirmed sequence %d does not exist", sequence)
				}
			}
			events = append(events, &TransferEvent{
				From:   t.Address,
				To:     t.FeeTo,
				Token:  param.Token,
				Amount: transfer.Fees,
			})
		case Failed:
			transfer, exist := t.UnconfirmedOuts[sequence]
			if exist {
				transfer.Hash = ""
				t.TransferOuts[sequence] = transfer
			}
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
		AcrossSeqs:  make([]*Across, 0),
		FinishSeq:   make([]uint64, 0),
		ContinueSeq: t.ContinueSeq,
		Sequence:    t.Sequence,
		InAmount:    t.InAmount,
		OutAmount:   t.OutAmount,
	}
	for _, transfer := range t.TransferOuts {
		rlpTh.Transfers = append(rlpTh.Transfers, transfer)
	}
	for _, transfer := range t.UnconfirmedOuts {
		rlpTh.Unconfirmed = append(rlpTh.Unconfirmed, transfer)
	}
	for seq, hash := range t.AcrossSeqs {
		rlpTh.AcrossSeqs = append(rlpTh.AcrossSeqs, &Across{
			Seq:  seq,
			Hash: hash,
		})
	}
	for seq, _ := range t.FinishSeq {
		rlpTh.FinishSeq = append(rlpTh.FinishSeq, seq)
	}

	sort.Slice(rlpTh.Transfers, func(i, j int) bool {
		return rlpTh.Transfers[i].Sequence < rlpTh.Transfers[j].Sequence
	})

	sort.Slice(rlpTh.Unconfirmed, func(i, j int) bool {
		return rlpTh.Unconfirmed[i].Sequence < rlpTh.Unconfirmed[j].Sequence
	})

	sort.Slice(rlpTh.AcrossSeqs, func(i, j int) bool {
		return rlpTh.AcrossSeqs[i].Seq < rlpTh.AcrossSeqs[j].Seq
	})

	sort.Slice(rlpTh.FinishSeq, func(i, j int) bool {
		return rlpTh.FinishSeq[i] < rlpTh.FinishSeq[j]
	})
	return rlpTh
}

func (t *TokenHub) Bytes() []byte {
	bytes, _ := rlp.EncodeToBytes(t.ToRlp())
	return bytes
}

func (t *TokenHub) GetSymbol() string {
	return ""
}

type TransferEvent struct {
	From   hasharry.Address
	To     hasharry.Address
	Token  hasharry.Address
	Amount uint64
}

type Across struct {
	Seq  uint64
	Hash string
}

type RlpTokenHub struct {
	Address     string
	Setter      string
	Admin       string
	FeeTo       string
	FeeRate     string
	Transfers   []*Transfer
	Unconfirmed []*Transfer
	AcrossSeqs  []*Across
	FinishSeq   []uint64
	ContinueSeq uint64
	Sequence    uint64
	InAmount    uint64
	OutAmount   uint64
}

func DecodeToTokenHub(bytes []byte) (*TokenHub, error) {
	var rlpTh *RlpTokenHub
	rlp.DecodeBytes(bytes, &rlpTh)
	th := &TokenHub{
		Address:         hasharry.StringToAddress(rlpTh.Address),
		Setter:          hasharry.StringToAddress(rlpTh.Setter),
		Admin:           hasharry.StringToAddress(rlpTh.Admin),
		FeeTo:           hasharry.StringToAddress(rlpTh.FeeTo),
		FeeRate:         rlpTh.FeeRate,
		TransferOuts:    make(map[uint64]*Transfer),
		UnconfirmedOuts: make(map[uint64]*Transfer),
		AcrossSeqs:      make(map[uint64]string),
		FinishSeq:       make(map[uint64]bool),
		ContinueSeq:     rlpTh.ContinueSeq,
		Sequence:        rlpTh.Sequence,
		InAmount:        rlpTh.InAmount,
		OutAmount:       rlpTh.OutAmount,
	}

	for _, transfer := range rlpTh.Transfers {
		th.TransferOuts[transfer.Sequence] = transfer
	}
	for _, transfer := range rlpTh.Unconfirmed {
		th.UnconfirmedOuts[transfer.Sequence] = transfer
	}
	for _, across := range rlpTh.AcrossSeqs {
		th.AcrossSeqs[across.Seq] = across.Hash
	}
	for _, seq := range rlpTh.FinishSeq {
		th.FinishSeq[seq] = true
	}
	return th, nil
}
