package exchange_runner

import (
	"errors"
	"fmt"
	"github.com/UBChainNet/UBChain/common/codec"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/core/runner/library"
	"github.com/UBChainNet/UBChain/core/runner/method"
	"github.com/UBChainNet/UBChain/core/types"
	"github.com/UBChainNet/UBChain/core/types/contractv2"
	"github.com/UBChainNet/UBChain/core/types/contractv2/exchange"
	"github.com/UBChainNet/UBChain/core/types/functionbody/exchange_func"
	"github.com/UBChainNet/UBChain/crypto/base58"
	"github.com/UBChainNet/UBChain/param"
	"github.com/UBChainNet/UBChain/ut"
	"math/big"
	"strings"
)

type ExchangeState struct {
	header  *contractv2.ContractV2
	body    *exchange.Exchange
	library *library.RunnerLibrary
}

func NewExchangeState(runnerLibrary *library.RunnerLibrary, exAddress string) (*ExchangeState, error) {
	exHeader := runnerLibrary.GetContractV2(exAddress)
	if exHeader == nil {
		return nil, fmt.Errorf("exchange %s already exist", exAddress)
	}
	exBody, _ := exHeader.Body.(*exchange.Exchange)
	return &ExchangeState{
		header:  exHeader,
		body:    exBody,
		library: runnerLibrary,
	}, nil
}

func (es *ExchangeState) Methods() map[string]*method.MethodInfo {
	return method.ExMethods
}

func (es *ExchangeState) MethodExist(mth string) bool {
	_, exist := method.ExMethods[mth]
	return exist
}

func (es *ExchangeState) Pairs() []exchange.PairInfo {
	return es.body.Pairs()
}

func (es *ExchangeState) PairAddress(tokenA, tokenB string) string {
	token0, token1 := library.SortToken(hasharry.StringToAddress(tokenA), hasharry.StringToAddress(tokenB))
	return es.body.PairAddress(token0, token1).String()
}

func (es *ExchangeState) ExchangeRouter(tokenA, tokenB string) [][]string {
	return es.body.ExchangeRouter(tokenA, tokenB)
}

type Router struct {
	Path      []string `json:"path"`
	AmountOut float64  `json:"amountOut"`
	Error     string   `json:"error"`
}

type RouterIn struct {
	Path     []string `json:"path"`
	AmountIn float64  `json:"amountIn"`
	Error    string   `json:"error"`
}

func (es *ExchangeState) ExchangeRouterWithAmount(tokenA, tokenB string, amountIn float64) []Router {
	var routers []Router
	paths := es.body.ExchangeRouter(tokenA, tokenB)
	for _, path := range paths {
		var errInfo = ""
		amountOut, err := es.amountOut(path, amountIn)
		if err != nil {
			errInfo = err.Error()
		}
		routers = append(routers, Router{
			Path:      path,
			AmountOut: amountOut,
			Error:     errInfo,
		})
	}
	return routers
}

func (es *ExchangeState) ExchangeOptimalRouter(tokenA, tokenB string, amountIn float64) *Router {
	paths := es.body.ExchangeRouter(tokenA, tokenB)
	var maxOut float64
	var optimal []string
	for _, path := range paths {
		amountOut, err := es.amountOut(path, amountIn)
		if err != nil {
			continue
		}
		if amountOut > maxOut {
			maxOut = amountOut
			optimal = path
		}
	}
	if maxOut == float64(0) {
		return &Router{
			Path:      nil,
			AmountOut: 0,
			Error:     "unable to change",
		}
	}
	return &Router{
		Path:      optimal,
		AmountOut: maxOut,
		Error:     "",
	}
}

func (es *ExchangeState) ExchangeOptimalRouterAmountIn(tokenA, tokenB string, amountOut float64) *RouterIn {
	paths := es.body.ExchangeRouter(tokenA, tokenB)
	var minIn float64
	var optimal []string
	for _, path := range paths {
		amountIn, err := es.amountIn(path, amountOut)
		if err != nil {
			continue
		}
		if minIn == float64(0) {
			minIn = amountIn
			optimal = path
		} else if amountIn < minIn {
			minIn = amountIn
			optimal = path
		}
	}
	if minIn == float64(0) {
		return &RouterIn{
			Path:     nil,
			AmountIn: 0,
			Error:    "unable to change",
		}
	}
	return &RouterIn{
		Path:     optimal,
		AmountIn: minIn,
		Error:    "",
	}
}

