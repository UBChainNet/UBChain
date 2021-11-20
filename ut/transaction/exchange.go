package transaction

import (
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/core/runner/exchange_runner"
	"github.com/UBChainNet/UBChain/core/types"
	"github.com/UBChainNet/UBChain/core/types/contractv2"
	"github.com/UBChainNet/UBChain/core/types/functionbody/exchange_func"
	"github.com/UBChainNet/UBChain/param"
	"time"
)

func NewExchange(net, from, admin, feeTo, symbol string, nonce uint64, note string) (*types.Transaction, error) {
	contract, err := exchange_runner.ExchangeAddress(net, from, nonce)
	if err != nil {
		return nil, err
	}

	tx := &types.Transaction{
		TxHead: &types.TransactionHead{
			TxType:     types.ContractV2_,
			TxHash:     hasharry.Hash{},
			From:       hasharry.StringToAddress(from),
			Nonce:      nonce,
			Time:       uint64(time.Now().Unix()),
			Note:       note,
			SignScript: &types.SignScript{},
			Fees:       param.Fees,
		},
		TxBody: &types.TxContractV2Body{
			Contract:     hasharry.StringToAddress(contract),
			Type:         contractv2.Exchange_,
			FunctionType: contractv2.Exchange_Init,
			Function: &exchange_func.ExchangeInitBody{
				Admin:  hasharry.StringToAddress(admin),
				FeeTo:  hasharry.StringToAddress(feeTo),
				Symbol: symbol,
			},
		},
	}
	tx.SetHash()
	return tx, nil
}

func NewSetAdmin(from, exchange, admin string, nonce uint64, note string) (*types.Transaction, error) {
	tx := &types.Transaction{
		TxHead: &types.TransactionHead{
			TxType:     types.ContractV2_,
			TxHash:     hasharry.Hash{},
			From:       hasharry.StringToAddress(from),
			Nonce:      nonce,
			Time:       uint64(time.Now().Unix()),
			Note:       note,
			SignScript: &types.SignScript{},
			Fees:       param.Fees,
		},
		TxBody: &types.TxContractV2Body{
			Contract:     hasharry.StringToAddress(exchange),
			Type:         contractv2.Exchange_,
			FunctionType: contractv2.Exchange_SetAdmin,
			Function: &exchange_func.ExchangeAdmin{
				Address: hasharry.StringToAddress(admin),
			},
		},
	}
	tx.SetHash()
	return tx, nil
}

func NewSetFeeTo(from, exchange, feeTo string, nonce uint64, note string) (*types.Transaction, error) {
	tx := &types.Transaction{
		TxHead: &types.TransactionHead{
			TxType:     types.ContractV2_,
			TxHash:     hasharry.Hash{},
			From:       hasharry.StringToAddress(from),
			Nonce:      nonce,
			Time:       uint64(time.Now().Unix()),
			Note:       note,
			SignScript: &types.SignScript{},
			Fees:       param.Fees,
		},
		TxBody: &types.TxContractV2Body{
			Contract:     hasharry.StringToAddress(exchange),
			Type:         contractv2.Exchange_,
			FunctionType: contractv2.Exchange_SetFeeTo,
			Function: &exchange_func.ExchangeFeeTo{
				Address: hasharry.StringToAddress(feeTo),
			},
		},
	}
	tx.SetHash()
	return tx, nil
}

