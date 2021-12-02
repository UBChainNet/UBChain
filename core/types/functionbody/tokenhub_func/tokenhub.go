package tokenhub_func

import (
	"errors"
	"fmt"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/param"
	"github.com/UBChainNet/UBChain/ut"
	"github.com/ethereum/go-ethereum/common"
	"strconv"
)

const maxFeeRate float64 = 1
const minTransferAmount uint64 = 100000000

type TokenHubInitBody struct {
	Setter  hasharry.Address
	Admin   hasharry.Address
	FeeTo   hasharry.Address
	FeeRate string
}

func (t *TokenHubInitBody) Verify() error {
	if ok := ut.CheckUBCAddress(param.Net, t.Setter.String()); !ok {
		return errors.New("wrong setter address")
	}
	if ok := ut.CheckUBCAddress(param.Net, t.Admin.String()); !ok {
		return errors.New("wrong admin address")
	}
	if ok := ut.CheckUBCAddress(param.Net, t.FeeTo.String()); !ok {
		return errors.New("wrong feeTo address")
	}
	feeRate, err := strconv.ParseFloat(t.FeeRate, 64)
	if err != nil {
		return err
	}
	if feeRate >= maxFeeRate {
		return fmt.Errorf("the fees rate cannot be greater than %f%%", maxFeeRate*100)
	}
	return nil
}

type TokenHubTransferOutBody struct {
	To     string
	Amount uint64
}

func (t *TokenHubTransferOutBody) Verify() error {
	if !common.IsHexAddress(t.To) {
		return errors.New("incorrect to address")
	}
	if t.Amount < minTransferAmount {
		return fmt.Errorf("the minimum allowable transfer amount is %d", minTransferAmount)
	}
	return nil
}

type TokenHubTransferInBody struct {
	To        hasharry.Address
	Amount    uint64
	AcrossSeq uint64
}

func (t *TokenHubTransferInBody) Verify() error {
	if !ut.CheckUBCAddress(param.Net, t.To.String()) {
		return errors.New("incorrect to address")
	}
	return nil
}

type TokenHubFinishAcrossBody struct {
	AcrossSeqs []uint64
}

func (t *TokenHubFinishAcrossBody) Verify() error {
	if len(t.AcrossSeqs) == 0 {
		return errors.New("invalid across seq")
	}
	return nil
}

type TokenHubAckBody struct {
	Sequences []uint64
	AckTypes  []uint8
	Hashes 	  []string
}

func (t *TokenHubAckBody) Verify() error {
	if len(t.Sequences) != len(t.AckTypes) {
		return errors.New("invalid ack data")
	}
	if len(t.Sequences) != len(t.Hashes) {
		return errors.New("invalid ack data")
	}
	if len(t.Sequences) == 0 {
		return errors.New("invalid ack data")
	}
	for _, hash := range t.Hashes{
		if len(hash) > 100{
			return errors.New("invalid ack hash")
		}
	}
	return nil
}