func (es *ExchangeState) LegalPair(tokenA, tokenB string) (bool, error) {
	return es.body.LegalPair(tokenA, tokenB)
}

func (es *ExchangeState) AmountOut(paths string, amountIn float64) (float64, error) {
	pathsList := strings.Split(paths, ",")
	return es.amountOut(pathsList, amountIn)
}

func (es *ExchangeState) amountOut(paths []string, amountIn float64) (float64, error) {
	arryPaths := make([]hasharry.Address, len(paths))
	for i, path := range paths {
		arryPaths[i] = hasharry.StringToAddress(path)
	}
	iAmountIn, _ := types.NewAmount(amountIn)
	outs, err := es.getAmountsOut(iAmountIn, arryPaths)
	if err != nil {
		return 0, err
	}
	return types.Amount(outs[len(outs)-1]).ToCoin(), nil
}

func (es *ExchangeState) AmountIn(paths string, amountOut float64) (float64, error) {
	pathsList := strings.Split(paths, ",")
	arryPaths := make([]hasharry.Address, len(pathsList))
	for i, path := range pathsList {
		arryPaths[i] = hasharry.StringToAddress(path)
	}
	iAmountOut, _ := types.NewAmount(amountOut)
	ins, err := es.getAmountsIn(iAmountOut, arryPaths)
	if err != nil {
		return 0, err
	}
	return types.Amount(ins[0]).ToCoin(), nil
}

func (es *ExchangeState) amountIn(paths []string, amountOut float64) (float64, error) {
	arryPaths := make([]hasharry.Address, len(paths))
	for i, path := range paths {
		arryPaths[i] = hasharry.StringToAddress(path)
	}
	iAmountOut, _ := types.NewAmount(amountOut)
	ins, err := es.getAmountsIn(iAmountOut, arryPaths)
	if err != nil {
		return 0, err
	}
	return types.Amount(ins[0]).ToCoin(), nil
}

func (es *ExchangeState) getAmountsOut(amountIn uint64, path []hasharry.Address) ([]uint64, error) {
	amounts := make([]uint64, len(path))
	amounts[0] = amountIn
	for i := 0; i < len(path)-1; i++ {
		// ???????????????
		token0, token1 := library.SortToken(path[i], path[i+1])
		pairAddress := es.body.PairAddress(token0, token1)
		reserveIn, reserveOut, err := es.library.GetReservesByPairAddress(pairAddress, path[i], path[i+1])
		if err != nil {
			return nil, err
		}
		// ??????????????? =  ???????????????????????????
		amounts[i+1], err = GetAmountOut(amounts[i], reserveIn, reserveOut)
		if err != nil {
			return amounts, err
		}
	}
	return amounts, nil
}

// getAmountsIn performs chained getAmountIn calculations on any number of pairs
func (es *ExchangeState) getAmountsIn(amountOut uint64, path []hasharry.Address) ([]uint64, error) {
	amounts := make([]uint64, len(path))
	amounts[len(amounts)-1] = amountOut
	for i := len(path) - 1; i > 0; i-- {
		// ???????????????
		token0, token1 := library.SortToken(path[i-1], path[i])
		pairAddress := es.body.PairAddress(token0, token1)
		reserveIn, reserveOut, err := es.library.GetReservesByPairAddress(pairAddress, path[i-1], path[i])
		if err != nil {
			return nil, err
		}
		amounts[i-1], err = GetAmountIn(amounts[i], reserveIn, reserveOut)
		if err != nil {
			return amounts, err
		}
	}
	return amounts, nil
}

type ExchangeRunner struct {
	exState      *ExchangeState
	address      hasharry.Address
	tx           types.ITransaction
	contractBody *types.TxContractV2Body
	pairList     []*contractv2.ContractV2
	events       []*types.Event
	height       uint64
}

func NewExchangeRunner(lib *library.RunnerLibrary, tx types.ITransaction, height uint64) *ExchangeRunner {
	var ex *exchange.Exchange
	address := tx.GetTxBody().GetContract()
	exHeader := lib.GetContractV2(address.String())
	if exHeader != nil {
		ex = exHeader.Body.(*exchange.Exchange)
	}

	contractBody := tx.GetTxBody().(*types.TxContractV2Body)
	return &ExchangeRunner{
		exState: &ExchangeState{
			header:  exHeader,
			body:    ex,
			library: lib,
		},
		address:      address,
		tx:           tx,
		contractBody: contractBody,
		events:       make([]*types.Event, 0),
		height:       height,
	}
}

