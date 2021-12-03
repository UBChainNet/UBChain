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
		TokenHubAckCmd,
		TokenHubTransferOutCmd,
		TokenHubTransferInCmd,
		TokenHubFinishAcrossCmd,
	}
	RootCmd.AddCommand(tokenHubCmds...)
	RootSubCmdGroups["tokenhub"] = tokenHubCmds

}

var TokenHubInitCmd = &cobra.Command{
	Use:     "TokenHubInit {from} {contract} {setter} {admin} {feeTo} {feeRate} {password} {nonce}; Init tokenhub contract;",
	Aliases: []string{"TokenHubInit", "tokenhubinit", "thi", "THI"},
	Short:   "TokenHubInit {from} {contract} {setter} {admin} {feeTo} {feeRate} {password} {nonce}; Init tokenhub contract;",
	Example: `
	TokenHubInit 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE "0.001" 123456
		OR
	TokenHubInit 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE "0.001" 123456 1
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
	feeRate := args[5]

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
	Use:     "TokenHubAck {from} {contract} {sequence:ackType:hash|sequence:ackType:hash} {password} {nonce}; tokenhub ack transfer;",
	Aliases: []string{"TokenHubAck", "tokenhuback", "tha", "THA"},
	Short:   "TokenHubAck {from} {contract} {sequence:ackType:hash|sequence:ackType:hash} {password} {nonce}; tokenhub ack transfer;",
	Example: `
	TokenHubAck 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE "1:3:0x78427897dd3aa5116953b00ffa3445a273fec75b3864fd2e73766ea59a757bbf" 123456
		OR
	TokenHubAck 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE "1:3:0x78427897dd3aa5116953b00ffa3445a273fec75b3864fd2e73766ea59a757bbf" 123456 1
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
	list := strings.Split(args[2], "|")
	sequences := make([]uint64, 0)
	ackTypes := make([]uint8, 0)
	hashes := make([]string, 0)
	for _, sequenceAndType := range list {
		strs := strings.Split(sequenceAndType, ":")
		if len(strs) != 3 {
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
		hashes = append(hashes, strs[2])
	}
	if len(args) > 4 {
		nonce, err = strconv.ParseUint(args[4], 10, 64)
		if err != nil {
			return nil, errors.New("wrong nonce")
		}
	}
	tx, err := transaction.NewTokenHubAck(from, contract, sequences, ackTypes, hashes, nonce, "")
	if err != nil {
		return nil, err
	}
	return tx, nil
}


var TokenHubTransferOutCmd = &cobra.Command{
	Use:     "TokenHubTransferOut {from} {contract} {to} {amount} {password} {nonce}; transfer by tokenhub;",
	Aliases: []string{"TokenHubTransferOut", "tokenhubtransferout", "thto", "THTO"},
	Short:   "TokenHubTransferOut {from} {contract} {to} {amount} {password} {nonce}; transfer by tokenhub;",
	Example: `
	TokenHubTransferOut 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 0x8a26d5a655aafd799fdf49c8bbe546e50fce39d1 100 123456
		OR
	TokenHubTransferOut 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 0x8a26d5a655aafd799fdf49c8bbe546e50fce39d1 100 123456 1
	`,
	Args: cobra.MinimumNArgs(4),
	Run:  TokenHubTransferOut,
}

func TokenHubTransferOut(cmd *cobra.Command, args []string) {
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

	tx, err := parseTHTParams(args, account.Nonce+1)
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

func parseTHTParams(args []string, nonce uint64) (*types.Transaction, error) {
	var err error
	from := args[0]
	contract := args[1]
	to := args[2]
	amountf, err := strconv.ParseFloat(args[3], 64)
	if err != nil{
		return nil, err
	}
	amount, _ := types.NewAmount(amountf)
	if len(args) > 5 {
		nonce, err = strconv.ParseUint(args[5], 10, 64)
		if err != nil {
			return nil, errors.New("wrong nonce")
		}
	}
	tx, err := transaction.NewTokenHubTransferOut(from, contract, to, amount, nonce, "")
	if err != nil {
		return nil, err
	}
	return tx, nil
}


var TokenHubTransferInCmd = &cobra.Command{
	Use:     "TokenHubTransferIn {from} {contract} {to} {amount} {across seq} {password} {nonce}; transfer in by tokenhub;",
	Aliases: []string{"TokenHubTransferIn", "tokenhubtransferin", "thti", "THTI"},
	Short:   "TokenHubTransferIn {from} {contract} {to} {amount} {across seq} {password} {nonce}; transfer in by tokenhub;",
	Example: `
	TokenHubTransferIn 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 100 100 123456
		OR
	TokenHubTransferIn 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE 100 100 123456 1
	`,
	Args: cobra.MinimumNArgs(5),
	Run:  TokenHubTransferIn,
}

func TokenHubTransferIn(cmd *cobra.Command, args []string) {
	var passwd []byte
	var err error
	if len(args) > 5 {
		passwd = []byte(args[5])
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

	tx, err := parseTHTIParams(args, account.Nonce+1)
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

func parseTHTIParams(args []string, nonce uint64) (*types.Transaction, error) {
	var err error
	from := args[0]
	contract := args[1]
	to := args[2]
	amountf, err := strconv.ParseFloat(args[3], 64)
	if err != nil{
		return nil, err
	}
	amount, _ := types.NewAmount(amountf)
	acrossSeq, err := strconv.ParseUint(args[4], 10, 64)
	if err != nil{
		return nil, err
	}
	if len(args) > 6 {
		nonce, err = strconv.ParseUint(args[6], 10, 64)
		if err != nil {
			return nil, errors.New("wrong nonce")
		}
	}
	tx, err := transaction.NewTokenHubTransferIn(from, contract, to, amount, acrossSeq, nonce, "")
	if err != nil {
		return nil, err
	}
	return tx, nil
}


var TokenHubFinishAcrossCmd = &cobra.Command{
	Use:     "TokenHubFinishAcross {from} {contract} {acrossSeq|acrossSeq}; finish across;",
	Aliases: []string{"TokenHubFinishAcross", "THFA", "thfa"},
	Short:   "TokenHubFinishAcross {from} {contract} {acrossSeq|acrossSeq}; finish across;",
	Example: `
	TokenHubFinishAcross 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE "1|2" 123456
		OR
	TokenHubFinishAcross 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE "1|2" 123456 1
	`,
	Args: cobra.MinimumNArgs(3),
	Run:  TokenHubFinishAcross,
}

func TokenHubFinishAcross(cmd *cobra.Command, args []string) {
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

	tx, err := parseTHFAParams(args, account.Nonce+1)
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

func parseTHFAParams(args []string, nonce uint64) (*types.Transaction, error) {
	var err error
	from := args[0]
	contract := args[1]
	str := args[2]
	seqs := strings.Split(str, "|")
	acrossSeqs := make([]uint64, 0)
	for _, seq := range seqs{
		acrossSeq, err := strconv.ParseUint(seq, 10, 64)
		if err != nil{
			return nil, err
		}
		acrossSeqs = append(acrossSeqs, acrossSeq)
	}

	if len(args) > 4 {
		nonce, err = strconv.ParseUint(args[4], 10, 64)
		if err != nil {
			return nil, errors.New("wrong nonce")
		}
	}
	tx, err := transaction.NewTokenHubFinishAcross(from, contract, acrossSeqs, nonce, "")
	if err != nil {
		return nil, err
	}
	return tx, nil
}
