package types

import (
	"fmt"
	"github.com/UBChainNet/UBChain/param"
	"testing"
)

func TestCalCoinBase(t *testing.T) {
	heights := []uint64{
		0, 1, 100, 1000, param.CoinHeight - 1, param.CoinHeight, 21588999, 21589000, 42612999, 42613000,
		63636999, 63637000, 84660999, 84661000, 105684999, 105685000, 105685001}

	for _, height := range heights {
		coinbase := CalCoinBase(height, param.CoinHeight)
		fmt.Println(height, " = ", coinbase)
	}
}