func (e *ExchangeRunner) PreInitVerify() error {
	if e.exState.header != nil {
		return fmt.Errorf("exchange %s already exist", e.address.String())
	}
	funcBody, _ := e.contractBody.Function.(*exchange_func.ExchangeInitBody)
	_, err := e.exState.library.GetSymbolContract(funcBody.Symbol)
	return err
}

func (e *ExchangeRunner) PreSetVerify() error {
	if e.exState.header == nil {
		return fmt.Errorf("exchange %s is not exist", e.address.String())
	}
	return e.exState.body.VerifySetter(e.tx.From())
}

func (e *ExchangeRunner) PreExactInVerify(lastHeight uint64) error {
	if e.exState.header == nil {
		return fmt.Errorf("exchange is not exist")
	}
	funcBody, _ := e.contractBody.Function.(*exchange_func.ExactIn)
	if funcBody == nil {
		return errors.New("wrong contractV2 function")
	}
	if len(funcBody.Path) < 2 {
		return errors.New("invalid path")
	}
	for i := 0; i < len(funcBody.Path)-1; i++ {
		if exist := e.exState.body.Exist(library.SortToken(funcBody.Path[i], funcBody.Path[i+1])); !exist {
			return fmt.Errorf("the pair of %s and %s does not exist", funcBody.Path[i].String(), funcBody.Path[i+1].String())
		}
	}
	if funcBody.Deadline != 0 && funcBody.Deadline < lastHeight {
		return fmt.Errorf("past the deadline")
	}
	balance := e.exState.library.GetBalance(e.tx.From(), funcBody.Path[0])
	if funcBody.Path[0].IsEqual(param.Token) {
		if balance < funcBody.AmountIn+e.tx.GetFees() {
			return errors.New("balance not enough")
		}
	} else {
		if balance < funcBody.AmountIn {
			return errors.New("balance not enough")
		}
	}
	return nil
}

func (e *ExchangeRunner) PreExactOutVerify(lastHeight uint64) error {
	if e.exState.header == nil {
		return fmt.Errorf("exchange is not exist")
	}
	funcBody, _ := e.contractBody.Function.(*exchange_func.ExactOut)
	if funcBody == nil {
		return errors.New("wrong contractV2 function")
	}
	if len(funcBody.Path) < 2 {
		return errors.New("invalid path")
	}
	for i := 0; i < len(funcBody.Path)-1; i++ {
		if exist := e.exState.body.Exist(library.SortToken(funcBody.Path[i], funcBody.Path[i+1])); !exist {
			return fmt.Errorf("the pair of %s and %s does not exist", funcBody.Path[i].String(), funcBody.Path[i+1].String())
		}
	}
	if funcBody.Deadline != 0 && funcBody.Deadline < lastHeight {
		return fmt.Errorf("past the deadline")
	}
	balance := e.exState.library.GetBalance(e.tx.From(), funcBody.Path[0])
	if funcBody.Path[0].IsEqual(param.Token) {
		if balance < funcBody.AmountInMax+e.tx.GetFees() {
			return errors.New("balance not enough")
		}
	} else {
		if balance < funcBody.AmountInMax {
			return errors.New("balance not enough")
		}
	}
	return nil
}

func (e *ExchangeRunner) Init() {
	var ERR error
	state := &types.ContractV2State{State: types.Contract_Success}
	defer func() {
		if ERR != nil {
			state.State = types.Contract_Failed
			state.Error = ERR.Error()
		} else {
			state.Event = e.events
		}
		e.exState.library.SetContractV2State(e.tx.Hash().String(), state)
	}()

	contract := &contractv2.ContractV2{
		Address:    e.contractBody.Contract,
		CreateHash: e.tx.Hash(),
		Type:       e.contractBody.Type,
		Body:       nil,
	}
	if e.exState.header != nil {
		ERR = fmt.Errorf("exchange %s already exist", contract.Address.String())
		return
	}
	initBody := e.contractBody.Function.(*exchange_func.ExchangeInitBody)
	contract.Body, _ = exchange.NewExchange(initBody.Admin, initBody.FeeTo, initBody.Symbol)
	e.exState.library.SetSymbol(initBody.Symbol, contract.Address.String())
	e.exState.library.SetContractV2(contract)
}