func NewPairAddLiquidity(net, from, to, exchange, tokenA, tokenB string, amountADesired, amountBDesired, amountAMin, amountBMin, deadline, nonce uint64, note string) (*types.Transaction, error) {
	contract, err := exchange_runner.PairAddress(net, hasharry.StringToAddress(tokenA), hasharry.StringToAddress(tokenB), hasharry.StringToAddress(exchange))
	if err != nil {
		return nil, err
	}
	tx := &types.Transaction{
		TxHead: &types.TransactionHead{
			TxType:     types.ContractV2_,
			TxHash:     hasharry.Hash{},
			From:       hasharry.StringToAddress(from),
			Nonce:      nonce,
			Time:       uint64(time.Now().Unix()),
			Note:       note,
			SignScript: &types.SignScript{},
			Fees:       param.Fees,
		},
		TxBody: &types.TxContractV2Body{
			Contract:     hasharry.StringToAddress(contract),
			Type:         contractv2.Pair_,
			FunctionType: contractv2.Pair_AddLiquidity,
			Function: &exchange_func.ExchangeAddLiquidity{
				Exchange:       hasharry.StringToAddress(exchange),
				TokenA:         hasharry.StringToAddress(tokenA),
				TokenB:         hasharry.StringToAddress(tokenB),
				To:             hasharry.StringToAddress(to),
				AmountADesired: amountADesired,
				AmountBDesired: amountBDesired,
				AmountAMin:     amountAMin,
				AmountBMin:     amountBMin,
				Deadline:       deadline,
			},
		},
	}
	tx.SetHash()
	return tx, nil
}

func NewPairRemoveLiquidity(net, from, to, exchange, tokenA, tokenB string, amountAMin, amountBMin, liquidity, deadline, nonce uint64, note string) (*types.Transaction, error) {
	contract, err := exchange_runner.PairAddress(net, hasharry.StringToAddress(tokenA), hasharry.StringToAddress(tokenB), hasharry.StringToAddress(exchange))
	if err != nil {
		return nil, err
	}
	tx := &types.Transaction{
		TxHead: &types.TransactionHead{
			TxType:     types.ContractV2_,
			TxHash:     hasharry.Hash{},
			From:       hasharry.StringToAddress(from),
			Nonce:      nonce,
			Time:       uint64(time.Now().Unix()),
			Note:       note,
			SignScript: &types.SignScript{},
			Fees:       param.Fees,
		},
		TxBody: &types.TxContractV2Body{
			Contract:     hasharry.StringToAddress(contract),
			Type:         contractv2.Pair_,
			FunctionType: contractv2.Pair_RemoveLiquidity,
			Function: &exchange_func.ExchangeRemoveLiquidity{
				Exchange:   hasharry.StringToAddress(exchange),
				TokenA:     hasharry.StringToAddress(tokenA),
				TokenB:     hasharry.StringToAddress(tokenB),
				To:         hasharry.StringToAddress(to),
				Liquidity:  liquidity,
				AmountAMin: amountAMin,
				AmountBMin: amountBMin,
				Deadline:   deadline,
			},
		},
	}
	tx.SetHash()
	return tx, nil
}

func NewSwapExactIn(from, to, exchange string, amountIn, amountOutMin uint64, path []string, deadline, nonce uint64, note string) (*types.Transaction, error) {
	address := make([]hasharry.Address, 0)
	for _, addr := range path {
		address = append(address, hasharry.StringToAddress(addr))
	}
	tx := &types.Transaction{
		TxHead: &types.TransactionHead{
			TxType:     types.ContractV2_,
			TxHash:     hasharry.Hash{},
			From:       hasharry.StringToAddress(from),
			Nonce:      nonce,
			Time:       uint64(time.Now().Unix()),
			Note:       note,
			SignScript: &types.SignScript{},
			Fees:       param.Fees,
		},
		TxBody: &types.TxContractV2Body{
			Contract:     hasharry.StringToAddress(exchange),
			Type:         contractv2.Exchange_,
			FunctionType: contractv2.Exchange_ExactIn,
			Function: &exchange_func.ExactIn{
				AmountIn:     amountIn,
				AmountOutMin: amountOutMin,
				Path:         address,
				To:           hasharry.StringToAddress(to),
				Deadline:     deadline,
			},
		},
	}
	tx.SetHash()
	return tx, nil
}

