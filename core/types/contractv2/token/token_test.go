package token

import (
	"fmt"
	"testing"
)

func TestToken_Mint(t1 *testing.T) {

	token, err := NewToken("1", "2", 10000000, 1000000, 10, 3, 5)
	if err != nil{
		t1.Fatal(err.Error())
	}
	fmt.Println(token.Mint("1", 9996))
	fmt.Println(token.Mint("1", 9997))
	fmt.Println(token.Mint("1", 9998))
	fmt.Println(token.Mint("1", 9999))
	fmt.Println(token.Mint("1", 10000))
	fmt.Println(token.Mint("1", 1000000))
	fmt.Println(token.Mint("1", 1000000))
	fmt.Println(token.Mint("1", 1000000))
	fmt.Println(token.Mint("1", 1000000000000))
	fmt.Println(token.Mint("1", 10))
	fmt.Println(token.Mint("1", 10))
	fmt.Println(token.Mint("1", 10))
	fmt.Println(token.Mint("1", 10))
	fmt.Println(token.Mint("1", 10))
	fmt.Println(token.Mint("1", 10))
	fmt.Println(token.Mint("1", 10))
}
