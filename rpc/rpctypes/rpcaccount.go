package rpctypes

import (
	"github.com/jhdriver/UBChain/common/hasharry"
	"github.com/jhdriver/UBChain/core/types"
)

type Account struct {
	Address         string         `json:"address"`
	Nonce           uint64         `json:"nonce"`
	Time            uint64         `json:"time"`
	Coins           []*CoinAccount `json:"coins"`
	ConfirmedHeight uint64         `json:"confirmedheight"`
	ConfirmedNonce  uint64         `json:"confirmednonce"`
	ConfirmedTime   uint64         `json:"confirmedtime"`
}

type CoinAccount struct {
	Contract  string  `json:"contract"`
	Balance   float64 `json:"balance"`
	LockedOut float64 `json:"lockedout"`
	LockedIn  float64 `json:"lockedin"`
}

func TranslateAccountToRpcAccount(account *types.Account) *Account {
	coins := []*CoinAccount{}
	for _, coinAccount := range *account.Coins {
		coins = append(coins, &CoinAccount{
			Contract:  coinAccount.Contract,
			LockedOut: types.Amount(coinAccount.LockOut).ToCoin(),
			LockedIn:  types.Amount(coinAccount.LockIn).ToCoin(),
			Balance:   types.Amount(coinAccount.Balance).ToCoin(),
		})
	}
	rpcAccount := &Account{
		Address:         account.Address.String(),
		Nonce:           account.Nonce,
		Time:            account.Time,
		Coins:           coins,
		ConfirmedHeight: account.ConfirmedHeight,
		ConfirmedNonce:  account.ConfirmedNonce,
		ConfirmedTime:   account.ConfirmedTime,
	}
	return rpcAccount
}

func TranslateAddressListToRpcAddressList(list []hasharry.Address) []string {
	address := []string{}
	for _, ls := range list {
		address = append(address, ls.String())
	}
	return address
}
