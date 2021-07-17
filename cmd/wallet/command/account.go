package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/UBChainNet/UBChain/common/keystore"
	"github.com/UBChainNet/UBChain/crypto/ecc/secp256k1"
	"github.com/UBChainNet/UBChain/p2p"
	"github.com/UBChainNet/UBChain/rpc"
	"github.com/UBChainNet/UBChain/rpc/rpctypes"
	"github.com/UBChainNet/UBChain/ut"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func init() {
	accountCmds := []*cobra.Command{
		CreateAccountCmd,
		GetAccountCmd,
		ShowAccountCmd,
		DecryptAccountCmd,
		MnemonicToAccountCmd,
		EcToAccountCmd,
		P2PIDCmd,
	}

	RootCmd.AddCommand(accountCmds...)
	RootSubCmdGroups["account"] = accountCmds
}

var GetAccountCmd = &cobra.Command{
	Use:     "GetAccount {address};Get account status;",
	Aliases: []string{"getaccount", "ga", "GA"},
	Short:   "GetAccount {address};Get account status;",
	Example: `
	GetAccount 23zE69fmaqK2LCHQrMQifTASSF1U 
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  GetAccount,
}

func GetAccount(cmd *cobra.Command, args []string) {
	resp, err := GetAccountByRpc(args[0])
	if err != nil {
		outputError(cmd.Use, err)
		return
	}
	if resp.Code == 0 {
		account := &rpctypes.Account{}
		json.Unmarshal(resp.Result, account)
		if account.Address != args[0] {
			account.Address = args[0]
		}
		bytes, _ := json.Marshal(account)
		output(string(bytes))
		return
	} else {
		outputRespError(cmd.Use, resp)
	}
}
func GetAccountRpc(addr string) (string, error) {
	resp, err := GetAccountByRpc(addr)
	if err != nil {
		return "", err
	}
	if resp.Code == 0 {
		account := &rpctypes.Account{}
		json.Unmarshal(resp.Result, account)
		if account.Address != addr {
			account.Address = addr
		}
		bytes, _ := json.Marshal(account)

		return string(bytes), nil
	} else {
		return resp.Err, nil
	}
}
func GetAccountByRpc(addr string) (*rpc.Response, error) {
	client, err := NewRpcClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()
	return client.Gc.GetAccount(ctx, &rpc.Address{Address: addr})

}

var P2PIDCmd = &cobra.Command{
	Use:     "P2PID {address};Get account P2P ID;",
	Aliases: []string{"P2P", "p2p"},
	Short:   "P2PID {address};Get account P2P ID;",
	Example: `
	P2PID 23zE69fmaqK2LCHQrMQifTASSF1U 
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  P2PID,
}

func P2PID(cmd *cobra.Command, args []string) {
	var passWd []byte
	var err error
	if len(args) == 2 && args[1] != "" {
		passWd = []byte(args[1])
	} else {
		fmt.Println("please set account password, cannot exceed 32 bytes：")
		passWd, err = readPassWd()
		if err != nil {
			outputError(cmd.Use, fmt.Errorf("read password failed! %s", err.Error()))
			return
		}
	}
	if len(passWd) > 32 {
		outputError(cmd.Use, fmt.Errorf("password too long! "))
		return
	}

	keyFile := getAddJsonPath(args[0])

	privKey, err := ReadAddrPrivate(keyFile, passWd)
	if err != nil {
		outputError(cmd.Use, fmt.Errorf("wrong password"))
		return
	}
	priv, err := secp256k1.ParseStringToPrivate(privKey.Private)
	if err != nil {
		outputError(cmd.Use, fmt.Errorf("wrong private key"))
		return
	}
	p2pId, err := p2p.GenerateP2pId(priv)
	if err != nil {
		outputError(cmd.Use, err)
		return
	}
	fmt.Println(p2pId.String())
}

var CreateAccountCmd = &cobra.Command{
	Use:     "CreateAccount {password}",
	Short:   "CreateAccount {password}; Create account;",
	Aliases: []string{"createaccount", "CA", "ca"},
	Example: `
	CreateAccount  
		OR
	CreateAccount 123456
	`,
	Args: cobra.MinimumNArgs(0),
	Run:  CreateAccount,
}

func CreateAccount(cmd *cobra.Command, args []string) {
	var passWd []byte
	var err error
	if len(args) == 1 && args[0] != "" {
		passWd = []byte(args[0])
	} else {
		fmt.Println("please set account password, cannot exceed 32 bytes：")
		passWd, err = readPassWd()
		if err != nil {
			outputError(cmd.Use, fmt.Errorf("read password failed! %s", err.Error()))
			return
		}
	}
	if len(passWd) > 32 {
		outputError(cmd.Use, fmt.Errorf("password too long! "))
		return
	}
	entropy, err := ut.Entropy()
	if err != nil {
		outputError(cmd.Use, err)
		return
	}
	mnemonicStr, err := ut.Mnemonic(entropy)
	if err != nil {
		outputError(cmd.Use, err)
		return
	}
	key, err := ut.MnemonicToEc(mnemonicStr)
	if err != nil {
		outputError(cmd.Use, fmt.Errorf("generate secp256k1 key failed! %s", err.Error()))
		return
	}
	p2pId, err := p2p.GenerateP2pId(key)
	if err != nil {
		outputError(cmd.Use, fmt.Errorf("generate p2p id failed! %s", err.Error()))
	}
	if j, err := keystore.GenerateKeyJson(Net, Cfg.KeyStoreDir, key, mnemonicStr, passWd); err != nil {
		outputError(cmd.Use, fmt.Errorf("generate key failed! %s", err.Error()))
	} else {
		j.P2pId = p2pId.String()
		bytes, _ := json.Marshal(j)
		output(string(bytes))
	}
}

