package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jhdriver/UBChain/common/hasharry"
	"github.com/jhdriver/UBChain/core/types"
	"github.com/jhdriver/UBChain/rpc"
	"github.com/jhdriver/UBChain/rpc/rpctypes"
	"github.com/jhdriver/UBChain/ut/transaction"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

func init() {
	exchangeCmds := []*cobra.Command{
		CreateExchangeCmd,
		SetExchangeAdminCmd,
		SetExchangeFeeToCmd,
		AddLiquidityCmd,
		RemoveLiquidityCmd,
		SwapExactInCmd,
		SwapExactOutCmd,
		GetAllPairsCmd,
	}
	RootCmd.AddCommand(exchangeCmds...)
	RootSubCmdGroups["exchange"] = exchangeCmds

}

var CreateExchangeCmd = &cobra.Command{
	Use:     "CreateExchange {from} {admin} {feeTo} {password} {nonce}; Create a decentralized exchange;",
	Aliases: []string{"CreateExchange", "createexchange", "ce", "CE"},
	Short:   "CreateExchange {from} {admin} {feeTo} {password} {nonce}; Create a decentralized exchange;",
	Example: `
	CreateExchange 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 123456
		OR
	CreateExchange 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 123456 1
	`,
	Args: cobra.MinimumNArgs(3),
	Run:  CreateExchange,
}

func CreateExchange(cmd *cobra.Command, args []string) {
	var passwd []byte
	var err error
	if len(args) > 3 {
		passwd = []byte(args[3])
	} else {
		fmt.Println("please input password：")
		passwd, err = readPassWd()
		if err != nil {
			outputError(cmd.Use, fmt.Errorf("read password failed! %s", err.Error()))
			return
		}
	}
	privKey, err := ReadAddrPrivate(getAddJsonPath(args[0]), passwd)
	if err != nil {
		outputError(cmd.Use, fmt.Errorf("wrong password"))
		return
	}
	resp, err := GetAccountByRpc(args[0])
	if err != nil {
		outputError(cmd.Use, err)
		return
	}
	if resp.Code != 0 {
		outputRespError(cmd.Use, resp)
		return
	}
	var account *rpctypes.Account
	if err := json.Unmarshal(resp.Result, &account); err != nil {
		outputError(cmd.Use, err)
		return
	}

	tx, err := parseCEParams(args, account.Nonce+1)
	if err != nil {
		outputError(cmd.Use, err)
		return
	}

	if !signTx(cmd, tx, privKey.Private) {
		outputError(cmd.Use, errors.New("signature failure"))
		return
	}

	rs, err := sendTx(cmd, tx)
	if err != nil {
		outputError(cmd.Use, err)
	} else if rs.Code != 0 {
		outputRespError(cmd.Use, rs)
	} else {
		fmt.Println()
		fmt.Println(string(rs.Result))
	}
}

func parseCEParams(args []string, nonce uint64) (*types.Transaction, error) {
	var err error
	from := hasharry.StringToAddress(args[0])
	admin := args[1]
	feeTo := args[2]
	if len(args) > 4 {
		nonce, err = strconv.ParseUint(args[4], 10, 64)
		if err != nil {
			return nil, errors.New("wrong nonce")
		}
	}
	tx, err := transaction.NewExchange(Net, from.String(), admin, feeTo, nonce, "")
	if err != nil {
		return nil, err
	}
	return tx, nil
}

var SetExchangeAdminCmd = &cobra.Command{
	Use:     "SetExchangeAdmin {from} {exchange} {admin} {password} {nonce}; Set exchange feeTo setter;",
	Aliases: []string{"setexchangeadmin", "sea", "SEA"},
	Short:   "SetExchangeAdmin {from} {exchange} {admin} {password} {nonce}; Set exchange feeTo setter;",
	Example: `
	SetExchangeAdmin UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UWTfBGxDMZX19vjnacXVkP51min9EjhYq43W UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw 123456
		OR
	SetExchangeAdmin UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UWTfBGxDMZX19vjnacXVkP51min9EjhYq43W UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw 123456 1
	`,
	Args: cobra.MinimumNArgs(3),
	Run:  SetExchangeAdmin,
}

