package exchange_runner

var exMethods = map[string]*MethodInfo{
	"Methods": &MethodInfo{
		Name:   "Methods",
		Params: nil,
		Returns: []Value{
			{
				Name: "Open methods",
				Type: "",
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
				Type: "",
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
				Type: "",
			},
		},
	},
	"AmountOut": &MethodInfo{
		Name: "AmountOut",
		Params: []Value{
			{
				Name: "paths(tokenA->tokenB->tokenC)",
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
				Name: "paths(tokenA->tokenB->tokenC)",
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
}

var pairMethods = map[string]*MethodInfo{
	"Methods": &MethodInfo{
		Name:   "Methods",
		Params: nil,
		Returns: []Value{
			{
				Name: "Open methods",
				Type: "",
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
}
