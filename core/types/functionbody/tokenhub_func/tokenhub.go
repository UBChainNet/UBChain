package tokenhub_func

import (
	"errors"
	"fmt"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/param"
	"github.com/UBChainNet/UBChain/ut"
)

const maxFeeRate float64 = 1

type TokenHubInitBody struct {
	Setter hasharry.Address
	Admin  hasharry.Address
	FeeTo  hasharry.Address
	FeeRate float64
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
	if t.FeeRate > maxFeeRate{
		return fmt.Errorf("the fees rate cannot be greater than %f%%", maxFeeRate * 100)
	}
	return nil
}

type TokenHubAckBody struct {
	Sequences []uint64
	AckTypes  []uint8
}

func (t *TokenHubAckBody) Verify() error {
	if len(t.Sequences) != len(t.AckTypes){
		return errors.New("invalid ack data")
	}
	if len(t.Sequences) == 0{
		return errors.New("invalid ack data")
	}
	return nil
}
