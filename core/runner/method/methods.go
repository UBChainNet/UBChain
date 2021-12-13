package method

type Value struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type MethodInfo struct {
	Name    string  `json:"name"`
	Params  []Value `json:"params"`
	Returns []Value `json:"returns"`
}

var ExMethods = map[string]*MethodInfo{
	"Methods": &MethodInfo{
		Name:   "Methods",
		Params: nil,
		Returns: []Value{
			{
				Name: "Open methods",
				Type: "json",
			},
		},
	},
	"MethodExist": &MethodInfo{
		Name: "MethodExist",
		Params: []Value{
			{
				Name: "method",
				Type: "string",
			},
		},
		Returns: []Value{
			{
				Name: "exist",
				Type: "bool",
			},
		},
	},
	"Pairs": &MethodInfo{
		Name:   "Pairs",
		Params: nil,
		Returns: []Value{
			{
				Name: "pair list",
				Type: "json",
			},
		},
	},
	"ExchangeRouter": &MethodInfo{
		Name: "ExchangeRouter",
		Params: []Value{
			{
				Name: "tokenA",
				Type: "string",
			},
			{
				Name: "tokenB",
				Type: "string",
			},
		},
		Returns: []Value{
			{
				Name: "paths",
				Type: "json",
			},
		},
	},

	"ExchangeRouterWithAmount": &MethodInfo{
		Name: "ExchangeRouterWithAmount",
		Params: []Value{
			{
				Name: "tokenA",
				Type: "string",
			},
			{
				Name: "tokenB",
				Type: "string",
			},
			{
				Name: "amountIn",
				Type: "float64",
			},
		},
		Returns: []Value{
			{
				Name: "path and amount",
				Type: "json",
			},
		},
	},

	"ExchangeOptimalRouter": &MethodInfo{
		Name: "ExchangeOptimalRouter",
		Params: []Value{
			{
				Name: "tokenA",
				Type: "string",
			},
			{
				Name: "tokenB",
				Type: "string",
			},
			{
				Name: "amountIn",
				Type: "float64",
			},
		},
		Returns: []Value{
			{
				Name: "optimal path",
				Type: "json",
			},
		},
	},

	"AmountOut": &MethodInfo{
		Name: "AmountOut",
		Params: []Value{
			{
				Name: "paths(tokenA,tokenB,tokenC)",
				Type: "string",
			},
			{
				Name: "amountIn",
				Type: "float64",
			},
		},
		Returns: []Value{
			{
				Name: "amountOut",
				Type: "float64",
			},
		},
	},
	"AmountIn": &MethodInfo{
		Name: "AmountIn",
		Params: []Value{
			{
				Name: "paths(tokenA,tokenB,tokenC)",
				Type: "string",
			},
			{
				Name: "amountOut",
				Type: "float64",
			},
		},
		Returns: []Value{
			{
				Name: "amountIn",
				Type: "float64",
			},
		},
	},
	"LegalPair": &MethodInfo{
		Name: "LegalPair",
		Params: []Value{
			{
				Name: "tokenA",
				Type: "string",
			},
			{
				Name: "tokenB",
				Type: "string",
			},
		},
		Returns: []Value{
			{
				Name: "is legal",
				Type: "bool",
			},
		},
	},
	"PairAddress": &MethodInfo{
		Name: "PairAddress",
		Params: []Value{
			{
				Name: "tokenA",
				Type: "string",
			},
			{
				Name: "tokenB",
				Type: "string",
			},
		},
		Returns: []Value{
			{
				Name: "pair address",
				Type: "string",
			},
		},
	},
}

