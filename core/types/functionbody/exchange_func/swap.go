package exchange_func

import (
	"errors"
	"fmt"
	"github.com/jhdriver/UBChain/common/hasharry"
	"github.com/jhdriver/UBChain/param"
	"github.com/jhdriver/UBChain/ut"
)

type ExactIn struct {
	AmountIn     uint64
	AmountOutMin uint64
	Path         []hasharry.Address
	To           hasharry.Address
	Deadline     uint64
}

func (e *ExactIn) Verify() error {
	if !ut.CheckUBCAddress(param.Net, e.To.String()) {
		return errors.New("wrong to address")
	}
	for _, addr := range e.Path {
		if !ut.IsValidContractAddress(param.Net, addr.String()) {
			return fmt.Errorf("wrong path address %s", addr.String())
		}
	}
	return nil
}

type ExactOut struct {
	AmountOut   uint64
	AmountInMax uint64
	Path        []hasharry.Address
	To          hasharry.Address
	Deadline    uint64
}

func (e *ExactOut) Verify() error {
	if !ut.CheckUBCAddress(param.Net, e.To.String()) {
		return errors.New("wrong to address")
	}
	for _, addr := range e.Path {
		if !ut.IsValidContractAddress(param.Net, addr.String()) {
			return fmt.Errorf("wrong path address %s", addr.String())
		}
	}
	return nil
}
