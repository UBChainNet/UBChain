package contractdb

import (
	"github.com/UBChainNet/UBChain/common/codec"
	"github.com/UBChainNet/UBChain/common/hasharry"
	"github.com/UBChainNet/UBChain/core/types"
	"github.com/UBChainNet/UBChain/core/types/contractv2"
	"github.com/UBChainNet/UBChain/database/triedb"
	"github.com/UBChainNet/UBChain/trie"
)

const synmbolPrefix = "s_"

type ContractStorage struct {
	trieDB       *triedb.TrieDB
	contractTrie *trie.Trie
}

func NewContractStorage(path string) *ContractStorage {
	trieDB := triedb.NewTrieDB(path)
	return &ContractStorage{trieDB, nil}
}

func (c *ContractStorage) InitTrie(contractRoot hasharry.Hash) error {
	contractTrie, err := trie.New(contractRoot, c.trieDB)
	if err != nil {
		return err
	}
	c.contractTrie = contractTrie
	return nil
}

func (c *ContractStorage) Commit() (hasharry.Hash, error) {
	return c.contractTrie.Commit()
}

func (c *ContractStorage) RootHash() hasharry.Hash {
	return c.contractTrie.Hash()
}

func (c *ContractStorage) Open() error {
	return c.trieDB.Open()
}

func (c *ContractStorage) Close() error {
	return c.trieDB.Close()
}

func (c *ContractStorage) GetContract(contractAddr string) *types.Contract {
	contract := types.NewContract()
	bytes := c.contractTrie.Get([]byte(contractAddr))
	err := codec.FromBytes(bytes, &contract)
	if err != nil {
		return nil
	}
	return contract
}

func (c *ContractStorage) SetContract(contract *types.Contract) {
	bytes, err := codec.ToBytes(contract)
	if err != nil {
		return
	}
	c.contractTrie.Update([]byte(contract.Contract), bytes)
}

func (c *ContractStorage) GetContractV2(contractAddr string) *contractv2.ContractV2 {
	bytes := c.contractTrie.Get(hasharry.StringToAddress(contractAddr).Bytes())
	contract, _ := contractv2.DecodeContractV2(bytes)
	return contract
}

func (c *ContractStorage) SetContractV2(contract *contractv2.ContractV2) {
	c.contractTrie.Update(contract.Address.Bytes(), contract.Bytes())
}

func (c *ContractStorage) SetContractV2State(txHash string, state *types.ContractV2State) {
	c.contractTrie.Update([]byte(txHash), state.Bytes())
}

func (c *ContractStorage) GetContractV2State(txHash string) *types.ContractV2State {
	bytes := c.contractTrie.Get([]byte(txHash))
	cs, _ := types.DecodeContractV2State(bytes)
	return cs
}

func (c *ContractStorage) SetSymbol(symbol string, contract string) {
	c.contractTrie.Update([]byte(synmbolPrefix + symbol), hasharry.StringToAddress(contract).Bytes())
}

func (c *ContractStorage) GetSymbol(symbol string) (string, bool){
	bytes := c.contractTrie.Get([]byte(synmbolPrefix + symbol))
	if bytes == nil{
		return "", false
	}
	return hasharry.BytesToAddress(bytes).String(), true
}
