package param

import "github.com/jhdriver/UBChain/common/hasharry"

const (
	// Mainnet logo
	MainNet = "MN"
	// Testnet logo
	TestNet = "TN"
)

var (
	// Program name
	AppName = "UBChain"
	// Current network
	Net = MainNet
	// Network logo
	UniqueNetWork = "_UBChain"
	// Token name
	Token = hasharry.StringToAddress("UBC")

	FeeAddress   = hasharry.StringToAddress("UBCWUP6SUEmr9A4zn5Zg32ECksunWYCK1pME")
	EaterAddress = hasharry.StringToAddress("UBCCoinEaterAddressDontSend000000000")
)

const (
	// Block interval period
	BlockInterval = uint64(15)
	// Re-election interval
	TermInterval = 60 * 60 * 24 * 365 * 100
	// Maximum number of super nodes
	MaxWinnerSize = 9
	// The minimum number of nodes required to confirm the transaction
	SafeSize = MaxWinnerSize*2/3 + 1
	// The minimum threshold at which a block is valid
	ConsensusSize                 = MaxWinnerSize*2/3 + 1
	SkipCurrentWinnerWaitTimeBase = BlockInterval * MaxWinnerSize * 1
)

const (
	// AtomsPerCoin is the number of atomic units in one coin.
	AtomsPerCoin = 1e8

	// Circulation is the total number of COINS issued.
	Circulation = 6300 * 1e4 * AtomsPerCoin

	// GenesisCoins is genesis Coins
	GenesisCoins = 5000 * 1e4 * AtomsPerCoin

	// CoinBaseCoins is reward
	CoinBaseCoins = 3 * AtomsPerCoin

	//MaxAddressTxs is address the maximum number of transactions in the trading pool
	MaxAddressTxs = 1000

	// MinFeesCoefficient is minimum fee required to process the transaction
	MinFeesCoefficient uint64 = 1e4

	// MaxFeesCoefficient is maximum fee required to process the transaction
	MaxFeesCoefficient uint64 = 1 * AtomsPerCoin

	// MinAllowedAmount is the minimum allowable amount for a transaction
	MinAllowedAmount uint64 = 0.005 * AtomsPerCoin

	// MaxAllContractCoin is the maximum allowable sum of contract COINS
	MaxAllContractCoin uint64 = 1e11 * AtomsPerCoin

	// MaxContractCoin is the maximum allowable contract COINS
	MaxContractCoin uint64 = 1e10 * AtomsPerCoin

	Fees uint64 = 0.005 * AtomsPerCoin

	TokenConsumption uint64 = 10.24 * AtomsPerCoin

	CoinHeight = 1

	MaximumReceiver = 1000
)

var (
	MainPubKeyHashAddrID  = [3]byte{0x03, 0x77, 0x7d} //UBC 3, 82, 32
	TestPubKeyHashAddrID  = [3]byte{0x06, 0xb5, 0xab} //ubc
	MainPubKeyHashTokenID = [3]byte{0x03, 0x77, 0xa2} //UBT 3, 82, 55
	TestPubKeyHashTokenID = [3]byte{0x06, 0xb5, 0xd2} //ubt
)

type MappingInfo struct {
	Address string
	Note    string
	Amount  uint64
}

var MappingCoin = []MappingInfo{
	{
		Address: "UBCVHn5fGP34Uf3iRCHJrw2HnvCWVVJmmZ3K",
		Note:    "",
		Amount:  5000 * 1e4 * 1e8,
	},
}
