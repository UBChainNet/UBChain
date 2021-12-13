package exchange

import (
	"fmt"
	"testing"
)

func TestGetRecordEveryHeight(t *testing.T) {
	Records := []*Record{
		&Record{
			Start:  9900,
			Amount: 3,
		},
		&Record{
			Start:  9980,
			Amount: 2,
		},
		&Record{
			Start:  10000,
			Amount: 1,
		},
	}
	rs := GetRecordEveryHeight(Records, 10010, 1, 10000)
	fmt.Println(rs)
}