func NewSwapExactOut(from, to, exchange string, amountOut, amountInMax uint64, path []string, deadline, nonce uint64, note string) (*types.Transaction, error) {
	address := make([]hasharry.Address, 0)
	for _, addr := range path {
		address = append(address, hasharry.StringToAddress(addr))
	}
	tx := &types.Transaction{
		TxHead: &types.TransactionHead{
			TxType:     types.ContractV2_,
			TxHash:     hasharry.Hash{},
			From:       hasharry.StringToAddress(from),
			Nonce:      nonce,
			Time:       uint64(time.Now().Unix()),
			Note:       note,
			SignScript: &types.SignScript{},
			Fees:       param.Fees,
		},
		TxBody: &types.TxContractV2Body{
			Contract:     hasharry.StringToAddress(exchange),
			Type:         contractv2.Exchange_,
			FunctionType: contractv2.Exchange_ExactOut,
			Function: &exchange_func.ExactOut{
				AmountOut:   amountOut,
				AmountInMax: amountInMax,
				Path:        address,
				To:          hasharry.StringToAddress(to),
				Deadline:    deadline,
			},
		},
	}
	tx.SetHash()
	return tx, nil
}

func NewPledgeInit(from, admin, contract, exchange, receive string, maxSupply, preMint,  nonce uint64, note string) (*types.Transaction, error) {
	tx := &types.Transaction{
		TxHead: &types.TransactionHead{
			TxType:     types.ContractV2_,
			TxHash:     hasharry.Hash{},
			From:       hasharry.StringToAddress(from),
			Nonce:      nonce,
			Time:       uint64(time.Now().Unix()),
			Note:       note,
			SignScript: &types.SignScript{},
			Fees:       param.Fees,
		},
		TxBody: &types.TxContractV2Body{
			Contract:     hasharry.StringToAddress(contract),
			Type:         contractv2.Pledge_,
			FunctionType: contractv2.Pledge_Init,
			Function: &exchange_func.PledgeInitBody{
				Exchange:         hasharry.StringToAddress(exchange),
				Receiver:         hasharry.StringToAddress(receive),
				Admin:            hasharry.StringToAddress(admin),
				PreMint:          preMint,
				MaxSupply:        maxSupply,
			},
		},
	}
	tx.SetHash()
	return tx, nil
}

func NewPledgeStart(from, contract string, blockMint, matureTime, nonce uint64, note string) (*types.Transaction, error) {
	tx := &types.Transaction{
		TxHead: &types.TransactionHead{
			TxType:     types.ContractV2_,
			TxHash:     hasharry.Hash{},
			From:       hasharry.StringToAddress(from),
			Nonce:      nonce,
			Time:       uint64(time.Now().Unix()),
			Note:       note,
			SignScript: &types.SignScript{},
			Fees:       param.Fees,
		},
		TxBody: &types.TxContractV2Body{
			Contract:     hasharry.StringToAddress(contract),
			Type:         contractv2.Pledge_,
			FunctionType: contractv2.Pledge_Start,
			Function: &exchange_func.PledgeStartBody{
				BlockMintAmount:  blockMint,
				PledgeMatureTime: matureTime,
			},
		},
	}
	tx.SetHash()
	return tx, nil
}

func NewAddPledgePool(from, contract, pair string, blockReward, nonce uint64, note string) (*types.Transaction, error) {
	tx := &types.Transaction{
		TxHead: &types.TransactionHead{
			TxType:     types.ContractV2_,
			TxHash:     hasharry.Hash{},
			From:       hasharry.StringToAddress(from),
			Nonce:      nonce,
			Time:       uint64(time.Now().Unix()),
			Note:       note,
			SignScript: &types.SignScript{},
			Fees:       param.Fees,
		},
		TxBody: &types.TxContractV2Body{
			Contract:     hasharry.StringToAddress(contract),
			Type:         contractv2.Pledge_,
			FunctionType: contractv2.Pledge_AddPool,
			Function: &exchange_func.PledgeAddPoolBody{
				Pair: hasharry.StringToAddress(pair),
				BlockReward: blockReward,
			},
		},
	}
	tx.SetHash()
	return tx, nil
}


