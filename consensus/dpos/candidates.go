package dpos

type CandidatesInfo struct {
	Address string
	PeerId  string
}

// initialCandidates the first super node of the block generation cycle.
// The first half is the address of the block, the second half is the id of the block node
var initialCandidates = []CandidatesInfo{
	{
		Address: "UBCRz5TqWKE3WvibnGPx5XEgbNN6iUZx7A8p",
		PeerId:  "16Uiu2HAmRhxt3cwXWDU8sTcHyKwVvAfFXdmv6jkfKFojMABrA99k",
	},
	{
		Address: "UBCdKrHzLifc4i5HdL7Tt5UsXTrsxP7yY4M8",
		PeerId:  "16Uiu2HAmMzmPptV4nXXpnasuqe3FRRxUP9ChvxPLjdre691BY1ev",
	},
	{
		Address: "UBCPEQHwZK5AhGhp8enFnAgZLXmQQaA4gfdN",
		PeerId:  "16Uiu2HAmGaUyLvLQudCJdiJBqarstgg5ndc5LTt3KJrHsP8Ssdru",
	},
	/*{
		Address: "UBCcSaX5hUhutRDC4EQT65DLeRJxBJWBWppX",
		PeerId:  "16Uiu2HAmK2PNwXpTM91ftN9ZF92CNLKf5m7XDApBaBa69cwPDpMv",
	},
	{
		Address: "UBCHR3U2FQDNYdQpgUrHw3LkrMiC23wetsb6",
		PeerId:  "16Uiu2HAmKCs2So8WRqZj5aWQKHAMj6DhLd1eQQBvF8MdjgbajUHL",
	},
	{
		Address: "UBCamE7HaJNvRrRDx3Rnmj9WKqHnDMFzQ6P3",
		PeerId:  "16Uiu2HAkvRSRt5VEeN8vXKg1sxDSeXwJfQzPqYPPneoaTNwYxQuM",
	},
	{
		Address: "UBCRKs3eg7deVcRo5pJVCsRAVQ1umXNM1DiD",
		PeerId:  "16Uiu2HAmABaSg6ZmBpKXoSCV5V89fNS4N3s91wRAbneq7DJPN81y",
	},
	{
		Address: "UBCZqKCCQLV2yTZ7crW6kZn7UFLPCpHtFdGi",
		PeerId:  "16Uiu2HAm9w9gZsnS7ZKGqdLfNgXuyMrfhBgSKHWePVfhHGVzDDxH",
	},
	{
		Address: "UBCKvSoMxcD7MbT1KZ3kX5xjnL4tWK7XXGCd",
		PeerId:  "16Uiu2HAmULDuqSBP9mueYjHCknUuex3yGVQM4QXgBNTHf8qiNFD6",
	},
	{
		Address: "UBCN2rVzEgJRVGyZFUQstL2y2JpSkGVz1kbY",
		PeerId:  "16Uiu2HAkxKA74FNnwY5hAJGW64GyDgfWFeCtkkuFJwXbSLeVKK1t",
	},
	{
		Address: "UBCToaCgPNCZ168ZnvUg38bRMm5U7DjwRVuA",
		PeerId:  "16Uiu2HAmRsZi1iWAx1efSv7o5dvzeYoRXhQusNtjzbhmmuqUo99q",
	},*/
}