func CreateAccountRpc(pwd string) (*keystore.Json, error) {
	var passWd []byte
	var err error

	passWd = []byte(pwd)
	if len(passWd) > 32 {

		return nil, fmt.Errorf("password too long! ")
	}
	entropy, err := ut.Entropy()
	if err != nil {
		return nil, err
	}
	mnemonicStr, err := ut.Mnemonic(entropy)
	if err != nil {
		return nil, err
	}
	key, err := ut.MnemonicToEc(mnemonicStr)
	if err != nil {
		return nil, fmt.Errorf("generate secp256k1 key failed! %s", err.Error())
	}
	p2pId, err := p2p.GenerateP2pId(key)
	if err != nil {
		return nil, fmt.Errorf("generate p2p id failed! %s", err.Error())
	}
	if j, err := keystore.GenerateKeyJson(Net, Cfg.KeyStoreDir, key, mnemonicStr, passWd); err != nil {
		return nil, fmt.Errorf("generate key failed! %s", err.Error())
	} else {
		j.P2pId = p2pId.String()
		return j, nil
	}
}

func readPassWd() ([]byte, error) {
	var passWd [33]byte

	n, err := os.Stdin.Read(passWd[:])
	if err != nil {
		return nil, err
	}
	if n <= 1 {
		return nil, errors.New("not read")
	}
	return passWd[:n-1], nil
}

var ShowAccountCmd = &cobra.Command{
	Use:     "ShowAccounts",
	Short:   "ShowAccounts; Show all account of the wallet;",
	Aliases: []string{"showaccounts", "sa", "SA"},
	Example: `
	ShowAccounts
	`,
	Args: cobra.MinimumNArgs(0),
	Run:  ShowAccount,
}

func ShowAccount(cmd *cobra.Command, args []string) {
	if addrList, err := keystore.ReadAllAccount(Cfg.KeyStoreDir); err != nil {
		outputError(cmd.Use, fmt.Errorf("read account failed! %s", err.Error()))
	} else {
		bytes, _ := json.Marshal(addrList)
		output(string(bytes))
	}
}

var DecryptAccountCmd = &cobra.Command{
	Use:     "DecryptAccount {address} {password} {key file}；Decrypting account json file generates the private key and mnemonic;；",
	Short:   "DecryptAccount {address} {password} {key file}; Decrypting account json file generates the private key and mnemonic;",
	Aliases: []string{"decryptaccount", "DA", "da"},

	Example: `
	DecryptAccount 3ajKPvYpncZ8YtmCXogJFkKSQJb2FeXYceBf
		OR
	DecryptAccount 3ajKPvYpncZ8YtmCXogJFkKSQJb2FeXYceBf 123456
		OR
	DecryptAccount 3ajKPvYpncZ8YtmCXogJFkKSQJb2FeXYceBf 123456 3ajKPvYpncZ8YtmCXogJFkKSQJb2FeXYceBf.json
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  DecryptAccount,
}

func DecryptAccount(cmd *cobra.Command, args []string) {
	var passWd []byte
	var keyFile string
	var err error
	if len(args) >= 2 && args[1] != "" {
		passWd = []byte(args[1])
	} else {
		fmt.Println("please input password：")
		passWd, err = readPassWd()
		if err != nil {
			outputError(cmd.Use, fmt.Errorf("read password failed! %s", err.Error()))
			return
		}
	}
	if len(args) == 3 && args[2] != "" {
		keyFile = args[2]
	} else {
		keyFile = getAddJsonPath(args[0])
	}

	privKey, err := ReadAddrPrivate(keyFile, passWd)
	if err != nil {
		outputError(cmd.Use, fmt.Errorf("wrong password"))
		return
	}

	bytes, _ := json.Marshal(privKey)
	output(string(bytes))
}

func DecryptAccountRpc(address string, pwd string) (*keystore.Private, error) {
	var passWd []byte
	var keyFile string
	var err error
	passWd = []byte(pwd)
	keyFile = getAddJsonPath(address)
	privKey, err := ReadAddrPrivate(keyFile, passWd)
	if err != nil {
		return nil, fmt.Errorf("wrong password")
	}
	return privKey, nil
}

var MnemonicToAccountCmd = &cobra.Command{
	Use:     "MnemonicToAccount {mnemonic} {password}；Restore address by mnemonic and set new password;",
	Short:   "MnemonicToAccount {mnemonic} {password}; Restore address by mnemonic and set new password;",
	Aliases: []string{"mnemonictoaccount", "MTA", "mta"},
	Example: `
	MnemonicToAccount "sadness ladder sister camp suspect sting height diagram confirm program twist ostrich blush bronze pass gasp resist random nothing recycle husband install business turtle"
		OR
	MnemonicToAccount "sadness ladder sister camp suspect sting height diagram confirm program twist ostrich blush bronze pass gasp resist random nothing recycle husband install business turtle" 123456
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  MnemonicToAccount,
}