func (e *ExchangeRunner) SetAdmin() {
	var ERR error
	state := &types.ContractV2State{State: types.Contract_Success}
	defer func() {
		if ERR != nil {
			state.State = types.Contract_Failed
			state.Error = ERR.Error()
		} else {
			state.Event = e.events
		}
		e.exState.library.SetContractV2State(e.tx.Hash().String(), state)
	}()

	if e.exState.header == nil {
		ERR = fmt.Errorf("exchanges %s is not exist", e.tx.GetTxBody().GetContract().String())
		return
	}
	funcBody, _ := e.contractBody.Function.(*exchange_func.ExchangeAdmin)
	ex, _ := e.exState.header.Body.(*exchange.Exchange)
	if err := ex.SetAdmin(funcBody.Address, e.tx.From()); err != nil {
		ERR = err
		return
	}
	e.exState.header.Body = ex
	e.exState.library.SetContractV2(e.exState.header)
}

func (e *ExchangeRunner) SetFeeTo() {
	var ERR error
	state := &types.ContractV2State{State: types.Contract_Success}
	defer func() {
		if ERR != nil {
			state.State = types.Contract_Failed
			state.Error = ERR.Error()
		} else {
			state.Event = e.events
		}
		e.exState.library.SetContractV2State(e.tx.Hash().String(), state)
	}()

	if e.exState.header == nil {
		ERR = fmt.Errorf("exchanges %s is not exist", e.tx.GetTxBody().GetContract().String())
		return
	}
	funcBody, _ := e.contractBody.Function.(*exchange_func.ExchangeFeeTo)
	ex, _ := e.exState.header.Body.(*exchange.Exchange)
	if err := ex.SetFeeTo(funcBody.Address, e.tx.From()); err != nil {
		ERR = err
		return
	}
	e.exState.header.Body = ex
	e.exState.library.SetContractV2(e.exState.header)
}

type SwapExactIn struct {
	AmountOut uint64 `json:"amountOut"`
}

func (e *ExchangeRunner) SwapExactIn(blockTime uint64) {
	var ERR error
	var err error
	var SwapInfo SwapExactIn
	var amounts []uint64
	state := &types.ContractV2State{State: types.Contract_Success}
	defer func() {
		if ERR != nil {
			state.State = types.Contract_Failed
			state.Error = ERR.Error()
		} else {
			state.Event = e.events
		}
		if e.height < param.UIPBlock2 {
			state.Event = e.events
		}
		e.exState.library.SetContractV2State(e.tx.Hash().String(), state)
	}()

	funcBody, _ := e.contractBody.Function.(*exchange_func.ExactIn)

	if funcBody.Deadline != 0 && funcBody.Deadline < e.height {
		ERR = fmt.Errorf("past the deadline")
		return
	}
	amounts, err = e.exState.getAmountsOut(funcBody.AmountIn, funcBody.Path)
	if err != nil {
		ERR = err
		return
	}
	SwapInfo.AmountOut = amounts[len(amounts)-1]
	if SwapInfo.AmountOut < funcBody.AmountOutMin {
		ERR = fmt.Errorf("outAmount %d is less than the minimum output %d", SwapInfo.AmountOut, funcBody.AmountOutMin)
		return
	}
	pair0 := e.exState.body.PairAddress(library.SortToken(funcBody.Path[0], funcBody.Path[1]))
	if err = e.swapAmounts(amounts, funcBody.Path, funcBody.To, blockTime); err != nil {
		ERR = err
		return
	}

	e.transferEvent(e.tx.From(), pair0, funcBody.Path[0], amounts[0])

	if err = e.runEvents(); err != nil {
		ERR = err
		return
	}

	e.update()
}

func (e *ExchangeRunner) update() {
	for _, pairContract := range e.pairList {
		e.exState.library.SetContractV2(pairContract)
	}
}

// requires the initial amount to have already been sent to the first pair
func (e *ExchangeRunner) swapAmounts(amounts []uint64, path []hasharry.Address, to hasharry.Address, blockTime uint64) error {
	var amount0Out, amount1Out uint64
	var amount0In, amount1In uint64
	for i := 0; i < len(path)-1; i++ {
		input, output := path[i], path[i+1]
		token0, _ := library.SortToken(input, output)
		amountOut := amounts[i+1]
		amountIn := amounts[i]
		if input.IsEqual(token0) {
			amount0Out, amount1Out = 0, amountOut
			amount0In, amount1In = amountIn, 0
		} else {
			amount0Out, amount1Out = amountOut, 0
			amount0In, amount1In = 0, amountIn
		}
		toAddr := to
		if i < len(path)-2 {
			toAddr = e.exState.body.PairAddress(library.SortToken(output, path[i+2]))
		}
		if err := e.swap(input, output, amount0In, amount1In, amount0Out, amount1Out, toAddr, blockTime); err != nil {
			return err
		}
	}
	return nil
}

