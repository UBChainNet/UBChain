package contractstate

import (
	"fmt"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/core/types"
	"github.com/UBChainNet/UBChain/core/types/contractv2"
	"github.com/UBChainNet/UBChain/database/contractdb"
	"github.com/UBChainNet/UBChain/param"
	"sync"
)

const contractSate = "contract_state"

// Contract status, used to store all published contract information records
type ContractState struct {
	contractDb      IContractStorage
	contractMutex   sync.RWMutex
	confirmedHeight uint64
}

func NewContractState(dataDir string) (*ContractState, error) {
	storage := contractdb.NewContractStorage(dataDir + "/" + contractSate)
	err := storage.Open()
	if err != nil {
		return nil, err
	}
	return &ContractState{
		contractDb: storage,
	}, nil
}

// Initialize the contract state tree
func (cs *ContractState) InitTrie(contractRoot hasharry.Hash) error {
	return cs.contractDb.InitTrie(contractRoot)
}

func (cs *ContractState) RootHash() hasharry.Hash {
	return cs.contractDb.RootHash()
}

// Commit contract status changes
func (cs *ContractState) ContractTrieCommit() (hasharry.Hash, error) {
	return cs.contractDb.Commit()
}

func (c *ContractState) GetContract(contractAddr string) *types.Contract {
	c.contractMutex.RLock()
	defer c.contractMutex.RUnlock()

	contract := c.contractDb.GetContract(contractAddr)
	return contract
}

func (c *ContractState) SetContract(contract *types.Contract) {
	c.contractMutex.Lock()
	defer c.contractMutex.Unlock()

	c.contractDb.SetContract(contract)
}

func (c *ContractState) GetContractV2(contractAddr string) *contractv2.ContractV2 {
	c.contractMutex.RLock()
	defer c.contractMutex.RUnlock()

	contract := c.contractDb.GetContractV2(contractAddr)
	return contract
}

func (c *ContractState) SetContractV2(contract *contractv2.ContractV2) {
	c.contractMutex.RLock()
	defer c.contractMutex.RUnlock()

	c.contractDb.SetContractV2(contract)
}

func (c *ContractState) GetContractV2State(txHash string) *types.ContractV2State {
	c.contractMutex.RLock()
	defer c.contractMutex.RUnlock()

	state := c.contractDb.GetContractV2State(txHash)
	return state
}

func (c *ContractState) SetContractV2State(txHash string, contract *types.ContractV2State) {
	c.contractMutex.Lock()
	defer c.contractMutex.Unlock()

	c.contractDb.SetContractV2State(txHash, contract)
}

func (c *ContractState) UpdateConfirmedHeight(height uint64) {
	c.confirmedHeight = height
}

// Verification contract
func (c *ContractState) VerifyState(tx types.ITransaction) error {
	c.contractMutex.RLock()
	defer c.contractMutex.RUnlock()

	if tx.GetTxType() != types.Contract_ {
		return nil
	}
	contractAddr := tx.GetTxBody().GetContract()
	contract := c.contractDb.GetContract(contractAddr.String())
	if contract != nil {
		return contract.Verify(tx)
	} else {
		_, exist := c.contractDb.GetSymbolContract(tx.GetTxBody().GetAbbr())
		if exist {
			return fmt.Errorf("%s already exists", tx.GetTxBody().GetAbbr())
		}
	}
	return nil
}

// Update contract status
func (c *ContractState) UpdateContract(tx types.ITransaction, blockHeight uint64) {
	c.contractMutex.Lock()
	defer c.contractMutex.Unlock()

	txBody := tx.GetTxBody()
	contractRecord := &types.ContractRecord{
		Height:   blockHeight,
		TxHash:   tx.Hash(),
		Time:     tx.GetTime(),
		Amount:   txBody.GetAmount(),
		Receiver: txBody.ToAddress().ReceiverList()[0].Address.String(),
	}
	contractAddr := txBody.GetContract()
	contract := c.contractDb.GetContract(contractAddr.String())
	if contract != nil {
		contract.AddContract(contractRecord)
	} else {
		contract = &types.Contract{
			Sender:         tx.From().String(),
			Contract:       contractAddr.String(),
			CoinName:       txBody.GetName(),
			CoinAbbr:       txBody.GetAbbr(),
			Description:    txBody.GetDescription(),
			IncreaseSwitch: txBody.GetIncreaseSwitch(),
			Records: &types.RecordList{
				contractRecord,
			},
		}
		c.contractDb.SetSymbol(txBody.GetAbbr(), contractAddr.String())
	}
	c.contractDb.SetContract(contract)
}

func (c *ContractState) SetSymbol(symbol string, contract string) {
	c.contractMutex.Lock()
	defer c.contractMutex.Unlock()

	c.contractDb.SetSymbol(symbol, contract)
}

func (c *ContractState) GetSymbolContract(symbol string) (string, bool) {
	if symbol == param.Token.String() {
		return param.Token.String(), true
	}
	c.contractMutex.RLock()
	defer c.contractMutex.RUnlock()

	return c.contractDb.GetSymbolContract(symbol)
}

func (c *ContractState) TokenList() []*types.Token {
	c.contractMutex.RLock()
	defer c.contractMutex.RUnlock()

	return c.contractDb.TokenList()
}

func (c *ContractState) Close() error {
	c.contractMutex.Lock()
	defer c.contractMutex.Unlock()

	return c.contractDb.Close()
}