func SetExchangeAdmin(cmd *cobra.Command, args []string) {
	var passwd []byte
	var err error
	if len(args) > 3 {
		passwd = []byte(args[3])
	} else {
		fmt.Println("please input password：")
		passwd, err = readPassWd()
		if err != nil {
			outputError(cmd.Use, fmt.Errorf("read password failed! %s", err.Error()))
			return
		}
	}
	privKey, err := ReadAddrPrivate(getAddJsonPath(args[0]), passwd)
	if err != nil {
		outputError(cmd.Use, fmt.Errorf("wrong password"))
		return
	}
	resp, err := GetAccountByRpc(args[0])
	if err != nil {
		outputError(cmd.Use, err)
		return
	}
	if resp.Code != 0 {
		outputRespError(cmd.Use, resp)
		return
	}
	var account *rpctypes.Account
	if err := json.Unmarshal(resp.Result, &account); err != nil {
		outputError(cmd.Use, err)
		return
	}

	tx, err := parseSEFTSParams(args, account.Nonce+1)
	if err != nil {
		outputError(cmd.Use, err)
		return
	}

	if !signTx(cmd, tx, privKey.Private) {
		outputError(cmd.Use, errors.New("signature failure"))
		return
	}

	rs, err := sendTx(cmd, tx)
	if err != nil {
		outputError(cmd.Use, err)
	} else if rs.Code != 0 {
		outputRespError(cmd.Use, rs)
	} else {
		fmt.Println()
		fmt.Println(string(rs.Result))
	}
}

func parseSEFTSParams(args []string, nonce uint64) (*types.Transaction, error) {
	var err error
	from := args[0]
	exchange := args[1]
	admin := args[2]
	if len(args) > 4 {
		nonce, err = strconv.ParseUint(args[4], 10, 64)
		if err != nil {
			return nil, errors.New("wrong nonce")
		}
	}
	tx, err := transaction.NewSetAdmin(from, exchange, admin, nonce, "")
	if err != nil {
		return nil, err
	}
	return tx, nil
}

var SetExchangeFeeToCmd = &cobra.Command{
	Use:     "SetExchangeFeeTo {from} {exchange} {feeTo} {password} {nonce}; Set exchange feeTo;",
	Aliases: []string{"setexchangefeeto", "seft", "SEFT"},
	Short:   "SetExchangeFeeTo {from} {exchange} {feeTo} {password} {nonce}; Set exchange feeTo;",
	Example: `
	SetExchangeFeeTo UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UWTfBGxDMZX19vjnacXVkP51min9EjhYq43W UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw 123456
		OR
	SetExchangeFeeTo UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UWTfBGxDMZX19vjnacXVkP51min9EjhYq43W UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw 123456 1
	`,
	Args: cobra.MinimumNArgs(3),
	Run:  SetExchangeFeeTo,
}

func SetExchangeFeeTo(cmd *cobra.Command, args []string) {
	var passwd []byte
	var err error
	if len(args) > 3 {
		passwd = []byte(args[3])
	} else {
		fmt.Println("please input password：")
		passwd, err = readPassWd()
		if err != nil {
			outputError(cmd.Use, fmt.Errorf("read password failed! %s", err.Error()))
			return
		}
	}
	privKey, err := ReadAddrPrivate(getAddJsonPath(args[0]), passwd)
	if err != nil {
		outputError(cmd.Use, fmt.Errorf("wrong password"))
		return
	}
	resp, err := GetAccountByRpc(args[0])
	if err != nil {
		outputError(cmd.Use, err)
		return
	}
	if resp.Code != 0 {
		outputRespError(cmd.Use, resp)
		return
	}
	var account *rpctypes.Account
	if err := json.Unmarshal(resp.Result, &account); err != nil {
		outputError(cmd.Use, err)
		return
	}

	tx, err := parseSEFTParams(args, account.Nonce+1)
	if err != nil {
		outputError(cmd.Use, err)
		return
	}

	if !signTx(cmd, tx, privKey.Private) {
		outputError(cmd.Use, errors.New("signature failure"))
		return
	}

	rs, err := sendTx(cmd, tx)
	if err != nil {
		outputError(cmd.Use, err)
	} else if rs.Code != 0 {
		outputRespError(cmd.Use, rs)
	} else {
		fmt.Println()
		fmt.Println(string(rs.Result))
	}
}

func parseSEFTParams(args []string, nonce uint64) (*types.Transaction, error) {
	var err error
	from := args[0]
	exchange := args[1]
	feeTo := args[2]
	if len(args) > 4 {
		nonce, err = strconv.ParseUint(args[4], 10, 64)
		if err != nil {
			return nil, errors.New("wrong nonce")
		}
	}
	tx, err := transaction.NewSetFeeTo(from, exchange, feeTo, nonce, "")
	if err != nil {
		return nil, err
	}
	return tx, nil
}

