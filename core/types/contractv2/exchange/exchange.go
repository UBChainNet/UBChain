package exchange

import (
	"errors"
	"fmt"
	"github.com/UBChainNet/UBChain/common/encode/rlp"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/common/utils"
	"github.com/UBChainNet/UBChain/ut"
	"math"
	"strings"
)

type ReadFunction string

const (
	Func_PairAddress ReadFunction = "pairAddress"
	Func_PairList                 = "pairAddress"

	maxPairTokenLen = 3
)

type PairAddress struct {
	Key     string
	Address hasharry.Address
	Symbol0 string
	Symbol1 string
}

type RlpExchange struct {
	Symbol   string
	FeeTo    hasharry.Address
	Admin    hasharry.Address
	AllPairs []PairAddress
}

type Exchange struct {
	FeeTo    hasharry.Address
	Symbol   string
	Admin    hasharry.Address
	Pair     map[hasharry.Address]map[hasharry.Address]hasharry.Address
	AllPairs []PairAddress
}

func NewExchange(admin, feeTo hasharry.Address, symbol string) (*Exchange, error) {
	if err := ut.CheckSymbol(symbol); err != nil {
		return nil, err
	}
	return &Exchange{
		FeeTo:    feeTo,
		Admin:    admin,
		Symbol:   symbol,
		Pair:     make(map[hasharry.Address]map[hasharry.Address]hasharry.Address),
		AllPairs: make([]PairAddress, 0),
	}, nil
}

func (e *Exchange) GetSymbol() string {
	return e.Symbol
}

func (e *Exchange) SetFeeTo(address hasharry.Address, sender hasharry.Address) error {
	if err := e.VerifySetter(sender); err != nil {
		return err
	}
	e.FeeTo = address
	return nil
}

func (e *Exchange) SetAdmin(address hasharry.Address, sender hasharry.Address) error {
	if err := e.VerifySetter(sender); err != nil {
		return err
	}
	e.Admin = address
	return nil
}

func (e *Exchange) VerifySetter(sender hasharry.Address) error {
	if !e.Admin.IsEqual(sender) {
		return errors.New("forbidden")
	}
	return nil
}

func (e *Exchange) Exist(token0, token1 hasharry.Address) bool {
	token1Map, ok := e.Pair[token0]
	if ok {
		_, ok := token1Map[token1]
		return ok
	}
	return false
}

func (e *Exchange) PairExist(pairAddress hasharry.Address) bool {
	for _, pair := range e.AllPairs {
		if pair.Address.IsEqual(pairAddress) {
			return true
		}
	}
	return false
}

func (e *Exchange) PairAddress(token0, token1 hasharry.Address) hasharry.Address {
	token1Map, ok := e.Pair[token0]
	if ok {
		address, _ := token1Map[token1]
		return address
	}
	return hasharry.Address{}
}

func (e *Exchange) AddPair(token0, token1, address hasharry.Address, symbol0, symbol1 string) {
	token1Addr, ok := e.Pair[token0]
	if ok {
		token1Addr[token1] = address
		e.Pair[token0] = token1Addr
	} else {
		e.Pair[token0] = map[hasharry.Address]hasharry.Address{token1: address}
	}

	e.AllPairs = append(e.AllPairs, PairAddress{
		Key:     pairKey(token0, token1),
		Address: address,
		Symbol0: symbol0,
		Symbol1: symbol1,
	})
}

func (e *Exchange) Bytes() []byte {
	elpEx := &RlpExchange{
		FeeTo:    e.FeeTo,
		Admin:    e.Admin,
		Symbol:   e.Symbol,
		AllPairs: e.AllPairs,
	}
	bytes, _ := rlp.EncodeToBytes(elpEx)
	return bytes
}

func DecodeToExchange(bytes []byte) (*Exchange, error) {
	var rlpEx *RlpExchange
	if err := rlp.DecodeBytes(bytes, &rlpEx); err != nil {
		return nil, err
	}
	ex, err := NewExchange(rlpEx.Admin, rlpEx.FeeTo, rlpEx.Symbol)
	if err != nil {
		return nil, err
	}
	ex.AllPairs = rlpEx.AllPairs
	for _, pair := range rlpEx.AllPairs {
		token0, token1 := ParseKey(pair.Key)
		token1Addr, ok := ex.Pair[token0]
		if ok {
			token1Addr[token1] = pair.Address
			ex.Pair[token0] = token1Addr
		} else {
			ex.Pair[token0] = map[hasharry.Address]hasharry.Address{token1: pair.Address}
		}
	}
	return ex, nil
}