func NewRemovePledgePool(from, contract, pair string, nonce uint64, note string) (*types.Transaction, error) {
	tx := &types.Transaction{
		TxHead: &types.TransactionHead{
			TxType:     types.ContractV2_,
			TxHash:     hasharry.Hash{},
			From:       hasharry.StringToAddress(from),
			Nonce:      nonce,
			Time:       uint64(time.Now().Unix()),
			Note:       note,
			SignScript: &types.SignScript{},
			Fees:       param.Fees,
		},
		TxBody: &types.TxContractV2Body{
			Contract:     hasharry.StringToAddress(contract),
			Type:         contractv2.Pledge_,
			FunctionType: contractv2.Pledge_RemovePool,
			Function: &exchange_func.PledgeRemovePoolBody{
				Pair: hasharry.StringToAddress(pair),
			},
		},
	}
	tx.SetHash()
	return tx, nil
}

func NewAddPledge(from, contract, pair string, amount, nonce uint64, note string) (*types.Transaction, error) {
	tx := &types.Transaction{
		TxHead: &types.TransactionHead{
			TxType:     types.ContractV2_,
			TxHash:     hasharry.Hash{},
			From:       hasharry.StringToAddress(from),
			Nonce:      nonce,
			Time:       uint64(time.Now().Unix()),
			Note:       note,
			SignScript: &types.SignScript{},
			Fees:       param.Fees,
		},
		TxBody: &types.TxContractV2Body{
			Contract:     hasharry.StringToAddress(contract),
			Type:         contractv2.Pledge_,
			FunctionType: contractv2.Pledge_Add,
			Function: &exchange_func.PledgeAddBody{
				Pair:   hasharry.StringToAddress(pair),
				Amount: amount,
			},
		},
	}
	tx.SetHash()
	return tx, nil
}

func NewRemovePledge(from, contract, pair string, amount, nonce uint64, note string) (*types.Transaction, error) {
	tx := &types.Transaction{
		TxHead: &types.TransactionHead{
			TxType:     types.ContractV2_,
			TxHash:     hasharry.Hash{},
			From:       hasharry.StringToAddress(from),
			Nonce:      nonce,
			Time:       uint64(time.Now().Unix()),
			Note:       note,
			SignScript: &types.SignScript{},
			Fees:       param.Fees,
		},
		TxBody: &types.TxContractV2Body{
			Contract:     hasharry.StringToAddress(contract),
			Type:         contractv2.Pledge_,
			FunctionType: contractv2.Pledge_Remove,
			Function: &exchange_func.PledgeRemoveBody{
				Pair:   hasharry.StringToAddress(pair),
				Amount: amount,
			},
		},
	}
	tx.SetHash()
	return tx, nil
}

func NewRemovePledgeReward(from, contract string, nonce uint64, note string) (*types.Transaction, error) {
	tx := &types.Transaction{
		TxHead: &types.TransactionHead{
			TxType:     types.ContractV2_,
			TxHash:     hasharry.Hash{},
			From:       hasharry.StringToAddress(from),
			Nonce:      nonce,
			Time:       uint64(time.Now().Unix()),
			Note:       note,
			SignScript: &types.SignScript{},
			Fees:       param.Fees,
		},
		TxBody: &types.TxContractV2Body{
			Contract:     hasharry.StringToAddress(contract),
			Type:         contractv2.Pledge_,
			FunctionType: contractv2.Pledge_RemoveReward,
			Function:     &exchange_func.PledgeRewardRemoveBody{},
		},
	}
	tx.SetHash()
	return tx, nil
}

func NewUpdatePledgeReward(from, contract string, nonce uint64, note string) (*types.Transaction, error) {
	tx := &types.Transaction{
		TxHead: &types.TransactionHead{
			TxType:     types.ContractV2_,
			TxHash:     hasharry.Hash{},
			From:       hasharry.StringToAddress(from),
			Nonce:      nonce,
			Time:       uint64(time.Now().Unix()),
			Note:       note,
			SignScript: &types.SignScript{},
			Fees:       param.Fees,
		},
		TxBody: &types.TxContractV2Body{
			Contract:     hasharry.StringToAddress(contract),
			Type:         contractv2.Pledge_,
			FunctionType: contractv2.Pledge_Update,
			Function:     &exchange_func.PledgeUpdateBody{},
		},
	}
	tx.SetHash()
	return tx, nil
}
