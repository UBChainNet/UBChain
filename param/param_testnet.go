package param

import (
	"github.com/UBChainNet/UBChain/common/hasharry"
)

var (
	FeeAddressTest   = hasharry.StringToAddress("ubcbaQvSw9oerHyJUWah4GhSkdUPAmqg4qpx")
	EaterAddressTest = hasharry.StringToAddress("ubcCoinEaterAddressDontSend000000000")
)

const (
	// Block interval period
	BlockIntervalTest = uint64(5)
	// Re-election interval
	TermIntervalTest uint64 = 60 * 60 * 24 * 365 * 100
	// Maximum number of super nodes
	MaxWinnerSizeTest = 3
	// The minimum number of nodes required to confirm the transaction
	SafeSizeTest = MaxWinnerSizeTest*2/3 + 1
	// The minimum threshold at which a block is valid
	ConsensusSizeTest                 = MaxWinnerSizeTest*2/3 + 1
	SkipCurrentWinnerWaitTimeBaseTest = int(BlockIntervalTest) * (MaxWinnerSizeTest) * 1
)

const (
	UIPBlock1Test = 0
	UIPBlock2Test = 0
	UIPBlock3Test = 0
	UIPBlock4Test = 238594
)

var MappingCoinTest = []MappingInfo{
	{
		Address: "ubcdkaH2gypSVGiq3bLG1NJsqwqXautd3xak",
		Note:    "",
		Amount:  5000 * 1e4 * 1e8,
	},
}

// initialCandidates the first super node of the block generation cycle.
// The first half is the address of the block, the second half is the id of the block node
var InitialCandidatesTest = []CandidatesInfo{
	{
		Address: "ubci8DDtULsHyYAa8XJMD5qADv2JdCDGcUkJ",
		Reward:  "ubci8DDtULsHyYAa8XJMD5qADv2JdCDGcUkJ",
		PeerId:  "16Uiu2HAkvMA9jU45K8ndvJ4rsRL8a91d4YtLc1kxa2LRaAeeUFYX",
	},
	{
		Address: "ubcbVHm347qfaFw9bpyDyztTyhvqDbVZbD5t",
		Reward:  "ubcbVHm347qfaFw9bpyDyztTyhvqDbVZbD5t",
		PeerId:  "16Uiu2HAmMhj8As8V7Q9oVyPB65Bb7E5oXzwe2BHK16HuqKKSCS5S",
	},
	{
		Address: "ubcgZrNvnQED7zHZhF5AHeMkBHYcqZCRmcWN",
		Reward:  "ubcgZrNvnQED7zHZhF5AHeMkBHYcqZCRmcWN",
		PeerId:  "16Uiu2HAmL3mrRHLRdiSovm6pXCuepeVE4zT76AKtiWGzMwFnXn7G",
	},
}

var BootsTest = []string{
	// ubcd5FTPtJiRWgmSYNDFzKEKqXoPxdqNJ5tF
	"/ip4/180.188.198.214/tcp/2211/ipfs/16Uiu2HAm39NFZjFVoauqtgDTPrDptPBCMADQ5pd6Ynyb9JmLX9my",
}