func pairKey(token0 hasharry.Address, token1 hasharry.Address) string {
	return fmt.Sprintf("%s-%s", token0.String(), token1.String())
}

func ParseKey(key string) (hasharry.Address, hasharry.Address) {
	strList := strings.Split(key, "-")
	if len(strList) != 2 {
		return hasharry.Address{}, hasharry.Address{}
	}
	return hasharry.StringToAddress(strList[0]), hasharry.StringToAddress(strList[1])
}

type PairInfo struct {
	Address string `json:"address"`
	Token0  string `json:"token0"`
	Symbol0 string `json:"symbol0"`
	Token1  string `json:"token1"`
	Symbol1 string `json:"symbol1"`
}

func (e *Exchange) Pairs() []PairInfo {
	var infoList []PairInfo
	for _, pair := range e.AllPairs {
		token0, token1 := ParseKey(pair.Key)
		infoList = append(infoList, PairInfo{
			Address: pair.Address.String(),
			Token0:  token0.String(),
			Symbol0: pair.Symbol0,
			Token1:  token1.String(),
			Symbol1: pair.Symbol1,
		})
	}
	return infoList
}

func (e *Exchange) ExchangeRouter(tokenA, tokenB string) [][]string {
	pairList := []map[string]string{}
	for token0, token1Addr := range e.Pair {
		for token1, _ := range token1Addr {
			pairList = append(pairList, map[string]string{
				token0.String(): token1.String(),
			})
		}
	}
	if len(pairList) == 0 {
		return nil
	}
	return CalculatePaths(tokenA, tokenB, pairList)
}

func (e *Exchange) LegalPair(tokenA, tokenB string) (bool, error) {
	/*mainToken := param.Token.String()
	if tokenA == mainToken {
		return true, nil
	}
	if tokenB == mainToken {
		return true, nil
	}
	paths := e.ExchangeRouter(tokenA, mainToken)
	if paths == nil || len(paths[0]) > maxPairTokenLen {
		paths := e.ExchangeRouter(tokenB, mainToken)
		if paths == nil || len(paths[0]) > maxPairTokenLen {
			return false, fmt.Errorf("the path of %s->%s must be smaller than %d", tokenA, mainToken, maxPairTokenLen)
		}
	}*/
	return true, nil
}

func CalculateShortestPath(tokenA, tokenB string, pairs []map[string]string) []string {
	paths := CalculateShortestPaths(tokenA, tokenB, pairs)
	if paths == nil {
		return nil
	}
	if len(paths) != 0 {
		return paths[0]
	} else {
		return nil
	}
}

func CalculateShortestPaths(tokenA, tokenB string, pairs []map[string]string) [][]string {
	g := utils.NewGraph()
	for _, pair := range pairs {
		for token0, token1 := range pair {
			g.AddEdge(utils.NewNode(token0, 0), utils.NewNode(token1, 0))
		}
	}
	paths, err := g.FindNodePath(utils.NewNode(tokenA, 0), utils.NewNode(tokenB, 0))
	if err != nil {
		return nil
	}
	minLen := math.MaxInt32
	pathMap := map[int][][]string{}
	for _, path := range paths {
		if len(path) < minLen {
			minLen = len(path)
		}
		pathString := []string{}
		for _, node := range path {
			pathString = append(pathString, node.String())
		}
		pathList, ok := pathMap[len(path)]
		if ok {
			pathList = append(pathList, pathString)
			pathMap[len(path)] = pathList
		} else {
			pathMap[len(path)] = [][]string{pathString}
		}
	}

	return pathMap[minLen]
}

func CalculatePaths(tokenA, tokenB string, pairs []map[string]string) [][]string {
	g := utils.NewGraph()
	for _, pair := range pairs {
		for token0, token1 := range pair {
			g.AddEdge(utils.NewNode(token0, 0), utils.NewNode(token1, 0))
		}
	}
	paths, err := g.FindNodePath(utils.NewNode(tokenA, 0), utils.NewNode(tokenB, 0))
	if err != nil {
		return nil
	}
	allPath := [][]string{}

	for _, path := range paths {
		pathList := make([]string, len(path))
		for i, node := range path {
			pathList[i] = node.String()
		}
		allPath = append(allPath, pathList)
	}

	return allPath
}