var PairMethods = map[string]*MethodInfo{
	"Methods": &MethodInfo{
		Name:   "Methods",
		Params: nil,
		Returns: []Value{
			{
				Name: "Open methods",
				Type: "json",
			},
		},
	},
	"MethodExist": &MethodInfo{
		Name: "MethodExist",
		Params: []Value{
			{
				Name: "method",
				Type: "string",
			},
		},
		Returns: []Value{
			{
				Name: "exist",
				Type: "bool",
			},
		},
	},
	"QuoteAmountB": &MethodInfo{
		Name: "QuoteAmountB",
		Params: []Value{
			{
				Name: "TokenA",
				Type: "string",
			},
			{
				Name: "AmountA",
				Type: "float64",
			},
		},
		Returns: []Value{
			{
				Name: "AmountB",
				Type: "float64",
			},
		},
	},

	"TotalValue": &MethodInfo{
		Name: "TotalValue",
		Params: []Value{
			{
				Name: "liquidity",
				Type: "float64",
			},
		},
		Returns: []Value{
			{
				Name: "value",
				Type: "json",
			},
		},
	},
	"Profit": &MethodInfo{
		Name: "Profit",
		Params: []Value{
			{
				Name: "liquidity",
				Type: "float64",
			},
		},
		Returns: []Value{
			{
				Name: "value",
				Type: "json",
			},
		},
	},
}

var PledgeMethods = map[string]*MethodInfo{
	"Methods": &MethodInfo{
		Name:   "Methods",
		Params: nil,
		Returns: []Value{
			{
				Name: "Open methods",
				Type: "json",
			},
		},
	},
	"MethodExist": &MethodInfo{
		Name: "MethodExist",
		Params: []Value{
			{
				Name: "method",
				Type: "string",
			},
		},
		Returns: []Value{
			{
				Name: "exist",
				Type: "bool",
			},
		},
	},
	"GetAccountRewards": &MethodInfo{
		Name:   "GetAccountRewards",
		Params: []Value{},
		Returns: []Value{
			{
				Name: "value",
				Type: "json",
			},
		},
	},
	"GetPledgeReward": &MethodInfo{
		Name: "GetPledgeReward",
		Params: []Value{
			{
				Name: "address",
				Type: "string",
			},
			{
				Name: "pair contract address",
				Type: "string",
			},
		},
		Returns: []Value{
			{
				Name: "value",
				Type: "json",
			},
		},
	},
	"GetPledgeRewards": &MethodInfo{
		Name: "GetPledgeRewards",
		Params: []Value{
			{
				Name: "address",
				Type: "string",
			},
		},
		Returns: []Value{
			{
				Name: "value",
				Type: "json",
			},
		},
	},
	"GetPledge": &MethodInfo{
		Name: "GetPledge",
		Params: []Value{
			{
				Name: "address",
				Type: "string",
			},
			{
				Name: "pair contract address",
				Type: "string",
			},
		},
		Returns: []Value{
			{
				Name: "value",
				Type: "json",
			},
		},
	},
	"GetPledges": &MethodInfo{
		Name: "GetPledges",
		Params: []Value{
			{
				Name: "address",
				Type: "string",
			},
		},
		Returns: []Value{
			{
				Name: "value",
				Type: "json",
			},
		},
	},
	"GetPairPool": &MethodInfo{
		Name:   "GetPairPool",
		Params: []Value{},
		Returns: []Value{
			{
				Name: "value",
				Type: "json",
			},
		},
	},
	"GetPledgeYields": &MethodInfo{
		Name:   "GetPledgeYields",
		Params: []Value{},
		Returns: []Value{
			{
				Name: "value",
				Type: "json",
			},
		},
	},
	"GetPoolInfos": &MethodInfo{
		Name:   "GetPoolInfos",
		Params: []Value{},
		Returns: []Value{
			{
				Name: "value",
				Type: "json",
			},
		},
	},
	"GetPoolInfo": &MethodInfo{
		Name: "GetPoolInfo",
		Params: []Value{
			{
				Name: "pair contract address",
				Type: "string",
			},
		},
		Returns: []Value{
			{
				Name: "value",
				Type: "json",
			},
		},
	},
}

var TokenHubMethods = map[string]*MethodInfo{
	"Methods": &MethodInfo{
		Name:   "Methods",
		Params: nil,
		Returns: []Value{
			{
				Name: "Open methods",
				Type: "json",
			},
		},
	},
	"MethodExist": &MethodInfo{
		Name: "MethodExist",
		Params: []Value{
			{
				Name: "method",
				Type: "string",
			},
		},
		Returns: []Value{
			{
				Name: "exist",
				Type: "bool",
			},
		},
	},
}
