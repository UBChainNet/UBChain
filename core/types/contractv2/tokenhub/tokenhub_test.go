package tokenhub

import (
	"fmt"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"testing"
)

func TestTokenHub_AcrossFinished(t1 *testing.T) {
	th := NewTokenHub(hasharry.Address{}, hasharry.Address{}, hasharry.Address{}, hasharry.Address{}, "0.1")
	th.AcrossFinished(hasharry.Address{}, []uint64{1, 2, 3, 4, 5})
	fmt.Println(th.FinishSeq)
	fmt.Println(th.ContinueSeq)
	th.AcrossFinished(hasharry.Address{}, []uint64{1, 2, 3, 4, 6,  9, 12})
	fmt.Println(th.FinishSeq)
	fmt.Println(th.ContinueSeq)
	th.AcrossFinished(hasharry.Address{}, []uint64{7, 6,  9, 12, 15})
	fmt.Println(th.FinishSeq)
	fmt.Println(th.ContinueSeq)
}