func (e *ExchangeRunner) swap(tokenA, tokenB hasharry.Address, amount0In, amount1In, amount0Out, amount1Out uint64, to hasharry.Address, blockTime uint64) error {
	if amount0Out <= 0 && amount1Out <= 0 {
		return errors.New("insufficient output amount")
	}
	_token0, _token1 := library.SortToken(tokenA, tokenB)
	pairAddress := e.exState.body.PairAddress(_token0, _token1)
	pairContract := e.exState.library.GetContractV2(pairAddress.String())
	pair := pairContract.Body.(*exchange.Pair)
	_reserve0, _reserve1, err := e.exState.library.GetReservesByPairAddress(pairAddress, _token0, _token1)
	if err != nil {
		return err
	}
	if amount0Out >= _reserve0 || amount1Out >= _reserve1 {
		return errors.New("insufficient liquidity")
	}

	var balance0, balance1 uint64
	if to.IsEqual(_token0) || to.IsEqual(_token1) {
		return errors.New("invalid to")
	}
	// ?????????to??????
	if amount0Out > 0 {
		e.transferEvent(pairAddress, to, _token0, amount0Out)
	}
	if amount1Out > 0 {
		e.transferEvent(pairAddress, to, _token1, amount1Out)
	}

	// ????????????
	if e.height < param.UIPBlock1 {
		// pair?????????lock?????????????????????
		balance0 = e.exState.library.GetBalance(pairAddress, _token0)
		balance1 = e.exState.library.GetBalance(pairAddress, _token1)
	} else {
		balance0 = _reserve0
		balance1 = _reserve1
	}

	if amount0In > 0 {
		balance0 = balance0 + amount0In
	} else {
		balance1 = balance1 + amount1In
	}

	if amount0Out > 0 {
		if balance0 < amount0Out {
			return errors.New("insufficient liquidity")
		}
		balance0 = balance0 - amount0Out
	} else {
		if balance1 < amount1Out {
			return errors.New("insufficient liquidity")
		}
		balance1 = balance1 - amount1Out
	}

	//????????????????????????????????????
	/*if balance0 > _reserve0-amount0Out {
		amount0In = balance0 - (_reserve0 - amount0Out)
	} else {
		amount0In = 0
	}

	if balance1 > _reserve1-amount1Out {
		amount1In = balance1 - (_reserve1 - amount1Out)
	} else {
		amount1In = 0
	}*/
	if amount0In <= 0 && amount1In <= 0 {
		return errors.New("insufficient input amount")
	}
	// balance0Adjusted = balance0 * 1000 - amount0In * 3
	balance0Adjusted := big.NewInt(0).Sub(big.NewInt(0).Mul(big.NewInt(int64(balance0)), big.NewInt(1000)),
		big.NewInt(0).Mul(big.NewInt(int64(amount0In)), big.NewInt(3)))
	// balance1Adjusted = balance1 * 1000 - amount1In * 3
	balance1Adjusted := big.NewInt(0).Sub(big.NewInt(0).Mul(big.NewInt(int64(balance1)), big.NewInt(1000)),
		big.NewInt(0).Mul(big.NewInt(int64(amount1In)), big.NewInt(3)))

	// ??????k?????????K?????????????????????????????????
	// x = balance0Adjusted * balance1Adjusted
	x := big.NewInt(0).Mul(balance0Adjusted, balance1Adjusted)
	// y = _reserve0 * _reserve1 * 1000^2
	y := big.NewInt(0).Mul(big.NewInt(0).Mul(big.NewInt(int64(_reserve0)), big.NewInt(int64(_reserve1))), big.NewInt(1000^2))
	if x.Cmp(y) < 0 {
		return errors.New("K")
	}
	pair.UpdateReserve(balance0, balance1, _reserve0, _reserve1, blockTime)
	pairContract.Body = pair
	e.pairList = append(e.pairList, pairContract)
	return nil
}