var AddLiquidityCmd = &cobra.Command{
	Use:     "AddLiquidity {from} {to} {exchange} {tokenA} {amountADesired} {amountAmin} {tokenB} {amountBDesired} {amountBMin} {password} {nonce}; Create and add liquidity;",
	Aliases: []string{"addliquidity", "al", "AL"},
	Short:   "AddLiquidity {from} {to} {exchange} {tokenA} {amountADesired} {amountAmin} {tokenB} {amountBDesired} {amountBMin} {password} {nonce}; Create and add liquidity;",
	Example: `
	AddLiquidity UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UWTfBGxDMZX19vjnacXVkP51min9EjhYq43W UWTXEqvUWik48uAHcJXZiyyWMy4GLtpGuttL 100 90 UBC 1 0.9 123456
		OR
	AddLiquidity UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UWTfBGxDMZX19vjnacXVkP51min9EjhYq43W UWTXEqvUWik48uAHcJXZiyyWMy4GLtpGuttL 100 90 UBC 1 0.9 123456 1
	`,
	Args: cobra.MinimumNArgs(9),
	Run:  AddLiquidity,
}

func AddLiquidity(cmd *cobra.Command, args []string) {
	var passwd []byte
	var err error
	if len(args) > 9 {
		passwd = []byte(args[9])
	} else {
		fmt.Println("please input password：")
		passwd, err = readPassWd()
		if err != nil {
			outputError(cmd.Use, fmt.Errorf("read password failed! %s", err.Error()))
			return
		}
	}
	privKey, err := ReadAddrPrivate(getAddJsonPath(args[0]), passwd)
	if err != nil {
		outputError(cmd.Use, fmt.Errorf("wrong password"))
		return
	}
	resp, err := GetAccountByRpc(args[0])
	if err != nil {
		outputError(cmd.Use, err)
		return
	}
	if resp.Code != 0 {
		outputRespError(cmd.Use, resp)
		return
	}
	var account *rpctypes.Account
	if err := json.Unmarshal(resp.Result, &account); err != nil {
		outputError(cmd.Use, err)
		return
	}

	tx, err := parseALParams(args, account.Nonce+1)
	if err != nil {
		outputError(cmd.Use, err)
		return
	}

	if !signTx(cmd, tx, privKey.Private) {
		outputError(cmd.Use, errors.New("signature failure"))
		return
	}

	rs, err := sendTx(cmd, tx)
	if err != nil {
		outputError(cmd.Use, err)
	} else if rs.Code != 0 {
		outputRespError(cmd.Use, rs)
	} else {
		fmt.Println()
		fmt.Println(string(rs.Result))
	}
}

func parseALParams(args []string, nonce uint64) (*types.Transaction, error) {
	var err error
	from := args[0]
	to := args[1]
	exchange := args[2]
	tokenA := args[3]
	amountADesiredf, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		return nil, errors.New("wrong amountADesired")
	}
	amountADesired, _ := types.NewAmount(amountADesiredf)
	amountAMinf, err := strconv.ParseFloat(args[5], 64)
	if err != nil {
		return nil, errors.New("wrong amountAMin")
	}
	amountAMin, _ := types.NewAmount(amountAMinf)
	tokenB := args[6]
	amountBDesiredf, err := strconv.ParseFloat(args[7], 64)
	if err != nil {
		return nil, errors.New("wrong amountBDesired")
	}
	amountBDesired, _ := types.NewAmount(amountBDesiredf)
	amountBMinf, err := strconv.ParseFloat(args[8], 64)
	if err != nil {
		return nil, errors.New("wrong amountBMin")
	}
	amountBMin, _ := types.NewAmount(amountBMinf)
	if len(args) > 10 {
		nonce, err = strconv.ParseUint(args[10], 10, 64)
		if err != nil {
			return nil, errors.New("wrong nonce")
		}
	}
	tx, err := transaction.NewPairAddLiquidity(Net, from, to, exchange, tokenA, tokenB, amountADesired, amountBDesired, amountAMin, amountBMin, nonce, "")
	if err != nil {
		return nil, err
	}
	return tx, nil
}

