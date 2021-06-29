package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jhdriver/UBChain/common/hasharry"
	"github.com/jhdriver/UBChain/core/types"
	"github.com/jhdriver/UBChain/crypto/ecc/secp256k1"
	"github.com/jhdriver/UBChain/rpc"
	"github.com/jhdriver/UBChain/rpc/rpctypes"
	"github.com/jhdriver/UBChain/ut/transaction"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
	"time"
)

func init() {
	txCmds := []*cobra.Command{
		GetTransactionCmd,
		SendTransactionCmd,
	}
	RootCmd.AddCommand(txCmds...)
	RootSubCmdGroups["transaction"] = txCmds

}

var SendTransactionCmd = &cobra.Command{
	Use:     "SendTransaction {from} {to:amount|to:amount} {contract} {note} {password} {nonce}; Send a transaction ;",
	Aliases: []string{"SendTransaction", "st", "ST"},
	Short:   "SendTransaction {from} {to:amount|to:amount} {contract} {note} {password} {nonce}; Send a transaction ;",
	Example: `
	SendTransaction 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ "3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ:10|3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndD:10" UBC  "transaction note"
		OR
	SendTransaction 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ "3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ:10|3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndD:10" UBC  "transaction note" 123456
		OR
	SendTransaction 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ "3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ:10|3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndD:10" UBC  "transaction note" 123456 1
	`,
	Args: cobra.MinimumNArgs(4),
	Run:  SendTransaction,
}

func SendTransaction(cmd *cobra.Command, args []string) {
	var passwd []byte
	var err error
	if len(args) > 4 {
		passwd = []byte(args[4])
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

	tx, err := parseV2Params(args)
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

func parseV2Params(args []string) (*types.Transaction, error) {
	var err error
	var nonce uint64
	from := hasharry.StringToAddress(args[0])
	to, err := parseReceiver(args[1])
	if err != nil {
		return nil, err
	}
	contract := hasharry.StringToAddress(args[2])
	note := args[3]
	if len(args) > 5 {
		nonce, err = strconv.ParseUint(args[5], 10, 64)
		if err != nil {
			return nil, errors.New("wrong nonce")
		}
	}
	return transaction.NewTransactionV2(from.String(), to, contract.String(), note, nonce), nil
}

func parseReceiver(toStr string) ([]map[string]uint64, error) {
	toList := []map[string]uint64{}
	receivers := strings.Split(toStr, "|")
	if len(receivers) == 0 {
		return nil, fmt.Errorf("no receiver")
	}
	for _, receiver := range receivers {
		strs := strings.Split(receiver, ":")
		if len(strs) != 2 {
			return nil, fmt.Errorf("wrong receiver %s", receiver)
		}
		fAmt, err := strconv.ParseFloat(strs[1], 64)
		if err != nil {
			return nil, fmt.Errorf("wrong receiver %s", receiver)
		}
		if amt, err := types.NewAmount(fAmt); err != nil {
			return nil, fmt.Errorf("wrong receiver %s", receiver)
		} else {
			toList = append(toList, map[string]uint64{strs[0]: amt})
		}
	}
	return toList, nil
}

func signTx(cmd *cobra.Command, tx *types.Transaction, key string) bool {
	tx.SetHash()
	priv, err := secp256k1.ParseStringToPrivate(key)
	if err != nil {
		outputError(cmd.Use, errors.New("[key] wrong"))
		return false
	}
	if err := tx.SignTx(priv); err != nil {
		outputError(cmd.Use, errors.New("sign failed"))
		return false
	}
	return true
}
func signTx1(tx *types.Transaction, key string) bool {
	tx.SetHash()
	priv, err := secp256k1.ParseStringToPrivate(key)
	if err != nil {
		return false
	}
	if err := tx.SignTx(priv); err != nil {
		return false
	}
	return true
}
func sendTx(cmd *cobra.Command, tx *types.Transaction) (*rpc.Response, error) {
	rpcTx, err := types.TranslateTxToRpcTx(tx)
	if err != nil {
		return nil, err
	}
	rpcClient, err := NewRpcClient()
	if err != nil {
		return nil, err
	}
	defer rpcClient.Close()

	jsonBytes, err := json.Marshal(rpcTx)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(jsonBytes))
	re := &rpc.Bytes{Bytes: jsonBytes}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()

	resp, err := rpcClient.Gc.SendTransaction(ctx, re)
	if err != nil {
		return nil, err
	}
	return resp, nil

}
func sendTx1(tx *types.Transaction) (*rpc.Response, error) {
	rpcTx, err := types.TranslateTxToRpcTx(tx)
	if err != nil {
		return nil, err
	}
	rpcClient, err := NewRpcClient()
	if err != nil {
		return nil, err
	}
	defer rpcClient.Close()

	jsonBytes, err := json.Marshal(rpcTx)
	if err != nil {
		return nil, err
	}
	re := &rpc.Bytes{Bytes: jsonBytes}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()

	resp, err := rpcClient.Gc.SendTransaction(ctx, re)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

var GetTransactionCmd = &cobra.Command{
	Use:     "GetTransaction {txhash}; Get Transaction by hash;",
	Aliases: []string{"gettransaction", "gt", "GT"},
	Short:   "GetTransaction {txhash}; Get Transaction by hash;",
	Example: `
	GetTransaction 0xef7b92e552dca02c97c9d596d1bf69d0044d95dec4cee0e6a20153e62bce893b
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  GetTransaction,
}

func GetTransaction(cmd *cobra.Command, args []string) {
	resp, err := GetTransactionRpc(args[0])
	if err != nil {
		outputError(cmd.Use, err)
		return
	}
	if resp.Code == 0 {
		output(string(resp.Result))
		return
	}
	outputRespError(cmd.Use, resp)
}

func GetTransactionRpc(hashStr string) (*rpc.Response, error) {
	client, err := NewRpcClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()
	h := &rpc.Hash{Hash: hashStr}
	resp, err := client.Gc.GetTransaction(ctx, h)
	return resp, err
}

func SendTransactionRpc(tx string) (*rpc.Response, error) {

	rpcClient, err := NewRpcClient()
	if err != nil {
		return nil, err
	}
	defer rpcClient.Close()

	re := &rpc.Bytes{Bytes: []byte(tx)}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()

	resp, err := rpcClient.Gc.SendTransaction(ctx, re)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
