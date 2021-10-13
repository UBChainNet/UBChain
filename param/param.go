package param

import (
	"github.com/UBChainNet/UBChain/common/hasharry"
)

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

	TokenConsumption uint64 = 1 * AtomsPerCoin

	CoinHeight = 565000

	MaximumReceiver = 1000
)

var (
	MainPubKeyHashAddrID  = [3]byte{0x03, 0x77, 0x7d} //UBC 3, 82, 32
	TestPubKeyHashAddrID  = [3]byte{0x06, 0xb5, 0xab} //ubc
	MainPubKeyHashTokenID = [3]byte{0x03, 0x77, 0xa2} //UBT 3, 82, 55
	TestPubKeyHashTokenID = [3]byte{0x06, 0xb5, 0xd2} //ubt
	MainPubKeyHashContractID = [3]byte{0x03, 0x77, 0xab} //UBX 3, 82, 55
	TestPubKeyHashContractID = [3]byte{0x06, 0xb5, 0xdc} //ubx
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


type CandidatesInfo struct {
	Address string
	Reward  string
	PeerId  string
}

var MinerReward = map[string]string{}

func InitMinerReward(){
	for _, candidate := range InitialCandidates{
		MinerReward[candidate.Address] = candidate.Reward
	}
}
// initialCandidates the first super node of the block generation cycle.
// The first half is the address of the block, the second half is the id of the block node
var InitialCandidates = []CandidatesInfo{
	{
		Address: "UBCUnJZqeyLurgWQ2PvQhMGus36t4585T25C",
		Reward:  "UBCdvXm83tnWDkVReQipXnKAkdH1qEnmFHf9",
		PeerId:  "16Uiu2HAkwzWaFqLs6Xn8kUhhXXtcwmRsLkzBrPREpUBEM3qjMwFL",
	},
	{
		Address: "UBCLDyD7zpzuMDJF8GagDgVDtqR7tosdjEWX",
		Reward:  "UBCNry75WY4T6jwdzaBV1VZ9YqQTDs3AZUWj",
		PeerId:  "16Uiu2HAmFTT8zNsHYywAQmcd8k35aDLbfvgKKku7D1Tbc8GApxUt",
	},
	{
		Address: "UBCNZ8bg8DhaoYEprADH3RpRgwvEKXvLPqBj",
		Reward:  "UBCcJ9uNnAbfwJu8dEp8PPmmJUQSJdD8k94X",
		PeerId:  "16Uiu2HAmVaCCCCXHvaCZ9jzeqeEk4WJXFgfDtV4ryPWzmsBwYG9P",
	},
	{
		Address: "UBChoP2pHAZusPKEzYDPUke7Vqn4KWi2UBNJ",
		Reward:  "UBCXxJsYzkUV5ChcBkHpBDvTLCYaTaVfnxHR",
		PeerId:  "16Uiu2HAm7KXydcXNuJp7rZD5VP7eRaqQneW2GJ485yYnKrnwz2LF",
	},
	{
		Address: "UBCLp9kXXGwU9h8Vs6VNQgkrUYEFmU5XqTeo",
		Reward:  "UBCbeMHqktshSb3az45Q9ot5rT37RNmgEdpx",
		PeerId:  "16Uiu2HAmEBBoCPP1CKmbND64pRR7FzNiaZdvMW1DwgPnmSxLFzLj",
	},
	{
		Address: "UBCgavpVj7cHqf9dzWxeaMwshkUNBo4ebpT8",
		Reward:  "UBCNNbMmyknMVoTthTYb5NPvADR2vSt9VJhc",
		PeerId:  "16Uiu2HAmT1U3mQVsNC9GPS7jZEeCeLrineNE98rC1oX5M3oMyvmP",
	},
	{
		Address: "UBCPxNiZiP8q4kdU9awjYNVtsFCdPBTfFwaK",
		Reward:  "UBCWaWPr2sQuw6NmvdPfTSkFbzsWzdDMHWED",
		PeerId:  "16Uiu2HAkxxhNJS8E5qZCHaN2STjRcd1iMzwKmrdhNCAVbKUyK8sm",
	},
	{
		Address: "UBCWvnbP2JJpP74gFejokjWxL9ggG4vd6Gzt",
		Reward:  "UBCT7A9prqVzmhgFa72LjZegYbRxkTTECdz6",
		PeerId:  "16Uiu2HAm97ceWF1xRcrhD9zwM3JMQZ2Ps2RLuekfSZ1RN8YgXian",
	},
	{
		Address: "UBCdN1SfngYus39EhvEsVrm8jmam32fqhjmP",
		Reward:  "UBCKuSDm39PA41NquWJLuw5jc1crDQWmVXgU",
		PeerId:  "16Uiu2HAmDtfuGh4J4yzXAvfy2bG7ttbeWUQ1hQ77YYcq65j921M8",
	},
}




