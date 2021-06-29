package dpos

type CandidatesInfo struct {
	Address string
	PeerId  string
}

// initialCandidates the first super node of the block generation cycle.
// The first half is the address of the block, the second half is the id of the block node
var initialCandidates = []CandidatesInfo{
	{
		Address: "UBCZ4BB2hGjMqKprGpmgKwBfzkQ7aUs3Jyj1",
		PeerId:  "16Uiu2HAmHcDq73pnXsWMHcFp7hKJsUkA8CXcKhVpWm4LcMzWx1ya",
	},
	{
		Address: "UBCaWBrLUF4Mg7eQqC1ufySAkd59kL5r1dz8",
		PeerId:  "16Uiu2HAmCfo2diJGqEmz3JfZxkqp6J2aUymwBvWX3Lptsx9aQEe9",
	},
	{
		Address: "UBCNL8raaNiriNJai14BJdcDLsQFyAms5XXd",
		PeerId:  "16Uiu2HAmVaAzh6wS1mmypZCJWkMY35ftZPVnxoLdfDNSCwv8DrLT",
	},
	{
		Address: "UBCcwMpxPt9XcjQSGJZLNoujhDsaGSsdRhZ7",
		PeerId:  "16Uiu2HAkz6fzCgpE7A2EhcNX176ULuEGvRNuKjJnVSXA9yohrQWV",
	},
	{
		Address: "UBCbnPFVh4GaSy8hzLefc9ppxoKuLyzmUCf6",
		PeerId:  "16Uiu2HAmJQY2tQk9CQjYRS1c2JXMTGyyNC166ZNenzNMtA7DA8Uu",
	},
	{
		Address: "UBCTZabxNmMscnk8fJhDkZRYCRXxraziFd7Y",
		PeerId:  "16Uiu2HAm62n2nprp2Kf15DhQe6RVCpj4rsjgAZQLeNo3HreM6TQ4",
	},
	{
		Address: "UBChEPgFDfGynFERskUXA8boDL7kjcfE7pQq",
		PeerId:  "16Uiu2HAm3TYCAZe5yv5mFBYJYwEpQvu7ygztyDgF8NGsNCB33i3f",
	},
	{
		Address: "UBCTYUj6YLA7hPYdRTGJcaAQudKyUEq98go2",
		PeerId:  "16Uiu2HAmHicRtXhQkQTs3rBV4QBAGKAmb6A78nACV3xR2ogyMw6m",
	},
	{
		Address: "UBCRcDoWaF2mfiWgNiybKqt5zuYXork6Rpwn",
		PeerId:  "16Uiu2HAmHpdStkoNqawJpR9DP47CAFQe78W9N1qsxbJSUHPpmjHv",
	},
}