var RemoveLiquidityCmd = &cobra.Command{
	Use:     "RemoveLiquidity {from} {to} {exchange} {tokenA} {amountAmin} {tokenB} {amountBMin} {liquidity} {deadline} {password} {nonce}; Removal of liquidity;",
	Aliases: []string{"removeliquidity", "rl", "RL"},
	Short:   "RemoveLiquidity {from} {to} {exchange} {tokenA} {amountAmin} {tokenB} {amountBMin} {liquidity} {deadline} {password} {nonce}; Removal of liquidity;",
	Example: `
	RemoveLiquidity UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UWTfBGxDMZX19vjnacXVkP51min9EjhYq43W UWTXEqvUWik48uAHcJXZiyyWMy4GLtpGuttL 100  UBC  0.9 10 100 123456
		OR
	RemoveLiquidity UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UWTfBGxDMZX19vjnacXVkP51min9EjhYq43W UWTXEqvUWik48uAHcJXZiyyWMy4GLtpGuttL 100  UBC  0.9 10 100 123456 1
	`,
	Args: cobra.MinimumNArgs(9),
	Run:  RemoveLiquidity,
}

func RemoveLiquidity(cmd *cobra.Command, args []string) {
	var passwd []byte
	var err error
	if len(args) > 9 {
		passwd = []byte(args[9])
	} else {
		fmt.Println("please input password：")
		passwd, err = readPassWd()
		if err != nil {
			outputError(cmd.Use, fmt.Errorf("read password failed! %s", err.Error()))
			return
		}
	}
	privKey, err := ReadAddrPrivate(getAddJsonPath(args[0]), passwd)
	if err != nil {
		outputError(cmd.Use, fmt.Errorf("wrong password"))
		return
	}
	resp, err := GetAccountByRpc(args[0])
	if err != nil {
		outputError(cmd.Use, err)
		return
	}
	if resp.Code != 0 {
		outputRespError(cmd.Use, resp)
		return
	}
	var account *rpctypes.Account
	if err := json.Unmarshal(resp.Result, &account); err != nil {
		outputError(cmd.Use, err)
		return
	}

	tx, err := parseRLParams(args, account.Nonce+1)
	if err != nil {
		outputError(cmd.Use, err)
		return
	}

	if !signTx(cmd, tx, privKey.Private) {
		outputError(cmd.Use, errors.New("signature failure"))
		return
	}

	rs, err := sendTx(cmd, tx)
	if err != nil {
		outputError(cmd.Use, err)
	} else if rs.Code != 0 {
		outputRespError(cmd.Use, rs)
	} else {
		fmt.Println()
		fmt.Println(string(rs.Result))
	}
}

func parseRLParams(args []string, nonce uint64) (*types.Transaction, error) {
	var err error
	from := args[0]
	to := args[1]
	exchange := args[2]
	tokenA := args[3]

	amountAMinf, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		return nil, errors.New("wrong amountAMin")
	}
	amountAMin, _ := types.NewAmount(amountAMinf)
	tokenB := args[5]

	amountBMinf, err := strconv.ParseFloat(args[6], 64)
	if err != nil {
		return nil, errors.New("wrong amountBMin")
	}
	amountBMin, _ := types.NewAmount(amountBMinf)
	liquidityf, err := strconv.ParseFloat(args[7], 64)
	if err != nil {
		return nil, errors.New("wrong amountBMin")
	}
	fmt.Println(args[7])
	fmt.Println(liquidityf)
	liquidity, _ := types.NewAmount(liquidityf)
	deadline, err := strconv.ParseUint(args[8], 10, 64)
	if err != nil {
		return nil, errors.New("wrong amountBMin")
	}
	if len(args) > 10 {
		nonce, err = strconv.ParseUint(args[10], 10, 64)
		if err != nil {
			return nil, errors.New("wrong nonce")
		}
	}
	tx, err := transaction.NewPairRemoveLiquidity(Net, from, to, exchange, tokenA, tokenB, amountAMin, amountBMin, liquidity, deadline, nonce, "")
	if err != nil {
		return nil, err
	}
	return tx, nil
}

var SwapExactInCmd = &cobra.Command{
	Use:     "SwapExactIn {from} {to} {exchange} {tokenA} {tokenB} {amountIn} {amountOutMin} {deadline} {password} {nonce};Swap exact input tokens for tokens;",
	Aliases: []string{"swapexactin", "sei", "SEI"},
	Short:   "SwapExactIn {from} {to} {exchange} {tokenA} {tokenB} {amountIn} {amountOutMin} {deadline} {password} {nonce}; Swap exact input tokens for tokens;",
	Example: `
	SwapExactIn UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UWTfBGxDMZX19vjnacXVkP51min9EjhYq43W UWTXEqvUWik48uAHcJXZiyyWMy4GLtpGuttL UBC 100 1 100 123456
		OR
	SwapExactIn UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UWTfBGxDMZX19vjnacXVkP51min9EjhYq43W UWTXEqvUWik48uAHcJXZiyyWMy4GLtpGuttL UBC 100 1 100 123456 1
	`,
	Args: cobra.MinimumNArgs(8),
	Run:  SwapExactIn,
}