func MnemonicToAccount(cmd *cobra.Command, args []string) {
	var passWd []byte
	var err error
	priv, err := ut.MnemonicToEc(args[0])
	if err != nil {
		outputError(cmd.Use, errors.New("[mnemonic] wrong"))
		return
	}
	if len(args) == 2 && args[1] != "" {
		passWd = []byte(args[1])
	} else {
		fmt.Println("please set address password, cannot exceed 32 bytes：")
		passWd, err = readPassWd()
		if err != nil {
			outputError(cmd.Use, fmt.Errorf("read pass word failed! %s", err.Error()))
			return
		}
	}
	if len(passWd) > 32 {
		outputError(cmd.Use, fmt.Errorf("password too long! "))
		return
	}
	p2pId, err := p2p.GenerateP2pId(priv)
	if err != nil {
		outputError(cmd.Use, fmt.Errorf("generate p2p id failed! %s", err.Error()))
	}
	if j, err := keystore.GenerateKeyJson(Net, Cfg.KeyStoreDir, priv, args[0], passWd); err != nil {
		outputError(cmd.Use, fmt.Errorf("generate key failed! %s", err.Error()))
	} else {
		j.P2pId = p2pId.String()
		bytes, _ := json.Marshal(j)
		output(string(bytes))
	}
}

var EcToAccountCmd = &cobra.Command{
	Use:     "EcToAccount {mnemonic} {password}；Restore address by private and set new password;",
	Short:   "EcToAccount {mnemonic} {password}; Restore address by private and set new password;",
	Aliases: []string{"ectoaccount", "ETA", "eta"},
	Example: `
	EcToAccount "4c2cee98b562b2a63fb76b416768bf6052fc177cb9cadc55e4021eeac9bb26d0"
		OR
	EcToAccount "4c2cee98b562b2a63fb76b416768bf6052fc177cb9cadc55e4021eeac9bb26d0" 123456
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  EcToAccount,
}

func EcToAccount(cmd *cobra.Command, args []string) {
	var passWd []byte
	var err error
	priv, err := secp256k1.ParseStringToPrivate(args[0])
	if err != nil {
		outputError(cmd.Use, errors.New("[priavte] wrong"))
		return
	}
	if len(args) == 2 && args[1] != "" {
		passWd = []byte(args[1])
	} else {
		fmt.Println("please set address password, cannot exceed 32 bytes：")
		passWd, err = readPassWd()
		if err != nil {
			outputError(cmd.Use, fmt.Errorf("read pass word failed! %s", err.Error()))
			return
		}
	}
	if len(passWd) > 32 {
		outputError(cmd.Use, fmt.Errorf("password too long! "))
		return
	}
	p2pId, err := p2p.GenerateP2pId(priv)
	if err != nil {
		outputError(cmd.Use, fmt.Errorf("generate p2p id failed! %s", err.Error()))
	}
	if j, err := keystore.GenerateKeyJson(Net, Cfg.KeyStoreDir, priv, "", passWd); err != nil {
		outputError(cmd.Use, fmt.Errorf("generate key failed! %s", err.Error()))
	} else {
		j.P2pId = p2pId.String()
		bytes, _ := json.Marshal(j)
		output(string(bytes))
	}
}

func EcToAccountRpc(passWd []byte, private string) (*keystore.Json, error) {
	var err error
	priv, err := secp256k1.ParseStringToPrivate(private)
	if err != nil {
		return nil, errors.New("[priavte] wrong")
	}
	if len(passWd) > 32 {
		return nil, fmt.Errorf("password too long! ")
	}
	p2pId, err := p2p.GenerateP2pId(priv)
	if err != nil {
		return nil, fmt.Errorf("generate p2p id failed! %s", err.Error())
	}
	if j, err := keystore.GenerateKeyJson(Net, Cfg.KeyStoreDir, priv, "", passWd); err != nil {
		return nil, fmt.Errorf("generate key failed! %s", err.Error())
	} else {
		j.P2pId = p2pId.String()

		return j, nil
	}
}
func getAddJsonPath(addr string) string {
	return Cfg.KeyStoreDir + "/" + addr + ".json"
}

func ReadAddrPrivate(jsonFile string, password []byte) (*keystore.Private, error) {
	j, err := keystore.ReadJson(jsonFile)
	if err != nil {
		return nil, err
	}
	return keystore.DecryptPrivate(password, j)
}