// GetAmountOut given an input amount of an asset and pair reserves, returns the maximum output amount of the other asset
func GetAmountOut(amountIn, reserveIn, reserveOut uint64) (uint64, error) {
	if amountIn <= 0 {
		return 0, errors.New("insufficient input amount")
	}
	if reserveIn <= 0 || reserveOut <= 0 {
		return 0, errors.New("insufficient liquidity")
	}
	// amountInWithFee = amountIn * 997
	// 0.5% fees
	amountInWithFee := big.NewInt(0).Mul(big.NewInt(int64(amountIn)), big.NewInt(997))
	// numerator = amountInWithFee * reserveOut
	numerator := big.NewInt(0).Mul(amountInWithFee, big.NewInt(int64(reserveOut)))
	// denominator = reserveIn * 1000 + amountInWithFee

	denominator := big.NewInt(0).Add(big.NewInt(0).Mul(big.NewInt(int64(reserveIn)), big.NewInt(1000)), amountInWithFee)
	amountOut := big.NewInt(0).Div(numerator, denominator)
	return amountOut.Uint64(), nil
}

// GetAmountIn given an output amount of an asset and pair reserves, returns a required input amount of the other asset
func GetAmountIn(amountOut, reserveIn, reserveOut uint64) (uint64, error) {
	if amountOut <= 0 {
		return 0, errors.New("insufficient output amount")
	}
	if reserveIn <= 0 || reserveOut <= 0 {
		return 0, errors.New("insufficient liquidity")
	}
	if reserveOut < amountOut {
		return 0, errors.New("insufficient liquidity")
	}
	/*	amountOut = amountOut / 1000000
		reserveIn = reserveIn / 10000000
		reserveOut = reserveOut / 10000000*/
	// numerator = amountOut * reserveIn * 1000
	numerator := big.NewInt(0).Mul(big.NewInt(0).Mul(big.NewInt(int64(amountOut)), big.NewInt(int64(reserveIn))), big.NewInt(1000))
	// denominator = (reserveOut - amountOut) (* 997)
	denominator := big.NewInt(0).Mul(big.NewInt(0).Sub(big.NewInt(int64(reserveOut)), big.NewInt(int64(amountOut))), big.NewInt(997))
	// amountIn = (numerator\denominator) + 1
	x := big.NewInt(0).Div(numerator, denominator)

	amountIn := big.NewInt(0).Add(x, big.NewInt(1))
	return amountIn.Uint64(), nil
}

func (e *ExchangeRunner) SwapExactOut(blockTime uint64) {
	var ERR error
	var err error
	var amounts []uint64
	state := &types.ContractV2State{State: types.Contract_Success}
	defer func() {
		if ERR != nil {
			state.State = types.Contract_Failed
			state.Error = ERR.Error()
		} else {
			state.Event = e.events
		}
		e.exState.library.SetContractV2State(e.tx.Hash().String(), state)
	}()

	funcBody, _ := e.contractBody.Function.(*exchange_func.ExactOut)

	if funcBody.Deadline != 0 && funcBody.Deadline < e.height {
		ERR = fmt.Errorf("past the deadline")
		return
	}
	amounts, err = e.exState.getAmountsIn(funcBody.AmountOut, funcBody.Path)
	if err != nil {
		ERR = err
		return
	}
	amountIn := amounts[0]
	if amountIn > funcBody.AmountInMax {
		ERR = fmt.Errorf("amountIn %d is greater than the maximum input amount %d", amounts[0], funcBody.AmountInMax)
		return
	}
	pair0 := e.exState.body.PairAddress(library.SortToken(funcBody.Path[0], funcBody.Path[1]))
	if err := e.swapAmounts(amounts, funcBody.Path, funcBody.To, blockTime); err != nil {
		ERR = err
		return
	}

	e.transferEvent(e.tx.From(), pair0, funcBody.Path[0], amountIn)

	if err = e.runEvents(); err != nil {
		ERR = err
		return
	}
	e.update()
}

func (e *ExchangeRunner) transferEvent(from, to, token hasharry.Address, amount uint64) {
	e.events = append(e.events, &types.Event{
		EventType: types.Event_Transfer,
		From:      from,
		To:        to,
		Token:     token,
		Amount:    amount,
		Height:    e.height,
	})
}

func (e *ExchangeRunner) runEvents() error {
	for _, event := range e.events {
		if err := e.exState.library.PreRunEvent(event); err != nil {
			return err
		}
	}
	for _, event := range e.events {
		e.exState.library.RunEvent(event)
	}
	return nil
}

func ExchangeAddress(net, from string, nonce uint64) (string, error) {
	bytes := make([]byte, 0)
	nonceBytes := codec.Uint64toBytes(nonce)
	bytes = append(base58.Decode(from), nonceBytes...)
	return ut.GenerateContractV2Address(net, bytes)
}
