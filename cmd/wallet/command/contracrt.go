package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/core/types"
	"github.com/UBChainNet/UBChain/rpc"
	"github.com/UBChainNet/UBChain/rpc/rpctypes"
	"github.com/UBChainNet/UBChain/ut"
	"github.com/UBChainNet/UBChain/ut/transaction"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

func init() {
	contractCmds := []*cobra.Command{
		GetContractCmd,
		SendContractCmd,
	}
	RootCmd.AddCommand(contractCmds...)
	RootSubCmdGroups["contract"] = contractCmds

}

var SendContractCmd = &cobra.Command{
	Use:     "SendContract {from} {to} {name} {abbr} {Increase} {description} {amount} {note} {password} {nonce}; Send contract of coin publish;",
	Aliases: []string{"sendcontract", "SC", "sc"},
	Short:   "SendContract {from} {to} {name} {abbr} {Increase} {description} {amount} {note} {password} {nonce}; Send contract of coin publish;",
	Example: `
	SendContract 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE "Test Coin" TC false "description" 1000  "transaction note"
		OR
	SendContract 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE "Test Coin" TC false "description" 1000  "transaction note" 123456
		OR
	SendContract 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE "Test Coin" TC false "description" 1000  "transaction note" 123456 0
	`,
	Args: cobra.MinimumNArgs(8),
	Run:  SendContract,
}

func SendContract(cmd *cobra.Command, args []string) {
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

	tx, err := parseSCParams(args)
	if err != nil {
		outputError(cmd.Use, err)
		return
	}
	resp, err := GetAccountByRpc(tx.From().String())
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
	if tx.TxHead.Nonce == 0 {
		tx.TxHead.Nonce = account.Nonce + 1
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

func parseSCParams(args []string) (*types.Transaction, error) {
	var err error
	var amount, nonce uint64
	from := hasharry.StringToAddress(args[0])
	to := hasharry.StringToAddress(args[1])
	name := args[2]
	abbr := args[3]
	increase, err := strconv.ParseBool(args[4])
	description := args[5]
	if err != nil {
		return nil, fmt.Errorf("wrong increase, %s", err.Error())
	}
	if fAmount, err := strconv.ParseFloat(args[6], 64); err != nil {
		return nil, errors.New("wrong amount")
	} else {
		if fAmount < 0 {
			return nil, errors.New("wrong amount")
		}
		if amount, err = types.NewAmount(fAmount); err != nil {
			return nil, errors.New("wrong amount")
		}
	}
	if err := ut.CheckSymbol(abbr); err != nil {
		return nil, err
	}
	contract, err := ut.GenerateContractAddress(Net, abbr)
	if err != nil {
		return nil, err
	}
	fmt.Println("\ncontract address is ", contract)

	note := args[7]
	if len(args) > 9 {
		nonce, err = strconv.ParseUint(args[9], 10, 64)
		if err != nil {
			return nil, errors.New("wrong nonce")
		}
	}
	tx := transaction.NewContract(from.String(), to.String(), contract, note, amount, nonce, name, abbr, increase, description)
	return tx, nil
}

var GetContractCmd = &cobra.Command{
	Use:     "GetContract {contract address}; Get a contract;",
	Aliases: []string{"getcontract", "gc", "GC"},
	Short:   "GetContract {contract address}; Get a contract;",
	Example: `
	GetContract 2KwjygFUZ8oWbWAzY7mT5tvpHC8ohtG9h3h3xjxmtqYD
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  GetContract,
}

func GetContract(cmd *cobra.Command, args []string) {
	resp, err := GetContractByRpc(args[0])
	if err != nil {
		outputError(cmd.Use, err)
		return
	}
	if resp.Code == 0 {
		output(string(resp.Result))
		return
	} else {
		outputRespError(cmd.Use, resp)
	}
}

func GetContractByRpc(contractAddr string) (*rpc.Response, error) {
	client, err := NewRpcClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()
	return client.Gc.GetContract(ctx, &rpc.Address{Address: contractAddr})
}
