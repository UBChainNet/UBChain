package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/UBChainNet/UBChain/core/runner/tokenhub_runner"
	"github.com/UBChainNet/UBChain/core/types"
	"github.com/UBChainNet/UBChain/rpc/rpctypes"
	"github.com/UBChainNet/UBChain/ut/transaction"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

func init() {
	tokenHubCmds := []*cobra.Command{
		TokenHubInitCmd,
	}
	RootCmd.AddCommand(tokenHubCmds...)
	RootSubCmdGroups["tokenhub"] = tokenHubCmds

}

var TokenHubInitCmd = &cobra.Command{
	Use:     "TokenHubInit {from} {contract} {setter} {admin} {feeTo} {feeRate} {password} {nonce}; Init tokenhub contract;",
	Aliases: []string{"TokenHubInit", "tokenhubinit", "thi", "THI"},
	Short:   "TokenHubInit {from} {contract} {setter} {admin} {feeTo} {feeRate} {password} {nonce}; Init tokenhub contract;",
	Example: `
	TokenHubInit 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 0.001 123456
		OR
	TokenHubInit 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 0.001 123456 1
	`,
	Args: cobra.MinimumNArgs(6),
	Run:  TokenHubInit,
}

func TokenHubInit(cmd *cobra.Command, args []string) {
	var passwd []byte
	var err error
	if len(args) > 6 {
		passwd = []byte(args[6])
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

	tx, err := parseTHIParams(args, account.Nonce+1)
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

func parseTHIParams(args []string, nonce uint64) (*types.Transaction, error) {
	var err error
	from := args[0]
	contract := args[1]
	setter := args[2]
	admin := args[3]
	feeTo := args[4]
	feeRatef := args[5]
	feeRate, err := strconv.ParseFloat(feeRatef, 64)
	if err != nil {
		return nil, err
	}

	if len(args) > 7 {
		nonce, err = strconv.ParseUint(args[7], 10, 64)
		if err != nil {
			return nil, errors.New("wrong nonce")
		}
	}
	if contract == "" {
		contract, _ = tokenhub_runner.TokenHubAddress(Net, from, nonce)
	}
	tx, err := transaction.NewTokenHubInit(from, contract, setter, admin, feeTo, feeRate, nonce, "")
	if err != nil {
		return nil, err
	}
	return tx, nil
}

var TokenHubAckCmd = &cobra.Command{
	Use:     "TokenHubAck {from} {contract} {sequence:ackType|sequence:ackType} {password} {nonce}; tokenhub ack transfer;",
	Aliases: []string{"TokenHubAck", "tokenhuback", "tha", "THA"},
	Short:   "TokenHubAck {from} {contract} {sequence:ackType|sequence:ackType} {password} {nonce}; tokenhub ack transfer;",
	Example: `
	TokenHubAck 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE "1:3|2:2|3:1" 123456
		OR
	TokenHubAck 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE "1:3|2:2|3:1" 123456 1
	`,
	Args: cobra.MinimumNArgs(3),
	Run:  TokenHubAck,
}

func TokenHubAck(cmd *cobra.Command, args []string) {
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

	tx, err := parseTHAParams(args, account.Nonce+1)
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

func parseTHAParams(args []string, nonce uint64) (*types.Transaction, error) {
	var err error
	from := args[0]
	contract := args[1]
	list := strings.Split(args[1], "|")
	sequences := make([]uint64, 0)
	ackTypes := make([]uint8, 0)
	for _, sequenceAndType := range list {
		strs := strings.Split(sequenceAndType, ":")
		if len(strs) != 2 {
			return nil, fmt.Errorf("wrong sequence and ackType")
		}
		sequence, err := strconv.ParseUint(strs[0], 10, 64)
		if err != nil {
			return nil, err
		}
		ackType, err := strconv.ParseUint(strs[1], 10, 64)
		if err != nil {
			return nil, err
		}
		sequences = append(sequences, sequence)
		ackTypes = append(ackTypes, uint8(ackType))
	}
	if len(args) > 4 {
		nonce, err = strconv.ParseUint(args[4], 10, 64)
		if err != nil {
			return nil, errors.New("wrong nonce")
		}
	}
	tx, err := transaction.NewTokenHubAck(from, contract, sequences, ackTypes, nonce, "")
	if err != nil {
		return nil, err
	}
	return tx, nil
}