func SwapExactIn(cmd *cobra.Command, args []string) {
	var passwd []byte
	var err error
	if len(args) > 8 {
		passwd = []byte(args[8])
	} else {
		fmt.Println("please input password：")
		passwd, err = readPassWd()
		if err != nil {
			outputError(cmd.Use, fmt.Errorf("read password failed! %s", err.Error()))
			return
		}
	}
	privKey, err := ReadAddrPrivate(getAddJsonPath(args[0]), passwd)
	if err != nil {
		outputError(cmd.Use, fmt.Errorf("wrong password"))
		return
	}
	resp, err := GetAccountByRpc(args[0])
	if err != nil {
		outputError(cmd.Use, err)
		return
	}
	if resp.Code != 0 {
		outputRespError(cmd.Use, resp)
		return
	}
	var account *rpctypes.Account
	if err := json.Unmarshal(resp.Result, &account); err != nil {
		outputError(cmd.Use, err)
		return
	}

	tx, err := parseSEIParams(args, account.Nonce+1)
	if err != nil {
		outputError(cmd.Use, err)
		return
	}

	if !signTx(cmd, tx, privKey.Private) {
		outputError(cmd.Use, errors.New("signature failure"))
		return
	}

	rs, err := sendTx(cmd, tx)
	if err != nil {
		outputError(cmd.Use, err)
	} else if rs.Code != 0 {
		outputRespError(cmd.Use, rs)
	} else {
		fmt.Println()
		fmt.Println(string(rs.Result))
	}
}

func parseSEIParams(args []string, nonce uint64) (*types.Transaction, error) {
	var err error
	from := args[0]
	to := args[1]
	exchange := args[2]
	tokenA := args[3]
	tokenB := args[4]
	amountInf, err := strconv.ParseFloat(args[5], 64)
	if err != nil {
		return nil, errors.New("wrong amountIn")
	}
	amountIn, _ := types.NewAmount(amountInf)
	amountOutMinf, err := strconv.ParseFloat(args[6], 64)
	if err != nil {
		return nil, errors.New("wrong amountOutMin")
	}
	amountOutMin, _ := types.NewAmount(amountOutMinf)

	deadline, err := strconv.ParseUint(args[7], 10, 64)
	if err != nil {
		return nil, errors.New("wrong deadline")
	}
	if len(args) > 9 {
		nonce, err = strconv.ParseUint(args[9], 10, 64)
		if err != nil {
			return nil, errors.New("wrong nonce")
		}
	}
	pairs, err := GetAllPairByRpc(exchange)
	if err != nil {
		return nil, err
	}
	paths := transaction.CalculateShortestPath(hasharry.StringToAddress(tokenA), hasharry.StringToAddress(tokenB), pairs)
	if len(paths) == 0 {
		return nil, fmt.Errorf("not found")
	}
	tx, err := transaction.NewSwapExactIn(from, to, exchange, amountIn, amountOutMin, paths, deadline, nonce, "")
	if err != nil {
		return nil, err
	}
	return tx, nil
}

var SwapExactOutCmd = &cobra.Command{
	Use:     "SwapExactOut {from} {to} {exchange} {tokenA} {tokenB} {amountOut} {amountInMax} {deadline} {password} {nonce};Swap exact output tokens for tokens;",
	Aliases: []string{"swapexactout", "seo", "SEO"},
	Short:   "SwapExactOut {from} {to} {exchange} {tokenA} {tokenB} {amountOut} {amountInMax} {deadline} {password} {nonce}; Swap exact output tokens for tokens;",
	Example: `
	SwapExactOut UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UWTfBGxDMZX19vjnacXVkP51min9EjhYq43W UWTXEqvUWik48uAHcJXZiyyWMy4GLtpGuttL UBC 100 1 100 123456
		OR
	SwapExactOut UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UBCGLmQMfEeF6Fh8CGztrSktnHVpCxLiheYw UWTfBGxDMZX19vjnacXVkP51min9EjhYq43W UWTXEqvUWik48uAHcJXZiyyWMy4GLtpGuttL UBC 100 1 100 123456 1
	`,
	Args: cobra.MinimumNArgs(8),
	Run:  SwapExactOut,
}

