package dpos

type CandidatesInfo struct {
	Address string
	PeerId  string
}

// initialCandidates the first super node of the block generation cycle.
// The first half is the address of the block, the second half is the id of the block node
var initialCandidates = []CandidatesInfo{
	{
		Address: "UBCUnJZqeyLurgWQ2PvQhMGus36t4585T25C",
		PeerId:  "16Uiu2HAkwzWaFqLs6Xn8kUhhXXtcwmRsLkzBrPREpUBEM3qjMwFL",
	},
	{
		Address: "UBCLDyD7zpzuMDJF8GagDgVDtqR7tosdjEWX",
		PeerId:  "16Uiu2HAmFTT8zNsHYywAQmcd8k35aDLbfvgKKku7D1Tbc8GApxUt",
	},
	{
		Address: "UBCNZ8bg8DhaoYEprADH3RpRgwvEKXvLPqBj",
		PeerId:  "16Uiu2HAmVaCCCCXHvaCZ9jzeqeEk4WJXFgfDtV4ryPWzmsBwYG9P",
	},
	{
		Address: "UBChoP2pHAZusPKEzYDPUke7Vqn4KWi2UBNJ",
		PeerId:  "16Uiu2HAm7KXydcXNuJp7rZD5VP7eRaqQneW2GJ485yYnKrnwz2LF",
	},
	{
		Address: "UBCLp9kXXGwU9h8Vs6VNQgkrUYEFmU5XqTeo",
		PeerId:  "16Uiu2HAmEBBoCPP1CKmbND64pRR7FzNiaZdvMW1DwgPnmSxLFzLj",
	},
	{
		Address: "UBCgavpVj7cHqf9dzWxeaMwshkUNBo4ebpT8",
		PeerId:  "16Uiu2HAmT1U3mQVsNC9GPS7jZEeCeLrineNE98rC1oX5M3oMyvmP",
	},
	{
		Address: "UBCPxNiZiP8q4kdU9awjYNVtsFCdPBTfFwaK",
		PeerId:  "16Uiu2HAkxxhNJS8E5qZCHaN2STjRcd1iMzwKmrdhNCAVbKUyK8sm",
	},
	{
		Address: "UBCWvnbP2JJpP74gFejokjWxL9ggG4vd6Gzt",
		PeerId:  "16Uiu2HAm97ceWF1xRcrhD9zwM3JMQZ2Ps2RLuekfSZ1RN8YgXian",
	},
	{
		Address: "UBCdN1SfngYus39EhvEsVrm8jmam32fqhjmP",
		PeerId:  "16Uiu2HAmDtfuGh4J4yzXAvfy2bG7ttbeWUQ1hQ77YYcq65j921M8",
	},
}
