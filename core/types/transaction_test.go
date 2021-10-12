package types

import (
	"fmt"
	"github.com/UBChainNet/UBChain/param"
	"testing"
)

func TestCalCoinBase(t *testing.T) {
	heights := []uint64{
		0, 1, 100, 1000, param.CoinHeight -1, param.CoinHeight, 21080499, 21080500, 42104499, 42104500, 63128499, 63128500, 84152499, 84152500, 105176499, 105176500, 105176501}

	for _, height := range heights{
		coinbase := CalCoinBase(height, param.CoinHeight)
		fmt.Println(height, " = ", coinbase)
	}
}