func SwapExactOut(cmd *cobra.Command, args []string) {
	var passwd []byte
	var err error
	if len(args) > 8 {
		passwd = []byte(args[8])
	} else {
		fmt.Println("please input password：")
		passwd, err = readPassWd()
		if err != nil {
			outputError(cmd.Use+" err: ", fmt.Errorf("read password failed! %s", err.Error()))
			return
		}
	}
	privKey, err := ReadAddrPrivate(getAddJsonPath(args[0]), passwd)
	if err != nil {
		log.Error(cmd.Use+" err: ", fmt.Errorf("wrong password"))
		return
	}
	resp, err := GetAccountByRpc(args[0])
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	if resp.Code != 0 {
		outputRespError(cmd.Use, resp)
		return
	}
	var account *rpctypes.Account
	if err := json.Unmarshal(resp.Result, &account); err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}

	tx, err := parseSEOParams(args, account.Nonce+1)
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}

	if !signTx(cmd, tx, privKey.Private) {
		log.Error(cmd.Use+" err: ", errors.New("signature failure"))
		return
	}

	rs, err := sendTx(cmd, tx)
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
	} else if rs.Code != 0 {
		outputRespError(cmd.Use, rs)
	} else {
		fmt.Println()
		fmt.Println(string(rs.Result))
	}
}

func parseSEOParams(args []string, nonce uint64) (*types.Transaction, error) {
	var err error
	from := args[0]
	to := args[1]
	exchange := args[2]
	tokenA := args[3]
	tokenB := args[4]
	amountOutf, err := strconv.ParseFloat(args[5], 64)
	if err != nil {
		return nil, errors.New("wrong amountOut")
	}
	amountOut, _ := types.NewAmount(amountOutf)
	amountInMaxf, err := strconv.ParseFloat(args[6], 64)
	if err != nil {
		return nil, errors.New("wrong amountInMax")
	}
	amountInMax, _ := types.NewAmount(amountInMaxf)

	deadline, err := strconv.ParseUint(args[7], 10, 64)
	if err != nil {
		return nil, errors.New("wrong deadline")
	}
	if len(args) > 9 {
		nonce, err = strconv.ParseUint(args[9], 10, 64)
		if err != nil {
			return nil, errors.New("wrong nonce")
		}
	}
	pairs, err := GetAllPairByRpc(exchange)
	if err != nil {
		return nil, err
	}
	paths := transaction.CalculateShortestPath(hasharry.StringToAddress(tokenA), hasharry.StringToAddress(tokenB), pairs)
	if len(paths) == 0 {
		return nil, fmt.Errorf("not found")
	}
	tx, err := transaction.NewSwapExactOut(from, to, exchange, amountOut, amountInMax, paths, deadline, nonce, "")
	if err != nil {
		return nil, err
	}
	return tx, nil
}

var GetAllPairsCmd = &cobra.Command{
	Use:     "GetAllPairs {exchange};Get all pairs for exchange;",
	Aliases: []string{"getallpairs", "gap", "GAP"},
	Short:   "GetAllPairs {exchange}; Get all pairs for exchanges;",
	Example: `
	GetAllPairs UWTfBGxDMZX19vjnacXVkP51min9EjhYq43W
		OR
	GetAllPairs UWTfBGxDMZX19vjnacXVkP51min9EjhYq43W
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  GetAllPairs,
}

func GetAllPairs(cmd *cobra.Command, args []string) {
	client, err := NewRpcClient()
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()

	resp, err := client.Gc.GetExchangePairs(ctx, &rpc.Address{Address: args[0]})
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	if resp.Code == 0 {
		output(string(resp.Result))
		return
	}
	outputRespError(cmd.Use, resp)
}

func GetAllPairByRpc(addr string) ([]*types.RpcPair, error) {
	client, err := NewRpcClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()
	rs, err := client.Gc.GetExchangePairs(ctx, &rpc.Address{Address: addr})
	if err != nil {
		return nil, err
	}
	if rs.Code != rpctypes.RpcSuccess {
		return nil, errors.New(rs.Err)
	}
	pairs := make([]*types.RpcPair, 0)
	if err := json.Unmarshal(rs.Result, &pairs); err != nil {
		return nil, err
	}
	return pairs, nil
}
